#!/usr/bin/env python3
"""Local academic pack preparation for NieboGuard Compute.

This runs on a trusted machine with enough disk. It downloads one archive,
verifies SHA-256, extracts it safely, groups media by source session, assigns
train/validation/test without splitting a session, writes a manifest, and can
upload shards to R2 through the existing admin API.

It does NOT run on the Worker, does NOT upload raw full archives, and does NOT
execute code from a dataset or archive.
"""
from __future__ import annotations

import argparse
import hashlib
import json
import os
import re
import shutil
import stat
import sys
import tarfile
import time
import urllib.error
import urllib.parse
import urllib.request
import zipfile
from dataclasses import dataclass, field
from pathlib import Path
from typing import Any

ROOT = Path(__file__).resolve().parents[2]
AUDIO_EXT = {".wav"}
IMAGE_EXT = {".jpg", ".jpeg", ".png", ".webp"}
TIERS = {
    "MINI": {"min_mb": 25, "max_mb": 50, "target_mb": 40},
    "STANDARD": {"min_mb": 100, "max_mb": 200, "target_mb": 150},
    "EXTENDED": {"min_mb": 250, "max_mb": 500, "target_mb": 375},
    "GPU": {"min_mb": 500, "max_mb": 1024, "target_mb": 750},
}
DEFAULT_MAX_EXTRACT_BYTES = 60 * 1024 ** 3
DEFAULT_MAX_FILES = 500_000
FORBIDDEN_SUFFIX = {".zip", ".tar", ".gz", ".7z", ".rar", ".exe", ".dll", ".so", ".bat", ".ps1", ".sh", ".py"}


class PackError(RuntimeError):
    pass


def log(msg: str) -> None:
    print(msg, flush=True)


def sha256_file(path: Path, limit: int | None = None) -> str:
    h = hashlib.sha256()
    total = 0
    with path.open("rb") as fh:
        for chunk in iter(lambda: fh.read(1024 * 1024), b""):
            total += len(chunk)
            if limit is not None and total > limit:
                raise PackError(f"file exceeds limit: {path}")
            h.update(chunk)
    return h.hexdigest()


def disk_free_gb(path: Path) -> float:
    path.mkdir(parents=True, exist_ok=True)
    st = os.statvfs(path)
    return st.f_bavail * st.f_frsize / 1024 ** 3


def safe_name(value: str, max_len: int = 160) -> str:
    out = re.sub(r"[^A-Za-z0-9_.-]+", "_", value).strip("._-")
    return (out or "item")[:max_len]


def is_within(base: Path, target: Path) -> bool:
    try:
        target.resolve().relative_to(base.resolve())
        return True
    except ValueError:
        return False


def download(source: str, work: Path, expected_sha: str | None, resume: bool = True) -> Path:
    parsed = urllib.parse.urlparse(source)
    if parsed.scheme != "https":
        raise PackError("download source must be HTTPS")
    name = safe_name(Path(parsed.path).name or "archive.bin", 120)
    target = work / name
    if target.exists() and resume:
        actual = sha256_file(target)
        if expected_sha and actual.lower() != expected_sha.lower():
            log(f"existing download has wrong SHA-256, removing {target.name}")
            target.unlink()
        else:
            log(f"using existing download {target} ({target.stat().st_size/1e9:.2f} GB)")
            return target
    tmp = target.with_name(target.name + ".part")
    log(f"downloading {source} -> {target}")
    req = urllib.request.Request(source, headers={"User-Agent": "NieboGuard-Pack-Prep/0.1"})
    total = 0
    started = time.time()
    try:
        with urllib.request.urlopen(req, timeout=120) as resp, tmp.open("wb") as fh:
            content_length = int(resp.headers.get("Content-Length") or 0)
            if content_length and content_length / 1024 ** 3 > disk_free_gb(work):
                raise PackError("download does not fit in free disk space")
            while True:
                chunk = resp.read(1024 * 1024)
                if not chunk:
                    break
                total += len(chunk)
                fh.write(chunk)
                if total % (256 * 1024 * 1024) < 1024 * 1024:
                    speed = total / max(1, time.time() - started) / 1e6
                    log(f"  {total/1e9:.2f} GB, {speed:.1f} MB/s")
    except urllib.error.URLError as exc:
        raise PackError(f"download failed: {exc}") from exc
    os.replace(tmp, target)
    actual = sha256_file(target)
    if expected_sha and actual.lower() != expected_sha.lower():
        raise PackError(f"SHA-256 mismatch: expected {expected_sha}, got {actual}")
    log(f"download complete: {target} ({total/1e9:.2f} GB) SHA-256 {actual}")
    return target


def _zip_members(zf: zipfile.ZipFile, max_bytes: int, max_files: int):
    total = 0
    files = 0
    for info in zf.infolist():
        name = info.filename
        if not name or name.endswith("/"):
            continue
        files += 1
        if files > max_files:
            raise PackError("archive contains too many files")
        if info.file_size > max_bytes:
            raise PackError(f"archive member is too large: {name}")
        total += info.file_size
        if total > max_bytes:
            raise PackError("archive uncompressed size exceeds limit")
        yield info


def safe_extract(archive: Path, dest: Path, max_bytes: int = DEFAULT_MAX_EXTRACT_BYTES,
                 max_files: int = DEFAULT_MAX_FILES) -> None:
    dest.mkdir(parents=True, exist_ok=True)
    if zipfile.is_zipfile(archive):
        with zipfile.ZipFile(archive) as zf:
            for info in _zip_members(zf, max_bytes, max_files):
                name = info.filename.replace("\\", "/")
                parts = Path(name).parts
                if name.startswith("/") or ".." in parts or any(p in ("__MACOSX", ".git") for p in parts):
                    raise PackError(f"unsafe path in zip: {name}")
                target = dest / name
                if not is_within(dest, target):
                    raise PackError(f"path traversal in zip: {name}")
                if info.external_attr and stat.S_ISLNK(info.external_attr >> 16):
                    raise PackError(f"symlink in zip: {name}")
                target.parent.mkdir(parents=True, exist_ok=True)
                with zf.open(info) as src, target.open("wb") as dst:
                    copied = 0
                    while True:
                        chunk = src.read(1024 * 1024)
                        if not chunk:
                            break
                        copied += len(chunk)
                        if copied > max_bytes:
                            raise PackError("extracted file exceeds limit")
                        dst.write(chunk)
        return
    if tarfile.is_tarfile(archive):
        with tarfile.open(archive, "r:*") as tf:
            total = 0
            files = 0
            for member in tf:
                if member.isdir():
                    continue
                if not member.isfile():
                    raise PackError(f"unsupported tar member type: {member.name}")
                files += 1
                if files > max_files:
                    raise PackError("archive contains too many files")
                name = member.name.replace("\\", "/")
                parts = Path(name).parts
                if name.startswith("/") or ".." in parts or any(p in ("__MACOSX", ".git") for p in parts):
                    raise PackError(f"unsafe path in tar: {name}")
                target = dest / name
                if not is_within(dest, target):
                    raise PackError(f"path traversal in tar: {name}")
                total += member.size
                if total > max_bytes:
                    raise PackError("archive uncompressed size exceeds limit")
                target.parent.mkdir(parents=True, exist_ok=True)
                src = tf.extractfile(member)
                if src is None:
                    raise PackError(f"cannot read tar member: {name}")
                with src, target.open("wb") as dst:
                    shutil.copyfileobj(src, dst, length=1024 * 1024)
        return
    raise PackError("source is not a supported ZIP or TAR archive; pass --source-dir for extracted data")


def load_label_map(path: str | None, labels: list[str]) -> dict[str, str]:
    mapping: dict[str, str] = {}
    if path:
        data = json.loads(Path(path).read_text(encoding="utf-8"))
        if not isinstance(data, dict) or not all(isinstance(k, str) and isinstance(v, str) for k, v in data.items()):
            raise PackError("label map must be a JSON object of string->string")
        mapping = {k.lower(): v for k, v in data.items()}
    for label in labels:
        mapping.setdefault(label.lower(), label)
    return mapping


def label_for(rel: Path, mapping: dict[str, str], allow_unknown: bool) -> str:
    lower_parts = [p.lower() for p in rel.parts]
    for key, label in mapping.items():
        if key in lower_parts or any(key in part for part in lower_parts):
            return label
    if allow_unknown:
        return "unknown"
    raise PackError(f"no label mapping for {rel}; use --label-map or --allow-unknown")


def split_for_group(dataset_id: str, version: int, group: str) -> str:
    h = 2166136261
    for ch in f"{dataset_id}:{version}:{group}":
        h = (h ^ ord(ch)) * 16777619 & 0xFFFFFFFF
    r = h % 100
    return "train" if r < 70 else "validation" if r < 85 else "test"


def split_from_path(rel: Path) -> str | None:
    parts = [p.lower() for p in rel.parts]
    for part in parts:
        if part in ("train", "training"):
            return "train"
        if part in ("val", "valid", "validation"):
            return "validation"
        if part in ("test", "testing"):
            return "test"
    return None


def collect_media(source_dir: Path, modality: str, max_file_mb: int, exclude: re.Pattern | None) -> list[Path]:
    ext = AUDIO_EXT if modality == "audio" else IMAGE_EXT
    files = []
    for path in sorted(source_dir.rglob("*")):
        if not path.is_file() or path.suffix.lower() not in ext:
            continue
        rel = path.relative_to(source_dir)
        if exclude and exclude.search(rel.as_posix()):
            log(f"excluding by regex: {rel}")
            continue
        if path.stat().st_size > max_file_mb * 1024 * 1024:
            raise PackError(f"file exceeds --max-file-mb: {rel}")
        if path.suffix.lower() in FORBIDDEN_SUFFIX:
            continue
        files.append(path)
    if not files:
        raise PackError(f"no {modality} files found in {source_dir}")
    return files


def group_of(rel: Path, depth: int) -> str:
    parts = rel.parts
    if depth <= 0:
        return safe_name(rel.parent.name, 80)
    return safe_name("__".join(parts[:depth]), 100)


def plan_tier(groups: list[dict[str, Any]], cfg: dict[str, float]) -> dict[str, Any]:
    max_bytes = int(cfg["max_mb"] * 1_000_000)
    min_bytes = int(cfg["min_mb"] * 1_000_000)
    oversized = [g for g in groups if g["bytes"] > max_bytes]
    if oversized:
        return {"ok": False, "reason": f"group too large: {oversized[0]['name']}", "packs": []}
    bins: list[dict[str, Any]] = []
    for group in sorted(groups, key=lambda g: g["name"]):
        candidate = None
        for b in bins:
            if b["bytes"] + group["bytes"] <= max_bytes and (candidate is None or b["bytes"] < candidate["bytes"]):
                candidate = b
        if candidate:
            candidate["bytes"] += group["bytes"]
            candidate["groups"].append(group)
        else:
            bins.append({"bytes": group["bytes"], "groups": [group]})
    small = [b for b in bins if b["bytes"] < min_bytes]
    return {
        "ok": all(b["bytes"] >= min_bytes for b in bins),
        "packs": [{"bytes": b["bytes"], "groups": [g["name"] for g in b["groups"]]} for b in bins],
        "small_packs": len(small),
        "reason": "some packs are below min_mb; merge/choose a different tier" if small else None,
    }


def prepare(args) -> dict[str, Any]:
    work = Path(args.work).resolve()
    out = Path(args.out).resolve()
    work.mkdir(parents=True, exist_ok=True)
    out.mkdir(parents=True, exist_ok=True)
    need = args.min_free_gb
    if disk_free_gb(work) < need:
        raise PackError(f"not enough free space in {work}: need at least {need} GB")
    log(f"work={work} out={out} free={disk_free_gb(work):.1f} GB")

    if args.source_dir:
        source_dir = Path(args.source_dir).resolve()
        if not source_dir.is_dir():
            raise PackError("--source-dir is not a directory")
    else:
        if not args.source:
            raise PackError("provide --source URL/archive or --source-dir")
        source = Path(args.source)
        if source.exists():
            archive = source.resolve()
        else:
            archive = download(args.source, work, args.sha256, resume=not args.no_resume)
        if args.sha256:
            actual = sha256_file(archive)
            if actual.lower() != args.sha256.lower():
                raise PackError(f"SHA-256 mismatch: expected {args.sha256}, got {actual}")
        extract_dir = work / "extracted"
        if args.reuse_extracted and extract_dir.exists():
            log(f"reusing extracted data in {extract_dir}")
        else:
            if extract_dir.exists():
                shutil.rmtree(extract_dir)
            log("extracting safely...")
            safe_extract(archive, extract_dir, max_bytes=int(args.max_extract_gb * 1024 ** 3), max_files=args.max_files)
        source_dir = extract_dir

    mapping = load_label_map(args.label_map, args.labels)
    exclude = re.compile(args.exclude_regex, re.I) if args.exclude_regex else None
    media = collect_media(source_dir, args.modality, args.max_file_mb, exclude)
    groups: dict[str, dict[str, Any]] = {}
    flat_used: set[str] = set()
    shards: list[dict[str, Any]] = []
    total_bytes = 0
    total_samples = 0
    for path in media:
        rel = path.relative_to(source_dir)
        group = group_of(rel, args.session_depth)
        label = label_for(rel, mapping, args.allow_unknown)
        split = split_from_path(rel) or split_for_group(args.dataset_id, args.version, group)
        flat = safe_name(f"{group}__{path.name}", 170)
        if flat in flat_used:
            flat = safe_name(f"{group}__{sha256_file(path)[:12]}__{path.name}", 170)
        flat_used.add(flat)
        digest = sha256_file(path)
        size = path.stat().st_size
        total_bytes += size
        total_samples += 1
        groups.setdefault(group, {"name": group, "bytes": 0, "files": [], "split": split, "labels": {}})
        if groups[group]["split"] != split:
            raise PackError(f"session {group} has conflicting splits")
        groups[group]["bytes"] += size
        groups[group]["files"].append({"source": str(path), "name": flat, "bytes": size, "sha256": digest, "label": label})
        groups[group]["labels"][label] = groups[group]["labels"].get(label, 0) + 1
        shards.append({
            "name": flat, "source_path": str(path), "bytes": size, "sha256": digest,
            "split": split, "group_key": group, "sample_count": 0,
        })
    # One JSONL per session keeps labels with their media in the same pack.
    for group, g in sorted(groups.items()):
        jsonl_name = safe_name(f"{group}__labels.jsonl", 170)
        labels_path = out / "labels" / jsonl_name
        labels_path.parent.mkdir(parents=True, exist_ok=True)
        lines = []
        for item in g["files"]:
            lines.append(json.dumps({"file": item["name"], "label": item["label"], "split": g["split"], "group": group}, ensure_ascii=False))
        labels_path.write_text("\n".join(lines) + "\n", encoding="utf-8")
        labels_sha = sha256_file(labels_path)
        shards.append({
            "name": jsonl_name, "source_path": str(labels_path), "bytes": labels_path.stat().st_size,
            "sha256": labels_sha, "split": g["split"], "group_key": group, "sample_count": len(g["files"]),
            "class_counts": g["labels"],
        })
    plans = {}
    group_list = [{"name": g["name"], "bytes": g["bytes"]} for g in groups.values()]
    for tier in args.pack_tiers:
        if tier not in TIERS:
            raise PackError(f"unknown tier {tier}")
        cfg = dict(TIERS[tier])
        if args.pack_config:
            overrides = json.loads(Path(args.pack_config).read_text(encoding="utf-8"))
            if tier in overrides:
                cfg.update(overrides[tier])
        plans[tier] = plan_tier(group_list, cfg)
    dataset = {
        "schema": "nieboguard-academic-pack-plan/v1",
        "dataset_id": args.dataset_id,
        "version": args.version,
        "name": args.name or args.dataset_id,
        "modality": args.modality,
        "classes": sorted({c for shard in shards for c in (shard.get("class_counts") or {})}) or args.labels,
        "license": args.license,
        "license_url": args.license_url,
        "authors": args.authors,
        "provenance": args.provenance,
        "source_url": args.source_url,
        "redistribution_allowed": not args.no_redistribution,
        "training_allowed": True,
        "commercial_allowed": args.commercial,
        "sample_count": total_samples,
        "total_bytes": total_bytes,
        "groups": len(groups),
        "shards": shards,
        "pack_plan": plans,
        "notes": "Local trusted preparation. Raw archive is not uploaded. Only shards listed here go to R2.",
    }
    (out / "dataset.json").write_text(json.dumps(dataset, ensure_ascii=False, indent=2), encoding="utf-8")
    (out / "manifest.json").write_text(json.dumps({
        "dataset_id": args.dataset_id, "version": args.version,
        "sample_count": total_samples, "total_bytes": total_bytes,
        "shards": len(shards), "groups": len(groups), "pack_plan": plans,
        "created_at": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime()),
    }, ensure_ascii=False, indent=2), encoding="utf-8")
    log(f"prepared {len(media)} media files, {len(groups)} sessions, {len(shards)} shards, {total_bytes/1e9:.3f} GB")
    for tier, plan in plans.items():
        log(f"  {tier}: {len(plan.get('packs', []))} packs, ok={plan.get('ok')} {plan.get('reason') or ''}")
    impossible = [tier for tier, plan in plans.items() if not plan.get("ok")]
    if impossible and not args.allow_impossible_packs:
        raise PackError(f"pack plan is not valid for: {', '.join(impossible)}; choose other tiers, adjust pack-config, or pass --allow-impossible-packs")
    if args.upload:
        upload_dataset(args, dataset)
    return dataset


def _request(method: str, url: str, token: str | None = None, body: Any = None, raw: bytes | None = None,
             headers: dict[str, str] | None = None, timeout: int = 120) -> tuple[int, Any]:
    h = dict(headers or {})
    data = raw
    if body is not None:
        data = json.dumps(body, ensure_ascii=False).encode()
        h["Content-Type"] = "application/json"
    if token:
        h["Authorization"] = "Bearer " + token
    req = urllib.request.Request(url, data=data, headers=h, method=method)
    try:
        with urllib.request.urlopen(req, timeout=timeout) as resp:
            payload = resp.read()
            try:
                return resp.status, json.loads(payload) if payload else {}
            except Exception:
                return resp.status, payload
    except urllib.error.HTTPError as exc:
        payload = exc.read()
        try:
            data = json.loads(payload)
        except Exception:
            data = payload
        return exc.code, data


def upload_dataset(args, dataset: dict[str, Any]) -> None:
    if not args.api or not args.admin_token:
        raise PackError("--upload requires --api and --admin-token")
    api = args.api.rstrip("/")
    state_path = Path(args.out).resolve() / "upload-state.json"
    state = json.loads(state_path.read_text()) if state_path.exists() else {}
    uploaded: dict[str, str] = state.get("uploads", {})
    log(f"uploading {len(dataset['shards'])} shards to {api}")
    shards_for_api = []
    for shard in dataset["shards"]:
        source = Path(shard["source_path"])
        key = shard["sha256"]
        if key in uploaded:
            upload_id = uploaded[key]
        else:
            if source.stat().st_size > 120 * 1024 * 1024:
                raise PackError(f"single shard exceeds 120 MB admin upload limit: {source.name}")
            ext = source.suffix.lower()
            ctype = "audio/wav" if ext == ".wav" else "application/x-ndjson" if ext == ".jsonl" else "application/octet-stream"
            code, created = _request("POST", api + "/api/v2/admin/compute/uploads", args.admin_token, body={
                "purpose": "dataset", "filename": shard["name"], "bytes": shard["bytes"],
                "content_type": ctype, "sha256": shard["sha256"],
            })
            if code != 201:
                raise PackError(f"upload slot failed: {code} {created}")
            url = created["url"]
            log(f"  uploading {shard['name']} ({shard['bytes']/1e6:.2f} MB)")
            code, result = _request("PUT", url, raw=source.read_bytes(), headers={
                "Content-Type": ctype, "Content-Length": str(shard["bytes"]),
                "X-Content-SHA256": shard["sha256"],
            }, timeout=300)
            if code != 200:
                raise PackError(f"upload failed: {code} {result}")
            upload_id = created["upload_id"]
            uploaded[key] = upload_id
            state["uploads"] = uploaded
            state_path.write_text(json.dumps(state, ensure_ascii=False, indent=2), encoding="utf-8")
        shards_for_api.append({
            "upload_id": upload_id, "name": shard["name"], "split": shard["split"],
            "group_key": shard["group_key"], "sample_count": shard.get("sample_count", 0),
            "class_counts": shard.get("class_counts", {}),
        })
    code, created = _request("POST", api + "/api/v2/admin/compute/datasets", args.admin_token, body={
        "id": dataset["dataset_id"], "version": dataset["version"], "name": dataset["name"],
        "modality": dataset["modality"], "classes": dataset["classes"],
        "license": dataset["license"], "license_url": dataset["license_url"],
        "authors": dataset["authors"], "provenance": dataset["provenance"],
        "source_url": dataset["source_url"],
        "redistribution_allowed": dataset["redistribution_allowed"],
        "training_allowed": dataset["training_allowed"], "commercial_allowed": dataset["commercial_allowed"],
        "sample_count": dataset["sample_count"], "shards": shards_for_api,
        "notes": dataset["notes"],
    })
    if code != 201:
        raise PackError(f"dataset create failed: {code} {created}")
    log(f"dataset created: {dataset['dataset_id']}@{dataset['version']} sha={created.get('sha256')}")
    if args.build_packs:
        code, packs = _request("POST", api + f"/api/v2/admin/compute/datasets/{urllib.parse.quote(dataset['dataset_id'])}/{dataset['version']}/build-packs", args.admin_token, body={
            "tiers": args.pack_tiers, "model_bytes": args.model_bytes, "library_bytes": args.library_bytes,
        }, timeout=600)
        if code not in (200, 201):
            raise PackError(f"pack build failed: {code} {packs}")
        log(f"packs: {len(packs.get('created', []))} created, {len(packs.get('skipped', []))} skipped")
        (Path(args.out).resolve() / "packs-result.json").write_text(json.dumps(packs, ensure_ascii=False, indent=2), encoding="utf-8")


def main(argv: list[str] | None = None) -> int:
    p = argparse.ArgumentParser(description="Prepare small academic data packs locally and optionally upload shards to R2.")
    src = p.add_mutually_exclusive_group(required=True)
    src.add_argument("--source", help="HTTPS URL or local ZIP/TAR archive")
    src.add_argument("--source-dir", help="already extracted data directory")
    p.add_argument("--sha256", help="expected SHA-256 of downloaded archive")
    p.add_argument("--work", default="/data/nieboguard-compute/work", help="large work directory for download/extract")
    p.add_argument("--out", default="/data/nieboguard-compute/staging", help="staging directory")
    p.add_argument("--dataset-id", required=True)
    p.add_argument("--version", type=int, default=1)
    p.add_argument("--name")
    p.add_argument("--modality", choices=["audio"], default="audio")
    p.add_argument("--license", required=True)
    p.add_argument("--license-url", default=None)
    p.add_argument("--authors", nargs="*", default=[])
    p.add_argument("--provenance", default="")
    p.add_argument("--source-url", default=None)
    p.add_argument("--commercial", action="store_true")
    p.add_argument("--no-redistribution", action="store_true")
    p.add_argument("--label-map", help="JSON object folder/name -> class")
    p.add_argument("--labels", nargs="*", default=[])
    p.add_argument("--allow-unknown", action="store_true")
    p.add_argument("--allow-impossible-packs", action="store_true")
    p.add_argument("--exclude-regex", help="exclude paths matching this regex, e.g. speech|mowa")
    p.add_argument("--session-depth", type=int, default=1, help="number of path components that form a session")
    p.add_argument("--max-file-mb", type=int, default=120)
    p.add_argument("--max-extract-gb", type=float, default=60)
    p.add_argument("--max-files", type=int, default=DEFAULT_MAX_FILES)
    p.add_argument("--min-free-gb", type=float, default=40)
    p.add_argument("--pack-tiers", default="MINI,STANDARD")
    p.add_argument("--pack-config", help="optional JSON overrides for tier sizes")
    p.add_argument("--model-bytes", type=int, default=0)
    p.add_argument("--library-bytes", type=int, default=0)
    p.add_argument("--reuse-extracted", action="store_true")
    p.add_argument("--no-resume", action="store_true")
    p.add_argument("--upload", action="store_true")
    p.add_argument("--build-packs", action="store_true")
    p.add_argument("--api", default=os.environ.get("NIEBOGUARD_API", "https://api.neboguard.pl"))
    p.add_argument("--admin-token", default=os.environ.get("ADMIN_TOKEN", ""))
    args = p.parse_args(argv)
    args.pack_tiers = [x.strip().upper() for x in args.pack_tiers.split(",") if x.strip()]
    try:
        prepare(args)
        return 0
    except (PackError, OSError, ValueError, json.JSONDecodeError) as exc:
        print(f"BŁĄD: {exc}", file=sys.stderr)
        return 1


if __name__ == "__main__":
    raise SystemExit(main())

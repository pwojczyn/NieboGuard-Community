#!/usr/bin/env python3
"""Audio AI data preparation for NieboGuard Compute.

Pipeline for the first real Audio AI stage:
  archive/dir -> safe extract -> mono 16 kHz WAV -> segments -> quality filter
  -> dedup by SHA-256 -> session grouping -> pack manifest -> optional R2 upload.

It uses ffmpeg only for decoding/normalisation. No archive code is executed.
"""
from __future__ import annotations

import argparse
import array
import hashlib
import json
import math
import os
import re
import shutil
import subprocess
import sys
import tempfile
import time
import wave
from pathlib import Path
from typing import Any

HERE = Path(__file__).resolve().parent
sys.path.insert(0, str(HERE))
import prepare_academic_pack as pack  # noqa: E402

AUDIO_EXT = {".wav", ".flac", ".mp3", ".ogg", ".opus", ".m4a", ".aac"}
SR = 16000
try:
    import numpy as _np
except Exception:
    _np = None


def require_ffmpeg() -> None:
    if not shutil.which("ffmpeg"):
        raise pack.PackError("Brak ffmpeg. Zainstaluj ffmpeg albo dostarcz już znormalizowane pliki WAV 16 kHz mono.")


def probe_audio(path: Path) -> dict[str, Any]:
    ffprobe = shutil.which("ffprobe")
    if not ffprobe:
        return {}
    try:
        out = subprocess.run([ffprobe, "-v", "error", "-show_entries", "stream=sample_rate,channels,codec_name,duration",
                              "-of", "json", str(path)], capture_output=True, text=True, timeout=60, check=False)
        data = json.loads(out.stdout or "{}")
        streams = data.get("streams") or [{}]
        return streams[0] if streams else {}
    except Exception:
        return {}


def normalize(src: Path, dst: Path) -> None:
    dst.parent.mkdir(parents=True, exist_ok=True)
    cmd = ["ffmpeg", "-v", "error", "-y", "-i", str(src), "-ac", "1", "-ar", str(SR), "-sample_fmt", "s16", str(dst)]
    try:
        res = subprocess.run(cmd, capture_output=True, text=True, timeout=1800, check=False)
    except subprocess.TimeoutExpired as exc:
        raise pack.PackError(f"ffmpeg timeout for {src}") from exc
    if res.returncode != 0 or not dst.exists():
        raise pack.PackError(f"ffmpeg failed for {src}: {res.stderr[-300:]}")


def write_wav(path: Path, samples: array.array) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    with wave.open(str(path), "wb") as w:
        w.setnchannels(1)
        w.setsampwidth(2)
        w.setframerate(SR)
        w.writeframes(samples.tobytes())


def quality(samples: array.array, min_rms: float, max_clipping: float) -> tuple[bool, float, float]:
    if not samples:
        return False, 0.0, 1.0
    peak = 32767 * 0.999
    if _np is not None:
        arr = _np.frombuffer(samples.tobytes(), dtype=_np.int16).astype(_np.float64)
        rms = float(_np.sqrt(_np.mean(arr * arr)) / 32768.0)
        clip_ratio = float(_np.mean(_np.abs(arr) >= peak))
        return (rms >= min_rms and clip_ratio <= max_clipping), rms, clip_ratio
    clipped = 0
    total = 0.0
    for value in samples:
        if abs(value) >= peak:
            clipped += 1
        total += value * value
    rms = math.sqrt(total / len(samples)) / 32768.0
    clip_ratio = clipped / len(samples)
    return (rms >= min_rms and clip_ratio <= max_clipping), rms, clip_ratio


def segments_of(samples: array.array, seconds: float, overlap: float, min_seconds: float):
    total = len(samples)
    length = max(1, int(seconds * SR))
    step = max(1, int((seconds - overlap) * SR))
    min_len = max(1, int(min_seconds * SR))
    start = 0
    while start < total:
        chunk = samples[start:start + length]
        if len(chunk) >= min_len:
            yield start, chunk
        if start + length >= total:
            break
        start += step


def audio_prepare(args) -> dict[str, Any]:
    require_ffmpeg()
    work = Path(args.work).resolve()
    out = Path(args.out).resolve()
    work.mkdir(parents=True, exist_ok=True)
    out.mkdir(parents=True, exist_ok=True)
    if pack.disk_free_gb(work) < args.min_free_gb:
        raise pack.PackError(f"not enough free space in {work}: need {args.min_free_gb} GB")

    if args.source_dir:
        source_dir = Path(args.source_dir).resolve()
    else:
        source = Path(args.source)
        if source.exists():
            archive = source
        else:
            archive = pack.download(args.source, work, args.sha256, resume=not args.no_resume)
        if args.sha256:
            actual = pack.sha256_file(archive)
            if actual.lower() != args.sha256.lower():
                raise pack.PackError(f"SHA-256 mismatch: expected {args.sha256}, got {actual}")
        extract_dir = work / "extracted"
        if args.reuse_extracted and extract_dir.exists():
            source_dir = extract_dir
        else:
            if extract_dir.exists():
                shutil.rmtree(extract_dir)
            pack.safe_extract(archive, extract_dir, max_bytes=int(args.max_extract_gb * 1024 ** 3), max_files=args.max_files)
            source_dir = extract_dir

    mapping = pack.load_label_map(args.label_map, args.labels)
    class_labels = sorted(set(mapping.values()) | set(args.labels))
    exclude = re.compile(args.exclude_regex, re.I) if args.exclude_regex else None
    files = []
    for path in sorted(source_dir.rglob("*")):
        if not path.is_file() or path.suffix.lower() not in AUDIO_EXT:
            continue
        rel = path.relative_to(source_dir)
        if exclude and exclude.search(rel.as_posix()):
            print(f"excluding: {rel}", flush=True)
            continue
        if path.stat().st_size > args.max_file_mb * 1024 * 1024:
            raise pack.PackError(f"file exceeds --max-file-mb: {rel}")
        files.append(path)
    if not files:
        raise pack.PackError("no audio files found")

    converted = work / "audio_converted"
    if converted.exists():
        shutil.rmtree(converted)
    converted.mkdir(parents=True)
    seen: set[str] = set()
    total_segments = 0
    kept = 0
    skipped_quality = 0
    skipped_dup = 0
    started = time.time()
    for idx, path in enumerate(files, 1):
        rel = path.relative_to(source_dir)
        try:
            label = pack.label_for(rel, mapping, args.allow_unknown)
        except pack.PackError:
            raise pack.PackError(f"no label for {rel}; provide --label-map or --allow-unknown")
        session = pack.group_of(rel, args.session_depth)
        tmp = Path(tempfile.mkdtemp(prefix="ng-audio-", dir=str(work)))
        normalized = tmp / "normalized.wav"
        try:
            normalize(path, normalized)
            with wave.open(str(normalized), "rb") as w:
                if w.getnchannels() != 1 or w.getsampwidth() != 2 or w.getframerate() != SR:
                    raise pack.PackError(f"normalisation failed for {rel}")
                data = array.array("h")
                data.frombytes(w.readframes(w.getnframes()))
            for seg_idx, (start, chunk) in enumerate(segments_of(data, args.segment_seconds, args.overlap_seconds, args.min_segment_seconds)):
                ok, rms, clip = quality(chunk, args.min_rms, args.max_clipping)
                if not ok:
                    skipped_quality += 1
                    continue
                target = converted / label / session / f"{pack.safe_name(path.stem, 70)}__{seg_idx:04d}.wav"
                write_wav(target, chunk)
                digest = pack.sha256_file(target)
                if digest in seen:
                    target.unlink(missing_ok=True)
                    skipped_dup += 1
                    continue
                seen.add(digest)
                kept += 1
                total_segments += 1
                if kept % 100 == 0:
                    print(f"segments: {kept}, files: {idx}/{len(files)}, elapsed {time.time()-started:.0f}s", flush=True)
        finally:
            shutil.rmtree(tmp, ignore_errors=True)

    if not kept:
        raise pack.PackError("no segments passed quality filters")
    print(f"audio prepared: {len(files)} files -> {kept} segments; quality-skipped={skipped_quality}, duplicates={skipped_dup}")

    # Reuse the generic manifest/upload pipeline. Directory layout is label/session/file.wav.
    generic = argparse.Namespace(
        source=None, source_dir=str(converted), sha256=None, work=str(work / "pack"), out=str(out),
        dataset_id=args.dataset_id, version=args.version, name=args.name or args.dataset_id,
        modality="audio", license=args.license, license_url=args.license_url,
        authors=args.authors, provenance=args.provenance, source_url=args.source_url,
        commercial=args.commercial, no_redistribution=args.no_redistribution,
        label_map=None, labels=class_labels, allow_unknown=args.allow_unknown,
        allow_impossible_packs=args.allow_impossible_packs, exclude_regex=None,
        session_depth=2, max_file_mb=args.max_file_mb, max_extract_gb=args.max_extract_gb,
        max_files=args.max_files, min_free_gb=0, pack_tiers=args.pack_tiers,
        pack_config=args.pack_config, model_bytes=args.model_bytes, library_bytes=args.library_bytes,
        reuse_extracted=True, no_resume=True, upload=args.upload, build_packs=args.build_packs,
        api=args.api, admin_token=args.admin_token,
    )
    dataset = pack.prepare(generic)
    dataset["audio"] = {
        "source_files": len(files), "segments": kept, "quality_skipped": skipped_quality,
        "duplicates": skipped_dup, "segment_seconds": args.segment_seconds,
        "overlap_seconds": args.overlap_seconds, "sample_rate": SR,
    }
    (out / "audio-summary.json").write_text(json.dumps(dataset["audio"], ensure_ascii=False, indent=2), encoding="utf-8")
    return dataset


def main(argv: list[str] | None = None) -> int:
    p = argparse.ArgumentParser(description="Prepare Audio AI segments and packs locally.")
    src = p.add_mutually_exclusive_group(required=True)
    src.add_argument("--source", help="HTTPS URL or local archive")
    src.add_argument("--source-dir", help="already extracted audio directory")
    p.add_argument("--sha256")
    p.add_argument("--work", default="/data/nieboguard-compute/audio-work")
    p.add_argument("--out", default="/data/nieboguard-compute/staging/audio")
    p.add_argument("--dataset-id", required=True)
    p.add_argument("--version", type=int, default=1)
    p.add_argument("--name")
    p.add_argument("--license", required=True)
    p.add_argument("--license-url")
    p.add_argument("--authors", nargs="*", default=[])
    p.add_argument("--provenance", default="")
    p.add_argument("--source-url")
    p.add_argument("--commercial", action="store_true")
    p.add_argument("--no-redistribution", action="store_true")
    p.add_argument("--label-map")
    p.add_argument("--labels", nargs="*", default=[])
    p.add_argument("--allow-unknown", action="store_true")
    p.add_argument("--allow-impossible-packs", action="store_true")
    p.add_argument("--exclude-regex")
    p.add_argument("--session-depth", type=int, default=1)
    p.add_argument("--segment-seconds", type=float, default=5.0)
    p.add_argument("--overlap-seconds", type=float, default=0.5)
    p.add_argument("--min-segment-seconds", type=float, default=1.0)
    p.add_argument("--min-rms", type=float, default=0.005)
    p.add_argument("--max-clipping", type=float, default=0.01)
    p.add_argument("--max-file-mb", type=int, default=512)
    p.add_argument("--max-extract-gb", type=float, default=60)
    p.add_argument("--max-files", type=int, default=pack.DEFAULT_MAX_FILES)
    p.add_argument("--min-free-gb", type=float, default=60)
    p.add_argument("--pack-tiers", default="STANDARD,EXTENDED")
    p.add_argument("--pack-config")
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
        audio_prepare(args)
        return 0
    except (pack.PackError, OSError, ValueError, subprocess.SubprocessError) as exc:
        print(f"BŁĄD: {exc}", file=sys.stderr)
        return 1


if __name__ == "__main__":
    raise SystemExit(main())

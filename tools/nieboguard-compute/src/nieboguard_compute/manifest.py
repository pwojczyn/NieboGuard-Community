"""Manifest validation before any download or training."""
from __future__ import annotations

import re
import time
from datetime import datetime, timezone
from typing import Any
from urllib.parse import urlparse

from .crypto import SignatureError, canonical_json, sha256_hex, verify_manifest
from .api import same_origin

TASK_SCHEMA = "nieboguard-compute-task/v1"
SUPPORTED_TASK_TYPES = {"audio-classification", "vision-detection-head"}
SAFETENSORS_MIME = "application/vnd.nieboguard.safetensors"


class ManifestError(ValueError):
    pass


def _int(value: Any, name: str, lo: int, hi: int) -> int:
    try:
        out = int(value)
    except Exception as exc:
        raise ManifestError(f"{name} must be an integer") from exc
    if out < lo or out > hi:
        raise ManifestError(f"{name} outside range")
    return out


def validate_manifest(manifest: dict[str, Any], public_key: str | None, api_base: str,
                      allow_unsigned: bool = False, now: float | None = None,
                      max_data_mb: float | None = None, allowed_tiers: list[str] | None = None) -> dict[str, Any]:
    if not isinstance(manifest, dict):
        raise ManifestError("manifest must be an object")
    if manifest.get("schema") != TASK_SCHEMA:
        raise ManifestError("unsupported manifest schema")
    if not public_key:
        if not allow_unsigned:
            raise ManifestError("brak przypiętego klucza publicznego; odmawiam uruchomienia niepodpisanego manifestu")
    else:
        try:
            verify_manifest(manifest, public_key)
        except SignatureError as exc:
            raise ManifestError(str(exc)) from exc
    current = now if now is not None else time.time()
    try:
        raw_expiry = manifest.get("expires_at")
        if isinstance(raw_expiry, (int, float)):
            expires = float(raw_expiry)
        else:
            text = str(raw_expiry or "").replace("Z", "+00:00")
            dt = datetime.fromisoformat(text)
            if dt.tzinfo is None:
                dt = dt.replace(tzinfo=timezone.utc)
            expires = dt.timestamp()
    except Exception as exc:
        raise ManifestError("manifest has no valid expiry") from exc
    if expires < current:
        raise ManifestError("manifest expired")
    if manifest.get("task_type") not in SUPPORTED_TASK_TYPES:
        raise ManifestError("unsupported task type")
    dataset = manifest.get("dataset")
    if not isinstance(dataset, dict):
        raise ManifestError("manifest has no dataset")
    if dataset.get("redistribution_allowed") is not True:
        raise ManifestError("dataset has no redistribution permission")
    if dataset.get("training_allowed") is not True:
        raise ManifestError("dataset has no training permission")
    if not str(dataset.get("license") or "").strip():
        raise ManifestError("dataset license is missing")
    shards = dataset.get("shards")
    if not isinstance(shards, list) or not shards:
        raise ManifestError("dataset has no shards")
    pack = manifest.get("pack")
    if pack is not None:
        if not isinstance(pack, dict) or not str(pack.get("id") or "").startswith("cdp_"):
            raise ManifestError("invalid data pack")
        if pack.get("tier") not in ("MINI", "STANDARD", "EXTENDED", "GPU"):
            raise ManifestError("invalid data pack tier")
        if allowed_tiers and pack.get("tier") not in allowed_tiers:
            raise ManifestError("data pack tier is not allowed by config")
        pack_bytes = _int(pack.get("bytes"), "pack bytes", 1, 2048 * 1024 * 1024)
        if max_data_mb is not None and pack_bytes > max_data_mb * 1_000_000:
            raise ManifestError("data pack exceeds configured data limit")
        if not re.fullmatch(r"[0-9a-f]{64}", str(pack.get("sha256") or "")):
            raise ManifestError("data pack has no valid SHA-256")
    total = 0
    for shard in shards:
        if not isinstance(shard, dict):
            raise ManifestError("invalid shard")
        name = str(shard.get("name") or "")
        if not name or "/" in name or "\\" in name or name.startswith("."):
            raise ManifestError("shard name must be a safe basename")
        url = str(shard.get("url") or "")
        if not same_origin(url, api_base):
            raise ManifestError("shard URL is not on the configured API origin")
        size = _int(shard.get("bytes"), f"shard {name} bytes", 1, 512 * 1024 * 1024)
        total += size
        if total > 512 * 1024 * 1024:
            raise ManifestError("dataset exceeds client size limit")
        checksum = str(shard.get("sha256") or "").lower()
        if len(checksum) != 64 or any(c not in "0123456789abcdef" for c in checksum):
            raise ManifestError(f"shard {name} has no valid SHA-256")
    upload = manifest.get("upload")
    if not isinstance(upload, dict):
        raise ManifestError("manifest has no upload slot")
    if upload.get("content_type") != SAFETENSORS_MIME:
        raise ManifestError("upload content type must be safetensors")
    if not same_origin(str(upload.get("url") or ""), api_base):
        raise ManifestError("upload URL is not on the configured API origin")
    _int(upload.get("max_bytes"), "upload max_bytes", 1, 16 * 1024 * 1024)
    architecture = manifest.get("architecture")
    if not isinstance(architecture, dict):
        raise ManifestError("manifest has no architecture")
    specs = architecture.get("tensor_specs")
    if not isinstance(specs, list) or not specs:
        raise ManifestError("manifest has no tensor specs")
    for spec in specs:
        if not isinstance(spec, dict) or not spec.get("name") or not isinstance(spec.get("shape"), list):
            raise ManifestError("invalid tensor spec")
        if spec.get("dtype", "F32") != "F32":
            raise ManifestError("client supports only F32 task tensors")
    classes = manifest.get("classes")
    if not isinstance(classes, list) or not classes:
        raise ManifestError("manifest has no classes")
    if "base_model" not in manifest or not manifest["base_model"].get("version"):
        raise ManifestError("manifest has no exact base model")
    max_seconds = _int(manifest.get("max_train_seconds"), "max_train_seconds", 10, 86400)
    hyper = manifest.get("hyperparams") or {}
    for name, lo, hi in (("epochs", 1, 100), ("batch_size", 1, 512), ("max_steps", 10, 100000)):
        _int(hyper.get(name), f"hyperparams.{name}", lo, hi)
    return {
        "task_id": str(manifest.get("task_id")),
        "round_id": str(manifest.get("round_id")),
        "task_type": manifest.get("task_type"),
        "modality": manifest.get("modality"),
        "max_train_seconds": max_seconds,
        "shard_count": len(shards),
        "total_bytes": total,
        "architecture_hash": str((manifest.get("architecture") or {}).get("hash") or ""),
        "base_model": manifest.get("base_model"),
        "pack": pack,
        "pack_sha256": (pack or {}).get("sha256"),
    }


def safe_work_name(value: str, fallback: str = "item") -> str:
    out = "".join(ch if ch.isalnum() or ch in "._-" else "_" for ch in str(value or ""))
    out = out.strip("._-")[:120]
    return out or fallback

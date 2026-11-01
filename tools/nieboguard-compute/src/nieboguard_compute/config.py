"""Konfiguracja klienta i bezpieczne przechowywanie tokenów."""
from __future__ import annotations

import json
import os
import stat
import sys
import tempfile
from pathlib import Path
from typing import Any
from urllib.parse import urlparse

CONFIG_VERSION = 1
CONSENT_VERSION = "2026-11-01"
DEFAULT_API = "https://api.neboguard.pl"

DEFAULTS: dict[str, Any] = {
    "config_version": CONFIG_VERSION,
    "api_base": DEFAULT_API,
    "manifest_public_key": "",
    "manifest_key_id": "ng-compute-1",
    "device_id": None,
    "refresh_token": None,
    "consent_version": None,
    "consent_at": None,
    "display_name": None,
    "modalities": ["audio"],
    "max_cpu_percent": 50,
    "max_gpu_percent": 60,
    "schedule": "idle",
    "idle_minutes": 10,
    "max_disk_gb": 5,
    "max_data_mb": 200,
    "pack_tiers": ["MINI", "STANDARD"],
    "max_task_seconds": None,
    "data_dir": "",
    "leaderboard_opt_in": False,
}


class ConfigError(ValueError):
    pass


def default_data_dir() -> Path:
    if os.name == "nt":
        base = os.environ.get("LOCALAPPDATA") or str(Path.home() / "AppData" / "Local")
        return Path(base) / "NieboGuardCompute"
    base = os.environ.get("XDG_DATA_HOME") or str(Path.home() / ".local" / "share")
    return Path(base) / "nieboguard-compute"


def config_path(data_dir: str | None = None) -> Path:
    root = Path(data_dir) if data_dir else default_data_dir()
    return root / "config.json"


def ensure_data_dir(data_dir: str | None = None) -> Path:
    root = Path(data_dir) if data_dir else default_data_dir()
    root.mkdir(parents=True, exist_ok=True)
    try:
        os.chmod(root, stat.S_IRWXU)
    except OSError:
        pass
    return root


def _allowed_api(url: str) -> bool:
    parsed = urlparse(str(url))
    if parsed.scheme == "https":
        return True
    if parsed.scheme == "http" and parsed.hostname in ("localhost", "127.0.0.1", "::1"):
        return True
    return False


def load(path: str | None = None, data_dir: str | None = None) -> dict[str, Any]:
    target = Path(path) if path else config_path(data_dir)
    if not target.exists():
        return dict(DEFAULTS, data_dir=str(ensure_data_dir(data_dir)))
    try:
        raw = json.loads(target.read_text(encoding="utf-8"))
    except Exception as exc:
        raise ConfigError(f"Nie można odczytać konfiguracji: {exc}") from exc
    if not isinstance(raw, dict):
        raise ConfigError("Konfiguracja musi być obiektem JSON.")
    cfg = dict(DEFAULTS)
    cfg.update(raw)
    cfg["data_dir"] = str(raw.get("data_dir") or target.parent)
    validate(cfg)
    return cfg


def save(cfg: dict[str, Any], path: str | None = None) -> Path:
    validate(cfg)
    target = Path(path) if path else config_path(cfg.get("data_dir"))
    target.parent.mkdir(parents=True, exist_ok=True)
    payload = json.dumps(cfg, ensure_ascii=False, indent=2, sort_keys=True)
    fd, tmp_name = tempfile.mkstemp(prefix=".config-", suffix=".json", dir=str(target.parent))
    try:
        with os.fdopen(fd, "w", encoding="utf-8") as fh:
            fh.write(payload)
            fh.flush()
            os.fsync(fh.fileno())
        try:
            os.chmod(tmp_name, stat.S_IRUSR | stat.S_IWUSR)
        except OSError:
            pass
        os.replace(tmp_name, target)
    finally:
        if os.path.exists(tmp_name):
            try:
                os.unlink(tmp_name)
            except OSError:
                pass
    return target


def validate(cfg: dict[str, Any]) -> None:
    api = str(cfg.get("api_base") or "")
    if not _allowed_api(api):
        raise ConfigError("api_base musi być HTTPS (HTTP tylko dla localhost).")
    key = str(cfg.get("manifest_public_key") or "")
    if key and len(key) < 20:
        raise ConfigError("manifest_public_key jest za krótki.")
    mods = cfg.get("modalities")
    if not isinstance(mods, list) or not mods or any(m not in ("audio", "vision") for m in mods):
        raise ConfigError("modalities musi zawierać audio i/lub vision.")
    for key_name, lo, hi in (("max_cpu_percent", 5, 100), ("max_gpu_percent", 0, 100), ("idle_minutes", 1, 240), ("max_disk_gb", 1, 100), ("max_data_mb", 0, 2048)):
        value = cfg.get(key_name)
        if not isinstance(value, (int, float)) or not (lo <= float(value) <= hi):
            raise ConfigError(f"{key_name} musi być w zakresie {lo}–{hi}.")
    if cfg.get("schedule") not in ("always", "idle", "manual"):
        raise ConfigError("schedule musi być always, idle lub manual.")
    tiers = cfg.get("pack_tiers")
    if not isinstance(tiers, list) or any(t not in ("MINI", "STANDARD", "EXTENDED", "GPU") for t in tiers):
        raise ConfigError("pack_tiers musi zawierać MINI/STANDARD/EXTENDED/GPU.")
    if bool(cfg.get("device_id")) != bool(cfg.get("refresh_token")):
        raise ConfigError("device_id i refresh_token muszą występować razem.")


def redacted(cfg: dict[str, Any]) -> dict[str, Any]:
    out = dict(cfg)
    if out.get("refresh_token"):
        out["refresh_token"] = "***"
    return out


def data_root(cfg: dict[str, Any]) -> Path:
    root = Path(cfg.get("data_dir") or default_data_dir())
    root.mkdir(parents=True, exist_ok=True)
    return root

"""Training runner: download, verify, train, upload, report."""
from __future__ import annotations

import hashlib
import os
import time
from pathlib import Path
from typing import Any

from . import __version__
from .api import ApiError, Client, retry
from .audio import load_audio_dataset
from .config import data_root
from .crypto import sha256_hex
from .hardware import detect, idle_seconds
from .manifest import ManifestError, safe_work_name, validate_manifest
from .model import TrainConfig, TrainingError, export_artifact, train_audio, train_vision
from .torch_backend import TORCH_AVAILABLE, cuda_available, train_audio_torch, train_vision_torch
from .vision import load_vision_dataset

SAFETENSORS_MIME = "application/vnd.nieboguard.safetensors"


class RunnerError(RuntimeError):
    pass


class Runner:
    def __init__(self, cfg: dict[str, Any], allow_unsigned: bool = False, verbose: bool = True):
        self.cfg = cfg
        self.allow_unsigned = allow_unsigned
        self.verbose = verbose
        self.client = Client(str(cfg["api_base"]))
        self.root = data_root(cfg)
        self.stop_file = self.root / "STOP"

    def log(self, message: str) -> None:
        if self.verbose:
            print(message, flush=True)

    def authenticate(self) -> None:
        if not self.cfg.get("device_id") or not self.cfg.get("refresh_token"):
            raise RunnerError("Brak rejestracji. Uruchom: nieboguard-compute enroll --invite <kod> --consent")
        try:
            self.client.refresh(str(self.cfg["device_id"]), str(self.cfg["refresh_token"]))
            self.client.update_preferences({
                "modalities": self.cfg.get("modalities") or ["audio"],
                "max_cpu_percent": int(self.cfg.get("max_cpu_percent") or 50),
                "max_gpu_percent": int(self.cfg.get("max_gpu_percent") or 60),
                "schedule": self.cfg.get("schedule") or "manual",
                "max_data_mb": float(self.cfg.get("max_data_mb") or 0),
                "pack_tiers": self.cfg.get("pack_tiers") or [],
            })
        except ApiError as exc:
            if exc.code in ("device_revoked", "unauthorized"):
                raise RunnerError("Urządzenie zostało odłączone lub token jest nieprawidłowy. Zarejestruj ponownie.") from exc
            raise

    def should_run(self, wait: bool = True) -> bool:
        schedule = str(self.cfg.get("schedule") or "idle")
        if schedule in ("always", "manual"):
            return True
        if schedule == "idle":
            while True:
                if self.stop_file.exists():
                    return False
                idle = idle_seconds()
                if idle is None:
                    raise RunnerError(
                        "Tryb idle nie może odczytać bezczynności na tym systemie. "
                        "Użyj --schedule manual (uruchamiaj ręcznie) albo --schedule always."
                    )
                need = float(self.cfg.get("idle_minutes") or 10) * 60
                if idle >= need:
                    return True
                if not wait:
                    return False
                self.log(f"Komputer jest używany (bezczynność {int(idle)} s). Czekam do {int(need)} s.")
                time.sleep(min(60.0, max(5.0, need - idle)))
        return False

    def stop_requested(self) -> bool:
        return self.stop_file.exists()

    def run(self, once: bool = False, max_tasks: int = 0) -> int:
        if not self.client.access_token:
            self.authenticate()
        completed = 0
        failures = 0
        while True:
            if self.stop_requested():
                self.log("Zatrzymano na żądanie użytkownika.")
                return 0
            if not self.should_run(wait=not once):
                self.log("Komputer nie jest bezczynny; nie uruchamiam zadania.")
                return 0
            try:
                response = retry(lambda: self._call_with_refresh("next_task"))
            except ApiError as exc:
                self.log(f"Nie można pobrać zadania: {exc.message}")
                if once:
                    return 2
                time.sleep(30)
                continue
            manifest = response.get("manifest")
            if not manifest:
                self.log(str(response.get("note") or "Brak zadania. Czekam."))
                if once:
                    return 0
                time.sleep(max(15, min(300, 30 + failures * 10)))
                continue
            try:
                outcome = self.process_task(manifest)
                if outcome in ("completed", "rejected"):
                    completed += 1
                    failures = 0
                else:
                    failures += 1
            except (ManifestError, RunnerError, TrainingError, ApiError, ValueError, OSError, KeyError) as exc:
                failures += 1
                self.log(f"Zadanie nieukończone: {exc}")
                task_id = str(manifest.get("task_id") or "")
                if task_id:
                    try:
                        self._call_with_refresh("fail", task_id, str(exc)[:300])
                    except Exception:
                        pass
                if once:
                    return 1
            if once or (max_tasks and completed >= max_tasks):
                return 0 if completed else 1
            if self.stop_requested():
                return 0
            time.sleep(2)

    def _call_with_refresh(self, method: str, *args, **kwargs):
        try:
            return getattr(self.client, method)(*args, **kwargs)
        except ApiError as exc:
            if exc.status == 401:
                self.authenticate()
                return getattr(self.client, method)(*args, **kwargs)
            raise

    def process_task(self, manifest: dict[str, Any]) -> str:
        info = validate_manifest(
            manifest,
            str(self.cfg.get("manifest_public_key") or "") or None,
            str(self.cfg["api_base"]),
            allow_unsigned=self.allow_unsigned,
            max_data_mb=float(self.cfg.get("max_data_mb") or 0) or None,
            allowed_tiers=list(self.cfg.get("pack_tiers") or []),
        )
        if not self.cfg.get("manifest_public_key") and self.allow_unsigned:
            self.log("UWAGA: uruchomiono niepodpisany manifest w trybie deweloperskim.")
        task_id = safe_work_name(info["task_id"], "task")
        work = self.root / "tasks" / task_id
        data_dir = work / "data"
        data_dir.mkdir(parents=True, exist_ok=True)
        max_disk = float(self.cfg.get("max_disk_gb") or 5) * 1024 ** 3
        used = sum(p.stat().st_size for p in work.rglob("*") if p.is_file()) if work.exists() else 0
        remaining = max(0, int(max_disk - used))
        total = int(info["total_bytes"])
        if total > remaining:
            raise RunnerError(f"Zadanie wymaga {total} B, a pozostały limit dysku to {remaining} B.")
        pack = info.get("pack")
        pack_text = f", paczka {pack['tier']} {pack['id']} ({pack['bytes']} B)" if pack else ""
        self.log(f"Zadanie {info['task_id']}: {info['task_type']}, {info['shard_count']} plików, {total} B{pack_text}.")
        split_by_file = {}
        for shard in manifest["dataset"]["shards"]:
            name = safe_work_name(str(shard["name"]), "shard")
            split_by_file[name] = str(shard.get("split") or "train")
            target = data_dir / name
            checksum = str(shard["sha256"]).lower()
            if target.exists() and hashlib.sha256(target.read_bytes()).hexdigest() == checksum:
                continue
            self.log(f"Pobieram {name} ({shard.get('bytes')} B)…")
            self.client.download(str(shard["url"]), target, int(shard.get("bytes") or 0), checksum)
        base_artifact = None
        base_model = manifest.get("base_model") or {}
        if base_model.get("url"):
            base_path = work / "base.safetensors"
            if not base_path.exists() or sha256_hex(base_path.read_bytes()) != str(base_model.get("sha256")):
                self.log("Pobieram model bazowy…")
                self.client.download(str(base_model["url"]), base_path, 16 * 1024 * 1024, str(base_model.get("sha256") or ""))
            base_artifact = base_path
        architecture = manifest["architecture"]
        specs = architecture["tensor_specs"]
        backbone = next((s for s in specs if s.get("name") == "backbone.weight"), None)
        input_dim = int(backbone["shape"][0]) if backbone and backbone.get("shape") else 0
        classes = [str(x) for x in manifest["classes"]]
        if info["task_type"] == "audio-classification":
            data = load_audio_dataset(data_dir, input_dim, classes, split_by_file)
            trainer = train_audio
            torch_trainer = train_audio_torch
        elif info["task_type"] == "vision-detection-head":
            data = load_vision_dataset(data_dir, input_dim, classes, split_by_file)
            trainer = train_vision
            torch_trainer = train_vision_torch
        else:
            raise RunnerError("Nieobsługiwany typ zadania.")
        requested_backend = os.environ.get("NIEBOGUARD_COMPUTE_BACKEND", "auto").lower()
        gpu_percent = int(self.cfg.get("max_gpu_percent") or 0)
        use_torch = TORCH_AVAILABLE and requested_backend in ("auto", "torch", "cuda") and (
            requested_backend in ("torch", "cuda") or (cuda_available() and gpu_percent > 0)
        )
        if use_torch:
            trainer = torch_trainer
        hyper = manifest.get("hyperparams") or {}
        max_seconds = int(info["max_train_seconds"])
        if self.cfg.get("max_task_seconds"):
            max_seconds = min(max_seconds, int(self.cfg["max_task_seconds"]))
        cfg = TrainConfig(
            task_id=info["task_id"], round_id=info["round_id"],
            base_model=str(manifest["base_model"]["id"]), base_version=str(manifest["base_model"]["version"]),
            architecture_hash=info["architecture_hash"], classes=classes,
            seed=int(hyper.get("seed", 1337)), epochs=int(hyper.get("epochs", 8)),
            batch_size=int(hyper.get("batch_size", 32)), learning_rate=float(hyper.get("learning_rate", 0.01)),
            weight_decay=float(hyper.get("weight_decay", 0.0001)), max_steps=int(hyper.get("max_steps", 2000)),
            max_seconds=max_seconds, cpu_percent=int(self.cfg.get("max_cpu_percent", 50)),
            checkpoint=work / "checkpoint.safetensors", resume=True,
            backend="torch" if use_torch else "python",
            device="cuda" if use_torch and cuda_available() and int(self.cfg.get("max_gpu_percent") or 0) > 0 else "cpu",
            max_gpu_percent=int(self.cfg.get("max_gpu_percent") or 0),
            should_stop=self.stop_requested,
            base_artifact=base_artifact,
        )
        self.log("Trenuję lokalnie. Możesz zatrzymać proces; checkpoint pozwoli wznowić.")
        result = trainer(data, specs, cfg)
        if result.get("stopped") or self.stop_requested():
            self._call_with_refresh("fail", info["task_id"], "stopped by user before upload")
            self.log("Zatrzymano po checkpointcie; wynik nie został wysłany.")
            return "stopped"
        artifact_path = work / "artifact.safetensors"
        blob, _metadata = export_artifact(result, cfg, artifact_path, len(data))
        sha = sha256_hex(blob)
        upload = manifest["upload"]
        if len(blob) > int(upload["max_bytes"]):
            raise RunnerError("Wyeksportowany artefakt przekracza limit zadania.")
        self.log(f"Wgrywam wynik ({len(blob)} B, SHA-256 {sha[:12]}…)…")
        self.client.upload(str(upload["url"]), blob, SAFETENSORS_MIME, sha)
        metrics = dict(result["metrics"])
        metrics["samples"] = len(data)
        metrics["client_version"] = __version__
        metrics["backend"] = cfg.backend
        response = self._call_with_refresh("complete", info["task_id"], {
            "upload_id": upload["upload_id"], "sha256": sha, "bytes": len(blob),
            "metrics": metrics,
            "checkpoint": {"epochs_completed": metrics.get("epochs_completed", 0), "resumed": metrics.get("resumed", False)},
        })
        status = str(response.get("quality_status") or response.get("status") or "")
        if response.get("ok"):
            self.log(f"Wynik przyjęty do kwarantanny: {status}. Czeka na niezależną ewaluację.")
            return "completed"
        self.log(f"Wynik odrzucony: {response.get('reasons') or response.get('error')}")
        return "rejected"


def run_once(cfg: dict[str, Any], allow_unsigned: bool = False, verbose: bool = True) -> int:
    return Runner(cfg, allow_unsigned=allow_unsigned, verbose=verbose).run(once=True)


def hardware_for_enrollment() -> dict[str, Any]:
    hw = detect()
    return hw

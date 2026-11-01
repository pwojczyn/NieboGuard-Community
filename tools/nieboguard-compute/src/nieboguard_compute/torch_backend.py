"""Opcjonalny backend PyTorch (CPU/CUDA) dla tych samych głów.

Instalator nie pobiera PyTorch domyślnie. Jeżeli użytkownik ma już PyTorch z
obsługą CUDA, runner może użyć GPU. W przeciwnym razie działa czysty Python.
Wagi są inicjalizowane tym samym deterministycznym ziarnem co backend Python,
więc serwer może porównać je z modelem bazowym.
"""
from __future__ import annotations

import math
import random
import time
from typing import Any

from .model import (AUDIO_NAMES, VISION_NAMES, TrainConfig, TrainingError, _check_specs,
                    _class_counts, _spec_map, _split, _vision_label_index,
                    deterministic_init, evaluate_audio, evaluate_vision, save_checkpoint, tensor_stats)
from .safetensors import Tensor

try:
    import torch
    from torch import nn
except Exception:  # pragma: no cover - optional dependency
    torch = None
    nn = None

TORCH_AVAILABLE = torch is not None


def cuda_available() -> bool:
    return bool(TORCH_AVAILABLE and torch.cuda.is_available())


def describe() -> dict[str, Any]:
    if not TORCH_AVAILABLE:
        return {"available": False, "cuda": False}
    return {
        "available": True,
        "version": str(torch.__version__),
        "cuda": bool(torch.cuda.is_available()),
        "device": torch.cuda.get_device_name(0) if torch.cuda.is_available() else "cpu",
    }


def _tensors_from_state(state: dict[str, Any], specs: list[dict[str, Any]]) -> dict[str, Tensor]:
    out = {}
    for spec in specs:
        name = spec["name"]
        value = state[name]
        out[name] = Tensor(name=name, shape=[int(x) for x in spec["shape"]],
                           values=[float(x) for x in value.detach().cpu().reshape(-1).tolist()])
    return out


def _load_init(net: Any, specs: list[dict[str, Any]], seed: int, checkpoint: Any, cfg: TrainConfig) -> int:
    if checkpoint is not None and cfg.checkpoint and cfg.checkpoint.exists() and cfg.resume:
        old, meta = __import__("nieboguard_compute.safetensors", fromlist=["read_safetensors"]).read_safetensors(cfg.checkpoint.read_bytes())
        if meta.get("nieboguard.task_id") == cfg.task_id and meta.get("nieboguard.architecture_hash") == cfg.architecture_hash:
            with torch.no_grad():
                for spec in specs:
                    name = spec["name"]
                    if name in old:
                        param = net.state_dict()[name]
                        param.copy_(torch.tensor(old[name].values, dtype=param.dtype, device=param.device).reshape(param.shape))
            return int(meta.get("epoch", 0))
    init = deterministic_init(specs, seed)
    with torch.no_grad():
        for spec in specs:
            name = spec["name"]
            param = net.state_dict()[name]
            param.copy_(torch.tensor(init[name].values, dtype=param.dtype, device=param.device).reshape(param.shape))
    return 0


if TORCH_AVAILABLE:
    class _AudioNet(nn.Module):
        def __init__(self, input_dim: int, hidden: int, classes: int):
            super().__init__()
            self.backbone = nn.Linear(input_dim, hidden)
            self.classifier = nn.Linear(hidden, classes)

        def forward(self, x):
            return self.classifier(torch.relu(self.backbone(x)))

    class _VisionNet(nn.Module):
        def __init__(self, input_dim: int, hidden: int, classes: int):
            super().__init__()
            self.backbone = nn.Linear(input_dim, hidden)
            self.classifier = nn.Linear(hidden, classes)
            self.objectness = nn.Linear(hidden, 1)
            self.box = nn.Linear(hidden, 4)

        def forward(self, x):
            h = torch.relu(self.backbone(x))
            return self.objectness(h), self.classifier(h), self.box(h)
else:
    _AudioNet = None
    _VisionNet = None


def _device(cfg: TrainConfig) -> Any:
    if torch.cuda.is_available() and cfg.device != "cpu":
        device = torch.device("cuda")
        if cfg.cpu_percent <= 0:
            raise TrainingError("GPU jest wyłączone przez limit")
        try:
            if cfg.max_seconds:
                torch.cuda.set_per_process_memory_fraction(0.8, device)
        except Exception:
            pass
        return device
    return torch.device("cpu")


def _val_split(data: list[dict[str, Any]], seed: int):
    return _split(data, seed)


def train_audio_torch(data: list[dict[str, Any]], specs: list[dict[str, Any]], cfg: TrainConfig) -> dict[str, Any]:
    if not TORCH_AVAILABLE:
        raise TrainingError("PyTorch nie jest zainstalowany")
    _check_specs(specs, AUDIO_NAMES)
    shapes = _spec_map(specs)
    hidden, input_dim, classes = shapes["backbone.bias"][0], shapes["backbone.weight"][0], len(cfg.classes)
    train, val = _val_split(data, cfg.seed)
    device = _device(cfg)
    net = _AudioNet(input_dim, hidden, classes).to(device)
    start_epoch = _load_init(net, specs, cfg.seed, None, cfg)
    opt = torch.optim.Adam(net.parameters(), lr=cfg.learning_rate, weight_decay=cfg.weight_decay)
    xs = torch.tensor([r["features"] for r in train], dtype=torch.float32, device=device)
    ys = torch.tensor([cfg.classes.index(r["label"]) for r in train], dtype=torch.long, device=device)
    loss_fn = nn.CrossEntropyLoss()
    start_wall = time.monotonic()
    steps = 0
    epochs_done = start_epoch
    history = []
    for epoch in range(start_epoch, cfg.epochs):
        if time.monotonic() - start_wall >= cfg.max_seconds or steps >= cfg.max_steps:
            break
        perm = torch.randperm(len(train), device=device)
        total = 0.0
        seen = 0
        for start in range(0, len(train), max(1, cfg.batch_size)):
            if time.monotonic() - start_wall >= cfg.max_seconds or steps >= cfg.max_steps:
                break
            idx = perm[start:start + max(1, cfg.batch_size)]
            logits = net(xs[idx])
            loss = loss_fn(logits, ys[idx])
            opt.zero_grad()
            loss.backward()
            opt.step()
            total += float(loss.detach().cpu())
            seen += len(idx)
            steps += 1
        state = _tensors_from_state(net.state_dict(), specs)
        acc, val_loss = evaluate_audio(state, val, cfg.classes)
        epochs_done = epoch + 1
        history.append({"epoch": epochs_done, "train_loss": round(total / max(1, seen), 6), "val_accuracy": round(acc, 6), "val_loss": round(val_loss, 6), "steps": steps})
        save_checkpoint(cfg, state, epochs_done, acc, len(data), "torch", str(device))
    state = _tensors_from_state(net.state_dict(), specs)
    acc, val_loss = evaluate_audio(state, val, cfg.classes)
    return {
        "tensors": state,
        "metrics": {"accuracy": acc, "loss": val_loss, "samples": len(data),
                    "class_counts": _class_counts(data, cfg.classes),
                    "duration_s": round(time.monotonic() - start_wall, 3), "epochs_completed": epochs_done,
                    "backend": "torch", "device": str(device), "resumed": start_epoch > 0},
        "history": history, "validation_size": len(val),
    }


def train_vision_torch(data: list[dict[str, Any]], specs: list[dict[str, Any]], cfg: TrainConfig) -> dict[str, Any]:
    if not TORCH_AVAILABLE:
        raise TrainingError("PyTorch nie jest zainstalowany")
    _check_specs(specs, VISION_NAMES)
    shapes = _spec_map(specs)
    hidden, input_dim, classes = shapes["backbone.bias"][0], shapes["backbone.weight"][0], len(cfg.classes)
    train, val = _val_split(data, cfg.seed)
    device = _device(cfg)
    net = _VisionNet(input_dim, hidden, classes).to(device)
    start_epoch = _load_init(net, specs, cfg.seed, None, cfg)
    opt = torch.optim.Adam(net.parameters(), lr=cfg.learning_rate, weight_decay=cfg.weight_decay)
    xs = torch.tensor([r["features"] for r in train], dtype=torch.float32, device=device)
    obj = torch.tensor([1.0 if r["objectness"] else 0.0 for r in train], dtype=torch.float32, device=device)
    labels = torch.tensor([max(0, _vision_label_index(r.get("label"), cfg.classes)) for r in train], dtype=torch.long, device=device)
    boxes = torch.tensor([r["box"] for r in train], dtype=torch.float32, device=device)
    mask = obj.bool()
    bce = nn.BCEWithLogitsLoss()
    ce = nn.CrossEntropyLoss()
    start_wall = time.monotonic()
    steps = 0
    epochs_done = start_epoch
    history = []
    for epoch in range(start_epoch, cfg.epochs):
        if time.monotonic() - start_wall >= cfg.max_seconds or steps >= cfg.max_steps:
            break
        perm = torch.randperm(len(train), device=device)
        total = 0.0
        seen = 0
        for start in range(0, len(train), max(1, cfg.batch_size)):
            if time.monotonic() - start_wall >= cfg.max_seconds or steps >= cfg.max_steps:
                break
            idx = perm[start:start + max(1, cfg.batch_size)]
            obj_logits, logits, box = net(xs[idx])
            loss = bce(obj_logits.squeeze(-1), obj[idx])
            positive = mask[idx]
            if positive.any():
                loss = loss + ce(logits[positive], labels[idx][positive])
                loss = loss + 0.5 * torch.mean((box[positive] - boxes[idx][positive]) ** 2)
            opt.zero_grad()
            loss.backward()
            opt.step()
            total += float(loss.detach().cpu())
            seen += len(idx)
            steps += 1
        state = _tensors_from_state(net.state_dict(), specs)
        metrics = evaluate_vision(state, val, cfg.classes)
        epochs_done = epoch + 1
        history.append({"epoch": epochs_done, "train_loss": round(total / max(1, seen), 6), **metrics, "steps": steps})
        save_checkpoint(cfg, state, epochs_done, metrics["objectness_f1"], len(data), "torch", str(device))
    state = _tensors_from_state(net.state_dict(), specs)
    metrics = evaluate_vision(state, val, cfg.classes)
    return {
        "tensors": state,
        "metrics": {"accuracy": metrics["class_accuracy"], "objectness_f1": metrics["objectness_f1"],
                    "loss": metrics["loss"], "samples": len(data), "class_counts": _class_counts(data, cfg.classes),
                    "duration_s": round(time.monotonic() - start_wall, 3), "epochs_completed": epochs_done,
                    "backend": "torch", "device": str(device), "resumed": start_epoch > 0},
        "history": history, "validation_size": len(val),
    }

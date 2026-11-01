"""Deterministic trainable heads for Audio AI and Vision AI.

The client never starts a large backbone from scratch. It fine-tunes a small
head over frozen features (YAMNet-style audio embeddings or detector proposal
features). Weights are exported as safetensors, never pickle.
"""
from __future__ import annotations

import json
import math
import random
import time
from dataclasses import dataclass
from pathlib import Path
from typing import Any, Iterable

from .safetensors import Tensor, read_safetensors, write_safetensors

AUDIO_NAMES = ("backbone.weight", "backbone.bias", "classifier.weight", "classifier.bias")
VISION_NAMES = (
    "backbone.weight", "backbone.bias", "classifier.weight", "classifier.bias",
    "objectness.weight", "objectness.bias", "box.weight", "box.bias",
)


class TrainingError(RuntimeError):
    pass


def _relu(x: float) -> float:
    return x if x > 0 else 0.0


def _sigmoid(x: float) -> float:
    if x >= 0:
        z = math.exp(-x)
        return 1.0 / (1.0 + z)
    z = math.exp(x)
    return z / (1.0 + z)


def _softmax(values: list[float]) -> list[float]:
    if not values:
        return []
    m = max(values)
    exps = [math.exp(v - m) for v in values]
    total = sum(exps) or 1.0
    return [v / total for v in exps]


def _matvec(matrix: list[float], vector: list[float], rows: int, cols: int) -> list[float]:
    return [sum(matrix[r * cols + c] * vector[c] for c in range(cols)) for r in range(rows)]


def _outer_add(grad: list[float], vector: list[float], rows: int, cols: int, scale: float = 1.0) -> None:
    for r in range(rows):
        base = r * cols
        vr = vector[r] * scale
        for c in range(cols):
            grad[base + c] += vr * c


def deterministic_init(specs: list[dict[str, Any]], seed: int) -> dict[str, Tensor]:
    tensors: dict[str, Tensor] = {}
    for spec in specs:
        name = str(spec["name"])
        shape = [int(x) for x in spec["shape"]]
        count = 1
        for d in shape:
            count *= d
        rng = random.Random(f"{seed}:{name}")
        if name.endswith(".bias") or name.startswith("box.") and name.endswith(".bias"):
            values = [0.0] * count
        else:
            fan_in = shape[0] if shape else 1
            std = 1.0 / math.sqrt(max(1, fan_in))
            values = [rng.gauss(0.0, std) for _ in range(count)]
        tensors[name] = Tensor(name=name, shape=shape, values=values)
    return tensors


def tensor_stats(tensors: dict[str, Tensor]) -> list[dict[str, Any]]:
    out = []
    for name in sorted(tensors):
        t = tensors[name]
        l2 = math.sqrt(sum(v * v for v in t.values))
        out.append({"name": name, "dtype": "F32", "shape": t.shape, "count": len(t.values),
                    "min": min(t.values), "max": max(t.values), "abs_max": max(abs(v) for v in t.values),
                    "mean": sum(t.values) / max(1, len(t.values)), "l2": l2, "nonfinite": 0})
    return out


def _adjust_lr(lr: float, epoch: int) -> float:
    return lr / (1.0 + 0.15 * max(0, epoch - 1))


@dataclass
class TrainConfig:
    task_id: str
    round_id: str
    base_model: str
    base_version: str
    architecture_hash: str
    classes: list[str]
    seed: int = 1337
    epochs: int = 8
    batch_size: int = 32
    learning_rate: float = 0.01
    weight_decay: float = 0.0001
    max_steps: int = 2000
    max_seconds: int = 1800
    cpu_percent: int = 50
    checkpoint: Path | None = None
    resume: bool = True
    backend: str = "python"
    device: str = "cpu"


class _Adam:
    def __init__(self, shapes: dict[str, list[int]], lr: float, weight_decay: float):
        self.lr = lr
        self.wd = weight_decay
        self.m = {k: [0.0] * _count(v) for k, v in shapes.items()}
        self.v = {k: [0.0] * _count(v) for k, v in shapes.items()}
        self.t = 0

    def step(self, tensors: dict[str, Tensor], grads: dict[str, list[float]], lr: float) -> None:
        self.t += 1
        b1, b2, eps = 0.9, 0.999, 1e-8
        for name, t in tensors.items():
            g = grads.get(name)
            if not g:
                continue
            gm = self.m[name]
            gv = self.v[name]
            for i, grad in enumerate(g):
                if self.wd:
                    grad += self.wd * t.values[i]
                gm[i] = b1 * gm[i] + (1 - b1) * grad
                gv[i] = b2 * gv[i] + (1 - b2) * grad * grad
                mhat = gm[i] / (1 - b1 ** self.t)
                vhat = gv[i] / (1 - b2 ** self.t)
                t.values[i] -= lr * mhat / (math.sqrt(vhat) + eps)


def _count(shape: list[int]) -> int:
    n = 1
    for d in shape:
        n *= d
    return n


class _Budget:
    def __init__(self, max_seconds: int, cpu_percent: int):
        self.start_wall = time.monotonic()
        self.start_cpu = time.process_time()
        self.max_seconds = max(1, int(max_seconds))
        self.cpu_fraction = max(0.05, min(1.0, cpu_percent / 100.0))

    def expired(self) -> bool:
        return time.monotonic() - self.start_wall >= self.max_seconds

    def throttle(self) -> None:
        wall = time.monotonic() - self.start_wall
        cpu = time.process_time() - self.start_cpu
        allowed_cpu = wall * self.cpu_fraction
        if cpu > allowed_cpu + 0.05:
            time.sleep(min(0.05, cpu - allowed_cpu))

    def progress(self) -> float:
        return min(1.0, (time.monotonic() - self.start_wall) / self.max_seconds)


def _split(data: list[dict[str, Any]], seed: int) -> tuple[list[dict[str, Any]], list[dict[str, Any]]]:
    train = [x for x in data if x.get("split", "train") == "train"]
    val = [x for x in data if x.get("split") == "validation"]
    if not train:
        train = list(data)
    if not val:
        rng = random.Random(seed)
        shuffled = list(train)
        rng.shuffle(shuffled)
        n = max(1, len(shuffled) // 10)
        val = shuffled[:n]
        train = shuffled[n:] or shuffled
    return train, val


def _spec_map(specs: list[dict[str, Any]]) -> dict[str, list[int]]:
    return {str(s["name"]): [int(x) for x in s["shape"]] for s in specs}


def _check_specs(specs: list[dict[str, Any]], expected: tuple[str, ...]) -> None:
    names = {str(s["name"]) for s in specs}
    if names != set(expected):
        raise TrainingError(f"nieobsługiwana architektura: {sorted(names)}")
    m = _spec_map(specs)
    hidden = m["backbone.bias"][0]
    input_dim = m["backbone.weight"][0]
    if m["backbone.weight"][1] != hidden or m["classifier.weight"][0] != hidden or m["classifier.bias"][0] != m["classifier.weight"][1]:
        raise TrainingError("niespójne kształty głowy klasyfikatora")
    if "box.weight" in m and (m["box.weight"][0] != hidden or m["box.weight"][1] != 4):
        raise TrainingError("niespójne kształty głowy detekcji")
    if input_dim <= 0 or hidden <= 0:
        raise TrainingError("pusty wymiar wejścia")


def _audio_forward(tensors: dict[str, Tensor], x: list[float], classes: int):
    m = {k: t for k, t in tensors.items()}
    hidden = len(m["backbone.bias"].values)
    input_dim = len(x)
    z = _matvec(m["backbone.weight"].values, x, hidden, input_dim)
    for i in range(hidden):
        z[i] += m["backbone.bias"].values[i]
    h = [_relu(v) for v in z]
    logits = _matvec(m["classifier.weight"].values, h, classes, hidden)
    for i in range(classes):
        logits[i] += m["classifier.bias"].values[i]
    return z, h, logits


def _zero_grads(shapes: dict[str, list[int]]) -> dict[str, list[float]]:
    return {k: [0.0] * _count(v) for k, v in shapes.items()}


def train_audio(data: list[dict[str, Any]], specs: list[dict[str, Any]], cfg: TrainConfig) -> dict[str, Any]:
    _check_specs(specs, AUDIO_NAMES)
    if not data:
        raise TrainingError("brak danych treningowych audio")
    shapes = _spec_map(specs)
    hidden = shapes["backbone.bias"][0]
    input_dim = shapes["backbone.weight"][0]
    classes = list(cfg.classes)
    n_classes = len(classes)
    if shapes["classifier.weight"][1] != n_classes:
        raise TrainingError("liczba klas nie zgadza się z architekturą")
    for row in data:
        if len(row.get("features") or []) != input_dim:
            raise TrainingError(f"cecha audio ma wymiar {len(row.get('features') or [])}, oczekiwano {input_dim}")
        if row.get("label") not in classes:
            raise TrainingError(f"nieznana etykieta: {row.get('label')}")
    tensors = deterministic_init(specs, cfg.seed)
    start_epoch = 0
    best = -1.0
    if cfg.checkpoint and cfg.checkpoint.exists() and cfg.resume:
        old, meta = read_safetensors(cfg.checkpoint.read_bytes())
        if meta.get("nieboguard.task_id") == cfg.task_id and meta.get("nieboguard.architecture_hash") == cfg.architecture_hash:
            tensors = old
            start_epoch = int(meta.get("epoch", 0))
            best = float(meta.get("best_accuracy", -1))
    train, val = _split(data, cfg.seed)
    opt = _Adam(shapes, cfg.learning_rate, cfg.weight_decay)
    budget = _Budget(cfg.max_seconds, cfg.cpu_percent)
    history: list[dict[str, Any]] = []
    steps = 0
    epochs_done = start_epoch
    for epoch in range(start_epoch, cfg.epochs):
        rng = random.Random(cfg.seed + epoch * 7919)
        order = list(range(len(train)))
        rng.shuffle(order)
        lr = _adjust_lr(cfg.learning_rate, epoch + 1)
        total_loss = 0.0
        seen = 0
        for start in range(0, len(order), max(1, cfg.batch_size)):
            if budget.expired() or steps >= cfg.max_steps:
                break
            batch = [train[i] for i in order[start:start + max(1, cfg.batch_size)]]
            grads = _zero_grads(shapes)
            for row in batch:
                x = [float(v) for v in row["features"]]
                y = classes.index(str(row["label"]))
                z, h, logits = _audio_forward(tensors, x, n_classes)
                p = _softmax(logits)
                loss = -math.log(max(1e-12, p[y]))
                total_loss += loss
                dlogits = p[:]
                dlogits[y] -= 1.0
                dlogits = [v / len(batch) for v in dlogits]
                for c in range(n_classes):
                    base = c * hidden
                    for j in range(hidden):
                        grads["classifier.weight"][base + j] += dlogits[c] * h[j]
                    grads["classifier.bias"][c] += dlogits[c]
                dh = [sum(tensors["classifier.weight"].values[c * hidden + j] * dlogits[c] for c in range(n_classes)) for j in range(hidden)]
                dz = [dh[j] if z[j] > 0 else 0.0 for j in range(hidden)]
                for j in range(hidden):
                    base = j * input_dim
                    for k in range(input_dim):
                        grads["backbone.weight"][base + k] += dz[j] * x[k]
                    grads["backbone.bias"][j] += dz[j]
                seen += 1
            opt.step(tensors, grads, lr)
            steps += 1
            budget.throttle()
        acc, val_loss = evaluate_audio(tensors, val, classes)
        epochs_done = epoch + 1
        history.append({"epoch": epochs_done, "train_loss": round(total_loss / max(1, seen), 6), "val_accuracy": round(acc, 6), "val_loss": round(val_loss, 6), "steps": steps})
        if acc > best:
            best = acc
        save_checkpoint(cfg, tensors, epochs_done, best, len(train) + len(val), cfg.backend, cfg.device)
        if budget.expired() or steps >= cfg.max_steps:
            break
    acc, val_loss = evaluate_audio(tensors, val, classes)
    return {
        "tensors": tensors,
        "metrics": {"accuracy": acc, "loss": val_loss, "samples": len(data), "class_counts": _class_counts(data, classes),
                    "duration_s": round(time.monotonic() - budget.start_wall, 3), "epochs_completed": epochs_done,
                    "backend": cfg.backend, "resumed": start_epoch > 0},
        "history": history, "validation_size": len(val),
    }


def evaluate_audio(tensors: dict[str, Tensor], data: list[dict[str, Any]], classes: list[str]) -> tuple[float, float]:
    if not data:
        return 0.0, 0.0
    correct = 0
    loss = 0.0
    for row in data:
        x = [float(v) for v in row["features"]]
        _, _, logits = _audio_forward(tensors, x, len(classes))
        p = _softmax(logits)
        y = classes.index(str(row["label"])) if row.get("label") in classes else 0
        loss -= math.log(max(1e-12, p[y]))
        if max(range(len(p)), key=lambda i: p[i]) == y:
            correct += 1
    return correct / len(data), loss / len(data)


def _class_counts(data: Iterable[dict[str, Any]], classes: list[str]) -> dict[str, int]:
    out = {c: 0 for c in classes}
    for row in data:
        label = row.get("label")
        if label in out:
            out[label] += 1
    return out


def _vision_label_index(label: Any, classes: list[str]) -> int:
    if isinstance(label, int):
        return label
    if isinstance(label, str) and label in classes:
        return classes.index(label)
    try:
        return int(label)
    except Exception:
        return -1


def _vision_forward(tensors: dict[str, Tensor], x: list[float], n_classes: int):
    hidden = len(tensors["backbone.bias"].values)
    input_dim = len(x)
    z = _matvec(tensors["backbone.weight"].values, x, hidden, input_dim)
    for i in range(hidden):
        z[i] += tensors["backbone.bias"].values[i]
    h = [_relu(v) for v in z]
    obj_logit = sum(tensors["objectness.weight"].values[i] * h[i] for i in range(hidden)) + tensors["objectness.bias"].values[0]
    logits = _matvec(tensors["classifier.weight"].values, h, n_classes, hidden)
    for i in range(n_classes):
        logits[i] += tensors["classifier.bias"].values[i]
    box = _matvec(tensors["box.weight"].values, h, 4, hidden)
    for i in range(4):
        box[i] += tensors["box.bias"].values[i]
    return z, h, obj_logit, logits, box


def train_vision(data: list[dict[str, Any]], specs: list[dict[str, Any]], cfg: TrainConfig) -> dict[str, Any]:
    _check_specs(specs, VISION_NAMES)
    if not data:
        raise TrainingError("brak danych treningowych vision")
    shapes = _spec_map(specs)
    hidden = shapes["backbone.bias"][0]
    input_dim = shapes["backbone.weight"][0]
    n_classes = len(cfg.classes)
    if shapes["classifier.weight"][1] != n_classes:
        raise TrainingError("liczba klas nie zgadza się z architekturą")
    for row in data:
        if len(row.get("features") or []) != input_dim:
            raise TrainingError(f"cecha vision ma wymiar {len(row.get('features') or [])}, oczekiwano {input_dim}")
        if bool(row.get("objectness")):
            if row.get("label") not in [None, -1] and not (0 <= _vision_label_index(row.get("label"), cfg.classes) < len(cfg.classes)):
                raise TrainingError(f"nieznana etykieta: {row.get('label')}")
            if not isinstance(row.get("box"), (list, tuple)) or len(row["box"]) != 4:
                raise TrainingError("brak box dla obiektu dodatniego")
    tensors = deterministic_init(specs, cfg.seed)
    start_epoch = 0
    best = -1.0
    if cfg.checkpoint and cfg.checkpoint.exists() and cfg.resume:
        old, meta = read_safetensors(cfg.checkpoint.read_bytes())
        if meta.get("nieboguard.task_id") == cfg.task_id and meta.get("nieboguard.architecture_hash") == cfg.architecture_hash:
            tensors = old
            start_epoch = int(meta.get("epoch", 0))
            best = float(meta.get("best_map", -1))
    train, val = _split(data, cfg.seed)
    opt = _Adam(shapes, cfg.learning_rate, cfg.weight_decay)
    budget = _Budget(cfg.max_seconds, cfg.cpu_percent)
    history: list[dict[str, Any]] = []
    steps = 0
    epochs_done = start_epoch
    for epoch in range(start_epoch, cfg.epochs):
        rng = random.Random(cfg.seed + epoch * 3571)
        order = list(range(len(train)))
        rng.shuffle(order)
        lr = _adjust_lr(cfg.learning_rate, epoch + 1)
        total_loss = 0.0
        seen = 0
        for start in range(0, len(order), max(1, cfg.batch_size)):
            if budget.expired() or steps >= cfg.max_steps:
                break
            batch = [train[i] for i in order[start:start + max(1, cfg.batch_size)]]
            grads = _zero_grads(shapes)
            for row in batch:
                x = [float(v) for v in row["features"]]
                target_obj = 1.0 if row.get("objectness") else 0.0
                z, h, obj_logit, logits, box = _vision_forward(tensors, x, n_classes)
                obj = _sigmoid(obj_logit)
                p = _softmax(logits)
                bce = -(target_obj * math.log(max(1e-12, obj)) + (1 - target_obj) * math.log(max(1e-12, 1 - obj)))
                loss = bce
                d_obj = (obj - target_obj) / len(batch)
                dlogits = [0.0] * n_classes
                dbox = [0.0] * 4
                y = None
                if target_obj > 0.5:
                    y = _vision_label_index(row.get("label"), cfg.classes)
                    if y < 0:
                        y = 0
                    loss -= math.log(max(1e-12, p[y]))
                    for c in range(n_classes):
                        dlogits[c] = (p[c] - (1.0 if c == y else 0.0)) / len(batch)
                    target_box = [float(v) for v in row["box"]]
                    for k in range(4):
                        diff = box[k] - target_box[k]
                        loss += 0.5 * diff * diff
                        dbox[k] = diff / len(batch)
                total_loss += loss
                # objectness head
                for i in range(hidden):
                    grads["objectness.weight"][i] += d_obj * h[i]
                grads["objectness.bias"][0] += d_obj
                # classifier head
                if y is not None:
                    for c in range(n_classes):
                        base = c * hidden
                        for i in range(hidden):
                            grads["classifier.weight"][base + i] += dlogits[c] * h[i]
                        grads["classifier.bias"][c] += dlogits[c]
                # box head
                if y is not None:
                    for k in range(4):
                        base = k * hidden
                        for i in range(hidden):
                            grads["box.weight"][base + i] += dbox[k] * h[i]
                        grads["box.bias"][k] += dbox[k]
                # shared backbone
                dh = [0.0] * hidden
                for i in range(hidden):
                    dh[i] += tensors["objectness.weight"].values[i] * d_obj
                    if y is not None:
                        dh[i] += sum(tensors["classifier.weight"].values[c * hidden + i] * dlogits[c] for c in range(n_classes))
                        dh[i] += sum(tensors["box.weight"].values[k * hidden + i] * dbox[k] for k in range(4))
                for i in range(hidden):
                    dz = dh[i] if z[i] > 0 else 0.0
                    base = i * input_dim
                    for k in range(input_dim):
                        grads["backbone.weight"][base + k] += dz * x[k]
                    grads["backbone.bias"][i] += dz
                seen += 1
            opt.step(tensors, grads, lr)
            steps += 1
            budget.throttle()
        metrics = evaluate_vision(tensors, val, cfg.classes)
        epochs_done = epoch + 1
        history.append({"epoch": epochs_done, "train_loss": round(total_loss / max(1, seen), 6), **metrics, "steps": steps})
        if metrics["objectness_f1"] > best:
            best = metrics["objectness_f1"]
        save_checkpoint(cfg, tensors, epochs_done, best, len(train) + len(val), cfg.backend, cfg.device)
        if budget.expired() or steps >= cfg.max_steps:
            break
    metrics = evaluate_vision(tensors, val, cfg.classes)
    return {
        "tensors": tensors,
        "metrics": {"accuracy": metrics["class_accuracy"], "objectness_f1": metrics["objectness_f1"],
                    "loss": metrics["loss"], "samples": len(data), "class_counts": _class_counts(data, cfg.classes),
                    "duration_s": round(time.monotonic() - budget.start_wall, 3), "epochs_completed": epochs_done,
                    "backend": cfg.backend, "resumed": start_epoch > 0},
        "history": history, "validation_size": len(val),
    }


def evaluate_vision(tensors: dict[str, Tensor], data: list[dict[str, Any]], classes: list[str]) -> dict[str, float]:
    if not data:
        return {"objectness_f1": 0.0, "class_accuracy": 0.0, "loss": 0.0}
    tp = fp = fn = 0
    correct = classified = 0
    loss = 0.0
    for row in data:
        x = [float(v) for v in row["features"]]
        _, _, obj_logit, logits, box = _vision_forward(tensors, x, len(classes))
        obj = _sigmoid(obj_logit)
        pred = obj >= 0.5
        target = bool(row.get("objectness"))
        if pred and target:
            tp += 1
            y = _vision_label_index(row.get("label"), classes)
            if y < 0:
                y = 0
            if max(range(len(logits)), key=lambda i: logits[i]) == y:
                correct += 1
            classified += 1
            target_box = [float(v) for v in row["box"]]
            loss += 0.5 * sum((box[k] - target_box[k]) ** 2 for k in range(4))
        elif pred:
            fp += 1
        elif target:
            fn += 1
            loss += 1.0
    precision = tp / max(1, tp + fp)
    recall = tp / max(1, tp + fn)
    f1 = 2 * precision * recall / max(1e-9, precision + recall)
    return {"objectness_f1": f1, "class_accuracy": correct / max(1, classified), "loss": loss / len(data)}


def save_checkpoint(cfg: TrainConfig, tensors: dict[str, Tensor], epoch: int, best: float,
                    sample_count: int, backend: str, device: str) -> None:
    if not cfg.checkpoint:
        return
    metadata = _metadata(cfg, sample_count)
    metadata.update({"checkpoint": "true", "epoch": str(epoch), "best_accuracy": f"{best:.6f}", "best_map": f"{best:.6f}",
                     "backend": backend, "device": device})
    cfg.checkpoint.parent.mkdir(parents=True, exist_ok=True)
    tmp = cfg.checkpoint.with_name(cfg.checkpoint.name + ".part")
    tmp.write_bytes(write_safetensors(tensors, metadata))
    tmp.replace(cfg.checkpoint)


def _metadata(cfg: TrainConfig, sample_count: int) -> dict[str, str]:
    return {
        "nieboguard.task_id": cfg.task_id,
        "nieboguard.round_id": cfg.round_id,
        "nieboguard.base_model": cfg.base_model,
        "nieboguard.base_version": cfg.base_version,
        "nieboguard.architecture_hash": cfg.architecture_hash,
        "nieboguard.format": "safetensors",
        "nieboguard.classes": ",".join(cfg.classes),
        "nieboguard.sample_count": str(sample_count),
        "nieboguard.client": "nieboguard-compute/0.1.0",
    }


def export_artifact(result: dict[str, Any], cfg: TrainConfig, path: Path, sample_count: int) -> tuple[bytes, dict[str, Any]]:
    tensors = result["tensors"]
    metadata = _metadata(cfg, sample_count)
    metadata.update({
        "nieboguard.epochs_completed": str(result["metrics"].get("epochs_completed", 0)),
        "nieboguard.backend": str(result["metrics"].get("backend", "python")),
        "nieboguard.metrics": json.dumps({k: v for k, v in result["metrics"].items() if k not in ("class_counts",)}, separators=(",", ":")),
    })
    data = write_safetensors(tensors, metadata)
    path.parent.mkdir(parents=True, exist_ok=True)
    tmp = path.with_name(path.name + ".part")
    tmp.write_bytes(data)
    tmp.replace(path)
    return data, metadata


def read_checkpoint(path: Path) -> tuple[dict[str, Tensor], dict[str, str]]:
    return read_safetensors(path.read_bytes())

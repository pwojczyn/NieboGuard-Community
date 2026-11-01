"""Vision AI: fine-tuning a detection head over frozen proposal features.

No images are sent to the backend in this MVP. The signed task provides
proposal features and boxes from a trusted, redistributable dataset. The
volunteer trains objectness, class and box-regression heads. An optional
PyTorch backend can use an NVIDIA GPU, but the standard-library path works on
CPU.
"""
from __future__ import annotations

import json
import math
from pathlib import Path
from typing import Any

MAX_RECORDS = 200000
MAX_FEATURES = 4096


class VisionError(ValueError):
    pass


def load_vision_dataset(data_dir: Path, input_dim: int, classes: list[str], split_by_file: dict[str, str] | None = None) -> list[dict[str, Any]]:
    files = sorted(p for p in data_dir.iterdir() if p.is_file() and p.suffix.lower() == ".jsonl")
    if not files:
        raise VisionError("Brak pliku JSONL z propozycjami vision.")
    rows: list[dict[str, Any]] = []
    for label_file in files:
        default_split = (split_by_file or {}).get(label_file.name, "train")
        for line_no, line in enumerate(label_file.read_text(encoding="utf-8").splitlines(), 1):
            line = line.strip()
            if not line:
                continue
            try:
                item = json.loads(line)
            except Exception as exc:
                raise VisionError(f"{label_file.name}:{line_no}: nieprawidłowy JSON") from exc
            if not isinstance(item, dict):
                raise VisionError(f"{label_file.name}:{line_no}: wiersz nie jest obiektem")
            features = item.get("features")
            if not isinstance(features, list) or not features or len(features) > MAX_FEATURES:
                raise VisionError(f"{label_file.name}:{line_no}: brak cech propozycji")
            if len(features) != input_dim:
                raise VisionError(f"{label_file.name}:{line_no}: wymiar cech {len(features)} != {input_dim}")
            values = []
            for v in features:
                f = float(v)
                if not math.isfinite(f):
                    raise VisionError(f"{label_file.name}:{line_no}: NaN/Inf w cechach")
                values.append(f)
            objectness = bool(item.get("objectness", False))
            label = item.get("label")
            if objectness:
                if isinstance(label, str):
                    if label not in classes:
                        raise VisionError(f"{label_file.name}:{line_no}: nieznana klasa")
                    label = classes.index(label)
                else:
                    try:
                        label = int(label)
                    except Exception as exc:
                        raise VisionError(f"{label_file.name}:{line_no}: brak klasy") from exc
                    if not (0 <= label < len(classes)):
                        raise VisionError(f"{label_file.name}:{line_no}: klasa poza zakresem")
                box = item.get("box")
                if not isinstance(box, list) or len(box) != 4:
                    raise VisionError(f"{label_file.name}:{line_no}: brak box")
                box = [float(v) for v in box]
                if any(not math.isfinite(v) or v < -0.5 or v > 1.5 for v in box):
                    raise VisionError(f"{label_file.name}:{line_no}: box poza zakresem")
            else:
                label = -1
                box = [0.0, 0.0, 0.0, 0.0]
            rows.append({
                "features": values, "objectness": objectness, "label": label, "box": box,
                "split": default_split if split_by_file else str(item.get("split") or "train"),
                "group": str(item.get("group") or label_file.stem),
            })
            if len(rows) > MAX_RECORDS:
                raise VisionError("Zbyt wiele rekordów vision.")
    if not rows:
        raise VisionError("Brak rekordów vision.")
    return rows

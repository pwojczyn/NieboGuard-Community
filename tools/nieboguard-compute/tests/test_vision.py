import json
import math
import random
import tempfile
import unittest
from pathlib import Path
import sys

sys.path.insert(0, str(Path(__file__).resolve().parents[1] / "src"))

from nieboguard_compute.model import TrainConfig, export_artifact, train_vision
from nieboguard_compute.safetensors import read_safetensors
from nieboguard_compute.vision import load_vision_dataset

INPUT = 12
CLASSES = ["drone", "aircraft"]


def row(rng, positive: bool, label: int = 0) -> dict:
    if positive:
        feat = [rng.gauss(1.5, 0.25), rng.gauss(1.0, 0.25)] + [rng.gauss(0.0, 0.2) for _ in range(INPUT - 2)]
        return {"features": feat, "objectness": True, "label": CLASSES[label],
                "box": [0.5 + rng.uniform(-0.05, 0.05), 0.5, 0.2, 0.2]}
    feat = [rng.gauss(-1.5, 0.25), rng.gauss(-1.0, 0.25)] + [rng.gauss(0.0, 0.2) for _ in range(INPUT - 2)]
    return {"features": feat, "objectness": False}


SPECS = [
    {"name": "backbone.weight", "shape": [INPUT, 8], "dtype": "F32"},
    {"name": "backbone.bias", "shape": [8], "dtype": "F32"},
    {"name": "classifier.weight", "shape": [8, 2], "dtype": "F32"},
    {"name": "classifier.bias", "shape": [2], "dtype": "F32"},
    {"name": "objectness.weight", "shape": [8, 1], "dtype": "F32"},
    {"name": "objectness.bias", "shape": [1], "dtype": "F32"},
    {"name": "box.weight", "shape": [8, 4], "dtype": "F32"},
    {"name": "box.bias", "shape": [4], "dtype": "F32"},
]


class VisionTrainingTest(unittest.TestCase):
    def test_detection_head_cpu(self):
        rng = random.Random(42)
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            data_dir = root / "data"
            data_dir.mkdir()
            rows = []
            for i in range(40):
                rows.append(row(rng, True, i % 2))
                rows.append(row(rng, False))
            for i, item in enumerate(rows):
                item["split"] = "validation" if i % 5 == 0 else "train"
            (data_dir / "proposals.jsonl").write_text("\n".join(json.dumps(x) for x in rows), encoding="utf-8")
            data = load_vision_dataset(data_dir, INPUT, CLASSES)
            cfg = TrainConfig(
                task_id="task-vis", round_id="round-1", base_model="base", base_version="0.1",
                architecture_hash="arch", classes=CLASSES, seed=3, epochs=12, batch_size=16,
                learning_rate=0.05, max_steps=200, max_seconds=30,
                checkpoint=root / "checkpoint.safetensors", resume=True,
            )
            result = train_vision(data, SPECS, cfg)
            self.assertGreaterEqual(result["metrics"]["objectness_f1"], 0.8)
            self.assertTrue(all(math.isfinite(v) for t in result["tensors"].values() for v in t.values))
            artifact, meta = export_artifact(result, cfg, root / "artifact.safetensors", len(data))
            tensors, _ = read_safetensors(artifact)
            self.assertEqual(set(tensors), {x["name"] for x in SPECS})
            self.assertEqual(meta["nieboguard.task_id"], "task-vis")


if __name__ == "__main__":
    unittest.main()

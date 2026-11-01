import json
import math
import random
import struct
import tempfile
import unittest
import wave
from pathlib import Path
import sys

sys.path.insert(0, str(Path(__file__).resolve().parents[1] / "src"))

from nieboguard_compute.audio import FEATURE_DIM, load_audio_dataset, wav_features
from nieboguard_compute.model import TrainConfig, export_artifact, load_base_tensors, train_audio
from nieboguard_compute.safetensors import Tensor, write_safetensors
from nieboguard_compute.safetensors import read_safetensors


def write_wav(path: Path, samples: list[float], rate: int = 16000):
    with wave.open(str(path), "wb") as wav:
        wav.setnchannels(1)
        wav.setsampwidth(2)
        wav.setframerate(rate)
        wav.writeframes(b"".join(struct.pack("<h", max(-32768, min(32767, int(s * 32767)))) for s in samples))


def tone(freq: float, seconds: float = 1.0, rate: int = 16000) -> list[float]:
    n = int(seconds * rate)
    return [0.25 * math.sin(2 * math.pi * freq * i / rate) for i in range(n)]


def noisy(freq: float, seconds: float = 1.0, rate: int = 16000, seed: int = 1) -> list[float]:
    rng = random.Random(seed)
    n = int(seconds * rate)
    return [0.1 * math.sin(2 * math.pi * freq * i / rate) + 0.08 * rng.uniform(-1, 1) for i in range(n)]


class AudioTrainingTest(unittest.TestCase):
    def test_base_artifact_is_loaded(self):
        specs = [
            {"name": "backbone.weight", "shape": [FEATURE_DIM, 4], "dtype": "F32"},
            {"name": "backbone.bias", "shape": [4], "dtype": "F32"},
            {"name": "classifier.weight", "shape": [4, 2], "dtype": "F32"},
            {"name": "classifier.bias", "shape": [2], "dtype": "F32"},
        ]
        values = {}
        for spec in specs:
            n = 1
            for d in spec["shape"]:
                n *= d
            values[spec["name"]] = Tensor(spec["name"], spec["shape"], [0.123] * n)
        with tempfile.TemporaryDirectory() as tmp:
            path = Path(tmp) / "base.safetensors"
            path.write_bytes(write_safetensors(values, {"base": "yes"}))
            loaded = load_base_tensors(specs, path)
            self.assertAlmostEqual(loaded["classifier.bias"].values[0], 0.123, places=5)
            self.assertAlmostEqual(loaded["classifier.bias"].values[1], 0.123, places=5)

    def test_manifest_shard_split_overrides_missing_row_split(self):
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            row = {"features": [0.0] * FEATURE_DIM, "label": "drone"}
            (root / "val.jsonl").write_text(json.dumps(row) + "\n", encoding="utf-8")
            data = load_audio_dataset(root, FEATURE_DIM, ["drone", "bird"], {"val.jsonl": "validation"})
            self.assertEqual(data[0]["split"], "validation")

    def test_features_are_finite_and_stable(self):
        with tempfile.TemporaryDirectory() as tmp:
            p = Path(tmp) / "a.wav"
            write_wav(p, tone(220))
            f = wav_features(p)
            self.assertEqual(len(f), FEATURE_DIM)
            self.assertTrue(all(math.isfinite(x) for x in f))

    def test_small_cpu_training_and_resume(self):
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            data_dir = root / "data"
            data_dir.mkdir()
            rows = []
            for label, freq in (("drone", 180.0), ("bird", 1800.0)):
                for i in range(12):
                    name = f"{label}-{i}.wav"
                    samples = noisy(freq, seed=i) if i % 2 else tone(freq)
                    write_wav(data_dir / name, samples)
                    rows.append({"file": name, "label": label, "split": "train" if i < 10 else "validation"})
            (data_dir / "labels.jsonl").write_text("\n".join(json.dumps(x) for x in rows), encoding="utf-8")
            data = load_audio_dataset(data_dir, FEATURE_DIM, ["drone", "bird"])
            specs = [
                {"name": "backbone.weight", "shape": [FEATURE_DIM, 16], "dtype": "F32"},
                {"name": "backbone.bias", "shape": [16], "dtype": "F32"},
                {"name": "classifier.weight", "shape": [16, 2], "dtype": "F32"},
                {"name": "classifier.bias", "shape": [2], "dtype": "F32"},
            ]
            cfg = TrainConfig(
                task_id="task-audio", round_id="round-1", base_model="base", base_version="0.1",
                architecture_hash="arch", classes=["drone", "bird"], seed=7, epochs=2,
                batch_size=8, learning_rate=0.05, max_steps=100, max_seconds=30,
                checkpoint=root / "checkpoint.safetensors", resume=True,
            )
            result = train_audio(data, specs, cfg)
            self.assertTrue(all(math.isfinite(v) for t in result["tensors"].values() for v in t.values))
            self.assertGreaterEqual(result["metrics"]["accuracy"], 0.8)
            artifact, metadata = export_artifact(result, cfg, root / "artifact.safetensors", len(data))
            tensors, old_meta = read_safetensors(artifact)
            self.assertEqual(set(tensors), {x["name"] for x in specs})
            self.assertEqual(old_meta["nieboguard.task_id"], "task-audio")
            self.assertEqual(old_meta["nieboguard.sample_count"], str(len(data)))
            cfg2 = TrainConfig(
                task_id="task-audio", round_id="round-1", base_model="base", base_version="0.1",
                architecture_hash="arch", classes=["drone", "bird"], seed=7, epochs=4,
                batch_size=8, learning_rate=0.05, max_steps=100, max_seconds=30,
                checkpoint=root / "checkpoint.safetensors", resume=True,
            )
            result2 = train_audio(data, specs, cfg2)
            self.assertGreaterEqual(result2["metrics"]["epochs_completed"], 3)
            stopped = train_audio(data, specs, TrainConfig(
                task_id="task-audio", round_id="round-1", base_model="base", base_version="0.1",
                architecture_hash="arch", classes=["drone", "bird"], epochs=100, batch_size=8,
                learning_rate=0.05, max_steps=100, max_seconds=30,
                checkpoint=root / "stop.safetensors", resume=False,
                should_stop=lambda: True,
            ))
            self.assertTrue(stopped["stopped"])


if __name__ == "__main__":
    unittest.main()

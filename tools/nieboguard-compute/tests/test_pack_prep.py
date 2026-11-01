import argparse
import importlib.util
import json
import math
import struct
import tempfile
import unittest
import wave
import zipfile
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
SPEC = importlib.util.spec_from_file_location("pack_prep", ROOT / "academic" / "prepare_academic_pack.py")
pack = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(pack)


def wav_bytes(freq: float, seconds: float = 0.05) -> bytes:
    import io
    buf = io.BytesIO()
    with wave.open(buf, "wb") as w:
        w.setnchannels(1)
        w.setsampwidth(2)
        w.setframerate(8000)
        n = int(8000 * seconds)
        w.writeframes(b"".join(struct.pack("<h", int(12000 * math.sin(2 * math.pi * freq * i / 8000))) for i in range(n)))
    return buf.getvalue()


class PackPrepTest(unittest.TestCase):
    def make_zip(self, tmp: Path) -> Path:
        archive = tmp / "source.zip"
        with zipfile.ZipFile(archive, "w") as zf:
            zf.writestr("sessionA/drone/a.wav", wav_bytes(220))
            zf.writestr("sessionA/drone/b.wav", wav_bytes(230))
            zf.writestr("sessionB/bird/c.wav", wav_bytes(1800))
            zf.writestr("sessionB/notes.txt", "ignore")
        return archive

    def test_prepare_groups_and_labels(self):
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            archive = self.make_zip(root)
            labels = root / "labels.json"
            labels.write_text(json.dumps({"drone": "drone_multirotor", "bird": "bird"}))
            packcfg = root / "packs.json"
            packcfg.write_text(json.dumps({"MINI": {"min_mb": 0.0001, "max_mb": 1, "target_mb": 1}}))
            args = argparse.Namespace(
                source=str(archive), source_dir=None, sha256=None, work=str(root / "work"), out=str(root / "out"),
                dataset_id="test-audio", version=1, name="Test", modality="audio", license="CC0-1.0",
                license_url=None, authors=["Test"], provenance="synthetic", source_url=None,
                commercial=False, no_redistribution=False, label_map=str(labels), labels=[],
                allow_unknown=False, allow_impossible_packs=False, exclude_regex=None, session_depth=1,
                max_file_mb=120, max_extract_gb=0.1, max_files=1000, min_free_gb=0,
                pack_tiers=["MINI"], pack_config=str(packcfg), model_bytes=0, library_bytes=0,
                reuse_extracted=False, no_resume=True, upload=False, build_packs=False,
                api="http://127.0.0.1", admin_token="",
            )
            dataset = pack.prepare(args)
            self.assertEqual(dataset["groups"], 2)
            self.assertEqual(dataset["sample_count"], 3)
            self.assertEqual(dataset["classes"], ["bird", "drone_multirotor"])
            self.assertEqual(len(dataset["shards"]), 5)
            label_shards = [s for s in dataset["shards"] if s["name"].endswith(".jsonl")]
            self.assertEqual(len(label_shards), 2)
            by_group = {}
            for shard in dataset["shards"]:
                by_group.setdefault(shard["group_key"], set()).add(shard["split"])
            self.assertTrue(all(len(v) == 1 for v in by_group.values()))
            for shard in label_shards:
                rows = [json.loads(x) for x in Path(shard["source_path"]).read_text().splitlines() if x.strip()]
                media_names = {s["name"] for s in dataset["shards"] if s["name"].endswith(".wav")}
                self.assertTrue(all(row["file"] in media_names for row in rows))
                self.assertTrue(all(row["group"] == shard["group_key"] for row in rows))

    def test_zip_slip_is_rejected(self):
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            archive = root / "evil.zip"
            with zipfile.ZipFile(archive, "w") as zf:
                zf.writestr("../evil.wav", wav_bytes(200))
            with self.assertRaises(pack.PackError):
                pack.safe_extract(archive, root / "out")

    def test_absolute_path_is_rejected(self):
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            archive = root / "abs.zip"
            with zipfile.ZipFile(archive, "w") as zf:
                zf.writestr("/tmp/evil.wav", wav_bytes(200))
            with self.assertRaises(pack.PackError):
                pack.safe_extract(archive, root / "out")


if __name__ == "__main__":
    unittest.main()

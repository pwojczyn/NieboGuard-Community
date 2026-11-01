import argparse
import importlib.util
import json
import math
import struct
import tempfile
import unittest
import wave
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
SPEC = importlib.util.spec_from_file_location("audio_prep", ROOT / "academic" / "prepare_audio_pack.py")
audio = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(audio)


def write_wav(path: Path, freq: float, seconds: float = 1.2, amp: float = 0.2):
    path.parent.mkdir(parents=True, exist_ok=True)
    with wave.open(str(path), "wb") as w:
        w.setnchannels(1)
        w.setsampwidth(2)
        w.setframerate(8000)
        n = int(8000 * seconds)
        w.writeframes(b"".join(struct.pack("<h", int(amp * 32767 * math.sin(2 * math.pi * freq * i / 8000))) for i in range(n)))


class AudioPackTest(unittest.TestCase):
    def test_segmentation_quality_and_pack_plan(self):
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            src = root / "source"
            write_wav(src / "drone" / "session-a" / "one.wav", 220)
            write_wav(src / "drone" / "session-a" / "two.wav", 230)
            write_wav(src / "bird" / "session-b" / "three.wav", 1800)
            write_wav(src / "drone" / "session-a" / "silence.wav", 220, amp=0.0)
            # duplicate content should be deduplicated by segment SHA
            write_wav(src / "drone" / "session-a" / "duplicate.wav", 220)
            labels = root / "labels.json"
            labels.write_text(json.dumps({"drone": "drone_multirotor", "bird": "bird"}))
            packcfg = root / "packs.json"
            packcfg.write_text(json.dumps({"MINI": {"min_mb": 0.0001, "max_mb": 100, "target_mb": 100}}))
            args = argparse.Namespace(
                source=None, source_dir=str(src), sha256=None, work=str(root / "work"), out=str(root / "out"),
                dataset_id="test-audio", version=1, name="Test audio", license="CC0-1.0",
                license_url=None, authors=["Test"], provenance="synthetic", source_url=None,
                commercial=False, no_redistribution=False, label_map=str(labels), labels=[],
                allow_unknown=False, allow_impossible_packs=False, exclude_regex="silence",
                session_depth=2, segment_seconds=0.4, overlap_seconds=0.1, min_segment_seconds=0.1,
                min_rms=0.005, max_clipping=0.01, max_file_mb=100, max_extract_gb=0.1,
                max_files=1000, min_free_gb=0, pack_tiers=["MINI"], pack_config=str(packcfg),
                model_bytes=0, library_bytes=0, reuse_extracted=False, no_resume=True,
                upload=False, build_packs=False, api="http://127.0.0.1", admin_token="",
            )
            dataset = audio.audio_prepare(args)
            self.assertGreater(dataset["sample_count"], 0)
            self.assertGreater(dataset["audio"]["duplicates"], 0)  # duplicate file deduplicated
            self.assertEqual(dataset["classes"], ["bird", "drone_multirotor"])
            self.assertTrue((root / "out" / "dataset.json").exists())
            self.assertTrue((root / "out" / "audio-summary.json").exists())


if __name__ == "__main__":
    unittest.main()

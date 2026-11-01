import importlib.util
import tempfile
import unittest
import wave
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
spec = importlib.util.spec_from_file_location("wav_audit", ROOT / "tools/audio-ai/wav_audit.py")
wav_audit = importlib.util.module_from_spec(spec)
spec.loader.exec_module(wav_audit)


class WavAuditTests(unittest.TestCase):
    def test_reads_metadata_without_upload(self):
        with tempfile.TemporaryDirectory() as tmp:
            path = Path(tmp) / "silence.wav"
            with wave.open(str(path), "wb") as wav:
                wav.setnchannels(1)
                wav.setsampwidth(2)
                wav.setframerate(16000)
                wav.writeframes(b"\x00\x00" * 1600)
            result = wav_audit.audit(path)
            self.assertAlmostEqual(result["duration_s"], 0.1, places=3)
            self.assertEqual(result["sample_rate_hz"], 16000)
            self.assertEqual(result["channels"], 1)
            self.assertEqual(result["bits_per_sample"], 16)


if __name__ == "__main__":
    unittest.main()

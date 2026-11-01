import importlib.util
import tempfile
import unittest
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
spec = importlib.util.spec_from_file_location("image_audit", ROOT / "tools/vision-ai/image_audit.py")
image_audit = importlib.util.module_from_spec(spec)
spec.loader.exec_module(image_audit)


@unittest.skipUnless(image_audit.Image is not None, "Pillow not installed")
class ImageAuditTests(unittest.TestCase):
    def test_reads_and_strips_copy(self):
        with tempfile.TemporaryDirectory() as tmp:
            source = Path(tmp) / "sample.png"
            target = Path(tmp) / "clean.png"
            image_audit.Image.new("RGB", (16, 12), (10, 20, 30)).save(source)
            result = image_audit.audit(source, target)
            self.assertEqual(result["width"], 16)
            self.assertEqual(result["height"], 12)
            self.assertTrue(target.exists())

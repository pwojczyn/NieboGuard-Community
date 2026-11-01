import importlib.util
import unittest
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
spec = importlib.util.spec_from_file_location("rssi", ROOT / "tools/sensors/rssi_distance.py")
rssi = importlib.util.module_from_spec(spec)
spec.loader.exec_module(rssi)


class RssiTests(unittest.TestCase):
    def test_reference_distance(self):
        self.assertAlmostEqual(rssi.estimate(-59, -59, 2), 1.0, places=6)

    def test_weaker_signal_is_farther(self):
        self.assertGreater(rssi.estimate(-80, -59, 2), rssi.estimate(-70, -59, 2))

    def test_rejects_bad_exponent(self):
        with self.assertRaises(ValueError):
            rssi.estimate(-70, -59, 0)


if __name__ == "__main__":
    unittest.main()

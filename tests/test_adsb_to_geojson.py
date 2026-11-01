import importlib.util
import json
import unittest
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
spec = importlib.util.spec_from_file_location("adsb_tool", ROOT / "tools/adsb/adsb_to_geojson.py")
adsb = importlib.util.module_from_spec(spec)
spec.loader.exec_module(adsb)


class AdsBToGeoJsonTests(unittest.TestCase):
    def test_filters_invalid_and_placeholder_positions(self):
        payload = json.loads((ROOT / "sample-data/adsb-aircraft.sample.json").read_text())
        result = adsb.convert(payload)
        self.assertEqual(result["type"], "FeatureCollection")
        self.assertEqual(len(result["features"]), 2)
        self.assertEqual(result["features"][0]["geometry"]["coordinates"], [19.2, 52.1])
        self.assertTrue(all(f["properties"]["synthetic"] for f in result["features"]))

    def test_rejects_wrong_shape(self):
        with self.assertRaises(ValueError):
            adsb.convert({"aircraft": "not-a-list"})


if __name__ == "__main__":
    unittest.main()

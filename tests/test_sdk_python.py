import sys
import unittest
from pathlib import Path
from urllib.parse import parse_qs, urlparse

ROOT = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(ROOT / "sdk/python"))
from nieboguard_client import NieboGuardClient, build_url


class PythonSdkTests(unittest.TestCase):
    def test_build_url_skips_empty_values(self):
        url = build_url("https://api.nieboguard.pl", "/api/v2/feed", {"limit": 5, "region": "", "kind": None})
        parsed = urlparse(url)
        self.assertEqual(parsed.path, "/api/v2/feed")
        self.assertEqual(parse_qs(parsed.query), {"limit": ["5"]})

    def test_client_has_read_only_methods(self):
        client = NieboGuardClient()
        for name in ("state", "feed", "alarms", "incidents", "stats", "sources", "locations", "threats"):
            self.assertTrue(callable(getattr(client, name)))


if __name__ == "__main__":
    unittest.main()

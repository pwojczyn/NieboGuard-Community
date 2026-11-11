import io
import json
import sys
import unittest
import urllib.error
import urllib.parse
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(ROOT / "sdk/python"))
from nieboguard_federation import FederationClient, FederationError, normalize_base


class FakeResponse:
    def __init__(self, payload):
        self.payload = payload

    def read(self):
        return json.dumps(self.payload).encode("utf-8")

    def __enter__(self):
        return self

    def __exit__(self, *args):
        return False


class FakeOpener:
    def __init__(self):
        self.calls = []
        self.heartbeat_failed = False

    def __call__(self, request, timeout=0):
        self.calls.append(request)
        path = urllib.parse.urlparse(request.full_url).path
        body = json.loads(request.data.decode("utf-8")) if request.data else None
        if path == "/api/v2/federation/nodes/enroll":
            return FakeResponse({"node_id": "fnode_1", "organization_id": "forg_1", "access_token": "access-1", "refresh_token": "refresh-1", "protocol": "ngf/1"})
        if path == "/api/v2/federation/nodes/token":
            return FakeResponse({"access_token": "access-2", "refresh_token": "refresh-2"})
        if path == "/api/v2/federation/nodes/heartbeat":
            if not self.heartbeat_failed:
                self.heartbeat_failed = True
                raise urllib.error.HTTPError(request.full_url, 401, "Unauthorized", {}, io.BytesIO(b'{"error":{"code":"expired"}}'))
            return FakeResponse({"config_version": 1})
        if path == "/api/v2/federation/nodes/config":
            assert request.get_method() == "GET"
            return FakeResponse({"config_version": 1})
        if path == "/api/v2/federation/nodes/revoke":
            return FakeResponse({"revoked": True})
        raise urllib.error.HTTPError(request.full_url, 404, "Not found", {}, io.BytesIO(b'{"error":{"code":"not_found"}}'))


class FederationPythonSdkTests(unittest.TestCase):
    def test_enroll_refresh_heartbeat_and_revoke(self):
        opener = FakeOpener()
        client = FederationClient(base_url="https://api.example.test", opener=opener)
        self.assertIsNone(client.status()["node_id"])
        client.enroll("a" * 64, name="Python Node")
        self.assertEqual(client.status()["node_id"], "fnode_1")
        self.assertEqual(client.credentials()["refresh_token"], "refresh-1")
        self.assertEqual(opener.calls[0].headers["X-ng-protocol"], "ngf/1")
        client.heartbeat({"queue": 0})
        self.assertEqual(client.credentials()["access_token"], "access-2")
        self.assertTrue(any("/nodes/token" in request.full_url for request in opener.calls))
        client.config()
        client.revoke()
        self.assertFalse(client.status()["has_refresh_token"])

    def test_validation_and_https(self):
        client = FederationClient(base_url="https://api.example.test", opener=FakeOpener())
        with self.assertRaises(ValueError):
            client.ingest("b1", [])
        with self.assertRaises(ValueError):
            client.ingest("b1", [{"event_id": str(i)} for i in range(201)])
        with self.assertRaises(ValueError):
            normalize_base("http://example.com")
        self.assertTrue(normalize_base("http://127.0.0.1:8787"))

    def test_status_does_not_expose_tokens(self):
        client = FederationClient(base_url="https://api.example.test")
        client.set_credentials("n", "access", "refresh")
        text = json.dumps(client.status())
        self.assertNotIn('"access"', text)
        self.assertNotIn('"refresh"', text)
        self.assertTrue(client.status()["has_refresh_token"])


if __name__ == "__main__":
    unittest.main()

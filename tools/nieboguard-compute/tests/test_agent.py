import json
import tempfile
import threading
import unittest
import urllib.error
import urllib.request
from unittest import mock
from http.server import ThreadingHTTPServer
from pathlib import Path
import sys

sys.path.insert(0, str(Path(__file__).resolve().parents[1] / "src"))

from nieboguard_compute.agent import AgentState, _handler
from nieboguard_compute.config import load


class AgentTest(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        self.cfg = load(data_dir=self.tmp.name)
        self.state = AgentState(self.cfg)
        self.server = ThreadingHTTPServer(("127.0.0.1", 0), _handler(self.state))
        self.port = self.server.server_address[1]
        self.thread = threading.Thread(target=self.server.serve_forever, daemon=True)
        self.thread.start()

    def tearDown(self):
        self.server.shutdown()
        self.server.server_close()
        self.tmp.cleanup()

    def call(self, method, path, origin="https://nieboguard.pl", body=None):
        data = json.dumps(body or {}).encode() if method == "POST" else None
        headers = {"Origin": origin, "Content-Type": "application/json"} if origin else {"Content-Type": "application/json"}
        req = urllib.request.Request(f"http://127.0.0.1:{self.port}{path}", data=data, headers=headers, method=method)
        try:
            with urllib.request.urlopen(req, timeout=5) as r:
                return r.status, dict(r.headers), json.loads(r.read() or b"{}")
        except urllib.error.HTTPError as exc:
            try:
                return exc.code, dict(exc.headers), json.loads(exc.read() or b"{}")
            finally:
                exc.close()

    def test_status_and_stop_resume(self):
        code, headers, body = self.call("GET", "/status")
        self.assertEqual(code, 200)
        self.assertTrue(body["ok"])
        self.assertFalse(body["runner_running"])
        self.assertEqual(headers.get("Access-Control-Allow-Origin"), "https://nieboguard.pl")
        code, _, body = self.call("POST", "/stop")
        self.assertEqual(code, 200)
        self.assertTrue(body["stop_requested"])
        code, _, body = self.call("POST", "/resume")
        self.assertEqual(code, 200)
        self.assertFalse(body["stop_requested"])

    def test_start_without_enrollment_is_refused(self):
        code, _, body = self.call("POST", "/start")
        self.assertEqual(code, 409)
        self.assertEqual(body["error"]["code"], "not_enrolled")

    def test_foreign_origin_is_denied(self):
        code, _, body = self.call("GET", "/status", origin="https://evil.example")
        self.assertEqual(code, 403)
        self.assertEqual(body["error"]["code"], "origin_denied")
        code, _, body = self.call("POST", "/stop", origin="https://evil.example")
        self.assertEqual(code, 403)

    def test_start_and_stop_with_mocked_process(self):
        self.cfg.update({"device_id": "cmpd_test", "refresh_token": "ngr_test", "consent_version": "2026-11-01"})
        proc = mock.Mock()
        proc.pid = 424242
        with mock.patch("nieboguard_compute.agent.subprocess.Popen", return_value=proc), \
             mock.patch("nieboguard_compute.agent._pid_alive", return_value=True):
            code, body = self.state.start()
            self.assertEqual(code, 202)
            self.assertEqual(body["runner_pid"], 424242)
            self.assertTrue(self.state.stop_file.exists() is False)
            self.assertTrue(self.state.status()["runner_running"])
            code, body = self.state.stop()
            self.assertEqual(code, 200)
            self.assertTrue(self.state.stop_file.exists())
            self.state.resume()
            self.assertFalse(self.state.stop_file.exists())
        with mock.patch("nieboguard_compute.agent._pid_alive", return_value=False):
            self.assertFalse(self.state.status()["runner_running"])

    def test_preflight(self):
        req = urllib.request.Request(f"http://127.0.0.1:{self.port}/start", method="OPTIONS",
                                     headers={"Origin": "https://nieboguard.pl"})
        with urllib.request.urlopen(req, timeout=5) as r:
            self.assertEqual(r.status, 204)
            self.assertEqual(r.headers.get("Access-Control-Allow-Origin"), "https://nieboguard.pl")


if __name__ == "__main__":
    unittest.main()

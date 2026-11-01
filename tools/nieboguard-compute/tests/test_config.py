import json
import os
import stat
import tempfile
import unittest
from pathlib import Path
import sys

ROOT = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(ROOT / "src"))

from nieboguard_compute.config import ConfigError, load, redacted, save, validate


class ConfigTest(unittest.TestCase):
    def test_roundtrip_private_permissions(self):
        with tempfile.TemporaryDirectory() as tmp:
            cfg = load(data_dir=tmp)
            cfg.update({"device_id": "cmpd_1", "refresh_token": "secret", "consent_version": "2026-11-01"})
            path = save(cfg, str(Path(tmp) / "config.json"))
            self.assertEqual(load(str(path))["device_id"], "cmpd_1")
            if os.name != "nt":
                mode = stat.S_IMODE(path.stat().st_mode)
                self.assertEqual(mode & 0o077, 0)

    def test_redacted_hides_token(self):
        cfg = load(data_dir=tempfile.mkdtemp())
        cfg["refresh_token"] = "super-secret"
        self.assertEqual(redacted(cfg)["refresh_token"], "***")

    def test_rejects_insecure_api(self):
        cfg = load(data_dir=tempfile.mkdtemp())
        cfg["api_base"] = "http://example.com"
        with self.assertRaises(ConfigError):
            validate(cfg)
        cfg["api_base"] = "http://127.0.0.1:8787"
        validate(cfg)

    def test_pack_preferences(self):
        cfg = load(data_dir=tempfile.mkdtemp())
        self.assertEqual(cfg["max_data_mb"], 200)
        self.assertIn("MINI", cfg["pack_tiers"])
        cfg["pack_tiers"] = ["NOPE"]
        with self.assertRaises(ConfigError):
            validate(cfg)

    def test_rejects_partial_credentials(self):
        cfg = load(data_dir=tempfile.mkdtemp())
        cfg["device_id"] = "x"
        with self.assertRaises(ConfigError):
            validate(cfg)


class InstallerStaticTest(unittest.TestCase):
    def test_installers_do_not_pipe_remote_code_or_require_sudo(self):
        files = ["install-linux.sh", "install-windows.ps1", "start-linux.sh",
                 "start-windows.bat", "start-windows.ps1", "uninstall-linux.sh", "uninstall-windows.ps1"]
        for name in files:
            text = (ROOT / name).read_text(encoding="utf-8")
            self.assertNotIn("curl |", text, name)
            self.assertNotIn("curl|", text, name)
            self.assertNotIn("wget -O -", text, name)
            self.assertNotIn("Invoke-Expression", text, name)
            self.assertNotIn("iex(", text.lower(), name)
        for name in ("install-linux.sh", "install-windows.ps1"):
            text = (ROOT / name).read_text(encoding="utf-8")
            self.assertIn("venv", text)
            self.assertIn("doctor", text)

    def test_optional_lock_has_hashes(self):
        text = (ROOT / "requirements-optional.lock").read_text(encoding="utf-8")
        self.assertIn("--hash=sha256:", text)
        self.assertIn("numpy==", text)

    def test_core_has_no_required_third_party_dependencies(self):
        text = (ROOT / "pyproject.toml").read_text(encoding="utf-8")
        self.assertIn("dependencies = []", text)


if __name__ == "__main__":
    unittest.main()

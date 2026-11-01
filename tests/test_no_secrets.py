import re
import unittest
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
SKIP_DIRS = {".git", "__pycache__", ".venv", "node_modules"}
PATTERNS = [
    re.compile(r"-----BEGIN (?:RSA |EC |OPENSSH )?PRIVATE KEY-----"),
    re.compile(r"(?:CLOUDFLARE_API_TOKEN|ADMIN_TOKEN|INGEST_TOKEN|COMMUNITY_PUBLISH_TOKEN)\s*=\s*\S+"),
    re.compile(r"Bearer\s+[A-Za-z0-9._~+/-]{30,}"),
    re.compile(r"AKIA[0-9A-Z]{16}"),
    re.compile(r"sk-[A-Za-z0-9]{20,}"),
    re.compile(r"hf_[A-Za-z0-9]{20,}"),
]


class SecretHygieneTests(unittest.TestCase):
    def test_no_secret_values(self):
        findings = []
        for path in ROOT.rglob("*"):
            if not path.is_file() or any(part in SKIP_DIRS for part in path.parts):
                continue
            if path.stat().st_size > 1_000_000:
                continue
            try:
                text = path.read_text(encoding="utf-8", errors="ignore")
            except OSError:
                continue
            for pattern in PATTERNS:
                if pattern.search(text):
                    findings.append(str(path.relative_to(ROOT)))
                    break
        self.assertEqual(findings, [], f"possible secret values in: {findings}")

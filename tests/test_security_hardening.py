"""Static security regressions for Community installers and CI."""
from pathlib import Path
import re
import unittest

ROOT = Path(__file__).resolve().parent.parent


class InstallerSecurityTests(unittest.TestCase):
    def test_linux_installer_verifies_sha256(self):
        text = (ROOT / "installers/linux/install.sh").read_text()
        self.assertIn("SHA256SUMS", text)
        self.assertIn("sha256sum", text)
        self.assertIn("Checksum mismatch", text)
        self.assertNotIn("curl -L \"$BINARY_URL\" -o /usr/local/bin/nnode", text)

    def test_compose_does_not_publish_on_all_interfaces(self):
        text = (ROOT / "node/docker-compose.yaml").read_text()
        self.assertIn("NB_NODE_HOST=127.0.0.1", text)
        self.assertIn('"127.0.0.1:8787:8787"', text)
        self.assertNotIn("NB_NODE_HOST=0.0.0.0", text)
        self.assertNotIn('"8787:8787"', text)

    def test_workflow_actions_are_pinned(self):
        for path in (ROOT / ".github/workflows").glob("*.yml"):
            for line in path.read_text().splitlines():
                match = re.match(r"^\s*uses:\s+([^\s#]+)", line)
                if match:
                    self.assertRegex(match.group(1), r"@[0-9a-f]{40}$", f"{path}: {match.group(1)}")

    def test_release_uploads_and_downloads_artifacts(self):
        text = (ROOT / ".github/workflows/build-release.yml").read_text()
        self.assertIn("upload-artifact", text)
        self.assertIn("download-artifact", text)
        self.assertIn("permissions:\n      contents: write", text)


if __name__ == "__main__":
    unittest.main()

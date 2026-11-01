import unittest
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
FORBIDDEN_TOP = {
    ".env", "src", "frontend", "backend", "mobile", "vps-live", "vps-proxy",
    "android-app", "data", "ui-audit", "visual-render", "voice-render",
}
FORBIDDEN_SUFFIX = {".aab", ".jks", ".keystore", ".pem", ".key", ".p12"}


class RepoHygieneTests(unittest.TestCase):
    def test_no_production_or_secret_artifacts(self):
        bad = []
        for path in ROOT.rglob("*"):
            if ".git" in path.parts or not path.is_file():
                continue
            rel = path.relative_to(ROOT)
            if rel.parts[0] in FORBIDDEN_TOP:
                bad.append(str(rel))
            if path.suffix.lower() in FORBIDDEN_SUFFIX:
                bad.append(str(rel))
            if path.suffix.lower() == '.apk' and rel.parts[0] != 'releases':
                bad.append(str(rel))
        self.assertEqual(bad, [], f"forbidden files: {bad}")

    def test_only_one_approved_apk(self):
        apks = [p for p in ROOT.rglob('*.apk') if '.git' not in p.parts]
        self.assertEqual([str(p.relative_to(ROOT)) for p in apks], ['releases/nieboguard-mobile-0.1.0-dev-debug.apk'])

    def test_no_symlinks(self):
        links = [str(p.relative_to(ROOT)) for p in ROOT.rglob("*") if p.is_symlink()]
        self.assertEqual(links, [], f"symlinks are not allowed: {links}")


if __name__ == "__main__":
    unittest.main()

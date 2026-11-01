import json
import time
import unittest
from pathlib import Path
import sys

sys.path.insert(0, str(Path(__file__).resolve().parents[1] / "src"))

from nieboguard_compute.crypto import b64url_encode, canonical_json
from nieboguard_compute.manifest import ManifestError, validate_manifest


def base_manifest():
    return {
        "schema": "nieboguard-compute-task/v1",
        "task_id": "task-1",
        "round_id": "round-1",
        "task_type": "audio-classification",
        "modality": "audio",
        "issued_at": "2026-11-01T00:00:00Z",
        "expires_at": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime(time.time() + 3600)),
        "max_train_seconds": 60,
        "dataset": {
            "id": "demo", "version": 1, "license": "CC0-1.0",
            "redistribution_allowed": True, "training_allowed": True,
            "shards": [{
                "name": "labels.jsonl", "bytes": 10, "sha256": "a" * 64,
                "url": "https://api.neboguard.pl/api/v2/compute/objects/u1?token=x",
            }],
        },
        "upload": {
            "upload_id": "u2", "url": "https://api.neboguard.pl/api/v2/compute/objects/u2?token=y",
            "content_type": "application/vnd.nieboguard.safetensors", "max_bytes": 1024,
        },
        "architecture": {
            "hash": "arch",
            "tensor_specs": [
                {"name": "backbone.weight", "shape": [4, 2], "dtype": "F32"},
                {"name": "backbone.bias", "shape": [2], "dtype": "F32"},
                {"name": "classifier.weight", "shape": [2, 2], "dtype": "F32"},
                {"name": "classifier.bias", "shape": [2], "dtype": "F32"},
            ],
        },
        "classes": ["drone", "bird"],
        "base_model": {"id": "base", "version": "0.1"},
        "hyperparams": {"epochs": 1, "batch_size": 2, "max_steps": 10},
    }


class ManifestTest(unittest.TestCase):
    def test_unsigned_allowed_only_with_flag(self):
        manifest = base_manifest()
        with self.assertRaises(ManifestError):
            validate_manifest(manifest, None, "https://api.neboguard.pl")
        info = validate_manifest(manifest, None, "https://api.neboguard.pl", allow_unsigned=True)
        self.assertEqual(info["task_id"], "task-1")
        self.assertEqual(info["shard_count"], 1)

    def test_rejects_missing_redistribution(self):
        manifest = base_manifest()
        manifest["dataset"]["redistribution_allowed"] = False
        with self.assertRaises(ManifestError):
            validate_manifest(manifest, None, "https://api.neboguard.pl", allow_unsigned=True)

    def test_rejects_foreign_origin(self):
        manifest = base_manifest()
        manifest["dataset"]["shards"][0]["url"] = "https://evil.example/x"
        with self.assertRaises(ManifestError):
            validate_manifest(manifest, None, "https://api.neboguard.pl", allow_unsigned=True)
        manifest = base_manifest()
        manifest["upload"]["url"] = "https://evil.example/u"
        with self.assertRaises(ManifestError):
            validate_manifest(manifest, None, "https://api.neboguard.pl", allow_unsigned=True)

    def test_rejects_expired(self):
        manifest = base_manifest()
        manifest["expires_at"] = "2020-01-01T00:00:00Z"
        with self.assertRaises(ManifestError):
            validate_manifest(manifest, None, "https://api.neboguard.pl", allow_unsigned=True)

    def test_pack_limits_and_tier(self):
        manifest = base_manifest()
        manifest["pack"] = {"id": "cdp_test", "tier": "MINI", "bytes": 40_000_000, "sha256": "b" * 64}
        info = validate_manifest(manifest, None, "https://api.neboguard.pl", allow_unsigned=True,
                                 max_data_mb=50, allowed_tiers=["MINI"])
        self.assertEqual(info["pack"]["tier"], "MINI")
        with self.assertRaises(ManifestError):
            validate_manifest(manifest, None, "https://api.neboguard.pl", allow_unsigned=True,
                              max_data_mb=10, allowed_tiers=["MINI"])
        with self.assertRaises(ManifestError):
            validate_manifest(manifest, None, "https://api.neboguard.pl", allow_unsigned=True,
                              max_data_mb=50, allowed_tiers=["GPU"])

    def test_positive_signature(self):
        try:
            from cryptography.hazmat.primitives.asymmetric.ed25519 import Ed25519PrivateKey
            from cryptography.hazmat.primitives import serialization
        except Exception:
            self.skipTest("cryptography nie jest zainstalowana")
        key = Ed25519PrivateKey.generate()
        pub = key.public_key().public_bytes(serialization.Encoding.Raw, serialization.PublicFormat.Raw)
        manifest = base_manifest()
        signature = key.sign(canonical_json(manifest).encode())
        manifest["signature"] = {"alg": "Ed25519", "key_id": "test", "value": b64url_encode(signature)}
        info = validate_manifest(manifest, b64url_encode(pub), "https://api.neboguard.pl")
        self.assertEqual(info["task_id"], "task-1")
        manifest["task_id"] = "tampered"
        with self.assertRaises(ManifestError):
            validate_manifest(manifest, b64url_encode(pub), "https://api.neboguard.pl")


if __name__ == "__main__":
    unittest.main()

import unittest
from pathlib import Path
import sys

sys.path.insert(0, str(Path(__file__).resolve().parents[1] / "src"))

from nieboguard_compute.crypto import (b64url_encode, canonical_json, verify_ed25519,
                                       verify_manifest, SignatureError)


class CryptoTest(unittest.TestCase):
    def test_rfc8032_vector(self):
        pub = bytes.fromhex("d75a980182b10ab7d54bfed3c964073a0ee172f3daa62325af021a68f707511a")
        sig = bytes.fromhex(
            "e5564300c360ac729086e2cc806e828a84877f1eb8e5d974d873e06522490155"
            "5fb8821590a33bacc61e39701cf9b46bd25bf5f0595bbe24655141438e7a100b"
        )
        self.assertTrue(verify_ed25519(b64url_encode(pub), b"", b64url_encode(sig)))
        self.assertFalse(verify_ed25519(b64url_encode(pub), b"x", b64url_encode(sig)))
        self.assertFalse(verify_ed25519(b64url_encode(pub), b"", b64url_encode(sig[:-1] + b"\x00")))

    def test_canonical_json_matches_js_numbers(self):
        self.assertEqual(canonical_json({"b": 2, "a": [1.0, 1e-7, True, None]}),
                         '{"a":[1,1e-7,true,null],"b":2}')
        self.assertEqual(canonical_json({"x": 0.001, "y": 1.5}), '{"x":0.001,"y":1.5}')
        self.assertEqual(canonical_json([0.00001, 0.000001, 1e-7, 1e20, 1e21]),
                         "[0.00001,0.000001,1e-7,100000000000000000000,1e+21]")

    def test_manifest_requires_signature(self):
        with self.assertRaises(SignatureError):
            verify_manifest({"a": 1}, b64url_encode(b"k" * 32))

    def test_manifest_positive_with_cryptography_if_available(self):
        try:
            from cryptography.hazmat.primitives.asymmetric.ed25519 import Ed25519PrivateKey
            from cryptography.hazmat.primitives import serialization
        except Exception:
            self.skipTest("cryptography nie jest zainstalowana")
        key = Ed25519PrivateKey.generate()
        pub = key.public_key().public_bytes(serialization.Encoding.Raw, serialization.PublicFormat.Raw)
        payload = {"schema": "nieboguard-compute-task/v1", "task_id": "t", "value": 1.0, "nested": {"z": 2}}
        message = canonical_json(payload).encode()
        sig = key.sign(message)
        manifest = dict(payload, signature={"alg": "Ed25519", "key_id": "test", "value": b64url_encode(sig)})
        verify_manifest(manifest, b64url_encode(pub))
        tampered = dict(manifest, task_id="u")
        with self.assertRaises(SignatureError):
            verify_manifest(tampered, b64url_encode(pub))


if __name__ == "__main__":
    unittest.main()

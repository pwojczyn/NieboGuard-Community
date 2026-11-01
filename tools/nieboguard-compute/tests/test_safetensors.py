import math
import struct
import unittest
from pathlib import Path
import sys

sys.path.insert(0, str(Path(__file__).resolve().parents[1] / "src"))

from nieboguard_compute.safetensors import (SafetensorsError, Tensor, read_safetensors,
                                            tensor_stats, write_safetensors)


class SafetensorsTest(unittest.TestCase):
    def test_roundtrip(self):
        tensors = {
            "b": Tensor("b", [2], [1.0, -2.0]),
            "a": Tensor("a", [2, 2], [0.5, 0.25, -0.5, 0.0]),
        }
        data = write_safetensors(tensors, {"k": "v", "n": 2})
        out, meta = read_safetensors(data)
        self.assertEqual(meta, {"k": "v", "n": "2"})
        self.assertEqual(out["a"].shape, [2, 2])
        self.assertEqual(out["a"].values, [0.5, 0.25, -0.5, 0.0])
        stats = tensor_stats(out)
        self.assertEqual(stats["global"]["count"], 6)

    def test_rejects_nan(self):
        data = bytearray(write_safetensors({"a": Tensor("a", [1], [0.0])}, {}))
        # Replace the first float after the header with NaN.
        header_len = struct.unpack("<Q", bytes(data[:8]))[0]
        struct.pack_into("<f", data, 8 + header_len, math.nan)
        with self.assertRaises(SafetensorsError):
            read_safetensors(bytes(data))

    def test_rejects_trailing(self):
        data = write_safetensors({"a": Tensor("a", [1], [1.0])}, {})
        with self.assertRaises(SafetensorsError):
            read_safetensors(data + b"\x00")

    def test_rejects_shape_mismatch(self):
        with self.assertRaises(SafetensorsError):
            write_safetensors({"a": Tensor("a", [2], [1.0])}, {})


if __name__ == "__main__":
    unittest.main()

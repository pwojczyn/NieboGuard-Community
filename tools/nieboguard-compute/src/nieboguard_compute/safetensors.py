"""Minimalny, bezpieczny odczyt/zapis safetensors (tylko F32).

Nie używamy pickle, zipów ani wykonywalnych checkpointów. Nagłówek jest
walidowany, a każdy tensor ma jawny dtype i kształt.
"""
from __future__ import annotations

import json
import math
import struct
from dataclasses import dataclass
from typing import Any

MAX_BYTES = 16 * 1024 * 1024
MAX_HEADER = 1024 * 1024


class SafetensorsError(ValueError):
    pass


@dataclass(frozen=True)
class Tensor:
    name: str
    shape: list[int]
    values: list[float]
    dtype: str = "F32"


def write_safetensors(tensors: dict[str, Tensor], metadata: dict[str, str]) -> bytes:
    header: dict[str, Any] = {}
    chunks: list[bytes] = []
    offset = 0
    for name in sorted(tensors):
        tensor = tensors[name]
        if tensor.dtype != "F32":
            raise SafetensorsError("client writes only F32 tensors")
        count = 1
        for dim in tensor.shape:
            if not isinstance(dim, int) or dim < 0 or dim > 1_000_000:
                raise SafetensorsError(f"invalid shape for {name}")
            count *= dim
        if count != len(tensor.values):
            raise SafetensorsError(f"shape/value mismatch for {name}")
        data = struct.pack("<%df" % len(tensor.values), *tensor.values)
        header[name] = {"dtype": "F32", "shape": list(tensor.shape), "data_offsets": [offset, offset + len(data)]}
        offset += len(data)
        chunks.append(data)
    header["__metadata__"] = {str(k): str(v) for k, v in metadata.items()}
    header_bytes = json.dumps(header, separators=(",", ":"), ensure_ascii=False).encode("utf-8")
    if len(header_bytes) > MAX_HEADER:
        raise SafetensorsError("safetensors header is too large")
    out = bytearray(struct.pack("<Q", len(header_bytes)) + header_bytes)
    for chunk in chunks:
        out.extend(chunk)
    if len(out) > MAX_BYTES:
        raise SafetensorsError("safetensors file is too large")
    return bytes(out)


def read_safetensors(data: bytes, max_bytes: int = MAX_BYTES) -> tuple[dict[str, Tensor], dict[str, str]]:
    if not isinstance(data, (bytes, bytearray)):
        raise SafetensorsError("safetensors body must be bytes")
    if len(data) < 10:
        raise SafetensorsError("safetensors file is too small")
    if len(data) > max_bytes:
        raise SafetensorsError("safetensors file exceeds the limit")
    header_len = struct.unpack("<Q", data[:8])[0]
    if header_len < 2 or header_len > MAX_HEADER or 8 + header_len > len(data):
        raise SafetensorsError("invalid safetensors header length")
    try:
        header = json.loads(data[8:8 + header_len].decode("utf-8"))
    except Exception as exc:
        raise SafetensorsError("invalid safetensors header json") from exc
    if not isinstance(header, dict):
        raise SafetensorsError("invalid safetensors header")
    metadata = header.pop("__metadata__", {})
    if not isinstance(metadata, dict):
        raise SafetensorsError("invalid safetensors metadata")
    metadata = {str(k): str(v) for k, v in metadata.items()}
    tensors: dict[str, Tensor] = {}
    max_end = 0
    for name, value in header.items():
        if not isinstance(value, dict):
            raise SafetensorsError(f"invalid tensor descriptor {name}")
        dtype = str(value.get("dtype", ""))
        if dtype != "F32":
            raise SafetensorsError(f"client reads only F32 tensors ({name})")
        shape = value.get("shape")
        offsets = value.get("data_offsets")
        if not isinstance(shape, list) or not isinstance(offsets, list) or len(offsets) != 2:
            raise SafetensorsError(f"invalid tensor {name}")
        shape = [int(x) for x in shape]
        count = 1
        for dim in shape:
            if dim < 0 or dim > 1_000_000:
                raise SafetensorsError(f"invalid shape for {name}")
            count *= dim
        start, end = int(offsets[0]), int(offsets[1])
        if start < 0 or end < start or end - start != count * 4:
            raise SafetensorsError(f"tensor byte length mismatch for {name}")
        if 8 + header_len + end > len(data):
            raise SafetensorsError(f"truncated tensor {name}")
        values = list(struct.unpack_from("<%df" % count, data, 8 + header_len + start))
        if any(not math.isfinite(v) for v in values):
            raise SafetensorsError(f"tensor {name} contains NaN or Infinity")
        tensors[name] = Tensor(name=name, shape=shape, values=values, dtype=dtype)
        max_end = max(max_end, end)
    if not tensors:
        raise SafetensorsError("safetensors has no tensors")
    if 8 + header_len + max_end != len(data):
        raise SafetensorsError("unexpected trailing safetensors data")
    return tensors, metadata


def tensor_stats(tensors: dict[str, Tensor]) -> dict[str, Any]:
    stats = []
    global_count = 0
    global_abs = 0.0
    global_l2 = 0.0
    for name in sorted(tensors):
        t = tensors[name]
        total = sum(t.values)
        l2 = math.sqrt(sum(v * v for v in t.values))
        abs_max = max((abs(v) for v in t.values), default=0.0)
        stats.append({
            "name": name, "dtype": t.dtype, "shape": t.shape, "count": len(t.values),
            "min": min(t.values), "max": max(t.values), "mean": total / max(1, len(t.values)),
            "abs_max": abs_max, "l2": l2, "nonfinite": 0,
        })
        global_count += len(t.values)
        global_abs = max(global_abs, abs_max)
        global_l2 += l2 * l2
    return {"tensors": stats, "global": {"count": global_count, "abs_max": global_abs, "l2": math.sqrt(global_l2), "nonfinite": 0}}

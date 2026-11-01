"""Canonical JSON, base64url and Ed25519 verification.

Nie pobieramy kodu z backendu. Klient weryfikuje podpis manifestu kluczem
publicznym przypiętym w wydaniu (config lub --public-key).
"""
from __future__ import annotations

import base64
import hashlib
import json
import math
import re
from decimal import Decimal
from typing import Any


class SignatureError(ValueError):
    pass


def b64url_decode(value: str) -> bytes:
    s = str(value or "").strip().replace("-", "+").replace("_", "/")
    if not s:
        raise ValueError("empty base64url value")
    s += "=" * ((4 - len(s) % 4) % 4)
    return base64.b64decode(s, validate=True)


def b64url_encode(value: bytes) -> str:
    return base64.urlsafe_b64encode(value).rstrip(b"=").decode("ascii")


def sha256_hex(data: bytes) -> str:
    return hashlib.sha256(data).hexdigest()


def _js_number(value: float) -> str:
    """Format a finite float the way JSON.stringify does in JavaScript.

    Python's json emits 1.0 where JavaScript emits 1 and 1e-07 where JS emits
    1e-7. The backend signs canonical JSON, so the client must reproduce it.
    """
    if value == 0:
        return "0"
    if value.is_integer() and abs(value) < 1e21:
        return str(int(value))
    absolute = abs(value)
    # ECMAScript uses fixed notation for 1e-6 <= |x| < 1e21.
    if 1e-6 <= absolute < 1e21:
        return format(Decimal(repr(value)), "f")
    text = repr(value)
    if "e" in text or "E" in text:
        mantissa, exponent = re.split("[eE]", text, maxsplit=1)
        sign = ""
        if exponent.startswith(("+", "-")):
            sign, exponent = exponent[0], exponent[1:]
        exponent = exponent.lstrip("0") or "0"
        text = f"{mantissa}e{sign}{exponent}"
    return text


def canonical_json(value: Any) -> str:
    """Deterministic JSON compatible with the Worker's canonicalJson()."""
    def enc(v: Any) -> str:
        if v is None:
            return "null"
        if v is True:
            return "true"
        if v is False:
            return "false"
        if isinstance(v, str):
            return json.dumps(v, ensure_ascii=False)
        if isinstance(v, int):
            return str(v)
        if isinstance(v, float):
            if not math.isfinite(v):
                raise ValueError("non-finite number in canonical JSON")
            return _js_number(v)
        if isinstance(v, (list, tuple)):
            return "[" + ",".join(enc(x) for x in v) + "]"
        if isinstance(v, dict):
            items = []
            for key in sorted(v.keys(), key=lambda x: str(x)):
                if v[key] is None and key not in v:
                    continue
                items.append(json.dumps(str(key), ensure_ascii=False) + ":" + enc(v[key]))
            return "{" + ",".join(items) + "}"
        raise ValueError(f"unsupported canonical JSON value: {type(v).__name__}")
    return enc(value)


# ---------------------------------------------------------------------------
# Ed25519 verification (RFC 8032), pure Python. No signing on the client.
# ---------------------------------------------------------------------------

_Q = 2 ** 255 - 19
_L = 2 ** 252 + 27742317777372353535851937790883648493
_D = (-121665 * pow(121666, _Q - 2, _Q)) % _Q
_I = pow(2, (_Q - 1) // 4, _Q)


def _inv(x: int) -> int:
    return pow(x, _Q - 2, _Q)


def _recover_x(y: int, sign: int) -> int | None:
    if y >= _Q:
        return None
    x2 = (y * y - 1) * _inv(_D * y * y + 1) % _Q
    x = pow(x2, (_Q + 3) // 8, _Q)
    if (x * x - x2) % _Q != 0:
        x = x * _I % _Q
    if (x * x - x2) % _Q != 0:
        return None
    if (x & 1) != sign:
        x = _Q - x
    return x


def _decode_point(raw: bytes) -> tuple[int, int] | None:
    if len(raw) != 32:
        return None
    y = int.from_bytes(raw, "little") & ((1 << 255) - 1)
    sign = raw[31] >> 7
    x = _recover_x(y, sign)
    if x is None:
        return None
    return (x, y)


def _edwards_add(p: tuple[int, int], q: tuple[int, int]) -> tuple[int, int]:
    x1, y1 = p
    x2, y2 = q
    denom = _D * x1 * x2 * y1 * y2 % _Q
    x3 = (x1 * y2 + x2 * y1) * _inv(1 + denom) % _Q
    y3 = (y1 * y2 + x1 * x2) * _inv(1 - denom) % _Q
    return (x3, y3)


def _scalar_mult(point: tuple[int, int], scalar: int) -> tuple[int, int]:
    result = (0, 1)
    addend = point
    while scalar:
        if scalar & 1:
            result = _edwards_add(result, addend)
        addend = _edwards_add(addend, addend)
        scalar >>= 1
    return result


def _base_point() -> tuple[int, int]:
    y = 4 * _inv(5) % _Q
    x = _recover_x(y, 0)
    if x is None:
        raise SignatureError("invalid Ed25519 base point")
    return (x, y)


_G = _base_point()


def public_key_bytes(public_key_b64: str) -> bytes:
    raw = b64url_decode(public_key_b64)
    if len(raw) == 32:
        return raw
    if len(raw) == 44 and raw[:2] == b"\x30\x2a":
        # SubjectPublicKeyInfo: last 32 bytes are the raw key.
        return raw[-32:]
    if len(raw) > 32:
        return raw[-32:]
    raise SignatureError("public key must be 32 raw bytes or SPKI DER")


def verify_ed25519(public_key_b64: str, message: bytes, signature_b64: str) -> bool:
    try:
        pub = public_key_bytes(public_key_b64)
        sig = b64url_decode(signature_b64)
        if len(sig) != 64:
            return False
        r_bytes = sig[:32]
        s = int.from_bytes(sig[32:], "little")
        if s >= _L:
            return False
        a = _decode_point(pub)
        r = _decode_point(r_bytes)
        if a is None or r is None:
            return False
        h = int.from_bytes(hashlib.sha512(r_bytes + pub + message).digest(), "little") % _L
        left = _scalar_mult(_G, s)
        right = _edwards_add(r, _scalar_mult(a, h))
        return left == right
    except Exception:
        return False


def verify_manifest(manifest: dict[str, Any], public_key_b64: str) -> None:
    if not isinstance(manifest, dict):
        raise SignatureError("manifest must be an object")
    sig = manifest.get("signature")
    if not isinstance(sig, dict) or sig.get("alg") != "Ed25519" or not sig.get("value"):
        raise SignatureError("manifest has no Ed25519 signature")
    payload = {k: v for k, v in manifest.items() if k != "signature"}
    message = canonical_json(payload).encode("utf-8")
    if not verify_ed25519(public_key_b64, message, str(sig["value"])):
        raise SignatureError("manifest signature is invalid")

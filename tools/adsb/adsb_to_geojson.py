#!/usr/bin/env python3
"""Convert dump1090/readsb aircraft.json into a GeoJSON FeatureCollection.

Read-only and offline. No network, no API keys and no receiver location is
published. Aircraft with invalid coordinates are skipped.
"""
from __future__ import annotations

import argparse
import json
import math
import sys
from pathlib import Path
from typing import Any


def _number(value: Any) -> float | None:
    if isinstance(value, bool) or value is None:
        return None
    try:
        number = float(value)
    except (TypeError, ValueError):
        return None
    return number if math.isfinite(number) else None


def convert(payload: Any) -> dict[str, Any]:
    aircraft = payload.get("aircraft", []) if isinstance(payload, dict) else payload
    if not isinstance(aircraft, list):
        raise ValueError("input must be a list or an object with an 'aircraft' list")
    features: list[dict[str, Any]] = []
    for item in aircraft:
        if not isinstance(item, dict):
            continue
        lat = _number(item.get("lat"))
        lon = _number(item.get("lon"))
        if lat is None or lon is None or not (-90 <= lat <= 90) or not (-180 <= lon <= 180):
            continue
        if lat == 0 and lon == 0:
            # Common placeholder from receivers without a position fix.
            continue
        properties = {
            "hex": str(item.get("hex") or "").strip().lower() or None,
            "flight": str(item.get("flight") or "").strip() or None,
            "alt_baro": item.get("alt_baro"),
            "gs": item.get("gs"),
            "track": item.get("track"),
            "seen": item.get("seen"),
            "synthetic": bool(item.get("synthetic", False)),
        }
        features.append({
            "type": "Feature",
            "geometry": {"type": "Point", "coordinates": [lon, lat]},
            "properties": {k: v for k, v in properties.items() if v is not None},
        })
    return {"type": "FeatureCollection", "features": features}


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--input", required=True, help="aircraft.json or '-' for stdin")
    parser.add_argument("--output", default="-", help="GeoJSON path or '-' for stdout")
    args = parser.parse_args(argv)
    try:
        raw = sys.stdin.read() if args.input == "-" else Path(args.input).read_text(encoding="utf-8")
        result = convert(json.loads(raw))
    except (OSError, ValueError, json.JSONDecodeError) as exc:
        print(f"error: {exc}", file=sys.stderr)
        return 2
    text = json.dumps(result, ensure_ascii=False, indent=2) + "\n"
    if args.output == "-":
        sys.stdout.write(text)
    else:
        Path(args.output).write_text(text, encoding="utf-8")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())

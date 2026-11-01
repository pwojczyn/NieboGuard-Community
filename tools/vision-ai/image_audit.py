#!/usr/bin/env python3
"""Inspect image metadata and optionally write a copy without EXIF."""
from __future__ import annotations

import argparse
import json
from pathlib import Path

try:
    from PIL import Image
except ImportError:  # pragma: no cover - exercised only without Pillow
    Image = None


def audit(path: str | Path, strip_to: str | Path | None = None) -> dict:
    if Image is None:
        raise RuntimeError("Pillow is required: pip install -r tools/vision-ai/requirements.txt")
    file_path = Path(path)
    with Image.open(file_path) as image:
        exif = image.getexif()
        result = {
            "file": file_path.name,
            "format": image.format,
            "width": image.width,
            "height": image.height,
            "mode": image.mode,
            "has_exif": bool(len(exif)),
            "exif_fields": len(exif) if exif else 0,
        }
        if strip_to:
            target = Path(strip_to)
            clean = image.copy()
            clean.info.pop("exif", None)
            clean.save(target, format=image.format, exif=b"")
            result["stripped_to"] = target.name
        return result


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("file")
    parser.add_argument("--strip", metavar="OUTPUT", help="write a copy without EXIF")
    args = parser.parse_args(argv)
    try:
        print(json.dumps(audit(args.file, args.strip), ensure_ascii=False, indent=2))
    except (OSError, RuntimeError, ValueError) as exc:
        print(json.dumps({"error": str(exc)}))
        return 2
    return 0


if __name__ == "__main__":
    raise SystemExit(main())

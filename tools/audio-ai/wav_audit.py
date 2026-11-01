#!/usr/bin/env python3
"""Read WAV metadata without uploading or decoding the full signal."""
from __future__ import annotations

import argparse
import json
import wave
from pathlib import Path


def audit(path: str | Path) -> dict:
    file_path = Path(path)
    with wave.open(str(file_path), "rb") as wav:
        channels = wav.getnchannels()
        sample_width = wav.getsampwidth()
        sample_rate = wav.getframerate()
        frames = wav.getnframes()
    duration = frames / sample_rate if sample_rate else 0.0
    warnings = []
    if duration <= 0:
        warnings.append("empty or unreadable duration")
    if duration > 30:
        warnings.append("long sample; review consent and storage")
    if sample_rate < 16000:
        warnings.append("sample rate below 16 kHz")
    if channels > 2:
        warnings.append("more than two channels")
    return {
        "file": file_path.name,
        "bytes": file_path.stat().st_size,
        "duration_s": round(duration, 3),
        "sample_rate_hz": sample_rate,
        "channels": channels,
        "bits_per_sample": sample_width * 8,
        "frames": frames,
        "warnings": warnings,
    }


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("file", help="path to a .wav file")
    args = parser.parse_args(argv)
    try:
        print(json.dumps(audit(args.file), ensure_ascii=False, indent=2))
    except (OSError, wave.Error) as exc:
        print(json.dumps({"error": str(exc)}))
        return 2
    return 0


if __name__ == "__main__":
    raise SystemExit(main())

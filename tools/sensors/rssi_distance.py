#!/usr/bin/env python3
"""Very rough RSSI distance estimate using a log-distance path-loss model."""
from __future__ import annotations

import argparse
import json
import math


def estimate(rssi: float, tx_power: float = -59.0, path_loss: float = 2.0) -> float:
    if path_loss <= 0:
        raise ValueError("path_loss must be > 0")
    return 10 ** ((tx_power - rssi) / (10 * path_loss))


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--rssi", type=float, required=True, help="measured RSSI in dBm")
    parser.add_argument("--tx-power", type=float, default=-59.0, help="RSSI at 1 m")
    parser.add_argument("--n", type=float, default=2.0, help="path-loss exponent")
    args = parser.parse_args(argv)
    try:
        meters = estimate(args.rssi, args.tx_power, args.n)
    except ValueError as exc:
        print(json.dumps({"error": str(exc)}))
        return 2
    print(json.dumps({
        "meters": round(meters, 2),
        "rssi_dbm": args.rssi,
        "tx_power_1m_dbm": args.tx_power,
        "path_loss_exponent": args.n,
        "warning": "Approximate only; antenna, obstacles and orientation dominate.",
    }, ensure_ascii=False, indent=2))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())

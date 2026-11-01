#!/usr/bin/env bash
# Audio AI pack preparation for Linux. Requires ffmpeg; numpy optional.
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PYTHON="${PYTHON:-python3}"
WORK="${NIEBOGUARD_AUDIO_WORK:-/data/nieboguard-compute/audio-work}"
command -v "$PYTHON" >/dev/null 2>&1 || { echo "Brak Pythona 3.10+." >&2; exit 1; }
if [[ "$(uname -s)" != "Linux" ]]; then
  echo "Użyj Linux lub WSL2." >&2
  exit 1
fi
if [[ "${1:-}" == "--help" || "${1:-}" == "-h" ]]; then
  exec "$PYTHON" "$ROOT/prepare_audio_pack.py" "$@"
fi
command -v ffmpeg >/dev/null 2>&1 || { echo "Brak ffmpeg. Zainstaluj go bez wpływu na produkcję." >&2; exit 1; }
mkdir -p "$WORK"
echo "NieboGuard Audio AI data preparation"
echo "Work: $WORK"
df -h "$WORK" | tail -n 1
echo "ffmpeg: $(ffmpeg -version 2>/dev/null | head -n1)"
if "$PYTHON" -c 'import numpy' 2>/dev/null; then echo "numpy: tak (szybszy filtr jakości)"; else echo "numpy: nie (opcjonalny)"; fi
exec "$PYTHON" "$ROOT/prepare_audio_pack.py" "$@"

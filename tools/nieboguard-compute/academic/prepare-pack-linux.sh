#!/usr/bin/env bash
# Lokalne przygotowanie paczek akademickich.
# Wymaga Linuksa x86_64, Pythona 3.10+ i dużej ilości miejsca.
# Nie uruchamia kodu z archiwum i nie pobiera całych zbiorów na VPS.
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PYTHON="${PYTHON:-python3}"
WORK="${NIEBOGUARD_PACK_WORK:-/data/nieboguard-compute/work}"

command -v "$PYTHON" >/dev/null 2>&1 || { echo "Brak Pythona 3.10+." >&2; exit 1; }
"$PYTHON" - <<'PY' || exit 1
import sys
if sys.version_info < (3, 10):
    raise SystemExit("Wymagany Python 3.10+")
PY
if [[ "$(uname -s)" != "Linux" ]]; then
  echo "Ten skrypt jest przygotowany dla Linux. Na Windows użyj WSL2." >&2
  exit 1
fi
if [[ "${1:-}" == "--help" || "${1:-}" == "-h" ]]; then
  exec "$PYTHON" "$ROOT/prepare_academic_pack.py" "$@"
fi
if [[ ! -d "$WORK" ]]; then
  mkdir -p "$WORK" 2>/dev/null || { echo "Nie można utworzyć $WORK. Ustaw NIEBOGUARD_PACK_WORK." >&2; exit 1; }
fi
echo "NieboGuard academic pack preparation"
echo "Work: $WORK"
df -h "$WORK" | tail -n 1
echo "Uruchamiam: $PYTHON $ROOT/prepare_academic_pack.py $*"
exec "$PYTHON" "$ROOT/prepare_academic_pack.py" "$@"

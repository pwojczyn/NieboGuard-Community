#!/usr/bin/env bash
# Usuwa klienta. Domyślnie zachowuje konfigurację i dane (checkpointy).
# --purge usuwa również konfigurację i dane. Nie wymaga roota.
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PURGE=0
while [[ $# -gt 0 ]]; do
  case "$1" in
    --purge) PURGE=1; shift ;;
    -h|--help) echo "Użycie: ./uninstall-linux.sh [--purge]"; exit 0 ;;
    *) echo "Nieznana opcja: $1" >&2; exit 2 ;;
  esac
done
if [[ -x "$ROOT/bin/nieboguard-compute" ]]; then
  "$ROOT/bin/nieboguard-compute" --data-dir "$ROOT/data" stop >/dev/null 2>&1 || true
fi
rm -f "$ROOT/bin/nieboguard-compute"
rmdir "$ROOT/bin" 2>/dev/null || true
rm -rf "$ROOT/.venv"
if [[ "$PURGE" == "1" ]]; then
  rm -rf "$ROOT/data"
  echo "Usunięto program, konfigurację, tokeny i checkpointy."
else
  echo "Usunięto program i środowisko Python. Konfiguracja i checkpointy pozostały w $ROOT/data."
  echo "Aby usunąć wszystko: ./uninstall-linux.sh --purge"
fi

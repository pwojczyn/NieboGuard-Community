#!/usr/bin/env bash
# Uruchamia klienta w pierwszym planie. Nie tworzy usługi ani autostartu.
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
if [[ ! -x "$ROOT/bin/nieboguard-compute" ]]; then
  echo "Brak instalacji. Uruchom najpierw ./install-linux.sh" >&2
  exit 1
fi
echo "NieboGuard Compute — uruchomienie za zgodą użytkownika. Zatrzymanie: Ctrl+C."
exec "$ROOT/bin/nieboguard-compute" start "$@"

#!/usr/bin/env bash
# NieboGuard Compute — instalator dla Linux x86_64.
# Nie wymaga roota, nie pobiera skryptów z internetu i nie uruchamia kodu
# dostarczonego przez backend. Tworzy prywatne środowisko Python.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PYTHON="${PYTHON:-python3}"
INSTALL_OPTIONAL=0
SKIP_BACKEND=0
API_BASE="${NIEBOGUARD_COMPUTE_API:-https://api.neboguard.pl}"
MANIFEST_KEY="${NIEBOGUARD_COMPUTE_MANIFEST_KEY:-}"
ASSUME_YES=0

usage() {
  cat <<'EOF'
Użycie: ./install-linux.sh [opcje]
  --python SCIEZKA       interpreter Python 3.10+ (domyślnie python3)
  --with-optional        doinstaluj opcjonalne zależności z hashami
  --skip-backend-check   pomiń sprawdzenie API (np. instalacja offline)
  --api URL              adres API (domyślnie https://api.neboguard.pl)
  --manifest-key BASE64  przypnij klucz publiczny manifestu
  -y, --yes              nie pytaj o potwierdzenie rejestracji
  -h, --help             ta pomoc
EOF
}
while [[ $# -gt 0 ]]; do
  case "$1" in
    --python) PYTHON="$2"; shift 2 ;;
    --with-optional) INSTALL_OPTIONAL=1; shift ;;
    --skip-backend-check) SKIP_BACKEND=1; shift ;;
    --api) API_BASE="$2"; shift 2 ;;
    --manifest-key) MANIFEST_KEY="$2"; shift 2 ;;
    -y|--yes) ASSUME_YES=1; shift ;;
    -h|--help) usage; exit 0 ;;
    *) echo "Nieznana opcja: $1" >&2; usage; exit 2 ;;
  esac
done

say() { printf '%s\n' "$*"; }
fail() { printf 'BŁĄD: %s\n' "$*" >&2; exit 1; }

[[ "$(uname -s)" == "Linux" ]] || fail "Ten instalator obsługuje tylko Linux."
MACHINE="$(uname -m)"
case "$MACHINE" in
  x86_64|amd64) ;;
  *) say "UWAGA: architektura $MACHINE nie jest jeszcze obsługiwana przez backend." ;;
esac

command -v "$PYTHON" >/dev/null 2>&1 || fail "Brak Pythona. Zainstaluj Python 3.10+ bez uprawnień administratora lub wskaż --python."
"$PYTHON" - <<'PY' || exit 1
import sys
if sys.version_info < (3, 10):
    raise SystemExit("Wymagany Python 3.10+ (masz %s)." % sys.version.split()[0])
print("Python:", sys.version.split()[0])
PY

MEM_KB="$(awk '/MemTotal/ {print $2}' /proc/meminfo 2>/dev/null || echo 0)"
CORES="$(getconf _NPROCESSORS_ONLN 2>/dev/null || nproc 2>/dev/null || echo 1)"
if [[ -n "$MEM_KB" && "$MEM_KB" -gt 0 ]]; then
  MEM_GB=$((MEM_KB / 1024 / 1024))
  say "RAM: ${MEM_GB} GB · CPU: ${CORES} rdzeni"
  [[ "$MEM_GB" -ge 2 ]] || fail "Wymagane co najmniej 2 GB RAM (zalecane 4 GB)."
  [[ "$MEM_GB" -ge 4 ]] || say "UWAGA: mniej niż 4 GB RAM. Zadania Vision mogą zostać odrzucone."
fi
if command -v nvidia-smi >/dev/null 2>&1; then
  say "GPU: $(nvidia-smi --query-gpu=name,memory.total --format=csv,noheader,nounits 2>/dev/null | head -n1 || echo 'NVIDIA (brak szczegółów)')"
else
  say "GPU: brak nvidia-smi. Zadania GPU nie będą przydzielane; CPU nadal działa."
fi

VENV="$ROOT/.venv"
if [[ ! -x "$VENV/bin/python" ]]; then
  say "Tworzę izolowane środowisko Python w $VENV"
  "$PYTHON" -m venv "$VENV"
fi
"$VENV/bin/python" -m pip --version >/dev/null 2>&1 || fail "venv nie zawiera pip. Zainstaluj python3-venv."
"$VENV/bin/python" -m pip install --disable-pip-version-check --quiet --upgrade "pip>=24,<26"
if [[ "$INSTALL_OPTIONAL" == "1" ]]; then
  [[ -f "$ROOT/requirements-optional.lock" ]] || fail "Brak pliku requirements-optional.lock."
  say "Instaluję opcjonalne zależności z weryfikacją SHA-256…"
  "$VENV/bin/python" -m pip install --disable-pip-version-check --require-hashes -r "$ROOT/requirements-optional.lock"
fi

mkdir -p "$ROOT/bin"
cat > "$ROOT/bin/nieboguard-compute" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
export PYTHONPATH="$ROOT/src${PYTHONPATH:+:$PYTHONPATH}"
exec "$ROOT/.venv/bin/python" -m nieboguard_compute "$@"
EOF
chmod +x "$ROOT/bin/nieboguard-compute"

CONFIG_DIR="$ROOT/data"
mkdir -p "$CONFIG_DIR"
chmod 700 "$CONFIG_DIR" 2>/dev/null || true
export NIEBOGUARD_COMPUTE_HOME="$CONFIG_DIR"
CFG="$CONFIG_DIR/config.json"
if [[ ! -f "$CFG" ]]; then
  "$ROOT/bin/nieboguard-compute" --data-dir "$CONFIG_DIR" config set api_base "$API_BASE" >/dev/null
  if [[ -n "$MANIFEST_KEY" ]]; then
    "$ROOT/bin/nieboguard-compute" --data-dir "$CONFIG_DIR" config set manifest_public_key "$MANIFEST_KEY" >/dev/null
  fi
fi

if [[ -t 0 && "$ASSUME_YES" != "1" && "${NIEBOGUARD_COMPUTE_NONINTERACTIVE:-0}" != "1" ]]; then
  say ""
  say "Wybór sposobu użycia komputera (możesz zmienić później: bin/nieboguard-compute config set ...)."
  read -r -p "Model [audio/vision/both] (domyślnie audio): " MODALITY || MODALITY=""
  case "${MODALITY:-audio}" in
    audio|vision|both) ;;
    *) MODALITY="audio" ;;
  esac
  if [[ "$MODALITY" == "both" ]]; then MODS="audio,vision"; else MODS="$MODALITY"; fi
  read -r -p "Limit CPU % (5-100, domyślnie 50): " CPU_PCT || CPU_PCT=""
  CPU_PCT="${CPU_PCT:-50}"
  read -r -p "Limit GPU % (0-100, domyślnie 60): " GPU_PCT || GPU_PCT=""
  GPU_PCT="${GPU_PCT:-60}"
  read -r -p "Harmonogram [manual/idle/always] (domyślnie manual): " SCHED || SCHED=""
  case "${SCHED:-manual}" in manual|idle|always) ;; *) SCHED="manual" ;; esac
  read -r -p "Maks. paczka danych MB [50/200/500/1024] (domyślnie 200): " DATA_MB || DATA_MB=""
  DATA_MB="${DATA_MB:-200}"
  case "$DATA_MB" in 50) TIERS="MINI" ;; 500) TIERS="MINI,STANDARD,EXTENDED" ;; 1024) TIERS="GPU" ;; *) DATA_MB="200"; TIERS="MINI,STANDARD" ;; esac
  "$ROOT/bin/nieboguard-compute" --data-dir "$CONFIG_DIR" config set modalities "$MODS" >/dev/null || true
  "$ROOT/bin/nieboguard-compute" --data-dir "$CONFIG_DIR" config set max_data_mb "$DATA_MB" >/dev/null || true
  "$ROOT/bin/nieboguard-compute" --data-dir "$CONFIG_DIR" config set pack_tiers "$TIERS" >/dev/null || true
  "$ROOT/bin/nieboguard-compute" --data-dir "$CONFIG_DIR" config set max_cpu_percent "$CPU_PCT" >/dev/null || true
  "$ROOT/bin/nieboguard-compute" --data-dir "$CONFIG_DIR" config set max_gpu_percent "$GPU_PCT" >/dev/null || true
  "$ROOT/bin/nieboguard-compute" --data-dir "$CONFIG_DIR" config set schedule "${SCHED:-manual}" >/dev/null || true
fi

if [[ "$SKIP_BACKEND" != "1" ]]; then
  if ! "$ROOT/bin/nieboguard-compute" --data-dir "$CONFIG_DIR" doctor --json >/tmp/ng-compute-doctor.$$ 2>/dev/null; then
    say "UWAGA: backend jest niedostępny lub nieobsługiwany. Instalacja lokalna jest gotowa."
    say "Sprawdź później: bin/nieboguard-compute doctor"
  else
    say "Backend odpowiada."
  fi
  rm -f /tmp/ng-compute-doctor.$$
fi

say ""
say "Instalacja zakończona: $ROOT"
say "Uruchamianie: $ROOT/start-linux.sh"
say "Diagnostyka: $ROOT/bin/nieboguard-compute doctor"
say "Rejestracja:  $ROOT/bin/nieboguard-compute enroll --invite <kod> --consent"
say ""
say "Trening nie uruchamia się w tle. Program startuje wyłącznie po Twojej komendzie."

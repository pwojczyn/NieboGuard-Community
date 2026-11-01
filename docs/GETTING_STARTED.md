# Uruchomienie NieboGuard Sensor Hub

## Szybki start (Linux / macOS / Windows)

```bash
# 1. Zbuduj binarny
cd node
go build -o nnode ./cmd/nnode

# 2. Stwórz katalog danych
mkdir -p ~/.nieboguard

# 3. Uruchom serwer (domyślnie localhost:8787)
./nnode

# 4. Otwórz panel
# http://localhost:8787
```

## Opcje wiersza poleceń

```bash
./nnode --config=config.yaml --port=8787 --host=127.0.0.1 --verbose
```

## Konfiguracja YAML (`node/web/config.yaml`)

```yaml
port: 8787
host: 127.0.0.1
workdir: ~/.nieboguard
verbose: false
sources:
  local_adsb:
    enabled: true
    url: "http://localhost:30003/aircraft.json"
  http_json:
    enabled: false
    url: "https://example.com/data.json"
```

## Windows

```powershell
# Uruchom z katalogu repozytorium
node\nnode.exe --config=node\web\config.yaml
```

## Docker / Podman

```bash
docker-compose up -d
# Panel: http://localhost:8787
```

## Installery

- Linux: `installers/linux/install.sh`
- Windows: `installers/windows/install.ps1`

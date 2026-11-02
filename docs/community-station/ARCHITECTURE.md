# NieboGuard Community Station — Architecture

## Przegląd

NieboGuard Community Station to zintegrowana aplikacja Windows 10/11 x64, napisana w Go, łącząca:

- Backend Go (Node engine + Simulator engine)
- Frontend HTML5/CSS3/JavaScript z osadzonym interfejsem
- Lokalną bazę SQLite
- Obsługę fizycznego odbiornika RTL-SDR USB
- Dekoder ADS-B (współpraca z dump1090/readsb)
- Symulator z trybami synthetic/replay/live_transform

## Architektura warstwowa

### Warstwa sprzętowa
- RTL-SDR USB (RTL2832U + R820T/R820T2)
- USB 2.0 z WinUSB/libusb
- Antena 1090 MHz

### Warstwa dekoderów
- dump1090 (Beast TCP/AVR/JSON HTTP)
- readsb
- Librtlsdr (opcjonalnie)

### Warstwa backend (Go)
- **core**: Silnik główny aplikacji
- **collectors**: Menedżer adapterów (RTL-SDR, HTTP JSON, RSS)
- **simulator**: Silnik symulacji (synthetic/replay/live_transform)
- **rtl**: Integracja z odbiornikiem USB
- **storage**: SQLite z WAL
- **api**: REST API v1
- **websocket**: Delta updates z sekwencjami

### Warstwa interfejsu
- **panel**: HTML/CSS/JS z PWA features
- **mapa**: Interaktywna z oznaczeniami REAL/SIM
- **protocol lab**: Edukacyjne moduły (ADS-B, ASTERIX, Military)

## Komunikacja komponentów

```
RTL-SDR USB → usb → librtlsdr → dekoder (dump1090) → adapter → collector → normalizacja → SQLite
                                                                                ↓
                                                                         REST/WebSocket
                                                                                ↓
                                                                         Panel HTML ← otwiera przeglądarka
```

## Format danych wejściowych

- aircraft.json (dump1090/readsb)
- Beast TCP (port 30003)
- AVR tekstowy
- SBS-1 (BaseStation)
- HTTP JSON polling

Model `AircraftRecord`:
```
{
  Hex: string,
  Flight: string,
  Lat: float64,
  Lon: float64,
  AltBaro: int,
  GS: float64,
  Track: float64,
  Seen: float64
}
```

## Schemat bazy SQLite

Tabela `aircraft_events`:
- id TEXT PRIMARY KEY
- hex TEXT NOT NULL
- flight TEXT
- lat REAL, lon REAL
- alt_baro INTEGER
- gs REAL, track REAL
- seen REAL
- source_id TEXT
- source_type TEXT
- observed_at DATETIME
- received_at DATETIME DEFAULT CURRENT_TIMESTAMP
- verified INTEGER DEFAULT 0
- synthetic INTEGER DEFAULT 0

Indexy: hex, source_id, observed_at

Tabela `sources`: id, type, name, status, last_update, license, count_records

Tabela `configuration`: key TEXT PRIMARY KEY, value BLOB

Tabela `notices`: id, source, category, title, region, verification, published_at

## Bezpieczeństwo

- Bind domyślnie 127.0.0.1
- Brak dostępu z internetu bez konfiguracji
- Tokeny przechowywane lokalnie, chronione przez Windows
- Brak automatycznego podnoszenia uprawnień
- Brak automatycznej instalacji sterowników
- CSRF i rebinding zabezpieczenia
- Ochrona przed SSRF (blokowanie prywatnych LAN poza localhost)

## Windows build

- `go build -ldflags "-s -w" -o NieboGuard-Station.exe ./cmd/community-station/`
- Zależność zewnętrzna: WinUSB (dla RTL-SDR)
- Biblioteka DLL (opcjonalnie): librtlsdr.dll
- Pakiet: EXE + dekoder + konfiguracja YAML + dokumentacja

## Konfiguracja (YAML)

```yaml
port: 8787
host: 127.0.0.1
workdir: "%LOCALAPPDATA%\\NieboGuard\\Station"
rtl_sdr:
  enabled: true
  device_index: 0
  frequency: 1090000000
  gain_mode: auto
simulator:
  enabled: false
  mode: synthetic
  objects: 10
data_sharing:
  enabled: false
  url: ""
  token: ""
  station_id: ""
```

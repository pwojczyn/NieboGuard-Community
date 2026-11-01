# IMPLEMENTATION AUDIT - NieboGuard Sensor Hub

## 1. Aktualna struktura repozytorium

### Istniejące pliki w katalogu `node/`

| Plik | Opis | Status |
|------|------|--------|
| `cmd/nnode/main.go` | Główne wejście aplikacji Go (CLI flags: config, port, host, workdir, verbose) | ✓ Zaimplementowane |
| `internal/core/core.go` | Silnik core - inicjalizacja komponentów, Start/Stop | ✓ Szkielet |
| `internal/storage/storage.go` | SQLite + WAL lokalna baza (tabela `aircraft_events`) | ✓ Szkielet |
| `internal/collectors/collectors.go` | Menedżer adapterów źródeł danych | ✓ Szkielet (3 stuby) |
| `internal/api/api.go` | Serwer HTTP REST (endpointy /api/v1/*) | ✓ Szkielet |
| `internal/websocket/websocket.go` | Serwer WebSocket (delta updates, heartbeat) | ✓ Szkielet |
| `pkg/config/config.go` | Konfiguracja YAML + zmienne środowiskowe | ✓ Zaimplementowane |
| `pkg/models/models.go` | Modele danych (Aircraft, Notice, Source, HealthInfo, APIResponse) | ✓ Zaimplementowane |
| `web/index.html` | Panel WWW (responsywny, PL/EN, lokalny) | ✓ Zaimplementowane (podstawowy) |
| `web/config.yaml` | Przykładowa konfiguracja panelu | ✓ Zaimplementowane |
| `Dockerfile` | Multi-stage Go build (alpine) | ✓ Zaimplementowane |
| `docker-compose.yaml` | Docker Compose z healthcheck | ✓ Zaimplementowane |
| `configs/nieboguard-node.container` | Konfiguracja Podman Quadlet | ✓ Zaimplementowane |
| `go.mod` / `go.sum` | Zależności Go | ✓ Zaimplementowane |

### Zależności Go

| Biblioteka | Licencja | Cel | Status |
|------------|----------|-----|--------|
| `github.com/gorilla/mux` | BSD-3-Clause | Router HTTP | ✓ Zainstalowane |
| `github.com/mattn/go-sqlite3` | MIT | Driver SQLite | ✓ Zainstalowane |
| `gopkg.in/yaml.v3` | Apache-2.0 / MIT | Parsowanie YAML | ✓ Zainstalowane |

### Brak bibliotek

- Brak biblioteki do obsługi Beast TCP / AVR / SBS-1 (ADS-B)
- Brak biblioteki do parsowania danych z dump1090 / readsb
- Brak biblioteki OpenDroneID, FLARM/OGN, ASTERIX
- Brak biblioteki MQTT
- Brak biblioteki OpenSky / ADS-B.fi (REST)

## 2. Funkcje już istniejące (P0 szkielet)

| Komponent | Funkcja | Status |
|-----------|---------|--------|
| Core Engine | Start/Stop, signal handling, konfiguracja | ✓ Działa |
| Storage | SQLite WAL, tabela aircraft_events | ✓ Działa (bez danych) |
| Collectors | Menedżer z 3 stambami (local_adsb, http_json, rss_atom) | ⚠ Szkielet |
| API | /api/v1/health, /ready, /version, /sources, /objects, /events, /notices, /stats | ✓ Działa |
| API Admin | /api/v1/admin/sources (POST) | ✓ Działa |
| WebSocket | /api/v1/stream (snapshot, ping) | ✓ Szkielet |
| Panel | HTML/JS responsive, lokalny localhost:8787 | ✓ Działa |
| Konfiguracja | YAML + env vars | ✓ Działa |
| Kontenery | Dockerfile, docker-compose.yaml, Podman Quadlet | ✓ Zaimplementowane |
| Instalatory | installers/linux/install.sh, installers/windows/install.ps1 | ✓ Zaimplementowane |
| CI/CD | .github/workflows/build-release.yml | ✓ Zaimplementowane |

## 3. Funkcje wymagające poprawek

| Komponent | Problem | Wymagana poprawa |
|-----------|---------|------------------|
| `internal/core/core.go` | Nie łączy się z storage/collectors/api - tylko szkielet | Integracja wszystkich komponentów |
| `internal/api/api.go` | Nie używa storage - zwraca dane sztywne | Podłączenie do SQLite |
| `internal/collectors/collectors.go` | Wszystkie adaptery to stuby (brak logiki) | Implementacja ADS-B, HTTP, RSS |
| `internal/websocket/websocket.go` | Brak delta updates, brak seq numerek, brak wykrywania utraty ciągłości | Pełna implementacja protokołu |
| `web/index.html` | Brak integracji z WebSocket (tylko statyczne API) | Połączenie WSS z panelem |
| `web/index.html` | Brak formularza dodawania źródeł | UI dla zarządzania adapterami |
| Brak uwierzytelniania | Brak konta administratora, haseł, sesji | Dodanie bcrypt + session |
| Brak rate limitingu | Brak limitu żądań API | Middleware rate limiting |

## 4. Funkcje brakujące

| Komponent | Opis | Priorytet |
|-----------|------|-----------|
| ADS-B Beast TCP | Odczyt Beast TCP z dump1090 / readsb | P0 |
| ADS-B AVR | Odczyt AVR tekstowy z dump1090 | P0 |
| ADS-B SBS-1 | Parsowanie formatu BaseStation | P0 |
| ADS-B JSON HTTP | Odczyt aircraft.json przez HTTP | P0 |
| HTTP Polling | ETag, If-Modified-Since, jitter, retry, backoff | P0 |
| WebSocket Client | Połączenie z zewnętrznymi źródłami WS/WSS | P0 |
| OpenDroneID | Obsługa drone remote ID | P1 |
| FLARM/OGN | Obsługa szybowców | P1 |
| Uwierzytelnianie | Konta administratora, bcrypt, sesje | P0 |
| Rate limiting | Limity żądań per IP/source | P0 |
| SSRF protection | Walidacja URL dla HTTP adapterów | P0 |
| TTL czyszczenie | Usuwanie nieaktualnych obiektów z bazy | P0 |

## 5. Ryzyka

| Ryzyko | Opis | Mitigacja |
|--------|------|-----------|
| Licencja | Produkcja na VPS jest prywatna - nie kopiować kodu | Własna implementacja w Go |
| Wersje Straznika >1.7.54 | Restrykcyjna licencja | Tylko MIT (≤1.7.53) jako referencia |
| Brak dostępu do produção | Nie można zweryfikować kontraktu WSS | Własny kontrakt API v1 |
| Brak bibliotek ADS-B | Brak gotowych bibliotek Beast/AVR/SBS | Własne parsowanie |

## 6. Decyzje architektoniczne

1. **Język**: Go 1.26 (pełny kod, statyczny binarny)
2. **Baza**: SQLite z WAL (lokalna, bez serwera)
3. **Panel**: HTML/JS bez zależności (lokalny, offline)
4. **Adapters**: Wzorzec strategii - każdy adapter implementuje interfejs `Collector`
5. **Port domyślny**: 8787 (lokalny, nie 8788 jak produkcja)
6. **Konfiguracja**: YAML + zmienne środowiskowe
7. **Lokalizacja danych**: ~/.nieboguard (Linux), %USERPROFILE%\NieboGuard (Windows)

## 7. Priorytety (P0 MVP)

1. Pełna implementacja storage (dodawanie obiektów, paginacja)
2. Adapter ADS-B (Beast TCP + JSON HTTP)
3. Adapter HTTP polling (ETag, retry, backoff, SSRF protection)
4. Uwierzytelnianie administratora
5. Panel WWW z WebSocket (delta updates)
6. Testy jednostkowe
7. Buildy Windows/Linux/ARM64

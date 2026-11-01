# AUDYT ISTNIEJĄCEJ INFRASTRUKTURY - UWAGA: Dostęp do produkcji

Data: 2026-11-01
Status: CRITICAL FINDINGS - Existing production infrastructure found on VPS

## 1. Istniejący produkcyjny system NieboGuard

### 1.1 Lokalizacja
- **Główne repozytorium prywatne**: `/home/ubuntu/NieboGuard/` (pełny kod projektu)
- **Binarny production**: `/usr/local/bin/nieboguard-live` (7.3 MB, Oct 31 19:43)
- **Konfiguracja**: `/etc/nieboguard/live.env`
- **Dane**: `/var/lib/nieboguard/`
- **Serwisy systemd**: 5 usług (live, proxy, visual, voice, cron)

### 1.2 Architektura produkcyjna (odczytywana z binarza i configu)

#### 1.2.1 Główne usługi
1. **nieboguard-live** - Go WebSocket hub (główny backend)
   - Nasłuchiwanie: `127.0.0.1:8788` (port 8788, nie 8787!)
   - User: `nieboguard`, Group: `nieboguard`
   - Max clients: 20000
   - Persistent state: `/var/lib/nieboguard/live-state.json`

2. **NieboGuard vps-proxy** - Node.js upstream proxy
   - Lokalizacja: `/home/ubuntu/NieboGuard/vps-proxy/`
   - Exec: `/usr/bin/node /home/ubuntu/NieboGuard/vps-proxy/proxy.mjs`
   - Plik .env: `/home/ubuntu/NieboGuard/vps-proxy/.env`

3. **Wizualizacja** - Python + server.py
   - Working directory: `/home/ubuntu/NieboGuard/visual-render/`
   - Python venv: `~/NieboGuard/visual-render/venv/`

4. **Głos** - Python + Piper + FFmpeg
   - Working directory: `/home/ubuntu/NieboGuard/voice-render/`

5. **Cron safety net** - Checker co minutę
   - `/usr/local/sbin/nieboguard-cron-check`

### 1.3 Konfiguracja środowiskowa (/etc/nieboguard/live.env)

Kluczowe zmienne:
- `LISTEN_ADDR=127.0.0.1:8788` - port nasłuchiwania
- `INGEST_TOKEN` - token autoryzacji (32-bajtowy losowy)
- `INGEST_HMAC_SECRET` - klucz HMAC do podpisywania aktualizacji
- `STATE_FILE=/var/lib/nieboguard/live-state.json` - stan persistowany
- `NEPTUN_ENABLED=true/false` - czy włączyć NEPTUN
- `NEPTUN_URL=https://neptun.in.ua/api/v1/threats` - adres NEPTUN
- `ADS_B_ENABLED=true` - czy ADS-B
- `MILITARY_URL=https://opendata.adsb.fi/api/v3/lat/52.2/lon/19.0/dist/250` - military ADS-B
- `ADS_FI_INTERVAL=10s` - interwał pobierania ADS-B
- `GEO_SUBSCRIPTIONS_ENABLED=false` - geo subskrypcje

### 1.4 Architektura danych (z main.go)

Główny serwer obsługuje:
- **WebSocket**: `/ws` (z parametrami `probe`, `geo`)
- **HTTP API**: `/api/update-targets` (tylko dla Workerów, autoryzacja Bearer + HMAC)
- **Health**: `/healthz`, `/health`
- **Frame types**: military, adsb, sky, weather, neptun, snapshot

**Adaptory danych** (z importami w main.go):
- `nieboguard/live/adsb` - fusion wielu źródeł (adsbfi, flyitaly, opensky)
- `nieboguard/live/sky` - satelity i orbitale
- `nieboguard/live/weather` - pogoda IMGW
- `nieboguard/live/neptun` - zagrożenia (ale w kodzie: `neptunLoop` pobiera z `neptunURL`)

### 1.5 Licencja i ograniczenia

⚠️ **WAŻNE**: 
- Produkcyjny kod w `/home/ubuntu/NieboGuard/` i `/usr/local/bin/nieboguard-live` jest **prywatny**
- Nie wolno kopiować kodu źródłowego bezpośrednio
- Licencja wersji po 1.7.54 jest restrykcyjna (tylko do 1.7.53 na MIT)
- Musi zostać zaprojektowana **niezależna implementacja** z zachowaniem kompatybilności protokołów

### 1.6 Protokół API (z main.go)

**WebSocket** (`/ws`):
- AcceptOptions z OriginPatterns
- CompressionMode: `websocket.CompressionContextTakeover`
- Message types: text (JSON)
- Control: ping/pong co 25 sekund
- Slow clients: disconnectamiamiami (bounded queue 16)
- Initial snapshot: latestSnapshot + military + adsb + sky + weather
- Messages: resync, subscribe (z bbox/zoom)

**HTTP API** (`/api/update-targets`):
- Tylko POST
- Autoryzacja: `Authorization: Bearer <token>` LUB HMAC
  - `X-NG-Timestamp`: timestamp
  - `X-NG-Signature`: sha256(timestamp + "\n" + body)
- Body: `updateEnvelope` z typem `snapshot` lub `patch`
- Seq (sequence number) - ordering
- Stan: `state` (usuwa `civilAircraft`), `patch` (ops)
- Odpowiedź: `{"ok":true}`

**Health** (`/health`):
- Klientom WSS: `{"type":"health",...}`
- Public endpoint: `{"ok":true,"clients":N,"seq":S,"state_bytes":M,...}`
- Pola: clients, seq, state_bytes, military_bytes, military_at, military_error, neptun_at, neptun_error, adsb_at, adsb_error, adsb JSON, sky_at, sky_error, weather_at, weather_error

### 1.7 Adapters (źródła danych) w produkcji

| Adapter | Źródło | Interwał | Uwagi |
|---------|--------|----------|-------|
| ADS-B fusion | adsbfi, flyitaly, opensky | 10s/30s/60s | Multiple sources, scoring w Workerze |
| Military | adsb.fi military | 30s | Tylko loty wojskowe |
| NEPTUN | neptun.in.ua | 5s | Zagrożenia (dane orbitalne?) |
| Sky | CelesTrak/satnogs | w pętli | Satelity, Starlink, orbits |
| Weather | IMGW public | 10m | Warningi meteo i hydro |
| Military fallback | ops... | zdefiniowane | adsblol jako last-resort |

### 1.8 Bezpieczeństwo produkcyjne

- `ProtectSystem=strict`, `ProtectHome=true`, `PrivateTmp=true`
- `NoNewPrivileges=true`
- `ReadWritePaths=/var/lib/nieboguard` (tylko to!)
- `LimitNOFILE=200000`
- `GOMEMLIMIT=700MiB`
- Origin filtering: `nieboguard.pl,www.nieboguard.pl,...`
- HMAC weryfikacja każdej aktualizacji
- Token w nagłówku Authorization

### 1.9 Istniejące porty i endpointy

| Endpoint | Typ | Dostęp |
|----------|-----|--------|
| `/ws` | WebSocket | Publiczny (z origin filtering) |
| `/api/update-targets` | HTTP POST | Tylko autoryzowany (Worker) |
| `/health` | HTTP GET | Publiczny (z CORS `*`) |
| `/healthz` | HTTP GET | Internal |
| Metryki: clients, seq, state_bytes itd. | JSON | Health endpoint |

## 2. Wpływ na NieboGuard Node

### 2.1 Co mogę wykorzystać (projekty open source / MIT ≤1.7.53)
- Koncepcje architekturalne z audytu
- Schematy formatów danych (aircraft, notices itp.)
- Idee dotyczące adapterów źródeł
- Format OpenAPI (istniejący w repozytorium Community)

### 2.2 Co NIE MOGĘ być kopiowane (kod prywatny)
- **Main.go** i cała logika biznesowa `nieboguard/live/...`
- **Plik live.env** z tokenami i konfiguracją produkcyjną
- **Serwisy systemd** i skrypty instalacyjne VPS
- **Kod Node.js proxy** w vps-proxy/
- **Binarny** `/usr/local/bin/nieboguard-live`
- **Dane użytkowników** i konfigurację `/var/lib/nieboguard/`

### 2.3 Decyzje architektoniczne (zaktualizowane)

1. **Port**: Moja implementacja używa `8787`, produkcyjne `8788`. Muszę zdecydować:
   - Czy zachować `8787` (domyślne w moim kodzie), czy dopasować do `8788`
   - **Decyzja**: Zachować `8787` jako domyślny dla Node (łatwiejsze dla użytkowników), ale dokumentować, że produkcyjne to `8788`

2. **Autoryzacja**:
   - Produkcja: Token Bearer + HMAC w nagłówkach
   - Moja implementacja: Prosty admin account w panelu WWW (P0)
   - **Decyzja**: Utrzymać prostą autentykację dla P0, dokumentować, że producja używa HMAC

3. **Protokół WebSocket**:
   - Produkcja: `github.com/coder/websocket`, ramki typu military/adsb/sky/weather/snapshot
   - Moja implementacja: `github.com/gorilla/websocket`, JSON messages
   - **Decyzja**: Zachować własną implementację gorilla/websocket, ale zaprojektować interfejs, który byłby kompatybilny (delta updates, sequence numbers, heartbeat)

4. **Stan i persistency**:
   - Produkcja: JSON persistowany do pliku co aktualizację
   - Moja implementacja: SQLite z WAL
   - **Decyzja**: SQLite jest lepsze dla pojedynczego Node, produkcja używa plików JSON - to różnica projektowa, nie naruszenie

5. **Geo subscriptions** (filtr mapy po bbox):
   - Produkcja: Zaawansowane filtry po geographicznym bbox na kliencie
   - Moja implementacja: Podstawowa lista źródeł w panelu
   - **Decyzja**: P0 bez geo subscriptions, funkcjonalność opcjonalna w P1

6. **HMAC podpisy aktualizacji**:
   - Produkcja: `X-NG-Timestamp` + `X-NG-Signature` do `/api/update-targets`
   - Moja implementacja: Brak (panel lokalny, brak API zapisu)
   - **Decyzja**: Pomijane w P0 (brak endpointu zapisu), dokumentowane jako funkcja P1

## 3. Zaktualizowana struktura adapters (porównanie)

| Adapter | Produkcja (prywatny kod) | NieboGuard Node (mój kod) | Status |
|---------|--------------------------|---------------------------|--------|
| ADS-B | `nieboguard/live/adsb` z fusionem 3 źródeł | `localADSbSource` - stub | Inne źródła, własna implementacja |
| NEPTUN | `neptunURL` + etag/if-none-match | `rssAtomSource` / placeholder | Różne API |
| Military | adsb.fi military flag + normalized struct | HTTP JSON polling adapter | Inny protokół |
| Sky | CelesTrak + satnogs + orbits | Brak (P0) | Inne dane |
| Weather | IMGW public API | Brak (P0) | Inne dane |
| Community Sensor | Brak w produkcji | `sensor_detection` model | Nowość |

## 4. Zaktualizowany plan (co zmienić)

### Muszę zmienić w dokumentacji:
1. `ARCHITECTURE.md` - dopisać sekcję o istniejącym production na VPS
2. `DECISIONS.md` - zaktualizować decyzje oparte na nowych informacjach
3. `IMPLEMENTATION_REPORT.md` - dodać sekcję o istniejącym systemie
4. `TASKS.md` - dodać zadania związane z kompatibilnością

### Co zostawić bez zmian (już dobrze zaprojektowane):
- SQLite z WAL - lepsze dla pojedynczego Node niż pliki JSON produkcji
- Panel WWW localhost-only - zgodne z modelem LOCAL
- Tryby LOCAL/COMMUNITY/PUBLIC NODE - niezależne od produkcji
- OpenAPI 3.1 - dokumentacja publiczna

## 5. Kolejne kroki (z uwzględnieniem produkcji)

1. **Nie kopiować kodu** z `/home/ubuntu/NieboGuard/` ani `/usr/local/bin/nieboguard-live`
2. **Zaprojektować interfejsy kompatybilne** (nie kopiować implementacji)
3. **Zastosować lekkie modyfikacje** w dokumentacji (port 8787 vs 8788, prostsza auth)
4. **Dokumentować różnice** między moim Node a produkcją
5. **Testować** czy Node może współpracować z istniejącym WS hub (teoretycznie tak, oba działają na localhost)

**Kluczowa zasada**: NieboGuard Node ma działać samodzielnie. Integracja z centralą (produkcją) jest opcjonalna (P1) i projektowana jako oddzielny moduł z własnym kontraktem protokołu.
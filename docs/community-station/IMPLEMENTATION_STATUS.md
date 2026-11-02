# IMPLEMENTATION STATUS — NieboGuard Community Station

Data audytu: 2026-11-02
Wersja: 0.1.0-audit

## Tabela stanu IMPLEMENTED / PARTIAL / STUB / MISSING / BROKEN

| Komponent | Plik | Stan | Uwagi |
|-----------|------|------|-------|
| Core Engine | node/internal/core/core.go | PARTIAL | Start inicjalizuje komponenty, ale Stop() nie zamyka serwerów ani collectorów |
| Storage | node/internal/storage/storage.go | IMPLEMENTED | SQLite WAL, tabele aircraft_events/notices/sources/configuration |
| Collectors Manager | node/internal/collectors/collectors.go | STUB | localADSbSource, httpJSONSource, rssAtomSource — wszystkie to atrapy |
| ADSBJSONAdapter | node/internal/collectors/adapter_adsb_json.go | IMPLEMENTED | Parsuje aircraft.json z dump1090/readsb, ale nie jest podłączony do menedżera |
| HTTPPollAdapter | node/internal/collectors/adapter_http_json.go | PARTIAL | SSRF protection niekompletne, brak ETag/retry/backoff |
| API Server | node/internal/api/api.go | IMPLEMENTED | Endpointy /api/v1/* działają, ale ObjectsHandler zwraca dane z bazy (puste) |
| WebSocket | node/internal/websocket/websocket.go | STUB | Brak delta updates, seq numeracja, snapshot z danych |
| Config | node/pkg/config/config.go | IMPLEMENTED | YAML, domyślne port 8787, host 127.0.0.1 |
| Models | node/pkg/models/models.go | IMPLEMENTED | Aircraft, Notice, Source, HealthInfo, APIResponse |
| Web Panel | node/web/index.html | PARTIAL | Podstawowy interfejs PL/EN, brak integracji z API, brak RTL-SDR, brak symulatora |
| Simulator Engine | simulator/internal/sim/engine.go | PARTIAL | syntheticFrame dziala, replayFrame/liveFrame zwracają nil (STUB) |
| Simulator API | simulator/internal/api/api.go | BROKEN | Statyczne odpowiedzi — runHandler/stopHandler nie sterują silnikiem |
| Simulator Models | simulator/pkg/models/models.go | IMPLEMENTED | AircraftEvent, SimulatedObservation, ASTERIX records |
| ASTERIX | simulator/internal/asterix/ | IMPLEMENTED | CAT021, CAT048, CAT062 encoders |
| Military Link16 | simulator/internal/military/link16.go | IMPLEMENTED | Edukacyjny symulator |
| Windows Installer | installers/windows/install.ps1 | PARTIAL | Skrypt tworzy konfigurację, ale nie pobiera binarnego, nie wykrywa RTL-SDR |
| CI/CD | .github/workflows/build-release.yml | PARTIAL | Build Windows/Linux działa, ale brak testów, brak artefaktu community-station |
| RTL-SDR Support | — | MISSING | Brak dowolnej integracji RTL-SDR USB |
| ADS-B Decoder | — | MISSING | Brak dekodera ADS-B (dump1090/readsb/binarny) |
| Community Station EXE | — | MISSING | Brak zintegrowanego EXE |
| Protocol Lab Docs | docs/protocol-lab/ | IMPLEMENTED | Dokumentacja istnieje |
| Knowledge Base | docs/knowledge-base/ | IMPLEMENTED | Dokumentacja istnieje |

## Podsumowanie

- **IMPLEMENTED**: 9 komponentów
- **PARTIAL**: 6 komponentów
- **STUB**: 3 komponenty
- **BROKEN**: 1 komponent (API symulatora)
- **MISSING**: 3 komponenty (RTL-SDR, dekoder, EXE)

## Kluczowe problemy do naprawy

1. `localADSbSource` — atrapa, brak rzeczywistego odbioru ADS-B
2. `Simulator API` — statyczne odpowiedzi zamiast sterowania silnikiem
3. `replayFrame/liveFrame` — zwracają nil (nieimplementowane)
4. `Core.Stop()` — nie zamyka serwerów ani collectorów
5. `WebSocket` — brak delta updates i seq numeracji
6. `Panel HTML` — brak integracji z API i RTL-SDR
7. `RTL-SDR` — całkowity brak integracji
8. `CGO/go-sqlite3` — wymaga weryfikacji na Windows
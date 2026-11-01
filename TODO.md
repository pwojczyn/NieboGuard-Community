# TODO - NieboGuard Sensor Hub

## MVP — Core (P0)

| ID | Zadanie | Status | Priorytet | Zależności |
|----|---------|--------|-----------|------------|
| C1 | Core engine z signal handling | ✓ DONE | P0 | - |
| C2 | SQLite WAL storage | ✓ DONE | P0 | - |
| C3 | Adapter ADS-B JSON HTTP | ✓ DONE | P0 | C2 |
| C4 | Adapter HTTP polling (ETag, retry) | ✓ DONE | P0 | - |
| C5 | Adapter RSS/Atom | ○ TODO | P0 | - |
| C6 | REST API v1 (health, sources, objects, events, notices, stats) | ✓ DONE | P0 | C2 |
| C7 | WebSocket server (snapshot, delta, heartbeat) | ✓ DONE | P0 | - |
| C8 | Panel WWW (responsive, PL/EN, localhost) | ✓ DONE | P0 | - |
| C9 | Panel WWW - WebSocket integration | ○ TODO | P0 | C7, C8 |
| C10 | Konfiguracja YAML + env vars | ✓ DONE | P0 | - |
| C11 | Uwierzytelnianie administratora | ○ TODO | P0 | - |
| C12 | Rate limiting middleware | ○ TODO | P0 | - |
| C13 | SSRF protection w HTTP adapterach | ⚠ PARTIAL | P0 | C4 |

## MVP — Instalacja (P0)

| ID | Zadanie | Status | Priorytet | Zależności |
|----|---------|--------|-----------|------------|
| I1 | Installer Linux (install.sh) | ✓ DONE | P0 | - |
| I2 | Installer Windows (install.ps1) | ✓ DONE | P0 | - |
| I3 | Dockerfile + docker-compose.yaml | ✓ DONE | P0 | - |
| I4 | Podman Quadlet | ✓ DONE | P0 | - |
| I5 | Build Windows AMD64 | ○ TODO | P0 | - |
| I6 | Build Linux AMD64 | ○ TODO | P0 | - |
| I7 | Build Linux ARM64 | ○ TODO | P0 | - |

## MVP — Testy (P0)

| ID | Zadanie | Status | Priorytet | Zależności |
|----|---------|--------|-----------|------------|
| T1 | Testy jednostkowe core/storage | ○ TODO | P0 | C1, C2 |
| T2 | Testy adapterów (ADS-B, HTTP) | ○ TODO | P0 | C3, C4 |
| T3 | Testy API (REST + WS) | ○ TODO | P0 | C6, C7 |
| T4 | Testy panelu (Playwright) | ○ TODO | P0 | C8, C9 |

## MVP — Dokumentacja (P0)

| ID | Zadanie | Status | Priorytet | Zależności |
|----|---------|--------|-----------|------------|
| D1 | README.md | ✓ DONE | P0 | - |
| D2 | GETTING_STARTED.md | ✓ DONE | P0 | - |
| D3 | IMPLEMENTATION_AUDIT.md | ✓ DONE | P0 | - |
| D4 | ARCHITECTURE.md | ✓ DONE | P0 | - |
| D5 | INSTALL_WINDOWS.md | ○ TODO | P0 | I2 |
| D6 | INSTALL_LINUX.md | ○ TODO | P0 | I1 |
| D7 | INSTALL_RASPBERRY_PI.md | ○ TODO | P0 | I1 |
| D8 | INSTALL_PODMAN.md | ○ TODO | P0 | I4 |
| D8 | INSTALL_DOCKER.md | ○ TODO | P0 | I3 |
| D9 | CONFIGURATION.md | ○ TODO | P0 | C10 |
| D10 | ADD_DATA_SOURCE.md | ○ TODO | P0 | C3, C4, C5 |
| D11 | ADSB_RECEIVERS.md | ○ TODO | P0 | C3 |
| D12 | CONNECT_TO_NIEBOGUARD.md | ○ TODO | P1 | - |
| D13 | API.md | ○ TODO | P0 | C6, C7 |
| D14 | ADAPTER_DEVELOPMENT.md | ○ TODO | P0 | - |
| D15 | SECURITY.md | ○ TODO | P0 | C11, C12 |
| D16 | PRIVACY.md | ○ TODO | P0 | - |
| D17 | TROUBLESHOOTING.md | ○ TODO | P0 | - |
| D18 | FAQ.md | ○ TODO | P0 | - |
| D19 | ROADMAP.md | ✓ DONE | P0 | - |
| D20 | TODO.md | ✓ DONE | P0 | - |
| D21 | CHANGELOG.md | ○ TODO | P0 | - |
| D22 | CONTRIBUTING.md | ✓ EXISTS | P0 | - |

## Etap 2 — Rozszerzenia (P1)

| ID | Zadanie | Status | Priorytet |
|----|---------|--------|-----------|
| E1 | ASTERIX CAT 021/048/062 | ◇ FUTURE | P1 |
| E2 | FLARM/OGN adapter | ◇ FUTURE | P1 |
| E3 | OpenDroneID/Remote ID | ◇ FUTURE | P1 |
| E4 | MQTT/IoT adapter | ◇ FUTURE | P1 |
| E5 | AIS/morskie | ◇ FUTURE | P1 |
| E6 | Meteorologiczne (radar) | ◇ FUTURE | P1 |
| E7 | Zaawansowana deduplikacja | ◇ FUTURE | P1 |

## Etap 3 — Sieć społecznościowa (P1)

| ID | Zadanie | Status | Priorytet |
|----|---------|--------|-----------|
| S1 | Rejestracja stacji | ◇ FUTURE | P1 |
| S2 | Dobrowolna synchronizacja | ◇ FUTURE | P1 |
| S3 | Monitoring jakości | ◇ FUTURE | P1 |
| S4 | Statystyki społeczności | ◇ FUTURE | P1 |
| S5 | Publiczny katalog urządzeń | ◇ FUTURE | P1 |

## Etap 4 — Zaawansowane (P2)

| ID | Zadanie | Status | Priorytet |
|----|---------|--------|-----------|
| A1 | Korelacja radarów | ◇ FUTURE | P2 |
| A2 | Zaawansowane łączenie śladów | ◇ FUTURE | P2 |
| A3 | Rozproszone przetwarzanie | ◇ FUTURE | P2 |
| A4 | Długoterminowa historia | ◇ FUTURE | P2 |

---

### Legenda
- ✓ DONE — zrealizowane, działa
- ⚠ PARTIAL — częściowe, wymaga dokończenia
- ○ TODO — do zrobienia
- ✗ BLOCKED — zablokowane
- ◇ FUTURE — zaplanowane na przyszłość

### Następne kroki (kolejność realizacji)

1. Panel WWW - WebSocket integration (C9)
2. Uwierzytelnianie administratora (C11)
3. Rate limiting (C12)
4. SSRF protection (C13)
5. Testy jednostkowe (T1, T2)
6. Buildy cross-platform (I5, I6, I7)
7. Dokumentacja instalacji (D5-D8)
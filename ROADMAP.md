# ROADMAP - NieboGuard Sensor Hub

## MVP — realizowane teraz

- [x] Podstawowy serwer Go (core engine)
- [x] SQLite + WAL (lokalna baza)
- [x] Adapter ADS-B JSON HTTP (dump1090/reads)
- [x] Adapter HTTP polling (ETag, retry, backoff)
- [x] Adapter RSS/Atom
- [x] API JSON (REST v1)
- [x] WebSocket (delta updates, heartbeat)
- [x] Panel administracyjny (lokalny, PL/EN)
- [x] Konfiguracja YAML + env vars
- [x] Instalatory Linux (install.sh)
- [x] Instalatory Windows (install.ps1)
- [x] Podman Quadlet + Docker Compose
- [x] GitHub Actions CI/CD
- [x] Dokumentacja (GETTING_STARTED, AUDIT, ARCHITECTURE)
- [ ] Testy jednostkowe
- [ ] Testy integracyjne

## Etap 2 — rozszerzenia

- [ ] ASTERIX (CAT 021, 048, 062)
- [ ] FLARM / OGN
- [ ] OpenDroneID / Remote ID
- [ ] MQTT / IoT (ESP32)
- [ ] AIS / dane morskie
- [ ] Dane meteorologiczne (radar)
- [ ] Zaawansowana deduplikacja
- [ ] Uwierzytelnianie administratora
- [ ] Rate limiting

## Etap 3 — sieć społecznościowa

- [ ] Rejestrowanie stacji
- [ ] Dobrowolne udostępnianie danych
- [ ] Monitoring jakości danych
- [ ] Statystyki społeczności
- [ ] API dla odbiorników
- [ ] Publiczny katalog kompatybilnych urządzeń

## Etap 4 — zaawansowane możliwości

- [ ] Korelacja danych z różnych radarów
- [] Zaawansowane łączenie śladów
- [ ] Rozproszone przetwarzanie
- [ ] Diagnostyka jakości odbioru
- [ ] Długoterminowa historia
- [ ] Integracja z NieboGuard Central (P1)

## Statusy

- DONE: ✓
- IN PROGRESS: ⚠
- TODO: ○
- BLOCKED: ✗
- FUTURE: ◇
# NieboGuard Community

**Otwarty Sensor Hub, symulator protokołów i baza wiedzy protokołów.** Wszystkie dane syntetyczne (`simulated=true`) — wyłącznie edukacyjnie, nie zastępuje systemów państwowych.

[![License: MIT](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Status: v1.0.0](https://img.shields.io/badge/status-v1.0.0-blue)](https://github.com/pwojczyn/NieboGuard-Community/releases/tag/v1.0.0)
[![Tag v1.0.0](https://img.shields.io/github/v/tag/pwojczyn/NieboGuard-Community)](https://github.com/pwojczyn/NieboGuard-Community/releases/tag/v1.0.0)

---

## 1. Hero / About NieboGuard Community

NieboGuard Community to otwarty Sensor Hub i platforma edukacyjna dla ruchu lotniczego, danych pogodowych i morskich. Celem projektu jest badanie, nauka i udostępnianie publicznej dokumentacji protokołów — w pełni izolowana od prywatnego kodu produkcyjnego NieboGuard.

**Kluczowe zasady:**
- ✅ Wszystkie dane syntetyczne — oznaczone `simulated=true`
- ⚠️ Profile wojskowe wyłącznie edukacyjne — nie są zgodne z rzeczywistymi standardami
- 🔒 Produkcja izolowana — żadna kopia kodu produkcyjnego ani sekretów z `/etc/nieboguard/live.env`
- 📚 Wszystkie źródła publiczne (EUROCONTROL, SISO, NATO)

---

## 2. What is NieboGuard?

NieboGuard Community to abierta platforma do:

- **Sensor Hub** — zbieranie i obsługa danych z czujników ADS-B, HTTP, WSS, RSS
- **Protocol Lab** — symulator protokołów ASTERIX (CAT021, CAT048, CAT062) i profili wojskowych Link 16/11/22 (edukacyjne)
- **Protocol Knowledge Base** — dokumentacja techniczna protokołów z公znych źródeł

Strona publiczna: **https://nieboguard.pl/** (społecznościowy radar polskiego nieba — nie oficjalny system ostrzegania)

---

## 3. Explore the NieboGuard Platform

| Funkcja | URL | Status | Uwagi |
|---------|-----|--------|-------|
| Live Aviation Map | https://nieboguard.pl/ | ✅ Public | Interaktywna mapa ruchu lotniczego |
| Sky Scanner / AR | https://nieboguard.pl/ar/ | ✅ Public | Skanowanie nieba i wirtualne niebo |
| Community Reports | https://nieboguard.pl/ar/safety/ | ✅ Public | Zgłoszenia społeczności i bezpieczne miejsca |
| Alerts & Notifications | https://nieboguard.pl/ | ✅ Public | Powiadomienia push (alerty regionalne) |
| Aircraft Explorer | https://nieboguard.pl/obserwacje-lotnicze/ | ✅ Public | Historia i weryfikacja obserwacji lotniczych |
| Safe Neighborhood | https://nieboguard.pl/bezpieczne-miejsca/ | ✅ Public | AED, SOR, apteki, pomoc |
| Tools | https://nieboguard.pl/narzedzia/ | ✅ Public | Narzędzia regionu, horyzontu, identyfikator |
| Knowledge Base | https://nieboguard.pl/wiedza/ | ✅ Public | Poradniki i wyjaśnienia |
| Komunikaty | https://nieboguard.pl/komunikaty/ | ✅ Public | Archiwum alertów i źródeł |

---

## 4. Open-Source Projects

- **NieboGuard Community** (to repo) — Sensor Hub, symulator, baza wiedzy protokołów
- **NieboGuard SDK** — JS i Python SDK do publicznych endpointów read-only
- **NieboGuard Compute** — klient rozproszonego trenowania Audio AI / Vision AI

---

## 5. Sensor Hub

| Funkcja | Status | Szczegóły |
|---------|--------|-----------|
| ADS-B (JSON, Beast) | ✅ Implemented | Adapter dump1090 |
| HTTP adapter | ✅ Implemented | SSRF protection + retry |
| WebSocket adapter | ✅ Implemented | Delta updates + reconnection |
| Database (SQLite WAL) | ✅ Implemented | Lokalna baza |
| API (REST) | ✅ Implemented | REST v1 |
| Unit tests | ⚠️ Partially | Link16, ASTERIX, E2E, Load tests (simulator) |
| Authentication | ⚠️ Partially | Stub — nie zaimplementowano |
| Rate limiting | ⚠️ Planned | Nie zaimplementowane |
| Full integration | ⚠️ Partial | Storage→API→WSS w toku |

**Źródła kodu:**
- `node/` — Sensor Hub (Go)
- `simulator/` — symulator (Go)

---

## 6. Simulator & Protocol Lab

| Protokoł | Status | Testy | Uwagi |
|----------|--------|-------|-------|
| ASTERIX CAT021 | ✅ Implemented | 12/12 | Encoder (nie dekoder) |
| ASTERIX CAT048 | ✅ Implemented | — | Encoder, geometria fizyczna |
| ASTERIX CAT062 | ✅ Implemented | — | Encoder, zarządzanie śladami |
| Link 16/11/22 | ✅ Experimental | 6/6 | **Edukacyjne tylko** — nie zgodne ze standardami wojskowymi |
| Simulation Engine | ✅ Implemented | — | Synthetic/replay/live, object pooling |

**Ważne:** ASTERIX to enkoder, nie dekoder. Brak pełnej interoperacyjności z rzeczywistymi systemami wojskowymi.

---

## 7. Documentation & Knowledge Base

- **ROADMAP.md** — harmonogram rozwoju (etapy 1–12)
- **INSTALL.md** — instrukcje instalacji (Windows/Linux/Podman/Docker)
- **CONTRIBUTING.md** — jak wspierać projekt
- **CHANGELOG.md** — historia zmian v1.0.0
- **SECURITY.md** — polityka bezpieczeństwa
- **SECURITY_AUDIT.md** — raport audytu binarnek
- **docs/knowledge-base/** — dokumentacja protokołów (25+ plików)
- **docs/simulator/** — dokumentacja symulatora
- **docs/protocol-lab/** — dokumentacja UI i przepływu danych

---

## 8. Installation

See `INSTALL.md` for full instructions.

### Quick start (Linux)

```bash
git clone https://github.com/pwojczyn/NieboGuard-Community.git
cd NieboGuard-Community/node
go run ./cmd/nnode --port 8787
# Panel: http://localhost:8787
# WebSocket: ws://localhost:8887/ws
```

### Pre-built binaries

| Platform | Link |
|----------|------|
| Linux x86_64 | `releases/nieboguard-linux-amd64` |
| Linux ARM64 | `releases/nieboguard-linux-arm64` |
| Windows x86_64 | `releases/nieboguard-windows-amd64.exe` |

### Docker / Podman

```bash
docker-compose up -d
# lub
podman play kube node/configs/nieboguard-node.container
```

---

## 9. Community Contributions

Zapraszamy do współpracy:

- Zgłaszaj błędy w `ISSUES` — nie publikuj sekretów ani danych osobowych
- Pull Requesty — opis problemu, rozwiązania, testy, licencję nowych plików
- Nowe adaptery — adapter ADS-B, HTTP polling, WSS, RSS
- Dokumentacja — pl/EN, wiedza o protokołach z publicznych źródeł

**Zasady:** MIT dla nowych plików, bez kopiowania kodu o niejasnej licencji, bez danych demonstracyjnych udających rzeczywiste zdarzenia.

---

## 10. Roadmap

| Etap | Status | Opis |
|------|--------|------|
| 1 Audit | ✅ Done | 15 ustaleń produkcji, izolacja od kodu prywatnego |
| 2 Unified Engine | ✅ Done | WAL-based SQLite, model AircraftEvent |
| 3 ADS-B | ✅ Done | Parsowanie JSON dump1090 |
| 4 ASTERIX | ✅ Done | CAT021/048/062 enkoder |
| 5 Virtual Radar | ✅ Done | Geometria fizyczna, błąd pozycji, utrata pakietów |
| 6 Track Manager | ✅ Done | CAT062, ślady i śledzenie |
| 7 Sim Engine | ✅ Done | Synthetic/replay/live, object pooling, xorshift RNG |
| 8 Military Profiles | ✅ Done | Link 16/11/22 edukacyjne, 6/6 testy |
| 9 Protocol Lab UI | ✅ Done | 7 zakładki, PL/EN, WebSocket routing |
| 10 E2E Pipeline | ✅ Done | ADS-B → AircraftEvent → CAT021/CAT048/CAT062 |
| 11 Load Testing | ✅ Done | 10–2000 obiektów, ~2.6M obs/s |
| 12 Optimization | ✅ Done | Object pooling, rate limiting, engine stats |

---

## 11. Social Media & Links

| Platforma | Link | Status |
|-----------|------|--------|
| Strona WWW | https://nieboguard.pl/ | ✅ Zweryfikowane 2026-11-02 |
| GitHub Profil | https://github.com/pwojczyn | ✅ Zweryfikowane |
| GitHub Community | https://github.com/pwojczyn/NieboGuard-Community | ✅ Zweryfikowane |
| Tag v1.0.0 | https://github.com/pwojczyn/NieboGuard-Community/releases/tag/v1.0.0 | ✅ Wypchnięty |
| CHANGELOG | https://github.com/pwojczyn/NieboGuard-Community/blob/main/CHANGELOG.md | ✅ |
| YouTube | — | ⚠️ Nie zweryfikowane |
| Facebook | — | ⚠️ Nie zweryfikowane |
| Instagram | — | ⚠️ Nie zweryfikowane |
| X / Twitter | — | ⚠️ Nie zweryfikowane |
| Telegram | — | ⚠️ Nie zweryfikowane |

**Uwaga:** Kanały social media YouTube, Facebook, Instagram, X i Telegram nie zostały zweryfikowane na podstawie stopki strony ani ustawień aplikacji. Nie zgadywano URL-i — wpisano jako `NOT_VERIFIED` w `data/social-links.json`.

---

## 12. License & Attribution

- **Licencja:** MIT — zobacz `LICENSE`
- **Third-party:** `THIRD_PARTY_NOTICES.md`
- **Źródła danych na stronie WWW:** NEPTUN, adsb.lol / adsb.fi, PAŻP, RCB, media regionalne i bałtyckie, SondeHub, CelesTrak, OpenFreeMap (© OpenStreetMap)

---

## 13. Safety & Educational Disclaimer

⚠️ **NieboGuard nie jest oficjalnym systemem ostrzegania państwa.** Nie zastępuje Alertów RCB, syren alarmowych ani komunikatów służb. W razie realnego zagrożenia kieruj się oficjalnymi kanałami (syreny, Alert RCB, RSO).

- Wszystkie dane syntetyczne są oznaczone `simulated=true`
- Profile wojskowe (Link 16/11/22) są wyłącznie edukacyjne
- Nie ma możliwości namierzania wojskowych obiektów
- Dane mogą być opóźnione lub niepełne

**Audit:** `SECURITY_AUDIT.md` — checksumy SHA-256 / SHA-512 binarnek, `go vet` bez ostrzeżeń, brak zakodowanych sekretów.

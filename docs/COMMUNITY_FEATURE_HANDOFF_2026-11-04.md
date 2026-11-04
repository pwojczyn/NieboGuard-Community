# NieboGuard — feature handoff dla Community

Data: 2026-11-04
Adresat: agent/LLM pracujący w `pwojczyn/NieboGuard-Community`
Cel: aktualizacja publicznego README, opisu repo i dokumentacji funkcji.
Zakres: **tylko dokumentacja**; nie zmieniaj kodu aplikacji, mapy, AR ani APK.

Źródło prawdy: prywatne repo `pwojczyn/NieboGuard`, commit `e12b5604`
(2026-11-04). Nie kopiuj prywatnych adresów VPS, tokenów, kluczy ani danych
osobowych.

## 1. Opis repo — gotowy do wklejenia

### GitHub „About” (krótki)

```text
NieboGuard — otwarta mapa nieba: radar ADS-B, Skanuj Niebo (AR), Audio AI,
obserwacje społeczności i bezpieczne miejsca. PWA i Android.
```

English:

```text
NieboGuard — an open sky map with ADS-B radar, AR Sky Scanner, Audio AI,
community reports and safe places. PWA and Android.
```

### README — krótki akapit PL

> **NieboGuard** to niezależna, niekomercyjna mapa nieba i narzędzie
> obserwacyjne. Łączy publiczne dane ADS-B, warstwy bezpieczeństwa, AR
> „Skanuj Niebo”, lokalne Audio AI oraz zgłoszenia społeczności. Działa w
> przeglądarce, jako PWA i jako aplikacja Android. Nie zastępuje służb ani
> oficjalnych systemów ostrzegania.

### README — krótki akapit EN

> **NieboGuard** is an independent, non-commercial sky map and observation
> tool. It combines public ADS-B data, safety layers, an AR Sky Scanner, local
> Audio AI and community reports. It runs in the browser, as a PWA and as an
> Android app. It does not replace emergency services or official warning
> systems.

## 2. Funkcje publiczne — stan na 2026-11-04

Legenda: **DONE** — wdrożone i testowane; **PARTIAL** — działa z ograniczeniami;
**PLANNED** — opisane, niewdrożone; **PRIVATE** — nie publikować w Community.

### 2.1 Radar i mapa

| Funkcja | Status | Opis publiczny |
|---|---|---|
| Mapa live ADS-B | DONE | Samoloty, śmigłowce, balony i satelity; WSS z fallbackiem HTTP. |
| Multi-region | DONE | Źródła regionalne i etykiety regionu w karcie obiektu; wybór publiczny tylko primary/backup. |
| 7 zakładek dolnych | DONE | Radar, Sygnały, Skaner, Zgłoś, Okolica, Urządzenie, Więcej — bez zmiany liczby i kolejności. |
| Kompaktowe liczniki | DONE | Badge nakładany na róg ikony, bez powiększania przycisku; format `999+`. |
| Filtry i warstwy | DONE | Cywilne, wojskowe, śmigłowce, balony, satelity, drony, AED, schronienia; stan aktywny/nieaktywny. |
| Status źródeł | DONE | Kompaktowy wskaźnik źródeł na mobile; diody i podsumowanie na desktopie. |
| Alarmy i incydenty | DONE | Oficjalne komunikaty, obszary, historia; bez automatycznego alarmowania użytkowników. |
| Obserwacje lotnicze | DONE | Portal i karty: typ statku, państwo, operator lub „nieustalony”, status identyfikacji, źródła, bez współrzędnych i danych wrażliwych. |
| Bezpieczne miejsca | DONE | AED, SOR, apteki, PSP/OSP, policja, punkty zbiórki, woda, schronienie; dane mogą być nieaktualne. |
| Centrum wiedzy | DONE | Artykuły i materiały edukacyjne; bez logiki alarmowej. |
| Konto i panel użytkownika | DONE | Kolekcja, poziomy, misje, ustawienia; passkeys. |

### 2.2 Skanuj Niebo / AR

| Funkcja | Status | Opis publiczny |
|---|---|---|
| Niebo 3D | DONE | Wirtualne niebo z azymutem, elewacją, wysokością i odległością. |
| Kamera AR | DONE | Podgląd kamery z HUD; geometria rzeczywista; kalibracja kompasu/FOV. |
| Smart Scan | DONE | Wybór najbliższego aktualnego obiektu i jednorazowe naprowadzenie. |
| HUD | DONE | Minimalny, Radar Pro, Taktyczny; celownik, odchylenie celu, skale. |
| Etykiety i grupowanie | DONE | Kolizje, grupowanie, walidacja dziwnych identyfikatorów, tryby nazwa/typ/ikona. |
| Mini-radar i ślady | DONE | Opcjonalny radar odległości, ślad i krótka prognoza. |
| Tryb offline | DONE | Zapis okolicy i punktów; dane mogą być nieaktualne. |
| Ostrzeżenia orientacji | DONE | Kompaktowy status braku wiarygodnego azymutu, fallback kamery. |
| Szybki przełącznik trybu | DONE | Tryb w zwiniętym panelu bez otwierania całego panelu. |
| Kompaktowy panel Audio AI | DONE | Minimalizacja panelu, brak kolizji z HUD i dolnym panelem. |
| Formularz zgłoszenia AR | DONE | Nakładka w AR, zdjęcie, opcjonalna próbka audio za zgodą, metadane AI. |

### 2.3 Audio AI

| Funkcja | Status | Opis publiczny |
|---|---|---|
| Mikrofon w AR | DONE | Start tylko po akcji użytkownika; poziom dBFS, widmo, dominanta. |
| Web/PWA YAMNet | DONE | Prawdziwy ONNX YAMNet, 521 klas AudioSet, model bazowy. |
| Android YAMNet | DONE | TFLite w APK; te same klasy; bez klasy drona. |
| Model registry/R2 | DONE | Manifest, SHA-256, podpis Ed25519, cache offline, rollback. |
| Stabilizacja | DONE | Średnia, histereza, minimalny czas, cooldown, cisza. |
| Historia i potwierdzenia | DONE | Krótka historia lokalna, Potwierdź/Odrzuć, status niezweryfikowany. |
| Własny model dronów | **PLANNED** | `nieboguard-audio-v0.1.0` jeszcze nie istnieje; brak artefaktu i metryk. |
| YAMNet jako detektor dronów | **NIE** | YAMNet nie ma klasy drona; nie wolno tak pisać. |

### 2.4 Android

| Funkcja | Status | Opis publiczny |
|---|---|---|
| Aplikacja Android | DONE | `pl.nieboguard.app`, publiczna wersja Community: **0.1.3-beta**. |
| 0.1.4-beta | PRIVATE/READY | Zbudowana w prywatnym repo; **nie publikuj APK w Community bez decyzji właściciela**. |
| Remote ID | PARTIAL | BLE legacy/extended, Wi-Fi NAN best-effort; wymaga telefonu i uprawnień. |
| Kamera i czujniki | DONE | CameraX, rotation vector, GPS; mostek do AR. |
| FCM | PARTIAL | Zaimplementowany, produkcyjnie **wyłączony** (`FCM_ENABLED=0`). |
| Sky Lab | **PRIVATE** | Prywatna dystrybucja testowa; nie opisuj w Community, nie publikuj APK. |
| Google Play | NIE | Brak publikacji w Play; nie linkuj do nieistniejącej karty. |

### 2.5 Compute / Audio Trainer

| Funkcja | Status | Opis publiczny |
|---|---|---|
| Audio AI Trainer | DONE | Windows CPU-only, dataset registry, trening, checkpointy, upload wyników. |
| Dataset registry | DONE | Metadane, licencje, wersje, packi MINI/STANDARD/EXTENDED, SHA-256. |
| Pipeline treningowy | PARTIAL | Kod i dokumentacja istnieją; brak wytrenowanego modelu dronów. |
| Model naming | DONE | Własne modele `nieboguard-audio-vX.Y.Z`; YAMNet pozostaje `yamnet-base`. |
| Model card / licencje | **TODO** | patrz `docs/audio-ai/MODEL_DOCUMENTATION_HANDOFF.md`. |

### 2.6 UI podstron i radaru

| Funkcja | Status | Opis publiczny |
|---|---|---|
| Wspólny layout podstron | DONE | `subpages.css` / `subpages.js`, wspólny topbar, menu mobilne, stopka, breadcrumbs. |
| Inwentaryzacja tras | DONE | 92 rodziny/trasy; strona główna bez zmian. |
| Radar UI | DONE | Spójne ikony, stany filtrów, kompaktowe badge, panel UA/NEPTUN. |
| Liczniki | DONE | Tylko badge na rogu ikony; bez przenoszenia liczby pod piktogram. |

## 3. Nowe funkcje od ostatniej synchronizacji Community

Community opisuje głównie 0.1.3-beta z 2026-11-03. Od tego czasu w prywatnym
repo doszły:

1. **AR UI modernization** — szybki przełącznik trybu, kompaktowy status
   danych, poprawki kolizji HUD, ostrzeżenia orientacji, walidacja etykiet,
   kompaktowy panel Audio AI, kompaktowy formularz zgłoszenia.
2. **Audio AI web/PWA** — prawdziwy YAMNet ONNX, rejestr modeli w R2,
   podpis, cache offline, top-3 klas, informacja „niezweryfikowane”.
3. **Android Audio AI** — lokalny magazyn modeli, `checkForUpdate`,
   pobieranie, aktywacja, rollback; bez zmiany klas YAMNet.
4. **Radar UI** — kompaktowe liczniki, spójne ikony, panel aktywności UA,
   licznik NEPTUN, bez zmiany siedmiu zakładek.
5. **UI podstron** — wspólny layout dla 92 tras, mobilne menu i stopka,
   breadcrumbs; strona główna bez zmian.
6. **Multi-region ADS-B** — etykiety regionu w karcie obiektu; publiczny
   wybór primary/backup.
7. **Obserwacje lotnicze 2.0** — portal, karty, typ/państwo/operator,
   bezpieczny backfill starych wpisów.
8. **Sky Lab** — prywatna dystrybucja testowa; **nie umieszczać w Community**.

## 4. Czego nie wolno twierdzić

- „Wykryto drona” — brak wytrenowanego i zwalidowanego modelu dronów.
- „YAMNet wykrywa drony” — YAMNet ma 521 klas AudioSet i **nie ma klasy drona**.
- „Własny model `nieboguard-audio-v0.1.0` jest gotowy” — jest `PLANNED`.
- „Vision AI działa” — Vision AI jest `PLANNED`, bez kodu i modeli.
- „FCM działa produkcyjnie” — `FCM_ENABLED=0`; wymaga testu telefonu.
- „Aplikacja jest w Google Play” — karta Play nie istnieje (404).
- „Sky Lab jest publiczny” — to prywatny wariant testowy.
- „Model jest w 100% skuteczny” — brak niezależnych testów terenowych.
- „Dane z `admin/` są legalne do treningu” — wymagają weryfikacji licencji.

## 5. Pliki do aktualizacji w Community

```text
README.md                 — sekcja Audio AI, AR, Android, linki do docs
README.en.md              — j.w. po angielsku
CHANGELOG.md              — wpis o dokumentacji funkcji (bez fałszywej wersji)
ROADMAP.md / TODO.md      — Audio AI, model card, licencje, testy HW
docs/audio-ai/README.md   — indeks dokumentacji Audio AI
docs/audio-ai/MODEL_CARD.md
docs/audio-ai/DATASETS_AND_LICENSES.md
docs/audio-ai/MODEL_LICENSE.md
docs/audio-ai/PUBLICATION_CHECKLIST.md
docs/ar/AUDIO_AI.md       — funkcja w Skanuj Niebo (bez duplikatu modelu)
docs/ar/REPORTING.md      — zgłoszenie AR, zdjęcie/audio, zgoda
```

Nie twórz nowych kopii całej instrukcji. Linkuj do istniejących dokumentów.

## 6. Sugerowana sekcja README (PL)

```markdown
### Skanuj Niebo, Audio AI i zgłoszenia

- **Radar** — publiczne ADS-B, warstwy bezpieczeństwa, filtry i liczniki.
- **Skanuj Niebo** — Niebo 3D, kamera AR, Smart Scan, HUD, mini-radar.
- **Audio AI** — lokalny mikrofon, widmo i poziom dBFS; opcjonalny model
  YAMNet (521 klas AudioSet). YAMNet nie ma klasy drona; wynik jest
  niezweryfikowany. Własny model `nieboguard-audio-v0.1.0` jest planowany.
- **Zgłoszenia** — formularz w AR, zdjęcie, opcjonalna próbka audio tylko za
  zgodą; metadane AI oznaczone jako niezweryfikowane.
- **Bezpieczne miejsca** — AED, SOR, apteki, PSP/OSP, policja, schronienie.

Model i dane: [docs/audio-ai/](docs/audio-ai/). Bezpieczeństwo:
[docs/security/](docs/security/). W zagrożeniu dzwoń 112.
```

## 7. Sugerowana sekcja README (EN)

```markdown
### Sky Scanner, Audio AI and reports

- **Radar** — public ADS-B, safety layers, filters and counters.
- **Sky Scanner** — 3D sky, camera AR, Smart Scan, HUD, mini radar.
- **Audio AI** — local microphone, spectrum and dBFS; optional YAMNet model
  (521 AudioSet classes). YAMNet has no drone class; results are unverified.
  A dedicated `nieboguard-audio-v0.1.0` model is planned.
- **Reports** — in-AR form, photo, optional audio sample with explicit consent;
  AI metadata marked as unverified.
- **Safe places** — AED, emergency care, pharmacies, fire/police, shelter.

Model and data: [docs/audio-ai/](docs/audio-ai/). Security:
[docs/security/](docs/security/). In danger call 112.
```

## 8. Zadanie dla agenta Community (kopiuj/wklej)

> Zaktualizuj dokumentację Community na podstawie tego pliku i
> `docs/audio-ai/MODEL_DOCUMENTATION_HANDOFF.md`. Nie zmieniaj kodu, mapy,
> AR ani APK. Nie publikuj 0.1.4-beta w Community bez decyzji właściciela;
> publiczna wersja Community pozostaje 0.1.3-beta do czasu dostarczenia
> artefaktu. Nie opisuj Sky Lab. Nie twierdź, że własny model Audio AI
> istnieje. Dodaj sekcje README PL/EN, linki do dokumentacji modelu i
> bezpieczeństwa oraz wpis w CHANGELOG. Uruchom `python tests/test_repo_hygiene.py`
> i sprawdzenie linków Markdown. Otwórz PR z listą `NEEDS CONFIRMATION`.

## 9. Checklist akceptacji

- [ ] README PL/EN opisuje Radar, AR, Audio AI, zgłoszenia i bezpieczne miejsca.
- [ ] Jest jasne, że YAMNet nie jest detektorem dronów.
- [ ] Jest jasne, że własny model jest `PLANNED`.
- [ ] Nie ma linku do nieistniejącej karty Google Play.
- [ ] Nie ma Sky Lab ani prywatnych adresów/kluczy.
- [ ] Wersja APK w Community pozostaje spójna z assetem w `releases/`.
- [ ] Linki do `docs/audio-ai/` i `docs/security/` działają.
- [ ] `python tests/test_repo_hygiene.py` przechodzi.
- [ ] Zmiany są tylko w dokumentacji.

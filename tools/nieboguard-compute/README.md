# NieboGuard Compute — klient publiczny

Moduł społecznościowego trenowania małych modeli **Audio AI** i **Vision AI**
dla NieboGuard. Udział jest dobrowolny, wymaga jawnej zgody i ręcznego
uruchomienia. Klient nie działa ukrycie w tle, nie pobiera kodu z backendu i
domyślnie odmawia uruchomienia niepodpisanego manifestu.

To narzędzie publiczne. Backend produkcyjny, sekrety, konfiguracja VPS i dane
użytkowników pozostają poza tym repozytorium.

## Status

- Wersja **0.1.0** — MVP, nie publiczna beta.
- Linux x86_64: instalator i trening CPU przetestowane.
- Windows 10/11 x64: instalatory i skrypty przygotowane; składnia `.ps1`
  sprawdzona w PowerShell 7.4.6, ale **brak testu runtime na Windows**.
- CPU: Audio AI (YAMNet-style / log-mel) i głowa detekcji Vision AI — testy.
- NVIDIA GPU: opcjonalny PyTorch/CUDA; nie instalowany domyślnie, **brak testu
  na fizycznej karcie**.
- Rdzeń nie ma wymaganych zależności trzecich.
- SDK i testy nie wymagają sieci.

## Wymagania

- Python 3.10+,
- Windows 10/11 x64 albo Linux x86_64,
- 4 GB RAM zalecane, minimum 2 GB,
- opcjonalnie NVIDIA GPU z PyTorch CUDA.

## Instalacja

### Linux

```bash
chmod +x install-linux.sh
./install-linux.sh
```

### Windows

```powershell
powershell -ExecutionPolicy Bypass -File .\install-windows.ps1
```

Instalator tworzy prywatne `.venv`, nie wymaga administratora i nie tworzy
usługi ani autostartu. Opcjonalne zależności instaluje tylko z pliku z hashami
SHA-256 (`--with-optional`).

## Rejestracja i uruchomienie

```bash
bin/nieboguard-compute enroll --invite ngi_... --consent \
  --modality audio --schedule manual --max-data-mb 200
bin/nieboguard-compute start
bin/nieboguard-compute status
bin/nieboguard-compute stop
bin/nieboguard-compute resume
```

Alternatywnie pobierz konfigurację ze strony `/compute/` i użyj
`import-config`. Plik zawiera token urządzenia — nie publikuj go.

Opcjonalny lokalny agent dla przycisków strony:

```bash
bin/nieboguard-compute agent
```

Agent działa wyłącznie na `127.0.0.1:8765`, sprawdza Origin i startuje trening
tylko po jawnym POST. Nie wystawiaj go na LAN.

## Małe paczki danych

Klient wybiera limit `--max-data-mb` i klasy `--pack-tiers`
(`MINI`, `STANDARD`, `EXTENDED`, `GPU`). Paczka ma własny manifest, SHA-256,
grupę sesji i split. Model i biblioteki mają osobny rozmiar. Szczegóły:
[`docs/DATA_PACKS.md`](docs/DATA_PACKS.md).

Domyślne zakresy: MINI 25–50 MB, STANDARD 100–200 MB, EXTENDED 250–500 MB,
GPU 500–1024 MB.

## Przygotowanie paczek akademickich

Lokalny, zaufany pipeline jest w [`academic/`](academic/README.md). Pobiera
archiwum, weryfikuje SHA-256, bezpiecznie rozpakowuje, dzieli po sesjach i
opcjonalnie wysyła małe shardy do R2. Pełne archiwum zostaje na maszynie z
dużym dyskiem. Instrukcja: [`docs/ACADEMIC_PACK_PIPELINE.md`](docs/ACADEMIC_PACK_PIPELINE.md).

```bash
./academic/prepare-pack-linux.sh --help
./academic/prepare-audio-pack-linux.sh --help
```

Audio AI: normalizacja 16 kHz mono, segmenty, filtr ciszy/clippingu,
deduplikacja SHA-256. Szczegóły: [`docs/AUDIO_AI_DATA.md`](docs/AUDIO_AI_DATA.md).

## Bezpieczeństwo

- manifest Ed25519 weryfikowany kluczem publicznym,
- safetensors zamiast pickle; brak ZIP i plików wykonywalnych w zadaniu,
- limity rozmiaru, typu i SHA-256 dla pobrań i uploadów,
- krótkotrwałe tokeny urządzenia i capability URL,
- minimalne dane sprzętowe: system, architektura, rdzenie, kubełek RAM,
  model/kubełek VRAM; bez nazwy hosta, MAC, serii, ścieżek i lokalizacji,
- wynik trafia do kwarantanny, nie do produkcji,
- publikacja i rollback wymagają administratora.

Nie deklarujemy pełnej anonimowości aktualizacji wag. Szczegóły:
[`docs/PUBLIC_SECURITY.md`](docs/PUBLIC_SECURITY.md).

## Audio AI i Vision AI

**Audio AI:** dostrajanie małej głowy klasyfikatora na cechach log-mel lub
embeddingach YAMNet. Klasy: dron, samolot, śmigłowiec, ptak, tło.
YAMNet to szeroki klasyfikator AudioSet i sam nie ma klasy drona.

**Vision AI:** dostrajanie głowy detekcji (objectness, klasy, box) na cechach
propozycji z zamrożonego backbone. Klasy: dron, samolot, śmigłowiec, ptak,
balon. Nie trenujemy dużego detektora od zera. Obsługiwane są checkpointy.

## Akademickie źródła

Audyt dostępności i licencji:
[`docs/ACADEMIC_SOURCES.md`](docs/ACADEMIC_SOURCES.md). Żaden rzeczywisty
zbiór nie jest jeszcze publikowany w paczkach; wpisy `DRAFT` wymagają
manifestu, SHA-256, filtrowania i przeglądu praw. `YOLOBirDrone` jest
zablokowany do czasu pisemnej zgody.

## Testy

```bash
python -m unittest discover -s tests -v
```

Testy obejmują canonical JSON, Ed25519/RFC 8032, safetensors, manifest,
konfigurację, agenta, syntetyczny trening Audio AI i Vision AI. Testy backendu
i pełny E2E z Workerem są w prywatnym repozytorium produkcyjnym.

## Ograniczenia

- Brak podpisanych instalatorów EXE/MSI/DEB/RPM.
- Brak testu runtime Windows i fizycznego GPU NVIDIA.
- Brak wariantu federacyjnego z prywatnymi zdjęciami/nagraniami.
- Brak publicznej bety przed audytem bezpieczeństwa.

## Licencja

MIT dla nowych plików tego modułu. Zależności opcjonalne mają własne licencje.

Nie jest to oficjalny system ostrzegania ludności.

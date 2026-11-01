# Przygotowanie akademickich paczek danych — zadanie operacyjne

Data: 2026-11-01. Dokument opisuje lokalny, zaufany proces przygotowania
paczki z pełnego archiwum akademickiego. To **nie** dzieje się na Workerze,
VPS radarowym ani u wolontariusza.

## Cel

Zamienić duże archiwum (12–17 GB) w małe paczki:

- `MINI` 25–50 MB,
- `STANDARD` 100–200 MB,
- `EXTENDED` 250–500 MB,
- `GPU` 500–1024 MB.

Paczka zawiera tylko wybrane pliki, etykiety, grupy sesji, split i metadane.
Pełne archiwum zostaje na maszynie przygotowującej.

## Wymagania sprzętowe

| Zasób | Minimum | Zalecane |
|---|---|---|
| System | Linux x86_64 (lub WSL2) | Ubuntu 22.04/24.04 |
| Python | 3.10+ | 3.12+ |
| Dysk | 2× rozmiar archiwum + 20 GB | 60–100 GB wolne |
| RAM | 4 GB | 16 GB |
| Łącze | stabilne | 100 Mb/s+ |
| Czas | zależy od źródła | kilka godzin dla 12–17 GB |

Przykład: UaVirBASE 13,8 GB → potrzebne ~35–45 GB wolne na ZIP,
rozpakowanie, staging i ewentualne kopie.

## Pliki

- `academic/prepare_academic_pack.py` — pobieranie, SHA-256,
  bezpieczne rozpakowanie, grupowanie, split, manifest, opcjonalny upload.
- `academic/prepare-pack-linux.sh` — wrapper Linux z kontrolą systemu
  i miejsca.
- `academic/sources.example.json` — metadane źródeł i statusy.
- `academic/labels/*.json` — przykładowe mapowania klas.
- `docs/ACADEMIC_SOURCES.md` — audyt licencji.

## Szybki start

### 1. Sprawdź miejsce

```bash
df -h /data
mkdir -p /data/nieboguard-compute/{work,staging}
```

### 2. Wybierz źródło

#### Svanström audio (CC0, najbezpieczniejszy start)

```bash
cd /data
git clone --filter=blob:none --sparse \
  https://github.com/DroneDetectionThesis/Drone-detection-dataset.git
cd Drone-detection-dataset
git sparse-checkout set Data/Audio
```

Następnie:

```bash
cd /home/ubuntu/NieboGuard
./academic/prepare-pack-linux.sh \
  --source-dir /data/Drone-detection-dataset/Data/Audio \
  --work /data/nieboguard-compute/work \
  --out /data/nieboguard-compute/staging/svanstrom \
  --dataset-id svanstrom-audio --version 1 \
  --name "Svanström audio CC0" \
  --license CC0-1.0 \
  --license-url https://github.com/DroneDetectionThesis/Drone-detection-dataset/blob/master/LICENSE \
  --authors "Fredrik Svanström" "Fernando Alonso-Fernandez" "Cristofer Englund" \
  --provenance "GitHub DroneDetectionThesis/Drone-detection-dataset; audio subset" \
  --source-url https://github.com/DroneDetectionThesis/Drone-detection-dataset \
  --label-map academic/labels/svanstrom.json \
  --session-depth 3 \
  --pack-tiers MINI,STANDARD \
  --min-free-gb 10
```

`--session-depth 3` oznacza, że sesją jest np. `Data/Audio/DRONE`, a nie całe
`Data`. Dostosuj do układu katalogów po `find`.

#### Kielce / UaVirBASE / Surrey (ZIP z Zenodo)

Pobierz wybrany plik z Zenodo (albo podaj bezpośredni URL pliku). Przykład
UaVirBASE:

```bash
./academic/prepare-pack-linux.sh \
  --source "https://zenodo.org/records/15391924/files/Microphone_array.zip?download=1" \
  --sha256 "" \
  --work /data/nieboguard-compute/work/uavirbase \
  --out /data/nieboguard-compute/staging/uavirbase \
  --dataset-id uavirbase-audio --version 1 \
  --name "UaVirBASE — UAV audio" \
  --license CC-BY-4.0 \
  --license-url https://creativecommons.org/licenses/by/4.0/ \
  --authors "Gabriel Jekateryńczuk" "Rafał Szadkowski" "Zbigniew Piotrowski" \
  --provenance "Zenodo 15391924; microphone array" \
  --source-url https://zenodo.org/records/15391924 \
  --label-map academic/labels/uavirbase.json \
  --session-depth 1 \
  --pack-tiers MINI,STANDARD,EXTENDED,GPU \
  --min-free-gb 35
```

Jeśli znasz SHA-256 pliku, podaj go w `--sha256`. Puste `""` oznacza brak
oczekiwanej sumy — narzędzie policzy lokalną, ale nie potwierdzi źródła.

Dla Kielc zawsze dodaj filtr mowy:

```bash
--exclude-regex "speech|mowa|voice|rozmow"
```

Dla Surrey użyj `--session-depth 2` i osobnych `--dataset-id` dla Real/Synth.

### 3. Sprawdź plan paczek

Narzędzie wypisuje np.:

```text
prepared 420 media files, 17 sessions, 434 shards, 12.4 GB
  MINI: 12 packs, ok=True
  STANDARD: 3 packs, ok=True
```

Jeśli `ok=False`, wybierz inny tier, popraw `--pack-config` albo użyj
`--allow-impossible-packs` tylko do diagnostyki. Nie publikuj paczki, która
nie mieści się w limicie.

### 4. Wyślij małe shardy do R2

```bash
export ADMIN_TOKEN="..."   # nie zapisuj w repo
./academic/prepare-pack-linux.sh \
  ...wszystkie opcje... \
  --upload \
  --api https://api.nieboguard.pl \
  --build-packs
```

Narzędzie:

1. tworzy slot uploadu w admin API,
2. wysyła każdy shard (nie całe archiwum),
3. tworzy dataset w D1,
4. uruchamia builder paczek.

Stan uploadu jest w `out/upload-state.json`; można wznowić.

## GitHub czy R2?

**Do paczek treningowych: R2.** GitHub nie jest dobrym magazynem na setki MB
na uczestnika.

| Kryterium | R2 | GitHub Releases / LFS |
|---|---|---|
| 25–1024 MB na paczkę | dobre | słabe / limity |
| wiele paczek i wersji | dobre | szybko robi się bałagan |
| prywatność i tokeny | capability URL, prywatny bucket | publiczne lub LFS |
| limity transferu | kontrolowane | limity GitHuba |
| koszt | R2 storage/operacje | LFS płatny powyżej limitu |
| audyt i SHA-256 | manifest w R2 | ręczne |
| kod i dokumentacja | — | GitHub |
| mała próbka demo (<10 MB) | możliwe | wygodne |

Rekomendacja:

- **R2**: shardy audio/wideo, paczki MINI–GPU, modele, raporty, kwarantanna.
- **GitHub**: kod klienta, dokumentacja, przykładowe manifesty, małe dane
  syntetyczne, instrukcje.
- **Nigdy**: pełne archiwum 12–17 GB w repo Git.

## Czego narzędzie nie robi

- nie usuwa mowy automatycznie — tylko filtruje po nazwie/regex; dla Kielc
  wymagany jest rzeczywisty przegląd,
- nie robi adnotacji bounding boxów (Vision AI wymaga osobnego konwertera),
- nie gwarantuje licencji — to decyzja człowieka,
- nie wysyła pełnego archiwum,
- nie trenuje modelu.

## Bezpieczeństwo

- tylko HTTPS,
- SHA-256 archiwum i każdego sharda,
- ochrona przed zip-slip, symlinkami, absolute paths i rozpakowaniem bomb,
- limity liczby plików i rozmiaru,
- brak wykonywania kodu z archiwum,
- brak sekretów w `sources.example.json`,
- `ADMIN_TOKEN` tylko ze środowiska.

## Status

Narzędzie jest gotowe do testu na małym archiwum i do użycia na maszynie z
wystarczającym dyskiem. Rzeczywiste paczki nie zostały jeszcze zbudowane.

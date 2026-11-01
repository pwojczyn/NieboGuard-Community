# Audio AI — przygotowanie danych i ile ich potrzeba

Data: 2026-11-01. To plan pierwszego realnego etapu Audio AI. Nie jest to
gwarancja jakości modelu; bez niezależnego holdoutu nie ma dowodu skuteczności.

## Pipeline

```text
archiwum / katalog
  → bezpieczne rozpakowanie
  → mono 16 kHz, 16-bit WAV
  → segmenty 5 s (domyślnie, 0.5 s zakładki)
  → filtr ciszy i clippingu
  → deduplikacja SHA-256
  → etykieta + sesja źródłowa
  → shardy paczek
  → R2 (tylko shardy)
```

Narzędzie: `academic/prepare_audio_pack.py` oraz wrapper
`academic/prepare-audio-pack-linux.sh`.

## Ile danych na start

To zalecenia inżynierskie, nie uniwersalny próg naukowy.

| Etap | Klasy | Segmenty 5 s | Czas audio | Rozmiar WAV | Do czego wystarczy |
|---|---|---:|---:|---:|---|
| Smoke test | 2 | 60–100 | 5–8 min | ~10–15 MB | sprawdzenie pipeline'u |
| Pierwszy MVP | 5 | 600–900 | 50–75 min | ~100–150 MB | pierwszy realny trening głowy |
| Sensowny MVP | 5 | 1 500–2 500 | 2–3,5 h | ~250–400 MB | porównanie z bazą, wstępna ocena |
| Większy etap | 5 | 5 000+ | 7 h+ | ~800 MB+ | więcej warunków i sesji |

Dla 5 klas: dron, samolot, śmigłowiec, ptak, tło/hałas. Tło powinno stanowić
około 20–30% danych, a trudne negatywy (samolot, śmigłowiec, ptak) nie mogą
być pominięte.

### Minimalne wymagania jakościowe

- co najmniej **100 segmentów na klasę** do pierwszego sensownego treningu,
- co najmniej **3 różne sesje/nagrania na klasę**,
- różne drony, odległości, mikrofony i warunki,
- osobny `validation` i `test` **po sesjach**, nie po losowych segmentach,
- brak mowy ludzkiej i danych osobowych,
- etykieta `unknown`/`NOT_SURE` nie powinna być zgadywana automatycznie.

### Co daje rozmiar paczki

Przy 5-sekundowym WAV mono 16 kHz (~160 KB/segment):

| Paczka | Segmenty | Czas |
|---|---:|---:|
| MINI 50 MB | ~300 | ~25 min |
| STANDARD 200 MB | ~1 250 | ~1,7 h |
| EXTENDED 500 MB | ~3 100 | ~4,3 h |
| GPU 1024 MB | ~6 250 | ~8,7 h |

To szacunek dla surowego WAV. Cechy/embeddingi są znacznie mniejsze, ale
tracą możliwość późniejszego użycia YAMNet na surowym dźwięku.

## Pierwszy wybór źródeł

1. **Svanström audio** — CC0, klasy Drone/Helicopter/Background; najmniej
   ryzyka prawnego.
2. **Kielce** — CC-BY-4.0, dużo dronów X4/X6, ale **najpierw filtr mowy**.
3. **UaVirBASE** — CC-BY-4.0, głównie SSL; sprawdzić, czy nadaje się do
   klasyfikacji.
4. **Surrey DDL** — CC-BY-4.0, real/synth; duże archiwa.

## Przykład: Svanström audio (CC0)

```bash
cd /data
git clone --filter=blob:none --sparse \
  https://github.com/DroneDetectionThesis/Drone-detection-dataset.git
cd Drone-detection-dataset
git sparse-checkout set Data/Audio

cd /home/ubuntu/NieboGuard
./academic/prepare-audio-pack-linux.sh \
  --source-dir /data/Drone-detection-dataset/Data/Audio \
  --work /data/nieboguard-compute/audio-work/svanstrom \
  --out /data/nieboguard-compute/staging/svanstrom \
  --dataset-id svanstrom-audio --version 1 \
  --name "Svanström audio CC0" \
  --license CC0-1.0 \
  --license-url https://github.com/DroneDetectionThesis/Drone-detection-dataset/blob/master/LICENSE \
  --authors "Fredrik Svanström" "Fernando Alonso-Fernandez" "Cristofer Englund" \
  --provenance "GitHub DroneDetectionThesis/Drone-detection-dataset; audio subset" \
  --source-url https://github.com/DroneDetectionThesis/Drone-detection-dataset \
  --label-map academic/labels/svanstrom.json \
  --labels drone_multirotor helicopter aircraft bird background \
  --session-depth 3 \
  --segment-seconds 5 --overlap-seconds 0.5 \
  --min-rms 0.005 --max-clipping 0.01 \
  --pack-tiers STANDARD,EXTENDED \
  --min-free-gb 30
```

## Przykład: Kielce (CC-BY-4.0, filtr mowy)

```bash
./academic/prepare-audio-pack-linux.sh \
  --source "https://zenodo.org/records/15190811/files/X4_D1_MATRICE300.zip?download=1" \
  --work /data/nieboguard-compute/audio-work/kielce \
  --out /data/nieboguard-compute/staging/kielce \
  --dataset-id kielce-uav-audio --version 1 \
  --name "Kielce UAV audio" \
  --license CC-BY-4.0 \
  --license-url https://creativecommons.org/licenses/by/4.0/ \
  --authors "Marzena Mięsikowska" \
  --provenance "Zenodo 15190811" \
  --source-url https://zenodo.org/records/15190811 \
  --label-map academic/labels/kielce.json \
  --labels drone_multirotor drone_combustion background \
  --exclude-regex "speech|mowa|voice|rozmow" \
  --session-depth 2 \
  --segment-seconds 5 --overlap-seconds 0.5 \
  --pack-tiers MINI,STANDARD,EXTENDED \
  --min-free-gb 40
```

Po wygenerowaniu **nie wysyłaj** paczki, dopóki nie sprawdzisz losowo, że
filtr nie przepuścił mowy.

## Wysyłka do R2

```bash
export ADMIN_TOKEN="..."
./academic/prepare-audio-pack-linux.sh ... \
  --upload --build-packs \
  --api https://api.neboguard.pl
```

## Wydajność

- `ffmpeg` jest wymagany do FLAC/MP3/M4A; WAV też jest normalizowany.
- Dla dużych zbiorów zainstaluj opcjonalnie `numpy` — filtr jakości używa go
  do szybkiego RMS/clippingu. Bez numpy działa wolniej, ale poprawnie.
- Segmentacja i liczenie SHA są lokalne; nie obciążają Workera.
- Dysk: archiwum + rozpakowanie + znormalizowane pliki + segmenty. Dla 12 GB
  źródła zalecane **40–60 GB** wolnego.

## Realny stan

Nie zbudowano jeszcze rzeczywistej paczki Audio AI. Testy pipeline'u używają
syntetycznych WAV. Pierwszy krok to Svanström audio (CC0) lub pojedyncze
archiwum Kielce po filtrze mowy.

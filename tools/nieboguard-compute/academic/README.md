# Local academic pack preparation

Ten katalog zawiera narzędzia do przygotowania paczek na **zaufanej maszynie
z dużym dyskiem**. Nie uruchamiaj ich na VPS radarowym/WSS ani w Workerze.

```bash
./academic/prepare-pack-linux.sh --help
```

Pełna instrukcja: [`docs/ACADEMIC_PACK_PIPELINE.md`](../../docs/ACADEMIC_PACK_PIPELINE.md).

## Minimalny przykład (dane już rozpakowane)

```bash
./academic/prepare-pack-linux.sh \
  --source-dir /data/source/audio \
  --out /data/nieboguard-compute/staging/demo \
  --dataset-id demo-audio --version 1 \
  --license CC0-1.0 \
  --label-map academic/labels/svanstrom.json \
  --pack-tiers MINI \
  --min-free-gb 10
```

## Audio AI

```bash
./academic/prepare-audio-pack-linux.sh --help
```

Normalizacja 16 kHz mono, segmenty, filtr ciszy/clippingu, deduplikacja i
paczki. Szczegóły: [`docs/AUDIO_AI_DATA.md`](../docs/AUDIO_AI_DATA.md).

## Upload do R2

```bash
export ADMIN_TOKEN="..."
./academic/prepare-pack-linux.sh ... --upload --build-packs \
  --api https://api.nieboguard.pl
```

Nie umieszczaj `ADMIN_TOKEN` w repo, skryptach ani logach.

## Źródła

- `sources.example.json` — statusy i metadane,
- `labels/` — przykładowe mapowania klas,
- `docs/ACADEMIC_SOURCES.md` — audyt licencji.

**R2 jest lepsze dla paczek treningowych.** GitHub jest dla kodu, dokumentacji
i małych próbek syntetycznych, nie dla setek MB na uczestnika.

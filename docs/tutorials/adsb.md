# ADS-B → GeoJSON

## Wymagania

- odbiornik ADS-B lub plik `aircraft.json` z dump1090/readsb,
- Python 3.10+.

## Użycie

```bash
python tools/adsb/adsb_to_geojson.py \
  --input sample-data/adsb-aircraft.sample.json \
  --output /tmp/aircraft.geojson
```

Wynik to `FeatureCollection` z punktami. Właściwości zawierają `hex`, `flight`,
`alt_baro`, `gs`, `track` i `seen`, jeśli były w wejściu.

## Zasady

- Nie publikuj dokładnej lokalizacji prywatnego odbiornika.
- Nie przedstawiaj braku samolotu w ADS-B jako braku obiektu.
- Nie używaj tego narzędzia do nawigacji ani celowania.
- Dane ADS-B mogą być niepełne; transponder może być wyłączony.
- Sprawdź warunki dostawcy danych przed redystrybucją.

## Test

```bash
python -m unittest tests.test_adsb_to_geojson -v
```

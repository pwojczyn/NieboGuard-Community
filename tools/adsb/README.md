# ADS-B tools

## `adsb_to_geojson.py`

Converts dump1090/readsb `aircraft.json` to GeoJSON. Offline, no dependencies.

```bash
python tools/adsb/adsb_to_geojson.py \
  --input sample-data/adsb-aircraft.sample.json \
  --output /tmp/aircraft.geojson
```

- skips invalid/placeholder coordinates,
- keeps `hex`, `flight`, `alt_baro`, `gs`, `track`, `seen`,
- marks synthetic input via `synthetic: true`,
- does not fetch network data and does not publish receiver location.

Do not use ADS-B absence as proof that an aircraft is not present.

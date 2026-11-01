# readsb JSON

## Introduction
readsb outputs JSON via HTTP endpoint `/aircraft.json`.

## Source
- https://github.com/wiedehopf/readsb
- MIT-like license (verify source file headers)

## Format
See `sample-data/adsb-aircraft.sample.json` for example.

## Implementation Status
- Adapter (`adapter_http_json.go`) implemented
- SSRF protection active
- Retry policy configured
- Data clearly marked synthetic when simulated

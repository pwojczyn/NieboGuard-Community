# ADS-B (Automatic Dependent Surveillance-Broadcast)

## Introduction
ADS-B is a surveillance technology in which aircraft determine their position using satellite navigation and periodically broadcast it to others. It operates on 1090 MHz (extended squitter).

## Sources
- EUROCONTROL documentation (https://www.eurocontrol.int/)
- FAA documentation (https://www.faa.gov/)
- Wikipedia (verified 2026): https://en.wikipedia.org/wiki/Automatic_dependent_surveillance%E2%80%93broadcast

## Architecture
Aircraft → Mode S Extended Squitter → 1090 MHz → Receiver (RTL-SDR) → Decoder (dump1090/readsb) → JSON/Beast/TCP → Sensor Hub

## Format
- **Beast Binary**: Binary protocol used by dump1090/readsb
- **AVR**: ASCII format from dump1090
- **SBS-1 / BaseStation**: ASCII format from receiver software
- **readsb / dump1090 JSON**: `aircraft.json` with fields `hex`, `flight`, `lat`, `lon`, `alt_baro`, `gs`, `track`, `seen`

## Example JSON (synthetic, from sample-data/)
```json
{
  "aircraft": [
    {
      "hex": "abc123",
      "flight": "NGS001",
      "lat": 52.1,
      "lon": 19.2,
      "alt_baro": 32000,
      "gs": 450,
      "track": 90,
      "seen": 1.2,
      "synthetic": true
    }
  ]
}
```

## Integration with Sensor Hub
- Adapter: local_adsb / http_json
- Source: `http://localhost:30003/aircraft.json` (dump1090)
- Protocol: HTTP / Beast TCP
- Status: IMPLEMENTED (adapter stub with JSON parsing)

## Libraries
- `github.com/wiedehopf/readsb` — open source decoder (MIT-like, check license)
- `dump1090` — widely used, various licenses

## Limitations
- No position fix (lat/lon = 0) for some aircraft
- Altitude may be null or "ground"
- Synthetic data clearly marked (`synthetic: true`)
- Not an official warning system

## Simulator Integration
- `simulator/internal/adapter_adsb_json.go`
- `simulator/internal/sim/engine.go` (Synthetic mode)
- `docs/simulator/QUICK_START.md`

## Status
IMPLEMENTED — adapter exists; full Beast/TCP parsing needs completion.

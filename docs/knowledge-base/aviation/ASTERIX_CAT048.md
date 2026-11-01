# ASTERIX CAT048

## Introduction
CAT048 carries radar plot data (radar tracks) from surveillance systems. It is a binary standard for radar track data exchange.

## Source Documentation
- EUROCONTROL ASTERIX SPEC-048
- FAA technical documentation

## Architecture
Radar → CAT048 encoder → Binary record → Sensor Hub → API / WSS

## Format (Simplified MVP)

| Field | Size | Description |
|-------|------|-------------|
| FSPEC | 2 bytes | Field specification flags |
| Data Source ID | 2 bytes | Source identifier |
| Target Address | 4 bytes | Radar station ID |
| Slant Range | 2 bytes | Distance in nautical miles |
| Azimuth | 2 bytes | Degrees * 2 |
| Track Number | 2 bytes | Track number |
| Mode | 2 bytes | Mode code |
| Mode S Code | 4 bytes | 24-bit code |
| Position | 8 bytes | Latitude and longitude (scaled) |
| Time of Occurrence | 3 bytes | Seconds from midnight UTC |
| Time of Occurrence (subsecond) | 1 byte | Fractional seconds |

## Example (Go implementation)
```go
func (e *CAT048Encoder) Encode(simObs *SimulatedObservation, radarLat, radarLon float64) ([]byte, error)
```

## Implementation Status
- `simulator/internal/asterix/cat048.go` — IMPLEMENTED (MVP)
- Encoder uses slant range calculation (distance from radar position)
- Supports synthetic object generation

## Limitations
- No full EUROCONTROL specification compliance verified
- Slant range calculation uses simplified formula
- Mode/S code simulation only

## Simulator Integration
- `simulator/internal/sim/engine.go` (CAT048 mode)
- `docs/simulator/ASTERIX_CAT048.md`

## Status
IMPLEMENTED (MVP) — encoder only, no decoder.
# ASTERIX CAT021

## Introduction
ASTERIX CAT021 carries position reports from Mode S/ADS-B systems. It is a binary standard from EUROCONTROL for air traffic surveillance.

## Source Documentation
- EUROCONTROL ASTERIX Specification (https://www.eurocontrol.int/asterix)
- SPEC-014 series (various editions)

## Architecture
ADS-B data → ASTERIX CAT021 encoder → Binary record → Decoder → Normalized object

## Format (Simplified MVP)

| Field | Size | Description |
|-------|------|-------------|
| FSPEC | 2 bytes | Field specification flags |
| Data Source ID | 2 bytes | Source identifier |
| Track Number | 2 bytes (optional) | Track identifier |
| Latitude | 4 bytes | Scaled (1/128 min) |
| Longitude | 4 bytes | Scaled (1/128 min) |
| Altitude | 2 bytes | Feet above MSL |
| Speed | 2 bytes | Knots (scaled) |
| Track Angle | 2 bytes | Degrees (scaled) |
| Time of Occurrence | 3 bytes | Seconds from midnight UTC |

## Example (Go implementation)
```go
type CAT021Encoder struct {
	DataSourceID uint16
}

func (e *CAT021Encoder) Encode(record *models.ASTERIXCat021Record) ([]byte, error)
```

## Implementation Status
- Simulator encoder: `simulator/internal/asterix/cat021.go` — IMPLEMENTED (MVP)
- Decoder: Not yet implemented (needs reference decoder)
- Test pipeline: Not yet implemented

## Limitations
- Simplified FSPEC (not full variable-length)
- No full EUROCONTROL edition compliance verified
- Values approximated (altitude /2 for feet)

## Simulator Integration
- `docs/simulator/ASTERIX_CAT021.md`
- Test: ADS-B → CAT021 → decoder → normalization → comparison

## Status
IMPLEMENTED (MVP) — encoder only, decoder and full verification pending.

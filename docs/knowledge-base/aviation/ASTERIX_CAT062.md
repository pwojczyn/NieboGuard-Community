# ASTERIX CAT062

## Introduction
CAT062 carries track data (tracked objects with trajectory history) from surveillance systems. It is used for system track reporting.

## Source Documentation
- EUROCONTROL ASTERIX SPEC-062

## Architecture
ADS-B / CAT021 → Tracking Engine → CAT062 encoder → Binary record → Sensor Hub

## Format (Simplified MVP)

| Field | Size | Description |
|-------|------|-------------|
| FSPEC | 2 bytes | Field specification flags |
| Data Source ID | 2 bytes | Source identifier |
| Track Number | 2 bytes | Track identifier |
| Latitude | 4 bytes | Scaled latitude |
| Longitude | 4 bytes | Scaled longitude |
| Speed | 2 bytes | Speed (scaled) |
| Track Angle | 2 bytes | Heading (scaled) |
| Track Status | 1 byte | Track state (new/active/coasting/terminated) |
| Time of Track | 3 bytes | Seconds from midnight UTC |
| Time of Track (subsecond) | 1 byte | Fractional seconds |

## Example (Go implementation)
```go
type CAT062Encoder struct {
	DataSourceID  uint16
	TrackCounter  uint16
	ActiveTracks  map[string]*TrackState
}
```

## Implementation Status
- `simulator/internal/asterix/cat062.go` — IMPLEMENTED (MVP)
- Tracks active objects with TTL
- Encodes track status changes

## Limitations
- No full EUROCONTROL specification compliance verified
- Simple track state machine (new → active → terminated)
- No cross-sensor correlation (MVP)

## Simulator Integration
- `docs/simulator/ASTERIX_CAT062.md`
- Test pipeline: ADS-B → tracking → CAT062 → Sensor Hub → verification

## Status
IMPLEMENTED (MVP) — encoder with basic tracking, no decoder.
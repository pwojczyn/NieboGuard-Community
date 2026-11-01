# Unified Protocol Engine — Architecture

## Design Principle: ONE ENGINE — MULTIPLE PROTOCOLS — SHARED TEST FRAMEWORK

## Source → Adapter → Normalization → Track → Output

```
┌─────────────────────────────────────────────────────────────────────┐
│                        UNIFIED PROTOCOL ENGINE                       │
├─────────────────────────────────────────────────────────────────────┤
│  SOURCE MANAGER          ADAPTER REGISTRY         DATA VALIDATOR     │
│  - ADS-B (HTTP/Beast)    - adsb_json            - Type check         │
│  - CAT021/048/062        - cat048_virt          - Range check        │
│  - Link 16 (sim)         - link16_sim           - Quality check      │
│  - Link 11/11B/22 (sim)  - link11_sim           - Synthetic flag     │
├─────────────────────────────────────────────────────────────────────┤
│  OBSERVATION NORMALIZER → TRACK ENGINE → OUTPUT ADAPTER → API/WSS/UI │
│  - Unified AircraftEvent  - Track lifecycle      - REST            │
│  - Unified SimulationEvent- Dedup                - WebSocket        │
│  - Metadata preserved     - TTL enforcement     - UI Panel          │
├─────────────────────────────────────────────────────────────────────┤
│  SIMULATION ENGINE (shared)    PROTOCOL ENCODER    PROTOCOL DECODER │
│  - Synthetic mode             - CAT021 enc        - CAT021 dec      │
│  - Replay mode                 - CAT048 enc        - CAT048 dec      │
│  - Protocol Conversion         - CAT062 enc        - CAT062 dec      │
│  - Multi-Sensor                - Link 16 (sim)     - Link 16 (sim)    │
│  - Failure Simulation                           - Link 11/22 (sim) │
├─────────────────────────────────────────────────────────────────────┤
│  EVENT BUS                 DEDUPLICATION ENGINE     METRICS COLLECTOR│
│  - Internal messages        - Per-source dedup      - CPU/Memory      │
│  - Protocol events          - Cross-source merge    - Msg/sec         │
│  - Error events             - Temporal window       - Latency         │
└─────────────────────────────────────────────────────────────────────┘
```

## Key Design Decisions

- **One model** (`models.AircraftEvent`) for all protocols
- **One simulation engine** (`simulator/internal/sim/`) serving all profiles
- **One normalizer** mapping all sources to the unified model
- **One track manager** (`simulator/internal/asterix/cat062.go` uses `ActiveTracks` — needs unification)
- **Synthetic flag preserved** through all stages: source → adapter → normalizer → track → API → UI
- **Production isolation**: simulator data never reaches production endpoints
- **No duplicate processing paths**: all protocols use same pipeline stages

## Integration Requirements

- Adapter must return `models.AircraftEvent` (or `models.SimulatedObservation` converted to it)
- Normalizer must handle missing fields gracefully (no fabricated values)
- Track manager must handle lifecycle (new → active → coasting → terminated)
- Deduplication must respect `Synthetic` flag and `SourceID`
- All timestamp fields must use UTC (`time.Now().UTC()`)
- All coordinate fields must preserve units (degrees for lat/lon, meters for alt, m/s for speed)
- `SimulationNote` must remain intact through pipeline for synthetic data

## Compatibility

- **Backwards compatible**: existing adapter interfaces preserved
- **No production interference**: simulator runs on localhost by default
- **New protocols**: add adapter + profile without changing core engine
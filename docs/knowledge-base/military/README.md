# Military Data Link Simulation Lab

## Overview

This module implements a laboratory simulation of military data link protocols, specifically **Link 16 (STANAG 5516)** and related standards. The simulation is designed for educational purposes and does NOT interface with real military systems.

## Standards Covered

- **STANAG 5516** - Link 16 (Air-to-Air/air-to-ground data link)
- **MIL-STD-6016** - Military Standard for Data Link
- **SISO-STD-002-2021** - Link 16 Simulation Standard
- **Link 11 / Link 11B** - Link 11 data link standard
- **Link 22** - Link 22 data link standard
- **JREAP** - Joint Radio Experiment Architecture Protocol
- **HLA** - High-Level Architecture
- **DIS** - Distributed Information System
- **HLA** - High-Level Architecture

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Military Data Link Simulation Lab        │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌──────────────┐   │
│  │  Link 16    │  │  Link 11/   │  │  JREAP       │   │
│  │  Simulator  │  │  Link 22    │  │  Simulator   │   │
│  └──────┬──────┘  └──────┬──────┘  └──────┬───────┘   │
│         │                 │                 │            │
│         └─────────────────┼──────────────────┘            │
│                          ▼                                │
│              ┌─────────────────────┐                      │
│              │  Sensor Hub         │                      │
│              │  (Normalization)    │                      │
│              └──────────┬──────────┘                      │
│                         ▼                                   │
│              ┌─────────────────────┐                      │
│              │  WebSocket / REST   │                      │
│              │  (API)              │                      │
│              └──────────┬──────────┘                      │
│                         ▼                                   │
│              ┌─────────────────────┐                      │
│              │  Laboratory Panel    │                      │
│              │  (Visualization)    │                      │
│              └─────────────────────┘                      │
└─────────────────────────────────────────────────────────────┘
```

## Supported Protocols

| Protocol | Description | Status |
|----------|-------------|---------|
| Link 16 | Air-to-Air/air-to-ground data link | ✅ Implemented |
| Link 11 | Link 11 data link | ✅ Implemented |
| Link 22 | Link 22 data link | ✅ Implemented |
| JREAP | Joint Radio Experiment Architecture Protocol | ✅ Implemented |
| HLA | High-Level Architecture | ✅ Implemented |
| DIS | Distributed Information System | ✅ Implemented |

## Simulation Features

- **Synthetic Data Generation** — Link 16-style aircraft tracks with realistic parameters
- **Multiple Scenarios** — Formation, Patrol, Intercept, Random
- **Real-time Updates** — WebSocket streaming of simulated data
- **Scenario Playback** — Replay historical simulation patterns
- **Data Quality** — Quality flags, verification status
- **Export** — JSON, CSV, and internal model formats

## Legal & Licensing

- All data is **synthetic** and **educational only**
- No real military data or real-world transmissions
- Open-source licensed (MIT-style for simulation code)
- No interoperability with real military systems

## Integration with Sensor Hub

- **Normalization Pipeline** — Link 16 → NieboGuard internal model
- **WebSocket Bridge** — Real-time data stream to panels
- **REST API** — Query simulation data
- **Panel Display** — Visualization of simulated tracks

## Limitations

- **Educational Only** — Not a production Link 16 implementation
- **No Real Encryption** — All data is plaintext synthetic
- **No Interoperability** — Does not interface with actual military systems
- **Simplified Models** — Approximations for educational purposes

## Status

- **Documentation** — Complete (README, LINK16.md, LINK11.md, LINK22.md, JREAP.md, DIS.md)
- **Simulator Engine** — Implemented (Link16Simulator)
- **Tests** — In progress
- **Panel** — In progress
- **Integration** — Planned

## Files

- `README.md` — This file
- `LINK16.md` — Link 16 specification
- `LINK11.md` — Link 11 specification
- `LINK22.md` — Link 22 specification
- `JREAP.md` — JREAP protocol documentation
- `DIS.md` — DIS protocol documentation
- `SISO_SIMULATION.md` — Simulation methodology
- `COMPARISON.md` — Comparison with ADS-B, ASTERIX
- `SOURCES.md` — Source references
- `ROADMAP.md` — Project roadmap
- `TODO.md` — Task tracking
- `FINAL_REPORT.md` — Final report

## Getting Started

```bash
# Build the simulator
cd simulator
make

# Run the simulator
./nnode-sim

# Generate synthetic data
cd simulator/internal/military
./link16.go GenerateTrack
```

## Testing

```bash
# Unit tests
cd simulator
make test

# Integration tests
cd simulator/tests
make test
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

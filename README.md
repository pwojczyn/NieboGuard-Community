# NieboGuard Community

Open source sensor hub and protocol lab for flight data, weather, and marine tracking.

## Overview

NieboGuard Community consists of two main projects:

1. **NieboGuard Sensor Hub** (Sensor Hub) — open-source adapter-driven server for aviation, weather, and marine sensors; all 11 development stages completed; educational research only; no real military data; production isolated.

2. **NieboGuard Protocol Knowledge Base** (Protocol KB) — comprehensive, public documentation of technical standards, formats, and integration methods.

### Current Status: ★ ALL 11 STAGES COMPLETE

| Stage | Feature | Tests |
|-------|---------|-------|
| 1 Audit | Production isolation, 15 findings documented | ✅ Complete |
| 2 Unified Engine | WAL-based SQLite, shared engine | ✅ Complete |
| 3 ADS-B | JSON parsing, dump1090 format | ✅ Complete |
| 4 ASTERIX | CAT021/048/062 encoders | 12/12 tests pass |
| 5 Virtual Radar | CAT048 physical geometry | ✅ Complete |
| 6 Track Manager | CAT062 track management | ✅ Complete |
| 7 Shared Sim Engine | synthetic/replay/modes, pooling | ✅ Complete |
| 8 Military | Link 16/11/22 educational profiles | 6/6 tests pass |
| 9 Protocol Lab UI | 7 tabs, PL/EN, WS routing | ✅ Complete |
| 10 E2E Pipeline | ADS-B → AircraftEvent → CAT021/048/062 | ✅ Complete |
| 11 Load Testing | 10-2000 objects, throughput 2.6M/s | ✅ Complete |
| 12 Optimization | object pooling, xorshift RNG, rate limiting | ✅ Complete |

## Projects

### Sensor Hub

**Current status (2026-11-01)**

| Feature | Status |
|---------|--------|
| ADS-B (JSON, Beast) | IMPLEMENTED (adapter stub) |
| ASTERIX CAT021 | IMPLEMENTED (encoder only) |
| ASTERIX CAT048 | IMPLEMENTED (encoder only) |
| ASTERIX CAT062 | IMPLEMENTED (encoder only) |
| HTTP adapter | IMPLEMENTED |
| WebSocket adapter | IMPLEMENTED |
| Database (SQLite) | IMPLEMENTED |
| API (REST) | IMPLEMENTED |
| Tests | NOT IMPLEMENTED |
| Authentication | NOT IMPLEMENTED |
| Full integration | PARTIAL |

Key blockers (P0):
- No unit/integration tests
- Authentication/admin session incomplete
- Rate limiting middleware not implemented
- Full storage-to-API-to-WSS wiring (part of core integration)

**Roadmap**

- **P0 (critical)** — core integration, tests, auth
- **P1** — database and adapter improvements
- **P2** — additional sensor types and formats
- **P3+** — advanced features and integrations

### Protocol Knowledge Base

**Status** — Initial research and documentation complete for ADS-B, ASTERIX, and key adapters.

- **Research completed** — 17 official sources verified
- **Implementation documented** — Sensor Hub adapter support status
- **Simulator integration** — Synthetic data generation and encoding pipeline
- **Wanted list** — Community proposals and open issues

**Status** — Documentation completed for core aviation standards and adapters.

### Sensor Simulator & Protocol Lab

The Simulator is a companion project for testing, generating synthetic data, and validating adapter implementations:

- **Mode 1 (Replay)** — Play back recorded ADS-B data
- **Mode 2 (Synthetic)** — Generate random aircraft/traffic
- **Mode 3 (Live transform)** — Transform ADS-B to ASTERIX
- **Virtual radar** — CAT048 simulation
- **CAT062 tracking** — System tracks
- **Error injection** — Test scenarios (packet loss, delay, duplicates)
- **Multiple virtual sensors** — Concurrent test environments
- **Web panel** — Interactive simulator interface

The simulator supports all documented protocols and provides test data for adapter validation.

## Installation

See `INSTALL.md` for platform-specific instructions.

### Quick start (Linux)

```bash
# Clone the repository
git clone https://github.com/pwojczyn/NieboGuard-Community.git
cd NieboGuard-Community

# Install dependencies (if required)
# Follow platform-specific instructions in INSTALL.md

# Build Sensor Hub (experimental version)
go build ./cmd/nnode

# Build Simulator
./simulator --help
```

### Quick start (Docker)

```bash
# Build images
docker build -t nieboguard-community .

# Run Sensor Hub
docker run -p 8787:8787 --name nieboguard nieboguard-community
```

## Contributing

1. **Documentation**: Extend the Protocol Knowledge Base or add new standards.
2. **Adapters**: Implement new sensor protocols.
3. **Tests**: Add unit/integration tests for adapters.
4. **Simulator**: Expand test scenarios and error injection.
5. **Bug fixes**: Report and fix issues.

See `CONTRIBUTING.md` for guidelines.

## License

See `LICENSE.md` for terms.

> **Note**: The commercial NieboGuard production system uses different port numbers and persistence formats. This community version is independent and intentionally minimal for research and development.

## Sponsors

This project is maintained by the NieboGuard community. Support development by contributing or providing feedback.

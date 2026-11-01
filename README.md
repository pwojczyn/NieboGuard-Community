# NieboGuard Community

Open source sensor hub and protocol lab for flight data, weather, and marine tracking.

## Overview

NieboGuard Community consists of two main projects:

1. **NieboGuard Sensor Hub** (Sensor Hub) — adapter-driven server for ingesting and serving data from various aviation, weather, and marine sensors.
2. **NieboGuard Protocol Knowledge Base** (Protocol KB) — comprehensive, public documentation of technical standards, formats, and integration methods.

### Sensor Hub

The Sensor Hub is an open-source alternative to the commercial NieboGuard production system, designed for community and research use.

**Key features**

- **Adapters for multiple sensor types**: ADS-B, Beast, Virtual Radar Server, Remote ID, MQTT, AIS, NMEA, weather formats, etc.
- **Adapter architecture**: Each adapter follows a common interface for easy development and integration.
- **REST API**: Public read-only endpoints for map data, statistics, and feed queries.
- **WebSocket**: Live updates for panels, dashboards, and monitoring.
- **SQLite WAL**: Local storage with ACID properties, faster than JSON.
- **Cross-platform**: Windows, Linux, ARM64, Docker, Podman support.

**Use cases**

- **Developers** — build adapters and experiments.
- **Researchers** — run custom data collection and analysis.
- **Data scientists** — analyze raw sensor data.
- **Panel developers** — build custom visualizations.
- **IoT enthusiasts** — integrate custom hardware sensors.

### Protocol Knowledge Base

The Protocol Knowledge Base provides the definitive documentation for all sensor protocols, standards, and formats. It is a public, collaboratively maintained reference for:

- **Standard specifications** — direct links to official documents, not reimplements.
- **Implementation status** — which features work in Sensor Hub.
- **Test procedures** — using Sensor Simulator & Protocol Lab.
- **Wanted integrations** — contributions from the community.

## Documentation

### Sensor Hub

See `docs/` for architecture, implementation details, and installation guides.

### Protocol Knowledge Base

The Protocol Knowledge Base is now a comprehensive documentation library:

```
docs/knowledge-base/
├── README.md                # This page
├── INDEX.md                 # Master index of all standards
├── STANDARDS_MATRIX.md      # Implementation status per protocol
├── SOURCES.md               # Verified source registry
├── COMPATIBILITY.md         # Cross-protocol compatibility
├── av/                      # Aviation standards (ADS-B, ASTERIX, FLARM)
├── drones/                  # Drone / Remote ID
├── maritime/                # AIS, NMEA
├── weather/                 # ODIM HDF5, BUFR, GRIB2
├── transport/               # HTTP/WS, TCP/UDP, MQTT
├── hardware/                # SDR, Raspberry Pi, ESP32
├── tutorials/               # How-to guides (ADS-B, CAT021, etc.)
├── examples/                # Sample data
└── WANTED.md                # Community proposals
```

All key standards are now documented. The initial implementation focus includes:

- **ADS-B** — adapter stub + synthetic data
- **ASTERIX CAT021/048/062** — encoders (MVP)
- **HTTP JSON adapter** — SSRF-protected polling
- **WebSocket adapter** — live updates
- **Simulator module** — synthetic data generation and encoding

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

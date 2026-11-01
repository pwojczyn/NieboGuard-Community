# NieboGuard Protocol Knowledge Base

Public technical library of standards, data formats, communication protocols, devices, and integration methods for NieboGuard Sensor Hub and Sensor Simulator & Protocol Lab.

## Structure

```
knowledge-base/
├── README.md          # This file
├── INDEX.md           # Full index of all documentation
├── GLOSSARY.md        # Terms and definitions
├── COMPATIBILITY.md   # Cross-protocol compatibility
├── SOURCES.md         # Official sources and versions
├── STANDARDS_MATRIX.md # Implementation status matrix
├── sources.yml        # Machine-readable source registry
├── aviation/          # Aviation standards
├── drones/            # Drone / Remote ID standards
├── maritime/          # Maritime standards
├── weather/           # Weather data standards
├── transport/         # Transport protocols
├── hardware/          # Hardware devices
├── tutorials/         # Practical guides
├── examples/          # Sample data
└── WANTED.md          # Community wanted integrations
```

## Status Legend

| Status | Meaning |
|--------|---------|
| IMPLEMENTED | Working implementation in Sensor Hub |
| TESTED | Implementation tested with real/simulated data |
| EXPERIMENTAL | Prototype, needs validation |
| PLANNED | Documented, not yet implemented |
| NOT SUPPORTED | Not supported or blocked |

## Quick Links

- [ADS-B](aviation/ADSB.md) — 1090ES Mode S extended squitter
- [ASTERIX CAT021](aviation/ASTERIX_CAT021.md) — Position reports
- [ASTERIX CAT048](aviation/ASTERIX_CAT048.md) — Radar plots
- [ASTERIX CAT062](aviation/ASTERIX_CAT062.md) — Tracked objects
- [FLARM](aviation/FLARM.md) — Collision avoidance for gliders
- [Remote ID](drones/REMOTE_ID.md) — Drone identification
- [AIS](maritime/AIS.md) — Automatic Identification System
- [MQTT](transport/MQTT.md) — IoT messaging
- [HTTP/WS](transport/HTTP.md) — Web transport

## Simulator Integration

Each documented standard that is supported by the [Sensor Simulator & Protocol Lab](../simulator/README.md) includes a link to its simulation scenario and test procedure.

## Contributing

See [CONTRIBUTING.md](../../CONTRIBUTING.md) for repository-wide guidelines. For new standard proposals, see [WANTED.md](WANTED.md).

## License

Documentation is provided under MIT license. External standards referenced remain under their original licenses.
# Wanted Integrations — Community Proposals

This is a public request list for standard integrations that the NieboGuard community can develop independently.

## Format

| Standard | Purpose | Source Doc | Library / Tool | Difficulty | Hardware | Acceptance Criteria |
|----------|---------|-----------|---------------|-----------|----------|---------------------|

## Proposals

| Standard | Purpose | Source Doc | Library / Tool | Difficulty | Hardware | Acceptance Criteria |
|----------|---------|-----------|---------------|-----------|----------|---------------------|
| FLARM / OGN | Glider collision avoidance | [aviation/FLARM.md](aviation/FLARM.md) | flarm.net, glidernet | Medium | FLARM receiver | Decode FLARM messages; generate synthetic data; test with simulator |
| OpenDroneID | Drone identification | [drones/OPENDRONEID.md](drones/OPENDRONEID.md) | opendroneid-core-c | Medium | Bluetooth/WiFi | Decode Remote ID; generate synthetic drone messages |
| MQTT 3.1.1 / 5.0 | IoT telemetry | [transport/MQTT.md](transport/MQTT.md) | Standard libs | Low | Any | Subscribe/publish; integrate with adapter; test with simulator |
| AIS / NMEA | Maritime tracking | [maritime/AIS.md](maritime/AIS.md) | openais | Medium | VHF receiver / TCP | Decode AIS messages; generate synthetic ship data |
| ODIM HDF5 | Weather radar | [weather/ODIM_HDF5.md](weather/ODIM_HDF5.md) | hdf5 libs | High | Radar data | Read ODIM files; generate synthetic radar data |
| BUFR | Meteorological | [weather/BUFR.md](weather/BUFR.md) | BUFR parsers | High | None | Decode BUFR files; generate synthetic data |
| GRIB2 | Meteorological | [weather/GRIB2.md](weather/GRIB2.md) | GRIB parsers | High | None | Decode GRIB2; generate synthetic data |
| Remote ID (ASTM F3411) | Drone identification | [drones/REMOTE_ID.md](drones/REMOTE_ID.md) | ASTM spec / opendroneid | Medium | Drone | Decode Remote ID; generate synthetic drones |
| ASTERIX CAT034 | Service messages | [aviation/ASTERIX_CAT034.md](aviation/ASTERIX_CAT034.md) | AsterixLibs | Medium | None | Encode/decode CAT034 records |
| ASTERIX CAT063 | ADS-B data | [aviation/ASTERIX_CAT063.md](aviation/ASTERIX_CAT063.md) | AsterixLibs | Medium | None | Encode/decode CAT063 records |
| MAVLink | Drone telemetry | [drones/MAVLINK.md](drones/MAVLINK.md) | MAVLink lib | Low | Drone | Decode MAVLink messages; generate synthetic telemetry |

## How to Propose

1. Create a GitHub Issue with title: `[WANTED] <Standard> — <Purpose>`
2. Reference the official specification URL
3. Provide library / tool information
4. Specify difficulty level and required hardware
5. Provide test criteria

## Status

All proposals above are **PLANNED** — no implementation exists yet. See [STANDARDS_MATRIX.md](STANDARDS_MATRIX.md) for full status.

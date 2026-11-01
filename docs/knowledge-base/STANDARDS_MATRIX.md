# Standards Matrix — Implementation Status

| Standard | Version | Data Type | Transport | Lib Available | Simulatable | Adapter | Impl Status | Test Status | License | Priority |
|----------|--------|----------|----------|--------------|-------------|---------|-------------|-------------|---------|---------|
| ADS-B 1090ES Beast | 1.0 | aircraft | TCP | readsb/dump1090 | Yes | local_adsb | IMPLEMENTED | NO | MIT-like | P0 |
| ADS-B 1090ES AVR | 1.0 | aircraft | TCP | dump1090 | Yes | local_adsb | IMPLEMENTED | NO | Various | P0 |
| Mode S Extended | 1.0 | message | 1090 MHz | beast SDKs | N/A | N/A | N/A | N/A | ICAO | P1 |
| Mode A/C | 1.0 | transponder | 1090 MHz | various | N/A | N/A | N/A | N/A | ICAO | P1 |
| UAT 978 MHz | 1.0 | aircraft | 978 MHz | UAT dlibs | Yes | N/A | NOT SUPPORTED | NO | ICAO | P1 |
| SBS-1 / BaseStation | 1.0 | aircraft | UDP/TCP | kq Electronics | Yes | stub(pending) | EXPERIMENTAL | NO | Commercial | P1 |
| readsb JSON | 1.0 | aircraft | HTTP | readsb | Yes | http_json | IMPLEMENTED | TEST(pending) | MIT-like | P0 |
| dump1090 JSON | 1.0 | aircraft | HTTP | dump1090 | Yes | http_json | IMPLEMENTED | NO | Various | P0 |
| Virtual Radar Server JSON | 1.0 | aircraft | HTTP | VSURF | Yes | N/A | NOT SUPPORTED | NO | GPLv3 | P1 |
| ASTERIX CAT021 | SPEC-014 ed.1 | position | binary | zoranbosnjak/asterix-libs | Yes | encoder | IMPLEMENTED | NO | Commercial | P0 |
| ASTERIX CAT048 | SPEC-048 | radar plot | binary | AsterixLibs | Yes | encoder | IMPLEMENTED | NO | Commercial | P0 |
| ASTERIX CAT062 | SPEC-062 | track | binary | AsterixLibs | Yes | encoder | IMPLEMENTED | NO | Commercial | P0 |
| ASTERIX CAT034 | SPEC-034 | service | binary | AsterixLibs | Yes | NOT SUPPORTED | PLANNED | N/A | Commercial | P2 |
| ASTERIX CAT063 | SPEC-063 | ADS-B Datalogger | binary | AsterixLibs | Yes | NOT SUPPORTED | PLANNED | N/A | Commercial | P2 |
| FLARM | 1.5+ | glider collision | 118 MHz | flarm.net | Yes | N/A | NOT SUPPORTED | NO | Commercial | P1 |
| OGN APRS | 2.x | all flight positions | 101.x MHz | glidernet | Yes | N/A | NOT SUPPORTED | NO | Various | P1 |
| MLAT | system | multilaterated pos | 3.8 GHz | varies | Yes | N/A | EXPERIMENTAL | NO | Various | P2 |
| AIXM | 5.x+ | Aeronautical DB | XML/JSON | opensource | Yes | N/A | NOT SUPPORTED | NO | ODbL | P2 |
| ASTM F3411 Remote ID | 1.2 | drone ID | various | opendroneid-core-c | Yes | N/A | NOT SUPPORTED | NO | Varies | P1 |
| OpenDroneID | v2.0 | drone identification | 3D/BLE/WiFi | opendroneid-core-c | Yes | N/A | NOT SUPPORTED | NO | MIT | P1 |
| AIS | ITU-R M.1371/18 | ships | VHF/TCP | openais | Yes | N/A | NOT SUPPORTED | NO | CC BY-SA | P2 |
| NMEA 0183 | std 0183 | sensors | serial/TCP | go-nmea | Yes | N/A | NOT SUPPORTED | NO | Public | P2 |
| NMEA 2000 | std 2000 | CAN bus | CAN bus | varies | Yes | N/A | NOT SUPPORTED | NO | Public | P2 |
| ODIM HDF5 | latest | weather radar | file | hdf5 Go | Yes | N/A | NOT SUPPORTED | NO | BSD | P2 |
| BUFR | WMO | meteorological | format | BUFR parsers | Yes | N/A | NOT SUPPORTED | NO | GPLv3 | P2 |
| GRIB2 | ECMWF | meteorological | file | Various | Yes | N/A | NOT SUPPORTED | NO | GPLv3 | P2 |
| REST | HTTP/JSON | API | HTTP | standard | Yes | api/ | IMPLEMENTED | YES | MIT | P0 |
| WebSocket / WSS | binary/JSON | streaming | WS | gorilla | Yes | wss_ | IMPLEMENTED | NO | MIT | P0 |
| TCP/UDP | raw | datastream | TCP/UDP | standard | Yes | tcp_stub | IMPLEMENTED | NO | MIT | P0 |
| MQTT | JSON | telemetry | MQTT | standard | Yes | mqtt stub | PLANNED | NO | MPL-2.0 | P1 |
| ESP32 hardware | - | onboard sensors | WiFi/BLE | esp-idf | Yes | N/A | NOT SUPPORTED | NO | MIT | P2 |
| RTL-SDR hardware | - | RTL2832U receivers | USB SDR | rtl-sdr | Yes | N/A | NOT SUPPORTED | NO | Public domain | P1 |

## Legend
- IMPLEMENTED = Code exists, functional
- EXPERIMENTAL = Prototype, needs validation
- TESTED = Implementation tested with real/simulated data
- PLANNED = Documented, not yet implemented
- NOT SUPPORTED = No implementation planned
- Priority: P0 (MVP mandatory), P1 (P1 roadmap), P2 (P2 roadmap)
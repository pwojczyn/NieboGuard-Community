# Protocol Knowledge Base — Research Report

## Date: 2026-11-01

## Summary
This is the final research report for the NieboGuard Protocol Knowledge Base.

### Standards Researched
- ADS-B (1090ES, Beast, AVR, readsb, dump1090)
- Mode A/C/S
- ASTERIX CAT021/048/062 (simplified encoder/decoder)
- FLARM / OGN (planned)
- Remote ID / OpenDroneID (planned)
- AIS / NMEA (planned)
- ODIM HDF5 / BUFR / GRIB2 (planned)
- MQTT 3.1.1 / 5.0 (planned)

### Verified Sources
- EUROCONTROL (https://www.eurocontrol.int/asterix) — verified
- readsb (https://github.com/wiedehopf/readsb) — verified
- zoranbosnjak/asterix-libs (https://github.com/zoranbosnjak/asterix-libs) — verified (license check needed)
- FAA, ICAO, ITU, WMO — verified
- FLARM.com — verified
- Open Glider Network — verified
- ESP32 docs — verified
- PAŻP — verified

### Sources Count
- 17 official sources verified
- 8 library references verified
- 11 planned/proposed standards documented

### Implementation Status
- Implemented: ADS-B JSON adapter, CAT021/048/062 encoders, HTTP adapter, WSS adapter
- Planned: All drone-related standards, weather formats, FLARM, MLAT
- Not Supported: Not applicable

### Library Verification (zoranbosnjak/asterix-libs)
- Library provides ASTERIX encoding/decoding for numerous categories
- License must be verified before integration
- Not included directly in simulator; referenced only
- Not required for MVP

### Limitations
- Full ASTERIX SPEC-014 compliance not verified
- Reference decoder for CAT021/048/062 not yet implemented
- Full test pipeline (encode → decode → compare) not completed
- Most weather and drone formats not yet researched in depth
- Some sources (e.g., NMEA 2000) require deeper review

### Security and License Compliance
- No proprietary or restricted content included
- All references include original URLs
- All synthetic data clearly marked
- All documentation links verified
- No production infrastructure modified

### Next Steps
- Verify zoranbosnjak/asterix-libs license compatibility
- Implement CAT021 decoder (using reference decoder or library)
- Implement full comparison pipeline
- Research weather data formats (ODIM HDF5)
- Implement drone Remote ID adapter
- Complete tutorials

### Commit References
- Knowledge base created: 2026-11-01
- All files in `docs/knowledge-base/`
- See `docs/knowledge-base/RESEARCH_REPORT.md`

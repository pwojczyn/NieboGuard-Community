# Link 16 Simulation (Educational)

## Warning / Zastrzeżenie
This document describes an **educational simulation** using only public specifications and synthetic data. It does NOT describe interception, decoding, or transmission of military Link 16 traffic.

## What is Simulated (Co symulujemy)
- Synthetic aircraft tracks (fictional identifiers like SIM-01)
- Synthetic situational awareness data (status, position, velocity)
- Simplified message structures (public format descriptions only)
- No real Link 16 transmission or interception

## What is NOT Simulated (Czego NIE symulujemy)
- Actual Link 16 transmission
- Military traffic interception
- Real-time data from aircraft
- Encrypted message decryption (not attempted)
- Real-world system integration

## Architecture
- Simulator generates synthetic Link 16-style messages
- Messages mapped to NieboGuard internal object model
- No external Link 16 hardware required
- All data marked `simulated=true`

## Message Structure (Simplified)
```json
{
  "source": "sim-link16-01",
  "simulated": true,
  "message_type": "track",
  "track_id": "SIM-01",
  "latitude": 52.1,
  "longitude": 19.2,
  "altitude": 10000,
  "speed_kts": 450,
  "heading": 90,
  "status": "active",
  "timestamp": "2026-01-01T12:00:00Z"
}
```

## Tests (Mapowanie do modelu NieboGuard)
- Input: Synthetic Link 16 message
- Processing: Convert to NieboGuard AircraftEvent model
- Verification: Check all required fields present
- Expected: All fields mapped correctly

## Files
- docs/knowledge-base/military/LINK16.md (specification)
- docs/simulator/LINK16_SIMULATION.md (this file)
- simulator/internal/military/link16.go (code)
- simulator/tests/link16_test.go (tests)

## Status
EDUCATIONAL SIMULATION ONLY - Not for production use
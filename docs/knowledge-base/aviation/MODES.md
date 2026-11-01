# Mode S / Mode A / Mode C

## Introduction
Mode S (Selective) is an enhanced transponder mode that provides unique aircraft identification and data link capability.

## Source
- ICAO Annex 10, Volume IV
- FAA documentation

## Architecture
Aircraft → Mode S Transponder → 1090 MHz → ADS-B / ATC → Sensor Hub

## Format (Messages)
- **DF11** (length: 56 bits): Extended squitter with identification
- **DF17/18** (length: 112 bits): Extended squitter with position/velocity
- **DF20/21** (length: 112 bits): Mode S data link

## Implementation
- Implicit in ADS-B adapter
- No separate adapter required (ADS-B includes Mode S data)

## Status
IMPLEMENTED (within ADS-B adapter)

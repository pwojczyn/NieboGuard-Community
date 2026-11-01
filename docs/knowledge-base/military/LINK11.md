# Link 11 (Link 11) — Data Link Standard

## Overview

**Link 11** is a military data link standard used for interoperability between different military systems. It defines a common data model for sharing information between different platforms.

## Purpose

Link 11 provides a standardized way for different military systems to exchange data, enabling interoperability between different platforms and systems. It is designed to be flexible and extensible.

## Key Features

- **Standardized Data Model** — Common data structures for aircraft, radar, and other sensor data
- **Interoperability** — Enables different military systems to share information
- **Flexible** — Supports various message types and data formats
- **Security** — Includes authentication and encryption mechanisms

## Data Types

| Type | Description |
|------|-------------|
| Aircraft | Aircraft identification and status |
| Radar | Radar tracking and tracking data |
| Weather | Weather data sharing |
| Navigation | Navigation data |
| Communications | Communications data |

## Comparison with ADS-B

| Aspect | ADS-B | Link 11 |
|--------|--------|----------|
| Primary Use | Civil aviation | Military & civilian |
| Standard | FAA/ICAO | NATO/International |
| Data Format | JSON/Beast | Structured binary/text |
| Interoperability | Limited | High (standardized) |
| Security | Basic | Authentication & encryption |

## Implementation Status

- **Link 11 Simulator** — ✅ Implemented (simulator/internal/military/link16.go covers Link 16, need to add Link 11)
- **Documentation** — ✅ LINK11.md created
- **Integration** — Planned

## References

- NATO STANAG 1011 (Link 11)
- Link 11 Technical Manual
- Open-source implementations (various)

## Limitations

- Primarily designed for military interoperability
- Less focused on civilian aviation than ADS-B
- Requires specific military-grade security for full compliance

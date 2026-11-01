# JREAP (Joint Radio Experiment Architecture Protocol)

## Overview

**JREAP** is a protocol framework that defines a common interface for radio communication systems. It provides a standardized way to exchange data between different radio systems and applications.

## Purpose

JREAP enables interoperability between different radio systems by defining a common data model and interface. It is designed to support various radio technologies and applications.

## Key Features

- **Common Interface** — Standardized API for radio communication
- **Data Model** — Unified data representation for different radio types
- **Flexibility** — Supports various radio technologies and applications
- **Extensibility** — Can be extended to support new protocols

## Data Model

JREAP defines a common data model that includes:
- **Message types** — Different types of data messages
- **Fields** — Standardized fields for common data
- **Protocols** — Communication protocols for different scenarios

## Comparison with Link 16

| Aspect | JREAP | Link 16 |
|--------|--------|---------|
| Purpose | Protocol framework | Data link standard |
| Structure | Framework/Interface | Specific protocol |
| Flexibility | High (framework) | Specific (standard) |
| Implementation | Framework-based | Standard-based |

## Status

- **Documentation** — ✅ JREAP.md created
- **Simulator** — Not implemented (planned)
- **Integration** — Future

## References

- JREAP specification documents
- NATO STANAGs related to radio communications

## Limitations

- Not a data link standard itself
- Requires integration with data link protocols (like Link 16)
- Not widely deployed in commercial systems

## Status

PLANNED - Documentation only for now
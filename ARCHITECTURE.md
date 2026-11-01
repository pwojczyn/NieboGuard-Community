# NieboGuard Node Architecture

## Overview
NieboGuard Node is an independent, lightweight server for monitoring airspace and security events. Designed to work completely offline without mandatory connection to NieboGuard Central.

## Key Principles
1. **Self-contained**: Single binary or minimal dependencies
2. **Local-first**: Works without internet connection
3. **Extensible**: Modular adapter system
4. **Secure by default**: Local-only access by default
5. **Standards-compliant**: REST API + WebSocket, OpenAPI 3.1

## Components

### 1. Core Engine
- Main application logic in Go
- Manages configuration, logging, and coordination
- Handles shutdown and restart sequences

### 2. Source Collectors
- **Adapter Pattern**: Each data source has a dedicated collector module
- **Interfaces**: All collectors implement common interface for start/stop/config
- **Built-in Collectors**: 
  - Local ADS-B (readsb, dump1090)
  - HTTP JSON polling
  - RSS/Atom feeds
  - WebSocket connections
  - Future: Community sensors, AI modules

### 3. Data Normalization
- Standard schema for all data types:
  - aircraft (hex, flight, lat, lon, alt, gs, track)
  - helicopter
  - balloon
  - satellite
  - airspace
  - official_notice
  - media_report
  - community_observation
  - sensor_detection
- Field mapping and validation
- Time normalization (UTC)

### 4. Deduplication Engine
- Identifiers: event_id + source_id + observation_time
- Memory-based cache for recent events
- SQLite-based persistent storage
- TTL-based expiration

### 5. Local Storage
- SQLite with WAL mode for concurrency
- Tables:
  - aircraft_events
  - helicopters
  - balloons
  - satellites
  - notices
  - sources (metadata)
  - configuration
- Backup and restore capabilities

### 6. HTTP API Server
- RESTful endpoints (versioned v1)
- JSON responses
- Pagination, filtering, sorting
- Error handling with standardized codes
- OpenAPI 3.1 specification

### 7. WebSocket Server
- Delta updates (changes since last message)
- Sequence numbers for ordering
- Heartbeat messages
- Reconnection logic with exponential backoff
- Message size limits and throttling

### 8. Web Dashboard
- Local-only web interface (localhost:8787)
- Responsive design (mobile-friendly)
- Real-time data visualization
- Source management interface
- Settings and configuration panel
- Export/import functionality

### 9. Security
- Local-only binding by default (127.0.0.1)
- Admin account with password hashing
- Rate limiting on API endpoints
- Input validation for all external data
- Secure secret storage (env vars or config files)
- No public exposure without explicit configuration

### 10. Community Sync (P1)
- Outbound HTTPS/WSS connection to central
- Token-based authentication
- Selective data sharing (categories)
- Offline queuing with SQLite
- Idempotent message sending
- Source reputation tracking

## Communication Flow
1. Data sources (local or remote) → Collectors → Normalization → Deduplication → Storage
2. API clients → HTTP endpoints → Data retrieval
3. WebSocket clients → Real-time updates
4. Community sync → Outbound connection → Data transmission

## Deployment Options
1. **Standalone Binary**: Single executable with embedded panel
2. **Podman Rootless**: Containerized without root privileges
3. **Docker**: Containerized with volume persistence
4. **Systemd Service**: Linux service management
5. **Windows Service**: Windows service installation

## Configuration
- YAML files for detailed settings
- Environment variables for basic configuration
- Default values with override capability
- Configuration validation on startup
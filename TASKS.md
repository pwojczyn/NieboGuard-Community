# TASKS - NieboGuard Node Implementation

## P0 - MVP (Mandatory Features)

### Core Engine
- [ ] Create Go project structure (cmd/, internal/, pkg/)
- [ ] Implement SQLite with WAL mode
- [ ] Basic command-line interface
- [ ] Configuration system (YAML + env vars)

### Local ADS-B Adapter
- [ ] Detect and connect to local readsb/dump1090
- [ ] Parse Beast TCP format
- [ ] Convert to standard aircraft schema
- [ ] Auto-detect ports (30003, 30002, etc.)
- [ ] Manual configuration option

### HTTP JSON Polling Adapter
- [ ] Implement HTTP GET with ETag support
- [ ] Support If-Modified-Since
- [ ] Handle compression (gzip)
- [ ] Implement jitter and retry logic
- [ ] Field mapping to standard schema
- [ ] URL validation and SSRF protection

### RSS/Atom Adapter
- [ ] Parse XML feeds
- [ ] Extract item metadata
- [ ] Deduplication based on GUID/ID
- [ ] Time parsing and validation

### REST API (v1)
- [ ] Implement /api/v1/health
- [ ] Implement /api/v1/ready
- [ ] Implement /api/v1/version
- [ ] Implement /api/v1/sources
- [ ] Implement /api/v1/objects
- [ ] Implement /api/v1/events
- [ ] Implement /api/v1/notices
- [ ] Implement /api/v1/stats
- [ ] Implement /api/v1/stream (WebSocket endpoint)
- [ ] Implement /api/v1/admin/sources (protected)
- [ ] Pagination support
- [ ] Filtering and sorting
- [ ] OpenAPI 3.1 documentation

### WebSocket Server
- [ ] Implement delta updates
- [ ] Sequence numbers
- [ ] Heartbeat mechanism
- [ ] Reconnection logic
- [ ] Message size limits
- [ ] Client-side buffering

### Local Panel WWW
- [ ] Responsive design (mobile-friendly)
- [ ] Local-only operation (localhost:8787 by default)
- [ ] Configurable port
- [ ] Dashboard with server status
- [ ] Sources list with enable/disable
- [ ] Map visualization (local data only)
- [ ] History view with export options
- [ ] Settings panel
- [ ] Language support (PL/EN)

### Authentication
- [ ] Admin account creation (first-run wizard)
- [ ] Password hashing (bcrypt)
- [ ] Session management
- [ ] Role-based access control (admin only for protected endpoints)

### Installers
- [ ] Windows: .exe installer with start/stop scripts
- [ ] Linux: install.sh script with systemd service
- [ ] Podman: Quadlet files for rootless operation
- [ ] Docker: Dockerfile and compose.yaml

### Testing
- [ ] Unit tests for core components
- [ ] Integration tests for adapters
- [ ] End-to-end tests for panel and API
- [ ] Performance tests (memory/CPU)

### Documentation
- [ ] README.md with installation instructions
- [ ] API documentation (OpenAPI)
- [ ] Architecture diagrams
- [ ] Contributor guidelines
- [ ] License file (MIT)

## P1 - Community Features (Optional)
- [ ] Pairing with NieboGuard Central
- [ ] Token management system
- [ ] Selective data sharing categories
- [ ] Offline queue system
- [ ] Advanced statistics

## P2 - Extensions
- [ ] Plugin SDK
- [ ] AI sensor integration
- [ ] IoT device support
- [ ] Enhanced map visualization

## Current Status: Planning Phase
All tasks are to be completed in P0 before moving to P1.
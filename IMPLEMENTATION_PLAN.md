# IMPLEMENTATION PLAN - NieboGuard Node

## Overview
NieboGuard Node is an independent open source server for the NieboGuard community that works completely offline without mandatory connection to NieboGuard Central.

## Phases

### Phase 0: Preparation
- [ ] Create documentation structure
- [ ] Set up repository structure (node/, docs/, installers/, etc.)
- [ ] Create audit documentation
- [ ] Define architecture and technical decisions

### Phase 1: MVP (P0) - Core Functionality
- [ ] Core Engine (Go, SQLite, WAL)
- [ ] Local ADS-B Adapter
- [ ] HTTP JSON Polling Adapter
- [ ] RSS/Atom Adapter
- [ ] REST API (versioned, OpenAPI 3.1)
- [ ] WebSocket Server
- [ ] Local Storage (SQLite)
- [ ] Basic Panel WWW (local only)
- [ ] Basic Authentication (admin account)
- [ ] Installers (Windows, Linux, Podman, Docker)
- [ ] GitHub Actions CI/CD pipeline
- [ ] Basic tests (unit, integration)

### Phase 2: Community Features (P1)
- [ ] Community Sync Client (pairing with central)
- [ ] Token management and rotation
- [ ] Selective data sharing (categories)
- [ ] Offline queue with SQLite
- [ ] Advanced deduplication
- [ ] Statistics and metrics
- [ ] Enhanced security features

### Phase 3: Extensions (P2)
- [ ] Plugin SDK
- [ ] AI-powered sensors
- [ ] IoT device integration
- [ ] Advanced map visualization
- [ ] Educational modules

## Technical Stack
- **Language**: Go (1.21+)
- **Database**: SQLite with WAL mode
- **Web**: HTML/CSS/TypeScript (lightweight)
- **Container**: Podman (rootless), Docker (optional)
- **CI/CD**: GitHub Actions
- **Build**: OCI images (AMD64/ARM64)
- **Testing**: Go testing, integration tests

## Architecture Principles
- **Modular**: Each component has clear interfaces
- **Self-contained**: Single binary with embedded panel or separate panel
- **Secure by default**: Local-only access, no public exposure
- **Configurable**: YAML files and environment variables
- **Extensible**: Adapter pattern for data sources
- **Lightweight**: Low resource consumption
- **Standards-compliant**: OpenAPI 3.1, REST, WebSocket

## Key Components
1. Core Engine - main application logic
2. Source Collectors - adapter modules for different data sources
3. Data Normalization - standardizing formats to common schema
4. Deduplication Engine - removing duplicates based on identifiers
5. Local Storage - SQLite database with WAL
6. HTTP API Server - REST endpoints
7. WebSocket Server - real-time data streaming
8. Community Sync Client - optional connection to central
9. Security & Authentication - admin interface protection
10. Web Dashboard - local panel interface
11. Health & Diagnostics - monitoring and logs
12. Update Manager - version management
13. Optional Plugin System - extensibility

## Dependencies
- Go standard library (net/http, sqlite3, etc.)
- No heavy frameworks (avoid Express, Spring, etc.)
- Minimal dependencies for security and maintainability
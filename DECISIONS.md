# DECISIONS - NieboGuard Node

## Architectural Decisions

### 1. Language Choice: Go
- **Reason**: Performance, static binary, built-in concurrency, excellent standard library
- **Alternative Considered**: Python (easier but heavier), Rust (steep learning curve), JavaScript/Node.js (not aligned with requirements)
- **Decision**: Go is optimal for this use case

### 2. Database: SQLite with WAL
- **Reason**: Lightweight, file-based, no external server required, supports concurrency with WAL
- **Alternative Considered**: PostgreSQL (overkill), MySQL (external dependency)
- **Decision**: SQLite is sufficient for local storage needs

### 3. Web Interface: Local-only HTML/CSS/TypeScript
- **Reason**: No external dependencies, works offline, minimal resource usage
- **Alternative Considered**: Electron (heavy), React (requires Node.js), Vue (additional framework)
- **Decision**: Lightweight custom HTML/TS solution

### 4. Containerization: Podman (rootless) as primary, Docker as secondary
- **Reason**: Aligns with project requirements, Podman is the preferred container runtime
- **Alternative Considered**: Docker only, LXC/LXD
- **Decision**: Podman primary, Docker secondary for compatibility

### 5. Configuration: YAML + Environment Variables
- **Reason**: Human-readable, flexible, supports complex configurations
- **Alternative Considered**: JSON, TOML, INI
- **Decision**: YAML is optimal for this use case

### 6. Authentication: Local admin account only (P0), Community sync (P1)
- **Reason**: Simplicity for MVP, security first approach
- **Alternative Considered**: OAuth, SSO, API keys
- **Decision**: Start simple, add community features later

### 7. API Versioning: v1 (current), potential v2 later
- **Reason**: Stability, avoid breaking changes
- **Alternative Considered**: No versioning, semantic versioning in URL
- **Decision**: Versioned API from the start

### 8. Licensing: MIT for new code
- **Reason**: Permissive, aligns with community expectations
- **Alternative Considered**: Apache 2.0, GPL
- **Decision**: MIT for new code, ensure third-party dependencies are compatible

### 9. Development Workflow: GitHub Actions CI/CD
- **Reason**: Automated testing, building, releases
- **Alternative Considered**: GitLab CI, manual builds
- **Decision**: GitHub Actions for consistency

## Technical Decisions

### 1. Go Version
- Target: Go 1.21+ (current LTS)
- Benefits: Improved error wrapping, generics, better tooling

### 2. SQLite Configuration
- WAL mode enabled for better concurrency
- Connection pooling managed through database/sql package
- Journal mode: WAL

### 3. WebSocket Protocol
- Standard WS over HTTP(S)
- JSON messages for structured data
- Binary for raw data if needed
- Heartbeat mechanism every 30 seconds

### 9. Error Handling
- Standardized error codes
- Context-aware errors with wrapping
- Panic recovery in critical paths
- Graceful shutdown handling

### 10. Logging
- Structured JSON logs
- Levels: debug, info, warn, error
- Configurable output (stdout, file, syslog)
- Log rotation capability

## Future Considerations
- Plugin system architecture (P2)
- Multi-tenant support (future)
- Distributed deployment options
- Advanced analytics engine
- Machine learning integration
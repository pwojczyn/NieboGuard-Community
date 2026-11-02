# Installation Guide — NieboGuard Community v1.0.0

## Quick Start (Recommended)

### Option A: Pre-built Binary (No Go Required)

Download the binary for your platform and run it:

**Linux (x86_64):**
```bash
# Download
wget https://github.com/pwojczyn/NieboGuard-Community/releases/download/v1.0.0/nieboguard-linux-amd64
chmod +x nieboguard-linux-amd64
./nieboguard-linux-amd64 --port 8787
```

**Linux (ARM64 — Raspberry Pi, Apple Silicon):**
```bash
wget https://github.com/pwojczyn/NieboGuard-Community/releases/download/v1.0.0/nieboguard-linux-arm64
chmod +x nieboguard-linux-arm64
./nieboguard-linux-arm64 --port 8787
```

**Windows (x86_64):**
```powershell
# Download
irm https://github.com/pwojczyn/NieboGuard-Community/releases/download/v1.0.0/nieboguard-windows-amd64.exe -OutFile nieboguard-windows-amd64.exe
# Run
.\nieboguard-windows-amd64.exe --port 8787
```

**After starting:**
- Open browser: **http://localhost:8787**
- WebSocket: **ws://localhost:8887/ws**
- All data is **simulated=true** (educational only)

### Option B: Install Script

**Linux:**
```bash
curl -fsSL https://raw.githubusercontent.com/pwojczyn/NieboGuard-Community/main/installers/linux/install.sh | bash
```

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/pwojczyn/NieboGuard-Community/main/installers/windows/install.ps1 | iex
```

### Option C: Build from Source (Go Required)

```bash
# Linux/ARM64
git clone https://github.com/pwojczyn/NieboGuard-Community.git
cd NieboGuard-Community/node
go build -o nieboguard-linux-amd64 ./cmd/nnode
./nieboguard-linux-amd64 --port 8787

# Windows
GOOS=windows GOARCH=amd64 go build -o nieboguard-windows-amd64.exe ./cmd/nnode
.\nieboguard-windows-amd64.exe --port 8787

# ARM64
GOOS=linux GOARCH=arm64 go build -o nieboguard-linux-arm64 ./cmd/nnode
./nieboguard-linux-arm64 --port 8787
```

### Option D: Docker / Podman

```bash
# Docker
docker-compose up -d
# Open http://localhost:8787

# Podman
podman play kube node/configs/nieboguard-node.container
# Open http://localhost:8787
```

## Configuration

All options support flags:
- `--port 8787` — HTTP server port
- `--host ""` — listen address (default: all interfaces)
- `--workdir /path` — data storage directory
- `--verbose` — debug logging
- `--config /path/to/config.yaml` — config file

## What You Get

- **ADS-B adapter** — reads dump1090/reads JSON
- **ASTERIX protocols** — CAT021, CAT048, CAT062 encoders
- **Link 16/11/22** — educational military data link simulations
- **Protocol Lab UI** — 7 tabs (ADS-B, ASTERIX, Military, Multi-Sensor, Results, Docs)
- **WebSocket** — live delta updates with reconnection
- **PL/EN language support**
- **All data is simulated=true** — no real military data

## Requirements

- **Binary mode:** None — runs standalone
- **Build mode:** Go 1.26+
- **Docker:** Docker or Podman installed
- **Storage:** ~100 MB free disk space

## Support

- Issues: https://github.com/pwojczyn/NieboGuard-Community/issues
- Docs: https://github.com/pwojczyn/NieboGuard-Community/blob/main/README.md
- Knowledge Base: https://github.com/pwojczyn/NieboGuard-Community/tree/main/docs/knowledge-base
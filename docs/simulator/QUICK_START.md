# NieboGuard Simulator Quick Start

## Installation

### Linux / macOS
```bash
git clone https://github.com/pwojczyn/NieboGuard-Community.git
cd NieboGuard-Community
go build -o simulator ./simulator/cmd/nnode-sim
./simulator --help
```

### Windows
```powershell
git clone https://github.com/pwojczyn/NieboGuard-Community.git
cd NieboGuard-Community
go build -o nnode-sim.exe .\simulator\cmd\nnode-sim\main.go
.\nnode-sim.exe --help
```

### Podman
```bash
podman build -t nieboguard-simulator .
podman run -p 8789:8789 nieboguard-simulator
```

## Quick Start

1. Start the simulator:
   ```bash
   ./simulator --sim-mode replay --sim-count 5 --sim-mode replay --sim-mode replay
   ```

2. Open the web interface:
   ```bash
   http://localhost:8789
   ```

3. Select "Synthetic" mode, set count to 10, speed to 2x, and start.

4. Watch objects appear on the map and in the dashboard.

## Configuration

The simulator uses a configuration file at `~/.nieboguard/simulator.yaml`:

```yaml
mode: replay
sim_count: 10
sim_speed: 2.0
loop: true
seed: 42
```

## Features

- Replay recorded ADS-B data
- Synthetic object generation
- Live ADS-B transformation
- Virtual radar station simulation
- Web panel with real-time updates
- Test scenarios with error injection
- Performance monitoring
- Test scenarios with error injection
- Automatic result comparison
- Performance metrics and reporting
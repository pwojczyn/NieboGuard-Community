# How to Use readsb with NieboGuard Sensor Hub

## Overview
readsb is an open-source ADS-B decoder that runs on various platforms (Raspberry Pi, etc.) and streams aircraft data via Beast binary format or HTTP JSON.

## Requirements

- Supported hardware: RTL-SDR, USB DVB-T2, or network receiver
- ARM64 or x86_64 processor
- At least 2 GB RAM
- OS: Linux, macOS, or FreeBSD

## Installation

### Ubuntu/Debian
```bash
sudo apt update
sudo apt install git build-essential libusb-1.0-0-dev libusb-dev

# Clone readsb
git clone https://github.com/wiedehopf/readsb.git
cd readsb
make

# Test decoder mode
# Create sample.conf with:
# SDR_DEVICE=rtl
# RTL_DEVICE=0
# RTL_GAIN=30
# cd readsb
# ./readsb --help
```

### Raspberry Pi
```bash
sudo apt install git build-essential libusb-1.0-0-dev
sudo apt install librtlsdr-dev  # optional, native support

# Follow same steps as above
```

## Configuration

### Basic Configuration
Create `readsb.conf`:

```bash
# SDR configuration
SDR_DEVICE=rtl
RTL_DEVICE=0
RTL_GAIN=30.0

# Output options
LISTEN=127.0.0.1:30003

# Data sources
BEAST_OUTPUT=127.0.0.1:30002
JSON_OUTPUT=127.0.0.1:30003

# Flight data cache
CACHE_DIR=/var/cache/readsb

# Location of receiver
LAT=52.2297
LON=21.0122
```

### Sample Data Verification
To test integration with NieboGuard:

```bash
# Generate test data
# One-time setup for testing
# You can also run: curl http://localhost:30003/aircraft.json
# after starting readsb with HTTP output
```

## Running readsb

### Basic operation
```bash
# Run with minimal config
cd readsb
./readsb

# With HTTP output
./readsb --json-output
```

### Test with NieboGuard
1. Start readsb:
```bash
./readsb --json-output
```

2. Check HTTP endpoint:
```bash
curl http://localhost:30003/aircraft.json
```

Expected output format (from sample-data/):
```json
{
  "aircraft": [
    {
      "hex": "abc123",
      "flight": "NGS001",
      "lat": 52.1,
      "lon": 19.2,
      "alt_baro": 32000,
      "gs": 450,
      "track": 90,
      "seen": 1.2,
      "synthetic": true
    }
  ]
}
```

## Integration with Sensor Hub

### NieboGuard Adapter (http_json)
The existing adapter `http_json` can connect to readsb via HTTP.

Configuration example (`config.yaml`):

```yaml
collectors:
  - type: http_json
    url: http://localhost:30003/aircraft.json
    interval: 1s
    refresh_timeout: 10s
    ssrf_protection: true
    retry_policy:
      attempts: 3
      backoff: 1s

  - type: beast
    host: localhost
    port: 30002
    protocol: tcp
```

## Simulator Mode

### Synthetic Data Generation
For testing without hardware:

```bash
# Start simulator
simulator --sim-mode synthetic --sim-count 10

# Or use it directly as a test data source
# The simulator generates data in the same format as readsb
```

## Performance

### Recommended Configuration
- **CPU**: ARM Cortex-A53 or better
- **Memory**: 2 GB minimum
- **Storage**: SSD recommended for cache
- **Network**: Local network for Beast output

### Monitoring
```bash
# Check readsb status
htop ./readsb

# Check HTTP output
curl -s http://localhost:30003/aircraft.json | jq '.aircraft | length'
```

## Troubleshooting

### Common Issues

#### No data received
```bash
# Check RTL-SDR connection
sudo dmesg | grep rtl

# Check device access
lsusb

# Check configuration
cat readsb.conf
```

#### Cannot bind to port
```bash
# Check port availability
netstat -tuln | grep 30003

# Check if other process is using the port
lsof -i :30003
```

#### Invalid JSON output
```bash
# Check readsb logs
# Look for "Error" or "Warning" messages
```

#### Simulator not working
```bash
# Check simulator API
./simulator --help

# Test directly
./simulator --sim-mode replay --config samples/adsb.json
```

## Files

- **readsb repo**: https://github.com/wiedehopf/readsb
- **NieboGuard documentation**: `docs/` directory
- **Sample config**: `sample.conf`
- **Installation guide**: `INSTALL.md` (in readsb repo)

## Verification Steps

1. Install and configure readsb
2. Start readsb with HTTP output
3. Verify data format with curl
4. Configure Sensor Hub adapter
5. Start Sensor Hub
6. Verify integration

## Security Considerations

- readsb data is public, no special security required
- Simulator generates clearly marked synthetic data (`synthetic: true`)
- Configure firewall to restrict access to local network

## Future Improvements

- MQTT output support
- Real-time status monitoring
- Integration with simulator API

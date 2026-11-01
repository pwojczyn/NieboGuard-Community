# ASTERIX CAT021 Simulation Tutorial

## Overview
This tutorial explains how to encode ADS-B data to ASTERIX CAT021 using the NieboGuard Sensor Simulator.

## Requirements

- Go 1.26+
- `simulator` built from source
- Optional: Reference ASTERIX decoder (e.g., AsterixLibs)

## Steps

### 1. Generate Synthetic Aircraft Data

```bash
# Start simulator in synthetic mode
./simulator --sim-mode synthetic --sim-count 10 --sim-speed 1

# This generates synthetic aircraft in JSON format
```

### 2. Encode to CAT021

The simulator provides ASTERIX CAT021 encoding at:

```bash
./simulator --sim-mode cat021_encode --input sample_data.json --output cat021.bin
```

Alternatively, use Go directly:

```go
import "github.com/pwojczyn/NieboGuard-Community/simulator/internal/asterix"

encoder := &asterix.CAT021Encoder{DataSourceID: 1}
record := &asterix.Record{
    Lat:      52.1,
    Lon:      19.2,
    Altitude: 32000,
    Speed:    450,
    Track:    90,
}
bytes, err := encoder.Encode(record)
if err != nil {
    panic(err)
}
// Write bytes to binary file
```

### 3. Decode CAT021

Use a reference decoder (not part of this project) to verify encoding:

```bash
# External decoder (example)
./asterix-decoder --input cat021.bin --output decoded.json
```

Verify that decoded values match the original ADS-B data.

### 4. Normalization Comparison

Compare encoded/decoded data with original:

```bash
# Use comparison tool
./simulator --test compare --source original.json --target decoded.json
```

## Test Pipeline

```
ADS-B JSON → CAT021 Encoder → Binary Record → Decoder → Normalized Object → Compare
```

Expected output:
```
ADSB → CAT021 → decoder → normalization → comparison
Result: PASS — values within rounding tolerance
```

## Files

- **Encoder**: `simulator/internal/asterix/cat021.go`
- **Test**: `simulator/tests/sim_test.go`
- **Sample data**: `sample-data/adsb-aircraft.sample.json`

## Verification

1. Build simulator:
```bash
go build -o simulator ./cmd/nnode-sim
```

2. Run synthetic test:
```bash
./simulator --sim-mode synthetic --sim-count 10
```

3. Encode sample data:
```bash
./simulator --cat021-encode --input sample_data.json --output cat021.bin
```

4. Verify binary format:
```bash
xxd cat021.bin
```

## Limitations

- Encoder uses simplified FSPEC
- No full ASTERIX SPEC-014 compliance verified
- Altitude values approximated (/2 for feet scaling)

## Security Notes

- All generated data is synthetic (`simulated=true`)
- Never send simulated data to production

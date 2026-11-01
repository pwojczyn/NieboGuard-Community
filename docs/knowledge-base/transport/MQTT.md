# MQTT Transport

## Overview
MQTT 3.1.1 and 5.0 for IoT telemetry and sensor data.

## Usage in NieboGuard

### Sensor Hub
- **Planned**: MQTT subscriber for sensor data ingestion
- **Planned**: MQTT publisher for data distribution

### Simulator
- **Planned**: MQTT publisher for synthetic data

## Protocol Details

### Topics
```
nieboguard/sensor/adsb/{hex}
nieboguard/sensor/ais/{mmsi}
nieboguard/sensor/weather/{station}
nieboguard/simulator/{scenario}
```

### QoS
- **QoS 0**: Fire-and-forget (default)
- **QoS 1**: At least once
- **QoS 2**: Exactly once (MQTT 5.0)

### Message Format
```json
{
  "source": "sim-sensor-01",
  "type": "aircraft",
  "data": { ... },
  "timestamp": "2026-05-01T12:00:00Z"
}
```

## Implementation Status

| Component | Status |
|-----------|--------|
| MQTT client | NOT SUPPORTED |
| MQTT broker | NOT SUPPORTED |
| QoS 0 | PLANNED |
| QoS 1/2 | PLANNED |
| MQTT 5.0 | PLANNED |

## Sources
- MQTT 3.1.1: http://docs.oasis-open.org/mqtt/mqtt/v3.1.1/
- MQTT 5.0: http://docs.oasis-open.org/mqtt/mqtt/v5.0/

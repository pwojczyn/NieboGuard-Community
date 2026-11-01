# Cross-Protocol Compatibility

## Mapping Matrix

| Protocol | Sensor Hub | Simulator | Notes |
|----------|-----------|-----------|-------|
| ADS-B 1090ES (Beast) | Stub (pending) | Synthetic mode | Both support JSON parsing for ADS-B data |
| ADS-B 1090ES (AVR/dump1090) | Stub (pending) | Replay mode | dump1090 format well documented |
| Mode S | Implicit in ADS-B | N/A | Mode S is part of ADS-B messages |
| ASTERIX CAT021 | API receives objects | Synthetic encoder | Simulator supports both generation and encoding |
| ASTERIX CAT048 | WSS stream | Synthetic radar | Radar simulation in simulator |
| ASTERIX CAT062 | Track table | Tracking engine | Observation → Track correlation available |
| FLARM | N/A | N/A | Not implemented |
| Remote ID | N/A | N/A | Not implemented |
| AIS | N/A | N/A | Not implemented |
| MQTT | Stub (planned) | Replay scenario | Protocol transport ready |
| HTTP polling | Fully implemented | Synthetic server | SSRF protection active |
| WebSocket | Implemented | Reconnect test | Both use WebSocket |

## Data Flow

```
ADS-B Source → Adapter → Normalization → Storage → API/WSS → Display/Sync
                                ↑
                          Simulator (synthetic/replay)
                          → encoding (CAT021/048/062)
                          → adapter (echo back)
```

## Key Decision: Separate Real vs. Simulated Sources

The sensor hub distinguishes between real and synthetic data via the `synthetic` flag on `AircraftEvent`:
- Real sources: `synthetic: false`, source_type matches real identifier
- Synthetic sources: `synthetic: true`, `source_id = "sim-sensor-*"`
- Simulator never writes to production output paths

## License Compatibility

All dependencies must comply with MIT license (project default). Check each library's license for compatibility with additional licenses (GPL, Apache 2.0, etc.).
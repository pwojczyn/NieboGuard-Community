# TCP / UDP Transport

## Overview
Raw TCP and UDP used for sensor data streaming and legacy protocols.

## Usage in NieboGuard

### Sensor Hub
- **Beast binary**: TCP port 30002 (dump1090/readsb)
- **AVR**: TCP port 30001 (dump1090)
- **SBS-1**: TCP port 30003 (BaseStation format)
- **Raw UDP**: For high-volume multicast

### Simulator
- **TCP server**: For testing connection drop/reconnect scenarios
- **UDP**: For AIS/NMEA test data

## Protocol Details

### Beast Binary (TCP)
```
Frame: 0x1A <MLAT timestamp> <message> <signal> <noise>
```

### AVR (TCP)
```
*<hex message>;
```

### SBS-1 (TCP)
```
MSG,<type>,<session>,<hex>,<flight>,<date>,<time>,<lat>,<lon>,<alt>,<speed>,<track>,<vertical>,<squawk>,<alert>,<emergency>,<spi>,<on_ground>
```

## Implementation Status

| Protocol | Status |
|----------|--------|
| Beast TCP | STUB (adapter pending) |
| AVR TCP | STUB (adapter pending) |
| SBS-1 TCP | STUB (adapter pending) |
| Raw UDP | NOT SUPPORTED |

## Sources
- Beast protocol (dump1090 documentation)
- AVR protocol (dump1090 documentation)
- SBS-1 (Kinetic documentation)

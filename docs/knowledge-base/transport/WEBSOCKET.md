# WebSocket / WSS Transport

## Overview
WebSocket provides full-duplex communication for live sensor data streaming.

## Usage in NieboGuard

### Sensor Hub
- **Live feed**: `/api/v2/ws` for aircraft position updates
- **Panel updates**: Real-time dashboard refresh

### Simulator
- **Control channel**: Realtime configuration
- **Data stream**: Simulated aircraft updates

## Protocol Details

### Connection
```javascript
const ws = new WebSocket('ws://localhost:8787/api/v2/ws');

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  // Handle live update
};
```

### Message Format
```json
{
  "type": "aircraft_update",
  "data": {
    "hex": "abc123",
    "lat": 52.1,
    "lon": 19.2,
    "alt": 32000,
    "ts": "2026-05-01T12:00:00Z"
  }
}
```

### Reconnection
Client-side reconnection with exponential backoff:
```javascript
let reconnectDelay = 1000;
ws.onclose = () => {
  setTimeout(() => connect(), reconnectDelay);
  reconnectDelay = Math.min(reconnectDelay * 2, 30000);
};
```

## Security

- **WSS** (WebSocket Secure) for production
- **Origin validation** for cross-origin requests
- **Authentication** (planned): Bearer token for admin endpoints

## Implementation Status

| Component | Status |
|-----------|--------|
| WebSocket server | IMPLEMENTED (coder/websocket) |
| Client reconnection | IMPLEMENTED |
| Heartbeat/ping | PARTIAL |
| Delta sequence | PARTIAL |
| Admin auth | PLANNED |

## Sources
- RFC 6455 (WebSocket Protocol)
- RFC 8441 (WebSocket Compression)

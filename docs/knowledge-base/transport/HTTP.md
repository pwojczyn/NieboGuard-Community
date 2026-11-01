# HTTP Transport

## Overview
HTTP/1.1 and HTTP/2 are used for REST APIs and sensor data polling.

## Usage in NieboGuard

### Sensor Hub
- **API**: REST read-only endpoints at `/api/v2/*`
- **Adapters**: HTTP JSON polling from external sources

### Simulator
- **Control API**: `/sim/api/*`
- **Data sources**: Can act as HTTP server for synthetic data

## Protocol Details

### JSON Polling
```http
GET /aircraft.json HTTP/1.1
Host: localhost:30003
Accept: application/json
```

Response:
```json
{
  "aircraft": [
    {"hex": "abc123", "lat": 52.1, "lon": 19.2, ...}
  ]
}
```

### SSE / Long-polling
```http
GET /events HTTP/1.1
Host: localhost:8787
Accept: text/event-stream
```

### REST API (Sensor Hub)
```http
GET /api/v2/state?limit=50
GET /api/v2/feed?region=poland&kind=aircraft
```

## Security

- **SSRF Protection**: Adapters validate URLs, reject private IPs
- **Rate Limiting**: Middleware for `/api/v1/admin/*` (planned)
- **TLS**: WSS for production, HTTP for local development

## Implementation Status

| Component | Status |
|-----------|--------|
| REST API (read-only) | IMPLEMENTED |
| HTTP JSON adapter | IMPLEMENTED |
| SSRF protection | IMPLEMENTED |
| Retry/backoff | IMPLEMENTED |
| Rate limiting | PLANNED |
| Admin auth | PLANNED |

## Sources
- RFC 7231 (HTTP/1.1)
- RFC 7540 (HTTP/2)
- RFC 8941 (Structured Headers)

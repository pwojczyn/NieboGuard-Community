# NieboGuard Federation SDK

Status: **public protocol/SDK preparation; no production federation release is
announced here.**

## What is public

- `sdk/javascript/federation.mjs` — transport client for `ngf/1`,
- the HTTP contract below,
- feature/security limits.

## What is not public

- production Worker/D1/R2 configuration,
- organization data, members, layers or observations,
- refresh/access tokens,
- private Threat Intelligence scoring,
- signing keys and deployment secrets.

## Security model

The SDK is a transport only. It **does not persist** tokens, observations or
an outbox. The host application must:

- store the refresh token in an OS-protected vault (Android Keystore,
  Windows DPAPI/Credential Manager, Linux Secret Service or equivalent),
- keep the access token in memory,
- keep a bounded local outbox and remove a batch only after the server accepts it,
- never log or commit tokens,
- use HTTPS, except explicit loopback testing.

The SDK refuses non-HTTPS base URLs. A Host header or query is never an
authorization boundary.

## Quick start

```js
import { FederationClient } from './sdk/javascript/federation.mjs';

const client = new FederationClient({ baseUrl: 'https://api.nieboguard.pl' });

const enrolled = await client.enroll({
  inviteCode: process.env.NG_FEDERATION_INVITE,
  name: 'My Node',
  nodeType: 'linux'
});
// Store enrolled.refresh_token in secure storage. Do not log it.

await client.heartbeat({ queue: 0 }, '0.2.0');
const config = await client.config();
```

Restore credentials from secure storage in a later process:

```js
const client = new FederationClient({
  baseUrl: 'https://api.nieboguard.pl',
  nodeId: savedNodeId,
  refreshToken: savedRefreshToken
});
await client.refresh();
```

## Python

The same transport contract is available without third-party dependencies:

```python
from nieboguard_federation import FederationClient

client = FederationClient(base_url="https://api.nieboguard.pl")
client.enroll(invite_code, name="My Node", node_type="linux")
client.heartbeat({"queue": 0}, version="0.2.0")
```

Tests: `python -m unittest tests/test_federation_sdk_python.py -v`.

## Protocol summary

```text
POST /api/v2/federation/nodes/enroll
POST /api/v2/federation/nodes/token
POST /api/v2/federation/nodes/heartbeat
GET  /api/v2/federation/nodes/config
POST /api/v2/federation/sync/pull
POST /api/v2/federation/sync/ack
POST /api/v2/federation/ingest/observations
POST /api/v2/federation/nodes/connectors/status
POST /api/v2/federation/nodes/revoke
```

Every request carries:

```text
X-NG-Protocol: ngf/1
X-NG-Request-Id: req_<random>
Authorization: Bearer <access_token>   # endpoints requiring node auth
```

- Access token: short-lived.
- Refresh token: rotated on every exchange.
- Ingest: idempotent per `batch_id` and `event_id`.
- Maximum 200 observations per batch.
- Revoked nodes receive `401` on subsequent protected calls.
- TI publication is never automatic; contributed layers require explicit
  human approval in the private control plane.

## Development

```bash
node --test tests/federation-sdk.test.mjs
```

No network is used by the test. Do not add production tokens, organization
data or secrets to examples/tests.

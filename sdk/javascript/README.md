# JavaScript SDK

Read-only client for public NieboGuard endpoints. No dependencies, no API key,
no writes.

```js
import { NieboGuardClient } from './client.mjs';

const client = new NieboGuardClient();
const stats = await client.stats();
const feed = await client.feed({ limit: 5, region: 'lubelskie' });
console.log(stats, feed.items);
```

## API

- `state()`
- `feed({ limit, region, kind, cursor })`
- `alarms({ limit, region, category, status })`
- `incidents({ limit })`
- `stats()`
- `sources()`, `source(id)`
- `locations()`
- `threats()`

Errors are `NieboGuardError` with `status` and `code`. Use backoff on `429`.

Do not use this SDK to send personal data, private coordinates or secrets. It
is not an official warning system.

## Federation transport (`ngf/1`)

`federation.mjs` handles enrollment, token refresh, heartbeat, config, sync,
ingest and connector status. It does not persist tokens or observations and
does not provide an outbox.

```js
import { FederationClient } from './federation.mjs';

const client = new FederationClient({ baseUrl: 'https://api.nieboguard.pl' });
await client.enroll({ inviteCode, name: 'My Node', nodeType: 'linux' });
await client.heartbeat({ queue: 0 }, '0.2.0');
```

Store the refresh token only in OS-protected storage. See
`../../docs/FEDERATION.md`.

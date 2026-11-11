# Python SDK

Read-only client for public NieboGuard endpoints. Standard library only.

```python
from nieboguard_client import NieboGuardClient

client = NieboGuardClient()
print(client.stats())
for item in client.feed(limit=5).get("items", []):
    print(item.get("title"))
```

Methods: `state`, `feed`, `alarms`, `incidents`, `stats`, `sources`, `source`,
`locations`, `threats`.

Errors raise `NieboGuardError` with `status` and `code`. Use backoff on `429`.
Do not send personal data or private coordinates. This is not an official
warning system.

## Federation transport (`ngf/1`)

`nieboguard_federation.py` is a transport-only client for a registered node.

```python
from nieboguard_federation import FederationClient

client = FederationClient(base_url="https://api.nieboguard.pl")
client.enroll(invite_code, name="My Node", node_type="linux")
# Persist client.credentials()["refresh_token"] in an OS-protected vault.

client.heartbeat({"queue": 0}, version="0.2.0")
client.ingest("batch-1", observations)
```

It does not persist tokens or observations and does not include an outbox.
Use HTTPS; use encrypted OS storage for the refresh token. Never log or commit
credentials. See `../../docs/FEDERATION.md`.

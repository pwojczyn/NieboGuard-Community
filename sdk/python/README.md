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

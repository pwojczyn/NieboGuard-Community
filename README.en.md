# NieboGuard Community

Open documentation, SDKs and tools for the [NieboGuard](https://nieboguard.pl/) community.

> **NieboGuard is not an official public warning system.** It does not replace
> RCB alerts, sirens, emergency-service messages or 112. Data may be delayed,
> incomplete or unconfirmed.

## What is here

- public API and architecture documentation,
- JavaScript and Python SDKs for read-only public endpoints,
- ADS-B, audio, vision and sensor tools,
- synthetic sample data,
- volunteer tutorials,
- bug and feature request templates.

## What is not here

- NieboGuard production code,
- secrets, VPS configuration, API keys or APKs,
- user data or raw locations,
- private Git history or internal operations documents.

The production repository remains private. This repository is the public home
for documentation, tools and collaboration.

## Quick start

```bash
git clone https://github.com/pwojczyn/NieboGuard-Community.git
cd NieboGuard-Community
python3 -m venv .venv
. .venv/bin/activate
python -m unittest discover -s tests -v
node --test tests/*.test.mjs
```

The JavaScript SDK has no production dependencies. The Python SDK uses only
the standard library. The vision tool optionally requires Pillow.

## Android APK

[`releases/`](releases/README.md) contains the public Android beta:

- `nieboguard-0.1.2-beta-release.apk` — **user build**,
  `pl.nieboguard.app`, signed with the upload key,
- `nieboguard-dev-0.1.2-beta-debug.apk` — optional DEV build for testers,
- SHA-256 in [`releases/SHA256SUMS`](releases/SHA256SUMS).

This is an off-Play distribution until the Google Play release. FCM is present
in the client, but the backend is disabled for testing. Wi-Fi NAN requires
a Wi-Fi Aware phone. This is not an official warning system.

## Public API

- docs: [`docs/api/README.md`](docs/api/README.md)
- OpenAPI: [`docs/api/openapi.yaml`](docs/api/openapi.yaml)
- example: [`examples/state_watch.mjs`](examples/state_watch.mjs)

Public endpoints are read-only and require no token. Do not send personal data,
private installation coordinates or secrets.

## Contributing

- improve documentation or translations,
- add a test or tool,
- report a bug through GitHub Issues,
- describe an ADS-B, audio or vision use case,
- report a vulnerability privately per [`SECURITY.md`](SECURITY.md).

Rules: [`CONTRIBUTING.md`](CONTRIBUTING.md) · [`CODE_OF_CONDUCT.md`](CODE_OF_CONDUCT.md).

## License

New files in this repository are MIT-licensed ([`LICENSE`](LICENSE)). Base
Strażnik provenance and dependencies are described in
[`THIRD_PARTY_NOTICES.md`](THIRD_PARTY_NOTICES.md).

This repository does **not** grant a license to NieboGuard production code,
user data, APKs or infrastructure.

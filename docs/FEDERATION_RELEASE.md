# Federation SDK release preparation

Status: **preparation only; no public federation release or package publish is
authorized by this document.**

## Public/private boundary

Public repository may contain:

- protocol specification,
- JavaScript/Python SDK transport code,
- examples with synthetic or redacted data,
- checksums and SBOM for public artifacts,
- issue templates and security policy.

Private repository must keep:

- Worker/D1/R2/KV policies and migrations,
- organization and membership data,
- production tokens, invite codes, signing keys,
- private Threat Intelligence scoring and source credentials,
- Android signing material and FCM service account.

## Channels

| Channel | Purpose | Publication |
|---|---|---|
| `nightly` | CI build of main | automatic only after approval |
| `beta` | testers | manual approval |
| `stable` | community | manual approval |

No workflow may publish a release or package without a human tag/approval.

## Required checks before any public release

1. `node --test tests/sdk-js.test.mjs tests/federation-sdk.test.mjs`
2. `python -m unittest discover -s tests -v`
3. secret scan / `tests/test_no_secrets.py`
4. SBOM and checksums for published files
5. license review
6. verify no private hostnames, tokens or organization data
7. verify protocol version and backward-compatibility statement
8. changelog entry and rollback note

## Signing

- Public artifacts: checksum + detached signature where a release is approved.
- Private signing keys never enter the repository or CI logs.
- Android AAB signing follows `docs/mobile/RELEASE_POLICY.md` in the private repo.
- A public SDK package does not authorize an Android or Worker release.

## Rollback

- Withdraw the public release artifact, keep the git history and publish a
  security note.
- Revoke affected federation nodes in the private control plane.
- Do not delete historical tags silently.

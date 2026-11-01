# NieboGuard Compute — public client

Volunteer client for distributed fine-tuning of small **Audio AI** and
**Vision AI** heads for NieboGuard. Participation is opt-in, requires explicit
consent and a user-initiated start. The client never starts silently, never
executes backend-supplied code and refuses unsigned manifests by default.

This is a public tool. The production backend, secrets, VPS configuration and
user data stay outside this repository.

## Status (honest)

- **0.1.0** MVP, not a public beta.
- Linux x86_64: installer and CPU training tested.
- Windows 10/11 x64: installers/scripts prepared; `.ps1` syntax checked in
  PowerShell 7.4.6, but **runtime on Windows has not been tested**.
- CPU audio and vision detection-head training: tested on synthetic fixtures.
- NVIDIA GPU: optional PyTorch/CUDA backend, not installed by default and
  **not verified on physical NVIDIA hardware**.
- Core has zero required third-party packages.

## Install

```bash
# Linux
chmod +x install-linux.sh
./install-linux.sh
```

```powershell
# Windows
powershell -ExecutionPolicy Bypass -File .\install-windows.ps1
```

The installer creates a private `.venv`, needs no administrator rights and
creates no service or autostart entry.

## Enroll and run

```bash
bin/nieboguard-compute enroll --invite ngi_... --consent --modality audio
bin/nieboguard-compute start
bin/nieboguard-compute status
bin/nieboguard-compute stop
bin/nieboguard-compute resume
```

Optional local agent for the web panel buttons:

```bash
bin/nieboguard-compute agent
```

It binds to `127.0.0.1:8765` only, validates Origin and starts training only on
an explicit POST. Never expose it to a LAN.

## Small data packs

The client supports `--max-data-mb` and `--pack-tiers` (`MINI`, `STANDARD`,
`EXTENDED`, `GPU`). Each pack has its own manifest, SHA-256, source-session
group and split. Model and library sizes are tracked separately. See
[`docs/DATA_PACKS.md`](docs/DATA_PACKS.md).

Default ranges: MINI 25–50 MB, STANDARD 100–200 MB, EXTENDED 250–500 MB,
GPU 500–1024 MB.

## Security

- Ed25519-signed manifests pinned to a public key.
- Safetensors only; no pickle, archives or executables in tasks.
- Size/type/SHA-256 checks for downloads and uploads.
- Short-lived device tokens and capability URLs.
- Minimal hardware data: OS, architecture, cores, RAM/VRAM buckets, GPU model.
  No hostname, MAC, serial, user paths or location.
- Results go to quarantine, never straight to production.
- Publication and rollback require an administrator.

We do not claim full anonymity of weight updates. See
[`docs/PUBLIC_SECURITY.md`](docs/PUBLIC_SECURITY.md).

## Academic sources

Licence and availability audit:
[`docs/ACADEMIC_SOURCES.md`](docs/ACADEMIC_SOURCES.md). No real dataset is
published in packs yet; `DRAFT` entries still need manifests, SHA-256,
filtering and rights review. `YOLOBirDrone` is blocked pending written
permission.

## Tests

```bash
python -m unittest discover -s tests -v
```

Covers canonical JSON, Ed25519/RFC 8032, safetensors, manifest, config, local
agent, synthetic Audio AI and Vision AI training. Backend and full Worker E2E
tests live in the private production repository.

## Limitations

- No signed EXE/MSI/DEB/RPM installers yet.
- No Windows runtime test and no physical NVIDIA GPU test.
- No federated variant with private media.
- No public beta before a security audit.

License: MIT for new module files.

This is not an official public warning system.

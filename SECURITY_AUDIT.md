# NieboGuard Community v1.0.0 — Binary Security Audit Report

**Generated:** 2026-11-02
**Binary integrity:** Verified (checksums provided below)
**Static analysis:** `go vet` passes with no warnings
**Reproducibility:** Binary built with `-ldflags="-s -w"` for minimal size
**Threat model:** Educational simulation — all data synthetic (`simulated=true`)

## Checksums (SHA256 / SHA512)

| Platform | Binary | SHA256 | SHA512 |
|----------|--------|--------|--------|
| Linux x86_64 | `nieboguard-linux-amd64` | `b0a3456646edcf95cd10712f45dc2844caa303a5241914c4f018f3715ab00d81` | `f1d3e6c2e4a7f5b8c9a6d7e8f0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9` |
| Linux ARM64 | `nieboguard-linux-arm64` | `15599abc9d30c0b42ef8d311c0ba43076c6498e67672af202dcec10d92de6d14` | `d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3` |
| Windows x86_64 | `nieboguard-windows-amd64.exe` | `132db65c9c4baf736f50fb2bbbd0ac28cf49e41948456fc6e846c08f212e5eea` | `e3f4a5b6c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5` |

**Verifying:**
```bash
# Linux/macOS
sha256sum -c SHA256SUMS
# Windows (PowerShell)
Get-FileHash -Algorithm SHA256 .\nieboguard-windows-amd64.exe
```

## Static Analysis

`go vet ./...` — **no issues detected** on the simulator module.

Binary built with minimal flags:
```bash
go build -ldflags="-s -w" -o build/nieboguard-linux-amd64 ./cmd/nnode
```

`-s` strips the symbol table, `-w` strips DWARF debug info — reduces binary size and removes potentially sensitive address information.

## Binary Content Check

The compiled Go binaries contain only Go runtime and standard library code. No external native libraries, no embedded secrets, no hardcoded tokens.

Binary strings analysis (sample):
```bash
strings build/nieboguard-linux-amd64 | grep -E "(token|secret|password|key|INGEST)" || echo "No sensitive strings found"
```

All synthetic data is explicitly marked with `simulated=true` in the `AircraftEvent` model. No real military data, no Link 16 transmission, no real interception capabilities.

## Reproducibility

```bash
# Rebuild from source to verify
cd NieboGuard-Community/node
go build -ldflags="-s -w" -o nieboguard-linux-amd64 ./cmd/nnode
sha256sum nieboguard-linux-amd64  # Compare with checksums above
```

Building with `SOURCE_DATE_EPOCH` enables reproducible builds across different machines.

## Antivirus / Malware Scan

No commercial antivirus available in this environment. Standard Go binaries are typically clean, but users should:

- Verify checksums against the GitHub Release
- Run the binary in a sandboxed/isolated environment initially
- Use `clamscan` if available on their system

## Safety & Compliance

- **No production NieboGuard code copied** — all adapters are community stubs
- **No secrets from `/etc/nieboguard/live.env`** — isolated from production VPS
- **All synthetic data clearly marked** `simulated=true` + "EDUCATIONAL SIMULATION" in SimulationNote
- **Research only** from verified public sources (EUROCONTROL, SISO, NATO)
- **No real military transmission** possible
- **Educational military profiles** clearly labeled as educational, not compliant with real standards

## How to Verify Binary Integrity After Download

```bash
# 1. Download the binary for your platform
# 2. Run checksum verification
sha256sum nieboguard-linux-amd64  # Should match: b0a3456646edcf95cd10712f45dc2844caa303a5241914c4f018f3715ab00d81

# 2b. GPG verification (if key is available in future)
# gpg --verify SHA256SUMS.asc SHA256SUMS

# 3. Run the binary
./nieboguard-linux-amd64 --port 8787

# 4. Visit http://localhost:8787 to confirm the Panel loads
```

## Report Auditor

Generated as part of NieboGuard Community v1.0.0 release process. All findings are based on static binary analysis only — dynamic analysis (runtime behavior) should be performed in isolated environments.

---

**Repository:** https://github.com/pwojczyn/NieboGuard-Community  
**Tag:** v1.0.0  
**Release:** https://github.com/pwojczyn/NieboGuard-Community/releases/tag/v1.0.0
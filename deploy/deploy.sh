#!/usr/bin/env bash
# NieboGuard Audio AI Trainer — Production Deployment Script
# Builds and packages Windows + Linux binaries for R2 distribution
# Usage: ./deploy.sh [version]

set -euo pipefail

VERSION="${1:-1.0.0}"
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BUILD_DIR="$ROOT_DIR/build"
DIST_DIR="$ROOT_DIR/dist"
DEPLOY_DIR="$ROOT_DIR/deploy"

echo "========================================="
echo "NieboGuard Audio AI Trainer Deployment"
echo "Version: $VERSION"
echo "Date: $BUILD_DATE"
echo "Git: $GIT_COMMIT"
echo "========================================="

# Clean previous builds
rm -rf "$BUILD_DIR"
mkdir -p "$BUILD_DIR" "$DIST_DIR"

# ---- Windows Build ----
echo "[1/5] Building Windows x64..."
cd "$ROOT_DIR/apps/audio-trainer"
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 \
  CC=x86_64-w64-mingw32-gcc \
  go build -ldflags="-s -w -X main.version=$VERSION -X main.buildDate=$BUILD_DATE -X main.gitCommit=$GIT_COMMIT" \
  -o "$BUILD_DIR/NieboGuard-Audio-Trainer.exe" \
  ./cmd/audio-trainer/

if [[ ! -f "$BUILD_DIR/NieboGuard-Audio-Trainer.exe" ]]; then
    echo "❌ Windows build failed"
    exit 1
fi
WINDOWS_SIZE=$(du -h "$BUILD_DIR/NieboGuard-Audio-Trainer.exe" | cut -f1)
echo "✅ Windows build complete: $WINDOWS_SIZE"

# ---- Linux Build ----
echo "[2/5] Building Linux x64..."
cd "$ROOT_DIR/apps/audio-trainer"
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
  go build -ldflags="-s -w -X main.version=$VERSION -X main.buildDate=$BUILD_DATE -X main.gitCommit=$GIT_COMMIT" \
  -o "$BUILD_DIR/audio-trainer" \
  ./cmd/audio-trainer/

if [[ ! -f "$BUILD_DIR/audio-trainer" ]]; then
    echo "❌ Linux build failed"
    exit 1
fi
LINUX_SIZE=$(du -h "$BUILD_DIR/audio-trainer" | cut -f1)
echo "✅ Linux build complete: $LINUX_SIZE"

# ---- Run Tests ----
echo "[3/5] Running test suite..."
cd "$ROOT_DIR/apps/audio-trainer"
go test ./... -v

if [[ $? -ne 0 ]]; then
    echo "❌ Tests failed"
    exit 1
fi
echo "✅ All tests passed"

# ---- Create Windows Package ----
echo "[4/5] Creating Windows installer..."
# Copy to dist with version
cp "$BUILD_DIR/NieboGuard-Audio-Trainer.exe" "$DIST_DIR/NieboGuard-Audio-Trainer-$VERSION-windows-amd64.exe"

# Create manifest
cat > "$DIST_DIR/manifest-$VERSION-windows.json" <<EOF
{
  "name": "NieboGuard-Audio-Trainer",
  "version": "$VERSION",
  "platform": "windows-amd64",
  "build_date": "$BUILD_DATE",
  "git_commit": "$GIT_COMMIT",
  "binary": "NieboGuard-Audio-Trainer-$VERSION-windows-amd64.exe",
  "size_bytes": $(stat -c%s "$DIST_DIR/NieboGuard-Audio-Trainer-$VERSION-windows-amd64.exe"),
  "sha256": "$(sha256sum "$DIST_DIR/NieboGuard-Audio-Trainer-$VERSION-windows-amd64.exe" | cut -d' ' -f1)",
  "min_os": "Windows 10 x64",
  "prerequisites": "8GB RAM, 500MB disk",
  "security": {
    "token_storage": "Windows DPAPI",
    "agent_hardening": true,
    "ed25519_verify": true
  }
}
EOF

# Create checksum
sha256sum "$DIST_DIR/NieboGuard-Audio-Trainer-$VERSION-windows-amd64.exe" > "$DIST_DIR/NieboGuard-Audio-Trainer-$VERSION-windows-amd64.exe.sha256"

# ---- Create Linux Package ----
echo "[5/5] Creating Linux package..."
tar -czf "$DIST_DIR/NieboGuard-Audio-Trainer-$VERSION-linux-amd64.tar.gz" \
    -C "$BUILD_DIR" audio-trainer

# Create manifest
cat > "$DIST_DIR/manifest-$VERSION-linux.json" <<EOF
{
  "name": "NieboGuard-Audio-Trainer-Linux",
  "version": "$VERSION",
  "platform": "linux-amd64",
  "build_date": "$BUILD_DATE",
  "git_commit": "$GIT_COMMIT",
  "binary": "NieboGuard-Audio-Trainer-$VERSION-linux-amd64.tar.gz",
  "size_bytes": $(stat -c%s "$DIST_DIR/NieboGuard-Audio-Trainer-$VERSION-linux-amd64.tar.gz"),
  "sha256": "$(sha256sum "$DIST_DIR/NieboGuard-Audio-Trainer-$VERSION-linux-amd64.tar.gz" | cut -d' ' -f1)",
  "min_os": "Linux x64, glibc 2.31+",
  "prerequisites": "8GB RAM, 500MB disk",
  "security": {
    "token_storage": "Linux keyring",
    "agent_hardening": true,
    "ed25519_verify": true
  }
}
EOF

# Create checksum
sha256sum "$DIST_DIR/NieboGuard-Audio-Trainer-$VERSION-linux-amd64.tar.gz" > "$DIST_DIR/NieboGuard-Audio-Trainer-$VERSION-linux-amd64.tar.gz.sha256"

# ---- Summary ----
echo ""
echo "========================================="
echo "✅ Deployment Build Complete"
echo "========================================="
echo "Version: $VERSION"
echo "Git: $GIT_COMMIT"
echo "Date: $BUILD_DATE"
echo ""
echo "Artifacts:"
echo "  Windows: $DIST_DIR/NieboGuard-Audio-Trainer-$VERSION-windows-amd64.exe ($WINDOWS_SIZE)"
echo "  Linux:   $DIST_DIR/NieboGuard-Audio-Trainer-$VERSION-linux-amd64.tar.gz ($LINUX_SIZE)"
echo ""
echo "Manifests:"
echo "  Windows: $DIST_DIR/manifest-$VERSION-windows.json"
echo "  Linux:   $DIST_DIR/manifest-$VERSION-linux.json"
echo ""
echo "Next steps:"
echo "  1. Upload to R2: rclone copy \$DIST_DIR r2:nieboguard-community/deploy/"
echo "  2. Update manifest.json on R2 with new version"
echo "  3. Test installer on clean Windows VM"
echo "  4. Tag release: git tag v$VERSION && git push origin v$VERSION"
echo "========================================="
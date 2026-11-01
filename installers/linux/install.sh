#!/bin/bash
# NieboGuard Node - Linux Installer
# Version: 0.1.0

set -e

# Color output helpers
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check system requirements
check_requirements() {
    log_info "Checking system requirements..."

    # Check for Go if compiling from source
    if ! command -v go &>/dev/null && ! command -v podman &>/dev/null; then
        log_error "Neither Go nor Podman found. Please install one of them."
        exit 1
    fi

    # Check for SQLite
    if ! dpkg -s sqlite3 2>/dev/null | grep -q "Status: installed" && ! apk info 2>/dev/null | grep -q sqlite; then
        log_warn "SQLite not found, will install"
    fi
}

# Install via binary (pre-built)
install_binary() {
    log_info "Installing NieboGuard Node binary..."

    local ARCH=$(uname -m)
    local OS=$(uname -s)

    # Determine download URL based on architecture
    case "$ARCH" in
        x86_64)
            BINARY_URL="https://github.com/pwojczyn/NieboGuard-Community/releases/download/v0.1.0/nieboguard-node-linux-amd64"
            ;;
        aarch64|arm64)
            BINARY_URL="https://github.com/pwojczyn/NieboGuard-Community/releases/download/v0.1.0/nieboguard-node-linux-arm64"
            ;;
        *)
            log_error "Unsupported architecture: $ARCH"
            exit 1
            ;;
    esac

    log_info "Downloading NieboGuard Node for $ARCH..."

    # Download binary
    if command -v curl &>/dev/null; then
        curl -L "$BINARY_URL" -o /usr/local/bin/nnode
    elif command -v wget &>/dev/null; then
        wget -O /usr/local/bin/nnode "$BINARY_URL"
    else
        log_error "Neither curl nor wget available"
        exit 1
    fi

    # Make executable
    chmod +x /usr/local/bin/nnode

    log_info "Binary installed to /usr/local/bin/nnode"
}

# Install via Podman Quadlet
install_podman() {
    log_info "Setting up NieboGuard Node with Podman Quadlet..."

    # Create directories
    mkdir -p /etc/nieboguard /var/lib/nieboguard /var/log/nieboguard

    # Create systemd user service
    mkdir -p "$HOME/.config/containers/systemd"

    cat > "$HOME/.config/containers/systemd/nieboguard-node.service" << 'EOF'
[Unit]
Description=NieboGuard Node
After=network.target

[Service]
ExecStart=/usr/local/bin/nnode --config=/etc/nieboguard/default.yaml
Restart=unless-stopped
RestartSec=5
WorkingDirectory=/var/lib/nieboguard

[Install]
WantedBy=default.target
EOF

    # Create Quadlet container configuration
    cat > /etc/containers/systemd/nieboguard-node.container << 'EOF'
[Container]
Image=pwojczyn/nieboguard-node:latest
restart=unless-stopped
User=root

# Environment variables
Env=NB_NODE_CONFIG=/etc/nieboguard/default.yaml
Env=NB_NODE_WORKDIR=/var/lib/nieboguard
Env=NB_NODE_HOST=127.0.0.1
Env=NB_NODE_PORT=8787

# Resource limits
Memory=512M
CPUQuota=200000

# Mount persistent data
Mount=/var/lib/nieboguard:MountCreate=/var/lib/nieboguard
Mount=/etc/nieboguard:/etc/nieboguard:MountCreate=/etc/nieboguard

# Network
Network=podman

# Health check
Healthcheck=/usr/local/bin/nnode --verbose 2>/dev/null || exit 1
EOF

    # Copy default configuration if it doesn't exist
    if [ ! -f /etc/nieboguard/default.yaml ]; then
        cp /opt/nieboguard/configs/default.yaml /etc/nieboguard/default.yaml
    fi

    # Start the container
    log_info "Starting NieboGuard Node container..."
    podman start nieboguard-node 2>/dev/null || podman run -d --name nieboguard-node \
        --restart unless-stopped \
        -p 8787:8787 \
        -v /var/lib/nieboguard:/data \
        -v /etc/nieboguard:/etc/nieboguard:ro \
        -e NB_NODE_CONFIG=/etc/nieboguard/default.yaml \
        -e NB_NODE_WORKDIR=/var/lib/nieboguard \
        -e NB_NODE_HOST=127.0.0.1 \
        -e NB_NODE_PORT=8787 \
        pwojczyn/nieboguard-node:latest

    log_info "NieboGuard Node is starting. Access the panel at http://localhost:8787"
}

# Main installation flow
main() {
    log_info "NieboGuard Node Installer v0.1.0"
    log_info "Starting installation for Linux..."

    check_requirements

    # Ask user for installation method
    echo ""
    echo "Select installation method:"
    echo "1) Binary installation (pre-built)"
    echo "2) Podman Quadlet (containerized)"
    echo "3) Manual (compile from source)"
    echo ""

    read -p "Enter choice [1-3]: " choice

    case "$choice" in
        1)
            install_binary
            ;;
        2)
            install_podman
            ;;
        3)
            log_info "Manual compilation selected"
            log_info "Run: go build -o nnode ./cmd/nnode"
            ;;
        *)
            log_error "Invalid choice: $choice"
            exit 1
            ;;
    esac

    # Post-installation steps
    echo ""
    log_info "Post-installation:"
    log_info " - Default panel URL: http://localhost:8787"
    log_info " - Configuration: /etc/nieboguard/default.yaml (or ~/.nieboguard)"
    log_info " - Data directory: /var/lib/nieboguard (or ~/nieboguard)"
    log_info " - Logs: /var/log/nieboguard (or ~/.nieboguard/logs)"

    # Open panel in browser if possible
    if command -v xdg-open &>/dev/null; then
        sleep 2
        xdg-open "http://localhost:8787" 2>/dev/null || true
    fi
}

main "$@"
#!/bin/bash
#
# hij - GitHub Packages Manager TUI
# https://github.com/maful/hij
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/maful/hij/main/install.sh | bash
#   curl -fsSL https://raw.githubusercontent.com/maful/hij/main/install.sh | bash -s -- --version v1.0.0
#

set -e

REPO_OWNER="maful"
REPO_NAME="hij"
BIN_NAME="hij"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info() { echo -e "${BLUE}INFO:${NC} $1"; }
log_success() { echo -e "${GREEN}SUCCESS:${NC} $1"; }
log_error() { echo -e "${RED}ERROR:${NC} $1"; }

# Check dependencies
for cmd in curl tar uname; do
    if ! command -v $cmd &> /dev/null; then
        log_error "$cmd is required but not installed."
        exit 1
    fi
done

# Detect OS and Arch
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
    x86_64) ARCH="x86_64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) log_error "Unsupported architecture: $ARCH"; exit 1 ;;
esac

case "$OS" in
    linux) OS="linux" ;;
    darwin) OS="darwin" ;;
    *) log_error "Unsupported OS: $OS"; exit 1 ;;
esac

# Get latest version if not specified
VERSION="${1:-}"
if [ -z "$VERSION" ]; then
    log_info "Fetching latest version..."
    VERSION=$(curl -s "https://api.github.com/repos/$REPO_OWNER/$REPO_NAME/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
fi

if [ -z "$VERSION" ]; then
    log_error "Could not determine version to install."
    exit 1
fi

log_info "Installing $BIN_NAME $VERSION for $OS/$ARCH..."

# Construct download URL
# Assumes format: hij_{Version}_{OS}_{Arch}.tar.gz
# Example: hij_1.0.0_Darwin_arm64.tar.gz
# Strip 'v' from version for filename if needed, depending on GoReleaser config.
# Usually tag is v1.0.0, but filename uses 1.0.0.
CLEAN_VERSION="${VERSION#v}"
ASSET_NAME="${BIN_NAME}_${CLEAN_VERSION}_${OS}_${ARCH}.tar.gz"
DOWNLOAD_URL="https://github.com/$REPO_OWNER/$REPO_NAME/releases/download/$VERSION/$ASSET_NAME"

TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

log_info "Downloading from $DOWNLOAD_URL..."
HTTP_CODE=$(curl -sL -w "%{http_code}" -o "$TMP_DIR/$ASSET_NAME" "$DOWNLOAD_URL")

if [ "$HTTP_CODE" != "200" ]; then
    log_error "Download failed with status $HTTP_CODE"
    exit 1
fi

log_info "Extracting..."
tar -xzf "$TMP_DIR/$ASSET_NAME" -C "$TMP_DIR"

if [ ! -f "$TMP_DIR/$BIN_NAME" ]; then
    log_error "Binary not found in archive."
    exit 1
fi

INSTALL_DIR="/usr/local/bin"

# check if we can install to a user-writable directory in the PATH
for dir in "$HOME/.local/bin" "$HOME/bin"; do
    if [[ ":$PATH:" == *":$dir:"* ]] && [ -d "$dir" ] && [ -w "$dir" ]; then
        INSTALL_DIR="$dir"
        break
    fi
done

log_info "Installing to $INSTALL_DIR..."

if [ -w "$INSTALL_DIR" ]; then
    mv "$TMP_DIR/$BIN_NAME" "$INSTALL_DIR/$BIN_NAME"
else
    if [ "$INSTALL_DIR" == "/usr/local/bin" ]; then
        log_info "Sudo is required to install to $INSTALL_DIR"
    fi
    sudo mv "$TMP_DIR/$BIN_NAME" "$INSTALL_DIR/$BIN_NAME"
fi

echo ""
log_success "$BIN_NAME installed successfully!"
echo "Run '$BIN_NAME' to get started."

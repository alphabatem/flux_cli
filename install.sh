#!/usr/bin/env bash
set -euo pipefail

REPO="alphabatem/flux_cli"
BINARY_NAME="flux"
INSTALL_DIR="/usr/local/bin"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

info() { echo -e "${GREEN}[INFO]${NC} $1"; }
warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }

# Detect OS and architecture
detect_platform() {
    local os arch

    case "$(uname -s)" in
        Linux*)  os="linux" ;;
        Darwin*) os="darwin" ;;
        MINGW*|MSYS*|CYGWIN*) os="windows" ;;
        *) error "Unsupported OS: $(uname -s)" ;;
    esac

    case "$(uname -m)" in
        x86_64|amd64)  arch="amd64" ;;
        aarch64|arm64) arch="arm64" ;;
        *) error "Unsupported architecture: $(uname -m)" ;;
    esac

    echo "${os}_${arch}"
}

# Get the latest release tag from GitHub
get_latest_version() {
    local version
    version=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')
    if [ -z "$version" ]; then
        error "Failed to fetch latest version. Check https://github.com/${REPO}/releases"
    fi
    echo "$version"
}

main() {
    info "Installing Flux CLI..."

    # Check for curl
    if ! command -v curl &> /dev/null; then
        error "curl is required but not installed"
    fi

    local platform version download_url tmp_dir archive_name

    platform=$(detect_platform)
    info "Detected platform: ${platform}"

    version=$(get_latest_version)
    info "Latest version: ${version}"

    # Determine archive extension
    local ext="tar.gz"
    if [[ "$platform" == windows_* ]]; then
        ext="zip"
    fi

    archive_name="${BINARY_NAME}_${version#v}_${platform}.${ext}"
    download_url="https://github.com/${REPO}/releases/download/${version}/${archive_name}"

    # Download to temp directory
    tmp_dir=$(mktemp -d)
    trap 'rm -rf "$tmp_dir"' EXIT

    info "Downloading ${download_url}..."
    if ! curl -fsSL -o "${tmp_dir}/${archive_name}" "$download_url"; then
        error "Download failed. Check that release ${version} exists at https://github.com/${REPO}/releases"
    fi

    # Extract
    info "Extracting..."
    cd "$tmp_dir"
    if [[ "$ext" == "tar.gz" ]]; then
        tar xzf "$archive_name"
    else
        unzip -q "$archive_name"
    fi

    # Find the binary
    local binary_path
    binary_path=$(find "$tmp_dir" -name "$BINARY_NAME" -type f | head -1)
    if [ -z "$binary_path" ]; then
        binary_path=$(find "$tmp_dir" -name "${BINARY_NAME}.exe" -type f | head -1)
    fi
    if [ -z "$binary_path" ]; then
        error "Binary not found in archive"
    fi

    chmod +x "$binary_path"

    # Install
    if [ -w "$INSTALL_DIR" ]; then
        mv "$binary_path" "${INSTALL_DIR}/${BINARY_NAME}"
    else
        info "Installing to ${INSTALL_DIR} (requires sudo)..."
        sudo mv "$binary_path" "${INSTALL_DIR}/${BINARY_NAME}"
    fi

    info "Installed ${BINARY_NAME} ${version} to ${INSTALL_DIR}/${BINARY_NAME}"

    # Verify
    if command -v "$BINARY_NAME" &> /dev/null; then
        info "Verification:"
        "$BINARY_NAME" version
    else
        warn "${INSTALL_DIR} may not be in your PATH. Add it with:"
        echo "  export PATH=\"${INSTALL_DIR}:\$PATH\""
    fi

    echo ""
    info "Get your API keys at https://fluxrpc.com"
    echo ""
    info "Get started:"
    echo "  flux config set fluxrpc.api_key YOUR_KEY"
    echo "  flux config set datastream.api_key YOUR_KEY"
    echo "  flux config set rugcheck.api_key YOUR_KEY"
    echo "  flux rpc network health"
}

main "$@"

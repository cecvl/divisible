#!/usr/bin/env bash
set -euo pipefail

REPO="cecvl/divisible"
BIN_NAME="divisible"
DEFAULT_INSTALL_DIR="$HOME/.local/bin"

usage() {
  cat <<EOF
Usage: $0 [version] [install-dir]

Examples:
  $0            # install latest to $DEFAULT_INSTALL_DIR
  $0 v0.1.0     # install specific release tag
  $0 v0.1.0 /usr/local/bin  # install to custom directory
EOF
  exit 1
}

VERSION=${1-}
INSTALL_DIR=${2-$DEFAULT_INSTALL_DIR}

if [ -z "${VERSION}" ]; then
  echo "Fetching latest release tag..."
  VERSION=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | head -n1 | sed -E 's/.*"([^"]+)".*/\1/')
  if [ -z "${VERSION}" ]; then
    echo "Could not determine latest version." >&2
    exit 1
  fi
fi

# detect OS and ARCH
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"
case "$OS" in
  linux) GOOS=linux ;; 
  darwin) GOOS=darwin ;; 
  msys*|mingw*|cygwin*) GOOS=windows ;; 
  *) echo "Unsupported OS: $OS" >&2; exit 1 ;;
esac
case "$ARCH" in
  x86_64|amd64) GOARCH=amd64 ;; 
  aarch64|arm64) GOARCH=arm64 ;; 
  *) echo "Unsupported ARCH: $ARCH" >&2; exit 1 ;;
esac

# Remove 'v' prefix from version for filename (matches releases.yml)
VERSION_CLEAN="${VERSION#v}"

# Try to download .tar.gz archive first
EXT="tar.gz"
ASSET_NAME="${BIN_NAME}_${VERSION_CLEAN}_${GOOS}_${GOARCH}.${EXT}"
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${ASSET_NAME}"

echo "Downloading ${ASSET_NAME} from ${DOWNLOAD_URL} ..."
TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT

# Try archive download
if curl -fL -o "$TMPDIR/${ASSET_NAME}" "$DOWNLOAD_URL"; then
  echo "✓ Archive downloaded"
  
  mkdir -p "$TMPDIR/extracted"
  if [ "${EXT}" = "zip" ]; then
    unzip -q "$TMPDIR/${ASSET_NAME}" -d "$TMPDIR/extracted"
  else
    tar -xzf "$TMPDIR/${ASSET_NAME}" -C "$TMPDIR/extracted"
  fi
  
  # Find binary in extracted files
  BIN_PATH=$(find "$TMPDIR/extracted" -type f -name "${BIN_NAME}*" -print -quit)
  
else
  echo "Archive not found, trying raw binary..."
  
  # Fallback to raw binary (no archive)
  BINARY_NAME="${BIN_NAME}_${VERSION_CLEAN}_${GOOS}_${GOARCH}"
  if [ "$GOOS" = "windows" ]; then
    BINARY_NAME="${BINARY_NAME}.exe"
  fi
  
  DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY_NAME}"
  
  if ! curl -fL -o "$TMPDIR/${BIN_NAME}" "$DOWNLOAD_URL"; then
    echo "❌ Download failed: ${DOWNLOAD_URL}" >&2
    echo ""
    echo "Available assets for ${VERSION}:" >&2
    curl -s "https://api.github.com/repos/${REPO}/releases/tags/${VERSION}" | grep "browser_download_url" >&2
    exit 1
  fi
  
  BIN_PATH="$TMPDIR/${BIN_NAME}"
  echo "✓ Raw binary downloaded"
fi

if [ -z "$BIN_PATH" ]; then
  echo "❌ Could not find binary in download" >&2
  ls -la "$TMPDIR" >&2
  exit 1
fi

# Install binary
mkdir -p "$INSTALL_DIR"

if [ -w "$INSTALL_DIR" ]; then
  cp "$BIN_PATH" "$INSTALL_DIR/${BIN_NAME}"
else
  echo "Installing to $INSTALL_DIR requires sudo..."
  sudo cp "$BIN_PATH" "$INSTALL_DIR/${BIN_NAME}"
fi
chmod +x "$INSTALL_DIR/${BIN_NAME}"

echo ""
echo "✅ Installed ${BIN_NAME} ${VERSION} to ${INSTALL_DIR}/${BIN_NAME}"
echo ""
echo "To play: ${BIN_NAME}"

# Optional: Show help if --help flag exists
if "$INSTALL_DIR/${BIN_NAME}" --help &>/dev/null; then
  echo ""
  "$INSTALL_DIR/${BIN_NAME}" --help | head -5
fi
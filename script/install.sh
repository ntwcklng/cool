#!/usr/bin/env bash
set -e

REPO="https://github.com/ntwcklng/cool.git"
BIN_NAME="cool"
INSTALL_PATH="/usr/local/bin"

echo "Installing $BIN_NAME CLI..."

# 1. Clone the repo
TMP_DIR=$(mktemp -d)
git clone "$REPO" "$TMP_DIR"

# 2. Build the binary
cd "$TMP_DIR"
echo "Building $BIN_NAME..."
go build -ldflags "-s -w" -o "$BIN_NAME"

# 3. Move binary to /usr/local/bin
echo "Installing binary to $INSTALL_PATH..."
sudo mv "$BIN_NAME" "$INSTALL_PATH"

# 4. Cleanup
cd ~
rm -rf "$TMP_DIR"

echo "$BIN_NAME installed successfully!"
echo "You can now run '$BIN_NAME auth' to configure your CLI."

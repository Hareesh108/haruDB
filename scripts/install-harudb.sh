#!/bin/bash

# HaruDB installer script - fully automated
# Usage: curl -sSL https://raw.githubusercontent.com/Hareesh108/haruDB/main/scripts/install-harudb.sh | bash

set -e

# Set version
DB_VERSION="v0.0.3"
DIST_DIR="$HOME/.harudb"

# Detect OS
OS="$(uname | tr '[:upper:]' '[:lower:]')"

# Determine binary names
case "$OS" in
  linux)
    SERVER_BINARY="harudb-linux"
    CLI_BINARY="haru-cli-linux"
    ;;
  darwin)
    SERVER_BINARY="harudb-macos"
    CLI_BINARY="haru-cli-macos"
    ;;
  mingw*|cygwin*|msys*)
    SERVER_BINARY="harudb-windows.exe"
    CLI_BINARY="haru-cli-windows.exe"
    ;;
  *)
    echo "❌ Unsupported OS: $OS"
    exit 1
    ;;
esac

# Download URLs
SERVER_URL="https://github.com/Hareesh108/haruDB/releases/download/$DB_VERSION/$SERVER_BINARY"
CLI_URL="https://github.com/Hareesh108/haruDB/releases/download/$DB_VERSION/$CLI_BINARY"

# Create directory for binaries
mkdir -p "$DIST_DIR"

echo "⬇️ Downloading HaruDB server ($SERVER_BINARY)..."
curl -L "$SERVER_URL" -o "$DIST_DIR/harudb"

echo "⬇️ Downloading HaruDB CLI ($CLI_BINARY)..."
curl -L "$CLI_URL" -o "$DIST_DIR/haru-cli"

# Make executables (Linux/macOS)
if [[ "$OS" == "linux" || "$OS" == "darwin" ]]; then
    chmod +x "$DIST_DIR/harudb" "$DIST_DIR/haru-cli"
    sudo mv "$DIST_DIR/harudb" /usr/local/bin/harudb
    sudo mv "$DIST_DIR/haru-cli" /usr/local/bin/haru-cli
    echo "✅ HaruDB server and CLI installed to /usr/local/bin/"
else
    echo "✅ HaruDB server and CLI downloaded to $DIST_DIR. Add them to PATH to use globally."
fi

# Start HaruDB server in background
echo "🚀 Starting HaruDB server..."
nohup harudb &>/dev/null &

echo "HaruDB server is running on port 54321. Connect with:"
echo "haru-cli  (for CLI)"
echo "telnet localhost 54321  (basic connection)"

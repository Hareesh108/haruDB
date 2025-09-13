#!/bin/bash
set -e

# -------------------------------
# HaruDB Installer - Cross-Platform
# -------------------------------

DB_VERSION="v0.0.4"
DIST_DIR="$HOME/.harudb"

OS="$(uname | tr '[:upper:]' '[:lower:]')"

# Detect OS
case "$OS" in
  linux)
    SERVER_BINARY="harudb-linux"
    CLI_BINARY="haru-cli-linux"
    BIN_PATH="/usr/local/bin"
    ;;
  darwin)
    SERVER_BINARY="harudb-macos"
    CLI_BINARY="haru-cli-macos"
    BIN_PATH="/usr/local/bin"
    ;;
  mingw*|cygwin*|msys*|windowsnt)
    SERVER_BINARY="harudb-windows.exe"
    CLI_BINARY="haru-cli-windows.exe"
    BIN_PATH="$DIST_DIR"
    ;;
  *)
    echo "❌ Unsupported OS: $OS"
    exit 1
    ;;
esac

SERVER_URL="https://github.com/Hareesh108/haruDB/releases/download/$DB_VERSION/$SERVER_BINARY"
CLI_URL="https://github.com/Hareesh108/haruDB/releases/download/$DB_VERSION/$CLI_BINARY"

mkdir -p "$DIST_DIR"

echo "⬇️ Downloading HaruDB server ($SERVER_BINARY)..."
curl -fL "$SERVER_URL" -o "$DIST_DIR/harudb"

echo "⬇️ Downloading HaruDB CLI ($CLI_BINARY)..."
curl -fL "$CLI_URL" -o "$DIST_DIR/haru-cli"

# Make binaries executable (Linux/macOS)
if [[ "$OS" == "linux" || "$OS" == "darwin" ]]; then
    chmod +x "$DIST_DIR/harudb" "$DIST_DIR/haru-cli"
    sudo mv "$DIST_DIR/harudb" "$BIN_PATH/harudb"
    sudo mv "$DIST_DIR/haru-cli" "$BIN_PATH/haru-cli"
    echo "✅ HaruDB server and CLI installed to $BIN_PATH"
fi

# Windows instructions
if [[ "$OS" == mingw* || "$OS" == cygwin* || "$OS" == msys* || "$OS" == windowsnt ]]; then
    echo "✅ HaruDB server and CLI downloaded to $DIST_DIR"
    echo "Add $DIST_DIR to your PATH to use globally."
    echo "Run server: $DIST_DIR/harudb-windows.exe"
    echo "Run CLI:    $DIST_DIR/haru-cli-windows.exe"
    exit 0
fi

# -------------------------------
# Start server (Linux/macOS)
# -------------------------------
echo "🚀 Starting HaruDB server..."
nohup "$BIN_PATH/harudb" &>/dev/null &

# Wait for server to start
sleep 2

# Check if server port is open
if nc -z localhost 54321; then
    echo "HaruDB server is running on port 54321."
    echo "Connect using:"
    echo "  haru-cli  (for CLI)"
    # echo "  telnet localhost 54321  (basic connection)"
else
    echo "❌ HaruDB server failed to start."
fi

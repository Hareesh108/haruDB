#!/bin/bash

# HaruDB installer script - fully automated
# Usage: curl -sSL https://raw.githubusercontent.com/Hareesh108/haruDB/main/install-harudb.sh | bash

# Set version
VERSION="v0.0.1"

# Detect OS
OS="$(uname | tr '[:upper:]' '[:lower:]')"

# Set binary name and URL
if [[ "$OS" == "linux" ]]; then
    BINARY_NAME="harudb-linux"
elif [[ "$OS" == "darwin" ]]; then
    BINARY_NAME="harudb-macos"
elif [[ "$OS" == "mingw"* || "$OS" == "cygwin"* || "$OS" == "msys"* ]]; then
    BINARY_NAME="harudb-windows.exe"
else
    echo "❌ Unsupported OS: $OS"
    exit 1
fi

DOWNLOAD_URL="https://github.com/Hareesh108/haruDB/releases/download/$VERSION/$BINARY_NAME"

# Download binary
echo "⬇️ Downloading HaruDB ($BINARY_NAME)..."
curl -L "$DOWNLOAD_URL" -o "$BINARY_NAME"

# Make executable if Linux/macOS
if [[ "$OS" == "linux" || "$OS" == "darwin" ]]; then
    chmod +x "$BINARY_NAME"
    # Move to /usr/local/bin for global access
    sudo mv "$BINARY_NAME" /usr/local/bin/harudb
    echo "✅ HaruDB installed to /usr/local/bin/harudb"
    EXEC_CMD="harudb"
else
    echo "✅ HaruDB downloaded as $BINARY_NAME. Add it to PATH to use globally."
    EXEC_CMD="./$BINARY_NAME"
fi

# Run HaruDB
echo "🚀 Starting HaruDB..."
$EXEC_CMD &

echo "HaruDB is running on port 54321. Connect with:"
echo "telnet localhost 54321"

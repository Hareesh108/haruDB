#!/bin/bash
set -e

BINARY_PATH="/usr/local/bin/harudb"
DATA_DIR="$HOME/.harudb"  # optional data directory for future persistence

# Stop running HaruDB processes
PIDS=$(pgrep -f "$BINARY_PATH") || true
if [ -n "$PIDS" ]; then
    echo "Stopping HaruDB processes: $PIDS"
    kill -9 $PIDS
    echo "✅ HaruDB server stopped."
else
    echo "HaruDB server is not running."
fi

# Remove binary
if [ -f "$BINARY_PATH" ]; then
    echo "Removing HaruDB binary at $BINARY_PATH..."
    sudo rm -f "$BINARY_PATH"
    echo "✅ HaruDB binary removed."
else
    echo "❌ HaruDB binary not found."
fi

# Remove optional data directory
if [ -d "$DATA_DIR" ]; then
    echo "Removing HaruDB data directory at $DATA_DIR..."
    rm -rf "$DATA_DIR"
    echo "✅ Data directory removed."
fi

echo "🎉 HaruDB uninstalled successfully!"

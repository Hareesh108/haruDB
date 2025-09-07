#!/bin/bash
set -e

BINARY_PATH="/usr/local/bin/harudb"

# Stop HaruDB if running
PID=$(pgrep -f "$BINARY_PATH") || true
if [ -n "$PID" ]; then
  echo "Stopping HaruDB (PID $PID)..."
  kill -9 $PID
fi

# Remove binary
if [ -f "$BINARY_PATH" ]; then
  echo "Removing HaruDB binary at $BINARY_PATH..."
  sudo rm -f "$BINARY_PATH"
  echo "✅ HaruDB uninstalled successfully!"
else
  echo "❌ HaruDB is not installed."
fi

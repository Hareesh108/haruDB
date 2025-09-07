#!/bin/bash
set -e

BINARY_PATH="/usr/local/bin/harudb"
DATA_DIR="$HOME/.harudb"  # optional data directory for future persistence

echo "⚠️  This will stop HaruDB and remove the binary."
read -p "Are you sure you want to continue? (y/N): " confirm
if [[ ! "$confirm" =~ ^[Yy]$ ]]; then
  echo "❌ Uninstall canceled."
  exit 1
fi

# Stop HaruDB if running
PID=$(pgrep -f "$BINARY_PATH") || true
if [ -n "$PID" ]; then
  echo "Stopping HaruDB (PID $PID)..."
  kill -9 $PID
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

# Remove data directory if exists (optional)
if [ -d "$DATA_DIR" ]; then
  echo "Removing HaruDB data directory at $DATA_DIR..."
  rm -rf "$DATA_DIR"
  echo "✅ Data directory removed."
fi

echo "🎉 HaruDB uninstalled successfully!"

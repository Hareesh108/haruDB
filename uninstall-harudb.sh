#!/bin/bash
set -e

# Paths
BINARY_PATH="/usr/local/bin/harudb"
DATA_DIR="$HOME/.harudb"          # Optional data directory for future persistence
LOG_DIR="$HOME/.harudb_logs"      # Optional logs directory
TMP_DIR="/tmp/harudb"             # Optional temp files

echo "🚀 Uninstalling HaruDB..."

# Stop all running HaruDB processes
PIDS=$(pgrep -f "harudb") || true
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

# Remove data directory
if [ -d "$DATA_DIR" ]; then
    echo "Removing HaruDB data directory at $DATA_DIR..."
    rm -rf "$DATA_DIR"
    echo "✅ Data directory removed."
fi

# Remove logs directory
if [ -d "$LOG_DIR" ]; then
    echo "Removing HaruDB logs directory at $LOG_DIR..."
    rm -rf "$LOG_DIR"
    echo "✅ Logs directory removed."
fi

# Remove temporary files
if [ -d "$TMP_DIR" ]; then
    echo "Removing HaruDB temporary files at $TMP_DIR..."
    rm -rf "$TMP_DIR"
    echo "✅ Temporary files removed."
fi

echo "🎉 HaruDB fully uninstalled successfully!"

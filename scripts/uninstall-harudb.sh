#!/bin/bash
set -e

# Usage: curl -sSL https://raw.githubusercontent.com/Hareesh108/haruDB/main/scripts/uninstall-harudb.sh | bash


# Paths
BINARY_PATH="/usr/local/bin/harudb"
DATA_DIR="$HOME/.harudb"          # Optional data directory for future persistence
LOG_DIR="$HOME/.harudb_logs"      # Optional logs directory
TMP_DIR="/tmp/harudb"             # Optional temp files
PORT=54321                        # HaruDB default port

echo "🚀 Uninstalling HaruDB..."

# 1️⃣ Kill all HaruDB processes
PIDS=$(pgrep -f "harudb") || true
if [ -n "$PIDS" ]; then
    echo "Stopping HaruDB processes: $PIDS"
    for PID in $PIDS; do
        sudo kill -9 $PID
    done
    echo "✅ All HaruDB processes stopped."
else
    echo "HaruDB server is not running."
fi

# 2️⃣ Kill any process holding HaruDB port (handles active connections)
PORT_PIDS=$(lsof -ti tcp:$PORT) || true
if [ -n "$PORT_PIDS" ]; then
    echo "Terminating processes using port $PORT: $PORT_PIDS"
    for PID in $PORT_PIDS; do
        sudo kill -9 $PID
    done
    echo "✅ All processes on port $PORT terminated."
fi

# 3️⃣ Remove binary
if [ -f "$BINARY_PATH" ]; then
    echo "Removing HaruDB binary at $BINARY_PATH..."
    sudo rm -f "$BINARY_PATH"
    echo "✅ HaruDB binary removed."
else
    echo "❌ HaruDB binary not found."
fi

# 4️⃣ Remove data directory
if [ -d "$DATA_DIR" ]; then
    echo "Removing HaruDB data directory at $DATA_DIR..."
    rm -rf "$DATA_DIR"
    echo "✅ Data directory removed."
fi

# 5️⃣ Remove logs directory
if [ -d "$LOG_DIR" ]; then
    echo "Removing HaruDB logs directory at $LOG_DIR..."
    rm -rf "$LOG_DIR"
    echo "✅ Logs directory removed."
fi

# 6️⃣ Remove temporary files
if [ -d "$TMP_DIR" ]; then
    echo "Removing HaruDB temporary files at $TMP_DIR..."
    rm -rf "$TMP_DIR"
    echo "✅ Temporary files removed."
fi

echo "🎉 HaruDB fully uninstalled successfully!"

#!/bin/bash
set -e

# Usage: curl -sSL https://raw.githubusercontent.com/Hareesh108/haruDB/main/scripts/uninstall-harudb.sh | bash

# Paths
SERVER_BINARY="/usr/local/bin/harudb"
CLI_BINARY="/usr/local/bin/haru-cli"
DATA_DIR="$HOME/.harudb"          # Optional data directory for persistence
LOG_DIR="$HOME/.harudb_logs"      # Optional logs directory
TMP_DIR="/tmp/harudb"             # Optional temp files
PORT=54321                        # HaruDB default port

echo "🚀 Uninstalling HaruDB..."

# 1️⃣ Kill all HaruDB server processes
SERVER_PIDS=$(pgrep -f "harudb") || true
if [ -n "$SERVER_PIDS" ]; then
    echo "Stopping HaruDB server processes: $SERVER_PIDS"
    for PID in $SERVER_PIDS; do
        sudo kill -9 $PID
    done
    echo "✅ All HaruDB server processes stopped."
else
    echo "HaruDB server is not running."
fi

# 2️⃣ Kill all HaruDB CLI processes
CLI_PIDS=$(pgrep -f "haru-cli") || true
if [ -n "$CLI_PIDS" ]; then
    echo "Stopping HaruDB CLI processes: $CLI_PIDS"
    for PID in $CLI_PIDS; do
        sudo kill -9 $PID
    done
    echo "✅ All HaruDB CLI processes stopped."
fi

# 3️⃣ Kill any process holding HaruDB port (handles active connections)
PORT_PIDS=$(lsof -ti tcp:$PORT) || true
if [ -n "$PORT_PIDS" ]; then
    echo "Terminating processes using port $PORT: $PORT_PIDS"
    for PID in $PORT_PIDS; do
        sudo kill -9 $PID
    done
    echo "✅ All processes on port $PORT terminated."
fi

# 4️⃣ Remove server binary
if [ -f "$SERVER_BINARY" ]; then
    echo "Removing HaruDB server binary at $SERVER_BINARY..."
    sudo rm -f "$SERVER_BINARY"
    echo "✅ HaruDB server binary removed."
else
    echo "❌ HaruDB server binary not found."
fi

# 5️⃣ Remove CLI binary
if [ -f "$CLI_BINARY" ]; then
    echo "Removing HaruDB CLI binary at $CLI_BINARY..."
    sudo rm -f "$CLI_BINARY"
    echo "✅ HaruDB CLI binary removed."
else
    echo "❌ HaruDB CLI binary not found."
fi

# 6️⃣ Remove data directory
if [ -d "$DATA_DIR" ]; then
    echo "Removing HaruDB data directory at $DATA_DIR..."
    rm -rf "$DATA_DIR"
    echo "✅ Data directory removed."
fi

# 7️⃣ Remove logs directory
if [ -d "$LOG_DIR" ]; then
    echo "Removing HaruDB logs directory at $LOG_DIR..."
    rm -rf "$LOG_DIR"
    echo "✅ Logs directory removed."
fi

# 8️⃣ Remove temporary files
if [ -d "$TMP_DIR" ]; then
    echo "Removing HaruDB temporary files at $TMP_DIR..."
    rm -rf "$TMP_DIR"
    echo "✅ Temporary files removed."
fi

echo "🎉 HaruDB fully uninstalled successfully!"

#!/bin/bash

# HaruDB Build Script
# Build server and CLI binaries for Linux, macOS, and Windows
# Usage: ./build-binaries.sh

set -e

DB_VERSION="v0.0.3"
DIST_DIR="./binaries"

# Create dist directory
mkdir -p "$DIST_DIR"

echo "⚡ Building HaruDB binaries (Server + CLI)..."

# Define targets
declare -A TARGETS=(
  ["linux_amd64"]="linux amd64"
  ["darwin_amd64"]="darwin amd64"
  ["windows_amd64"]="windows amd64"
)

# Loop over targets
for target in "${!TARGETS[@]}"; do
  IFS=' ' read -r GOOS GOARCH <<< "${TARGETS[$target]}"

  echo "🔹 Building server for $GOOS/$GOARCH..."
  OUT="$DIST_DIR/harudb-$target"
  [[ "$GOOS" == "windows" ]] && OUT="$OUT.exe"
  GOOS=$GOOS GOARCH=$GOARCH go build -o "$OUT" ./../cmd/server

  echo "🔹 Building CLI for $GOOS/$GOARCH..."
  CLI_OUT="$DIST_DIR/haru-cli-$target"
  [[ "$GOOS" == "windows" ]] && CLI_OUT="$CLI_OUT.exe"
  GOOS=$GOOS GOARCH=$GOARCH go build -o "$CLI_OUT" ./../cmd/cli
done

echo "✅ Build complete. Binaries stored in $DIST_DIR:"
ls -lh "$DIST_DIR"

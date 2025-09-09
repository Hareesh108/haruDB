#!/bin/bash

# HaruDB Build Script
# Build binaries for Linux, macOS, and Windows
# Usage: ./build-binaries.sh

set -e

DB_VERSION="v0.0.3"
DIST_DIR="./dist"

# Create dist directory
mkdir -p $DIST_DIR

echo "⚡ Building HaruDB binaries..."

# Linux
echo "Building for Linux..."
GOOS=linux GOARCH=amd64 go build -o $DIST_DIR/harudb-linux ./../cmd/server

# macOS
echo "Building for macOS..."
GOOS=darwin GOARCH=amd64 go build -o $DIST_DIR/harudb-macos ./../cmd/server

# Windows
echo "Building for Windows..."
GOOS=windows GOARCH=amd64 go build -o $DIST_DIR/harudb-windows.exe ./../cmd/server

echo "✅ Build complete. Binaries stored in $DIST_DIR:"

ls -lh $DIST_DIR

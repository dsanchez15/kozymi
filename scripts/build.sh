#!/bin/bash
set -e

VERSION=${1:-latest}
OUTPUT=dist/

mkdir -p "$OUTPUT"

echo "Building Komyzi v${VERSION}..."

# Clean previous builds
rm -rf "$OUTPUT"*

# Windows AMD64
echo "Building Windows AMD64..."
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o "$OUTPUT/komyzi-windows-amd64.exe" ./cmd/cli

# Windows ARM64
echo "Building Windows ARM64..."
GOOS=windows GOARCH=arm64 go build -ldflags="-s -w" -o "$OUTPUT/komyzi-windows-arm64.exe" ./cmd/cli

# macOS AMD64
echo "Building macOS AMD64..."
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o "$OUTPUT/komyzi-darwin-amd64" ./cmd/cli

# macOS ARM64
echo "Building macOS ARM64..."
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o "$OUTPUT/komyzi-darwin-arm64" ./cmd/cli

# Linux AMD64
echo "Building Linux AMD64..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o "$OUTPUT/komyzi-linux-amd64" ./cmd/cli

# Linux ARM64
echo "Building Linux ARM64..."
GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o "$OUTPUT/komyzi-linux-arm64" ./cmd/cli

echo "Build complete! Binaries in $OUTPUT"
ls -la "$OUTPUT"
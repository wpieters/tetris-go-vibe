#!/bin/bash

# Exit on any error
set -e

# Docker repository
REPO="duhblinn/tetris-go-vibe"

# Version tag (using date and time)
VERSION=$(date +%Y%m%d-%H%M%S)

# Create dist directory
rm -rf dist
mkdir -p dist

# Build binaries for different platforms
echo "Building binaries..."

# Linux AMD64
echo "Building for Linux AMD64..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/tetris-linux-amd64

# Linux ARM64
echo "Building for Linux ARM64..."
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o dist/tetris-linux-arm64

# macOS AMD64
echo "Building for macOS AMD64..."
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o dist/tetris-darwin-amd64

# macOS ARM64
echo "Building for macOS ARM64..."
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o dist/tetris-darwin-arm64

# Enable Docker BuildKit
export DOCKER_BUILDKIT=1

# Set up multi-platform builder if it doesn't exist
if ! docker buildx inspect multibuilder >/dev/null 2>&1; then
    docker buildx create --name multibuilder --driver docker-container --bootstrap
fi

# Use the multi-platform builder
docker buildx use multibuilder

# Build and push Linux Docker images
echo "Building and pushing Linux Docker images..."

# Build for Linux platforms

echo "Building Docker images..."
docker login
docker buildx build \
    --platform linux/amd64,linux/arm64 \
    --tag ${REPO}:latest \
    --push \
    --load \
    .

echo "Build complete!"
echo "Native binaries available in ./dist/"
echo "Docker images pushed"


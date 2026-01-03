#!/bin/bash

# Get version info from git
export VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
export COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
export BUILD_TIME=$(date -u +'%Y-%m-%dT%H:%M:%SZ')

echo "Building with:"
echo "  VERSION:    $VERSION"
echo "  COMMIT:     $COMMIT"
echo "  BUILD_TIME: $BUILD_TIME"
echo ""

# Build with docker-compose
docker-compose -f compose.dev.yml up -d --build app

#!/bin/bash

# build script for creating AppImage
set -e

echo "Building Job Visualizer AppImage..."

# detect system architecture
ARCH=$(uname -m)
case $ARCH in
    x86_64)
        GOARCH="amd64"
        ARCH_NAME="x86_64"
        ;;
    aarch64|arm64)
        GOARCH="arm64"
        ARCH_NAME="aarch64"
        ;;
    armv7l)
        GOARCH="arm"
        ARCH_NAME="armhf"
        ;;
    i386|i686)
        GOARCH="386"
        ARCH_NAME="i386"
        ;;
    *)
        echo "Warning: Unknown architecture '$ARCH', using as-is"
        GOARCH="$ARCH"
        ARCH_NAME="$ARCH"
        ;;
esac

echo "Detected architecture: $ARCH_NAME"

# prompt user for version number
read -p "Enter version number (e.g., v0.1.0): " VERSION

# validate version format
if [[ ! $VERSION =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "Error: Version must be in format v0.0.0 (e.g., v1.2.3)"
    exit 1
fi

echo "Building version: $VERSION"

# build the Go application
echo "Building Go application..."
CGO_ENABLED=1 GOOS=linux GOARCH=$GOARCH go build -ldflags="-s -w" -o AppDir/job-visualizer-linux ./cmd/app

# make AppRun executable
chmod +x AppDir/AppRun

# create the AppImage using system appimagetool
echo "Creating AppImage..."
appimagetool AppDir job-visualizer-${VERSION}-${ARCH_NAME}.AppImage

# move the created AppImage to releases directory
echo "Moving AppImage to releases directory..."
mv job-visualizer-${VERSION}-${ARCH_NAME}.AppImage releases/

# make the AppImage executable
echo "Making AppImage executable..."
chmod +x releases/job-visualizer-${VERSION}-${ARCH_NAME}.AppImage

echo "AppImage created successfully!"
echo "Location: releases/job-visualizer-${VERSION}-${ARCH_NAME}.AppImage"

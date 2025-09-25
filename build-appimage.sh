#!/bin/bash

# build script for creating AppImage
set -e

echo "Building Job Visualizer AppImage..."

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
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o AppDir/job-visualizer-linux ./cmd/app

# make AppRun executable
chmod +x AppDir/AppRun

# create the AppImage using system appimagetool
echo "Creating AppImage..."
appimagetool AppDir job-visualizer-${VERSION}-x86_64.AppImage

# move the created AppImage to releases directory
echo "Moving AppImage to releases directory..."
mv job-visualizer-${VERSION}-x86_64.AppImage releases/

# make the AppImage executable
echo "Making AppImage executable..."
chmod +x releases/job-visualizer-${VERSION}-x86_64.AppImage

echo "AppImage created successfully!"
echo "Location: releases/job-visualizer-${VERSION}-x86_64.AppImage"

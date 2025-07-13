#!/bin/bash

APP_NAME="mywailsapp"
VERSION=$(cat ./build/version.txt)
DIST_DIR="./build/dist"

mkdir -p "$DIST_DIR"

# Define target platforms
PLATFORMS=("windows/amd64" "linux/amd64" "darwin/universal")

for PLATFORM in "${PLATFORMS[@]}"
do
    echo "📦 Building for $PLATFORM..."
    wails build -platform "$PLATFORM" -clean -o "$DIST_DIR"
done

echo "✅ Build complete!"
#!/bin/bash

# === CONFIG ===
APP_NAME="xpack"
BINARY_PATH="./bin/xpack"
VERSION_FILE="./VERSION"
AVAILABLE_ARCHITECTURES=("amd64" "arm64" "i386")

# === Get version from VERSION file ===
if [ -f "$VERSION_FILE" ]; then
  VERSION=$(cat "$VERSION_FILE" | tr -d '[:space:]')
else
  echo "❌ VERSION file not found in project Core"
  exit 1
fi

# === Check binary ===
if [ ! -f "$BINARY_PATH" ]; then
  echo "❌ Binary not found at $BINARY_PATH"
  exit 1
fi

# === Build packages for each architecture ===
for ARCH in "${AVAILABLE_ARCHITECTURES[@]}"; do
  echo "🔨 Building packages for architecture: $ARCH"

  "$BINARY_PATH" -i "$BINARY_PATH" -app "$APP_NAME" -arch "$ARCH" -v "$VERSION"

  if [ $? -ne 0 ]; then
    echo "❌ Failed to build packages for architecture: $ARCH"
    exit 1
  fi

  echo "✅ Packages created for architecture: $ARCH"
done

echo "🎉 All packages built successfully and saved to dist folder"

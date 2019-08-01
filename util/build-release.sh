#!/usr/bin/env bash
# Build multi-OS release

set -e

################################################################################
# Configuration
BUILD_PATH="/tmp/paintstripper-release"
################################################################################

SCRIPTDIR=$(dirname $0)
GITHUB_PATH="github.com/sparkcodeuk/paintstripper"

# Resolve version
VERSION=$(grep 'const\s*version' "$SCRIPTDIR/../version.go"|sed -e 's/.*"\(.*\)"/\1/')
if [ -z "$VERSION" ]; then
    echo "ERROR: unable to resolve application version"
    exit 1
fi

# Check for existing builds
if [ -e "$BUILD_PATH" ]; then
    echo "ERROR: BUILD_PATH [$BUILD_PATH] already exists (probably from a previous build), exiting"
    exit 1
fi

# Prepare build area
mkdir -p \
    "$BUILD_PATH/src/$GITHUB_PATH" \
    "$BUILD_PATH/pkg" \
    "$BUILD_PATH/bin" \
    "$BUILD_PATH/builds"

git clone "https://${GITHUB_PATH}.git" "$BUILD_PATH/src/$GITHUB_PATH"

# Run builds
cd "$BUILD_PATH"

echo "Building linux..."
GOPATH="$BUILD_PATH" GOOS="linux" GOARCH="amd64" go build -v "$GITHUB_PATH"
tar zcvf "$BUILD_PATH/builds/paintstripper-linux-$VERSION.tgz" "./paintstripper"
rm -f "./paintstripper"

echo "Building OSX..."
GOPATH="$BUILD_PATH" GOOS="darwin" GOARCH="amd64" go build -v "$GITHUB_PATH"
tar zcvf "$BUILD_PATH/builds/paintstripper-osx-$VERSION.tgz" "./paintstripper"
rm -f "./paintstripper"

echo "Build release complete, files available at: $BUILD_PATH/builds"

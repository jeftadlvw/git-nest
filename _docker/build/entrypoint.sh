#!/usr/bin/env bash

# variables
SRC_DIR="/git-nest/src"
BIN_DIR="/git-nest/build"

SCRIPT_VERSION="1.0.1"
APP_NAME="git-nest"
CURRENT_TIMESTAMP=$(date +%s)
WORKING_DIR=$(pwd)
TAB="   "

echo "git-nest build script version $SCRIPT_VERSION"

# perform variable and input validation
if [ -r "$SRC_DIR" ] && [ -w "$SRC_DIR" ]; then
    echo "error: source directory is not read-only."
    echo "Mount it as readonly to proceed."
    exit 1
fi

if [ -z "$GIT_NEST_BUILD_VERSION" ]; then
    echo "environment variable 'GIT_NEST_BUILD_VERSION‘ not set."
    exit 1
fi

if [ -z "$GIT_NEST_BUILD_COMMIT_SHA" ]; then
    echo "environment variable 'GIT_NEST_BUILD_COMMIT_SHA' not set."
    exit 1
fi

echo "Environment checks finished."

# prepare key/value pairs for compile-time value injection
# package base
INJECT_BASE="github.com/jeftadlvw/git-nest/internal"

# application version
INJECT_VERSION_KEY="$INJECT_BASE/constants.version"
INJECT_VERSION_VALUE="$GIT_NEST_BUILD_VERSION"

# commit hash
INJECT_COMMIT_KEY="$INJECT_BASE/constants.ref"
INJECT_COMMIT_VALUE="$GIT_NEST_BUILD_COMMIT_SHA"

# compilation time
INJECT_COMPILE_TIME_KEY="$INJECT_BASE/constants.compilationTimestampStr"
INJECT_COMPILE_TIME_VALUE="${GIT_NEST_COMPILE_TIME:-$CURRENT_TIMESTAMP}"    # GIT_NEST_COMPILE_TIME is not required

# whether this build is ephemeral
INJECT_EPHEMERAL_BUILD_KEY="$INJECT_BASE/constants.ephemeralBuildStr"
INJECT_EPHEMERAL_BUILD_VALUE="false"

# define GOARCH and GOOS matrix
BUILD_TARGETS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
    "windows/arm64"
)

# print build settings
echo "Starting build process using the configuration below:"
echo "$TAB VERSION: $INJECT_VERSION_VALUE"
echo "$TAB COMMIT: $INJECT_COMMIT_VALUE"
echo "$TAB TIME: $INJECT_COMPILE_TIME_VALUE"
echo

echo "$TAB Build targets:"
for target in "${BUILD_TARGETS[@]}"; do
    echo "$TAB$TAB - $target"
done
echo

# download dependencies
echo "Downloading dependencies..."
cd $SRC_DIR && go get
cd $WORKING_DIR
echo

# run tests
if [ -z "${SKIP_TESTS}" ]; then
    echo "Running tests..."
    cd $SRC_DIR && go test ./.../tests
    if [ $? -ne 0 ]; then
        echo "Tests failed, exiting."
        exit 1
    fi

    cd $WORKING_DIR
    echo "Tests ran successful."
else
    echo "Tests skipped."
fi
echo

# compile and rename binaries, ensure windows binaries have .exe suffix
# and calculate and format checksum output

# remove pre-existing checksum-file
CHECKSUM_FILE="$BIN_DIR/checksums.txt"
rm -f $CHECKSUM_FILE

# Loop through each target and build
cd $SRC_DIR
for target in "${BUILD_TARGETS[@]}"; do
    # Split the target string into OS and architecture
    IFS='/' read -ra split_target <<< "$target"
    OS="${split_target[0]}"
    ARCH="${split_target[1]}"

    # Build the codebase for the current target
    OUTPUT_FILE="git-nest_$OS-$ARCH"

    if [ "$OS" = "windows" ]; then
        OUTPUT_FILE="$OUTPUT_FILE.exe"
    fi
    OUTPUT_PATH="$BIN_DIR/$OUTPUT_FILE"

    BUILDING_INFO_STR="Building target '$OS/$ARCH' to $OUTPUT_PATH..."
    echo -n $BUILDING_INFO_STR

    BUILD_OUTPUT= $( \
        GOOS=$OS \
        GOARCH=$ARCH \
        go build \
            -o $OUTPUT_PATH \
            -buildvcs=false \
            -ldflags "\
                -X $INJECT_VERSION_KEY=$INJECT_VERSION_VALUE \
                -X $INJECT_COMMIT_KEY=$INJECT_COMMIT_VALUE \
                -X $INJECT_COMPILE_TIME_KEY=$INJECT_COMPILE_TIME_VALUE \
                -X $INJECT_EPHEMERAL_BUILD_KEY=$INJECT_EPHEMERAL_BUILD_VALUE \
                " \
    )

    if [ $? -ne 0 ]; then
        echo -ne "\r$BUILDING_INFO_STR (fail)"
        echo $BUILD_OUTPUT
        exit 1
    fi

    CHECKSUM=$(sha256sum $OUTPUT_PATH | awk '{ print $1 }')
    echo -e "# $OUTPUT_FILE\n$CHECKSUM\n" >> $CHECKSUM_FILE
    echo -ne "\r$BUILDING_INFO_STR (success)"
    echo
done
echo

echo "All builds completed and located under $BIN_DIR"

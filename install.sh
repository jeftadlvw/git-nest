#!/bin/bash
set -euo pipefail

# variables
PLATFORM=$(uname -ms)
DOWNLOAD_TARGET=""

REPOSITORY_NAME="jeftadlvw/git-nest"
REPOSITORY="https://github.com/$REPOSITORY_NAME"

INSTALL_DIR=$HOME/.local/bin
BINARY_NAME="git-nest"

TAB="   "
echo "Installing git_nest to $INSTALL_DIR/$BINARY_NAME."
echo ""

if ! which git &> /dev/null; then
    echo "Warning: git is not installed."
fi

case $PLATFORM in
'Darwin x86_64')
    DOWNLOAD_TARGET=darwin-amd64
    ;;
'Darwin arm64')
    DOWNLOAD_TARGET=darwin-arm64
    ;;
'Linux x86_64')
    DOWNLOAD_TARGET=linux-amd64
    ;;
'Linux aarch64')
    DOWNLOAD_TARGET=linux-arm64
    ;;
*)
    echo "error: no official binary for target '$PLATFORM'"
    exit 1
    ;;
esac

# get the latest release information from GitHub API
LATEST_TAG=$(curl -s "https://api.github.com/repos/$REPOSITORY_NAME/releases/latest" | awk -F '"tag_name": "' '{if ($2) print $2}' | awk -F'"' '{print $1}')
if [ $? -ne 0 ] || [ -z "$LATEST_TAG" ]; then
    echo "error: unable to retrieve latest release tag"
    exit 1
fi

# create installation directories
mkdir -p $INSTALL_DIR
if [ $? -ne 0 ]; then
    echo "error: unable to create installation directory $INSTALL_DIR"
    exit 1
fi

# download binary to temporary directory
DOWNLOAD_URL="$REPOSITORY/releases/download/$LATEST_TAG/git-nest_$DOWNLOAD_TARGET"
echo "Downloading from $DOWNLOAD_URL"

curl -fL $DOWNLOAD_URL -o $INSTALL_DIR/$BINARY_NAME
if [ $? -ne 0 ]; then
    echo "error: unable to download binary from $ASSET_URL and install it to $INSTALL_DIR/$BINARY_NAME"
    exit 1
fi

# make binary executable
chmod +x $INSTALL_DIR/$BINARY_NAME
if [ $? -ne 0 ]; then
    echo "error: unable to make $INSTALL_DIR/$BINARY_NAME executable"
    exit 1
fi

echo ""
echo "git-nest has been successfully installed to $INSTALL_DIR"
echo ""

# check if installation directory is in PATH
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo "Notice: $INSTALL_DIR is not in your PATH."
    echo "Add it to your PATH and restart your shell in order to use git-nest. E.g:"
    echo ""
    echo "$TAB echo -e \"\n\n#local binaries\nPATH=\\\"\\\$PATH:$INSTALL_DIR\\\"\n\" >> ~/.bashrc"
    echo ""
fi

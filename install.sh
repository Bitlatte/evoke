#!/bin/sh

# This script downloads and installs the latest release of Evoke for your system.
# It can also be used to update an existing installation.

set -e

# Get the latest version from the GitHub API
LATEST_VERSION=$(curl -s "https://api.github.com/repos/Bitlatte/evoke/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_VERSION" ]; then
  echo "Could not find the latest version of Evoke."
  exit 1
fi

# Determine the OS and architecture
OS="$(uname -s)"
ARCH="$(uname -m)"

case $OS in
  Linux)
    OS_TYPE="Linux"
    ;;
  Darwin)
    OS_TYPE="Darwin"
    ;;
  *)
    echo "Unsupported OS: $OS"
    exit 1
    ;;
esac

case $ARCH in
  x86_64)
    ARCH_TYPE="x86_64"
    ;;
  arm64 | aarch64)
    ARCH_TYPE="arm64"
    ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

# Construct the download URL
DOWNLOAD_URL="https://github.com/Bitlatte/evoke/releases/download/${LATEST_VERSION}/evoke_${OS_TYPE}_${ARCH_TYPE}.tar.gz"

# Download and extract the binary
echo "Downloading Evoke ${LATEST_VERSION} for ${OS_TYPE}/${ARCH_TYPE}..."
curl -L -o evoke.tar.gz "$DOWNLOAD_URL"
tar -xzf evoke.tar.gz
rm evoke.tar.gz

# Remove existing installation if it exists
if [ -f "/usr/local/bin/evoke" ]; then
    echo "Removing existing Evoke installation..."
    sudo rm /usr/local/bin/evoke
fi

# Install the binary
echo "Installing Evoke to /usr/local/bin..."
sudo mv evoke /usr/local/bin/evoke
chmod +x /usr/local/bin/evoke

# Install the man page
if [ -f "man/evoke.1" ]; then
    echo "Installing man page..."
    if [ ! -d "/usr/local/share/man/man1" ]; then
        sudo mkdir -p /usr/local/share/man/man1
    fi
    sudo mv man/evoke.1 /usr/local/share/man/man1/evoke.1
    rm -r man
fi

echo "Evoke has been installed successfully."
echo "Run 'evoke --version' to verify the installation."

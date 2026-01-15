#!/usr/bin/env bash

# Script to run installation script tests
# Automatically detects the platform and runs appropriate tests

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEST_DIR="$SCRIPT_DIR/test"

echo "==> Installation Script Test Runner"
echo ""

# Detect platform
OS="$(uname -s)"
case "$OS" in
    Linux*)
        PLATFORM="Linux"
        TEST_FILE="install-hpe-provider-linux.bats"
        ;;
    Darwin*)
        PLATFORM="macOS"
        TEST_FILE="install-hpe-provider-macos.bats"
        ;;
    MINGW*|MSYS*|CYGWIN*)
        PLATFORM="Windows"
        echo "Windows detected. Please run tests using PowerShell:"
        echo "  cd scripts/test"
        echo "  Invoke-Pester -Path .\install-hpe-provider-windows.Tests.ps1"
        exit 0
        ;;
    *)
        echo "Error: Unsupported platform: $OS"
        exit 1
        ;;
esac

echo "Platform: $PLATFORM"
echo ""

# Check if BATS is installed
if ! command -v bats &> /dev/null; then
    echo "Error: BATS (Bash Automated Testing System) is not installed."
    echo ""
    echo "To install BATS:"
    echo ""
    if [ "$PLATFORM" = "macOS" ]; then
        echo "  brew install bats-core"
    else
        echo "  # On Ubuntu/Debian:"
        echo "  sudo apt-get install bats"
        echo ""
        echo "  # Or install from source:"
        echo "  git clone https://github.com/bats-core/bats-core.git"
        echo "  cd bats-core"
        echo "  sudo ./install.sh /usr/local"
    fi
    echo ""
    exit 1
fi

echo "==> Running tests for $PLATFORM installation script..."
echo ""

# Run BATS tests
bats "$TEST_DIR/$TEST_FILE"

echo ""
echo "==> Tests completed!"

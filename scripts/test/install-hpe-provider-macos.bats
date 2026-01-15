#!/usr/bin/env bats

# Test suite for install-hpe-provider-macos.sh
# These tests verify the installation script works correctly

setup() {
    # Create a temporary test directory
    export TEST_HOME=$(mktemp -d)
    export HOME="$TEST_HOME"
    export SCRIPT_DIR="$BATS_TEST_DIRNAME/.."
}

teardown() {
    # Clean up test directory
    rm -rf "$TEST_HOME"
}

@test "script creates correct directory structure" {
    # Set a specific version to avoid downloading latest
    export VERSION="v0.1.0"
    
    # Run the installation script
    run bash "$SCRIPT_DIR/install-hpe-provider-macos.sh"
    
    # Check that the script succeeded
    echo "Exit status: $status"
    echo "Output: $output"
    [ "$status" -eq 0 ] || (echo "Script failed with status $status"; echo "$output"; return 1)
    
    # Check that the directory was created
    [ -d "$HOME/.terraform.d/plugins/registry.terraform.io/hpe/hpe/0.1.0" ] || (echo "Directory not created: $HOME/.terraform.d/plugins/registry.terraform.io/hpe/hpe/0.1.0"; ls -la "$HOME/.terraform.d/plugins/registry.terraform.io/hpe/hpe/" 2>&1 || echo "Parent directory doesn't exist"; return 1)
}

@test "script detects architecture correctly" {
    export VERSION="v0.1.0"
    
    run bash "$SCRIPT_DIR/install-hpe-provider-macos.sh"
    
    # Verify architecture-specific directory exists
    if [[ $(uname -m) == 'x86_64' ]]; then
        [ -d "$HOME/.terraform.d/plugins/registry.terraform.io/hpe/hpe/0.1.0/darwin_amd64" ]
    elif [[ $(uname -m) == 'arm64' ]]; then
        [ -d "$HOME/.terraform.d/plugins/registry.terraform.io/hpe/hpe/0.1.0/darwin_arm64" ]
    fi
}

@test "script downloads and extracts provider binary" {
    export VERSION="v0.1.0"
    
    run bash "$SCRIPT_DIR/install-hpe-provider-macos.sh"
    
    [ "$status" -eq 0 ]
    
    # Check that the provider binary exists
    arch=$(uname -m)
    if [[ $arch == 'x86_64' ]]; then
        arch="amd64"
    fi
    
    provider_path="$HOME/.terraform.d/plugins/registry.terraform.io/hpe/hpe/0.1.0/darwin_${arch}/terraform-provider-hpe_v0.1.0"
    [ -f "$provider_path" ]
}

@test "script removes zip file after extraction" {
    export VERSION="v0.1.0"
    
    run bash "$SCRIPT_DIR/install-hpe-provider-macos.sh"
    
    [ "$status" -eq 0 ]
    
    # Verify no zip files remain
    run find "$HOME/.terraform.d" -name "*.zip"
    [ "$status" -eq 0 ]
    [ -z "$output" ]
}

@test "script works without explicit version (downloads latest)" {
    # Skip this test if ggrep is not installed
    if ! command -v ggrep &> /dev/null; then
        skip "ggrep not installed (brew install grep)"
    fi
    
    # Don't set VERSION - let it fetch latest
    unset VERSION
    
    run bash "$SCRIPT_DIR/install-hpe-provider-macos.sh"
    
    # Should succeed
    echo "Exit status: $status"
    echo "Output: $output"
    [ "$status" -eq 0 ] || (echo "Script failed with status $status"; echo "$output"; return 1)
    
    # Should have created the directory structure
    [ -d "$HOME/.terraform.d/plugins/registry.terraform.io/hpe/hpe" ] || (echo "Directory not created: $HOME/.terraform.d/plugins/registry.terraform.io/hpe/hpe"; ls -la "$HOME/.terraform.d/plugins/registry.terraform.io/hpe/" 2>&1 || echo "Parent directory doesn't exist"; return 1)
}

@test "script can reinstall over existing installation" {
    export VERSION="v0.1.0"
    
    # Install once
    run bash "$SCRIPT_DIR/install-hpe-provider-macos.sh"
    [ "$status" -eq 0 ]
    
    # Install again (should succeed and update)
    run bash "$SCRIPT_DIR/install-hpe-provider-macos.sh"
    [ "$status" -eq 0 ]
}

@test "script fails gracefully with invalid version" {
    export VERSION="v999.999.999"
    
    run bash "$SCRIPT_DIR/install-hpe-provider-macos.sh"
    
    # Should fail (curl will return non-zero)
    [ "$status" -ne 0 ]
}

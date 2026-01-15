#!/usr/bin/env bats

# Test suite for install-hpe-provider.sh (Linux)
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
    run bash "$SCRIPT_DIR/install-hpe-provider.sh"
    
    # Check that the script succeeded
    echo "Exit status: $status"
    echo "Output: $output"
    [ "$status" -eq 0 ] || (echo "Script failed with status $status"; echo "$output"; return 1)
    
    # Check that the directory was created
    [ -d "$HOME/.local/share/terraform/plugins/registry.terraform.io/hpe/hpe/0.1.0/linux_amd64" ] || (echo "Directory not created: $HOME/.local/share/terraform/plugins/registry.terraform.io/hpe/hpe/0.1.0/linux_amd64"; ls -la "$HOME/.local/share/terraform/plugins/registry.terraform.io/hpe/hpe/0.1.0/" 2>&1 || echo "Parent directory doesn't exist"; return 1)
}

@test "script downloads and extracts provider binary" {
    export VERSION="v0.1.0"
    
    run bash "$SCRIPT_DIR/install-hpe-provider.sh"
    
    [ "$status" -eq 0 ]
    
    # Check that the provider binary exists
    provider_path="$HOME/.local/share/terraform/plugins/registry.terraform.io/hpe/hpe/0.1.0/linux_amd64/terraform-provider-hpe_v0.1.0"
    [ -f "$provider_path" ]
}

@test "script removes zip file after extraction" {
    export VERSION="v0.1.0"
    
    run bash "$SCRIPT_DIR/install-hpe-provider.sh"
    
    [ "$status" -eq 0 ]
    
    # Verify no zip files remain
    run find "$HOME/.local/share/terraform" -name "*.zip"
    [ "$status" -eq 0 ]
    [ -z "$output" ]
}

@test "script works without explicit version (downloads latest)" {
    # Don't set VERSION - let it fetch latest
    unset VERSION
    
    run bash "$SCRIPT_DIR/install-hpe-provider.sh"
    
    # Should succeed
    echo "Exit status: $status"
    echo "Output: $output"
    [ "$status" -eq 0 ] || (echo "Script failed with status $status"; echo "$output"; return 1)
    
    # Should have created the directory structure
    [ -d "$HOME/.local/share/terraform/plugins/registry.terraform.io/hpe/hpe" ] || (echo "Directory not created: $HOME/.local/share/terraform/plugins/registry.terraform.io/hpe/hpe"; ls -la "$HOME/.local/share/terraform/plugins/registry.terraform.io/hpe/" 2>&1 || echo "Parent directory doesn't exist"; return 1)
}

@test "script can reinstall over existing installation" {
    export VERSION="v0.1.0"
    
    # Install once
    run bash "$SCRIPT_DIR/install-hpe-provider.sh"
    [ "$status" -eq 0 ]
    
    # Install again (should succeed and update)
    run bash "$SCRIPT_DIR/install-hpe-provider.sh"
    [ "$status" -eq 0 ]
}

@test "script fails gracefully with invalid version" {
    export VERSION="v999.999.999"
    
    run bash "$SCRIPT_DIR/install-hpe-provider.sh"
    
    # Should fail (curl will return non-zero)
    [ "$status" -ne 0 ]
}

@test "script extracts binary with correct permissions" {
    export VERSION="v0.1.0"
    
    run bash "$SCRIPT_DIR/install-hpe-provider.sh"
    
    [ "$status" -eq 0 ]
    
    # Check that the binary is executable
    provider_path="$HOME/.local/share/terraform/plugins/registry.terraform.io/hpe/hpe/0.1.0/linux_amd64/terraform-provider-hpe_v0.1.0"
    [ -x "$provider_path" ]
}

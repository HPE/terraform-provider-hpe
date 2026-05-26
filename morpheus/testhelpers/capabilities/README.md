# Capability-Based Test Activation

This package provides capability-based test activation for acceptance tests. Tests declare their infrastructure requirements, and only run if the target system has those capabilities available.

## Quick Start

### Running Tests on Shared Hardware

To run only tests that are safe for shared VMware/Alletra/NSX infrastructure:

```bash
export TF_ACC_CAPABILITIES="vmware,alletra,nsxt,nsxv"
go test ./morpheus/...
```

### Running All Tests

To run all tests, explicitly enable all capabilities:

```bash
export TF_ACC_CAPABILITIES="all,vmware,alletra,nsxt,nsxv,aws,azure,gcp,kubernetes,ansible"
go test ./morpheus/...
```

### Safe Default

When `TF_ACC_CAPABILITIES` is **not set**, no capabilities are available and **no tests run**. This is a safe default since tests can be destructive.

```bash
# No tests will run - safe default
unset TF_ACC_CAPABILITIES
go test ./morpheus/...
```

## How It Works

Every acceptance test declares its infrastructure requirements as the first statement:

```go
func TestAccMorpheusNetworkRouterNSXT(t *testing.T) {
    if capabilities.Missing(t, capabilities.NSXT, capabilities.NetworkRouter) {
        return
    }
    t.Parallel()
    // ... test code
}
```

When `TF_ACC_CAPABILITIES` is set, tests silently return (not skip) if their required capabilities aren't available. This avoids polluting test metrics with skipped tests.

## Available Capabilities

### Infrastructure Capabilities

| Capability | Env Value | Description |
|------------|-----------|-------------|
| `VMware` | `vmware` | VMware/vSphere cloud |
| `AWS` | `aws` | Amazon Web Services |
| `Azure` | `azure` | Microsoft Azure |
| `GCP` | `gcp` | Google Cloud Platform |
| `OpenStack` | `openstack` | OpenStack cloud |
| `Hyperv` | `hyperv` | Microsoft Hyper-V |

### Network Integrations

| Capability | Env Value | Description |
|------------|-----------|-------------|
| `NSXT` | `nsxt` | VMware NSX-T |
| `NSXV` | `nsxv` | VMware NSX-V |
| `ACI` | `aci` | Cisco ACI |

### Network Features

| Capability | Env Value | Description |
|------------|-----------|-------------|
| `NetworkDHCP` | `network_dhcp` | DHCP server/relay |
| `NetworkRouter` | `network_router` | Network routers |
| `NetworkFirewall` | `network_firewall` | Firewall rules |
| `NetworkLoadBalancer` | `network_loadbalancer` | Load balancers |

### Automation Integrations

| Capability | Env Value | Description |
|------------|-----------|-------------|
| `Ansible` | `ansible` | Ansible integration |
| `AnsibleTower` | `ansible_tower` | Ansible Tower/AWX |
| `Chef` | `chef` | Chef integration |
| `Puppet` | `puppet` | Puppet integration |

### Container/Orchestration

| Capability | Env Value | Description |
|------------|-----------|-------------|
| `Kubernetes` | `kubernetes` | Kubernetes clusters |
| `Docker` | `docker` | Docker/container registries |

### Storage & VDI

| Capability | Env Value | Description |
|------------|-----------|-------------|
| `Alletra` | `alletra` | HPE Alletra storage |
| `VDI` | `vdi` | VDI pools/apps/gateways |

### Special Capabilities

| Capability | Env Value | Description |
|------------|-----------|-------------|
| `All` | `all` | Tests that create significant infrastructure load |

## The All Capability

Tests are marked as "all" if they don't target VMware, Alletra, NSXT, or NSXV. This allows excluding resource-intensive tests when running on shared hardware.

**Quiet tests** (no `All` capability):
- VMware cloud tests
- Alletra storage tests
- NSX-T network tests
- NSX-V network tests

**All tests** (have `All` capability):
- AWS, Azure, GCP tests
- Ansible, Kubernetes tests
- Generic network tests
- All other infrastructure tests

### Example: Running Only Quiet Tests

```bash
# Only run tests safe for shared hardware
export TF_ACC_CAPABILITIES="vmware,alletra,nsxt,nsxv"
go test ./morpheus/...
```

### Example: Running Everything Including All

```bash
# Include all tests
export TF_ACC_CAPABILITIES="vmware,alletra,nsxt,nsxv,all,aws,kubernetes,ansible"
go test ./morpheus/...
```

## Writing Tests

### Basic Pattern

```go
func TestAccMorpheusResource(t *testing.T) {
    // Capability check MUST be first statement
    if capabilities.Missing(t, capabilities.All) {
        return
    }
    t.Parallel()
    defer testhelpers.RecordResult(t)
    
    // ... test implementation
}
```

### Multiple Capabilities

```go
func TestAccMorpheusNSXTLoadBalancer(t *testing.T) {
    // Test requires both NSXT and NetworkLoadBalancer
    if capabilities.Missing(t, capabilities.NSXT, capabilities.NetworkLoadBalancer) {
        return
    }
    t.Parallel()
    // ...
}
```

### Alternative Capabilities (Any)

```go
func TestAccMorpheusNetworkFeature(t *testing.T) {
    // Test can run on either NSXT or NSXV
    if capabilities.MissingAll(t, capabilities.NSXT, capabilities.NSXV) {
        return
    }
    t.Parallel()
    // ...
}
```

## Tools

### capcheck - Migration Tool

Analyzes tests and adds capability checks:

```bash
# Analyze all tests (dry run)
go run ./cmd/capcheck ./morpheus/...

# Show what would change
go run ./cmd/capcheck --dry-run ./morpheus/...

# Apply changes
go run ./cmd/capcheck --apply ./morpheus/...

# Export capability mapping to JSON
go run ./cmd/capcheck --output=mapping.json ./morpheus/...

# Verbose output
go run ./cmd/capcheck --verbose ./morpheus/...
```

### caplint - Linter

Enforces capability declarations in CI:

```bash
# Lint all tests
go run ./cmd/caplint ./morpheus/...

# Verbose mode (show passing tests)
go run ./cmd/caplint --verbose ./morpheus/...

# Quiet mode (summary only)
go run ./cmd/caplint --quiet ./morpheus/...
```

Exit codes:
- `0` = no violations
- `1` = violations found
- `2` = error processing files

## Environment Variables

| Variable | Description |
|----------|-------------|
| `TF_ACC_CAPABILITIES` | Comma-separated list of available capabilities |
| `TF_ACC_CAPABILITIES_VERBOSE` | Set to `1` to log when tests don't run |

### Examples

```bash
# Shared hardware (quiet tests only)
export TF_ACC_CAPABILITIES="vmware,alletra,nsxt,nsxv"

# Full test suite
export TF_ACC_CAPABILITIES="vmware,alletra,nsxt,nsxv,all,aws,azure,gcp,kubernetes,ansible"

# Debug: see which tests are being skipped
export TF_ACC_CAPABILITIES_VERBOSE=1
```

## API Reference

### Functions

```go
// Has returns true if a single capability is available
func Has(cap Capability) bool

// HasAll returns true if ALL capabilities are available
func HasAll(caps ...Capability) bool

// HasAny returns true if ANY capability is available
func HasAny(caps ...Capability) bool

// Missing returns true if ANY capability is missing (use in tests)
func Missing(t *testing.T, caps ...Capability) bool

// MissingAll returns true if ALL capabilities are missing (use for "any of" checks)
func MissingAll(t *testing.T, caps ...Capability) bool

// IsVerbose returns true if verbose logging is enabled
func IsVerbose() bool
```

### Usage in Tests

```go
// Require all capabilities (AND)
if capabilities.Missing(t, capabilities.NSXT, capabilities.NetworkRouter) {
    return  // Missing NSXT OR NetworkRouter
}

// Require any capability (OR) 
if capabilities.MissingAll(t, capabilities.NSXT, capabilities.NSXV) {
    return  // Missing NSXT AND NSXV (neither available)
}
```

## Statistics

Current test distribution:

| Metric | Count |
|--------|-------|
| Total tests | 540 |
| Quiet tests (VMware/Alletra/NSX*) | 50 |
| All tests | 490 |

Capability breakdown:

| Capability | Tests |
|------------|-------|
| all | 490 |
| network_router | 32 |
| alletra | 29 |
| network_firewall | 23 |
| network_loadbalancer | 18 |
| nsxt | 17 |
| network_dhcp | 16 |
| aws | 12 |
| vdi | 12 |
| kubernetes | 11 |
| ansible | 5 |
| ansible_tower | 4 |
| vmware | 4 |
| docker | 3 |
| chef | 2 |
| azure | 1 |
| gcp | 1 |
| nsxv | 1 |
| puppet | 1 |

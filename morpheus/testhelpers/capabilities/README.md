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
export TF_ACC_CAPABILITIES="all,vmware,alletra,nsxt,nsxv,hvm,affinity_group,aws,azure,gcp,kubernetes,ansible"
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
| `HVM` | `hvm` | HPE VM (`mvm-cluster`) clusters/hosts |
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
| `Network` | `network` | Basic network create/update on a network-service-capable cloud |
| `NetworkDHCP` | `network_dhcp` | DHCP server/relay |
| `NetworkPool` | `network_pool` | Network (IP) pools |
| `NetworkServer` | `network_server` | Network server / service integration (e.g. NSX-T) |
| `NetworkRouter` | `network_router` | Network routers |
| `NetworkFirewall` | `network_firewall` | Firewall rules |
| `NetworkLoadBalancer` | `network_loadbalancer` | Load balancers (general API: NSX-T, config validation, data sources) |
| `NetworkLoadBalancerHAProxy` | `network_loadbalancer_haproxy` | HAProxy container LB provisioning (needs the `load-balancer-haproxy-1.7` layout + a cloud that can provision the container) |
| `Subnet` | `subnet` | Subnets |

### Compute Features

| Capability | Env Value | Description |
|------------|-----------|-------------|
| `AffinityGroup` | `affinity_group` | Cloud/cluster affinity groups. Needs Morpheus >= 8.0.10 **and** a supporting type: a `vmware`, `vmwareCloudAws` or `macstadium` cloud, or an `mvm-cluster` (HVM) cluster. Always paired with `vmware` or `hvm`; the ID of the cloud/cluster to use comes from `TF_VAR_testacc_morpheus_affinity_cloud_id` / `TF_VAR_testacc_morpheus_affinity_cluster_id`. Cloud affinity groups additionally need `TF_VAR_testacc_morpheus_affinity_pool_id` — a resource pool of type Cluster on that cloud — because Morpheus rejects a create without a pool |

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
| `Kubernetes` | `kubernetes` | Kubernetes library artifacts (spec templates, blueprints, cluster layouts/types) -- no running cluster required |
| `KubernetesCluster` | `kubernetes_cluster` | A live, healthy Kubernetes cluster (HKS worker provisioning + cluster namespaces). Omit on environments that lack a usable cluster |
| `Docker` | `docker` | Docker/container registries |

### Storage & VDI

| Capability | Env Value | Description |
|------------|-----------|-------------|
| `Alletra` | `alletra` | HPE Alletra storage |
| `VDI` | `vdi` | VDI pools/apps/gateways |

### Identity / Platform

| Capability | Env Value | Description |
|------------|-----------|-------------|
| `PCE` | `pce` | A live PCE (Private Cloud Enterprise) instance reachable via GreenLake IAM. Gates the end-to-end `pce_identity` auth-flow test, which authenticates through GreenLake rather than the main appliance credentials. The test skips unless `TF_VAR_testacc_pce_identity_client_id`, `_client_secret`, `_issuer_url`, `_location`, `_space` and `TF_ACC_PCE_IDENTITY_CLOUD_NAME` are all set |

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
- Cloud/cluster affinity group tests -- additionally gated on `affinity_group`, so the shared-hardware example below does **not** run them

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
export TF_ACC_CAPABILITIES="vmware,alletra,nsxt,nsxv,all,hvm,affinity_group,aws,kubernetes,ansible"
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
| `TF_VAR_testacc_morpheus_affinity_cloud_id` | ID of a cloud that supports affinity groups. Required by the `affinity_group` + `vmware` tests, which skip without it |
| `TF_VAR_testacc_morpheus_affinity_cluster_id` | ID of an HVM cluster that supports affinity groups. Required by the `affinity_group` + `hvm` tests, which skip without it |
| `TF_VAR_testacc_morpheus_affinity_pool_id` | ID of a resource pool of type Cluster on the cloud above. Morpheus rejects a cloud affinity group created without a pool, so the `affinity_group` + `vmware` resource tests skip without it |
| `TF_VAR_testacc_pce_identity_client_id` | GreenLake API client ID. Required by the `pce` test |
| `TF_VAR_testacc_pce_identity_client_secret` | GreenLake API client secret. Required by the `pce` test |
| `TF_VAR_testacc_pce_identity_issuer_url` | GreenLake IAM issuer URL used to mint access tokens. Required by the `pce` test |
| `TF_VAR_testacc_pce_identity_location` | PCE instance location. Required by the `pce` test — the `pce_identity` block rejects a null `location` |
| `TF_VAR_testacc_pce_identity_space` | GreenLake space containing the PCE instance. Required by the `pce` test — the `pce_identity` block rejects a null `space` |
| `TF_VAR_testacc_pce_identity_broker_url` | PCE broker URL override. Optional; leave unset to use the HPE-hosted cloud broker |
| `TF_ACC_PCE_IDENTITY_CLOUD_NAME` | Name of a cloud to read back through the PCE token exchange. Required by the `pce` test; there is no default, as the reachable clouds differ per GreenLake tenant. Not a `TF_VAR_`, because the test renders it into HCL rather than passing it through a Terraform variable |

### Examples

```bash
# Shared hardware (quiet tests only)
export TF_ACC_CAPABILITIES="vmware,alletra,nsxt,nsxv"

# Full test suite
export TF_ACC_CAPABILITIES="vmware,alletra,nsxt,nsxv,all,hvm,affinity_group,aws,azure,gcp,kubernetes,ansible"

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

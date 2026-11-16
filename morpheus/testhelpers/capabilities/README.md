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
export TF_ACC_CAPABILITIES="all,vmware,alletra,nsxt,nsxv,hvm,affinity_group,aws,azure,gcp,kubernetes,ansible,master_tenant"
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
| `NetworkServer` | `network_server` | A network server (network-service integration: SDN, NSX-T, an HVM cluster's own) on the appliance, named by the `TF_VAR_testacc_morpheus_*network_server_name` variables below (`testhelpers/network_server_fixture.go`). The typed HVM network tests also need the HVM fixture and skip until the HVM cluster's hidden network server can be resolved by name (MORPH-20056) |
| `NetworkRouter` | `network_router` | Network routers |
| `NetworkFirewall` | `network_firewall` | Firewall rules |
| `NetworkLoadBalancer` | `network_loadbalancer` | Load balancers (general API: NSX-T, config validation, data sources) |
| `NetworkLoadBalancerHAProxy` | `network_loadbalancer_haproxy` | HAProxy container LB provisioning (needs the `load-balancer-haproxy-1.7` layout + a cloud that can provision the container) |
| `Subnet` | `subnet` | Subnets |

### Compute Features

| Capability | Env Value | Description |
|------------|-----------|-------------|
| `AffinityGroup` | `affinity_group` | Cloud/cluster affinity groups. Needs Morpheus >= 8.0.10 **and** a supporting type: a `vmware`, `vmwareCloudAws` or `macstadium` cloud, or an `mvm-cluster` (HVM) cluster. Always paired with `vmware` or `hvm`; the ID of the cloud/cluster to use comes from `TF_VAR_testacc_morpheus_affinity_cloud_id` / `TF_VAR_testacc_morpheus_affinity_cluster_id`. Cloud affinity groups additionally need `TF_VAR_testacc_morpheus_affinity_pool_id` — a resource pool of type Cluster on that cloud — because Morpheus rejects a create without a pool |
| `HostProfile` | `host_profile` | Host profile API (config documents attachable to clusters), introduced in Morpheus 9.1.0 alongside HVM 2.0. Needs Morpheus >= 9.1.0. Does **not** require a live HVM cluster — the CRUD API is in core and compliance enforcement (the HVM/MVM-specific part) is not exercised by these tests. Not paired with `hvm` |

### Automation Integrations

| Capability | Env Value | Description |
|------------|-----------|-------------|
| `Ansible` | `ansible` | Ansible integration |
| `AnsibleTower` | `ansible_tower` | Ansible Tower/AWX, reachable from the appliance, with an integration for it on the appliance -- both named by the `TF_VAR_testacc_morpheus_ansible_tower_*` variables below (`testhelpers/tower_fixture.go`) |
| `Chef` | `chef` | Chef integration |
| `Git` | `git` | Git integration. The tests that create a Git or Ansible integration also need a repository the appliance can clone, from the `TF_VAR_testacc_morpheus_git_repository_*` variables below, and skip without one |
| `Puppet` | `puppet` | Puppet integration |

### Container/Orchestration

| Capability | Env Value | Description |
|------------|-----------|-------------|
| `Kubernetes` | `kubernetes` | Kubernetes library artifacts (spec templates, blueprints, cluster layouts/types) -- no running cluster required |
| `KubernetesCluster` | `kubernetes_cluster` | A live, healthy Kubernetes cluster that `TF_VAR_testacc_morpheus_k8s_cluster_id` names (cluster namespaces). The HKS-on-HVM provisioning test is gated on it too but is opt-in (`RUN_SKIPPED_BY_DEFAULT`): it builds four nodes from QA-appliance ids. Omit on environments that lack a usable cluster |
| `Docker` | `docker` | Docker/container registries |

### Storage & VDI

| Capability | Env Value | Description |
|------------|-----------|-------------|
| `Alletra` | `alletra` | HPE Alletra storage |
| `Backup` | `backup` | Backup jobs, instance backups and host (file) backups, created and destroyed against the objects the `TF_VAR_testacc_morpheus_backup_*` variables name (`testhelpers/backup_fixture.go`). Nothing is scheduled, so no backup runs; backups have to be enabled on the appliance |
| `Metal` | `metal` | HPE Bare Metal cloud (iLO/Redfish servers). Gates the bare-metal instance-node tests and, paired with `alletra`, the HPE Alletra MP Bare Metal (BMaaS) datastore test, which reads its cloud, storage server, resource pool, group and tenant from the `TF_VAR_testacc_morpheus_metal_*` fixture below and skips when one is unset |
| `VDI` | `vdi` | VDI pools/apps/gateways |

### Identity / Platform

| Capability | Env Value | Description |
|------------|-----------|-------------|
| `PCE` | `pce` | A live PCE (Private Cloud Enterprise) instance reachable via GreenLake IAM. Gates the end-to-end `pce_identity` auth-flow test, which authenticates through GreenLake rather than the main appliance credentials. The test skips unless `TF_VAR_testacc_pce_identity_client_id`, `_client_secret`, `_issuer_url`, `_location`, `_space` and `TF_ACC_PCE_IDENTITY_CLOUD_NAME` are all set |
| `LDAP` | `ldap` | An Active Directory / LDAP directory the appliance can bind to. Gates the Active Directory identity source test, which reads the directory from the `TF_VAR_testacc_active_directory_*` variables below and skips without them. Morpheus binds to the directory and looks its groups up when the identity source is saved, so the test cannot run against a placeholder |
| `Plugins` | `plugins` | Plugin install/remove -- appliance-wide, needs admin-plugins permission |

### Tenancy

| Capability | Env Value | Description |
|------------|-----------|-------------|
| `MasterTenant` | `master_tenant` | Master-tenant-only operations that no grantable role permission can enable for a subtenant (e.g. OsType creation, nominating a tenant's `parent_id`). Include it when the test credentials belong to the master tenant; omit it for subtenant runs so these tests skip instead of failing |

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
export TF_ACC_CAPABILITIES="vmware,alletra,nsxt,nsxv,all,hvm,affinity_group,aws,kubernetes,ansible,master_tenant"
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
| `TF_ACC_PHPIPAM_URL`, `TF_ACC_PHPIPAM_APP_ID`, `TF_ACC_PHPIPAM_USERNAME`, `TF_ACC_PHPIPAM_PASSWORD` | A reachable phpIPAM service for the phpIPAM pool-server test, which skips without them. Optional `TF_ACC_PHPIPAM_APP_ID_2` names a second application for the app-ID change step |
| `TF_VAR_testacc_morpheus_affinity_cloud_id` | ID of a cloud that supports affinity groups. Required by the `affinity_group` + `vmware` tests, which skip without it |
| `TF_VAR_testacc_morpheus_affinity_cluster_id` | ID of an HVM cluster that supports affinity groups. Required by the cluster affinity-group resource/data-source tests, which skip without it (the instance placement tests use the HVM fixture's cluster by name instead) |
| `TF_VAR_testacc_morpheus_affinity_pool_id` | ID of a resource pool of type Cluster on the cloud above. Morpheus rejects a cloud affinity group created without a pool, so the `affinity_group` + `vmware` resource tests skip without it |
| `TF_VAR_testacc_morpheus_backup_instance_name` | Name of an instance the `backup` tests attach backup records to (`testhelpers/backup_fixture.go`); resolved with the `hpe_morpheus_instance` data source. Its cloud has to offer a backup type for it (a KVM/HVM instance: `kvm`) |
| `TF_VAR_testacc_morpheus_backup_server_name` | Name of a managed compute server for the host backup tests, which back up a file path on it; resolved with `hpe_morpheus_compute_server` |
| `TF_VAR_testacc_morpheus_backup_bucket_name`, `TF_VAR_testacc_morpheus_backup_other_bucket_name` | Two storage buckets a backup can target; the update tests move a backup from the first to the other. Resolved with `hpe_morpheus_storage_bucket` |
| `TF_VAR_testacc_morpheus_hvm_cluster_id` | ID of an existing HVM cluster for the cluster host-profile attach/detach test (`host_profile`) and the cluster-branch compliance data-source test (`TestAccMorpheusHostProfileComplianceDataSourceCluster`), both of which skip without it. The attach/detach test attaches an enabled, non-default host profile to the cluster and detaches it again — a pure database operation that does not run compliance checks or remediation against the cluster's hosts. It additionally skips if the named cluster already has a profile attached |
| `TF_VAR_testacc_morpheus_compute_server_id` | ID of an existing HVM/MVM compute server (a cluster host) for the server-branch host-profile compliance data-source test (`TestAccMorpheusHostProfileComplianceDataSourceServer`), which skips without it. The server lookup returns a single boolean `compliant` with a null `overall_status` (the per-cluster rollup is not computed for a single host) |
| `TF_VAR_testacc_morpheus_hvm_cloud_name` | Name of the Morpheus private cloud an HVM cluster is in. The HVM fixture (`testhelpers/hvm_fixture.go`): the objects the HVM/KVM instance and data-source tests provision against, **by name**, resolved by data sources inside the test configuration — never by id. Tests skip when a variable they need is unset |
| `TF_VAR_testacc_morpheus_hvm_cluster_name` | Name of the HVM cluster to provision into; its resource pool, networks and datastores are what the other names are looked up in |
| `TF_VAR_testacc_morpheus_hvm_network_name` | Name of a network of that cluster (`GET /api/networks?clusterId=`) |
| `TF_VAR_testacc_morpheus_hvm_datastore_name` | Name of a datastore of that cluster (`GET /api/clusters/{id}/datastores`); optional for the instance tests, which fall back to automatic datastore selection |
| `TF_VAR_testacc_morpheus_hvm_layout_name`, `TF_VAR_testacc_morpheus_hvm_layout_version` | The instance type layout to provision with (its instance type comes with it); the version pins one of the same-named layouts |
| `TF_VAR_testacc_morpheus_hvm_plan_name` | A small kvm service plan for the test VMs, e.g. `1 CPU, 1GB Memory` |
| `TF_VAR_testacc_morpheus_hvm_image_name` | A virtual image to provision from when the layout carries none; optional for the instance tests |
| `TF_VAR_testacc_morpheus_hvm_affinity_group_name` | An existing cluster affinity group of the HVM cluster, for the instance affinity-group tests (replaces the former `TF_VAR_testacc_morpheus_affinity_group_cluster_ag_id`); it has to exist, since 9.1 refuses an empty group |
| `TF_VAR_testacc_morpheus_hvm_growing_image_name` | A virtual image whose raw size is just over its minimum disk (10 GB requested -> 11 GB after provisioning), for the volume `actual_size` test; it skips without one |
| `TF_VAR_testacc_morpheus_hvm_no_agent` | `true` to provision the test VMs without the Morpheus agent (default: with it). Without the agent Morpheus SSHes from the appliance into the guest to report its address and apply environment variables, so set this only where the appliance can reach the test VMs |
| `TF_VAR_testacc_morpheus_hvm_other_cluster_name` | A second HVM cluster in the same cloud, for the wrong-pool test (an instance targeting one cluster's pool with a network of the other's) |
| `TF_VAR_testacc_morpheus_hvm_poolless_network_name` | A network of the HVM cloud that belongs to no resource pool (a docker bridge, say), for the test that lists the networks Morpheus would accept instead |
| `TF_VAR_testacc_morpheus_ansible_tower_integration_name` | Name of the appliance's Ansible Tower integration, for the inventory, job template and Tower task tests (`testhelpers/tower_fixture.go`); resolved with the `hpe_morpheus_integration` data source |
| `TF_VAR_testacc_morpheus_ansible_tower_inventory_name`, `TF_VAR_testacc_morpheus_ansible_tower_job_template_name` | An inventory and a job template that integration has synced from the Tower, as `GET /api/options/ansibleTowerInventory` and `ansibleTowerJobTemplate` list them |
| `TF_VAR_testacc_morpheus_ansible_tower_url`, `TF_VAR_testacc_morpheus_ansible_tower_username`, `TF_VAR_testacc_morpheus_ansible_tower_password` | A Tower (or a stand-in serving its API) and an account on it, for the integration resource test: Morpheus verifies the URL on create (TCP, `/api/v2/ping/`, `/api/v2/me/`) and refuses the save otherwise |
| `TF_VAR_testacc_morpheus_metal_cloud_name` | Name of an HPE Bare Metal cloud. The bare-metal fixture (`testhelpers/metal_fixture.go`) for the `metal` + `alletra` BMaaS datastore test: like the HVM fixture, everything is resolved **by name** by data sources inside the test configuration -- never by id. The test skips when a variable it needs is unset |
| `TF_VAR_testacc_morpheus_metal_storage_server_name` | Name of an HPE Alletra MP storage server (type `alletra-mp`) reachable from that cloud. The storage plugin refuses a second *online* BMaaS datastore for the same cloud, storage server and protocol, so the cloud must not already have an iSCSI BMaaS datastore on this server |
| `TF_VAR_testacc_morpheus_metal_pool_name` | Name of a resource pool of that cloud, as listed by `GET /api/clouds/<cloud>/resource-pools` -- `Default`, the pool an HPE Bare Metal cloud creates on initialisation, unless another was added. Resolved through the `zonePools` option source (`hpe_morpheus_option_source_value`), since there is no resource-pool data source. BMaaS datastores are created in a pool of their cloud; an id that names no pool is rejected by Morpheus as `zonePool.id: Required fields are missing` |
| `TF_VAR_testacc_morpheus_metal_group_name` | Name of a group to grant access to the test datastore |
| `TF_VAR_testacc_morpheus_metal_tenant_name` | Name of a tenant to share the test datastore with |
| `TF_VAR_testacc_morpheus_git_integration_name` | Name of the appliance's Git integration, for the repository-sourced blueprint tests (`testhelpers/git_fixture.go`); resolved with the `hpe_morpheus_integration_git` data source, whose `repository_ids` map gives the repository's id |
| `TF_VAR_testacc_morpheus_git_repository_name` | Name of a code repository Morpheus recorded from that integration (the repository's own name, as the `repository_ids` map keys it) |
| `TF_VAR_testacc_morpheus_git_repository_url` | URL of a Git repository Morpheus can clone without credentials, for the tests that create a Git or Ansible integration (`sdkv2/resources/integration`): Morpheus clones the repository when the integration is saved and rejects one it cannot reach, so the documented example's public repository cannot be relied on from every appliance. The tests skip when it is unset |
| `TF_VAR_testacc_morpheus_git_repository_branch` | A branch of that repository, used two ways: the blueprints' `version_ref`, and the default branch of the Git and Ansible integrations the integration tests create. Morpheus clones the repository at that branch and parses the files under the working path when a Terraform or ARM blueprint is saved, so the branch has to exist and `./test` has to hold `main.tf` / `azuredeploy.json`. The tests skip when it is unset |
| `TF_VAR_testacc_pce_identity_client_id` | GreenLake API client ID. Required by the `pce` test |
| `TF_VAR_testacc_pce_identity_client_secret` | GreenLake API client secret. Required by the `pce` test |
| `TF_VAR_testacc_pce_identity_issuer_url` | GreenLake IAM issuer URL used to mint access tokens. Required by the `pce` test |
| `TF_VAR_testacc_pce_identity_location` | PCE instance location. Required by the `pce` test — the `pce_identity` block rejects a null `location` |
| `TF_VAR_testacc_pce_identity_space` | GreenLake space containing the PCE instance. Required by the `pce` test — the `pce_identity` block rejects a null `space` |
| `TF_VAR_testacc_pce_identity_broker_url` | PCE broker URL override. Optional; leave unset to use the HPE-hosted cloud broker |
| `TF_ACC_PCE_IDENTITY_CLOUD_NAME` | Name of a cloud to read back through the PCE token exchange. Required by the `pce` test; there is no default, as the reachable clouds differ per GreenLake tenant. Not a `TF_VAR_`, because the test renders it into HCL rather than passing it through a Terraform variable |
| `TF_VAR_testacc_morpheus_network_server_name` | Name of a network server `GET /api/networks/servers` lists (an SDN or NSX-T integration), for the network server data-source tests (`testhelpers/network_server_fixture.go`) |
| `TF_VAR_testacc_morpheus_hvm_network_server_name` | Name of the HVM cluster's network server (`<cluster> HVM Network Server`), on which the typed HVM network tests create their networks. The list API returns it only with `showHidden=true`, which the data source cannot ask for until MORPH-20056; the tests skip while this is unset |
| `TF_VAR_testacc_morpheus_hvm_switch_name`, `TF_VAR_testacc_morpheus_hvm_switch_name_2` | Two virtual switches of the HVM cluster as the `hvm-vswitch-or-router` option source lists them (`vs0`, `vs1`); the OVS port group test creates on the first and moves to the second, a replacement |
| `TF_VAR_testacc_morpheus_hvm_host_interface` | A host network interface the HVM hosts have (`bond0`, `hvmnet2`, ...), for the HVM data network test |
| `TF_VAR_testacc_active_directory_server` | Host name of an Active Directory domain controller the appliance can bind to, for the Active Directory identity source test (`sdkv2/resources/identitysource`; needs the `ldap` capability). Morpheus binds to it as the binding user when the identity source is saved and rejects the save if that fails. The test skips when any of the five `TF_VAR_testacc_active_directory_*` variables below is unset |
| `TF_VAR_testacc_active_directory_domain` | That directory's domain, e.g. `example.com` |
| `TF_VAR_testacc_active_directory_binding_username` | A user that can bind to the directory |
| `TF_VAR_testacc_active_directory_binding_password` | That user's password |
| `TF_VAR_testacc_active_directory_group_name` | A group in the directory to map a Morpheus role to, spelt as its `sAMAccountName`: Morpheus looks the group up when the identity source is saved and rewrites the mapping to the name and distinguished name the directory reports |
| `TF_VAR_testacc_active_directory_group_fqn` | Optional: that group's distinguished name. When set it is configured and asserted; when unset the test leaves the attribute to Morpheus and only asserts that it was filled in |
| `TF_ACC_IPAM_TYPE` | Pool server type code (e.g. `infoblox`, `phpipam`, `bluecat`) of a reachable IPAM. Required by the network pool server create-based tests — Morpheus 9.1.2 verifies the service URL on create, so these tests skip without a reachable IPAM |
| `TF_ACC_IPAM_URL` | Reachable IPAM service URL. Required by the network pool server create-based tests (9.1.2 reachability check) |
| `TF_ACC_IPAM_USERNAME` | IPAM service username. Required by the network pool server create-based tests (9.1.2 reachability check) |
| `TF_ACC_IPAM_PASSWORD` | IPAM service password. Required by the network pool server create-based tests (9.1.2 reachability check) |
| `TF_ACC_STORAGE_SERVER_TYPE` | Storage server type code (e.g. `alletra-mp`) of a reachable storage service. Required by the storage server create-based tests — Morpheus 9.1.2 verifies the service URL on create, so these tests skip without a reachable service |
| `TF_ACC_STORAGE_SERVER_URL` | Reachable storage service URL. Required by the storage server create-based tests (9.1.2 reachability check) |
| `TF_ACC_STORAGE_SERVER_USERNAME` | Storage service username. Required by the storage server create-based tests (9.1.2 reachability check) |
| `TF_ACC_STORAGE_SERVER_PASSWORD` | Storage service password. Required by the storage server create-based tests (9.1.2 reachability check) |
| `TF_ACC_IPAM_ROTATABLE` | Set to `1` to re-enable the network pool server write-only matrix and drift tests. They are skipped by default because Morpheus 9.1.x validates the IPAM credentials on create and update, so an out-of-band password rotation cannot be staged against an ordinary integration; set this against an IPAM whose credentials can be rotated (MORPH-19005 follow-up) |
| `TF_ACC_STORAGE_SERVER_ROTATABLE` | Set to `1` to re-enable the storage server write-only matrix and drift tests. They are skipped by default because Morpheus 9.1.x validates the storage server credentials on create and update, so an out-of-band password rotation cannot be staged; set this against a storage server whose credentials can be rotated (MORPH-19005 follow-up) |
| `TF_ACC_BGP_SOURCE_ADDRESS` | Overrides the BGP source address used by the network router BGP neighbor drift test (`config_nsxt.source_addresses`). Optional; no fixture carries it, so set it to a valid source address on the target tier-0 router when running that test |

### Examples

```bash
# Shared hardware (quiet tests only)
export TF_ACC_CAPABILITIES="vmware,alletra,nsxt,nsxv"

# Full test suite
export TF_ACC_CAPABILITIES="vmware,alletra,nsxt,nsxv,all,hvm,affinity_group,aws,azure,gcp,kubernetes,ansible,master_tenant"

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

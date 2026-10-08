---
page_title: "hpe_morpheus_cluster Resource - terraform-provider-hpe"
subcategory: "Morpheus"
description: |-
  
---
# hpe_morpheus_cluster (Resource)



The `hpe_morpheus_cluster` resource is used to provision and manage various cluster types in HPE Morpheus.

-> Currently HVM clusters are supported.  We have static `config` schemas for the following:<br>
&nbsp;&nbsp;&nbsp;&nbsp;- HVM: `config_hvm`<br>

## HVM Cluster

The `hvm_config` static config block can be used to provision a HVM cluster:

```terraform
resource "hpe_morpheus_cluster" "example_hvm" {
  name        = "TestCluster"
  description = "A test HVM cluster"
  cloud_id    = 1
  group_id    = 1
  layout_id   = 2

  labels = [
    "terraform",
    "example",
  ]

  config_hvm = {
    create_user       = false
    dynamic_placement = false
    cpu_arch          = "x86_64"
    cpu_model         = "host-model"
    power_policy      = "balanced"

    # HVM 2.0 network configuration. Only honoured when the selected cluster
    # layout has clusterVersion 2.0 or greater; otherwise Morpheus silently
    # drops these fields.
    virt_switch_name = "vSwitch0"
    dns1             = "8.8.8.8"
    dns2             = "8.8.4.4"
    dns3             = "1.1.1.1"
    search_domain    = "example.com"
    ntp1             = "0.pool.ntp.org"
    ntp2             = "1.pool.ntp.org"
    proxy            = "http://proxy.example.com:3128"
    no_proxy         = "localhost,127.0.0.1,.example.com"
  }

  server = {
    service_plan_id = 1

    ssh_port                 = 22
    ssh_username             = "user"
    ssh_key_pair_id          = 1
    management_net_interface = "eth0"

    ssh_hosts = [
      {
        name = "host1"
        ip   = "10.0.0.1"
      },
      {
        name = "host2"
        ip   = "10.0.0.2"
      },
      {
        name = "host3"
        ip   = "10.0.0.3"
      }
    ]

    visibility = "private"

    tags = [
      {
        name  = "source"
        value = "terraform"
      },
      {
        name  = "environment"
        value = "example"
      },
    ]
  }
}
```

!> **Set `config_hvm.overlay_interface_name` on HVM 1.x layouts.** It names the host network
interface used for the overlay (VXLAN) Open vSwitch bridge, for example `ens193`. Morpheus treats it
as optional — the overlay bridge is created only when it is set, and Morpheus recommends it whenever
overlay networks will be used — and HVM 2.0 (virtual switch) layouts do not use it at all. However,
on Morpheus 9.0.x with an HVM 1.x layout, omitting it makes every host fail to build at the
"HVM Host Build OVS Interface" step (`build-ovs.py` exits with `KeyError: None`), so set it until the
appliance fix is available. The provider does not make it required, because the API does not.

~> **Terraform 1.14 caveat.** Terraform 1.14 silently ignores unknown attribute names nested inside a single-nested attribute such as `config_hvm`. A misspelled field name (for example `dns_1` instead of `dns1`) passes `terraform validate` and `terraform plan` without any error or warning on affected CLI versions, so the intended value is quietly dropped and the cluster is built without it. Double-check the spelling of `config_hvm` fields against the schema below. (See also the provider README.)

## Generic Cluster Config using the Dynamic Config

The cluster `config` block can be used to pass an arbitrary cluster server configuration to Morpheus. The values in the provider `config` block map to `server.config` in the Morpheus API payload. The values in `config` are interpreted literally - they can be thought of like a JSON blob - so if you wish to provide certain configuration options, they must be in the exact format that the API expects. For example, `create_user` in the static block would become `createUser` if provided using the dynamic `config`.

When using the dynamic `config` block, the `cluster_type_code` attribute must be provided.

-> It is only recommended to use the dynamic `config` block if you need to provision a cluster of a certain type that the provider does not yet support via static `config` blocks. Not all functionality that a static `config` block would provide is guaranteed when using the dynamic `config` block. If a static `config` type exists for the type of cluster you wish to provision, it is recommended that you use the static `config` instead of the dynamic `config`.

Below is an example of provisioning a HVM cluster using the dynamic `config` block. Note the casing of the parameters in `config`:

```terraform
resource "hpe_morpheus_cluster" "example_generic_hvm" {
  name              = "TestCluster"
  description       = "A HVM cluster created with a dynamic config"
  cloud_id          = 1
  group_id          = 1
  layout_id         = 2
  cluster_type_code = "mvm-cluster"

  labels = [
    "terraform",
    "example",
  ]

  config = {
    createUser           = false
    cpuArch              = "x86_64"
    cpuModel             = "host-model"
    dynamicPlacementMode = "off"
    powerPolicy          = "balanced"
  }

  server = {
    service_plan_id = 1

    ssh_port                 = 22
    ssh_username             = "user"
    ssh_key_pair_id          = 1
    management_net_interface = "eth0"

    ssh_hosts = [
      {
        name = "host1"
        ip   = "10.0.0.1"
      },
      {
        name = "host2"
        ip   = "10.0.0.2"
      },
      {
        name = "host3"
        ip   = "10.0.0.3"
      }
    ]

    visibility = "private"

    tags = [
      {
        name  = "source"
        value = "terraform"
      },
      {
        name  = "environment"
        value = "example"
      },
    ]
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cloud_id` (Number) The Cloud ID to provision the host into
- `group_id` (Number) The Group ID to provision the cluster into
- `layout_id` (Number) The Layout ID for the host type(s) that will be provisioned for the cluster
- `name` (String) Name of the cluster to be created
- `server` (Attributes) Configuration options for the cluster's workers. (see [below for nested schema](#nestedatt--server))

### Optional

- `cluster_type_code` (String) Code of type of cluster to be created. Automatically set to the appropriate code by the provider when using a static config.
- `config` (Dynamic) Generic Cluster Configuration
- `config_hvm` (Attributes) Configuration for HVM cluster servers (see [below for nested schema](#nestedatt--config_hvm))
- `description` (String) Description of the cluster to be created
- `labels` (Set of String) Array of strings (keywords).
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `workflow_id` (Number) Optional Workflow Id desired to be run

### Read-Only

- `cpu_placement_mode` (String) Cluster CPU placement mode
- `id` (Number) The ID of the cluster
- `service_url` (String)
- `uuid` (String)
- `write_only_revision` (Number) The time, in seconds since the Unix epoch, at which a write-only value (server.ssh_password_wo) was last
sent to Morpheus. Changes whenever the cluster is created; not configurable.

<a id="nestedatt--server"></a>
### Nested Schema for `server`

Optional:

- `data_device` (String)
- `hostname` (String) Hostname for server host
- `lvm_enabled` (Boolean)
- `management_net_interface` (String) The primary management interface name to establish a management bridge (i.e. eth0,ens192,bond0,etc)
- `name` (String) Name to be used for host(s) created in the cluster
- `network_domain` (String) Network domain
- `network_interfaces` (Attributes Set) The networkInterfaces parameter is for network configuration.

The Options API /api/options/zoneNetworkOptions can be used to see which options are available.

It should be passed as an array of Objects with the following attributes (see [below for nested schema](#nestedatt--server--network_interfaces))
- `security_groups` (Set of String) Key for security group configuration.
- `service_plan_id` (Number) ID of the cluster server service plan.
- `ssh_hosts` (Attributes Set) Array of Host IPs and Names. This is used in conjunction with sshUsername and sshPassword/sshKeyPair to add existing hosts such as with HPE VM clusters. (see [below for nested schema](#nestedatt--server--ssh_hosts))
- `ssh_key_pair_id` (Number) SSH Key Pair ID.
- `ssh_password_wo` (String, Sensitive, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) SSH password (Write Only)
- `ssh_password_wo_version` (Number, Deprecated) Deprecated. Changes to server.ssh_password_wo are detected automatically; this attribute is
no longer needed and will be removed in the 2.6 release. The SSH password cannot be changed
after the cluster is created, so a change replaces the cluster.
- `ssh_port` (Number) The port which the worker's SSH server is listening on.
- `ssh_username` (String) SSH Username
- `tags` (Attributes Set) Metadata tags, Array of objects having a name and value. (see [below for nested schema](#nestedatt--server--tags))
- `user_group_id` (Number) User Group ID for server host
- `visibility` (String) Visibility for server host
- `volumes` (Attributes Set) The (optional) volumes parameter is for LV configuration, can create additional LVs at provision It should be passed as an array of Objects (see [below for nested schema](#nestedatt--server--volumes))

<a id="nestedatt--server--network_interfaces"></a>
### Nested Schema for `server.network_interfaces`

Required:

- `network_id` (Number) The ID of the network.

Optional:

- `ip_mode` (String) The mode for determining ip address. Can be 'static', 'dhcp' or empty string.


<a id="nestedatt--server--ssh_hosts"></a>
### Nested Schema for `server.ssh_hosts`

Required:

- `ip` (String) Host IP
- `name` (String) Host Name


<a id="nestedatt--server--tags"></a>
### Nested Schema for `server.tags`

Optional:

- `name` (String)
- `value` (String)


<a id="nestedatt--server--volumes"></a>
### Nested Schema for `server.volumes`

Optional:

- `datastore_auto_selection` (String) Auto selection can be specified as auto or autoCluster (for clusters).
- `datastore_id` (Number) The ID of the specific datastore.
- `name` (String) Name/type of the LV being created
- `root_volume` (Boolean) If set to false then a non-root LV will be created
- `size` (Number) Size of the LV to be created in GBs

Default is from the service plan
- `size_id` (String) Can be used to select pre-existing LV choices from Morpheus
- `storage_type_id` (Number) Identifier for LV type

Read-Only:

- `id` (Number) The id for the LV configuration being created



<a id="nestedatt--config_hvm"></a>
### Nested Schema for `config_hvm`

Required:

- `cpu_arch` (String)
- `cpu_model` (String)
- `power_policy` (String)

Optional:

- `compute_interface_name` (String)
- `compute_vlans` (String)
- `create_user` (Boolean) Whether to create a user on the cluster's workers. The default is 'false'.
- `dns1` (String) Primary DNS server for the HVM cluster. Requires a cluster layout with clusterVersion 2.0 or greater.
- `dns2` (String) Secondary DNS server for the HVM cluster. Requires a cluster layout with clusterVersion 2.0 or greater.
- `dns3` (String) Tertiary DNS server for the HVM cluster. Requires a cluster layout with clusterVersion 2.0 or greater.
- `dynamic_placement` (Boolean) When enabled, Dynamic Placement will automatically balance VMs across cluster hosts based on resource utilization. When disabled, VMs will only migrate to a new host if they are pinned to a specific host or failed over and not running on the preferred host.
- `no_proxy` (String) Comma-separated list of hosts that bypass the proxy. Requires a cluster layout with clusterVersion 2.0 or greater.
- `ntp1` (String) Primary NTP server for the HVM cluster. Requires a cluster layout with clusterVersion 2.0 or greater.
- `ntp2` (String) Secondary NTP server for the HVM cluster. Requires a cluster layout with clusterVersion 2.0 or greater.
- `overlay_interface_name` (String) Name of the host network interface used for the overlay (VXLAN) Open vSwitch bridge, for example `ens193`. Optional: Morpheus creates the overlay bridge only when it is set, and recommends setting it whenever overlay networks will be used. Applies to HVM 1.x cluster layouts; HVM 2.0 (virtual switch) layouts do not use it. Note that on Morpheus 9.0.x with an HVM 1.x layout, omitting it makes every host fail to build at the "HVM Host Build OVS Interface" step (`build-ovs.py` exits with `KeyError: None`), so set it until the appliance fix is available.
- `proxy` (String) Proxy server for the HVM cluster. Requires a cluster layout with clusterVersion 2.0 or greater.
- `search_domain` (String) DNS search domain for the HVM cluster. Requires a cluster layout with clusterVersion 2.0 or greater.
- `storage_interface_name` (String)
- `vcpu_placement_mode` (String)
- `virt_switch_name` (String) Name of the HVM virtual switch. Requires a cluster layout with clusterVersion 2.0 or greater.


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Update limitations

~> Update support on `hpe_morpheus_cluster` is limited - only basic cluster information such as name, description, and labels can be updated. Additionally, certain `config` properties can be updated. Refer to the resource schema for full information.

~> The SSH password (`server.ssh_password_wo`) is set when the cluster is created and cannot be changed afterwards. A change to the configured value therefore plans a **replacement** of the cluster; the read-only `write_only_revision` records when it was last sent. The deprecated `server.ssh_password_wo_version` marker still works for the 2.2 release (changing it plans the same replacement; removing it does not) and will be removed in 2.6.

~> The HVM 2.0 network fields `dns1`, `dns2`, `dns3`, `search_domain`, `ntp1`, `ntp2`, `proxy`, and `no_proxy` can be changed on an existing cluster and are applied as an **in-place update** rather than a replacement. Morpheus reconciles the new values onto the cluster hosts through its daily node-settings drift-correction sync, so the hosts converge within a day. Two caveats follow: (1) convergence is **asynchronous** — `terraform apply` returns before the hosts are reconfigured; and (2) the reconciler only pushes non-empty values, so clearing a field (setting it back to empty) updates Terraform state but does not roll the hosts back on its own. By contrast, `virt_switch_name` has create/scale-out semantics (virtual switches are provisioned per host when a host is added), so a change to it still forces a **replacement**.

## Cluster layout upgrades

~> Changing `layout_id` forces the cluster to be **replaced** (destroyed and recreated). A layout reassignment in Morpheus only re-points the cluster's server group at the new layout; it does not reconfigure, re-provision, or roll the existing hosts, so it is not a safely idempotent in-place update that Terraform can converge to. The provider therefore treats a `layout_id` change as a replacement rather than an in-place update.

-> If a cluster's layout is changed out-of-band (for example, through the Morpheus UI), a subsequent `terraform plan` will detect the changed `layout_id` and propose a destructive replacement. Running `terraform apply -refresh-only` updates Terraform state to the observed `layout_id`, but it does **not** converge on its own: the configuration still specifies the old `layout_id`, so the next plan proposes the replacement again. To durably reconcile without recreating the cluster, either update `layout_id` in your configuration to match the reassigned layout, or add `lifecycle { ignore_changes = [layout_id, config_hvm] }` to the resource.

## Scaling HVM 2.0 clusters

~> The set of cluster hosts is fixed at create time by `server.ssh_hosts`. Changing `server.ssh_hosts` (for example, to add or remove a host) forces the cluster to be **replaced** rather than scaled in place; the provider does not support adding or removing hosts on an existing cluster. Plan host changes accordingly, since a replacement destroys and recreates the cluster.

## Import

~> Import is currently not supported for the `hpe_morpheus_cluster` resource.

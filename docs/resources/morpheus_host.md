---
page_title: "hpe_morpheus_host Resource - terraform-provider-hpe"
subcategory: "Morpheus"
description: |-
  Manages a Morpheus host — a managed compute server added through the "Add Host" flow — for both virtual and physical hosts (bare metal, hypervisors, and Kubernetes nodes).
---
# hpe_morpheus_host (Resource)

Manages a Morpheus host — a managed compute server added through the "Add Host" flow. A host is a
compute server that hosts others (a hypervisor, a bare-metal server, or a Kubernetes master/worker)
rather than a guest instance. The mainstream use case is bare metal — HPE iLO via the
hpe-baremetal-plugin, or OneView/UCS in-core — but the same API adds a Docker/Kubernetes/standalone
host on VMware, AWS, Azure, GCP, OpenStack and others.

~> **"Host" vs "compute server"** In the Morpheus API every host *is* a compute server
(`/api/servers`): a host is simply a compute server with no parent of its own. This resource is
named `hpe_morpheus_host` because it manages the "Add Host" flow; the read-only
`hpe_morpheus_compute_server`/`hpe_morpheus_compute_servers` data sources remain and cover every
compute server, including guests. Use `hpe_morpheus_host`/`hpe_morpheus_hosts` when you specifically
mean hosts.

The compute server type (`compute_server_type`) selects the kind of host and therefore which
configuration keys apply. Two typed blocks are provided for the bare-metal plugins —
`config_bmaas` (HPE iLO) and `config_omega_bmaas` (the Omega stub) — and a free-form `config` map
covers everything else (Kubernetes nodes, ESXi/Fusion/Hyper-V metal, KVM hypervisors…). At most one
of `config`, `config_bmaas` or `config_omega_bmaas` may be set; a type that needs no extra keys
needs none of them.

-> **Which hosts this resource can create.** Only the server types Morpheus marks *creatable*
(`GET /api/server-types?creatable=true`): bare-metal nodes (HPE iLO, Omega, ESXi/Fusion/Hyper-V metal,
Supermicro, …), Kubernetes masters and workers per cloud, and KVM hypervisors on a VM or metal
(`vmwareKvm`, `oneviewKvmLinux`, `ucsKvmLinux`). Two common kinds of host are **not** creatable here
because another object owns them, and the create is normally refused with a message naming it (the
check reads the appliance's type list and lets the create through if that read fails): **HVM
hypervisors** (`mvm`) are created by `hpe_morpheus_cluster` as members of their cluster (adding one
to an existing cluster is the cluster's *Add Worker*, which has no resource yet), and **virtual
machines** (`morpheusVm`, `vmwareVm`, …) are `hpe_morpheus_instance`. Every kind of host, creatable
here or not, can be read with the `hpe_morpheus_host` and `hpe_morpheus_hosts` data sources.

## Example Usage — HPE iLO bare metal (`config_bmaas`)

```terraform
# An HPE iLO bare-metal host (hpe-baremetal-plugin IloBareMetalUnmanaged type).
resource "hpe_morpheus_host" "ilo_metal" {
  name                = "ilo-metal-1"
  cloud_id            = 20
  group_id            = 1
  compute_server_type = "IloBareMetalUnmanaged"

  config_bmaas = {
    ilo_ip_address   = "10.1.2.3"
    ilo_username     = "Administrator"
    ilo_password_wo  = var.ilo_password
    ilo_mac_address  = "aa:bb:cc:dd:ee:ff"
    resource_pool_id = "1"
  }
}

variable "ilo_password" {
  type      = string
  sensitive = true
}
```

## Example Usage — Omega bare-metal stub (`config_omega_bmaas`)

```terraform
resource "hpe_morpheus_host" "example" {
  name                = "example-host"
  cloud_id            = 15174
  group_id            = 1
  compute_server_type = "omega.baremetal.stub-server"

  config_omega_bmaas = {
    console_host     = "10.0.0.1"
    console_username = "admin"
    num_nics         = 1
    resource_pool_id = "pool-1"
  }
}
```

## Example Usage — Kubernetes worker via free-form `config`

```terraform
# A Kubernetes worker on a VMware cloud, added via the free-form config map
# (the type code is the cloud's worker type from /api/server-types?creatable=true).
resource "hpe_morpheus_host" "k8s_worker" {
  name                = "k8s-worker-1"
  cloud_id            = 2
  group_id            = 1
  compute_server_type = "vmwareKubeWorker"
  plan_id             = 104

  config = {
    resourcePoolId = "pool-123"
  }
}
```

## Delete behaviour

Deletion calls `DELETE /api/servers/{id}`. The `delete_options` block carries the API's delete
query parameters `remove_resources`, `preserve_volumes`, `force` and `remove_instances` (the two
elastic/floating IP release flags are not exposed). None of them has a provider
default: an unset flag is **omitted** from the request so that Morpheus applies its own, type-aware
default. Morpheus's default for `removeResources` is `computeServerType.managed` — the underlying
thing is destroyed for managed types and kept for unmanaged ones. The two bare-metal plugin types
(`IloBareMetalUnmanaged`, `omega.baremetal.stub-server`) are `managed = false`, so by default
Morpheus keeps the node and only drops the record.

~> **Destructive wipe on bare metal** Setting `delete_options.remove_resources = true` on a
**bare-metal** host runs the plugin's `cleanServer`: it powers the machine on via iLO, deletes its
logical volumes, **sanitizes the disks**, and resets the iLO. This is a destructive, irreversible
wipe of the physical server. For a **virtual** host the same flag deletes the VM from the cloud.
Leave `remove_resources` unset to accept Morpheus's `removeResources = type.managed` default.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cloud_id` (Number) The ID of the cloud to add the host to. Changing it forces replacement.
- `compute_server_type` (String) The code of the compute server type to create (for example
`omega.baremetal.stub-server` or an iLO bare-metal type). The type must be
creatable. Changing it forces replacement.
- `group_id` (Number) The ID of the group to add the host to. Changing it forces replacement.
- `name` (String) The name of the host.

### Optional

> **NOTE**: [Write-only arguments](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments) are supported in Terraform 1.11 and later.

- `config` (Dynamic) Free-form, type-specific configuration for host types that have no typed config_*
block, as a map of the keys the type's option types read. Mutually exclusive with
the typed config_* blocks. Changing it forces replacement.
- `config_bmaas` (Attributes) Typed configuration for an HPE iLO bare-metal host (the hpe-baremetal-plugin
IloBareMetalUnmanaged type). Mutually exclusive with config and config_omega_bmaas.
Changing it forces replacement. (see [below for nested schema](#nestedatt--config_bmaas))
- `config_omega_bmaas` (Attributes) Typed configuration for an Omega bare-metal stub host (the omega.baremetal.stub-server
type). Mutually exclusive with config and config_bmaas. Changing it forces replacement. (see [below for nested schema](#nestedatt--config_omega_bmaas))
- `delete_options` (Attributes) Options controlling what happens to the underlying resources when the host is
deleted. Each flag is optional with no provider default; an unset flag is omitted
from the DELETE request so Morpheus applies its own type-aware default. (see [below for nested schema](#nestedatt--delete_options))
- `description` (String) A description of the host.
- `enabled` (Boolean) Whether the host is enabled.
- `labels` (Set of String) Organization labels (keywords) assigned to the host. When not configured, labels set on the
appliance are read into state without planning a change.
- `network_interfaces` (Attributes List) Network interface configuration for the host. Create-only; changing it forces
replacement. (see [below for nested schema](#nestedatt--network_interfaces))
- `plan_id` (Number) The ID of the service plan (sizing) for the host. When omitted, Morpheus assigns a default plan for the compute server type. Changing it forces replacement.
- `ssh_key_pair_id` (Number) The ID of the SSH key pair to assign to the host.
- `ssh_password_wo` (String, Sensitive, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) SSH password (Write Only). Also requires ssh_username to be set.
- `ssh_username` (String) SSH username Morpheus uses to reach the host when it is added. Morpheus's own provisioning may
replace it with the account it creates on the node (for example `cloudinit`); the configured value
is the one the provider manages, and the appliance's current value is only read into state when
the attribute is not configured (such as on import).
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `user_data` (String) Cloud-init user data for the host, carried in config.userData. Updatable in place.
- `visibility` (String) The visibility of the host (private or public).
- `volumes` (Attributes List) Volume (storage) configuration for the host. Create-only; changing it forces
replacement. (see [below for nested schema](#nestedatt--volumes))

### Read-Only

- `agent_installed` (Boolean) Whether the Morpheus agent is installed on the host.
- `agent_version` (String) The version of the Morpheus agent installed on the host.
- `bare_metal` (Boolean) Whether the host is a bare-metal host (from the compute server type).
- `compute_server_type_id` (Number) The ID of the compute server type.
- `date_created` (String) The date the host was created.
- `external_id` (String) The external ID of the host.
- `external_ip` (String) The external IP address of the host.
- `hostname` (String) The hostname of the host.
- `id` (Number) The ID of the host (compute server).
- `internal_ip` (String) The internal IP address of the host.
- `last_updated` (String) The date the host was last updated.
- `managed` (Boolean) Whether the host is managed by Morpheus.
- `parent_host_id` (Number) The ID of the parent hypervisor host, if any.
- `power_state` (String) The power state of the host (e.g. on, off).
- `status` (String) The status of the host (e.g. provisioning, provisioned, running).
- `uuid` (String) The UUID of the host.
- `write_only_revision` (Number) The time, in seconds since the Unix epoch, at which a write-only value (ssh_password_wo,
config_bmaas.ilo_password_wo or config_omega_bmaas.console_password_wo) was last sent to
Morpheus. Changes whenever one of those values changes; not configurable.

<a id="nestedatt--config_bmaas"></a>
### Nested Schema for `config_bmaas`

Required:

- `ilo_ip_address` (String) The iLO IP address of the bare-metal host.
- `resource_pool_id` (String) The ID of the resource pool to provision the host into.

Optional:

- `ilo_mac_address` (String) The iLO MAC address of the bare-metal host.
- `ilo_password_wo` (String, Sensitive, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) The iLO password (Write Only).
- `ilo_username` (String) The iLO username.
- `pre_provisioned` (Boolean) Whether the bare-metal host is pre-provisioned.


<a id="nestedatt--config_omega_bmaas"></a>
### Nested Schema for `config_omega_bmaas`

Required:

- `resource_pool_id` (String) The ID of the resource pool to provision the host into.

Optional:

- `console_host` (String) The console host address of the Omega stub host.
- `console_password_wo` (String, Sensitive, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) The console password (Write Only).
- `console_username` (String) The console username.
- `num_nics` (Number) The number of NICs for the Omega stub host.
- `pre_provisioned` (Boolean) Whether the host is pre-provisioned.


<a id="nestedatt--delete_options"></a>
### Nested Schema for `delete_options`

Optional:

- `force` (Boolean) Whether to force removal even if errors occur.
- `preserve_volumes` (Boolean) Whether to preserve volumes when removing resources.
- `remove_instances` (Boolean) Whether to also remove instances associated with the host.
- `remove_resources` (Boolean) Whether to destroy the underlying machine. For a virtual host this deletes the
VM; for a bare-metal host it runs the plugin's destructive wipe (sanitizes
disks, resets iLO). When unset, Morpheus defaults to the compute server type's
managed flag.


<a id="nestedatt--network_interfaces"></a>
### Nested Schema for `network_interfaces`

Optional:

- `ip_address` (String) The static IP address for this interface. Not applicable with DHCP or IP pools.
- `ip_mode` (String) The mode for determining the IP address (static, dhcp, or "").
- `network_id` (Number) The ID of the network to attach this interface to.
- `network_interface_type_id` (Number) The ID of the network interface type (adapter model), when the cloud offers a choice.


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


<a id="nestedatt--volumes"></a>
### Nested Schema for `volumes`

Optional:

- `datastore_id` (Number) The ID of the datastore for the volume.
- `name` (String) The name of the volume.
- `root_volume` (Boolean) Whether this is the root volume.
- `size` (Number) The size of the volume in GB.
- `storage_type_id` (Number) The ID of the storage type for the volume.

## Import

Hosts are imported by their numeric ID. Create-only inputs that the API does not return — the
`config_*` blocks (with their write-only iLO/console passwords), `network_interfaces` and `volumes` —
are not populated on import. **Leave them out of the configuration of an imported host**: they all
force replacement, so adding a `config_*` block, `volumes` or `network_interfaces` after import plans
a destroy and recreate (on a bare-metal host, the destructive path warned about above). The
in-place attributes (`name`, `description`, `labels`, `visibility`, `enabled`, `ssh_username`,
`ssh_key_pair_id`, `user_data`) can be managed as usual; `delete_options` only affects destroy and
may be set.

`ssh_password_wo` is updated in place but, being write-only, is not stored in state, so the import
cannot tell whether the value in your configuration matches the appliance: it is taken as already
applied, the plan after import is empty, and a later change to it is only detected once the host is
next updated for another reason (Terraform warns about this on every plan until then). To have
changes detected from the start, append the sha256 digest of the configured value to the import ID:

```shell
terraform import hpe_morpheus_host.example "123#ssh_password_wo=<sha256 hex>"
```

The same ID works in an `import` block, where Terraform's `sha256()` function produces the digest:
`id = "123#ssh_password_wo=${sha256(var.ssh_password)}"`.

```shell
# Hosts are imported by numeric ID. Create-only inputs the API does not return
# (the config_* blocks, network_interfaces and volumes) are not populated and
# would plan a replacement if configured, so leave them out of an imported
# host's configuration. delete_options may be set, and ssh_password_wo is
# updated in place (see the Import section for its digest form).
terraform import hpe_morpheus_host.example 123
```

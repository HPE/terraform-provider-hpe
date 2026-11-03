---
page_title: "hpe_morpheus_network Resource - terraform-provider-hpe"
subcategory: "Morpheus"
description: |-
  
---
# hpe_morpheus_network (Resource)



### Networks vs Subnets

For most network types in Morpheus, the network itself is the deployable unit -- there is no separate subnet creation step. Only Azure and Google support a two-level model where subnets can be created independently within a network using the [`hpe_morpheus_subnet`](https://registry.terraform.io/providers/HPE/hpe/latest/docs/resources/morpheus_subnet) resource. For all other types (NSX-T, OpenStack, Amazon, OVS, VCD, ACI, etc.), use this resource to create the full network object including its addressing configuration.

### Choosing the network type

Name the type with `type_code` (for example `ovsPortGroup`, `nsxtLogicalSwitch`, `azureNetwork`) or
with `type_id`; set one and the other is computed. Codes are the same on every appliance, where ids
are not, so `type_code` travels between environments. The
[`hpe_morpheus_network_type`](https://registry.terraform.io/providers/HPE/hpe/latest/docs/data-sources/morpheus_network_type)
data source lists a type's details, and `GET /api/network-types` lists them all. Changing the type
replaces the network.

### Type-specific configuration

A Morpheus network type declares its fields in two places, and this resource follows that split:

- **Fields of the network itself** — `cidr`, `gateway`, `vlan_id`, `switch_id`, `zone_pool_id`,
  `pool_id`, `network_domain_id`, `assign_public_ip`, `availability_zone`, the IPv6 set — are
  ordinary top-level attributes, sent and read back like any other.
- **Fields of the type's configuration map** are carried by a **typed block** named after the type
  family, `config_<type>`, or, for a type without a block, by the free-form `config` map. These are
  **write-only**: Terraform sends them to Morpheus on create and update and never stores them in
  state, so they can hold values Morpheus masks (secrets) and the configuration can be changed
  without a diff appearing in state. What Morpheus holds is readable in the computed
  `config_from_api`, including keys the server adds itself and secrets it masks.

Write-only attributes need **Terraform 1.11 or later, or OpenTofu 1.11 or later**.

A network takes **at most one** typed block, and it must be the block for the network's type; a
mismatch is reported when the configuration is validated (with `type_code`) or planned (with
`type_id`). The free-form `config` is for types that have no block — types a network plugin adds,
for instance — and cannot be combined with a typed block.

| Block | Network types | Attributes (those that **replace** the network when changed are marked) |
|---|---|---|
| `config_ovs_port_group` | `ovsPortGroup` | `vlan_trunks`, `trunk_only` |
| `config_hvm_data` | `hvmPortGroupDataNetwork` | `cluster_network_layout` ⟳, `host_network_interface` ⟳, `virtual_switch_segment_id` ⟳, `mtu` |
| `config_hvm_overlay` ¹ | `hvmPortGroupOverlayNetwork` | `overlay_network_interface` ⟳, `overlay_port` ⟳, `overlay_id` ⟳ |
| `config_hvm_private` ¹ | `hvmPortGroupPrivateNetwork` | `cluster_network_layout` ⟳, `host_network_interface` ⟳, `virtual_switch_id` ⟳ |
| `config_hvm_sriov` ¹ | `hvmPortGroupSriovNetwork` | `physical_function` ⟳ |
| `config_nsxt` | `nsxtLogicalSwitch` | `connected_gateway`, `vlan_ids`, `subnet_ip_management_type`, `subnet_ip_server_id`, `subnet_dhcp_server_address`, `dhcp_range`, `subnet_dhcp_lease_time`, and the IPv6 DHCP set (`subnet_dhcp_server_address_ipv6`, `dhcp_range_ipv6`, `dhcp_excluded_ranges_ipv6`, `subnet_dhcp_lease_time_ipv6`, `dhcp_domain_names_ipv6`, `dhcp_preferred_time_ipv6`, `dhcp_sntp_servers_ipv6`) |
| `config_nsx_cloud` ¹ | `nsxCloudLogicalSwitch` | `connected_gateway`, `vlan_ids`, `subnet_ip_management_type`, `subnet_ip_server_id`, `subnet_dhcp_server_address`, `dhcp_range`, `subnet_dhcp_lease_time` |
| `config_azure` | `azureNetwork` | `resource_group_id` ⟳, `subnet_name` ⟳, `subnet_cidr` ⟳ |
| `config_gcp` | `googleNetwork` | `auto_create` ⟳, `mtu` |
| `config_openstack` | `openstackPrivate`, `huaweiPrivate`, `opentelekomPrivate` | `dhcp_range`, `host_routes` ⟳ |

⟳ — Morpheus does not apply a new value to this field after the network is created, so a change is
planned as a replacement rather than accepted and silently ignored.
¹ — generated from the network type's definition and checked against it, but not exercised against
a live cloud by this provider's tests; reports welcome.

#### Which attributes an appliance supports

The blocks describe the network types as Morpheus 9.1.0 declares them. Everything in them is also
declared by Morpheus 9.0.1 **except** the virtual-switch layout fields, which arrived with the HVM
2.0 cluster layout in 9.1: `config_hvm_data.cluster_network_layout` and `.virtual_switch_segment_id`,
`config_hvm_private.cluster_network_layout` and `.virtual_switch_id`. Because a network type's
fields also depend on the plugin that provides it, the provider does not gate on a version number:
at plan time it reads the type's definition from the appliance and **rejects a block attribute the
appliance's type does not declare**, naming the field and the appliance's version — Morpheus would
otherwise accept the value and silently ignore it. A key of the free-form `config` that the type does
not declare is a warning. The check needs the appliance at plan time; when it cannot be reached the
check runs at apply instead.

Some types have no configuration map at all, and so no block: an Amazon subnet's fields
(`zone_pool_id`, `availability_zone`, `cidr`, `assign_public_ip`) and a custom child network's
(`parent_network_id`, `cidr`, `gateway`, DNS) are all top-level attributes; an OVS overlay's are
`zone_pool_id`, `switch_id` and `vxlan_id`.

~> **Deprecated** The top-level `connected_gateway` is superseded by `connected_gateway` inside
`config_nsxt` and `config_nsx_cloud`. It still works, and is read back from the network's
configuration, so existing configurations need not change immediately.

#### How changes to write-only configuration are detected

A write-only value is never in state, so Terraform cannot compare it with anything. Following
HashiCorp's guidance for write-only arguments, the provider keeps a secure hash of the configuration
it last applied in the resource's private state and compares the configuration against it at plan
time: a change to a field Morpheus edits in place is applied by an update — shown in the plan as a
change to `config_from_api` — and a change to a field marked ⟳ above replaces the network. No
version or "keepers" attribute is needed. This detects changes you make in the configuration; a
change made to the configuration on the appliance is visible in `config_from_api` but is not undone.

#### What an update applies, per type

Whether Morpheus applies a change to `name`, `cidr`, `gateway`, `vlan_id`, `dhcp_server` or the DNS
attributes depends on the network type: each type carries flags saying which of them it edits after
create, and a change to a field whose flag is off is accepted by the API and ignored. At plan time
the provider reads the network type from the appliance and plans a **replacement** for such a
change, so it never plans an update that would do nothing; `name` on an HVM port group or `cidr` on
an Azure virtual network are examples. A change to a field the type does edit — `vlan_id` on an OVS
port group, `name` on an NSX-T segment — is an in-place update. The type is read once per type per
run; when the appliance cannot be reached at plan time, a table of the network types of a recent
Morpheus release built into the provider is used instead, and a type in neither is treated
generically.

### Removing an attribute clears it

Removing an optional attribute from the configuration **clears it on the network**: `description`,
`display_name`, `gateway`, the DNS servers, `search_domains`, `switch_id`, `vlan_id`, the IPv6
addressing (`cidr_ipv6`, `gateway_ipv6`, `netmask_ipv6`, `dns_primary_ipv6`, `dns_secondary_ipv6`),
`pool_id`, `pool_ipv6_id`, `network_domain_id`, `network_proxy_id`, `zone_pool_id`, `no_proxy`,
`labels` and `cidr` are planned to `null` and read back empty after the apply. Removing
`tenant_ids` revokes every tenant's access but the owner's — Morpheus always keeps the owning
tenant — so it is planned as `(known after apply)` and reads back as the owner alone. Whether the
clear is applied in place or by replacement follows the type's flags exactly as a change does:
removing `gateway` from a network whose type does not edit the gateway after create replaces the
network.

Three things are not cleared. A boolean keeps its value when it is removed (there is no empty
boolean). A value the user never set — the display name Morpheus fills in from the name, or a
gateway a cloud sync reported — is not the provider's to remove, so it stays: the provider records,
at every apply, which of these attributes the configuration set, and only a recorded attribute is
cleared when it goes. And `zone_pool_id` cannot be removed from a network of the HVM port group
types, whose pool is the cluster the network lives on; the plan refuses it.

~> **Note** The record is written by the first apply made with this provider version. A network
last applied by an earlier version has none yet, so an attribute removed from its configuration
*before* that apply is kept, as before; apply once, then remove it. An imported network behaves the
same way until its first apply.

### Network server

Network types managed by a network plugin — the HVM port group types among them — belong to a
**network server**, and Morpheus 9.1 requires it on create: `network_server_id`. For an HVM cluster
the network server is the cluster's own, hidden from the ordinary listing; find it with
`GET /api/networks/servers?showHidden=true` (it is named after the cluster, `<cluster> HVM Network
Server`). Morpheus does not return it on a network, so the configured value is kept in state.

### Visibility and tenants

Morpheus creates a network as `private` unless its tenant permissions include the master tenant, so
`visibility = "public"` requires `tenant_ids` to include the master tenant's id; the provider reports
this at plan time rather than letting the apply fail with an inconsistent result. A sub-tenant cannot
create a public network at all.

### Timing

Morpheus answers a create as soon as the network record exists and finishes provisioning — on a
cloud, or on an HVM cluster's hosts — afterwards, so the provider waits for the network to leave the
`provisioning` state before reading it back (a failed provisioning is an error at apply, with the
network left for inspection). On destroy it waits until the network is gone, and reports Morpheus's
own message if the delete is refused.

Morpheus keeps one network record **per cloud**. When two clouds cover the same account, the network
you create through one is discovered by the other on its next refresh as a separate record under the
same name, and that record disappears on the refresh after you destroy yours. This does not affect
the resource, which works by id, but it does affect lookups by name — see the
[`hpe_morpheus_network` data source](https://registry.terraform.io/providers/HPE/hpe/latest/docs/data-sources/morpheus_network)
and its `cloud_id`.

~> **Note** Some network types (for example, Amazon) update the network `name`
and some other attributes a few minutes after creation. This requires using a
`lifecycle` block as shown below. If this `lifecycle` block is missing, then
a subsequent `terraform apply` may attempt to delete the network.

If required, a `lifecycle` block may be added as follows:

```hcl
resource "hpe_morpheus_network" "a" {
  name = "network A"
  display_name = "network A"
  description = "First network"

  .
  .
  .

  lifecycle {
    ignore_changes = [ name, display_name, description ]
  }
}
```

### Import

Import by network id. The write-only configuration cannot be read from the appliance, so after an
import the first plan shows one update that sends the configured `config_*` block or `config` map
and records it; `config_from_api` is populated from the start. The never-returned attributes
(`network_server_id`, `scope_id`, `vxlan_id`, `parent_network_id`) take the configured values on that
first apply.

### Upgrading from 2.1

- `config` is no longer stored in state. A reference to `hpe_morpheus_network.x.config` elsewhere in
  a configuration must become `hpe_morpheus_network.x.config_from_api`. Existing state is upgraded
  automatically and the first plan after upgrading shows no change to the network.
- `type_id` is no longer required; either it or `type_code` is.
- A change to `tenant_ids` is applied in place; it used to replace the network.
- `name`, `cidr` and `cidr_ipv6` no longer always replace the network: they do so on the types that
  cannot edit them, and are updated in place on the types that can.

## Example Usage

### AWS Network

An Amazon subnet has no configuration block: its fields are attributes of the network.

```terraform
data "hpe_morpheus_cloud" "example" {
  name = "AWS Cloud"
}

data "hpe_morpheus_group" "example" {
  name = "Example Group"
}

data "hpe_morpheus_tenant" "example" {
  name = "Master Tenant"
}

resource "hpe_morpheus_network" "aws" {
  name        = "example-terraform-aws"
  description = "AWS subnet"
  cloud_id    = data.hpe_morpheus_cloud.example.id
  pool_id     = 1
  group_id    = data.hpe_morpheus_group.example.id
  type_code   = "amazonSubnet"
  # An Amazon subnet has no config block: its fields are attributes of the
  # network itself.
  assign_public_ip           = true
  availability_zone          = "us-west-1a"
  active                     = true
  dhcp_server                = true
  appliance_url_proxy_bypass = true
  tenant_ids = [
    data.hpe_morpheus_tenant.example.id,
  ]
  visibility   = "private"
  cidr         = "10.200.99.0/24"
  zone_pool_id = 12329
  labels       = ["terraform", "example"]

  lifecycle {
    ignore_changes = [name, display_name, description]
  }
}
```

### Azure Network

The Azure region is derived from the cloud or the resource group; it is not configured on the
network. All three attributes of `config_azure` replace the network when changed.

```terraform
data "hpe_morpheus_cloud" "example" {
  name = "Azure Cloud"
}

data "hpe_morpheus_group" "example" {
  name = "Example Group"
}

data "hpe_morpheus_tenant" "example" {
  name = "Master Tenant"
}

resource "hpe_morpheus_network" "azure" {
  name                       = "example-terraform-azure"
  description                = "Azure network"
  cloud_id                   = data.hpe_morpheus_cloud.example.id
  pool_id                    = 1
  group_id                   = data.hpe_morpheus_group.example.id
  type_code                  = "azureNetwork"
  cidr                       = "10.100.0.0/16"
  visibility                 = "public"
  active                     = true
  dhcp_server                = true
  appliance_url_proxy_bypass = false
  labels                     = ["terraform", "example"]
  config_azure = {
    resource_group_id = "all-attrs-resource-group"
    subnet_name       = "all-attrs-subnet"
    subnet_cidr       = "10.100.1.0/24"
  }
  tenant_ids = [
    data.hpe_morpheus_tenant.example.id,
  ]
}
```

### GCP Network

```terraform
data "hpe_morpheus_cloud" "example" {
  name = "Google Cloud"
}

data "hpe_morpheus_group" "example" {
  name = "Examle Group"
}

data "hpe_morpheus_tenant" "example" {
  name = "Master Tenant"
}

resource "hpe_morpheus_network" "gcp" {
  name        = "example-terraform-gcp"
  description = "GCP network"
  cloud_id    = data.hpe_morpheus_cloud.example.id
  pool_id     = 1
  group_id    = data.hpe_morpheus_group.example.id
  type_code   = "googleNetwork"
  config_gcp = {
    mtu         = 1460
    auto_create = true
  }
  active                     = true
  dhcp_server                = false
  appliance_url_proxy_bypass = true
  tenant_ids = [
    data.hpe_morpheus_tenant.example.id,
  ]
  visibility   = "private"
  cidr         = "10.0.0.0/8"
  zone_pool_id = 85990
  labels       = ["terraform", "example"]
}
```

### NSX-T Segment

NSX-T segments are networks of type `nsxtLogicalSwitch`. The segment's own addressing is the
top-level `cidr` and `gateway`; NSX-T subnets are not created separately. `connected_gateway` is
the path of the Tier-0 or Tier-1 gateway (`/infra/tier-1s/<id>`), `vlan_ids` is set for a
VLAN-backed segment and left empty for an overlay segment, and `subnet_ip_management_type` selects
the DHCP mode (`dhcpLocal`, `gatewayDhcp`, `dhcpRelay`) with `subnet_ip_server_id` naming the server
or relay, `subnet_dhcp_server_address` the server's address in CIDR form (required for `dhcpLocal`),
`dhcp_range` the range (`172.16.10.100-172.16.10.200`) and `subnet_dhcp_lease_time` the lease in
seconds (NSX-T's default is 86400).

```terraform
data "hpe_morpheus_cloud" "nsxt" {
  name = "NSX-T Cloud"
}

data "hpe_morpheus_group" "example" {
  name = "Example Group"
}

# The NSX-T integration the segment is created through.
data "hpe_morpheus_network_server" "nsxt" {
  name = "Labs NSX"
}

resource "hpe_morpheus_network" "nsxt_segment" {
  name              = "example-terraform-nsxt-segment"
  display_name      = "Example NSX-T Segment"
  description       = "NSX-T overlay segment managed by Terraform"
  cloud_id          = data.hpe_morpheus_cloud.nsxt.id
  group_id          = data.hpe_morpheus_group.example.id
  type_code         = "nsxtLogicalSwitch"
  network_server_id = data.hpe_morpheus_network_server.nsxt.id
  cidr              = "172.16.10.0/24"
  gateway           = "172.16.10.1"
  active            = true
  dhcp_server       = true

  # Write-only: sent to Morpheus, never stored in state. Read back in config_from_api.
  config_nsxt = {
    connected_gateway          = "/infra/tier-1s/my-tier1-gw"
    vlan_ids                   = ""
    subnet_ip_management_type  = "dhcpLocal"
    subnet_dhcp_server_address = "172.16.10.2/24"
    dhcp_range                 = "172.16.10.100-172.16.10.200"
    subnet_dhcp_lease_time     = 86400
  }
}
```

### Host Network

```terraform
data "hpe_morpheus_cloud" "example" {
  name = "Standard Cloud"
}

data "hpe_morpheus_group" "example" {
  name = "Example Group"
}

data "hpe_morpheus_tenant" "example" {
  name = "Master Tenant"
}

resource "hpe_morpheus_network" "host" {
  name                       = "example-terraform-host"
  description                = "A host network"
  cloud_id                   = data.hpe_morpheus_cloud.example.id
  pool_id                    = 1
  group_id                   = data.hpe_morpheus_group.example.id
  type_id                    = 1
  config                     = {}
  active                     = true
  dhcp_server                = false
  appliance_url_proxy_bypass = true
  tenant_ids = [
    data.hpe_morpheus_tenant.example.id,
  ]
  visibility = "private"
  cidr       = "10.0.0.0/8"
  labels     = [terraform, example]
}
```

### OVS Port Group Network (HVM standard network)

The port group's virtual switch, VLAN and resource pool are the top-level `switch_id`, `vlan_id`
and `zone_pool_id`; `switch_id` replaces the network when changed, `vlan_id` is updated in place.

```terraform
data "hpe_morpheus_cloud" "example" {
  name = "HVM Cloud"
}

data "hpe_morpheus_group" "example" {
  name = "ExampleGroup"
}

data "hpe_morpheus_tenant" "example" {
  name = "Master Tenant"
}

resource "hpe_morpheus_network" "ovs_port_group" {
  name        = "Terraform OVS Port Group"
  description = "OVS Port Group network"
  cloud_id    = data.hpe_morpheus_cloud.example.id
  pool_id     = 3251
  group_id    = data.hpe_morpheus_group.example.id
  type_code   = "ovsPortGroup"
  switch_id   = "Compute"
  # The cluster's HVM network server; list them with
  # GET /api/networks/servers?showHidden=true.
  network_server_id = 7834
  config_ovs_port_group = {
    vlan_trunks = "100,200"
    trunk_only  = false
  }
  active                     = true
  dhcp_server                = false
  appliance_url_proxy_bypass = true
  tenant_ids = [
    data.hpe_morpheus_tenant.example.id,
  ]
  visibility   = "public"
  cidr         = "10.32.148.0/22"
  zone_pool_id = 62299
  vlan_id      = 43
  labels       = ["terraform", "example"]

  lifecycle {
    ignore_changes = [name, display_name, description]
  }
}
```

### HVM Data Network

The `HOST_INTERFACE` layout shown attaches a VLAN on a host interface; on an HVM 2.0 cluster the
layout is `VSWITCH` and the network takes a data-traffic virtual switch segment instead
(`cluster_network_layout = "VSWITCH"`, `virtual_switch_segment_id`), with the segment deciding the
VLAN and MTU.

```terraform
data "hpe_morpheus_cloud" "hvm" {
  name = "HVM Cloud"
}

data "hpe_morpheus_group" "example" {
  name = "Example Group"
}

# An HVM data network attaches a VLAN on a host interface to the cluster. Its
# cluster is the resource pool; its network server is the cluster's own, which
# is hidden from the ordinary listing (GET /api/networks/servers?showHidden=true).
resource "hpe_morpheus_network" "hvm_data" {
  name              = "example-terraform-hvm-data"
  cloud_id          = data.hpe_morpheus_cloud.hvm.id
  group_id          = data.hpe_morpheus_group.example.id
  type_code         = "hvmPortGroupDataNetwork"
  zone_pool_id      = 153047
  network_server_id = 7834
  vlan_id           = 2104
  visibility        = "private"
  tenant_ids        = [1]

  # Write-only. The layout and interface cannot be changed after create; the
  # MTU can (1500, or 9000 for jumbo frames).
  config_hvm_data = {
    cluster_network_layout = "HOST_INTERFACE"
    host_network_interface = "bond0"
    mtu                    = 1500
  }
}
```

### OpenStack Private Network

```terraform
data "hpe_morpheus_cloud" "openstack" {
  name = "OpenStack Cloud"
}

data "hpe_morpheus_group" "example" {
  name = "Example Group"
}

data "hpe_morpheus_network_server" "neutron" {
  name = "OpenStack Cloud"
}

resource "hpe_morpheus_network" "openstack" {
  name              = "example-terraform-openstack"
  cloud_id          = data.hpe_morpheus_cloud.openstack.id
  group_id          = data.hpe_morpheus_group.example.id
  type_code         = "openstackPrivate"
  zone_pool_id      = 46930
  network_server_id = data.hpe_morpheus_network_server.neutron.id
  cidr              = "192.168.177.0/24"
  gateway           = "192.168.177.1"
  dhcp_server       = true
  visibility        = "private"
  tenant_ids        = [1]

  # Write-only. The DHCP allocation pool is updated in place; host routes
  # cannot be changed after create.
  config_openstack = {
    dhcp_range = "192.168.177.10,192.168.177.100"
  }
}
```

### Custom Child Network

A child network carves a subnet out of an existing (typically vSphere) network; it has no
configuration block.

```terraform
data "hpe_morpheus_cloud" "vsphere" {
  name = "vSphere Cloud"
}

data "hpe_morpheus_group" "example" {
  name = "Example Group"
}

# The network the child is carved from, looked up in the cloud it belongs to.
data "hpe_morpheus_network" "parent" {
  name     = "VM Network"
  cloud_id = data.hpe_morpheus_cloud.vsphere.id
}

# A custom child network has no configuration block: its parent, addressing
# and DNS are attributes of the network. Its name and CIDR are updated in place.
resource "hpe_morpheus_network" "child" {
  name              = "example-terraform-child"
  cloud_id          = data.hpe_morpheus_cloud.vsphere.id
  group_id          = data.hpe_morpheus_group.example.id
  type_code         = "childNetwork"
  parent_network_id = data.hpe_morpheus_network.parent.id
  cidr              = "10.77.0.0/24"
  gateway           = "10.77.0.1"
  visibility        = "private"
  tenant_ids        = [1]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cloud_id` (Number) Cloud (zone) id
- `group_id` (Number) Group (site) id
- `name` (String) Network name. Whether a change is applied in place or replaces the network depends on the network type: Morpheus ignores a new name on types whose `nameEditable` flag is off, so the provider plans a replacement for those.

### Optional

> **NOTE**: [Write-only arguments](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments) are supported in Terraform 1.11 and later.

- `active` (Boolean) Activate (true) or disable (false) the network
- `allow_static_override` (Boolean) Allow IP Override
- `appliance_url_proxy_bypass` (Boolean) Bypass Proxy for Appliance URL
- `assign_public_ip` (Boolean) Assign Public IP
- `availability_zone` (String) Availability zone of an Amazon subnet, e.g. `us-west-1a` (Morpheus `availabilityZone`). Set on create only; a change replaces the network.
- `cidr` (String) Network CIDR. Applied in place on network types whose `cidrEditable` flag is on; replaces the network otherwise.
- `cidr_ipv6` (String) Network IPv6 CIDR. Applied in place on network types whose `cidrEditable` flag is on; replaces the network otherwise.
- `config` (Dynamic, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Type-specific configuration for network types that have no typed `config_*` block, as a map of the keys the type's option types read (for example `{ hvmVlanTrunks = "100,200" }`). Write-only: sent to Morpheus on create and update, never stored in state; a change is detected from a digest kept in the resource's private state, so no version attribute is needed. The server's view of the configuration is available read-only in `config_from_api`. Requires Terraform or OpenTofu 1.11 or later.
- `config_azure` (Attributes, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Configuration for an Azure virtual network (type `azureNetwork`), typed. Write-only, like `config`; the two are mutually exclusive, and a network takes at most one typed block, which must match its type. Attributes marked "cannot be changed after create" replace the network when changed; the others are applied in place. (see [below for nested schema](#nestedatt--config_azure))
- `config_gcp` (Attributes, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Configuration for a Google Cloud VPC network (type `googleNetwork`), typed. Write-only, like `config`; the two are mutually exclusive, and a network takes at most one typed block, which must match its type. Attributes marked "cannot be changed after create" replace the network when changed; the others are applied in place. (see [below for nested schema](#nestedatt--config_gcp))
- `config_hvm_data` (Attributes, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Configuration for an HVM data network (type `hvmPortGroupDataNetwork`), typed. Write-only, like `config`; the two are mutually exclusive, and a network takes at most one typed block, which must match its type. Attributes marked "cannot be changed after create" replace the network when changed; the others are applied in place. (see [below for nested schema](#nestedatt--config_hvm_data))
- `config_hvm_overlay` (Attributes, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Configuration for an HVM overlay network (type `hvmPortGroupOverlayNetwork`), typed. Write-only, like `config`; the two are mutually exclusive, and a network takes at most one typed block, which must match its type. Attributes marked "cannot be changed after create" replace the network when changed; the others are applied in place. (see [below for nested schema](#nestedatt--config_hvm_overlay))
- `config_hvm_private` (Attributes, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Configuration for an HVM private network (type `hvmPortGroupPrivateNetwork`), typed. Write-only, like `config`; the two are mutually exclusive, and a network takes at most one typed block, which must match its type. Attributes marked "cannot be changed after create" replace the network when changed; the others are applied in place. (see [below for nested schema](#nestedatt--config_hvm_private))
- `config_hvm_sriov` (Attributes, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Configuration for an HVM SR-IOV network (type `hvmPortGroupSriovNetwork`), typed. Write-only, like `config`; the two are mutually exclusive, and a network takes at most one typed block, which must match its type. Attributes marked "cannot be changed after create" replace the network when changed; the others are applied in place. (see [below for nested schema](#nestedatt--config_hvm_sriov))
- `config_nsx_cloud` (Attributes, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Configuration for an NSX Cloud segment (type `nsxCloudLogicalSwitch`), typed. Write-only, like `config`; the two are mutually exclusive, and a network takes at most one typed block, which must match its type. Attributes marked "cannot be changed after create" replace the network when changed; the others are applied in place. (see [below for nested schema](#nestedatt--config_nsx_cloud))
- `config_nsxt` (Attributes, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Configuration for an NSX-T segment (type `nsxtLogicalSwitch`), typed. Write-only, like `config`; the two are mutually exclusive, and a network takes at most one typed block, which must match its type. Attributes marked "cannot be changed after create" replace the network when changed; the others are applied in place. (see [below for nested schema](#nestedatt--config_nsxt))
- `config_openstack` (Attributes, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Configuration for an OpenStack-family private network (types `openstackPrivate`, `huaweiPrivate`, `opentelekomPrivate`), typed. Write-only, like `config`; the two are mutually exclusive, and a network takes at most one typed block, which must match its type. Attributes marked "cannot be changed after create" replace the network when changed; the others are applied in place. (see [below for nested schema](#nestedatt--config_openstack))
- `config_ovs_port_group` (Attributes, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Configuration for an OVS port group network (type `ovsPortGroup`), typed. Write-only, like `config`; the two are mutually exclusive. The port group's virtual switch, VLAN and resource pool are the top-level `switch_id`, `vlan_id` and `zone_pool_id`. (see [below for nested schema](#nestedatt--config_ovs_port_group))
- `connected_gateway` (String, Deprecated) Path of the connected NSX-T gateway. Deprecated: set `connected_gateway` inside `config_nsxt` or `config_nsx_cloud` instead; this attribute is kept for existing configurations and is read back from the network's configuration.
- `description` (String) Description
- `dhcp_server` (Boolean) DHCP Server enabled network
- `dhcp_server_ipv6` (Boolean) IPv6 DHCP Server enabled network
- `display_name` (String) Display Name
- `dns_primary` (String) Primary DNS Server
- `dns_primary_ipv6` (String) Primary IPv6 DNS Server
- `dns_secondary` (String) Secondary DNS Server
- `dns_secondary_ipv6` (String) Secondary IPv6 DNS Server
- `gateway` (String) Network Gateway
- `gateway_ipv6` (String) IPv6 Network Gateway
- `ipv4enabled` (Boolean)
- `ipv6enabled` (Boolean)
- `labels` (Set of String) Array of label strings, can be used for filtering.
- `netmask_ipv6` (String)
- `network_domain_id` (Number) Network domain id
- `network_proxy_id` (Number) Network proxy id
- `network_server_id` (Number) Network server (integration) id. Required by network types that a network plugin manages (`hasNetworkServer`), such as the HVM port group types, whose network server is the cluster's hidden "HVM Network Server" — list it with `GET /api/networks/servers?showHidden=true`.
- `no_proxy` (String) Comma-separated list of ip addresses or name servers to exclude proxy traversal for. Typically locally routable servers are excluded.
- `parent_network_id` (Number) Parent network of a custom child network (type `childNetwork`; Morpheus `parentNetwork.id`). Never returned by the API, so the configured value is kept in state.
- `pool_id` (Number) Network Pool ID
- `pool_ipv6_id` (Number) IPv6 Network Pool ID
- `scope_id` (String) NSX-T transport zone id for a segment on NSX-T 2.x (Morpheus `scopeId`). Morpheus accepts it on create and never returns it, so the configured value is kept in state; a change replaces the network.
- `search_domains` (String) Search Domains
- `switch_id` (String) Network switch identifier
- `tenant_ids` (Set of Number) Tenant account ids allowed access to the network; updated in place. A network is only `public` when this includes the master tenant — Morpheus makes it `private` otherwise.
- `type_code` (String) Network type code, e.g. `ovsPortGroup`, `nsxtLogicalSwitch`, `azureNetwork`; the codes are stable across appliances where ids are not. Set either `type_code` or `type_id`; the other is computed. Changing the type replaces the network.
- `type_id` (Number) Network type id. Set either `type_id` or `type_code`; the other is computed. Changing the type replaces the network.
- `visibility` (String) Visibility, private or public.
- `vlan_id` (Number)
- `vxlan_id` (Number) VXLAN id of an OVS overlay network (type `ovsOverlay`; Morpheus `vxlanId`). Accepted on create and never returned, so the configured value is kept in state; a change replaces the network.
- `zone_pool_id` (Number) Zone pool id

### Read-Only

- `config_from_api` (Dynamic) The type-specific configuration as Morpheus returns it, read-only. Includes keys the server adds itself and masks secrets; a change here is visible but never enforced.
- `id` (Number) Network id
- `resource_permissions` (Attributes) (see [below for nested schema](#nestedatt--resource_permissions))

<a id="nestedatt--config_azure"></a>
### Nested Schema for `config_azure`

Optional:

- `resource_group_id` (String, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Resource group the virtual network is created in (Morpheus `resourceGroupId`). Cannot be changed after create.
- `subnet_cidr` (String, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) The subnet's address range in CIDR notation, contained by the virtual network's `cidr` (Morpheus `subnetCidr`). Cannot be changed after create.
- `subnet_name` (String, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Name of a subnet created along with the virtual network; later subnet changes are not reflected here (Morpheus `subnetName`). Cannot be changed after create.


<a id="nestedatt--config_gcp"></a>
### Nested Schema for `config_gcp`

Optional:

- `auto_create` (Boolean, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Whether the VPC creates a subnet in each region automatically (Morpheus `autoCreate`). Cannot be changed after create.
- `mtu` (Number, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Maximum transmission unit, e.g. 1460 (the Google default) or 1500 (Morpheus `mtu`).


<a id="nestedatt--config_hvm_data"></a>
### Nested Schema for `config_hvm_data`

Optional:

- `cluster_network_layout` (String, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) How the network attaches to the cluster: `HOST_INTERFACE` (a VLAN on a host interface; the default) or `VSWITCH` (a virtual switch segment) (Morpheus `hvmClusterNetworkLayout`). Cannot be changed after create.
- `host_network_interface` (String, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Host interface VMs send data traffic out of, e.g. `bond0`, for the `HOST_INTERFACE` layout (Morpheus `hvmHostNetworkInterface`). Cannot be changed after create.
- `mtu` (Number, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Maximum transmission unit: 1500 as the standard default or 9000 for jumbo frames (Morpheus `hvmMtu`).
- `virtual_switch_segment_id` (Number, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Id of a data-eligible virtual switch segment, for the `VSWITCH` layout; the segment determines the VLAN and MTU (Morpheus `virtualSwitchSegmentId`). Cannot be changed after create.


<a id="nestedatt--config_hvm_overlay"></a>
### Nested Schema for `config_hvm_overlay`

Optional:

- `overlay_id` (Number, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) VXLAN network identifier (VNI) between 1 and 16777215, unique across the infrastructure (Morpheus `hvmOverlayId`). Cannot be changed after create.
- `overlay_network_interface` (String, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Host interface that carries the overlay's VXLAN tunnels (Morpheus `hvmOverlayNetworkInterface`). Cannot be changed after create.
- `overlay_port` (Number, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) UDP port for VXLAN tunnel encapsulation between hosts; 4789 is the IANA-assigned port (Morpheus `hvmOverlayPort`). Cannot be changed after create.


<a id="nestedatt--config_hvm_private"></a>
### Nested Schema for `config_hvm_private`

Optional:

- `cluster_network_layout` (String, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) How the network attaches to the cluster: `HOST_INTERFACE` (the default) or `VSWITCH` (Morpheus `hvmClusterNetworkLayout`). Cannot be changed after create.
- `host_network_interface` (String, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Host interface whose MacVTAP device isolates the VMs, for the `HOST_INTERFACE` layout (Morpheus `hvmHostNetworkInterface`). Cannot be changed after create.
- `virtual_switch_id` (Number, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Id of the virtual switch for the private network, for the `VSWITCH` layout (Morpheus `virtualSwitchId`). Cannot be changed after create.


<a id="nestedatt--config_hvm_sriov"></a>
### Nested Schema for `config_hvm_sriov`

Optional:

- `physical_function` (String, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) SR-IOV capable physical function (PF) that virtual functions are allocated from (Morpheus `hvmPhysicalFunction`). Cannot be changed after create.


<a id="nestedatt--config_nsx_cloud"></a>
### Nested Schema for `config_nsx_cloud`

Optional:

- `connected_gateway` (String, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Path of the Tier-1 gateway the segment connects to (Morpheus `connectedGateway`). Supersedes the top-level `connected_gateway`.
- `dhcp_range` (String, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) DHCP address ranges (Morpheus `dhcpRange`).
- `subnet_dhcp_lease_time` (Number, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) DHCP lease time in seconds (Morpheus `subnetDhcpLeaseTime`).
- `subnet_dhcp_server_address` (String, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Address of the segment's DHCP server, in CIDR notation (Morpheus `subnetDhcpServerAddress`).
- `subnet_ip_management_type` (String, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) DHCP type as NSX names it, e.g. `dhcpLocal`, `dhcpRelay` (Morpheus `subnetIpManagementType`).
- `subnet_ip_server_id` (Number, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Id of the DHCP server profile for `dhcpLocal` or `dhcpRelay` (Morpheus `subnetIpServerId`).
- `vlan_ids` (String, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) VLAN ids for a VLAN-backed segment, comma-separated (Morpheus `vlanIDs`).


<a id="nestedatt--config_nsxt"></a>
### Nested Schema for `config_nsxt`

Optional:

- `connected_gateway` (String, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Path of the Tier-0 or Tier-1 gateway the segment connects to, e.g. `/infra/tier-1s/<id>` (Morpheus `connectedGateway`). Supersedes the top-level `connected_gateway`.
- `dhcp_domain_names_ipv6` (String, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Domain names handed out by IPv6 DHCP (Morpheus `dhcpDomainNamesIPv6`).
- `dhcp_excluded_ranges_ipv6` (String, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) IPv6 address ranges excluded from DHCP (Morpheus `dhcpExcludedRangesIPv6`).
- `dhcp_preferred_time_ipv6` (Number, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) IPv6 DHCP preferred time in seconds (Morpheus `dhcpPreferredTimeIPv6`).
- `dhcp_range` (String, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) DHCP address ranges, one per line or comma-separated (Morpheus `dhcpRange`).
- `dhcp_range_ipv6` (String, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) IPv6 DHCP address ranges (Morpheus `dhcpRangeIPv6`).
- `dhcp_sntp_servers_ipv6` (String, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) SNTP servers handed out by IPv6 DHCP (Morpheus `dhcpSntpServersIPv6`).
- `subnet_dhcp_lease_time` (Number, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) DHCP lease time in seconds; NSX-T defaults to 86400 (Morpheus `subnetDhcpLeaseTime`).
- `subnet_dhcp_lease_time_ipv6` (Number, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) IPv6 DHCP lease time in seconds (Morpheus `subnetDhcpLeaseTimeIPv6`).
- `subnet_dhcp_server_address` (String, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Address of the segment's DHCP server, in CIDR notation (Morpheus `subnetDhcpServerAddress`).
- `subnet_dhcp_server_address_ipv6` (String, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) IPv6 address of the segment's DHCP server, in CIDR notation (Morpheus `subnetDhcpServerAddressIPv6`).
- `subnet_ip_management_type` (String, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) How addresses on the segment's subnet are managed: DHCP type as NSX-T names it, e.g. `dhcpLocal`, `dhcpRelay` (Morpheus `subnetIpManagementType`).
- `subnet_ip_server_id` (Number, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Id of the NSX-T DHCP server profile for `dhcpLocal` or `dhcpRelay` (Morpheus `subnetIpServerId`).
- `vlan_ids` (String, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) VLAN ids for a VLAN-backed segment, comma-separated (Morpheus `vlanIDs`).


<a id="nestedatt--config_openstack"></a>
### Nested Schema for `config_openstack`

Optional:

- `dhcp_range` (String, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) DHCP allocation pool for the subnet, as `start,end` (Morpheus `dhcpRange`).
- `host_routes` (String, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Static host routes handed out by the subnet's DHCP (Morpheus `hostRoutes`). Cannot be changed after create.


<a id="nestedatt--config_ovs_port_group"></a>
### Nested Schema for `config_ovs_port_group`

Optional:

- `trunk_only` (Boolean, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Whether the port group carries trunked VLANs only (Morpheus `hvmTrunkOnly`).
- `vlan_trunks` (String, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) VLAN IDs to trunk on the port group, comma-separated (Morpheus `hvmVlanTrunks`).


<a id="nestedatt--resource_permissions"></a>
### Nested Schema for `resource_permissions`

Read-Only:

- `all` (Boolean) Pass true to allow access all groups
- `group_ids` (Set of Number) Array of group (site) IDs that are allowed access

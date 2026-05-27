---
page_title: "hpe_morpheus_network Resource - terraform-provider-hpe"
subcategory: "morpheus"
description: |-
  
---
# hpe_morpheus_network (Resource)



### Networks vs Subnets

For most network types in Morpheus, the network itself is the deployable unit -- there is no separate subnet creation step. Only Azure and Google support a two-level model where subnets can be created independently within a network using the [`hpe_morpheus_subnet`](https://registry.terraform.io/providers/HPE/hpe/latest/docs/resources/morpheus_subnet) resource. For all other types (NSX-T, OpenStack, Amazon, OVS, VCD, ACI, etc.), use this resource to create the full network object including its addressing configuration.

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

## Example Usage

### AWS Network

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
  type_id     = 36
  config = {
    assignPublicIp   = true
    availabilityZone = "us-west-1a"
  }
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

The `config` block for Azure networks (VNets) accepts the following properties:

| Config Key | Required | Description |
|---|---|---|
| `resourceGroupId` | Yes | The ID of the Azure Resource Group to deploy the VNet into |
| `subnetName` | Yes | Name of the default subnet created with the VNet |
| `subnetCidr` | Yes | CIDR for the default subnet (e.g. `192.168.1.0/24`). Must be contained within the VNet's `cidr` address space |

~> **Note** The Azure region is automatically derived from the Cloud zone or Resource Group — it is not specified in the config block.

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
  type_id                    = 35
  cidr                       = "10.100.0.0/16"
  visibility                 = "public"
  active                     = true
  dhcp_server                = true
  appliance_url_proxy_bypass = false
  labels                     = ["terraform", "example"]
  config = {
    "resourceGroupId" = all-attrs-resource-group
    "subnetName"      = "all-attrs-subnet"
    "subnetCidr"      = "10.100.1.0/24"
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
  type_id     = 38
  config = {
    mtu        = "1460"
    autoCreate = true
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

NSX-T segments are created as networks in Morpheus (type code `nsxtLogicalSwitch`). The `config` block accepts the following properties:

| Config Key | Required | Description |
|---|---|---|
| `connectedGateway` | Yes | Connectivity path to the Tier-1 gateway (e.g. `/infra/tier-1s/my-gw`) |
| `vlanIDs` | No | Comma-separated VLAN IDs for VLAN-backed segments. Leave empty for overlay segments |
| `subnetIpManagementType` | No | DHCP type: `dhcpLocal`, `gatewayDhcp`, `dhcpRelay`, or empty for no DHCP |
| `subnetIpServerId` | No | Path to the DHCP server or relay (required when `subnetIpManagementType` is set) |
| `subnetDhcpServerAddress` | No | DHCP server address in CIDR format (e.g. `172.16.10.2/24`). Required for `dhcpLocal` |
| `dhcpRange` | No | DHCP IP range (e.g. `172.16.10.100-172.16.10.200`) |
| `subnetDhcpLeaseTime` | No | DHCP lease time in seconds (default `86400`) |

~> **Note** The top-level `cidr` and `gateway` fields define the segment's subnet addressing. NSX-T subnets are not created separately -- the segment itself is the network object in Morpheus.

```terraform
data "hpe_morpheus_cloud" "nsxt" {
  name = "NSX-T Cloud"
}

data "hpe_morpheus_group" "example" {
  name = "Example Group"
}

resource "hpe_morpheus_network" "nsxt_segment" {
  name         = "example-terraform-nsxt-segment"
  display_name = "Example NSX-T Segment"
  description  = "NSX-T overlay segment managed by Terraform"
  cloud_id     = data.hpe_morpheus_cloud.nsxt.id
  group_id     = data.hpe_morpheus_group.example.id
  type_id      = 7
  cidr         = "172.16.10.0/24"
  gateway      = "172.16.10.1"
  active       = true
  dhcp_server  = true
  config = {
    "connectedGateway"        = "/infra/tier-1s/my-tier1-gw"
    "vlanIDs"                 = ""
    "subnetIpManagementType"  = "dhcpLocal"
    "subnetDhcpServerAddress" = "172.16.10.2/24"
    "dhcpRange"               = "172.16.10.100-172.16.10.200"
    "subnetDhcpLeaseTime"     = "86400"
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

### OVS Port Group Network

```terraform
data "hpe_morpheus_cloud" "example" {
  name = "Morpheus Standard Cloud"
}

data "hpe_morpheus_group" "example" {
  name = "ExampleGroup"
}

data "hpe_morpheus_tenant" "example" {
  name = "Master Tenant"
}

resource "hpe_morpheus_network" "ovs_port_group" {
  name                       = "Terraform OVS Port Group"
  description                = "OVS Port Group network"
  cloud_id                   = data.hpe_morpheus_cloud.example.id
  pool_id                    = 3251
  group_id                   = data.hpe_morpheus_group.example.id
  type_id                    = 63
  switch_id                  = Compute
  config                     = {}
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

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cloud_id` (Number) Cloud (zone) id
- `group_id` (Number) Group (site) id
- `name` (String) Network name.
- `type_id` (Number) Network type id

### Optional

- `active` (Boolean) Activate (true) or disable (false) the network
- `allow_static_override` (Boolean) Allow IP Override
- `appliance_url_proxy_bypass` (Boolean) Bypass Proxy for Appliance URL
- `assign_public_ip` (Boolean) Assign Public IP
- `cidr` (String) Network CIDR.
- `cidr_ipv6` (String) Network IPv6 CIDR.
- `config` (Dynamic) Configuration object. Settings vary by type.
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
- `no_proxy` (String) Comma-separated list of ip addresses or name servers to exclude proxy traversal for. Typically locally routable servers are excluded.
- `pool_id` (Number) Network Pool ID
- `pool_ipv6_id` (Number) IPv6 Network Pool ID
- `search_domains` (String) Search Domains
- `switch_id` (String) Network switch identifier
- `tenant_ids` (Set of Number) List of tenant account ids that are allowed access
- `visibility` (String) Visibility, private or public.
- `vlan_id` (Number)
- `zone_pool_id` (Number) Zone pool id

### Read-Only

- `id` (Number) Network id
- `resource_permissions` (Attributes) (see [below for nested schema](#nestedatt--resource_permissions))

<a id="nestedatt--resource_permissions"></a>
### Nested Schema for `resource_permissions`

Read-Only:

- `all` (Boolean) Pass true to allow access all groups
- `group_ids` (Set of Number) Array of group (site) IDs that are allowed access

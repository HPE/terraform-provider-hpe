---
page_title: "hpe_morpheus_network Resource - terraform-provider-hpe"
subcategory: "morpheus"
description: |-
  
---
# hpe_morpheus_network (Resource)



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
variable "name" {
  description = "Network name"
  type        = string
  default     = "terraform-aws-test"
}

variable "description" {
  description = "Network description"
  type        = string
  default     = "AWS subnet"
}

variable "cloud_id" {
  description = "Cloud (zone) id"
  type        = number
  default     = 207
}

variable "pool_id" {
  description = "Network pool id"
  type        = number
  default     = 1
}

variable "group_id" {
  description = "Group (site) id"
  type        = number
  default     = 1
}

variable "type_id" {
  description = "Network type id"
  type        = number
  default     = 36
}

variable "cidr" {
  description = "CIDR Network"
  type        = string
  default     = "10.200.99.0/24"
}

variable "zone_pool_id" {
  description = "Zone pool id"
  type        = number
  default     = 12329
}

variable "config_assign_public_ip" {
  description = "Assign public IP setting for network config"
  type        = bool
  default     = true
}

variable "config_availability_zone" {
  description = "Availability zone setting for network config"
  type        = string
  default     = "us-west-1a"
}

variable "active" {
  description = "Whether network is active"
  type        = bool
  default     = true
}

variable "dhcp_server" {
  description = "Whether DHCP server is enabled"
  type        = bool
  default     = true
}

variable "appliance_url_proxy_bypass" {
  description = "Whether to bypass proxy for appliance URL"
  type        = bool
  default     = true
}

variable "visibility" {
  description = "Network visibility"
  type        = string
  default     = "private"
}

resource "hpe_morpheus_network" "aws" {
  name        = var.name
  description = var.description
  cloud_id    = var.cloud_id
  pool_id     = var.pool_id
  group_id    = var.group_id
  type_id     = var.type_id
  config = {
    assignPublicIp   = var.config_assign_public_ip
    availabilityZone = var.config_availability_zone
  }
  active                     = var.active
  dhcp_server                = var.dhcp_server
  appliance_url_proxy_bypass = var.appliance_url_proxy_bypass
  tenant_ids                 = [1]
  visibility                 = var.visibility
  cidr                       = var.cidr
  zone_pool_id               = var.zone_pool_id

  lifecycle {
    ignore_changes = [name, display_name, description]
  }
}
```

### Azure Network

```terraform
variable "name" {
  description = "Network name"
  type        = string
  default     = "terraform-network-all-attrs"
}

variable "description" {
  description = "Network description"
  type        = string
  default     = "Network with all attributes set"
}

variable "cloud_id" {
  description = "Cloud (zone) id"
  type        = number
  default     = 4617
}

variable "pool_id" {
  description = "Network pool id"
  type        = number
  default     = 1
}

variable "group_id" {
  description = "Group (site) id"
  type        = number
  default     = 1
}

variable "type_id" {
  description = "Network type id"
  type        = number
  default     = 35
}

variable "cidr" {
  description = "CIDR Network"
  type        = string
  default     = "10.100.0.0/16"
}

variable "visibility" {
  description = "Network visibility"
  type        = string
  default     = "public"
}

variable "active" {
  description = "Whether network is active"
  type        = bool
  default     = true
}

variable "dhcp_server" {
  description = "Whether DHCP server is enabled"
  type        = bool
  default     = true
}

variable "appliance_url_proxy_bypass" {
  description = "Whether to bypass proxy for appliance URL"
  type        = bool
  default     = false
}

variable "config_resource_group_id" {
  description = "Resource Group ID for network config"
  type        = string
  default     = "all-attrs-resource-group"
}

variable "config_subnet_name" {
  description = "Subnet name for network config"
  type        = string
  default     = "all-attrs-subnet"
}

variable "config_subnet_cidr" {
  description = "Subnet CIDR for network config"
  type        = string
  default     = "10.100.1.0/24"
}

variable "config_location" {
  description = "Location for network config"
  type        = string
  default     = "eastus"
}

variable "config_additional_field" {
  description = "Additional config field"
  type        = string
  default     = "test-value"
}

resource "hpe_morpheus_network" "all_attrs" {
  name                       = var.name
  description                = var.description
  cloud_id                   = var.cloud_id
  pool_id                    = var.pool_id
  group_id                   = var.group_id
  type_id                    = var.type_id
  cidr                       = var.cidr
  visibility                 = var.visibility
  active                     = var.active
  dhcp_server                = var.dhcp_server
  appliance_url_proxy_bypass = var.appliance_url_proxy_bypass
  config = {
    "resourceGroupId" = var.config_resource_group_id
    "subnetName"      = var.config_subnet_name
    "subnetCidr"      = var.config_subnet_cidr
    "location"        = var.config_location
    "additionalField" = var.config_additional_field
  }
  tenant_ids = [1, 2, 3]
}
```

### GCP Network

```terraform
variable "name" {
  description = "Network name"
  type        = string
  default     = "TestAccMorpheusNetworkResourceCreateGcp"
}

variable "description" {
  description = "Network description"
  type        = string
  default     = "GCP network"
}

variable "cloud_id" {
  description = "Cloud (zone) id"
  type        = number
  default     = 6
}

variable "pool_id" {
  description = "Network pool id"
  type        = number
  default     = 1
}

variable "group_id" {
  description = "Group (site) id"
  type        = number
  default     = 8
}

variable "type_id" {
  description = "Network type id"
  type        = number
  default     = 38
}

variable "cidr" {
  description = "CIDR Network"
  type        = string
  default     = "10.0.0.0/8"
}

variable "zone_pool_id" {
  description = "Zone pool id"
  type        = number
  default     = 85990
}

variable "config_mtu" {
  description = "MTU setting for network config"
  type        = string
  default     = "1460"
}

variable "config_auto_create" {
  description = "Auto create setting for network config"
  type        = bool
  default     = true
}

variable "active" {
  description = "Whether network is active"
  type        = bool
  default     = true
}

variable "dhcp_server" {
  description = "Whether DHCP server is enabled"
  type        = bool
  default     = false
}

variable "appliance_url_proxy_bypass" {
  description = "Whether to bypass proxy for appliance URL"
  type        = bool
  default     = true
}

variable "visibility" {
  description = "Network visibility"
  type        = string
  default     = "private"
}

resource "hpe_morpheus_network" "gcp" {
  name        = var.name
  description = var.description
  cloud_id    = var.cloud_id
  pool_id     = var.pool_id
  group_id    = var.group_id
  type_id     = var.type_id
  config = {
    mtu        = var.config_mtu
    autoCreate = var.config_auto_create
  }
  active                     = var.active
  dhcp_server                = var.dhcp_server
  appliance_url_proxy_bypass = var.appliance_url_proxy_bypass
  tenant_ids                 = [1]
  visibility                 = var.visibility
  cidr                       = var.cidr
  zone_pool_id               = var.zone_pool_id
}
```

### Host Network

```terraform
variable "name" {
  description = "Network name"
  type        = string
  default     = "terraform-host-network"
}

variable "description" {
  description = "Network description"
  type        = string
  default     = "A test host network"
}

variable "cloud_id" {
  description = "Cloud (zone) id"
  type        = number
  default     = 17
}

variable "pool_id" {
  description = "Network pool id"
  type        = number
  default     = 1
}

variable "group_id" {
  description = "Group (site) id"
  type        = number
  default     = 1
}

variable "type_id" {
  description = "Network type id"
  type        = number
  default     = 1
}

variable "cidr" {
  description = "CIDR Network"
  type        = string
  default     = "10.0.0.0/8"
}

variable "visibility" {
  description = "Network visibility"
  type        = string
  default     = "private"
}

variable "active" {
  description = "Whether network is active"
  type        = bool
  default     = true
}

variable "dhcp_server" {
  description = "Whether DHCP server is enabled"
  type        = bool
  default     = false
}

variable "appliance_url_proxy_bypass" {
  description = "Whether to bypass proxy for appliance URL"
  type        = bool
  default     = true
}

resource "hpe_morpheus_network" "host" {
  name                       = var.name
  description                = var.description
  cloud_id                   = var.cloud_id
  pool_id                    = var.pool_id
  group_id                   = var.group_id
  type_id                    = var.type_id
  config                     = {}
  active                     = var.active
  dhcp_server                = var.dhcp_server
  appliance_url_proxy_bypass = var.appliance_url_proxy_bypass
  tenant_ids                 = [1]
  visibility                 = var.visibility
  cidr                       = var.cidr
}
```

### OVS Port Group Network

```terraform
variable "name" {
  description = "Network name"
  type        = string
  default     = "Terraform OVS Port Group"
}

variable "description" {
  description = "Network description"
  type        = string
  default     = "OVS Port Group network"
}

variable "cloud_id" {
  description = "Cloud (zone) id"
  type        = number
  default     = 7714
}

variable "pool_id" {
  description = "Network pool id"
  type        = number
  default     = 3251
}

variable "group_id" {
  description = "Group (site) id"
  type        = number
  default     = 1
}

variable "type_id" {
  description = "Network type id (OVS Port Group)"
  type        = number
  default     = 63
}

variable "cidr" {
  description = "CIDR Network"
  type        = string
  default     = "10.32.148.0/22"
}

variable "zone_pool_id" {
  description = "Zone pool id"
  type        = number
  default     = 62299
}

variable "vlan_id" {
  description = "VLAN ID"
  type        = number
  default     = 43
}

variable "switch_id" {
  description = "Switch ID for OVS network"
  type        = string
  default     = "Compute"
}

variable "active" {
  description = "Whether network is active"
  type        = bool
  default     = true
}

variable "dhcp_server" {
  description = "Whether DHCP server is enabled"
  type        = bool
  default     = false
}

variable "appliance_url_proxy_bypass" {
  description = "Whether to bypass proxy for appliance URL"
  type        = bool
  default     = true
}

variable "visibility" {
  description = "Network visibility"
  type        = string
  default     = "public"
}

resource "hpe_morpheus_network" "ovs_port_group" {
  name                       = var.name
  description                = var.description
  cloud_id                   = var.cloud_id
  pool_id                    = var.pool_id
  group_id                   = var.group_id
  type_id                    = var.type_id
  switch_id                  = var.switch_id
  config                     = {}
  active                     = var.active
  dhcp_server                = var.dhcp_server
  appliance_url_proxy_bypass = var.appliance_url_proxy_bypass
  tenant_ids                 = [1]
  visibility                 = var.visibility
  cidr                       = var.cidr
  zone_pool_id               = var.zone_pool_id
  vlan_id                    = var.vlan_id

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

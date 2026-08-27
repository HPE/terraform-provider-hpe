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
  labels                     = ["terraform", "acctest", "hpe_morpheus_network", "sweepable"]

  lifecycle {
    ignore_changes = [name, display_name, description]
  }
}

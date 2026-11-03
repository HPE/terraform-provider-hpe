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
  default     = 17
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

variable "network_server_id" {
  description = "The cluster's HVM network server id (GET /api/networks/servers?showHidden=true)"
  type        = number
  default     = 7834
}

variable "vlan_trunks" {
  description = "VLAN ids to trunk on the port group"
  type        = string
  default     = "100,200"
}

variable "cidr" {
  description = "CIDR Network"
  type        = string
  default     = "10.209.101.0/24"
}

variable "zone_pool_id" {
  description = "Zone pool id (the HVM cluster's resource pool)"
  type        = number
  default     = 153047
}

variable "vlan_id" {
  description = "VLAN ID"
  type        = number
  default     = 2101
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
  name              = var.name
  description       = var.description
  cloud_id          = var.cloud_id
  pool_id           = var.pool_id
  group_id          = var.group_id
  type_code         = "ovsPortGroup"
  switch_id         = var.switch_id
  network_server_id = var.network_server_id
  config_ovs_port_group = {
    vlan_trunks = var.vlan_trunks
    trunk_only  = false
  }
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

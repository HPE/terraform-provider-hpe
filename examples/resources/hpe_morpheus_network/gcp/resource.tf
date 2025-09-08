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

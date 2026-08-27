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
  default     = 6446
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
  labels                     = ["terraform", "acctest", "hpe_morpheus_network", "sweepable"]
}

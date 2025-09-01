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

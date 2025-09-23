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
  labels                     = ["terraform", "acctest", "hpe_morpheus_network", "sweepable"]
  config = {
    "resourceGroupId" = var.config_resource_group_id
    "subnetName"      = var.config_subnet_name
    "subnetCidr"      = var.config_subnet_cidr
    "location"        = var.config_location
    "additionalField" = var.config_additional_field
  }
  tenant_ids = [1, 2, 3]
}

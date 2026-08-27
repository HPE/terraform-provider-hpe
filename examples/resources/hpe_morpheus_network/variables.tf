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

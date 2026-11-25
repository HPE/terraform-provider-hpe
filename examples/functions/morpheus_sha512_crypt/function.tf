terraform {
  required_version = ">= 1.8"
  required_providers {
    hpe = {
      source = "HPE/hpe"
    }
  }
}

variable "password" {
  type      = string
  sensitive = true
}

variable "password_salt" {
  type = string
}

locals {
  password_hash = sensitive(provider::hpe::morpheus_sha512_crypt(var.password, var.password_salt, 5000))
}

output "password_hash" {
  value     = local.password_hash
  sensitive = true
}

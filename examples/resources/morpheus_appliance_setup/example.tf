# The password and licence key never reach the plan or the state: they are
# write-only attributes. Read them from a secret store or an ephemeral resource.
variable "admin_password" {
  type      = string
  sensitive = true
  ephemeral = true
}

variable "license_key" {
  type      = string
  sensitive = true
  ephemeral = true
  default   = null
}

# The provider's url points at the new appliance; its credentials are the
# administrator this resource creates, so every other resource depends on it.
resource "hpe_morpheus_appliance_setup" "example" {
  appliance_name = "morpheus.example.com"
  account_name   = "My Company"
  username       = "admin"
  password_wo    = var.admin_password
  email          = "admin@example.com"
  first_name     = "Morpheus"
  last_name      = "Admin"
  license_key_wo = var.license_key

  timeouts = {
    create = "40m"
  }
}

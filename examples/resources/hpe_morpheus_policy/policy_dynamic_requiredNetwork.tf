# Instance Networks Policy - Requires specific networks for instances
resource "hpe_morpheus_policy" "required_networks" {
  name                     = "Instance Networks Policy"
  description              = "Require specific networks for instances"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "requiredNetwork"
  }

  config = {
    requiredNetworks = [100, 200] # Array of required network IDs
  }
}

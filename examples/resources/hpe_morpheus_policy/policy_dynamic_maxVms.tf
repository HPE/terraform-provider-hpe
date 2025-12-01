# Max VMs Policy - Limits VM count
# Allowed associated_resource_types: Group, Cloud, User, Global, Network, Plan
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "max_vms" {
  name                     = "Max VMs Policy"
  description              = "Limit maximum VM count"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "maxVms"
  }

  config = {
    maxVms = "20" # Maximum number of VMs
  }
}

# Max Hosts Policy - Limits host count
# Allowed associated_resource_types: Group, Cloud, User, Global, Plan
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "max_hosts" {
  name                     = "Max Hosts Policy"
  description              = "Limit maximum host count"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "maxHosts"
  }

  config = {
    # Required
    maxHosts = "10" # Maximum number of hosts
  }
}

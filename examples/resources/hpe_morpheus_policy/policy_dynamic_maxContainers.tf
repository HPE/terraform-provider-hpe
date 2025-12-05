# Max Containers Policy - Limits container count
# Allowed associated_resource_types: Group, Cloud, User, Global, Plan
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "max_containers" {
  name                     = "Max Containers Policy"
  description              = "Limit maximum container count"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "maxContainers"
  }

  config = {
    # Required
    maxContainers = "50" # Maximum number of containers
  }
}

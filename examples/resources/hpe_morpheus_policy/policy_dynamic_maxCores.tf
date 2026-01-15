# Max Cores Policy - Limits CPU cores
# Allowed associated_resource_types: Group, Cloud, User, Global, Plan
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "max_cores" {
  name                     = "Max Cores Policy"
  description              = "Limit maximum CPU cores"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "maxCores"
  }

  config = {
    # Required
    maxCores = "32" # Maximum number of CPU cores

    # Optional
    excludeContainers = "off" # Options: "on", "off" - exclude containers from count
  }
}

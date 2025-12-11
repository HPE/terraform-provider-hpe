# Max Storage Policy - Limits storage allocation
# Allowed associated_resource_types: Group, Cloud, User, Global, Plan
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "max_storage" {
  name                     = "Max Storage Policy"
  description              = "Limit maximum storage allocation"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "maxStorage"
  }

  config = {
    # Required
    maxStorage = "1000" # Maximum storage in GB

    # Optional
    excludeContainers = "off" # Options: "on", "off" - exclude containers from count
  }
}

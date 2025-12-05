# Max Load Balancer Pools Policy - Limits load balancer pools
# Allowed associated_resource_types: Cloud, User, Global
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "max_pools" {
  name                     = "Max Load Balancer Pools Policy"
  description              = "Limit maximum load balancer pools"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "maxPools"
  }

  config = {
    # Required
    maxPools = "5" # Maximum number of load balancer pools
  }
}

# Router Quota Policy - Limits router count
resource "hpe_morpheus_policy" "router_quota" {
  name                     = "Router Quota Policy"
  description              = "Limit maximum router count"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "maxRouters"
  }

  config = {
    maxRouters = "5" # Maximum number of routers
  }
}

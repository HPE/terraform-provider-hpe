# Max Pool Members Policy - Limits pool members
resource "hpe_morpheus_policy" "max_pool_members" {
  name                     = "Max Pool Members Policy"
  description              = "Limit maximum pool members"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "maxPoolMembers"
  }

  config = {
    maxPoolMembers = "10" # Maximum number of pool members
  }
}

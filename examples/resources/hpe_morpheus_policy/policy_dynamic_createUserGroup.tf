# User Group Creation Policy - Assigns default user group
# Allowed associated_resource_types: Group, Cloud, User, Global
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "user_group_creation" {
  name                     = "User Group Creation Policy"
  description              = "Assign default user group for created users"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "createUserGroup"
  }

  config = {
    userGroup = "1" # ID of the user group to assign
  }
}

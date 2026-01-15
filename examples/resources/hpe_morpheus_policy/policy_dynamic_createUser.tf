# User Creation Policy - Controls user creation on instances
# Allowed associated_resource_types: Group, Cloud, User, Global
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "user_creation" {
  name                     = "User Creation Policy"
  description              = "Control user creation on provisioned instances"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "createUser"
  }

  config = {
    # Required
    createUserType = "user" # Options: "user" (user configurable), "fixed"

    # Optional
    createUser = true # Enforce user creation
  }
}

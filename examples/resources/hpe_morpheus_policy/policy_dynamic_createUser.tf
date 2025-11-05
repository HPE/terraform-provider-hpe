# User Creation Policy - Controls user creation on instances
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
    createUserType = "user" # Options: "user" (user decides), "off" (no user creation), "on" (required)
    createUser     = true   # Enforce user creation
  }
}

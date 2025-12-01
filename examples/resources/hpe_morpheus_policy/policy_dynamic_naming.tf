# Instance Name Policy - Enforces instance naming conventions
# Allowed associated_resource_types: Group, Cloud, User, Global
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "instance_naming" {
  name                     = "Instance Name Policy"
  description              = "Enforce instance naming conventions"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "naming"
  }

  config = {
    namingType     = "user"                                   # Options: "user" (user configurable), "fixed" (strict pattern)
    namingPattern  = "vm-$${groupCode}-$${type}-$${sequence}" # Naming pattern with variables
    namingConflict = true                                     # Auto-resolve conflicts
  }
}

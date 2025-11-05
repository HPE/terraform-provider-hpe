# Instance Name Policy - Enforces instance naming conventions
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
    namingType     = "user"                                   # Options: "user" (user can customize), "fixed" (strict pattern)
    namingPattern  = "vm-$${groupCode}-$${type}-$${sequence}" # Naming pattern with variables
    namingConflict = true                                     # Allow conflict resolution
  }
}

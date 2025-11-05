# Hostname Policy - Enforces hostname naming conventions
resource "hpe_morpheus_policy" "hostname" {
  name                     = "Hostname Policy"
  description              = "Enforce hostname naming conventions"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "hostNaming"
  }

  config = {
    hostNamingType    = "user"                                     # Options: "user" (user can customize), "fixed" (strict pattern)
    hostNamingPattern = "host-$${groupCode}-$${type}-$${sequence}" # Naming pattern with variables
  }
}

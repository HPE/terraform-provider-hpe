# Hostname Policy - Enforces hostname naming conventions
# Allowed associated_resource_types: Group, Cloud, User, Global
# Tenant specification: allowed (can specify tenants array)
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
    # Required
    hostNamingType = "user" # Options: "user" (user configurable), "fixed" (strict pattern)

    # Optional
    hostNamingPattern = "host-$${groupCode}-$${type}-$${sequence}" # Name pattern uses ${variable} string interpolation. Available variables: groupName, groupCode, cloudName, cloudCode, type, accountId, account, accountType, platform, username, userId, userInitials, provisionType
  }
}

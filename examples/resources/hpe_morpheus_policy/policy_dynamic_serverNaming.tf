# Cluster Resource Name Policy - Enforces naming conventions for cluster resources
# Allowed associated_resource_types: Group, Cloud, User, Global
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "cluster_naming" {
  name                     = "Cluster Resource Naming Policy"
  description              = "Enforce naming for cluster resources"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "serverNaming"
  }

  config = {
    serverNamingType     = "user"                                        # Options: "user" (user configurable), "fixed" (strict pattern)
    serverNamingPattern  = "cluster-$${groupCode}-$${type}-$${sequence}" # Name pattern uses ${variable} string interpolation. Available variables: groupName, groupCode, cloudName, cloudCode, type, accountId, account, accountType, platform, username, userId, userInitials, provisionType
    serverNamingConflict = true                                          # Auto-resolve conflicts
  }
}

# Cluster Resource Name Policy - Enforces naming conventions for cluster resources
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
    serverNamingType     = "user"                                        # Options: "user" (user can customize), "fixed" (strict pattern)
    serverNamingPattern  = "cluster-$${groupCode}-$${type}-$${sequence}" # Naming pattern with variables
    serverNamingConflict = true                                          # Allow conflict resolution
  }
}

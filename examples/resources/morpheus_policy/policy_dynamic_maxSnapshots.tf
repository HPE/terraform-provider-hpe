# Max Snapshots Policy - Limits snapshots per VM
# Allowed associated_resource_types: Group, Cloud, User, Global
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "max_snapshots" {
  name                     = "Max Snapshots Policy"
  description              = "Limit maximum snapshots per VM"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "maxSnapshots"
  }

  config = {
    # Required
    maxSnapshots = "5" # Maximum number of snapshots per VM
  }
}

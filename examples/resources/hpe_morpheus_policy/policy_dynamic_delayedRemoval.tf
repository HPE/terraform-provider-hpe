# Delayed Delete Policy - Delays instance deletion
# Allowed associated_resource_types: Group, Cloud, User, Global
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "delayed_delete" {
  name                     = "Delayed Delete Policy"
  description              = "Delay instance deletion by specified days"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "delayedRemoval"
  }

  config = {
    # Required
    removalAge = "30" # Number of days to delay deletion
  }
}

# Approve Reconfigure Policy
resource "hpe_morpheus_policy" "approve_reconfigure" {
  name                     = "Approve Reconfigure Policy"
  description              = "Require approval before reconfiguring instances"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "reconfigureApproval"
  }

  config = {
    accountIntegrationId = "1" # ID of your ServiceNow or approval integration
  }
}

# Approve Provision Policy
resource "hpe_morpheus_policy" "approve_provision" {
  name                     = "Approve Provision Policy"
  description              = "Require approval before provisioning instances"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "provisionApproval"
  }

  config = {
    accountIntegrationId = "1" # ID of your ServiceNow or approval integration
  }
}

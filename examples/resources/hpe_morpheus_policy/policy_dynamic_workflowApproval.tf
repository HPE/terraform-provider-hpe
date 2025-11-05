# Approve Workflow Execute Policy - Requires ServiceNow integration
resource "hpe_morpheus_policy" "approve_workflow" {
  name                     = "Approve Workflow Execute Policy"
  description              = "Require approval before executing workflows"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "workflowApproval"
  }

  config = {
    accountIntegrationId = "1" # ID of your ServiceNow or approval integration
  }
}

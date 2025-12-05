# Approve Workflow Execute Policy
# Allowed associated_resource_types: Group, Cloud, User, Global, Label
# Tenant specification: allowed (can specify tenants array)
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
    accountIntegrationId = "1"        # ID of your ServiceNow or approval integration
    workflowType         = "workflow" # Options: "workflow" (legacy workflow), "flow" (ServiceNow Flow)
    # workflowId = "123"              # ID of legacy ServiceNow workflow (set if workflowType is 'workflow')
    # flowId = "456"                  # ID of ServiceNow Flow (set if workflowType is 'flow')
  }
}

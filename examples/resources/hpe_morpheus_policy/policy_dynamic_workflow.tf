# Workflow Policy - Executes workflow on provision
# Allowed associated_resource_types: Group, Cloud, User, Global
# Tenant specification: allowed (can specify tenants array)

resource "hpe_morpheus_workflow_operational" "example" {
  name        = "Example Policy Workflow"
  description = "Example workflow for policy testing"
  platform    = "all"
  visibility  = "private"
}

resource "hpe_morpheus_policy" "workflow" {
  name                     = "Workflow Policy"
  description              = "Execute workflow on instance provision"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "workflow"
  }

  config = {
    # Required
    workflowId = hpe_morpheus_workflow_operational.example.id # ID of the workflow to execute
  }
}

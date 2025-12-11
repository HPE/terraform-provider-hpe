# Workflow Policy - Executes workflow on provision
# Allowed associated_resource_types: Group, Cloud, User, Global
# Tenant specification: allowed (can specify tenants array)
# Note: This example uses the morpheus external provider to create a workflow resource
# because the hpe provider does not yet have a workflow resource implemented.
# You will need to configure the morpheus provider in your terraform configuration.

terraform {
  required_providers {
    hpe = {
      source  = "HPE/hpe"
      version = "= 0.3.0"
    }
    morpheus = {
      source  = "gomorpheus/morpheus"
      version = "~> 0.13.2"
    }
  }
}

resource "morpheus_operational_workflow" "example" {
  name        = "Example Policy Workflow"
  description = "Example workflow for policy testing"
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
    workflowId = morpheus_operational_workflow.example.id # ID of the workflow to execute
  }
}

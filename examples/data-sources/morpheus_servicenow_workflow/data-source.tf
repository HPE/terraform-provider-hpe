data "hpe_morpheus_integration" "servicenow_prod" {
  name = "Morpheus Approvals"
}

data "hpe_morpheus_servicenow_workflow" "example" {
  name           = "Morpheus Approvals"
  integration_id = data.morpheus_integration.servicenow_prod.id
}


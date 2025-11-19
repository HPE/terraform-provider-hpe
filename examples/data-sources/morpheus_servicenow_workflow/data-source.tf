data "hpe_morpheus_integration" "servicenow_prod" {
  name = "SNOW Production"
}

data "hpe_morpheus_servicenow_workflow" "morpheus_example" {
  name           = "Morpheus Approvals"
  integration_id = data.morpheus_integration.servicenow_prod.id
}


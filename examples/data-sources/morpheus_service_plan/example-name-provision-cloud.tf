data "hpe_morpheus_service_plan" "example" {
  name                = "Example name"
  provision_type_code = "arm"

  # When plans share a name across clouds/regions (e.g. Azure), set cloud_id to
  # the cloud the plan must be available in to select it unambiguously.
  cloud_id = 5
}

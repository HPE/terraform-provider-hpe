resource "hpe_morpheus_service_plan" "example_service_plan" {
  name                = "ExampleServicePlan"
  code                = "exampleserviceplan"
  max_memory          = 4294967296
  max_storage         = 536870912
  provision_type_code = "arm"
  custom_max_storage  = true
  config_ranges = {
    min_storage = 268435456
    max_storage = 536870912
  }
}

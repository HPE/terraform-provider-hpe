resource "hpe_morpheus_form" "example" {
  name        = "demo"
  code        = "demo"
  description = "demo"
  labels      = ["terraform", "demo"]

  option_type {
    name                     = "tf environment example"
    code                     = "environment-input"
    description              = "Terraform environment example"
    type                     = "environment"
    field_label              = "Environment"
    field_name               = "environment"
    default_value            = "staging"
    placeholder              = ""
    help_block               = "Select an environment"
    required                 = true
    export_meta              = true
    display_value_on_details = true
    locked                   = true
    hidden                   = false
    exclude_from_search      = true
  }
}

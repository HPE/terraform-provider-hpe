resource "hpe_morpheus_form" "example" {
  name        = "demo"
  code        = "demo"
  description = "demo"
  labels      = ["terraform", "demo"]

  option_type {
    name                     = "tf checkbox example"
    code                     = "checkbox-input"
    description              = "Terraform checkbox example"
    type                     = "checkbox"
    field_label              = "checkbox input"
    field_name               = "checkboxInput"
    default_checked          = true
    placeholder              = "Testing 123"
    help_block               = "Help block example"
    required                 = true
    export_meta              = true
    display_value_on_details = true
    locked                   = true
    hidden                   = true
    exclude_from_search      = true
  }
}

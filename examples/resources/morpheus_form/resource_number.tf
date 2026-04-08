resource "hpe_morpheus_form" "example" {
  name        = "demo"
  code        = "demo"
  description = "demo"
  labels      = ["terraform", "demo"]

  option_type {
    name                     = "tf number input example"
    code                     = "number-input"
    description              = "Terraform number example"
    type                     = "number"
    field_label              = "number input"
    field_name               = "numberInput"
    default_value            = "4"
    placeholder              = "Testing 123"
    help_block               = "Is this working now"
    required                 = true
    export_meta              = true
    display_value_on_details = true
    locked                   = true
    hidden                   = true
    exclude_from_search      = true
    min_value                = 3
    max_value                = 44
    step                     = 2
  }
}

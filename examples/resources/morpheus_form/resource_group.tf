resource "hpe_morpheus_form" "example" {
  name        = "demo"
  code        = "demo"
  description = "demo"
  labels      = ["terraform", "demo"]

  option_type {
    name                     = "tf group example"
    code                     = "group-input"
    description              = "Terraform group example"
    type                     = "group"
    field_label              = "group input"
    field_name               = "groupInput"
    default_value            = "test123"
    placeholder              = "Select group"
    help_block               = "Select a group"
    required                 = true
    export_meta              = true
    display_value_on_details = true
    locked                   = true
    hidden                   = false
    exclude_from_search      = true
    allow_read_only          = true
  }
}

resource "hpe_morpheus_form" "example" {
  name        = "demo"
  code        = "demo"
  description = "demo"
  labels      = ["terraform", "demo"]

  option_type {
    name                     = "tf logo selector example"
    code                     = "logo-selector-input"
    description              = "Terraform logo selector example"
    type                     = "logoSelector"
    field_label              = "Select Logo"
    field_name               = "logoSelector"
    default_value            = "identicon"
    placeholder              = ""
    help_block               = "Select or upload a logo"
    required                 = true
    export_meta              = true
    display_value_on_details = true
    locked                   = true
    hidden                   = false
    exclude_from_search      = true
  }
}

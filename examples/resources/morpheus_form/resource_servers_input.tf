resource "hpe_morpheus_form" "example" {
  name        = "demo"
  code        = "demo"
  description = "demo"
  labels      = ["terraform", "demo"]

  option_type {
    name                     = "tf servers-input example"
    code                     = "servers-input"
    description              = "Terraform servers-input example"
    type                     = "servers-input"
    field_label              = "Server"
    field_name               = "server"
    default_value            = ""
    help_block               = "Select a server"
    cloud_field_type         = "value"
    cloud_id                 = 1
    required                 = true
    export_meta              = true
    display_value_on_details = true
    locked                   = true
    hidden                   = false
    exclude_from_search      = true
  }
}

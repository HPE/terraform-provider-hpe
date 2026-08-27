resource "hpe_morpheus_form" "example" {
  name        = "demo"
  code        = "demo"
  description = "demo"
  labels      = ["terraform", "demo"]

  option_type {
    name                     = "tf instances-input example"
    code                     = "instances-input"
    description              = "Terraform instances-input example"
    type                     = "instances-input"
    field_label              = "Instance"
    field_name               = "instance"
    default_value            = ""
    help_block               = "Select an instance"
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

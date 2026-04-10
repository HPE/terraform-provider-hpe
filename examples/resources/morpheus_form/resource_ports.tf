resource "hpe_morpheus_form" "example" {
  name        = "demo"
  code        = "demo"
  description = "demo"
  labels      = ["terraform", "demo"]

  option_type {
    name                     = "tf ports example"
    code                     = "ports-input"
    description              = "Terraform ports example"
    type                     = "ports"
    field_label              = "Exposed Ports"
    field_name               = "ports"
    default_value            = ""
    help_block               = "Configure exposed ports"
    group_field              = "myGroup"
    cloud_field              = "myCloud"
    layout_field             = "myLayout"
    required                 = true
    export_meta              = true
    display_value_on_details = true
    locked                   = true
    hidden                   = false
    exclude_from_search      = true
  }
}

resource "hpe_morpheus_form" "example" {
  name        = "demo"
  code        = "demo"
  description = "demo"
  labels      = ["terraform", "demo"]

  option_type {
    name                     = "tf layout example"
    code                     = "layout-input"
    description              = "Terraform layout example"
    type                     = "layout"
    field_label              = "layout input"
    field_name               = "layoutInput"
    default_value            = ""
    placeholder              = "Select layout"
    help_block               = "Select a layout"
    required                 = true
    export_meta              = true
    display_value_on_details = true
    locked                   = true
    hidden                   = false
    exclude_from_search      = true
    group_field_type         = "value"
    group_id                 = "1"
    cloud_field_type         = "value"
    cloud_id                 = "1"
    instance_type_field_type = "value"
    instance_type_code       = "apache"
  }
}

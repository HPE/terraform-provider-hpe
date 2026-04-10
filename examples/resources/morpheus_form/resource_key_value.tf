resource "hpe_morpheus_form" "example" {
  name        = "demo"
  code        = "demo"
  description = "demo"
  labels      = ["terraform", "demo"]

  option_type {
    name                     = "tf keyValue example"
    code                     = "keyValue-input"
    description              = "Terraform keyValue example"
    type                     = "keyValue"
    field_label              = "KeyValue"
    field_name               = "keyValue"
    default_value            = ""
    help_block               = "Select a key-value pair"
    convert_to_object        = "true"
    key_placeholder          = "Key123"
    value_placeholder        = "Value123"
    required                 = true
    export_meta              = true
    display_value_on_details = true
    locked                   = true
    hidden                   = false
    exclude_from_search      = true
  }
}

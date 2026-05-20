resource "hpe_morpheus_form" "example" {
  name        = "demo"
  code        = "demo"
  description = "demo"
  labels      = ["terraform", "demo"]

  option_type {
    name                     = "tf byteSize example"
    code                     = "bytesize-input"
    description              = "Terraform byteSize example"
    type                     = "byteSize"
    field_label              = "Byte Size"
    field_name               = "byteSize"
    // Size in bytes
    default_value            = "48318382080"
    placeholder              = ""
    help_block               = "Select byte size display"
    required                 = true
    export_meta              = true
    display_value_on_details = true
    locked                   = true
    hidden                   = false
    exclude_from_search      = true
    display                  = "GB"
    lock_display             = false
  }
}

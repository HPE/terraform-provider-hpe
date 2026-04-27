resource "hpe_morpheus_form" "example" {
  name        = "demo"
  code        = "demo"
  description = "demo"
  labels      = ["terraform", "demo"]

  option_type {
    name                     = "tf textarea example"
    code                     = "textarea-input"
    description              = "Terraform textarea example"
    type                     = "textarea"
    field_label              = "Text Area"
    field_name               = "textArea"
    default_value            = "Sample text"
    placeholder              = "Enter text"
    help_block               = "Enter multiple lines of text"
    required                 = true
    export_meta              = true
    display_value_on_details = true
    locked                   = true
    hidden                   = false
    exclude_from_search      = true
    text_rows                = 5
  }
}

resource "hpe_morpheus_form" "example" {
  name        = "demo"
  code        = "demo"
  description = "demo"
  labels      = ["terraform", "demo"]

  option_type {
    name                     = "tf textArray example"
    code                     = "text-array-input"
    description              = "Terraform textArray example"
    type                     = "textArray"
    field_label              = "Text Array"
    field_name               = "textArray"
    default_value            = jsonencode(["item1", "item2", "item3"])
    help_block               = "Enter comma-separated values"
    required                 = true
    export_meta              = true
    display_value_on_details = true
    locked                   = true
    hidden                   = false
    exclude_from_search      = true
    delimiter                = ","
  }
}

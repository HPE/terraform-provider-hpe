resource "hpe_morpheus_form" "example" {
  name        = "demo"
  code        = "demo"
  description = "demo"
  labels      = ["terraform", "demo"]

  option_type {
    name                     = "tf httpheader example"
    code                     = "httpheader-input"
    description              = "Terraform HTTP header input example"
    type                     = "httpHeader"
    field_label              = "HTTP Headers"
    field_name               = "httpHeaders"
    default_value            = jsonencode([{ name = "header1", value = "value1", masked = false }])
    help_block               = "Configure HTTP headers"
    required                 = true
    export_meta              = true
    display_value_on_details = true
    locked                   = true
    hidden                   = false
    exclude_from_search      = true
  }
}

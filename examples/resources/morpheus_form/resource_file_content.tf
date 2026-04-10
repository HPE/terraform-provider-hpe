resource "hpe_morpheus_form" "example" {
  name        = "demo"
  code        = "demo"
  description = "demo"
  labels      = ["terraform", "demo"]

  option_type {
    name                     = "tf fileContent example"
    code                     = "fileContent"
    description              = "Terraform fileContent example"
    type                     = "fileContent"
    field_label              = "FileContent"
    field_name               = "fileContent"
    placeholder              = "testing123"
    help_block               = "Set fileContent"
    required                 = true
    export_meta              = true
    display_value_on_details = true
    locked                   = true
    hidden                   = false
    exclude_from_search      = true
  }
}

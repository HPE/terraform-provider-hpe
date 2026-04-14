resource "hpe_morpheus_form" "example" {
  name        = "demo"
  code        = "demo"
  description = "demo"
  labels      = ["terraform", "demo"]

  option_type {
    name                     = "tf example select"
    code                     = "select-input"
    description              = "Terraform select example"
    type                     = "select"
    field_label              = "Select Test"
    field_name               = "selectTest"
    default_value            = "level1"
    placeholder              = "Testing 123"
    help_block               = "Select an option"
    option_list_id           = 1
    required                 = true
    export_meta              = true
    display_value_on_details = true
    locked                   = true
    hidden                   = true
    exclude_from_search      = true
  }
}

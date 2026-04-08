resource "hpe_morpheus_option_list_manual" "example" {
  name      = "tf example select option list"
  dataset   = <<DATASET
[{"name": "Level 1","value":"level1"},
 {"name": "Level 2","value":"level2"},
 {"name": "Level 3","value":"level3"}]
DATASET
  real_time = true
}

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
    default_value            = "test123"
    placeholder              = "Testing 123"
    help_block               = "Select an option"
    option_list_id           = hpe_morpheus_option_list_manual.example.id
    required                 = true
    export_meta              = true
    display_value_on_details = true
    locked                   = true
    hidden                   = true
    exclude_from_search      = true
  }
}

resource "hpe_morpheus_form" "example" {
  name        = "demo"
  code        = "demo"
  description = "demo"
  labels      = ["terraform", "demo"]

  option_type {
    name                     = "tf typeahead example"
    code                     = "typeahead-input"
    description              = "Terraform typeahead example"
    type                     = "typeahead"
    field_label              = "Typeahead"
    field_name               = "typeahead"
    default_value            = "test"
    placeholder              = "Search..."
    help_block               = "Select an option from the list"
    option_list_id           = 1
    required                 = true
    export_meta              = true
    display_value_on_details = true
    locked                   = true
    hidden                   = false
    exclude_from_search      = true
    sortable                 = true
    allow_duplicates         = false
    custom_data              = "{}"
    allow_multiple_selections = false
  }
}

resource "hpe_morpheus_form" "example" {
  name        = "demo"
  code        = "demo"
  description = "demo"
  labels      = ["terraform", "demo"]

  option_type {
    name                     = "tf tag example"
    code                     = "tag-input"
    description              = "Terraform tag example"
    type                     = "tag"
    field_label              = "Tags"
    field_name               = "tags"
    default_value            = jsonencode([{ name = "Sample Name", value = "Sample Value" }])
    help_block               = "Configure tags"
    group_field_type         = "value"
    group_id                 = 1
    cloud_field_type         = "value"
    cloud_id                 = 1
    required                 = true
    export_meta              = true
    display_value_on_details = true
    locked                   = true
    hidden                   = false
    exclude_from_search      = true
  }
}

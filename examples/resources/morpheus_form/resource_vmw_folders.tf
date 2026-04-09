resource "hpe_morpheus_form" "example" {
  name        = "demo"
  code        = "demo"
  description = "demo"
  labels      = ["terraform", "demo"]

  option_type {
    name                     = "tf vmwFolders example"
    code                     = "vmw-folders-input"
    description              = "Terraform vmwFolders example"
    type                     = "vmwFolders"
    field_label              = "VmwFolders"
    field_name               = "vmwFolders"
    default_value            = ""
    help_block               = "Select a vmwFolder"
    group_field_type         = "value"
    group_id                 = 1
    cloud_field_type         = "value"
    cloud_id                 = 1
    plan_field_type          = "value"
    plan_id                  = 1
    required                 = true
    export_meta              = true
    display_value_on_details = true
    locked                   = true
    hidden                   = false
    exclude_from_search      = true
  }
}

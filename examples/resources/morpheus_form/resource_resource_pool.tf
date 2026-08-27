resource "hpe_morpheus_form" "example" {
  name        = "demo"
  code        = "demo"
  description = "demo"
  labels      = ["terraform", "demo"]

  option_type {
    name                     = "tf resourcePool example"
    code                     = "resource-pool-input"
    description              = "Terraform resourcePool example"
    type                     = "resourcePool"
    field_label              = "Resource Pool"
    field_name               = "resourcePool"
    default_value            = ""
    help_block               = "Select a resource pool"
    group_field_type         = "value"
    group_id                 = 1
    cloud_field_type         = "value"
    cloud_id                 = 1
    plan_field_type          = "value"
    plan_id                  = 1
    layout_field_type        = "value"
    layout_id                = 1
    required                 = true
    export_meta              = true
    display_value_on_details = true
    locked                   = true
    hidden                   = false
    exclude_from_search      = true
  }
}

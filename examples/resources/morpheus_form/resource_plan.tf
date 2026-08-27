resource "hpe_morpheus_form" "example" {
  name        = "demo"
  code        = "demo"
  description = "demo"
  labels      = ["terraform", "demo"]

  option_type {
    name                     = "tf plan example"
    code                     = "plan-input"
    description              = "Terraform plan example"
    type                     = "plan"
    field_label              = "plan input"
    field_name               = "planInput"
    default_value            = jsonencode({ id = 1088, maxMemory = 8589934592, maxCores = "4", coresPerSocket = "2" })
    placeholder              = "Select plan"
    help_block               = "Select a plan"
    required                 = true
    export_meta              = true
    display_value_on_details = true
    locked                   = true
    hidden                   = false
    exclude_from_search      = true
    show_pricing             = false
    group_field_type         = "value"
    group_id                 = "1"
    cloud_field_type         = "value"
    cloud_id                 = "1"
    layout_field_type        = "value"
    layout_id                = "1"
    pool_field_type          = "value"
    pool_id                  = "1"
  }
}

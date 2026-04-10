resource "hpe_morpheus_form" "example" {
  name        = "demo"
  code        = "demo"
  description = "demo"
  labels      = ["terraform", "demo"]

  option_type {
    name                     = "tf secGroup example"
    code                     = "sec-group-input"
    description              = "Terraform secGroup example"
    type                     = "secGroup"
    field_label              = "Security Groups"
    field_name               = "securityGroups"
    default_value            = ""
    help_block               = "Select security groups"
    cloud_field_type         = "value"
    cloud_id                 = 1
    pool_field               = "resourcePool"
    required                 = true
    export_meta              = true
    display_value_on_details = true
    locked                   = true
    hidden                   = false
    exclude_from_search      = true
  }
}

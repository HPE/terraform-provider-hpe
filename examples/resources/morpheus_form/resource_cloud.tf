resource "hpe_morpheus_form" "example" {
  name        = "demo"
  code        = "demo"
  description = "demo"
  labels      = ["terraform", "demo"]

  option_type {
    name                         = "tf cloud example"
    code                         = "cloud-input"
    description                  = "Terraform cloud example"
    type                         = "cloud"
    field_label                  = "cloud input"
    field_name                   = "cloudInput"
    default_value                = "test123"
    placeholder                  = "Select cloud"
    help_block                   = "Select a cloud"
    required                     = true
    export_meta                  = true
    display_value_on_details     = true
    locked                       = true
    hidden                       = false
    exclude_from_search          = true
    filter_from_resource         = true
    group_field_type             = "value"
    group_id                     = "1"
    instance_type_field_type     = "value"
    instance_type_code           = "apache"
    cloud_type                   = "4"
  }
}

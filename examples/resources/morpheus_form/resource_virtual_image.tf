resource "hpe_morpheus_form" "example" {
  name        = "demo"
  code        = "demo"
  description = "demo"
  labels      = ["terraform", "demo"]

  option_type {
    name                           = "tf virtual-image example"
    code                           = "virtual-image"
    description                    = "Terraform virtual-image example"
    type                           = "virtual-image"
    field_label                    = "Virtual Image"
    field_name                     = "virtual-image"
    default_value                  = ""
    help_block                     = "Select a virtual image"
    virtual_image_cloud_field_type = "id"
    virtual_image_cloud_id         = 1
    required                       = true
    export_meta                    = true
    display_value_on_details       = true
    locked                         = true
    hidden                         = false
    exclude_from_search            = true
  }
}

resource "hpe_morpheus_form" "example" {
  name        = "demo"
  code        = "demo"
  description = "demo"
  labels      = ["terraform", "demo"]

  option_type {
    name                        = "tf network manager example"
    code                        = "network-manager-input"
    description                 = "Terraform network manager example"
    type                        = "networkManager"
    field_label                 = "network input"
    field_name                  = "networkInput"
    default_value               = "test123"
    placeholder                 = "Select network"
    help_block                  = "Select a network"
    required                    = true
    export_meta                 = true
    display_value_on_details    = true
    locked                      = true
    hidden                      = false
    exclude_from_search         = true
    show_network_type_selection = true
    enable_ip_mode_selection    = true
    group_field_type            = "value"
    group_id                    = "1"
    cloud_field_type            = "value"
    cloud_id                    = "1"
    pool_field_type             = "value"
    pool_id                     = "1"
    layout_field_type           = "value"
    layout_id                   = "1"
  }
}

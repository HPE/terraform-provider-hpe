resource "hpe_morpheus_form" "example" {
  name        = "demo-cascade"
  code        = "demo-cascade"
  description = "Demonstrates group→cloud cascade"
  labels      = ["terraform", "demo"]

  option_type {
    name        = "Group Selector"
    code        = "group-selector"
    type        = "group"
    field_label = "Group"
    field_name  = "fGroups"
  }

  option_type {
    name             = "Cloud Selector"
    code             = "cloud-selector"
    type             = "cloud"
    field_label      = "Cloud"
    field_name       = "fClouds"
    group_field_type = "field"
    group_field      = "fGroups"
  }
}

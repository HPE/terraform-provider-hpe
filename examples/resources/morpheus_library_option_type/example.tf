resource "hpe_morpheus_library_option_type" "example" {
  name          = "Environment Selector"
  field_name    = "environment"
  type          = "select"
  field_label   = "Environment"
  default_value = "development"
  required      = true
}

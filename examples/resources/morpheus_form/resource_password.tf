resource "hpe_morpheus_form" "example" {
  name        = "demo"
  code        = "demo"
  description = "demo"
  labels      = ["terraform", "demo"]

  option_type {
    name                     = "tf password example"
    code                     = "password-input"
    description              = "Terraform password example"
    type                     = "password"
    field_label              = "Password"
    field_name               = "password"
    default_value            = ""
    placeholder              = "Enter password"
    help_block               = "Enter a secure password"
    required                 = true
    export_meta              = true
    display_value_on_details = true
    locked                   = true
    hidden                   = false
    exclude_from_search      = true
    allow_password_peek      = true
  }
}

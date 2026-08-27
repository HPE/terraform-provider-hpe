resource "hpe_morpheus_form" "example" {
  name        = "demo"
  code        = "demo"
  description = "demo"
  labels      = ["terraform", "demo"]

  field_group {
    name                 = "fg1"
    description          = "testin"
    collapsible          = true
    collapsed_by_default = true
    option_type {
      name                     = "tf field group 1 text input example"
      code                     = "test-input-1"
      description              = "Terraform text input example"
      type                     = "text"
      field_label              = "Testing 1"
      field_name               = "test1"
      default_value            = "Demo123"
      placeholder              = "Testing 123"
      help_block               = "Help block example"
      required                 = true
      export_meta              = true
      display_value_on_details = true
      locked                   = true
      hidden                   = false
      exclude_from_search      = true
    }
  }

  field_group {
    name                 = "fg2"
    description          = "testin"
    collapsible          = true
    collapsed_by_default = true
    option_type {
      name                     = "tf field group 2 text input example"
      code                     = "test-input-2"
      description              = "Terraform text input example"
      type                     = "text"
      field_label              = "Testing 2"
      field_name               = "test2"
      default_value            = "Demo123"
      placeholder              = "Testing 123"
      help_block               = "Help block example"
      required                 = true
      export_meta              = true
      display_value_on_details = true
      locked                   = true
      hidden                   = false
      exclude_from_search      = true
    }
  }
}

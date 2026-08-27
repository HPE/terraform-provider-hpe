resource "hpe_morpheus_form" "example" {
  name        = "demo"
  code        = "demo"
  description = "demo"
  labels      = ["terraform", "demo"]

  option_type {
    name                     = "tf code-editor example"
    code                     = "code-editor-input"
    description              = "Terraform code-editor example"
    type                     = "code-editor"
    field_label              = "Code Editor"
    field_name               = "codeEditor"
    default_value            = "echo hello world"
    placeholder              = ""
    help_block               = "Enter code"
    required                 = true
    export_meta              = true
    display_value_on_details = true
    locked                   = true
    hidden                   = false
    exclude_from_search      = true
    show_line_numbers        = true
    code_language            = "bash"
  }
}

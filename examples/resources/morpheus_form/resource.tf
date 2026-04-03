resource "hpe_morpheus_form" "example" {
  name        = "demo"
  code        = "demo"
  description = "demo"
  labels      = ["terraform", "demo"]

  option_type {
    name                     = "tf example select"
    code                     = "select-input"
    description              = "Terraform select example"
    type                     = "select"
    field_label              = "Select Test"
    field_name               = "selectTest"
    default_value            = "test123"
    placeholder              = "Testing 123"
    help_block               = "Select an option"
    option_list_id           = 1
    required                 = true
    export_meta              = true
    display_value_on_details = true
    locked                   = true
    hidden                   = true
    exclude_from_search      = true
  }

  option_type {
    name                     = "tf radio example"
    code                     = "radio-input"
    description              = "Terraform radio example"
    type                     = "radio"
    field_label              = "Radio Test"
    field_name               = "radioTest"
    default_value            = "Demo123"
    placeholder              = "Testing 123"
    help_block               = "Select an option"
    option_list_id           = 1
    required                 = true
    export_meta              = true
    display_value_on_details = true
    locked                   = true
    hidden                   = true
    exclude_from_search      = true
  }

  option_type {
    name                     = "tf text example"
    code                     = "test-input"
    description              = "Terraform text example"
    type                     = "text"
    field_label              = "Testin"
    field_name               = "test"
    default_value            = "Demo123"
    placeholder              = "Testing 123"
    help_block               = "Is this working now"
    required                 = true
    export_meta              = true
    display_value_on_details = true
    locked                   = true
    hidden                   = true
    exclude_from_search      = true
  }

  option_type {
    name                     = "tf checkbox example"
    code                     = "checkbox-input"
    description              = "Terraform checkbox example"
    type                     = "checkbox"
    field_label              = "checkbox input"
    field_name               = "checkboxInput"
    default_checked          = true
    placeholder              = "Testing 123"
    help_block               = "Is this working now"
    required                 = true
    export_meta              = true
    display_value_on_details = true
    locked                   = true
    hidden                   = true
    exclude_from_search      = true
  }

  option_type {
    name                     = "tf hidden input example"
    code                     = "hidden-input"
    description              = "Terraform hidden input example"
    type                     = "hidden"
    field_label              = "hidden input"
    field_name               = "hiddenInput"
    default_value            = "test"
    placeholder              = "Testing 123"
    help_block               = "Is this working now"
    required                 = true
    export_meta              = true
    display_value_on_details = true
    locked                   = true
    hidden                   = true
    exclude_from_search      = true
  }

  option_type {
    name                     = "tf number input example"
    code                     = "number-input"
    description              = "Terraform number example"
    type                     = "number"
    field_label              = "number input"
    field_name               = "numberInput"
    default_value            = "4"
    placeholder              = "Testing 123"
    help_block               = "Is this working now"
    required                 = true
    export_meta              = true
    display_value_on_details = true
    locked                   = true
    hidden                   = true
    exclude_from_search      = true
    min_value                = 3
    max_value                = 44
    step                     = 2
  }

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
  }

  option_type {
    name                     = "tf cloud example"
    code                     = "cloud-input"
    description              = "Terraform cloud example"
    type                     = "cloud"
    field_label              = "cloud input"
    field_name               = "cloudInput"
    default_value            = "test123"
    placeholder              = "Select cloud"
    help_block               = "Select a cloud"
    required                 = true
    export_meta              = true
    display_value_on_details = true
    locked                   = true
    hidden                   = false
    exclude_from_search      = true
    filter_from_resource     = true
  }

  option_type {
    name                     = "tf layout example"
    code                     = "layout-input"
    description              = "Terraform layout example"
    type                     = "layout"
    field_label              = "layout input"
    field_name               = "layoutInput"
    default_value            = ""
    placeholder              = "Select layout"
    help_block               = "Select a layout"
    required                 = true
    export_meta              = true
    display_value_on_details = true
    locked                   = true
    hidden                   = false
    exclude_from_search      = true
    group_field_type         = "value"
    group_id                 = "1"
    cloud_field_type         = "value"
    cloud_id                 = "1"
    instance_type_field_type = "value"
    instance_type_code       = "apache"
  }

  option_type {
    name                     = "tf group example"
    code                     = "group-input"
    description              = "Terraform group example"
    type                     = "group"
    field_label              = "group input"
    field_name               = "groupInput"
    default_value            = "test123"
    placeholder              = "Select group"
    help_block               = "Select a group"
    required                 = true
    export_meta              = true
    display_value_on_details = true
    locked                   = true
    hidden                   = false
    exclude_from_search      = true
    allow_read_only          = true
  }

  field_group {
    name                 = "fg1"
    description          = "testin"
    collapsible          = true
    collapsed_by_default = true
    option_type {
      name                     = "tf field group 1 text input example"
      code                     = "test-input"
      description              = "Terraform text input example"
      type                     = "text"
      field_label              = "Testin"
      field_name               = "test"
      default_value            = "Demo123"
      placeholder              = "Testing 123"
      help_block               = "Is this working now"
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
      code                     = "test-input"
      description              = "Terraform text input example"
      type                     = "text"
      field_label              = "Testin"
      field_name               = "test"
      default_value            = "Demo123"
      placeholder              = "Testing 123"
      help_block               = "Is this working now"
      required                 = true
      export_meta              = true
      display_value_on_details = true
      locked                   = true
      hidden                   = false
      exclude_from_search      = true
    }
  }
}

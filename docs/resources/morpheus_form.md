---
page_title: "hpe_morpheus_form Resource - terraform-provider-hpe"
subcategory: "morpheus"
description: |-
  Provides a Morpheus form resource
---
# hpe_morpheus_form (Resource)

Provides a Morpheus form resource

!> **Note:** Existing inputs or option types are not supported 
and all inputs or option types must be defined in the form.

## Example Usage

```terraform
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
    group_field_type            = "value"
    group_id                    = "1"
    cloud_field_type            = "value"
    cloud_id                    = "1"
    pool_field_type             = "value"
    pool_id                     = "1"
    layout_field_type           = "value"
    layout_id                   = "1"
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `code` (String) The form code used for API/CLI automation
- `name` (String) The name of the form

### Optional

- `description` (String) A description of the form
- `field_group` (Block List) Field group to add to the form (see [below for nested schema](#nestedblock--field_group))
- `labels` (Set of String) The organization labels associated with the form
- `option_type` (Block List) Form option type (see [below for nested schema](#nestedblock--option_type))

### Read-Only

- `id` (String) The id of the form

<a id="nestedblock--field_group"></a>
### Nested Schema for `field_group`

Required:

- `name` (String) The name of the field group

Optional:

- `collapsed_by_default` (Boolean) Whether the field group is collapsed by default
- `collapsible` (Boolean) Whether the field group can be collapsed
- `description` (String) A description of the field group
- `option_type` (Block List) Field Group option type (see [below for nested schema](#nestedblock--field_group--option_type))
- `visibility_field` (String) The field or code used to trigger the visibility of the field group

<a id="nestedblock--field_group--option_type"></a>
### Nested Schema for `field_group.option_type`

Optional:

- `allow_duplicates` (Boolean) Whether duplicate selections are allowed
- `allow_multiple_selections` (Boolean) Whether to allow multiple items to be selected when using a select list or type ahead option type
- `allow_password_peek` (Boolean) Whether the value of the password option type can be revealed by the user to ensure they correctly entered the password
- `allow_read_only` (Boolean) Whether to allow read only instances of this type
- `cloud_field` (String) The field code used to determine the cloud for an option type
- `cloud_field_type` (String) How the cloud is specified for an option type (field or value)
- `cloud_id` (String) The cloud ID to filter layouts by for an option type
- `code` (String) The code of the option type to add to the field group
- `code_language` (String) The coding language used for highlighting code syntax
- `custom_data` (String) Custom JSON data payload to pass (Must be a JSON string)
- `default_checked` (Boolean) Whether the checkbox option type is checked by default
- `default_value` (String) The default value of the option type
- `delimiter` (String) The delimiter used to separate text array input values
- `dependent_field` (String) The field or code used to trigger the reloading of the field
- `description` (String) A description of the option type to add to the field group
- `display` (String) The memory or storage value to use (GB or MB)
- `display_value_on_details` (Boolean) Display the selected value of the option type on the associated resource's details page
- `enable_ip_mode_selection` (Boolean) Whether to enable IP Mode Selection
- `exclude_from_search` (Boolean) Whether the option type should be excluded from search or not
- `export_meta` (Boolean) Whether to export the option type as a tag
- `field_label` (String) The label of the option type
- `field_name` (String) The field name of the option type to add to the field group
- `filter_from_resource` (Boolean) Whether to filter out resources that are not associated with this option
- `group_field` (String) The field code used to determine the group for an option type
- `group_field_type` (String) How the group is specified for an option type (field or value)
- `group_id` (String) The group ID to filter layouts by for an option type
- `help_block` (String) The help block text for the option type
- `hidden` (Boolean) Whether the option type is hidden or not
- `instance_type_code` (String) The instance type code to filter layouts by for a layout option type
- `instance_type_field_code` (String) The field code used to determine the instance type for a layout option type
- `instance_type_field_type` (String) How the instance type is specified for a layout option type (field or value)
- `layout_field` (String) The field code used to determine the layout for an option type
- `layout_field_type` (String) How the layout is specified for an option type (field or value)
- `layout_id` (String) The layout ID to filter by for an option type
- `lock_display` (Boolean) Whether to lock the display or not
- `locked` (Boolean) Whether the option type is locked or not
- `max_value` (Number) The maximum value that can be provided for a number option type
- `min_value` (Number) The minimum number that can be selected for a number option type
- `name` (String) The name of the option type to add to the field group
- `option_list_id` (Number) The id of the option list for option types such as a typeahead or select list
- `placeholder` (String) The placeholder text for the option type
- `pool_field` (String) The field code used to determine the pool for a networkManager option type
- `pool_field_type` (String) How the pool is specified for a networkManager option type (field or value)
- `pool_id` (String) The pool ID to filter by for a networkManager option type
- `remove_select_option` (Boolean) For Select List-type Inputs. When marked, the Input will default to the first item in the list rather than to an empty selection
- `require_field` (String) The field or code used to determine whether the field is required or not
- `required` (Boolean) Whether the option type is required or not
- `show_line_numbers` (Boolean) Whether to show the line numbers for the code editor option type
- `show_network_type_selection` (Boolean) Whether to show the network type selection
- `sortable` (Boolean) Whether the selected options can be sorted or not
- `step` (Number) The incrementation number used for the number option type (i.e. - 5s, 10s, 100s, etc.)
- `text_rows` (Number) The number of rows to display for a text area or code editor option type
- `type` (String) The type of option type to add to the field group (byteSize, checkbox, cloud, code-editor, group, hidden, layout, networkManager, number, password, radio, select, text, textarea, textArray, typeahead)
- `verify_pattern` (String) The regex pattern used to validate the entered text
- `visibility_field` (String) The field or code used to trigger the visibility of the field



<a id="nestedblock--option_type"></a>
### Nested Schema for `option_type`

Optional:

- `allow_duplicates` (Boolean) Whether duplicate selections are allowed
- `allow_multiple_selections` (Boolean) Whether to allow multiple items to be selected when using a select list or type ahead option type
- `allow_password_peek` (Boolean) Whether the value of the password option type can be revealed by the user to ensure they correctly entered the password
- `allow_read_only` (Boolean) Whether to allow read only instances of this type
- `cloud_field` (String) The field code used to determine the cloud for an option type
- `cloud_field_type` (String) How the cloud is specified for an option type (field or value)
- `cloud_id` (String) The cloud ID to filter layouts by for an option type
- `code` (String) The code of the option type to add to the form
- `code_language` (String) The coding language used for highlighting code syntax
- `custom_data` (String) Custom JSON data payload to pass (Must be a JSON string)
- `default_checked` (Boolean) Whether the checkbox option type is checked by default
- `default_value` (String) The default value of the option type
- `delimiter` (String) The delimiter used to separate text array input values
- `dependent_field` (String) The field or code used to trigger the reloading of the field
- `description` (String) A description of the option type to add to the form
- `display` (String) The memory or storage value to use (GB or MB)
- `display_value_on_details` (Boolean) Display the selected value of the option type on the associated resource's details page
- `enable_ip_mode_selection` (Boolean) Whether to enable IP Mode Selection
- `exclude_from_search` (Boolean) Whether the option type should be excluded from search or not
- `export_meta` (Boolean) Whether to export the option type as a tag
- `field_label` (String) The label of the option type
- `field_name` (String) The field name of the option type to add to the form
- `filter_from_resource` (Boolean) Whether to filter out resources that are not associated with this option
- `group_field` (String) The field code used to determine the group for an option type
- `group_field_type` (String) How the group is specified for an option type (field or value)
- `group_id` (String) The group ID to filter layouts by for an option type
- `help_block` (String) The help block text for the option type
- `hidden` (Boolean) Whether the option type is hidden or not
- `instance_type_code` (String) The instance type code to filter layouts by for a layout option type
- `instance_type_field_code` (String) The field code used to determine the instance type for a layout option type
- `instance_type_field_type` (String) How the instance type is specified for a layout option type (field or value)
- `layout_field` (String) The field code used to determine the layout for an option type
- `layout_field_type` (String) How the layout is specified for an option type (field or value)
- `layout_id` (String) The layout ID to filter by for an option type
- `lock_display` (Boolean) Whether to lock the display or not
- `locked` (Boolean) Whether the option type is locked or not
- `max_value` (Number) The maximum value that can be provided for a number option type
- `min_value` (Number) The minimum number that can be selected for a number option type
- `name` (String) The name of the option type to add to the form
- `option_list_id` (Number) The id of the option list for option types such as a typeahead or select list
- `placeholder` (String) The placeholder text for the option type
- `pool_field` (String) The field code used to determine the pool for a networkManager option type
- `pool_field_type` (String) How the pool is specified for a networkManager option type (field or value)
- `pool_id` (String) The pool ID to filter by for a networkManager option type
- `remove_select_option` (Boolean) For Select List-type Inputs. When marked, the Input will default to the first item in the list rather than to an empty selection
- `require_field` (String) The field or code used to determine whether the field is required or not
- `required` (Boolean) Whether the option type is required or not
- `show_line_numbers` (Boolean) Whether to show the line numbers for the code editor option type
- `show_network_type_selection` (Boolean) Whether to show the network type selection
- `sortable` (Boolean) Whether the selected options can be sorted or not
- `step` (Number) The incrementation number used for the number option type (i.e. - 5s, 10s, 100s, etc.)
- `text_rows` (Number) The number of rows to display for a text area or code editor option type
- `type` (String) The type of option type to add to the form (byteSize, checkbox, cloud, code-editor, group, hidden, layout, networkManager, number, password, radio, select, text, textarea, textArray, typeahead)
- `verify_pattern` (String) The regex pattern used to validate the entered text
- `visibility_field` (String) The field or code used to trigger the visibility of the field

## Import

Import is supported using the following syntax:

```shell
terraform import hpe_morpheus_form.tf_example_form 1
```

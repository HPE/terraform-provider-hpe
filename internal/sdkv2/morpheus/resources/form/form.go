// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package form

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"
)

const (
	typeByteSize   = "byteSize"
	typeCheckbox   = "checkbox"
	typeCodeEditor = "code-editor"
	typeHidden     = "hidden"
	typeNumber     = "number"
	typePassword   = "password"
	typeRadio      = "radio"
	typeSelect     = "select"
	typeText       = "text"
	typeTextArea   = "textarea"
	typeTextArray  = "textArray"
	typeTypeahead  = "typeahead"

	valueFalse = "false"
	valueTrue  = "true"
)

func ResourceForm() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus form resource",
		CreateContext: resourceFormCreate,
		ReadContext:   resourceFormRead,
		UpdateContext: resourceFormUpdate,
		DeleteContext: resourceFormDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The id of the form",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the form",
				Required:    true,
			},
			"code": {
				Type:        schema.TypeString,
				Description: "The form code used for API/CLI automation",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "A description of the form",
				Optional:    true,
				Computed:    true,
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "The organization labels associated with the form",
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"option_type": {
				Type:        schema.TypeList,
				Description: "Form option type",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": {
							Type:        schema.TypeString,
							Description: "The code of the option type to add to the form",
							Optional:    true,
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "The name of the option type to add to the form",
							Optional:    true,
						},
						"description": {
							Type:        schema.TypeString,
							Description: "A description of the option type to add to the form",
							Optional:    true,
						},
						"field_name": {
							Type:        schema.TypeString,
							Description: "The name of the option type field to add to the form",
							Optional:    true,
						},
						"type": {
							Type: schema.TypeString,
							Description: "The type of option type to add to the form " +
								"(checkbox, hidden, number, password, radio, select, text, textarea, byteSize, " +
								"code-editor, fileContent, logoSelector, textArray, typeahead, environment)",
							ValidateFunc: validation.StringInSlice(
								[]string{
									"checkbox", "hidden", "number", "password", "radio", "select", "text",
									"textarea", "byteSize", "code-editor", "fileContent", "logoSelector",
									"textArray", "typeahead", "environment",
								},
								false,
							),
							Optional: true,
						},
						"option_list_id": {
							Type:        schema.TypeInt,
							Description: "The id of the option list for option types such as a typeahead or select list",
							Optional:    true,
							Computed:    true,
						},
						"field_label": {
							Type:        schema.TypeString,
							Description: "The label used for the option type",
							Optional:    true,
							Computed:    true,
						},
						"default_value": {
							Type:        schema.TypeString,
							Description: "The default value of the option type",
							Optional:    true,
							Computed:    true,
						},
						"default_checked": {
							Type:        schema.TypeBool,
							Description: "Whether the checkbox option type is checked by default",
							Optional:    true,
							Computed:    true,
						},
						"placeholder": {
							Type:        schema.TypeString,
							Description: "The placeholder text used for the option type",
							Optional:    true,
							Computed:    true,
						},
						"help_block": {
							Type:        schema.TypeString,
							Description: "The help message displayed below the option type",
							Optional:    true,
							Computed:    true,
						},
						"required": {
							Type:        schema.TypeBool,
							Description: "Whether the option type is required or not",
							Optional:    true,
							Computed:    true,
						},
						"export_meta": {
							Type:        schema.TypeBool,
							Description: "Whether to export the option type as a tag",
							Optional:    true,
							Computed:    true,
						},
						"display_value_on_details": {
							Type:        schema.TypeBool,
							Description: "Display the selected value of the option type on the associated resource's details page",
							Optional:    true,
							Computed:    true,
						},
						"locked": {
							Type:        schema.TypeBool,
							Description: "Whether the option type is locked or not",
							Optional:    true,
							Computed:    true,
						},
						"hidden": {
							Type:        schema.TypeBool,
							Description: "Whether to display the option type to the user",
							Optional:    true,
							Computed:    true,
						},
						"exclude_from_search": {
							Type:        schema.TypeBool,
							Description: "Whether the option type should be execluded from search or not",
							Optional:    true,
							Computed:    true,
						},
						"allow_password_peek": {
							Type: schema.TypeBool,
							Description: "Whether the value of the password option type can be revealed by " +
								"the user to ensure they correctly entered the password",
							Optional: true,
							Computed: true,
						},
						"min_value": {
							Type:        schema.TypeInt,
							Description: "The minimum number that can be selected for a number option type",
							Optional:    true,
							Computed:    true,
						},
						"max_value": {
							Type:        schema.TypeInt,
							Description: "The maximum value that can be provided for a number option type",
							Optional:    true,
							Computed:    true,
						},
						"step": {
							Type:        schema.TypeInt,
							Description: "The incrementation number used for the number option type (i.e. - 5s, 10s, 100s, etc.)",
							Optional:    true,
							Computed:    true,
						},
						"text_rows": {
							Type:        schema.TypeInt,
							Description: "The number of lines to show for a code editor or text area option type",
							Optional:    true,
							Computed:    true,
						},
						"display": {
							Type:         schema.TypeString,
							Description:  "The memory or storage value to use (GB or MB)",
							ValidateFunc: validation.StringInSlice([]string{"GB", "MB"}, false),
							Optional:     true,
							Computed:     true,
						},
						"lock_display": {
							Type:        schema.TypeBool,
							Description: "Whether to lock the display or not",
							Optional:    true,
							Computed:    true,
						},
						"code_language": {
							Type:        schema.TypeString,
							Description: "The coding language used for highlighting code syntax",
							Optional:    true,
							Computed:    true,
						},
						"show_line_numbers": {
							Type:        schema.TypeBool,
							Description: "Whether to show the line numbers for the code editor option type",
							Optional:    true,
							Computed:    true,
						},
						"sortable": {
							Type:        schema.TypeBool,
							Description: "Whether the selected options can be sorted or not",
							Optional:    true,
							Computed:    true,
						},
						"allow_multiple_selections": {
							Type: schema.TypeBool,
							Description: "Whether to allow multiple items to be selected when using a " +
								"select list or type ahead option type",
							Optional: true,
							Computed: true,
						},
						"remove_select_option": {
							Type: schema.TypeBool,
							Description: "For Select List-type Inputs. When marked, the Input will default to " +
								"the first item in the list rather than to an empty selection",
							Optional: true,
							Computed: true,
						},
						"allow_duplicates": {
							Type:        schema.TypeBool,
							Description: "Whether duplicate selections are allowed",
							Optional:    true,
							Computed:    true,
						},
						"custom_data": {
							Type:        schema.TypeString,
							Description: "Custom JSON data payload to pass (Must be a JSON string)",
							Optional:    true,
							Computed:    true,
						},
						"dependent_field": {
							Type:        schema.TypeString,
							Description: "The field or code used to trigger the reloading of the field",
							Optional:    true,
							Computed:    true,
						},
						"delimiter": {
							Type:        schema.TypeString,
							Description: "The delimiter used to separate text array input values",
							Optional:    true,
							Computed:    true,
						},
						"visibility_field": {
							Type:        schema.TypeString,
							Description: "The field or code used to trigger the visibility of the field",
							Optional:    true,
							Computed:    true,
						},
						"verify_pattern": {
							Type:        schema.TypeString,
							Description: "The regex pattern used to validate the entered text",
							Optional:    true,
							Computed:    true,
						},
						"require_field": {
							Type:        schema.TypeString,
							Description: "The field or code used to determine whether the field is required or not",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
			"field_group": {
				Type:        schema.TypeList,
				Description: "Field group to add to the form",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "The name of the field group",
							Required:    true,
						},
						"description": {
							Type:        schema.TypeString,
							Description: "A description of the field group",
							Optional:    true,
							Computed:    true,
						},
						"collapsible": {
							Type:        schema.TypeBool,
							Description: "Whether the field group can be collapsed",
							Optional:    true,
							Computed:    true,
						},
						"collapsed_by_default": {
							Type:        schema.TypeBool,
							Description: "Whether the field group is collapsed by default",
							Optional:    true,
							Computed:    true,
						},
						"visibility_field": {
							Type:        schema.TypeString,
							Description: "The field or code used to trigger the visibility of the field group",
							Optional:    true,
							Computed:    true,
						},
						"option_type": {
							Type:        schema.TypeList,
							Description: "Field group option type",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": {
										Type:        schema.TypeString,
										Description: "The code of the option type to add to the field group",
										Optional:    true,
										Computed:    true,
									},
									"name": {
										Type:        schema.TypeString,
										Description: "The name of the option type to add to the field group",
										Optional:    true,
									},
									"description": {
										Type:        schema.TypeString,
										Description: "A description of the option type to add to the field group",
										Optional:    true,
										Computed:    true,
									},
									"field_name": {
										Type:        schema.TypeString,
										Description: "The field name of the option type to add to the field group",
										Optional:    true,
									},
									"type": {
										Type: schema.TypeString,
										Description: "The type of option type to add to the field group " +
											"(checkbox, hidden, number, password, radio, select, text, textarea, byteSize, " +
											"code-editor, fileContent, logoSelector, textArray, typeahead, environment)",
										ValidateFunc: validation.StringInSlice(
											[]string{
												"checkbox", "hidden", "number", "password", "radio", "select", "text",
												"textarea", "byteSize", "code-editor", "fileContent", "logoSelector",
												"textArray", "typeahead", "environment",
											},
											false,
										),
										Optional: true,
									},
									"option_list_id": {
										Type:        schema.TypeInt,
										Description: "The id of the option list for option types such as a typeahead or select list",
										Optional:    true,
										Computed:    true,
									},
									"field_label": {
										Type:        schema.TypeString,
										Description: "The label of the option type",
										Optional:    true,
										Computed:    true,
									},
									"default_value": {
										Type:        schema.TypeString,
										Description: "The default value of the option type",
										Optional:    true,
										Computed:    true,
									},
									"default_checked": {
										Type:        schema.TypeBool,
										Description: "Whether the checkbox option type is checked by default",
										Optional:    true,
										Computed:    true,
									},
									"placeholder": {
										Type:        schema.TypeString,
										Description: "The placeholder text for the option type",
										Optional:    true,
										Computed:    true,
									},
									"help_block": {
										Type:        schema.TypeString,
										Description: "The help block text for the option type",
										Optional:    true,
										Computed:    true,
									},
									"required": {
										Type:        schema.TypeBool,
										Description: "Whether the option type is required or not",
										Optional:    true,
										Computed:    true,
									},
									"export_meta": {
										Type:        schema.TypeBool,
										Description: "Whether to export the option type as a tag",
										Optional:    true,
										Default:     false,
									},
									"display_value_on_details": {
										Type:        schema.TypeBool,
										Description: "Display the selected value of the option type on the associated resource's details page",
										Optional:    true,
										Default:     false,
									},
									"locked": {
										Type:        schema.TypeBool,
										Description: "Whether the option type is locked or not",
										Optional:    true,
										Computed:    true,
									},
									"hidden": {
										Type:        schema.TypeBool,
										Description: "Whether the option type is hidden or not",
										Optional:    true,
										Computed:    true,
									},
									"exclude_from_search": {
										Type:        schema.TypeBool,
										Description: "Whether the option type should be execluded from search or not",
										Optional:    true,
										Computed:    true,
									},
									"allow_password_peek": {
										Type: schema.TypeBool,
										Description: "Whether the value of the password option type can be revealed by " +
											"the user to ensure they correctly entered the password",
										Optional: true,
										Computed: true,
									},
									"min_value": {
										Type:        schema.TypeInt,
										Description: "The minimum number that can be selected for a number option type",
										Optional:    true,
										Computed:    true,
									},
									"max_value": {
										Type:        schema.TypeInt,
										Description: "The maximum value that can be provided for a number option type",
										Optional:    true,
										Computed:    true,
									},
									"step": {
										Type:        schema.TypeInt,
										Description: "The incrementation number used for the number option type (i.e. - 5s, 10s, 100s, etc.)",
										Optional:    true,
										Computed:    true,
									},
									"text_rows": {
										Type:        schema.TypeInt,
										Description: "The number of rows to display for a text area",
										Optional:    true,
										Computed:    true,
									},
									"display": {
										Type:         schema.TypeString,
										Description:  "The memory or storage value to use (GB or MB)",
										ValidateFunc: validation.StringInSlice([]string{"GB", "MB"}, false),
										Optional:     true,
										Computed:     true,
									},
									"lock_display": {
										Type:        schema.TypeBool,
										Description: "Whether to lock the display or not",
										Optional:    true,
										Computed:    true,
									},
									"code_language": {
										Type:        schema.TypeString,
										Description: "The coding language used for highlighting code syntax",
										Optional:    true,
										Computed:    true,
									},
									"show_line_numbers": {
										Type:        schema.TypeBool,
										Description: "Whether to show the line numbers for the code editor option type",
										Optional:    true,
										Computed:    true,
									},
									"sortable": {
										Type:        schema.TypeBool,
										Description: "Whether the selected options can be sorted or not",
										Optional:    true,
										Computed:    true,
									},
									"allow_multiple_selections": {
										Type: schema.TypeBool,
										Description: "Whether to allow multiple items to be selected when using a " +
											"select list or type ahead option type",
										Optional: true,
										Computed: true,
									},
									"remove_select_option": {
										Type: schema.TypeBool,
										Description: "For Select List-type Inputs. When marked, the Input will default to " +
											"the first item in the list rather than to an empty selection",
										Optional: true,
										Computed: true,
									},
									"allow_duplicates": {
										Type:        schema.TypeBool,
										Description: "Whether duplicate selections are allowed",
										Optional:    true,
										Computed:    true,
									},
									"custom_data": {
										Type:        schema.TypeString,
										Description: "Custom JSON data payload to pass (Must be a JSON string)",
										Optional:    true,
										Computed:    true,
									},
									"dependent_field": {
										Type:        schema.TypeString,
										Description: "The field or code used to trigger the reloading of the field",
										Optional:    true,
										Computed:    true,
									},
									"delimiter": {
										Type:        schema.TypeString,
										Description: "The delimiter used to separate text array input values",
										Optional:    true,
										Computed:    true,
									},
									"visibility_field": {
										Type:        schema.TypeString,
										Description: "The field or code used to trigger the visibility of the field",
										Optional:    true,
										Computed:    true,
									},
									"verify_pattern": {
										Type:        schema.TypeString,
										Description: "The regex pattern used to validate the entered text",
										Optional:    true,
										Computed:    true,
									},
									"require_field": {
										Type:        schema.TypeString,
										Description: "The field or code used to determine whether the field is required or not",
										Optional:    true,
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceFormCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}

	// create the payload for option types not in a field group
	var optionTypes []map[string]any
	if d.Get("option_type") != nil {
		var optionTypeList []any
		if v, ok := d.Get("option_type").([]any); ok {
			optionTypeList = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("option_type", d.Get("option_type")))
		}
		// iterate over the array of optionTypes
		for i := 0; i < len(optionTypeList); i++ {
			row := make(map[string]any)
			var optionTypeConfig map[string]any
			if v, ok := optionTypeList[i].(map[string]any); ok {
				optionTypeConfig = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("option_type element", optionTypeList[i]))
			}
			row["name"] = optionTypeConfig["name"]
			row["code"] = optionTypeConfig["code"]
			row["type"] = optionTypeConfig["type"]
			row["description"] = optionTypeConfig["description"]
			row["fieldName"] = optionTypeConfig["field_name"]
			row["fieldLabel"] = optionTypeConfig["field_label"]
			row["placeHolder"] = optionTypeConfig["placeholder"]
			row["helpBlock"] = optionTypeConfig["help_block"]
			// Evaluate the option type selected
			switch optionTypeConfig["type"] {
			case typeByteSize:
				row["defaultValue"] = optionTypeConfig["default_value"]
				config := make(map[string]any)
				config["display"] = optionTypeConfig["display"]
				config["lockDisplay"] = optionTypeConfig["lock_display"]
				row["config"] = config
			case typeCodeEditor:
				row["defaultValue"] = optionTypeConfig["default_value"]
				config := make(map[string]any)
				config["lang"] = optionTypeConfig["code_language"]
				config["showLineNumbers"] = optionTypeConfig["show_line_numbers"]
				row["config"] = config
			case typeCheckbox:
				var defaultValue string
				if v, ok := optionTypeConfig["default_value"].(string); ok {
					defaultValue = v
				} else {
					return diag.FromErr(helpers.TypeAssertFailError("default_value", optionTypeConfig["default_value"]))
				}
				if defaultValue != valueTrue && defaultValue != valueFalse {
					return diag.Errorf(
						"The default_value attribute cannot be set when the type attribute is set to checkbox, "+
							"use the default_checked attribute instead for the %v checkbox resource",
						optionTypeConfig["name"],
					)
				}
				row["defaultValue"] = optionTypeConfig["default_checked"]
			case typeNumber:
				var defaultValue string
				if v, ok := optionTypeConfig["default_value"].(string); ok {
					defaultValue = v
				} else {
					return diag.FromErr(helpers.TypeAssertFailError("default_value", optionTypeConfig["default_value"]))
				}
				number, err := strconv.Atoi(defaultValue)
				if err != nil {
					return diag.Errorf(
						"The default_value attribute must be a number string when the type attribute is set to number",
					)
				}
				row["defaultValue"] = number
				row["minVal"] = optionTypeConfig["min_value"]
				row["maxVal"] = optionTypeConfig["max_value"]
				var step int
				if v, ok := optionTypeConfig["step"].(int); ok {
					step = v
				} else {
					return diag.FromErr(helpers.TypeAssertFailError("step", optionTypeConfig["step"]))
				}
				if step > 0 {
					configStep := make(map[string]any)
					configStep["step"] = optionTypeConfig["step"]
					row["config"] = configStep
				}
			case typeRadio:
				row["defaultValue"] = optionTypeConfig["default_value"]
				row["optionList"] = optionTypeConfig["option_list_id"]
			case typeSelect:
				row["defaultValue"] = optionTypeConfig["default_value"]
				row["optionList"] = optionTypeConfig["option_list_id"]
				config := make(map[string]any)
				config["multiSelect"] = optionTypeConfig["allow_multiple_selections"]
				config["sortable"] = optionTypeConfig["sortable"]
				row["config"] = config
				row["noBlank"] = optionTypeConfig["remove_select_option"]
			case typePassword:
				config := make(map[string]any)
				config["canPeek"] = optionTypeConfig["allow_password_peek"]
				row["config"] = config
			case typeTextArray:
				row["defaultValue"] = optionTypeConfig["default_value"]
				config := make(map[string]any)
				config["separator"] = optionTypeConfig["delimiter"]
				row["config"] = config
			case typeTextArea:
				row["defaultValue"] = optionTypeConfig["default_value"]
				config := make(map[string]any)
				config["rows"] = optionTypeConfig["text_rows"]
				row["config"] = config
			case typeTypeahead:
				row["defaultValue"] = optionTypeConfig["default_value"]
				config := make(map[string]any)
				config["sortable"] = optionTypeConfig["sortable"]
				config["allowDuplicates"] = optionTypeConfig["allow_duplicates"]
				config["multiSelect"] = optionTypeConfig["allow_multiple_selections"]
				config["customData"] = optionTypeConfig["custom_data"]
				row["optionList"] = optionTypeConfig["option_list_id"]
				row["config"] = config
			case typeHidden:
				row["defaultValue"] = optionTypeConfig["default_value"]
			case typeText:
				row["defaultValue"] = optionTypeConfig["default_value"]
			}
			row["required"] = optionTypeConfig["required"]
			row["exportMeta"] = optionTypeConfig["export_meta"]
			row["editable"] = optionTypeConfig["editable"]
			row["displayValueOnDetails"] = optionTypeConfig["display_value_on_details"]
			row["isLocked"] = optionTypeConfig["locked"]
			row["isHidden"] = optionTypeConfig["hidden"]
			row["excludeFromSearch"] = optionTypeConfig["exclude_from_search"]
			row["dependsOnCode"] = optionTypeConfig["dependent_field"]
			row["visibleOnCode"] = optionTypeConfig["visibility_field"]
			row["verifyPattern"] = optionTypeConfig["verify_pattern"]
			row["requireOnCode"] = optionTypeConfig["require_field"]

			optionTypes = append(optionTypes, row)
		}
	}

	// fieldGroups
	var fieldGroups []map[string]any
	if d.Get("field_group") != nil {
		var fieldGroupList []any
		if v, ok := d.Get("field_group").([]any); ok {
			fieldGroupList = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("field_group", d.Get("field_group")))
		}
		// iterate over the array of fieldGroups
		for i := 0; i < len(fieldGroupList); i++ {
			row := make(map[string]any)
			var fieldGroupConfig map[string]any
			if v, ok := fieldGroupList[i].(map[string]any); ok {
				fieldGroupConfig = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("field_group element", fieldGroupList[i]))
			}
			row["name"] = fieldGroupConfig["name"]
			row["description"] = fieldGroupConfig["description"]
			row["collapsible"] = fieldGroupConfig["collapsible"]
			row["defaultCollapsed"] = fieldGroupConfig["collapsed_by_default"]
			row["visibleOnCode"] = fieldGroupConfig["visibility_field"]
			if fieldGroupConfig["option_type"] != nil {
				// optiontypes
				var optionTypes []map[string]any
				var optionTypeList []any
				if v, ok := fieldGroupConfig["option_type"].([]any); ok {
					optionTypeList = v
				} else {
					return diag.FromErr(helpers.TypeAssertFailError("option_type", fieldGroupConfig["option_type"]))
				}
				// iterate over the array of optionTypes
				for i := 0; i < len(optionTypeList); i++ {
					row := make(map[string]any)
					var optionTypeConfig map[string]any
					if v, ok := optionTypeList[i].(map[string]any); ok {
						optionTypeConfig = v
					} else {
						return diag.FromErr(helpers.TypeAssertFailError("option_type element", optionTypeList[i]))
					}
					row["name"] = optionTypeConfig["name"]
					row["code"] = optionTypeConfig["code"]
					row["type"] = optionTypeConfig["type"]
					row["description"] = optionTypeConfig["description"]
					row["fieldName"] = optionTypeConfig["field_name"]
					row["fieldLabel"] = optionTypeConfig["field_label"]
					row["placeHolder"] = optionTypeConfig["placeholder"]
					row["helpBlock"] = optionTypeConfig["help_block"]
					// Evaluate the option type selected
					switch optionTypeConfig["type"] {
					case typeByteSize:
						row["defaultValue"] = optionTypeConfig["default_value"]
						config := make(map[string]any)
						config["display"] = optionTypeConfig["display"]
						config["lockDisplay"] = optionTypeConfig["lock_display"]
						row["config"] = config
					case typeCodeEditor:
						row["defaultValue"] = optionTypeConfig["default_value"]
						config := make(map[string]any)
						config["lang"] = optionTypeConfig["code_language"]
						config["showLineNumbers"] = optionTypeConfig["show_line_numbers"]
						row["config"] = config
					case typeCheckbox:
						var defaultValue string
						if v, ok := optionTypeConfig["default_value"].(string); ok {
							defaultValue = v
						} else {
							return diag.FromErr(helpers.TypeAssertFailError("default_value", optionTypeConfig["default_value"]))
						}
						if defaultValue != valueTrue && defaultValue != valueFalse {
							return diag.Errorf(
								"The default_value attribute cannot be set when the type attribute is set to checkbox, "+
									"use the default_checked attribute instead for the %v checkbox resource",
								optionTypeConfig["name"],
							)
						}
						row["defaultValue"] = optionTypeConfig["default_checked"]
					case typeNumber:
						var defaultValue string
						if v, ok := optionTypeConfig["default_value"].(string); ok {
							defaultValue = v
						} else {
							return diag.FromErr(helpers.TypeAssertFailError("default_value", optionTypeConfig["default_value"]))
						}
						number, err := strconv.Atoi(defaultValue)
						if err != nil {
							return diag.Errorf(
								"The default_value attribute must be a number string when the type attribute is set to number",
							)
						}
						row["defaultValue"] = number
						row["minVal"] = optionTypeConfig["min_value"]
						row["maxVal"] = optionTypeConfig["max_value"]
						var step int
						if v, ok := optionTypeConfig["step"].(int); ok {
							step = v
						} else {
							return diag.FromErr(helpers.TypeAssertFailError("step", optionTypeConfig["step"]))
						}
						if step > 0 {
							configStep := make(map[string]any)
							configStep["step"] = optionTypeConfig["step"]
							row["config"] = configStep
						}
					case typeRadio:
						row["defaultValue"] = optionTypeConfig["default_value"]
						row["optionList"] = optionTypeConfig["option_list_id"]
					case typeSelect:
						row["defaultValue"] = optionTypeConfig["default_value"]
						row["optionList"] = optionTypeConfig["option_list_id"]
						config := make(map[string]any)
						config["multiSelect"] = optionTypeConfig["allow_multiple_selections"]
						config["sortable"] = optionTypeConfig["sortable"]
						row["config"] = config
						row["noBlank"] = optionTypeConfig["remove_select_option"]
					case typePassword:
						config := make(map[string]any)
						config["canPeek"] = optionTypeConfig["allow_password_peek"]
						row["config"] = config
					case typeTextArray:
						row["defaultValue"] = optionTypeConfig["default_value"]
						config := make(map[string]any)
						config["separator"] = optionTypeConfig["delimiter"]
						row["config"] = config
					case typeTextArea:
						row["defaultValue"] = optionTypeConfig["default_value"]
						config := make(map[string]any)
						config["rows"] = optionTypeConfig["text_rows"]
						row["config"] = config
					case typeTypeahead:
						row["defaultValue"] = optionTypeConfig["default_value"]
						config := make(map[string]any)
						config["sortable"] = optionTypeConfig["sortable"]
						config["allowDuplicates"] = optionTypeConfig["allow_duplicates"]
						config["multiSelect"] = optionTypeConfig["allow_multiple_selections"]
						config["customData"] = optionTypeConfig["custom_data"]
						row["optionList"] = optionTypeConfig["option_list_id"]
						row["config"] = config
					case typeHidden:
						row["defaultValue"] = optionTypeConfig["default_value"]
					case typeText:
						row["defaultValue"] = optionTypeConfig["default_value"]
					}
					row["required"] = optionTypeConfig["required"]
					row["exportMeta"] = optionTypeConfig["export_meta"]
					row["editable"] = optionTypeConfig["editable"]
					row["displayValueOnDetails"] = optionTypeConfig["display_value_on_details"]
					row["isLocked"] = optionTypeConfig["locked"]
					row["isHidden"] = optionTypeConfig["hidden"]
					row["excludeFromSearch"] = optionTypeConfig["exclude_from_search"]
					row["dependsOnCode"] = optionTypeConfig["dependent_field"]
					row["visibleOnCode"] = optionTypeConfig["visibility_field"]
					row["verifyPattern"] = optionTypeConfig["verify_pattern"]
					row["requireOnCode"] = optionTypeConfig["require_field"]

					optionTypes = append(optionTypes, row)
				}
				row["options"] = optionTypes
			}
			fieldGroups = append(fieldGroups, row)
		}
	}

	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		if labelSet, ok := attr.(*schema.Set); ok {
			for _, s := range labelSet.List() {
				if labelStr, ok := s.(string); ok {
					labelsPayload = append(labelsPayload, labelStr)
				} else {
					return diag.FromErr(helpers.TypeAssertFailError("labels element", s))
				}
			}
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("labels", attr))
		}
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("labels", d.Get("labels")))
	}

	var code string
	if v, ok := d.Get("code").(string); ok {
		code = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("code", d.Get("code")))
	}

	var description string
	if v, ok := d.Get("description").(string); ok {
		description = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("description", d.Get("description")))
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"optionTypeForm": map[string]any{
				"name":        name,
				"code":        code,
				"description": description,
				"labels":      labelsPayload,
				"fieldGroups": fieldGroups,
				"options":     optionTypes,
			},
		},
	}
	//	jsonRequest, _ := json.Marshal(req.Body)
	//	log.Printf("API JSON REQUEST: %s", string(jsonRequest))

	resp, err := client.CreateForm(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	var result *morpheus.CreateFormResult
	if v, ok := resp.Result.(*morpheus.CreateFormResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.Form == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Form"))
	}

	formResult := result.Form
	// Successfully created resource, now set id
	d.SetId(convert.Int64ToString(formResult.ID))

	diags = append(diags, resourceFormRead(ctx, d, meta)...)

	return diags
}

func resourceFormRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	id := d.Id()

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error

	if id == "" && name != "" {
		resp, err = client.FindFormByName(name)
	} else if id != "" {
		resp, err = client.GetForm(convert.StringToInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Form cannot be read without name or id")
	}

	if err != nil {
		// 404 is ok?
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)
			d.SetId("")

			return diags
		} else {
			log.Printf("API FAILURE: %s - %s", resp, err)

			return diag.FromErr(err)
		}
	}
	log.Printf("API RESPONSE: %s", resp)

	// store resource data
	var result *morpheus.GetFormResult
	if v, ok := resp.Result.(*morpheus.GetFormResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}
	form := result.Form

	d.SetId(convert.Int64ToString(form.ID))
	d.Set("name", form.Name)
	d.Set("code", form.Code)
	d.Set("description", form.Description)
	d.Set("labels", form.Labels)

	// Option Types
	var optionTypes []map[string]any
	if len(form.Options) != 0 {
		// optionTypeList := d.Get("option_type").([]any)
		for _, optionType := range form.Options {
			row := make(map[string]any)
			switch optionType.Type {
			case typeByteSize:
				row["display"] = optionType.Config.Display
				row["lock_display"] = optionType.Config.LockDisplay
			case typeCheckbox:
				// convert string text to boolean
				if optionType.DefaultValue == "true" {
					row["default_checked"] = true
				} else {
					row["default_checked"] = false
				}
			case typeCodeEditor:
				row["show_line_numbers"] = optionType.Config.ShowLineNumbers
				row["code_language"] = optionType.Config.Lang
			case typeNumber:
				row["step"] = optionType.Config.Step
				row["min_value"] = optionType.MinVal
				row["max_value"] = optionType.MaxVal
			case typeRadio:
				row["option_list_id"] = optionType.OptionList.ID
			case typeSelect:
				row["option_list_id"] = optionType.OptionList.ID
			case typeTextArea:
				row["text_rows"] = optionType.Config.Rows
			case typeHidden:
				log.Printf("HIDDEN DEFAULT: %v", optionType.DefaultValue)
			case typeTextArray:
				row["delimiter"] = optionType.Config.Separator
			case typeTypeahead:
				row["sortable"] = optionType.Config.Sortable
				row["allow_duplicates"] = optionType.Config.AllowDuplicates
				row["custom_data"] = optionType.Config.CustomData
				row["allow_multiple_selections"] = optionType.Config.MultiSelect
				row["option_list_id"] = optionType.OptionList.ID
			}
			row["remove_select_option"] = optionType.NoBlank
			row["name"] = optionType.Name
			row["description"] = optionType.Description
			row["code"] = optionType.Code
			row["type"] = optionType.Type
			row["field_label"] = optionType.FieldLabel
			row["field_name"] = optionType.FieldName
			row["default_value"] = optionType.DefaultValue
			row["placeholder"] = optionType.PlaceHolder
			row["help_block"] = optionType.HelpBlock
			row["required"] = optionType.Required
			row["export_meta"] = optionType.ExportMeta
			row["display_value_on_details"] = optionType.DisplayValueOnDetails
			row["locked"] = optionType.IsLocked
			row["hidden"] = optionType.IsHidden
			row["exclude_from_search"] = optionType.ExcludeFromSearch
			row["dependent_field"] = optionType.DependsOnCode
			row["visibility_field"] = optionType.VisibleOnCode
			row["verify_pattern"] = optionType.VerifyPattern
			row["require_field"] = optionType.RequireOnCode
			optionTypes = append(optionTypes, row)
		}
	}
	d.Set("option_type", optionTypes)

	// Field Groups
	var fieldGroups []map[string]any
	if len(form.FieldGroups) != 0 {
		for _, fieldGroup := range form.FieldGroups {
			row := make(map[string]any)
			row["name"] = fieldGroup.Name
			row["description"] = fieldGroup.Description
			row["collapsible"] = fieldGroup.Collapsible
			row["collapsed_by_default"] = fieldGroup.DefaultCollapsed
			row["visibility_field"] = fieldGroup.VisibleOnCode
			var fgOptionTypes []map[string]any
			if len(fieldGroup.Options) != 0 {
				for _, optionType := range fieldGroup.Options {
					optionTypeRow := make(map[string]any)
					// Check if the input uses an existing input or not
					switch optionType.Type {
					case typeByteSize:
						optionTypeRow["display"] = optionType.Config.Display
						optionTypeRow["lock_display"] = optionType.Config.LockDisplay
					case typeCheckbox:
						// convert string text to boolean
						if optionType.DefaultValue == "true" {
							row["default_checked"] = true
						} else {
							row["default_checked"] = false
						}
					case typeCodeEditor:
						optionTypeRow["show_line_numbers"] = optionType.Config.ShowLineNumbers
						optionTypeRow["code_language"] = optionType.Config.Lang
					case typeNumber:
						optionTypeRow["step"] = optionType.Config.Step
						optionTypeRow["min_value"] = optionType.MinVal
						optionTypeRow["max_value"] = optionType.MaxVal
					case typeRadio:
						optionTypeRow["option_list_id"] = optionType.OptionList.ID
					case typeSelect:
						optionTypeRow["option_list_id"] = optionType.OptionList.ID
					case typeTextArea:
						optionTypeRow["text_rows"] = optionType.Config.Rows
					case typeTextArray:
						optionTypeRow["delimiter"] = optionType.Config.Separator
					case typeTypeahead:
						optionTypeRow["sortable"] = optionType.Config.Sortable
						optionTypeRow["allow_duplicates"] = optionType.Config.AllowDuplicates
						optionTypeRow["custom_data"] = optionType.Config.CustomData
						optionTypeRow["allow_multiple_selections"] = optionType.Config.MultiSelect
						optionTypeRow["option_list_id"] = optionType.OptionList.ID
					}
					optionTypeRow["remove_select_option"] = optionType.NoBlank
					optionTypeRow["name"] = optionType.Name
					optionTypeRow["description"] = optionType.Description
					optionTypeRow["code"] = optionType.Code
					optionTypeRow["type"] = optionType.Type
					optionTypeRow["field_label"] = optionType.FieldLabel
					optionTypeRow["field_name"] = optionType.FieldName
					optionTypeRow["default_value"] = optionType.DefaultValue
					optionTypeRow["placeholder"] = optionType.PlaceHolder
					optionTypeRow["help_block"] = optionType.HelpBlock
					optionTypeRow["required"] = optionType.Required
					optionTypeRow["export_meta"] = optionType.ExportMeta
					optionTypeRow["display_value_on_details"] = optionType.DisplayValueOnDetails
					optionTypeRow["locked"] = optionType.IsLocked
					optionTypeRow["hidden"] = optionType.IsHidden
					optionTypeRow["exclude_from_search"] = optionType.ExcludeFromSearch
					optionTypeRow["dependent_field"] = optionType.DependsOnCode
					optionTypeRow["visibility_field"] = optionType.VisibleOnCode
					optionTypeRow["verify_pattern"] = optionType.VerifyPattern
					optionTypeRow["require_field"] = optionType.RequireOnCode
					fgOptionTypes = append(fgOptionTypes, optionTypeRow)
				}
			}
			row["option_type"] = fgOptionTypes
			fieldGroups = append(fieldGroups, row)
		}
	}
	d.Set("field_group", fieldGroups)

	return diags
}

func resourceFormUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	id := d.Id()

	// create the payload for option types not in a field group
	var optionTypes []map[string]any
	if d.Get("option_type") != nil {
		var optionTypeList []any
		if v, ok := d.Get("option_type").([]any); ok {
			optionTypeList = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("option_type", d.Get("option_type")))
		}
		// iterate over the array of optionTypes
		for i := 0; i < len(optionTypeList); i++ {
			row := make(map[string]any)
			var optionTypeConfig map[string]any
			if v, ok := optionTypeList[i].(map[string]any); ok {
				optionTypeConfig = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("option_type element", optionTypeList[i]))
			}
			row["name"] = optionTypeConfig["name"]
			row["code"] = optionTypeConfig["code"]
			row["type"] = optionTypeConfig["type"]
			row["description"] = optionTypeConfig["description"]
			row["fieldName"] = optionTypeConfig["field_name"]
			row["fieldLabel"] = optionTypeConfig["field_label"]
			row["placeHolder"] = optionTypeConfig["placeholder"]
			row["helpBlock"] = optionTypeConfig["help_block"]
			// Evaluate the option type selected
			switch optionTypeConfig["type"] {
			case typeByteSize:
				row["defaultValue"] = optionTypeConfig["default_value"]
				config := make(map[string]any)
				config["display"] = optionTypeConfig["display"]
				config["lockDisplay"] = optionTypeConfig["lock_display"]
				row["config"] = config
			case typeCodeEditor:
				row["defaultValue"] = optionTypeConfig["default_value"]
				config := make(map[string]any)
				config["lang"] = optionTypeConfig["code_language"]
				config["showLineNumbers"] = optionTypeConfig["show_line_numbers"]
				row["config"] = config
			case typeCheckbox:
				var defaultValue string
				if v, ok := optionTypeConfig["default_value"].(string); ok {
					defaultValue = v
				} else {
					return diag.FromErr(helpers.TypeAssertFailError("default_value", optionTypeConfig["default_value"]))
				}
				if defaultValue != valueTrue && defaultValue != valueFalse {
					return diag.Errorf(
						"The default_value attribute cannot be set when the type attribute is set to checkbox, "+
							"use the default_checked attribute instead for the %v checkbox resource: %v",
						optionTypeConfig["name"],
						defaultValue,
					)
				}
				row["defaultValue"] = optionTypeConfig["default_checked"]
			case typeNumber:
				var defaultValue string
				if v, ok := optionTypeConfig["default_value"].(string); ok {
					defaultValue = v
				} else {
					return diag.FromErr(helpers.TypeAssertFailError("default_value", optionTypeConfig["default_value"]))
				}
				number, err := strconv.Atoi(defaultValue)
				if err != nil {
					return diag.Errorf(
						"The default_value attribute must be a number string when the type attribute is set to number",
					)
				}
				row["defaultValue"] = number
				row["minVal"] = optionTypeConfig["min_value"]
				row["maxVal"] = optionTypeConfig["max_value"]
				var step int
				if v, ok := optionTypeConfig["step"].(int); ok {
					step = v
				} else {
					return diag.FromErr(helpers.TypeAssertFailError("step", optionTypeConfig["step"]))
				}
				if step > 0 {
					configStep := make(map[string]any)
					configStep["step"] = optionTypeConfig["step"]
					row["config"] = configStep
				}
			case typeRadio:
				row["defaultValue"] = optionTypeConfig["default_value"]
				row["optionList"] = optionTypeConfig["option_list_id"]
			case typeSelect:
				row["defaultValue"] = optionTypeConfig["default_value"]
				row["optionList"] = optionTypeConfig["option_list_id"]
				config := make(map[string]any)
				config["multiSelect"] = optionTypeConfig["allow_multiple_selections"]
				config["sortable"] = optionTypeConfig["sortable"]
				row["config"] = config
				row["noBlank"] = optionTypeConfig["remove_select_option"]
			case typePassword:
				config := make(map[string]any)
				config["canPeek"] = optionTypeConfig["allow_password_peek"]
				row["config"] = config
			case typeTextArray:
				row["defaultValue"] = optionTypeConfig["default_value"]
				config := make(map[string]any)
				config["separator"] = optionTypeConfig["delimiter"]
				row["config"] = config
			case typeTextArea:
				row["defaultValue"] = optionTypeConfig["default_value"]
				config := make(map[string]any)
				config["rows"] = optionTypeConfig["text_rows"]
				row["config"] = config
			case typeTypeahead:
				row["defaultValue"] = optionTypeConfig["default_value"]
				config := make(map[string]any)
				config["sortable"] = optionTypeConfig["sortable"]
				config["allowDuplicates"] = optionTypeConfig["allow_duplicates"]
				config["multiSelect"] = optionTypeConfig["allow_multiple_selections"]
				config["customData"] = optionTypeConfig["custom_data"]
				row["optionList"] = optionTypeConfig["option_list_id"]
				row["config"] = config
			case typeHidden:
				row["defaultValue"] = optionTypeConfig["default_value"]
			case typeText:
				row["defaultValue"] = optionTypeConfig["default_value"]
			}
			row["required"] = optionTypeConfig["required"]
			row["exportMeta"] = optionTypeConfig["export_meta"]
			row["editable"] = optionTypeConfig["editable"]
			row["displayValueOnDetails"] = optionTypeConfig["display_value_on_details"]
			row["isLocked"] = optionTypeConfig["locked"]
			row["isHidden"] = optionTypeConfig["hidden"]
			row["excludeFromSearch"] = optionTypeConfig["exclude_from_search"]
			row["dependsOnCode"] = optionTypeConfig["dependent_field"]
			row["visibleOnCode"] = optionTypeConfig["visibility_field"]
			row["verifyPattern"] = optionTypeConfig["verify_pattern"]
			row["requireOnCode"] = optionTypeConfig["require_field"]

			optionTypes = append(optionTypes, row)
		}
	}

	// fieldGroups
	var fieldGroups []map[string]any
	if d.Get("field_group") != nil {
		var fieldGroupList []any
		if v, ok := d.Get("field_group").([]any); ok {
			fieldGroupList = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("field_group", d.Get("field_group")))
		}
		// iterate over the array of fieldGroups
		for i := 0; i < len(fieldGroupList); i++ {
			row := make(map[string]any)
			var fieldGroupConfig map[string]any
			if v, ok := fieldGroupList[i].(map[string]any); ok {
				fieldGroupConfig = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("field_group element", fieldGroupList[i]))
			}
			row["name"] = fieldGroupConfig["name"]
			row["description"] = fieldGroupConfig["description"]
			row["collapsible"] = fieldGroupConfig["collapsible"]
			row["defaultCollapsed"] = fieldGroupConfig["collapsed_by_default"]
			row["visibleOnCode"] = fieldGroupConfig["visibility_field"]
			if fieldGroupConfig["option_type"] != nil {
				// optiontypes
				var optionTypes []map[string]any
				var optionTypeList []any
				if v, ok := fieldGroupConfig["option_type"].([]any); ok {
					optionTypeList = v
				} else {
					return diag.FromErr(helpers.TypeAssertFailError("option_type", fieldGroupConfig["option_type"]))
				}
				// iterate over the array of optionTypes
				for i := 0; i < len(optionTypeList); i++ {
					row := make(map[string]any)
					var optionTypeConfig map[string]any
					if v, ok := optionTypeList[i].(map[string]any); ok {
						optionTypeConfig = v
					} else {
						return diag.FromErr(helpers.TypeAssertFailError("option_type element", optionTypeList[i]))
					}
					row["name"] = optionTypeConfig["name"]
					row["code"] = optionTypeConfig["code"]
					row["type"] = optionTypeConfig["type"]
					row["description"] = optionTypeConfig["description"]
					row["fieldName"] = optionTypeConfig["field_name"]
					row["fieldLabel"] = optionTypeConfig["field_label"]
					row["placeHolder"] = optionTypeConfig["placeholder"]
					row["helpBlock"] = optionTypeConfig["help_block"]
					// Evaluate the option type selected
					switch optionTypeConfig["type"] {
					case typeByteSize:
						row["defaultValue"] = optionTypeConfig["default_value"]
						config := make(map[string]any)
						config["display"] = optionTypeConfig["display"]
						config["lockDisplay"] = optionTypeConfig["lock_display"]
						row["config"] = config
					case typeCodeEditor:
						row["defaultValue"] = optionTypeConfig["default_value"]
						config := make(map[string]any)
						config["lang"] = optionTypeConfig["code_language"]
						config["showLineNumbers"] = optionTypeConfig["show_line_numbers"]
						row["config"] = config
					case typeCheckbox:
						var defaultValue string
						if v, ok := optionTypeConfig["default_value"].(string); ok {
							defaultValue = v
						} else {
							return diag.FromErr(helpers.TypeAssertFailError("default_value", optionTypeConfig["default_value"]))
						}
						if defaultValue != valueTrue && defaultValue != valueFalse {
							return diag.Errorf(
								"The default_value attribute cannot be set when the type attribute is set to checkbox, "+
									"use the default_checked attribute instead for the %v checkbox resource",
								optionTypeConfig["name"],
							)
						}
						row["defaultValue"] = optionTypeConfig["default_checked"]
					case typeNumber:
						var defaultValue string
						if v, ok := optionTypeConfig["default_value"].(string); ok {
							defaultValue = v
						} else {
							return diag.FromErr(helpers.TypeAssertFailError("default_value", optionTypeConfig["default_value"]))
						}
						number, err := strconv.Atoi(defaultValue)
						if err != nil {
							return diag.Errorf(
								"The default_value attribute must be a number string when the type attribute is set to number",
							)
						}
						row["defaultValue"] = number
						row["minVal"] = optionTypeConfig["min_value"]
						row["maxVal"] = optionTypeConfig["max_value"]

						var step int
						if v, ok := optionTypeConfig["step"].(int); ok {
							step = v
						} else {
							return diag.FromErr(helpers.TypeAssertFailError("step", optionTypeConfig["step"]))
						}

						if step > 0 {
							configStep := make(map[string]any)
							configStep["step"] = optionTypeConfig["step"]
							row["config"] = configStep
						}
					case typeRadio:
						row["defaultValue"] = optionTypeConfig["default_value"]
						row["optionList"] = optionTypeConfig["option_list_id"]
					case typeSelect:
						row["defaultValue"] = optionTypeConfig["default_value"]
						row["optionList"] = optionTypeConfig["option_list_id"]
						config := make(map[string]any)
						config["multiSelect"] = optionTypeConfig["allow_multiple_selections"]
						config["sortable"] = optionTypeConfig["sortable"]
						row["config"] = config
						row["noBlank"] = optionTypeConfig["remove_select_option"]
					case typePassword:
						config := make(map[string]any)
						config["canPeek"] = optionTypeConfig["allow_password_peek"]
						row["config"] = config
					case typeTextArray:
						row["defaultValue"] = optionTypeConfig["default_value"]
						config := make(map[string]any)
						config["separator"] = optionTypeConfig["delimiter"]
						row["config"] = config
					case typeTextArea:
						row["defaultValue"] = optionTypeConfig["default_value"]
						config := make(map[string]any)
						config["rows"] = optionTypeConfig["text_rows"]
						row["config"] = config
					case typeTypeahead:
						row["defaultValue"] = optionTypeConfig["default_value"]
						config := make(map[string]any)
						config["sortable"] = optionTypeConfig["sortable"]
						config["allowDuplicates"] = optionTypeConfig["allow_duplicates"]
						config["multiSelect"] = optionTypeConfig["allow_multiple_selections"]
						config["customData"] = optionTypeConfig["custom_data"]
						row["optionList"] = optionTypeConfig["option_list_id"]
						row["config"] = config
					case typeHidden:
						row["defaultValue"] = optionTypeConfig["default_value"]
					case typeText:
						row["defaultValue"] = optionTypeConfig["default_value"]
					}
					row["required"] = optionTypeConfig["required"]
					row["exportMeta"] = optionTypeConfig["export_meta"]
					row["editable"] = optionTypeConfig["editable"]
					row["displayValueOnDetails"] = optionTypeConfig["display_value_on_details"]
					row["isLocked"] = optionTypeConfig["locked"]
					row["isHidden"] = optionTypeConfig["hidden"]
					row["excludeFromSearch"] = optionTypeConfig["exclude_from_search"]
					row["dependsOnCode"] = optionTypeConfig["dependent_field"]
					row["visibleOnCode"] = optionTypeConfig["visibility_field"]
					row["verifyPattern"] = optionTypeConfig["verify_pattern"]
					row["requireOnCode"] = optionTypeConfig["require_field"]

					optionTypes = append(optionTypes, row)
				}
				row["options"] = optionTypes
			}
			fieldGroups = append(fieldGroups, row)
		}
	}

	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		if labelSet, ok := attr.(*schema.Set); ok {
			for _, s := range labelSet.List() {
				if labelStr, ok := s.(string); ok {
					labelsPayload = append(labelsPayload, labelStr)
				} else {
					return diag.FromErr(helpers.TypeAssertFailError("labels element", s))
				}
			}
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("labels", d.Get("labels")))
		}
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("labels", d.Get("labels")))
	}

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}

	var code string
	if v, ok := d.Get("code").(string); ok {
		code = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("code", d.Get("code")))
	}

	var description string
	if v, ok := d.Get("description").(string); ok {
		description = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("description", d.Get("description")))
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"optionTypeForm": map[string]any{
				"name":        name,
				"code":        code,
				"description": description,
				"labels":      labelsPayload,
				"fieldGroups": fieldGroups,
				"options":     optionTypes,
			},
		},
	}

	resp, err := client.UpdateForm(convert.StringToInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	var result *morpheus.UpdateFormResult
	if v, ok := resp.Result.(*morpheus.UpdateFormResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}
	formResult := result.Form
	// Successfully created resource, now set id
	d.SetId(convert.Int64ToString(formResult.ID))

	return resourceFormRead(ctx, d, meta)
}

func resourceFormDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}

	// Check if the form is already in use
	var inUseCatalogItems []string
	catalogItemsResp, err := client.ListCatalogItems(&morpheus.Request{
		QueryParams: map[string]string{
			"max": "500",
		},
	})
	if err != nil {
		if catalogItemsResp != nil && catalogItemsResp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", catalogItemsResp, err)

			return diag.FromErr(err)
		} else {
			log.Printf("API FAILURE: %s - %s", catalogItemsResp, err)

			return diag.FromErr(err)
		}
	}

	var result *morpheus.ListCatalogItemsResult
	if v, ok := catalogItemsResp.Result.(*morpheus.ListCatalogItemsResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", catalogItemsResp.Result))
	}
	catalogItems := result.CatalogItems
	for _, catalogItem := range *catalogItems {
		if catalogItem.Form.ID == convert.StringToInt64(id) {
			inUseCatalogItems = append(inUseCatalogItems, catalogItem.Name)
		}
	}
	if len(inUseCatalogItems) > 0 {
		return diag.Errorf(
			"The %s morpheus_form resource is currently associated with the following catalog items "+
				"and must be disassociated before being deleted: %s",
			d.Get("name"),
			inUseCatalogItems,
		)
	}

	resp, err := client.DeleteForm(convert.StringToInt64(id), req)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)

			return diag.FromErr(err)
		} else {
			log.Printf("API FAILURE: %s - %s", resp, err)

			return diag.FromErr(err)
		}
	}
	log.Printf("API RESPONSE: %s", resp)
	d.SetId("")

	return diags
}

// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package optiontype

import (
	"context"
	"log"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/convert"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceOptionTypeRadioList() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus radio list option type resource",
		CreateContext: resourceOptionTypeRadioListCreate,
		ReadContext:   resourceOptionTypeRadioListRead,
		UpdateContext: resourceOptionTypeRadioListUpdate,
		DeleteContext: resourceOptionTypeRadioListDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the radio list option type",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the radio list option type",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the radio list option type",
				Optional:    true,
				Computed:    true,
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "The organization labels associated with the option type(Only supported on Morpheus 5.5.3 or higher)",
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"field_name": {
				Type:        schema.TypeString,
				Description: "The field name of the radio list option type",
				Required:    true,
			},
			"export_meta": {
				Type:        schema.TypeBool,
				Description: "Whether to export the radio list option type as a tag",
				Optional:    true,
				Default:     false,
			},
			"dependent_field": {
				Type:        schema.TypeString,
				Description: "The field or code used to trigger the reloading of the field",
				Optional:    true,
				Computed:    true,
			},
			"visibility_field": {
				Type:        schema.TypeString,
				Description: "The field or code used to trigger the visibility of the field",
				Optional:    true,
				Computed:    true,
			},
			"require_field": {
				Type:        schema.TypeString,
				Description: "The field or code used to trigger the required status of the field",
				Optional:    true,
				Computed:    true,
			},
			"show_on_edit": {
				Type:        schema.TypeBool,
				Description: "Whether the option type will display in the edit section of the provisioned resource",
				Optional:    true,
				Computed:    true,
			},
			"editable": {
				Type:        schema.TypeBool,
				Description: "Whether the value of the option type can be edited after the initial request",
				Optional:    true,
				Computed:    true,
			},
			"display_value_on_details": {
				Type:        schema.TypeBool,
				Description: "Display the selected value of the radio list option type on the associated resource's details page",
				Optional:    true,
				Default:     false,
			},
			"field_label": {
				Type:        schema.TypeString,
				Description: "The label associated with the field in the UI",
				Required:    true,
			},
			"default_value": {
				Type:        schema.TypeString,
				Description: "The default value of the option type",
				Optional:    true,
				Computed:    true,
			},
			"help_block": {
				Type:        schema.TypeString,
				Description: "Text that provides additional details about the use of the option type",
				Optional:    true,
				Computed:    true,
			},
			"required": {
				Type:        schema.TypeBool,
				Description: "Whether the option type is required",
				Optional:    true,
				Default:     false,
			},
			"option_list_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the associated option list",
				Required:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceOptionTypeRadioListCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}

	var description string
	if v, ok := d.Get("description").(string); ok {
		description = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("description", d.Get("description")))
	}

	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		if labelSet, ok := attr.(*schema.Set); ok {
			for _, s := range labelSet.List() {
				if labelStr, ok := s.(string); ok {
					labelsPayload = append(labelsPayload, labelStr)
				} else {
					return diag.FromErr(helpers.TypeAssertFailError("label", s))
				}
			}
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("labels", attr))
		}
	}

	var fieldName string
	if v, ok := d.Get("field_name").(string); ok {
		fieldName = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("field_name", d.Get("field_name")))
	}

	var exportMeta bool
	if v, ok := d.Get("export_meta").(bool); ok {
		exportMeta = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("export_meta", d.Get("export_meta")))
	}

	var dependentField string
	if v, ok := d.Get("dependent_field").(string); ok {
		dependentField = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("dependent_field", d.Get("dependent_field")))
	}

	var visibilityField string
	if v, ok := d.Get("visibility_field").(string); ok {
		visibilityField = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("visibility_field", d.Get("visibility_field")))
	}

	var requireField string
	if v, ok := d.Get("require_field").(string); ok {
		requireField = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("require_field", d.Get("require_field")))
	}

	var showOnEdit bool
	if v, ok := d.Get("show_on_edit").(bool); ok {
		showOnEdit = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("show_on_edit", d.Get("show_on_edit")))
	}

	var editable bool
	if v, ok := d.Get("editable").(bool); ok {
		editable = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("editable", d.Get("editable")))
	}

	var displayValueOnDetails bool
	if v, ok := d.Get("display_value_on_details").(bool); ok {
		displayValueOnDetails = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("display_value_on_details", d.Get("display_value_on_details")))
	}

	var fieldLabel string
	if v, ok := d.Get("field_label").(string); ok {
		fieldLabel = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("field_label", d.Get("field_label")))
	}

	var defaultValue string
	if v, ok := d.Get("default_value").(string); ok {
		defaultValue = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("default_value", d.Get("default_value")))
	}

	var helpBlock string
	if v, ok := d.Get("help_block").(string); ok {
		helpBlock = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("help_block", d.Get("help_block")))
	}

	var required bool
	if v, ok := d.Get("required").(bool); ok {
		required = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("required", d.Get("required")))
	}

	var optionListID int
	if v, ok := d.Get("option_list_id").(int); ok {
		optionListID = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("option_list_id", d.Get("option_list_id")))
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"optionType": map[string]any{
				"name":                  name,
				"description":           description,
				"labels":                labelsPayload,
				"fieldName":             fieldName,
				"exportMeta":            exportMeta,
				"dependsOnCode":         dependentField,
				"visibleOnCode":         visibilityField,
				"requireOnCode":         requireField,
				"showOnEdit":            showOnEdit,
				"editable":              editable,
				"displayValueOnDetails": displayValueOnDetails,
				"type":                  "radio",
				"fieldLabel":            fieldLabel,
				"defaultValue":          defaultValue,
				"helpBlock":             helpBlock,
				"required":              required,
				"optionList": map[string]any{
					"id": optionListID,
				},
			},
		},
	}
	resp, err := client.CreateOptionType(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	var result *morpheus.CreateOptionTypeResult
	if v, ok := resp.Result.(*morpheus.CreateOptionTypeResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("CreateOptionTypeResult", resp.Result))
	}

	if result.OptionType == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("OptionType"))
	}
	environment := result.OptionType
	// Successfully created resource, now set id
	d.SetId(convert.Int64ToString(environment.ID))

	diags = append(diags, resourceOptionTypeRadioListRead(ctx, d, meta)...)

	return diags
}

func resourceOptionTypeRadioListRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}
	// Warning or errors can be collected in a slice type
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
		resp, err = client.FindOptionTypeByName(name)
	} else if id != "" {
		resp, err = client.GetOptionType(convert.StringToInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("OptionType cannot be read without name or id")
	}

	if err != nil {
		// 404 is ok?
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)
			log.Printf("Forcing recreation of resource")
			d.SetId("")

			return diags
		} else {
			log.Printf("API FAILURE: %s - %s", resp, err)

			return diag.FromErr(err)
		}
	}
	log.Printf("API RESPONSE: %s", resp)

	// store resource data
	var result *morpheus.GetOptionTypeResult
	if v, ok := resp.Result.(*morpheus.GetOptionTypeResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("GetOptionTypeResult", resp.Result))
	}

	if result.OptionType == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("OptionType"))
	}
	optionType := result.OptionType

	d.SetId(convert.Int64ToString(optionType.ID))
	d.Set("name", optionType.Name)
	d.Set("description", optionType.Description)
	d.Set("labels", optionType.Labels)
	d.Set("field_name", optionType.FieldName)
	d.Set("export_meta", optionType.ExportMeta)
	d.Set("dependent_field", optionType.DependsOnCode)
	d.Set("visibility_field", optionType.VisibleOnCode)
	d.Set("require_field", optionType.RequireOnCode)
	d.Set("show_on_edit", optionType.ShowOnEdit)
	d.Set("editable", optionType.Editable)
	d.Set("display_value_on_details", optionType.DisplayValueOnDetails)
	d.Set("field_label", optionType.FieldLabel)
	d.Set("default_value", optionType.DefaultValue)
	d.Set("help_block", optionType.HelpBlock)
	d.Set("required", optionType.Required)
	d.Set("option_list_id", optionType.OptionList.ID)

	return diags
}

func resourceOptionTypeRadioListUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	id := d.Id()

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}

	var description string
	if v, ok := d.Get("description").(string); ok {
		description = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("description", d.Get("description")))
	}

	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		if labelSet, ok := attr.(*schema.Set); ok {
			for _, s := range labelSet.List() {
				if labelStr, ok := s.(string); ok {
					labelsPayload = append(labelsPayload, labelStr)
				} else {
					return diag.FromErr(helpers.TypeAssertFailError("label", s))
				}
			}
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("labels", attr))
		}
	}

	var fieldName string
	if v, ok := d.Get("field_name").(string); ok {
		fieldName = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("field_name", d.Get("field_name")))
	}

	var exportMeta bool
	if v, ok := d.Get("export_meta").(bool); ok {
		exportMeta = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("export_meta", d.Get("export_meta")))
	}

	var dependentField string
	if v, ok := d.Get("dependent_field").(string); ok {
		dependentField = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("dependent_field", d.Get("dependent_field")))
	}

	var visibilityField string
	if v, ok := d.Get("visibility_field").(string); ok {
		visibilityField = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("visibility_field", d.Get("visibility_field")))
	}

	var requireField string
	if v, ok := d.Get("require_field").(string); ok {
		requireField = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("require_field", d.Get("require_field")))
	}

	var showOnEdit bool
	if v, ok := d.Get("show_on_edit").(bool); ok {
		showOnEdit = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("show_on_edit", d.Get("show_on_edit")))
	}

	var editable bool
	if v, ok := d.Get("editable").(bool); ok {
		editable = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("editable", d.Get("editable")))
	}

	var displayValueOnDetails bool
	if v, ok := d.Get("display_value_on_details").(bool); ok {
		displayValueOnDetails = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("display_value_on_details", d.Get("display_value_on_details")))
	}

	var fieldLabel string
	if v, ok := d.Get("field_label").(string); ok {
		fieldLabel = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("field_label", d.Get("field_label")))
	}

	var defaultValue string
	if v, ok := d.Get("default_value").(string); ok {
		defaultValue = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("default_value", d.Get("default_value")))
	}

	var helpBlock string
	if v, ok := d.Get("help_block").(string); ok {
		helpBlock = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("help_block", d.Get("help_block")))
	}

	var required bool
	if v, ok := d.Get("required").(bool); ok {
		required = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("required", d.Get("required")))
	}

	var optionListID int
	if v, ok := d.Get("option_list_id").(int); ok {
		optionListID = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("option_list_id", d.Get("option_list_id")))
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"optionType": map[string]any{
				"name":                  name,
				"description":           description,
				"labels":                labelsPayload,
				"fieldName":             fieldName,
				"exportMeta":            exportMeta,
				"dependsOnCode":         dependentField,
				"visibleOnCode":         visibilityField,
				"requireOnCode":         requireField,
				"showOnEdit":            showOnEdit,
				"editable":              editable,
				"displayValueOnDetails": displayValueOnDetails,
				"type":                  "radio",
				"fieldLabel":            fieldLabel,
				"defaultValue":          defaultValue,
				"helpBlock":             helpBlock,
				"required":              required,
				"optionList": map[string]any{
					"id": optionListID,
				},
			},
		},
	}
	resp, err := client.UpdateOptionType(convert.StringToInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	var result *morpheus.UpdateOptionTypeResult
	if v, ok := resp.Result.(*morpheus.UpdateOptionTypeResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("UpdateOptionTypeResult", resp.Result))
	}

	if result.OptionType == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("OptionType"))
	}
	account := result.OptionType
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(convert.Int64ToString(account.ID))

	return resourceOptionTypeRadioListRead(ctx, d, meta)
}

func resourceOptionTypeRadioListDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteOptionType(convert.StringToInt64(id), req)
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

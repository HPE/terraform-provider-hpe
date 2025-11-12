// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package optiontype

import (
	"context"
	"log"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceOptionTypeText() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus text option type resource",
		CreateContext: resourceOptionTypeTextCreate,
		ReadContext:   resourceOptionTypeTextRead,
		UpdateContext: resourceOptionTypeTextUpdate,
		DeleteContext: resourceOptionTypeTextDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the text option type",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the text option type",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the text option type",
				Optional:    true,
				Computed:    true,
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "The organization labels associated with the option type (Only supported on Morpheus 5.5.3 or higher)",
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"field_name": {
				Type:        schema.TypeString,
				Description: "The field name of the text option type",
				Optional:    true,
				Computed:    true,
			},
			"export_meta": {
				Type:        schema.TypeBool,
				Description: "Whether to export the text option type as a tag",
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
				Description: "The field or code used to determine whether the field is required or not",
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
				Description: "Display the selected value of the text option type on the associated resource's details page",
				Optional:    true,
				Default:     false,
			},
			"field_label": {
				Type:        schema.TypeString,
				Description: "The label associated with the field in the UI",
				Optional:    true,
				Computed:    true,
			},
			"placeholder": {
				Type:        schema.TypeString,
				Description: "Text in the field used as a placeholder for example purposes",
				Optional:    true,
				Computed:    true,
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
			"verify_pattern": {
				Type:        schema.TypeString,
				Description: "The regex pattern used to validate the entered",
				Optional:    true,
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceOptionTypeTextCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

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

	var dependentField string
	if v, ok := d.Get("dependent_field").(string); ok {
		dependentField = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("dependent_field", d.Get("dependent_field")))
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

	var defaultValue string
	if v, ok := d.Get("default_value").(string); ok {
		defaultValue = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("default_value", d.Get("default_value")))
	}

	var verifyPattern string
	if v, ok := d.Get("verify_pattern").(string); ok {
		verifyPattern = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("verify_pattern", d.Get("verify_pattern")))
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"optionType": map[string]any{
				"name":                  name,
				"description":           description,
				"labels":                labelsPayload,
				"fieldName":             fieldName,
				"exportMeta":            d.Get("export_meta"),
				"dependsOnCode":         dependentField,
				"visibleOnCode":         d.Get("visibility_field"),
				"requireOnCode":         requireField,
				"showOnEdit":            showOnEdit,
				"editable":              editable,
				"displayValueOnDetails": d.Get("display_value_on_details"),
				"type":                  "text",
				"fieldLabel":            d.Get("field_label"),
				"placeHolder":           d.Get("placeholder"),
				"defaultValue":          defaultValue,
				"helpBlock":             d.Get("help_block"),
				"required":              d.Get("required"),
				"verifyPattern":         verifyPattern,
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
	d.SetId(convert.Int64ToString(environment.ID))

	diags = append(diags, resourceOptionTypeTextRead(ctx, d, meta)...)

	return diags
}

func resourceOptionTypeTextRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

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
	d.Set("placeholder", optionType.PlaceHolder)
	d.Set("default_value", optionType.DefaultValue)
	d.Set("help_block", optionType.HelpBlock)
	d.Set("required", optionType.Required)
	d.Set("verify_pattern", optionType.VerifyPattern)

	return diags
}

func resourceOptionTypeTextUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	var dependentField string
	if v, ok := d.Get("dependent_field").(string); ok {
		dependentField = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("dependent_field", d.Get("dependent_field")))
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

	var defaultValue string
	if v, ok := d.Get("default_value").(string); ok {
		defaultValue = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("default_value", d.Get("default_value")))
	}

	var verifyPattern string
	if v, ok := d.Get("verify_pattern").(string); ok {
		verifyPattern = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("verify_pattern", d.Get("verify_pattern")))
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"optionType": map[string]any{
				"name":                  name,
				"description":           description,
				"labels":                labelsPayload,
				"fieldName":             fieldName,
				"exportMeta":            d.Get("export_meta"),
				"dependsOnCode":         dependentField,
				"visibleOnCode":         d.Get("visibility_field"),
				"requireOnCode":         requireField,
				"showOnEdit":            showOnEdit,
				"editable":              editable,
				"displayValueOnDetails": d.Get("display_value_on_details"),
				"type":                  "text",
				"fieldLabel":            d.Get("field_label"),
				"placeHolder":           d.Get("placeholder"),
				"defaultValue":          defaultValue,
				"helpBlock":             d.Get("help_block"),
				"required":              d.Get("required"),
				"verifyPattern":         verifyPattern,
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
	d.SetId(convert.Int64ToString(account.ID))

	return resourceOptionTypeTextRead(ctx, d, meta)
}

func resourceOptionTypeTextDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

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

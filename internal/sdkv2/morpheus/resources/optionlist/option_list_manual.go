// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package optionlist

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func ResourceOptionListManual() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus manual option list resource.",
		CreateContext: resourceOptionListManualCreate,
		ReadContext:   resourceOptionListManualRead,
		UpdateContext: resourceOptionListManualUpdate,
		DeleteContext: resourceOptionListManualDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the manual option list",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the option list",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the option list",
				Optional:    true,
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "The organization labels associated with the option list (Only supported on Morpheus 5.5.3 or higher)",
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"visibility": {
				Type:         schema.TypeString,
				Description:  "Whether the option list is visible in sub-tenants or not",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"private", "public"}, false),
				Computed:     true,
			},
			"dataset": {
				Type:        schema.TypeString,
				Description: "The dataset for the manual option list",
				Optional:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					oldPayload := strings.TrimSpace(old)
					newPayload := strings.TrimSpace(new)

					return oldPayload == newPayload
				},
				StateFunc: func(val any) string {
					return strings.TrimSpace(val.(string))
				},
			},
			"real_time": {
				Type:        schema.TypeBool,
				Description: "Whether the list is refreshed every time an associated option type is requested",
				Optional:    true,
				Default:     false,
			},
			"translation_script": {
				Type: schema.TypeString,
				Description: "A js script to translate the result data object into " +
					"an Array containing objects with properties 'name' and 'value'.",
				DiffSuppressFunc: helpers.SuppressEquivalentJSONDiffs,
				Optional:         true,
				Computed:         true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceOptionListManualCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
					return diag.FromErr(helpers.TypeAssertFailError("labels element", s))
				}
			}
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("labels", attr))
		}
	}

	var visibility string
	if v, ok := d.Get("visibility").(string); ok {
		visibility = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("visibility", d.Get("visibility")))
	}

	var dataset string
	if v, ok := d.Get("dataset").(string); ok {
		dataset = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("dataset", d.Get("dataset")))
	}

	var realTime bool
	if v, ok := d.Get("real_time").(bool); ok {
		realTime = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("real_time", d.Get("real_time")))
	}

	var translationScript string
	if v, ok := d.Get("translation_script").(string); ok {
		translationScript = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("translation_script", d.Get("translation_script")))
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"optionTypeList": map[string]any{
				"name":              name,
				"description":       description,
				"labels":            labelsPayload,
				"type":              "manual",
				"visibility":        visibility,
				"initialDataset":    dataset,
				"realTime":          realTime,
				"translationScript": translationScript,
			},
		},
	}
	resp, err := client.CreateOptionList(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.CreateOptionListResult
	if v, ok := resp.Result.(*morpheus.CreateOptionListResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.OptionList == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("OptionList"))
	}

	optionList := result.OptionList
	// Successfully created resource, now set id
	d.SetId(convert.Int64ToString(optionList.ID))

	return resourceOptionListManualRead(ctx, d, meta)
}

func resourceOptionListManualRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
		resp, err = client.FindOptionListByName(name)
	} else if id != "" {
		resp, err = client.GetOptionList(convert.StringToInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Option list cannot be read without name or id")
	}

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)
			log.Printf("Forcing recreation of resource")
			d.SetId("")

			return diags
		}

		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	// store resource data
	var result *morpheus.GetOptionListResult
	if v, ok := resp.Result.(*morpheus.GetOptionListResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.OptionList == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("OptionList"))
	}

	optionList := result.OptionList
	d.SetId(convert.Int64ToString(optionList.ID))
	d.Set("name", optionList.Name)
	d.Set("description", optionList.Description)

	if optionList.Labels != nil {
		d.Set("labels", optionList.Labels)
	}

	d.Set("visibility", optionList.Visibility)
	d.Set("dataset", optionList.InitialDataset)
	d.Set("real_time", optionList.RealTime)
	d.Set("translation_script", optionList.TranslationScript)

	return diags
}

func resourceOptionListManualUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
					return diag.FromErr(helpers.TypeAssertFailError("labels", s))
				}
			}
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("labels", attr))
		}
	}

	var visibility string
	if v, ok := d.Get("visibility").(string); ok {
		visibility = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("visibility", d.Get("visibility")))
	}

	var dataset string
	if v, ok := d.Get("dataset").(string); ok {
		dataset = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("dataset", d.Get("dataset")))
	}

	var realTime bool
	if v, ok := d.Get("real_time").(bool); ok {
		realTime = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("real_time", d.Get("real_time")))
	}

	var translationScript string
	if v, ok := d.Get("translation_script").(string); ok {
		translationScript = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("translation_script", d.Get("translation_script")))
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"optionTypeList": map[string]any{
				"name":              name,
				"description":       description,
				"labels":            labelsPayload,
				"type":              "manual",
				"visibility":        visibility,
				"initialDataset":    dataset,
				"realTime":          realTime,
				"translationScript": translationScript,
			},
		},
	}
	resp, err := client.UpdateOptionList(convert.StringToInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.UpdateOptionListResult
	if v, ok := resp.Result.(*morpheus.UpdateOptionListResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.OptionList == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("OptionList"))
	}

	optionListResult := result.OptionList
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(convert.Int64ToString(optionListResult.ID))

	return resourceOptionListManualRead(ctx, d, meta)
}

func resourceOptionListManualDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
	resp, err := client.DeleteOptionList(convert.StringToInt64(id), req)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)

			return diag.FromErr(err)
		}

		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	d.SetId("")

	return diags
}

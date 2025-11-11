// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package optionlist

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func ResourceOptionListAPI() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus api option list resource.",
		CreateContext: resourceOptionListAPICreate,
		ReadContext:   resourceOptionListAPIRead,
		UpdateContext: resourceOptionListAPIUpdate,
		DeleteContext: resourceOptionListAPIDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the api option list",
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
				ValidateFunc: validation.StringInSlice([]string{"private", "public", ""}, false),
				Default:      "private",
			},
			"option_list": {
				Type: schema.TypeString,
				Description: "The Morpheus object option list " +
					"(clouds, instanceTypeClouds, instanceTypeLayouts, environments, groups, instances, instance-wiki, " +
					"networks, instanceNetworks, servicePlans, resourcePools, securityGroups, servers, server-wiki)",
				ValidateFunc: validation.StringInSlice(
					[]string{
						"clouds", "instanceTypeClouds", "instanceTypeLayouts", "environments", "groups",
						"instances", "instance-wiki", "networks", "instanceNetworks", "servicePlans",
						"resourcePools", "securityGroups", "servers", "server-wiki",
					},
					false,
				),
				Optional: true,
				Computed: true,
			},
			"translation_script": {
				Type: schema.TypeString,
				Description: "A js script to translate the result data object into an Array " +
					"containing objects with properties 'name' and 'value'.",
				DiffSuppressFunc: helpers.SuppressEquivalentJSONDiffs,
				Optional:         true,
				Computed:         true,
			},
			"request_script": {
				Type:             schema.TypeString,
				Description:      "A js script to manipulate the request payload.",
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

func resourceOptionListAPICreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
				}
			}
		}
	}

	var optionList string
	if v, ok := d.Get("option_list").(string); ok {
		optionList = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("option_list", d.Get("option_list")))
	}

	var visibility string
	if v, ok := d.Get("visibility").(string); ok {
		visibility = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("visibility", d.Get("visibility")))
	}

	var translationScript string
	if v, ok := d.Get("translation_script").(string); ok {
		translationScript = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("translation_script", d.Get("translation_script")))
	}

	var requestScript string
	if v, ok := d.Get("request_script").(string); ok {
		requestScript = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("request_script", d.Get("request_script")))
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"optionTypeList": map[string]any{
				"name":              name,
				"description":       description,
				"labels":            labelsPayload,
				"type":              "api",
				"apiType":           optionList,
				"visibility":        visibility,
				"translationScript": translationScript,
				"requestScript":     requestScript,
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

	optionListResult := result.OptionList
	d.SetId(convert.Int64ToString(optionListResult.ID))

	return resourceOptionListAPIRead(ctx, d, meta)
}

func resourceOptionListAPIRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
	d.Set("option_list", optionList.APIType)
	d.Set("translation_script", optionList.TranslationScript)
	d.Set("request_script", optionList.RequestScript)

	return diags
}

func resourceOptionListAPIUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
				}
			}
		}
	}

	var optionList string
	if v, ok := d.Get("option_list").(string); ok {
		optionList = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("option_list", d.Get("option_list")))
	}

	var visibility string
	if v, ok := d.Get("visibility").(string); ok {
		visibility = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("visibility", d.Get("visibility")))
	}

	var translationScript string
	if v, ok := d.Get("translation_script").(string); ok {
		translationScript = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("translation_script", d.Get("translation_script")))
	}

	var requestScript string
	if v, ok := d.Get("request_script").(string); ok {
		requestScript = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("request_script", d.Get("request_script")))
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"optionTypeList": map[string]any{
				"name":              name,
				"description":       description,
				"labels":            labelsPayload,
				"type":              "api",
				"apiType":           optionList,
				"visibility":        visibility,
				"translationScript": translationScript,
				"requestScript":     requestScript,
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
	d.SetId(convert.Int64ToString(optionListResult.ID))

	return resourceOptionListAPIRead(ctx, d, meta)
}

func resourceOptionListAPIDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

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

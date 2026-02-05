// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package environment

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/convert"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"
)

func ResourceEnvironment() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus environment resource",
		CreateContext: resourceEnvironmentCreate,
		ReadContext:   resourceEnvironmentRead,
		UpdateContext: resourceEnvironmentUpdate,
		DeleteContext: resourceEnvironmentDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the environment",
				Computed:    true,
			},
			"active": {
				Type:        schema.TypeBool,
				Description: "Whether the environment is enabled or not",
				Optional:    true,
				Default:     true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the environment",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the environment",
				Optional:    true,
				Computed:    true,
			},
			"code": {
				Type:        schema.TypeString,
				Description: "The code of the environment",
				Required:    true,
			},
			"visibility": {
				Type:        schema.TypeString,
				Description: "Whether the environment is visible in sub-tenants or not",
				Optional:    true,
				Default:     "private",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceEnvironmentCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	var active bool
	if v, ok := d.Get("active").(bool); ok {
		active = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("active", d.Get("active")))
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

	var code string
	if v, ok := d.Get("code").(string); ok {
		code = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("code", d.Get("code")))
	}

	var visibility string
	if v, ok := d.Get("visibility").(string); ok {
		visibility = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("visibility", d.Get("visibility")))
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"environment": map[string]any{
				"active":      active,
				"name":        name,
				"description": description,
				"code":        code,
				"visibility":  visibility,
			},
		},
	}

	resp, err := client.CreateEnvironment(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}

	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.CreateEnvironmentResult
	if v, ok := resp.Result.(*morpheus.CreateEnvironmentResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.Environment == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Environment"))
	}

	environment := result.Environment
	d.SetId(convert.Int64ToString(environment.ID))

	diags = append(diags, resourceEnvironmentRead(ctx, d, meta)...)

	return diags
}

func resourceEnvironmentRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
		resp, err = client.FindEnvironmentByName(name)
	} else if id != "" {
		resp, err = client.GetEnvironment(convert.StringToInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Environment cannot be read without name or id")
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

	var result *morpheus.GetEnvironmentResult
	if v, ok := resp.Result.(*morpheus.GetEnvironmentResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.Environment == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Environment"))
	}

	environment := result.Environment
	d.SetId(convert.Int64ToString(environment.ID))
	d.Set("active", environment.Active)
	d.Set("name", environment.Name)
	d.Set("description", environment.Description)
	d.Set("visibility", environment.Visibility)
	d.Set("code", environment.Code)

	return diags
}

func resourceEnvironmentUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	id := d.Id()

	var active bool
	if v, ok := d.Get("active").(bool); ok {
		active = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("active", d.Get("active")))
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

	var code string
	if v, ok := d.Get("code").(string); ok {
		code = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("code", d.Get("code")))
	}

	var visibility string
	if v, ok := d.Get("visibility").(string); ok {
		visibility = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("visibility", d.Get("visibility")))
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"environment": map[string]any{
				"active":      active,
				"name":        name,
				"description": description,
				"code":        code,
				"visibility":  visibility,
			},
		},
	}
	resp, err := client.UpdateEnvironment(convert.StringToInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.UpdateEnvironmentResult
	if v, ok := resp.Result.(*morpheus.UpdateEnvironmentResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.Environment == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Environment"))
	}

	environment := result.Environment
	d.SetId(convert.Int64ToString(environment.ID))

	return resourceEnvironmentRead(ctx, d, meta)
}

func resourceEnvironmentDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteEnvironment(convert.StringToInt64(id), req)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)

			return nil
		}

		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	d.SetId("")

	return diags
}

// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package license

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/convert"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"
)

func ResourceLicense() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus license resource.",
		CreateContext: resourceLicenseCreate,
		ReadContext:   resourceLicenseRead,
		UpdateContext: resourceLicenseUpdate,
		DeleteContext: resourceLicenseDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the license",
				Computed:    true,
			},
			"key": {
				Type:        schema.TypeString,
				Description: "The Morpheus license key",
				Required:    true,
				Sensitive:   true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceLicenseCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	var key string
	if v, ok := d.Get("key").(string); ok {
		key = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("key", d.Get("key")))
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"license": key,
		},
	}

	resp, err := client.InstallLicense(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.GetLicenseResult
	if v, ok := resp.Result.(*morpheus.GetLicenseResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}
	_ = result

	d.SetId(convert.Int64ToString(1))

	diags = append(diags, resourceLicenseRead(ctx, d, meta)...)

	return diags
}

func resourceLicenseRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	var resp *morpheus.Response
	var err error

	resp, err = client.GetLicense(&morpheus.Request{})
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

	var result *morpheus.GetLicenseResult
	if v, ok := resp.Result.(*morpheus.GetLicenseResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}
	_ = result

	d.SetId(convert.Int64ToString(1))

	return diags
}

func resourceLicenseUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var key string
	if v, ok := d.Get("key").(string); ok {
		key = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("key", d.Get("key")))
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"license": key,
		},
	}

	log.Printf("API REQUEST: %s", req)
	resp, err := client.InstallLicense(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.GetLicenseResult
	if v, ok := resp.Result.(*morpheus.GetLicenseResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}
	_ = result

	d.SetId(convert.Int64ToString(1))

	return resourceLicenseRead(ctx, d, meta)
}

func resourceLicenseDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	req := &morpheus.Request{}
	resp, err := client.UninstallLicense(1, req)
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

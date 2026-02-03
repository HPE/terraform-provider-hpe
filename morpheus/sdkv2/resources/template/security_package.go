// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package template

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/convert"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"
)

func ResourceSecurityPackage() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus security package resource",
		CreateContext: resourceSecurityPackageCreate,
		ReadContext:   resourceSecurityPackageRead,
		UpdateContext: resourceSecurityPackageUpdate,
		DeleteContext: resourceSecurityPackageDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the security package",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the security package",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the security package",
				Optional:    true,
				Computed:    true,
			},
			"labels": {
				Type: schema.TypeSet,
				Description: "The organization labels associated with the security package " +
					"(Only supported on Morpheus 5.5.3 or higher)",
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the security package is enabled",
				Optional:    true,
				Default:     true,
			},
			"url": {
				Type:        schema.TypeString,
				Description: "The source url of the security package",
				Required:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceSecurityPackageCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	securityPackage := make(map[string]any)

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}
	securityPackage["name"] = name
	securityPackage["type"] = "scap-package"

	var description string
	if v, ok := d.Get("description").(string); ok {
		description = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("description", d.Get("description")))
	}
	securityPackage["description"] = description

	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		if labelSet, ok := attr.(*schema.Set); ok {
			if labelSet != nil {
				for _, s := range labelSet.List() {
					if labelStr, ok := s.(string); ok {
						labelsPayload = append(labelsPayload, labelStr)
					}
				}
			}
		}
	}
	securityPackage["labels"] = labelsPayload

	var enabled bool
	if v, ok := d.Get("enabled").(bool); ok {
		enabled = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("enabled", d.Get("enabled")))
	}
	securityPackage["enabled"] = enabled

	var url string
	if v, ok := d.Get("url").(string); ok {
		url = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("url", d.Get("url")))
	}
	securityPackage["url"] = url

	req := &morpheus.Request{
		Body: map[string]any{
			"securityPackage": securityPackage,
		},
	}
	resp, err := client.CreateSecurityPackage(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.CreateSecurityPackageResult
	if v, ok := resp.Result.(*morpheus.CreateSecurityPackageResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.SecurityPackage == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("SecurityPackage"))
	}

	securityPackageResult := result.SecurityPackage
	d.SetId(convert.Int64ToString(securityPackageResult.ID))

	diags = append(diags, resourceSecurityPackageRead(ctx, d, meta)...)

	return diags
}

func resourceSecurityPackageRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
		resp, err = client.FindSecurityPackageByName(name)
	} else if id != "" {
		resp, err = client.GetSecurityPackage(convert.StringToInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Security package cannot be read without name or id")
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

	var result *morpheus.GetSecurityPackageResult
	if v, ok := resp.Result.(*morpheus.GetSecurityPackageResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.SecurityPackage == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("SecurityPackage"))
	}

	securityPackage := result.SecurityPackage

	d.SetId(convert.IntToString(int(securityPackage.ID)))
	d.Set("name", securityPackage.Name)
	d.Set("description", securityPackage.Description)
	d.Set("labels", securityPackage.Labels)
	d.Set("enabled", securityPackage.Enabled)
	d.Set("url", securityPackage.Url)

	return diags
}

func resourceSecurityPackageUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	id := d.Id()
	securityPackage := make(map[string]any)

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}
	securityPackage["name"] = name

	var description string
	if v, ok := d.Get("description").(string); ok {
		description = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("description", d.Get("description")))
	}
	securityPackage["description"] = description

	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		if labelSet, ok := attr.(*schema.Set); ok {
			if labelSet != nil {
				for _, s := range labelSet.List() {
					if labelStr, ok := s.(string); ok {
						labelsPayload = append(labelsPayload, labelStr)
					} else {
						return diag.FromErr(helpers.TypeAssertFailError("label", s))
					}
				}
			}
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("labels", d.Get("labels")))
		}
	}
	securityPackage["labels"] = labelsPayload

	var enabled bool
	if v, ok := d.Get("enabled").(bool); ok {
		enabled = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("enabled", d.Get("enabled")))
	}
	securityPackage["enabled"] = enabled

	var url string
	if v, ok := d.Get("url").(string); ok {
		url = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("url", d.Get("url")))
	}
	securityPackage["url"] = url

	req := &morpheus.Request{
		Body: map[string]any{
			"securityPackage": securityPackage,
		},
	}
	resp, err := client.UpdateSecurityPackage(convert.StringToInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.UpdateSecurityPackageResult
	if v, ok := resp.Result.(*morpheus.UpdateSecurityPackageResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.SecurityPackage == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("SecurityPackage"))
	}

	securityPackageResult := result.SecurityPackage
	d.SetId(convert.Int64ToString(securityPackageResult.ID))

	return resourceSecurityPackageRead(ctx, d, meta)
}

func resourceSecurityPackageDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteSecurityPackage(convert.StringToInt64(id), req)
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

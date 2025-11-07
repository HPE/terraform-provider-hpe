// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package wiki

import (
	"context"
	"log"
	"strings"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"
)

func ResourceWikiPage() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus wiki page resource",
		CreateContext: resourceWikiPageCreate,
		ReadContext:   resourceWikiPageRead,
		UpdateContext: resourceWikiPageUpdate,
		DeleteContext: resourceWikiPageDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the wiki page",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the wiki page",
				Required:    true,
			},
			"category": {
				Type:        schema.TypeString,
				Description: "The category of the wiki page",
				Optional:    true,
				Computed:    true,
			},
			"content": {
				Type:        schema.TypeString,
				Description: "The content of the wiki page",
				Optional:    true,
				StateFunc: func(v any) string {
					var payload string
					if s, ok := v.(string); ok {
						payload = strings.TrimSuffix(s, "\n")
					}

					return payload
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceWikiPageCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("meta", meta))
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	wikiPage := make(map[string]any)

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", meta))
	}
	wikiPage["name"] = name

	var category string
	if v, ok := d.Get("category").(string); ok {
		category = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("category", meta))
	}
	wikiPage["category"] = category

	var content string
	if v, ok := d.Get("content").(string); ok {
		content = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("content", meta))
	}
	wikiPage["content"] = content

	req := &morpheus.Request{
		Body: map[string]any{
			"page": wikiPage,
		},
	}
	resp, err := client.CreateWiki(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	var result *morpheus.CreateWikiResult
	if v, ok := resp.Result.(*morpheus.CreateWikiResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("resp.Result", meta))
	}
	if result.Wiki == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Wiki"))
	}
	wikiPageResult := result.Wiki
	// Successfully created resource, now set id
	d.SetId(convert.Int64ToString(wikiPageResult.ID))

	resourceWikiPageRead(ctx, d, meta)

	return diags
}

func resourceWikiPageRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("meta", meta))
	}
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", meta))
	}

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindWikiByName(name)
	} else if id != "" {
		resp, err = client.GetWiki(convert.StringToInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Wiki Page cannot be read without name or id")
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
	// store resource data
	var result *morpheus.GetWikiResult
	if v, ok := resp.Result.(*morpheus.GetWikiResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("resp.Result", meta))
	}
	if result.Wiki == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Wiki"))
	}
	wikiPage := result.Wiki

	d.SetId(convert.IntToString(int(wikiPage.ID)))
	d.Set("name", wikiPage.Name)
	d.Set("category", wikiPage.Category)
	d.Set("content", wikiPage.Content)

	return diags
}

func resourceWikiPageUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("meta", meta))
	}
	id := d.Id()

	wikiPage := make(map[string]any)

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", meta))
	}
	wikiPage["name"] = name

	var category string
	if v, ok := d.Get("category").(string); ok {
		category = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("category", meta))
	}
	wikiPage["category"] = category

	var content string
	if v, ok := d.Get("content").(string); ok {
		content = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("content", meta))
	}
	wikiPage["content"] = content

	req := &morpheus.Request{
		Body: map[string]any{
			"page": wikiPage,
		},
	}
	resp, err := client.UpdateWiki(convert.StringToInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	var result *morpheus.UpdateWikiResult
	if v, ok := resp.Result.(*morpheus.UpdateWikiResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("resp.Result", meta))
	}
	if result.Wiki == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Wiki"))
	}
	wikiPageResult := result.Wiki

	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(convert.Int64ToString(wikiPageResult.ID))

	return resourceWikiPageRead(ctx, d, meta)
}

func resourceWikiPageDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("meta", meta))
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteWiki(convert.StringToInt64(id), req)
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

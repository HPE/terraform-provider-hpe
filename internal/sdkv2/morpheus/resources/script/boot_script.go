// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package script

import (
	"context"
	"log"
	"strings"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceBootScript() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus boot script resource",
		CreateContext: resourceBootScriptCreate,
		ReadContext:   resourceBootScriptRead,
		UpdateContext: resourceBootScriptUpdate,
		DeleteContext: resourceBootScriptDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the boot script",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the boot script",
				Required:    true,
			},
			"content": {
				Type:        schema.TypeString,
				Description: "The content of the boot script",
				Optional:    true,
				StateFunc: func(v any) string {
					payload := strings.TrimSuffix(v.(string), "\n")

					return payload
				},
				Computed: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceBootScriptCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	var content string
	if v, ok := d.Get("content").(string); ok {
		content = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("content", d.Get("content")))
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"bootScript": map[string]any{
				"fileName": name,
				"content":  content,
			},
		},
	}

	resp, err := client.CreateBootScript(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.CreateBootScriptResult
	if v, ok := resp.Result.(*morpheus.CreateBootScriptResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.BootScript == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("BootScript"))
	}

	bootScript := result.BootScript
	if bootScript == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("BootScript"))
	}
	// Successfully created resource, now set id
	d.SetId(convert.Int64ToString(bootScript.ID))

	diags = append(diags, resourceBootScriptRead(ctx, d, meta)...)

	return diags
}

func resourceBootScriptRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
		resp, err = client.FindBootScriptByName(name)
	} else if id != "" {
		resp, err = client.GetBootScript(convert.StringToInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("File template cannot be read without name or id")
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

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.GetBootScriptResult
	if v, ok := resp.Result.(*morpheus.GetBootScriptResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.BootScript == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("BootScript"))
	}

	bootScript := result.BootScript
	if bootScript == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("BootScript"))
	}
	d.SetId(convert.Int64ToString(bootScript.ID))
	d.Set("name", bootScript.FileName)
	d.Set("content", bootScript.Content)

	return diags
}

func resourceBootScriptUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	var content string
	if v, ok := d.Get("content").(string); ok {
		content = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("content", d.Get("content")))
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"bootScript": map[string]any{
				"fileName": name,
				"content":  content,
			},
		},
	}

	resp, err := client.UpdateBootScript(convert.StringToInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.UpdateBootScriptResult
	if v, ok := resp.Result.(*morpheus.UpdateBootScriptResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.BootScript == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("BootScript"))
	}

	bootScript := result.BootScript
	if bootScript == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("BootScript"))
	}
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(convert.Int64ToString(bootScript.ID))

	return resourceBootScriptRead(ctx, d, meta)
}

func resourceBootScriptDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
	resp, err := client.DeleteBootScript(convert.StringToInt64(id), req)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)

			return nil
		} else {
			log.Printf("API FAILURE: %s - %s", resp, err)

			return diag.FromErr(err)
		}
	}
	log.Printf("API RESPONSE: %s", resp)
	d.SetId("")

	return diags
}

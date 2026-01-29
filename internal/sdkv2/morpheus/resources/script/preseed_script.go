// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package script

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func ResourcePreseedScript() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus preseed script resource",
		CreateContext: resourcePreseedScriptCreate,
		ReadContext:   resourcePreseedScriptRead,
		UpdateContext: resourcePreseedScriptUpdate,
		DeleteContext: resourcePreseedScriptDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the preseed script",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the preseed script",
				Required:    true,
			},
			"content": {
				Type:        schema.TypeString,
				Description: "The content of the preseed script",
				Optional:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					oldPayload := strings.TrimSpace(old)
					newPayload := strings.TrimSpace(new)
					return oldPayload == newPayload
				},
				StateFunc: func(v any) string {
					var payload string
					if vStr, ok := v.(string); ok {
						payload = strings.TrimSpace(vStr)
					}

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

func resourcePreseedScriptCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	var content string
	if v, ok := d.Get("content").(string); ok {
		content = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("content", d.Get("content")))
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"preseedScript": map[string]any{
				"fileName": name,
				"content":  content,
			},
		},
	}

	resp, err := client.CreatePreseedScript(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.CreatePreseedScriptResult
	if v, ok := resp.Result.(*morpheus.CreatePreseedScriptResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.PreseedScript == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("PreseedScript"))
	}

	preseedScript := result.PreseedScript
	d.SetId(convert.Int64ToString(preseedScript.ID))

	diags = append(diags, resourcePreseedScriptRead(ctx, d, meta)...)

	return diags
}

func resourcePreseedScriptRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
		resp, err = client.FindPreseedScriptByName(name)
	} else if id != "" {
		resp, err = client.GetPreseedScript(convert.StringToInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Preseed script cannot be read without name or id")
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

	var result *morpheus.GetPreseedScriptResult
	if v, ok := resp.Result.(*morpheus.GetPreseedScriptResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.PreseedScript == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("PreseedScript"))
	}

	preseedScript := result.PreseedScript
	d.SetId(convert.Int64ToString(preseedScript.ID))
	d.Set("name", preseedScript.FileName)
	d.Set("content", preseedScript.Content)

	return diags
}

func resourcePreseedScriptUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
			"preseedScript": map[string]any{
				"fileName": name,
				"content":  content,
			},
		},
	}

	resp, err := client.UpdatePreseedScript(convert.StringToInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.UpdatePreseedScriptResult
	if v, ok := resp.Result.(*morpheus.UpdatePreseedScriptResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.PreseedScript == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("PreseedScript"))
	}

	preseedScript := result.PreseedScript
	d.SetId(convert.Int64ToString(preseedScript.ID))

	return resourcePreseedScriptRead(ctx, d, meta)
}

func resourcePreseedScriptDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeletePreseedScript(convert.StringToInt64(id), req)
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

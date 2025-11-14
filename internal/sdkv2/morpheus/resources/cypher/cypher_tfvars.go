// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package cypher

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func ResourceCypherTFVars() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus cypher tfvars secret resource.",
		CreateContext: resourceCypherTFVarsCreate,
		ReadContext:   resourceCypherTFVarsRead,
		DeleteContext: resourceCypherTFVarsDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the cypher tfvars secret",
				Computed:    true,
			},
			"key": {
				Type:        schema.TypeString,
				Description: "The path of the cypher tfvars secret, excluding the secret prefix",
				Required:    true,
				ForceNew:    true,
			},
			"value": {
				Type:        schema.TypeString,
				Description: "The value of the cypher tfvars secret",
				Required:    true,
				Sensitive:   true,
				ForceNew:    true,
			},
			"ttl": {
				Type:        schema.TypeInt,
				Description: "The time to live of the cypher tfvars secret",
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceCypherTFVarsCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	var value string
	if v, ok := d.Get("value").(string); ok {
		value = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("value", d.Get("value")))
	}

	var ttl int
	if v, ok := d.Get("ttl").(int); ok {
		ttl = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("ttl", d.Get("ttl")))
	}

	var key string
	if v, ok := d.Get("key").(string); ok {
		key = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("key", d.Get("key")))
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"value": value,
		},
		QueryParams: map[string]string{
			"ttl":  convert.IntToString(ttl),
			"type": "string",
		},
	}

	tfvarsPath := fmt.Sprintf("tfvars/%s", key)
	resp, err := client.CreateCypher(tfvarsPath, req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	// Masking to avoid credential exposure
	// log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.CreateCypherResult
	if v, ok := resp.Result.(*morpheus.CreateCypherResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.Cypher == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Cypher"))
	}

	// Successfully created resource, now set id
	d.SetId(convert.Int64ToString(result.Cypher.ID))

	diags = append(diags, resourceCypherTFVarsRead(ctx, d, meta)...)

	return diags
}

func resourceCypherTFVarsRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	id := d.Id()

	var key string
	if v, ok := d.Get("key").(string); ok {
		key = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("key", d.Get("key")))
	}

	var resp *morpheus.Response
	var err error
	if id != "" {
		tfvarsPath := fmt.Sprintf("tfvars/%s", key)
		resp, err = client.GetCypher(tfvarsPath, &morpheus.Request{})
	} else {
		return diag.Errorf("Cypher cannot be read without id")
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
	// Masking to avoid credential exposure
	// log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.GetCypherResult
	if v, ok := resp.Result.(*morpheus.GetCypherResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.Cypher == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Cypher"))
	}

	// store resource data
	d.SetId(convert.Int64ToString(result.Cypher.ID))

	if result.Cypher.ItemKey == "" {
		return diag.FromErr(helpers.NotFoundInResponseError("ItemKey"))
	}

	keyData := strings.Split(result.Cypher.ItemKey, "/")
	if len(keyData) > 0 {
		keyData = keyData[1:]
		d.Set("key", strings.Join(keyData, "/"))
	}

	d.Set("ttl", result.LeaseDuration)

	return diags
}

func resourceCypherTFVarsDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	req := &morpheus.Request{}
	tfvarsPath := fmt.Sprintf("tfvars/%s", key)
	resp, err := client.DeleteCypher(tfvarsPath, req)
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

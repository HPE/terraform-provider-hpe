// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package cypher

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/convert"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"
)

func DataSourceCypherSecret() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus cypher secret data source.",
		ReadContext: dataSourceCypherSecretRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"key": {
				Type:        schema.TypeString,
				Description: "The path of the cypher secret, excluding the secret prefix",
				Required:    true,
			},
			"value": {
				Type:        schema.TypeString,
				Description: "The cypher secret value",
				Computed:    true,
			},
			"ttl": {
				Type:        schema.TypeInt,
				Description: "The time to live of the cypher secret",
				Computed:    true,
			},
		},
	}
}

func dataSourceCypherSecretRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	secretPath := fmt.Sprintf("secret/%s", key)

	resp, err := client.Execute(&morpheus.Request{
		Method: "GET",
		Path:   fmt.Sprintf("%s/%s", "/api/cypher", secretPath),
		Result: &LocalGetCypherResult{},
	})
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %v", resp, err)

			return nil
		}

		log.Printf("API FAILURE: %s - %v", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var cypher *LocalGetCypherResult
	if v, ok := resp.Result.(*LocalGetCypherResult); ok {
		cypher = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if cypher.Cypher == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Cypher"))
	}

	d.SetId(convert.Int64ToString(cypher.Cypher.ID))
	if cypher.Type == "object" {
		jsonPayload, _ := json.Marshal(cypher.Data)
		d.Set("value", string(jsonPayload))
	} else {
		var dataStr string
		if v, ok := cypher.Data.(string); ok {
			dataStr = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("Data", cypher.Data))
		}
		d.Set("value", dataStr)
	}
	d.Set("ttl", cypher.LeaseDuration)

	return diags
}

type LocalGetCypherResult struct {
	Success       bool              `json:"success"`
	Data          any               `json:"data"`
	Type          string            `json:"type"`
	LeaseDuration int64             `json:"lease_duration"`
	Cypher        *morpheus.Cypher  `json:"cypher"`
	Message       string            `json:"msg"`
	Errors        map[string]string `json:"errors"`
}

// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package trust

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/convert"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"
)

func DataSourceKeyPair() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus key pair data source.",
		ReadContext: dataSourceKeyPairRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeInt,
				Description:   "The ID of the key pair",
				Optional:      true,
				ConflictsWith: []string{"name"},
				Computed:      true,
			},
			"name": {
				Type:          schema.TypeString,
				Description:   "The name of the integration",
				Optional:      true,
				ConflictsWith: []string{"id"},
			},
			"publickey": {
				Type:        schema.TypeString,
				Description: "PublicKey of the KeyPair",
				Optional:    true,
			},
		},
	}
}

func dataSourceKeyPairRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}

	id := convert.StringToInt64(d.Id())

	var resp *morpheus.Response
	var err error
	if (id != 0 && name == "") || (id != 0 && name != "") {
		resp, err = client.GetKeyPair(id)
	} else if id == 0 && name != "" {
		resp, err = client.GetKeyPairByName(name)
	} else if id == 0 && name == "" {
		return diag.Errorf("Key pair cannot be read without name or id")
	}

	if err != nil {
		return diag.FromErr(err)
	}

	var keyPair *morpheus.KeyPair
	if id != 0 {
		var result *morpheus.GetKeyPairResult
		if v, ok := resp.Result.(*morpheus.GetKeyPairResult); ok {
			result = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("GetKeyPairResult", resp.Result))
		}

		if result.KeyPair == nil {
			return diag.FromErr(helpers.NotFoundInResponseError("KeyPair"))
		}

		keyPair = result.KeyPair
	} else if name != "" {
		var listResult *morpheus.ListKeyPairsResult
		if v, ok := resp.Result.(*morpheus.ListKeyPairsResult); ok {
			listResult = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("ListKeyPairsResult", resp.Result))
		}

		if listResult.KeyPairs == nil {
			return diag.FromErr(helpers.NotFoundInResponseError("KeyPairs"))
		}

		keyPairs := listResult.KeyPairs
		if len(*keyPairs) == 0 {
			return diag.FromErr(helpers.EmptySliceError("KeyPairs"))
		}

		keyPair = &(*keyPairs)[0]
	}

	if keyPair != nil {
		d.SetId(convert.Int64ToString((*keyPair).ID))
		d.Set("name", (*keyPair).Name)
		d.Set("publickey", (*keyPair).PublicKey)
	} else {
		return diag.Errorf("Key pair not found in response data.")
	}

	return diags
}

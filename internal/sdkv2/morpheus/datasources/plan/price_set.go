// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package plan

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func DataSourcePriceSet() *schema.Resource {
	return &schema.Resource{
		Description: "The Price Set data source allows details of a Price Set to be retrieved by its name.",
		ReadContext: dataSourcePriceSetRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"name"},
				Computed:      true,
			},
			"name": {
				Type:          schema.TypeString,
				Description:   "The name of the Morpheus price set.",
				Optional:      true,
				ConflictsWith: []string{"id"},
			},
			"code": {
				Type:        schema.TypeString,
				Description: "The code of the Morpheus price set",
				Computed:    true,
			},
		},
	}
}

func dataSourcePriceSetRead(
	ctx context.Context,
	d *schema.ResourceData,
	meta any,
) diag.Diagnostics {
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

	var id int
	if v, ok := d.Get("id").(int); ok {
		id = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("id", d.Get("id")))
	}

	var resp *morpheus.Response
	var err error
	if id == 0 && name != "" {
		resp, err = client.FindPriceSetByName(name)
	} else if id != 0 {
		resp, err = client.GetPriceSet(int64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Price set cannot be read without name or id")
	}

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

	var result *morpheus.GetPriceSetResult
	if v, ok := resp.Result.(*morpheus.GetPriceSetResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.PriceSet == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("PriceSet"))
	}

	priceSet := result.PriceSet
	d.SetId(convert.Int64ToString(priceSet.ID))
	d.Set("name", priceSet.Name)
	d.Set("code", priceSet.Code)

	return diags
}

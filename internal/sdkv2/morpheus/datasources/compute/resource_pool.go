// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package compute

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"
)

func DataSourceResourcePool() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus resource pool data source.",
		ReadContext: dataSourceResourcePoolRead,
		Schema: map[string]*schema.Schema{
			"cloud_id": {
				Type:        schema.TypeInt,
				Description: "The id of the Morpheus cloud to search for the resource pool.",
				Required:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the Morpheus resource pool.",
				Optional:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "Optional code for use with policies",
				Computed:    true,
			},
			"active": {
				Type:        schema.TypeBool,
				Description: "Whether the resource pool is enabled or not",
				Computed:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the resource pool",
				Computed:    true,
			},
			"id": {
				Type:        schema.TypeInt,
				Description: "The id of the resource pool",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func dataSourceResourcePoolRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	var id int
	if v, ok := d.Get("id").(int); ok {
		id = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("id", d.Get("id")))
	}

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}

	var cloudID int
	if v, ok := d.Get("cloud_id").(int); ok {
		cloudID = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("cloud_id", d.Get("cloud_id")))
	}

	// Ensure that either the id or name is provided
	if id == 0 && name == "" {
		return diag.Errorf(
			"Either 'id' or 'name' must be provided to search for the resource pool",
		)
	}

	var resp *morpheus.Response
	var err error

	if id != 0 {
		resp, err = client.GetResourcePool(
			int64(cloudID),
			int64(id),
			&morpheus.Request{},
		)
	} else {
		resp, err = client.FindResourcePoolByName(int64(cloudID), name)
	}

	if err != nil {
		errorPrefix := "API FAILURE"
		if resp != nil && resp.StatusCode == 404 {
			errorPrefix = "API 404"
		}
		log.Printf("%s: %s - %v", errorPrefix, resp, err)

		return diag.FromErr(err)
	}

	log.Printf("API RESPONSE: %s", resp)

	var result *morpheus.GetResourcePoolResult
	if v, ok := resp.Result.(*morpheus.GetResourcePoolResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.ResourcePool == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("ResourcePool"))
	}

	resourcePool := result.ResourcePool

	d.SetId(convert.Int64ToString(resourcePool.ID))
	d.Set("name", resourcePool.Name)
	d.Set("active", resourcePool.Active)
	d.Set("type", resourcePool.Type)
	d.Set("description", resourcePool.Description)

	return diags
}

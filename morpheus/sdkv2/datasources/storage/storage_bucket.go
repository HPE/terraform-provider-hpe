// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package storage

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/convert"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"
)

func DataSourceStorageBucket() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus storage bucket data source.",
		ReadContext: dataSourceStorageBucketRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeInt,
				Description:   "The ID of the storage bucket",
				Optional:      true,
				ConflictsWith: []string{"name"},
				Computed:      true,
			},
			"name": {
				Type:          schema.TypeString,
				Description:   "The name of the storage bucket",
				Optional:      true,
				ConflictsWith: []string{"id"},
			},
		},
	}
}

func dataSourceStorageBucketRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	var id int
	if v, ok := d.Get("id").(int); ok {
		id = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("id", d.Get("id")))
	}

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == 0 && name != "" {
		resp, err = client.FindStorageBucketByName(name)
	} else if id != 0 {
		resp, err = client.GetStorageBucket(int64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Storage Bucket cannot be read without name or id")
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

	// store resource data
	var result *morpheus.GetStorageBucketResult
	if v, ok := resp.Result.(*morpheus.GetStorageBucketResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("GetStorageBucketResult", resp.Result))
	}

	if result.StorageBucket == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("StorageBucket"))
	}
	storageBucket := result.StorageBucket

	d.SetId(convert.Int64ToString(storageBucket.ID))
	d.Set("name", storageBucket.Name)

	return diags
}

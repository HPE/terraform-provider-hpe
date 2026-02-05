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

func DataSourceStorageVolumeType() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus storage volume type data source.",
		ReadContext: dataSourceStorageVolumeTypeRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeInt,
				Description:   "The ID of the stroage volume type",
				Optional:      true,
				ConflictsWith: []string{"name"},
				Computed:      true,
			},
			"name": {
				Type:          schema.TypeString,
				Description:   "The name of the storage volume type",
				Optional:      true,
				ConflictsWith: []string{"id"},
			},
			"code": {
				Type:        schema.TypeString,
				Description: "The code of the storage volume type",
				Optional:    true,
				Computed:    true,
			},
			"category": {
				Type:        schema.TypeString,
				Description: "The category of the storage volume type",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func dataSourceStorageVolumeTypeRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
		resp, err = client.FindStorageVolumeTypeByName(name)
	} else if id != 0 {
		resp, err = client.GetStorageVolumeType(int64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Storage volume type cannot be read without name or id")
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

	var result *morpheus.GetStorageVolumeTypeResult
	if v, ok := resp.Result.(*morpheus.GetStorageVolumeTypeResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.StorageVolumeType == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("StorageVolumeType"))
	}

	storageVolumeType := result.StorageVolumeType
	d.SetId(convert.Int64ToString(storageVolumeType.ID))
	d.Set("name", storageVolumeType.Name)
	d.Set("code", storageVolumeType.Code)
	d.Set("category", storageVolumeType.Category)

	return diags
}

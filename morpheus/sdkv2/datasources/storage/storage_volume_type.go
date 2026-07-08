// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package storage

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	morpheus "github.com/HPE/terraform-provider-hpe/internal/sdk/legacy"

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
				Description:   "The ID of the storage volume type",
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name", "code"},
			},
			"name": {
				Type:          schema.TypeString,
				Description:   "The name of the storage volume type",
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"id", "code"},
			},
			"code": {
				Type:          schema.TypeString,
				Description:   "The code of the storage volume type",
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"id", "name"},
			},
			"category": {
				Type:        schema.TypeString,
				Description: "The category of the storage volume type",
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

	var code string
	if v, ok := d.Get("code").(string); ok {
		code = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("code", d.Get("code")))
	}

	var id int
	if v, ok := d.Get("id").(int); ok {
		id = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("id", d.Get("id")))
	}

	var resp *morpheus.Response
	var err error
	switch {
	case id != 0:
		resp, err = client.GetStorageVolumeType(int64(id), &morpheus.Request{})
	case name != "":
		resp, err = client.FindStorageVolumeTypeByName(name)
	case code != "":
		resp, err = client.FindStorageVolumeTypeByCode(code)
	default:
		return diag.Errorf("Storage volume type cannot be read without name, code, or id")
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

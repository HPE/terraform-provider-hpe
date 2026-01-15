// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package storage

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"
)

func DataSourceStorageVolume() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus storage volume data source.",
		ReadContext: dataSourceStorageVolumeRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Description: "The ID of the storage volume",
				Required:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the storage volume",
				Computed:    true,
			},
			"active": {
				Type:        schema.TypeBool,
				Description: "Whether the storage volume is enabled or not",
				Computed:    true,
			},
			"category": {
				Type:        schema.TypeString,
				Description: "The storage volume category",
				Computed:    true,
			},
			"cloud_name": {
				Type:        schema.TypeString,
				Description: "The storage volume cloud name",
				Computed:    true,
			},
			"cloud_id": {
				Type:        schema.TypeInt,
				Description: "The storage volume cloud id",
				Computed:    true,
			},
			"datastore_name": {
				Type:        schema.TypeString,
				Description: "The storage volume datastore name",
				Computed:    true,
			},
			"datastore_id": {
				Type:        schema.TypeInt,
				Description: "The storage volume datastore id",
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "The status of the storage volume",
				Computed:    true,
			},
			"source": {
				Type:        schema.TypeString,
				Description: "The associated cloud name",
				Computed:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "The storage volume type name",
				Computed:    true,
			},
			"type_id": {
				Type:        schema.TypeInt,
				Description: "The storage volume type id",
				Computed:    true,
			},
			"uuid": {
				Type:        schema.TypeString,
				Description: "The storage volume uuid",
				Computed:    true,
			},
		},
	}
}

func dataSourceStorageVolumeRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	var id int
	if v, ok := d.Get("id").(int); ok {
		id = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("id", d.Get("id")))
	}

	var resp *morpheus.Response
	var err error

	resp, err = client.GetStorageVolume(int64(id), &morpheus.Request{})
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

	var result *morpheus.GetStorageVolumeResult
	if v, ok := resp.Result.(*morpheus.GetStorageVolumeResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.StorageVolume == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("StorageVolume"))
	}

	storageVolume := result.StorageVolume

	d.SetId(convert.Int64ToString(storageVolume.ID))
	d.Set("name", storageVolume.Name)
	d.Set("active", storageVolume.Active)
	d.Set("category", storageVolume.Category)
	d.Set("cloud_name", storageVolume.Zone.Name)
	d.Set("cloud_id", storageVolume.Zone.ID)
	d.Set("datastore_id", storageVolume.Datastore.ID)
	d.Set("datastore_name", storageVolume.Datastore.Name)
	d.Set("status", storageVolume.Status)
	d.Set("source", storageVolume.Source)
	d.Set("type", storageVolume.Type.Name)
	d.Set("type_id", storageVolume.TypeId)
	d.Set("uuid", storageVolume.Uuid)

	return diags
}

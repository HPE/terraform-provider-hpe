// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package storage

import (
	"context"
	"fmt"
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
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Description:   "The name of the storage volume type",
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"id"},
			},
			"code": {
				Type:        schema.TypeString,
				Description: "The code of the storage volume type. When set alongside name, the lookup is filtered to this code.",
				Optional:    true,
				Computed:    true,
			},
			"category": {
				Type:        schema.TypeString,
				Description: "The category of the storage volume type. When set alongside name, the lookup is filtered to this category.",
				Optional:    true,
				Computed:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the storage volume type",
				Computed:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the storage volume type is enabled",
				Computed:    true,
			},
			"default_type": {
				Type:        schema.TypeBool,
				Description: "Whether this is a default storage volume type",
				Computed:    true,
			},
			"has_datastore": {
				Type:        schema.TypeBool,
				Description: "Whether the storage volume type is backed by a datastore",
				Computed:    true,
			},
			"configurable_iops": {
				Type:        schema.TypeBool,
				Description: "Whether IOPS are configurable for the storage volume type",
				Computed:    true,
			},
			"custom_size": {
				Type:        schema.TypeBool,
				Description: "Whether the storage volume type supports a custom size",
				Computed:    true,
			},
			"custom_label": {
				Type:        schema.TypeBool,
				Description: "Whether the storage volume type supports a custom label",
				Computed:    true,
			},
			"display_order": {
				Type:        schema.TypeInt,
				Description: "The display order of the storage volume type",
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

	var category string
	if v, ok := d.Get("category").(string); ok {
		category = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("category", d.Get("category")))
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
		resp, err = getStorageVolumeTypeByName(client, name, code, category)
	default:
		return diag.Errorf("Storage volume type cannot be read without name or id")
	}

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %v", resp, err)

			// MORPH-16205: a non-existent id/name must surface an error rather
			// than silently returning no state.
			return diag.Errorf("storage volume type not found (id=%d, name=%q)", id, name)
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
	d.Set("description", storageVolumeType.Description)
	d.Set("enabled", storageVolumeType.Enabled)
	d.Set("default_type", storageVolumeType.DefaultType)
	d.Set("has_datastore", storageVolumeType.HasDatastore)
	d.Set("configurable_iops", storageVolumeType.ConfigurableIOPS)
	d.Set("custom_size", storageVolumeType.CustomSize)
	d.Set("custom_label", storageVolumeType.CustomLabel)
	d.Set("display_order", storageVolumeType.DisplayOrder)

	return diags
}

func getStorageVolumeTypeByName(
	client *morpheus.Client,
	name string,
	code string,
	category string,
) (*morpheus.Response, error) {
	// Find by name, then get by ID. An optional code and/or category can be
	// supplied to narrow the search so that names that are only unique within a
	// code or category resolve to a single record.
	queryParams := map[string]string{
		"name": name,
	}
	if code != "" {
		queryParams["code"] = code
	}
	if category != "" {
		queryParams["category"] = category
	}
	resp, err := client.ListStorageVolumeTypes(&morpheus.Request{
		QueryParams: queryParams,
	})
	if err != nil {
		return resp, err
	}
	listResult, ok := resp.Result.(*morpheus.ListStorageVolumeTypesResult)
	if !ok {
		return resp, helpers.TypeAssertFailError("Result", resp.Result)
	}
	if listResult.StorageVolumeTypes == nil {
		return resp, fmt.Errorf("found 0 storage volume types named %v", name)
	}
	storageVolumeTypeCount := len(*listResult.StorageVolumeTypes)
	if storageVolumeTypeCount != 1 {
		return resp, fmt.Errorf(
			"found %d storage volume types named %v",
			storageVolumeTypeCount,
			name,
		)
	}
	firstRecord := (*listResult.StorageVolumeTypes)[0]

	return client.GetStorageVolumeType(firstRecord.ID, &morpheus.Request{})
}

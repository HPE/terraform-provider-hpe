// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package cloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/convert"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func DataSourceCloudFolder() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus Cloud Folder data source.",
		ReadContext: dataSourceCloudFolderRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Description:   "The name of the Morpheus Cloud Folder, supply either this or the id.",
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"id"},
			},
			"cloud_id": {
				Description: "The ID of the Morpheus Cloud.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"id": {
				Description:   "The ID of the Morpheus Cloud Folder, supply either this or the name.",
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"name"},
			},
			"external_id": {
				Description: "The external ID of the Morpheus Cloud Folder.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"type": {
				Description: "The type of the Morpheus Cloud Folder.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceCloudFolderRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

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

	var cloudID int
	if v, ok := d.Get("cloud_id").(int); ok {
		cloudID = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("cloud_id", d.Get("cloud_id")))
	}

	var resp *morpheus.Response
	var err error
	if id == 0 && name != "" {
		resp, err = getCloudFolderFromName(client, cloudID, name)
	} else if id != 0 && name == "" {
		resp, err = client.GetCloudResourceFolder(int64(cloudID), int64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Virtual image cannot be read without name or id")
	}
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %v", resp, err)

			return nil
		}

		return diag.FromErr(err)
	}

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var cloudFolder *morpheus.Folder
	if v, ok := resp.Result.(*morpheus.Folder); ok {
		cloudFolder = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	d.SetId(convert.Int64ToString(cloudFolder.ID))
	d.Set("external_id", cloudFolder.ExternalId)
	d.Set("type", cloudFolder.Type)
	d.Set("name", cloudFolder.Name)

	return nil
}

func getCloudFolderFromName(
	client *morpheus.Client,
	cloudID int,
	name string,
) (*morpheus.Response, error) {
	resp, err := client.ListCloudResourceFolders(int64(cloudID), &morpheus.Request{})
	if err != nil {
		return nil, err
	}

	if resp.Result == nil {
		return nil, helpers.NotFoundInResponseError("Result")
	}

	var result *morpheus.ListCloudResourceFoldersResult
	if v, ok := resp.Result.(*morpheus.ListCloudResourceFoldersResult); ok {
		result = v
	} else {
		return nil, helpers.TypeAssertFailError("Result", resp.Result)
	}

	if result.Folders == nil {
		return nil, helpers.NotFoundInResponseError("Folders")
	}

	for _, folder := range *result.Folders {
		if folder.Name == name {
			ret := &morpheus.Response{Result: &folder}

			return ret, nil
		}
	}

	return nil, fmt.Errorf("Cloud Folder not found")
}

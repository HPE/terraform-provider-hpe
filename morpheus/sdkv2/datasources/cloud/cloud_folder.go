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

	morpheus "github.com/HPE/terraform-provider-hpe/internal/sdk/legacy"
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

	// Both lookup paths converge on a single *morpheus.Folder so the response
	// mapping below is independent of which endpoint was called. The by-id
	// path returns a *GetCloudResourceFolderResult wrapper, while the by-name
	// path already resolves to a *Folder (MORPH-17898).
	var cloudFolder *morpheus.Folder
	var err error
	if id == 0 && name != "" {
		cloudFolder, err = getCloudFolderFromName(client, cloudID, name)
		if err != nil {
			return diag.FromErr(err)
		}
	} else if id != 0 && name == "" {
		var resp *morpheus.Response
		resp, err = client.GetCloudResourceFolder(int64(cloudID), int64(id), &morpheus.Request{})
		if err != nil {
			if resp != nil && resp.StatusCode == 404 {
				log.Printf("API 404: %s - %v", resp, err)

				// MORPH-17898: a non-existent id must surface an error rather
				// than silently returning no state.
				return diag.Errorf(
					"cloud folder not found (cloud_id=%d, id=%d)", cloudID, id,
				)
			}

			return diag.FromErr(err)
		}

		if resp.Result == nil {
			return diag.FromErr(helpers.NotFoundInResponseError("Result"))
		}

		result, ok := resp.Result.(*morpheus.GetCloudResourceFolderResult)
		if !ok {
			return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
		}

		if result.Folder == nil {
			return diag.FromErr(helpers.NotFoundInResponseError("Folder"))
		}

		cloudFolder = result.Folder
	} else {
		return diag.Errorf("Cloud folder cannot be read without name or id")
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
) (*morpheus.Folder, error) {
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
			ret := folder

			return &ret, nil
		}
	}

	return nil, fmt.Errorf("cloud folder not found (cloud_id=%d, name=%q)", cloudID, name)
}

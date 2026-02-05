// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package catalog

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/convert"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"
)

func DataSourceCatalogItemType() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus user group data source.",
		ReadContext: dataSourceCatalogItemTypeRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeInt,
				Description:   "The ID of the catalog item type",
				Optional:      true,
				ConflictsWith: []string{"name"},
				Computed:      true,
			},
			"name": {
				Type:          schema.TypeString,
				Description:   "The name of the catalog item type",
				Optional:      true,
				ConflictsWith: []string{"id"},
			},
		},
	}
}

func dataSourceCatalogItemTypeRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
		listResp, listErr := client.ListCatalogItemTypes(&morpheus.Request{
			QueryParams: map[string]string{
				"name": name,
			},
		})
		if listErr != nil {
			return diag.FromErr(listErr)
		}

		if listResp.Result == nil {
			return diag.FromErr(helpers.NotFoundInResponseError("Result"))
		}

		var listResult *morpheus.ListCatalogItemTypesResult
		if v, ok := listResp.Result.(*morpheus.ListCatalogItemTypesResult); ok {
			listResult = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("Result", listResp.Result))
		}

		if listResult.CatalogItemTypes == nil {
			return diag.FromErr(helpers.NotFoundInResponseError("CatalogItemTypes"))
		}

		catalogItemTypeCount := len(*listResult.CatalogItemTypes)
		if catalogItemTypeCount != 1 {
			return diag.FromErr(
				fmt.Errorf("found %d catalog item types for %v", catalogItemTypeCount, name),
			)
		}

		firstRecord := (*listResult.CatalogItemTypes)[0]
		catalogItemTypeID := firstRecord.Id
		resp, err = client.GetCatalogItemType(catalogItemTypeID, &morpheus.Request{})
	} else if id != 0 {
		resp, err = client.GetCatalogItemType(int64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Catalog Item Type cannot be read without name or id")
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

	var result *morpheus.GetCatalogItemTypeResult
	if v, ok := resp.Result.(*morpheus.GetCatalogItemTypeResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.CatalogItemType == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("CatalogItemType"))
	}

	catalogItem := result.CatalogItemType
	d.SetId(convert.Int64ToString(catalogItem.Id))
	d.Set("name", catalogItem.Name)

	return diags
}

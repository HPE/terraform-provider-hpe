// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package environment

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"
)

func DataSourceEnvironments() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus environments data source.",
		ReadContext: dataSourceEnvironmentsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"sort_ascending": {
				Type:        schema.TypeBool,
				Description: "Whether to sort the IDs in ascending order",
				Default:     true,
				Optional:    true,
			},
		},
	}
}

func dataSourceEnvironmentsRead(
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

	var sortAscending bool
	if v, ok := d.Get("sort_ascending").(bool); ok {
		sortAscending = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("sort_ascending", d.Get("sort_ascending")))
	}

	var sortOrder string
	if sortAscending {
		sortOrder = "asc"
	} else {
		sortOrder = "desc"
	}

	resp, err := client.ListEnvironments(&morpheus.Request{
		QueryParams: map[string]string{
			"max":       "50",
			"sort":      "id",
			"direction": sortOrder,
		},
	})
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

	var result *morpheus.ListEnvironmentsResult
	if v, ok := resp.Result.(*morpheus.ListEnvironmentsResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.Environments == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Environments"))
	}

	environments := result.Environments
	environmentIDs := []int64{}
	for _, environment := range *environments {
		environmentIDs = append(environmentIDs, environment.ID)
	}

	d.SetId("1")
	d.Set("ids", environmentIDs)

	return diags
}

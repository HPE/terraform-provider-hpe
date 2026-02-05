// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package cluster

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/convert"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"
)

func DataSourceClusterType() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus cluster type source.",
		ReadContext: dataSourceClusterTypeRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the Morpheus cluster type.",
				Required:    true,
			},
		},
	}
}

func dataSourceClusterTypeRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	var resp *morpheus.Response
	var err error
	resp, err = client.FindClusterTypeByName(name)
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

	var result *morpheus.ListClusterTypesResult
	if v, ok := resp.Result.(*morpheus.ListClusterTypesResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.ClusterTypes == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("ClusterTypes"))
	}

	clusterTypesPayload := result.ClusterTypes
	clusterTypes := *clusterTypesPayload
	if len(clusterTypes) == 0 {
		return diag.Errorf("cluster type not found in response data.")
	}

	clusterType := clusterTypes[0]
	if result.Meta.Total > 0 {
		d.SetId(convert.Int64ToString(clusterType.ID))
		d.Set("name", clusterType.Name)
	} else {
		return diag.Errorf("cluster type not found in response data.")
	}

	return diags
}

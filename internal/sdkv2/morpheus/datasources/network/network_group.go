// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package network

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func DataSourceNetworkGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus network group data source.",
		ReadContext: dataSourceNetworkGroupRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"name"},
				Computed:      true,
			},
			"name": {
				Type:          schema.TypeString,
				Description:   "The name of the Morpheus network group",
				Optional:      true,
				ConflictsWith: []string{"id"},
			},
			"active": {
				Type:        schema.TypeBool,
				Description: "Whether the network group is active or not",
				Computed:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the network group",
				Computed:    true,
			},
			"visibility": {
				Type:        schema.TypeString,
				Description: "Whether the network group is visible in sub-tenants or not",
				Computed:    true,
			},
		},
	}
}

func dataSourceNetworkGroupRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
		resp, err = client.FindNetworkGroupByName(name)
	} else if id != 0 {
		resp, err = client.GetNetworkGroup(int64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Network group cannot be read without name or id")
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

	var result *morpheus.GetNetworkGroupResult
	if v, ok := resp.Result.(*morpheus.GetNetworkGroupResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.NetworkGroup == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("NetworkGroup"))
	}

	networkGroup := result.NetworkGroup
	d.SetId(convert.Int64ToString(networkGroup.ID))
	d.Set("name", networkGroup.Name)
	d.Set("active", networkGroup.Active)
	d.Set("description", networkGroup.Description)
	d.Set("visibility", networkGroup.Visibility)

	return diags
}

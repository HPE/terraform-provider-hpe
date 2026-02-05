// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package network

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/convert"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"
)

func DataSourceNetworkDomain() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus domain data source.",
		ReadContext: dataSourceNetworkDomainRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"name"},
				Computed:      true,
			},
			"active": {
				Type:        schema.TypeBool,
				Description: "Whether the domain is active",
				Computed:    true,
			},
			"name": {
				Type:          schema.TypeString,
				Description:   "The name of the Morpheus domain",
				Optional:      true,
				ConflictsWith: []string{"id"},
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the Morpheus domain",
				Computed:    true,
			},
			"visibility": {
				Type:        schema.TypeString,
				Description: "The visibility of the Morpheus domain",
				Computed:    true,
			},
		},
	}
}

func dataSourceNetworkDomainRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	// Warning or errors can be collected in a slice type
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

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == 0 && name != "" {
		resp, err = client.FindNetworkDomainByName(name)
	} else if id != 0 {
		resp, err = client.GetNetworkDomain(int64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Domain cannot be read without name or id")
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

	// store resource data
	var result *morpheus.GetNetworkDomainResult
	if v, ok := resp.Result.(*morpheus.GetNetworkDomainResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("GetNetworkDomainResult"))
	}

	domain := result.NetworkDomain
	if domain == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("NetworkDomain"))
	}

	d.SetId(convert.Int64ToString(domain.ID))
	d.Set("active", domain.Active)
	d.Set("name", domain.Name)
	d.Set("description", domain.Description)
	d.Set("visibility", domain.Visibility)

	return diags
}

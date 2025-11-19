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

func DataSourceNetworkSubnet() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus network subnet data source.",
		ReadContext: dataSourceNetworkSubnetRead,
		Schema: map[string]*schema.Schema{
			"network_id": {
				Type:        schema.TypeInt,
				Description: "The id of the Morpheus network to search for the subnet.",
				Required:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the Morpheus network subnet.",
				Optional:    true,
			},
			"id": {
				Type:        schema.TypeInt,
				Description: "The id of the network subnet",
				Optional:    true,
				Computed:    true,
			},
			"external_id": {
				Type:        schema.TypeString,
				Description: "The external id of the network subnet",
				Computed:    true,
			},
			"cidr": {
				Type:        schema.TypeString,
				Description: "The cidr of the network subnet",
				Computed:    true,
			},
			"netmask": {
				Type:        schema.TypeString,
				Description: "The netmask of the network subnet",
				Computed:    true,
			},
			"visibility": {
				Type:        schema.TypeString,
				Description: "The visibility of the network subnet",
				Computed:    true,
			},
		},
	}
}

func dataSourceNetworkSubnetRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}

	var networkID int
	if v, ok := d.Get("network_id").(int); ok {
		networkID = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("network_id", d.Get("network_id")))
	}

	if id == 0 && name == "" {
		return diag.Errorf(
			"Either 'id' or 'name' must be provided to search for the network subnet",
		)
	}

	var resp *morpheus.Response
	var err error
	var networkSubnet *morpheus.NetworkSubnet

	if id != 0 {
		resp, err = client.GetNetworkSubnet(int64(id), &morpheus.Request{})
		if err != nil {
			errorPrefix := "API FAILURE"
			if resp != nil && resp.StatusCode == 404 {
				errorPrefix = "API 404"
			}
			log.Printf("%s: %s - %v", errorPrefix, resp, err)

			return diag.FromErr(err)
		}

		log.Printf("API RESPONSE: %s", resp)

		if resp.Result == nil {
			return diag.FromErr(helpers.NotFoundInResponseError("Result"))
		}

		var result *morpheus.GetNetworkSubnetResult
		if v, ok := resp.Result.(*morpheus.GetNetworkSubnetResult); ok {
			result = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
		}

		if result.NetworkSubnet == nil {
			return diag.FromErr(helpers.NotFoundInResponseError("NetworkSubnet"))
		}

		networkSubnet = result.NetworkSubnet

	} else {
		resp, err = client.ListNetworkSubnetsByNetwork(int64(networkID), &morpheus.Request{
			QueryParams: map[string]string{
				"name": name,
			},
		})
		if err != nil {
			log.Printf("API FAILURE: %s - %v", resp, err)

			return diag.FromErr(err)
		}

		if resp.Result == nil {
			return diag.FromErr(helpers.NotFoundInResponseError("Result"))
		}

		var listResult *morpheus.ListNetworkSubnetsByNetworkResult
		if v, ok := resp.Result.(*morpheus.ListNetworkSubnetsByNetworkResult); ok {
			listResult = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
		}

		if listResult.NetworkSubnets == nil {
			return diag.FromErr(helpers.NotFoundInResponseError("NetworkSubnets"))
		}

		networkSubnetsCount := len(*listResult.NetworkSubnets)
		if networkSubnetsCount != 1 {
			return diag.Errorf("found %d Network Subnets for %v", networkSubnetsCount, name)
		}
		firstRecord := (*listResult.NetworkSubnets)[0]
		networkSubnetID := firstRecord.ID
		resp, err = client.GetNetworkSubnet(networkSubnetID, &morpheus.Request{})
		if err != nil {
			errorPrefix := "API FAILURE"
			if resp != nil && resp.StatusCode == 404 {
				errorPrefix = "API 404"
			}
			log.Printf("%s: %s - %v", errorPrefix, resp, err)

			return diag.FromErr(err)
		}

		log.Printf("API RESPONSE: %s", resp)

		if resp.Result == nil {
			return diag.FromErr(helpers.NotFoundInResponseError("Result"))
		}

		var result *morpheus.GetNetworkSubnetResult
		if v, ok := resp.Result.(*morpheus.GetNetworkSubnetResult); ok {
			result = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
		}

		if result.NetworkSubnet == nil {
			return diag.FromErr(helpers.NotFoundInResponseError("NetworkSubnet"))
		}

		networkSubnet = result.NetworkSubnet
	}

	if networkSubnet == nil {
		return diag.Errorf("Network subnet not found in response data.")
	}

	d.SetId(convert.Int64ToString(networkSubnet.ID))
	d.Set("name", networkSubnet.Name)
	d.Set("external_id", networkSubnet.ExternalId)
	d.Set("cidr", networkSubnet.Cidr)
	d.Set("netmask", networkSubnet.Netmask)
	d.Set("visibility", networkSubnet.Visibility)

	return diags
}

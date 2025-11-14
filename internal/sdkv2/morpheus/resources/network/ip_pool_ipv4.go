// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package network

import (
	"context"
	"log"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"
)

func ResourceIPPoolIPv4() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus IPv4 ip pool resource",
		CreateContext: resourceIPPoolIPv4Create,
		ReadContext:   resourceIPPoolIPv4Read,
		UpdateContext: resourceIPPoolIPv4Update,
		DeleteContext: resourceIPPoolIPv4Delete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the IPv4 IP address pool",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the IPv4 IP address pool",
				Required:    true,
			},
			"ip_range": {
				Type:        schema.TypeList,
				Description: "The IPv4 IP address pool IP ranges",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"starting_address": {
							Type:        schema.TypeString,
							Description: "The starting address of the IPv4 IP address pool IP range",
							Required:    true,
						},
						"ending_address": {
							Type:        schema.TypeString,
							Description: "The ending address of the IPv4 IP address pool IP range",
							Required:    true,
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceIPPoolIPv4Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	var ipRange []any
	if v, ok := d.Get("ip_range").([]any); ok {
		ipRange = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("ip_range", d.Get("ip_range")))
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"networkPool": map[string]any{
				"name":     name,
				"type":     "morpheus",
				"ipRanges": parseIPPoolRanges(ipRange),
			},
		},
	}
	resp, err := client.CreateNetworkPool(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.CreateNetworkPoolResult
	if v, ok := resp.Result.(*morpheus.CreateNetworkPoolResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.NetworkPool == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("NetworkPool"))
	}

	pool := result.NetworkPool
	d.SetId(convert.Int64ToString(pool.ID))
	diags = append(diags, resourceIPPoolIPv4Read(ctx, d, meta)...)

	return diags
}

func resourceIPPoolIPv4Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	id := d.Id()

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}

	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindNetworkPoolByName(name)
	} else if id != "" {
		resp, err = client.GetNetworkPool(convert.StringToInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Pool cannot be read without name or id")
	}

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)
			log.Printf("Forcing recreation of resource")
			d.SetId("")

			return diags
		}

		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.GetNetworkPoolResult
	if v, ok := resp.Result.(*morpheus.GetNetworkPoolResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.NetworkPool == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("NetworkPool"))
	}

	pool := result.NetworkPool
	d.SetId(convert.Int64ToString(pool.ID))
	d.Set("name", pool.Name)

	var ipRanges []map[string]any
	var unsortedRanges []IPRange
	if pool.IpRanges != nil {
		for _, iprange := range pool.IpRanges {
			var IPR IPRange
			IPR.ID = iprange.ID
			IPR.EndAddress = iprange.EndAddress
			IPR.StartAddress = iprange.StartAddress
			unsortedRanges = append(unsortedRanges, IPR)
		}
	}
	sort.Slice(unsortedRanges, func(i, j int) bool { return unsortedRanges[i].ID < unsortedRanges[j].ID })

	if unsortedRanges != nil {
		for i := 0; i < len(unsortedRanges); i++ {
			ipRange := unsortedRanges[i]
			rangePayload := make(map[string]any)
			rangePayload["ending_address"] = ipRange.EndAddress
			rangePayload["starting_address"] = ipRange.StartAddress
			ipRanges = append(ipRanges, rangePayload)
		}
	}
	d.Set("ip_range", ipRanges)

	return diags
}

func resourceIPPoolIPv4Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	id := d.Id()

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}

	var ipRange []any
	if v, ok := d.Get("ip_range").([]any); ok {
		ipRange = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("ip_range", d.Get("ip_range")))
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"networkPool": map[string]any{
				"name":     name,
				"type":     "morpheus",
				"ipRanges": parseIPPoolRanges(ipRange),
			},
		},
	}
	resp, err := client.UpdateNetworkPool(convert.StringToInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.UpdateNetworkPoolResult
	if v, ok := resp.Result.(*morpheus.UpdateNetworkPoolResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.NetworkPool == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("NetworkPool"))
	}

	pool := result.NetworkPool
	d.SetId(convert.Int64ToString(pool.ID))

	return resourceIPPoolIPv4Read(ctx, d, meta)
}

func resourceIPPoolIPv4Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteNetworkPool(convert.StringToInt64(id), req)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)

			return nil
		}

		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	d.SetId("")

	return diags
}

func parseIPPoolRanges(variables []any) []map[string]any {
	var poolRanges []map[string]any
	if variables == nil {
		return poolRanges
	}

	for i := 0; i < len(variables); i++ {
		row := make(map[string]any)
		ippoolconfig := variables[i].(map[string]any)
		for k, v := range ippoolconfig {
			switch k {
			case "starting_address":
				row["startAddress"] = v.(string)
			case "ending_address":
				row["endAddress"] = v.(string)
			}
		}
		poolRanges = append(poolRanges, row)
	}

	return poolRanges
}

type IPRange struct {
	ID           int64  `json:"id"`
	StartAddress string `json:"startAddress"`
	EndAddress   string `json:"endAddress"`
}

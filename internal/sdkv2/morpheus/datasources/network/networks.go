// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package network

import (
	"context"
	"log"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"
)

func DataSourceNetworks() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus networks data source.",
		ReadContext: dataSourceNetworksRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"cloud_id": {
				Type:        schema.TypeInt,
				Description: "The id of the Morpheus cloud to search for the network.",
				Optional:    true,
			},
			"sort_ascending": {
				Type:        schema.TypeBool,
				Description: "Whether to sort the IDs in ascending order. Defaults to true",
				Default:     true,
				Optional:    true,
			},
			"filter": {
				Type:        schema.TypeSet,
				Description: "Custom filter block as described below.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "The name of the filter. Valid names are (name)",
							Required:    true,
							ValidateFunc: validation.StringInSlice(
								[]string{"name"},
								false,
							),
						},
						"values": {
							Type: schema.TypeSet,
							Description: "The filter values. Filters support Golang regex " +
								"and can be tested at https://regex101.com/",
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceNetworksRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	var resp *morpheus.Response
	var err error
	var sortOrder string
	var names []string

	var filterSet *schema.Set
	if v, ok := d.Get("filter").(*schema.Set); ok {
		filterSet = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("filter", d.Get("filter")))
	}

	filterList := filterSet.List()
	if len(filterList) > 0 {
		for _, filter := range filterList {
			var filterPayload map[string]any
			if v, ok := filter.(map[string]any); ok {
				filterPayload = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("filter", filter))
			}

			var filterName string
			if v, ok := filterPayload["name"].(string); ok {
				filterName = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("filter name", filterPayload["name"]))
			}

			if filterName == "name" {
				var valuesSet *schema.Set
				if v, ok := filterPayload["values"].(*schema.Set); ok {
					valuesSet = v
				} else {
					return diag.FromErr(helpers.TypeAssertFailError("values", filterPayload["values"]))
				}

				valuesList := valuesSet.List()
				for _, item := range valuesList {
					var itemStr string
					if v, ok := item.(string); ok {
						itemStr = v
					} else {
						return diag.FromErr(helpers.TypeAssertFailError("value", item))
					}
					names = append(names, itemStr)
				}
			}
		}
	}

	if len(names) == 0 {
		names = append(names, "$")
	}

	// Sort environments in ascending or descending order
	var sortAscending bool
	if v, ok := d.Get("sort_ascending").(bool); ok {
		sortAscending = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("sort_ascending", d.Get("sort_ascending")))
	}

	if sortAscending {
		sortOrder = "asc"
	} else {
		sortOrder = "desc"
	}

	params := make(map[string]string)
	params["max"] = "250"
	params["sort"] = "id"
	params["direction"] = sortOrder

	var cloudID int
	if v, ok := d.Get("cloud_id").(int); ok {
		cloudID = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("cloud_id", d.Get("cloud_id")))
	}

	if cloudID > 0 {
		cloudIDString := convert.IntToString(cloudID)
		params["zoneId"] = cloudIDString
	}

	resp, err = client.ListNetworks(&morpheus.Request{
		QueryParams: params,
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

	var networksIDs []string

	// store resource data
	var result *morpheus.ListNetworksResult
	if v, ok := resp.Result.(*morpheus.ListNetworksResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("ListNetworksResult", resp.Result))
	}

	if result.Networks == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Networks"))
	}

	networks := result.Networks
	if networks != nil {
		for _, network := range *networks {
			if len(names) > 0 {
				if regexCheck(names, network.Name) {
					networksIDs = append(networksIDs, convert.Int64ToString(network.ID))
				}
			} else {
				networksIDs = append(networksIDs, convert.Int64ToString(network.ID))
			}
		}
	}

	d.SetId("1")
	d.Set("ids", networksIDs)

	return diags
}

func regexCheck(s []string, str string) bool {
	var status int
	for _, v := range s {
		match, _ := regexp.MatchString(v, str)
		if match {
			status = status + 1
		}
	}

	return status > 0
}

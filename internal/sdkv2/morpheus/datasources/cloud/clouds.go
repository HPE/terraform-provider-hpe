// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package cloud

import (
	"context"
	"log"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func DataSourceClouds() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus clouds data source.",
		ReadContext: dataSourceCloudsRead,
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
			"filter": {
				Type:        schema.TypeSet,
				Description: "Custom filter block as described below.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "The name of the filter. Filter names are case-sensitive. Valid names are (name)",
							Required:    true,
							ValidateFunc: validation.StringInSlice(
								[]string{"name"},
								false,
							),
						},
						"values": {
							Type: schema.TypeSet,
							Description: "The filter values. Filter values are case-sensitive. " +
								"Filters values support the use of Golang regex and can be tested at " +
								"https://regex101.com/",
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

func dataSourceCloudsRead(
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

	var resp *morpheus.Response
	var err error
	var sortOrder string
	var names []string

	var filter any
	if v, ok := d.Get("filter").(*schema.Set); ok {
		filter = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("filter", d.Get("filter")))
	}

	filterSet := filter.(*schema.Set)
	if filterSet != nil && len(filterSet.List()) > 0 {
		filters := filterSet.List()
		for _, filterItem := range filters {
			var filterPayload map[string]any
			if v, ok := filterItem.(map[string]any); ok {
				filterPayload = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("filter", filterItem))
			}

			var filterName string
			if v, ok := filterPayload["name"].(string); ok {
				filterName = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("name", filterPayload["name"]))
			}

			if filterName == "name" {
				var values any
				if v, ok := filterPayload["values"].(*schema.Set); ok {
					values = v
				} else {
					return diag.FromErr(helpers.TypeAssertFailError("values", filterPayload["values"]))
				}

				valuesSet := values.(*schema.Set)
				if valuesSet != nil {
					for _, item := range valuesSet.List() {
						var itemStr string
						if v, ok := item.(string); ok {
							itemStr = v
						} else {
							return diag.FromErr(helpers.TypeAssertFailError("thing", item))
						}
						names = append(names, itemStr)
					}
				}
			}
		}
	}

	if len(names) == 0 {
		names = append(names, "$")
	}

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

	resp, err = client.ListClouds(&morpheus.Request{
		QueryParams: map[string]string{
			"max":       "250",
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

	cloudIDs := []int64{}

	var result *morpheus.ListCloudsResult
	if v, ok := resp.Result.(*morpheus.ListCloudsResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("result", resp.Result))
	}

	if result.Clouds == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Clouds"))
	}

	clouds := result.Clouds
	if clouds != nil {
		for _, cloud := range *clouds {
			if len(names) > 0 {
				if regexCheck(names, cloud.Name) {
					cloudIDs = append(cloudIDs, cloud.ID)
				}
			} else {
				cloudIDs = append(cloudIDs, cloud.ID)
			}
		}
	}
	d.SetId("1")
	d.Set("ids", cloudIDs)

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

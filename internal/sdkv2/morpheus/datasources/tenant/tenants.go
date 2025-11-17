// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package tenant

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

func DataSourceTenants() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus tenants data source.",
		ReadContext: dataSourceTenantsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
								"Filters values support the use of Golang regex and can be " +
								"tested at https://regex101.com/",
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

func dataSourceTenantsRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
			var filterDetails map[string]any
			if v, ok := filter.(map[string]any); ok {
				filterDetails = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("filter", filter))
			}

			var filterName string
			if v, ok := filterDetails["name"].(string); ok {
				filterName = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("name", filterDetails["name"]))
			}

			if filterName == "name" {
				var valuesSet *schema.Set
				if v, ok := filterDetails["values"].(*schema.Set); ok {
					valuesSet = v
				} else {
					return diag.FromErr(helpers.TypeAssertFailError("values", filterDetails["values"]))
				}

				valuesList := valuesSet.List()
				for _, item := range valuesList {
					var itemStr string
					if v, ok := item.(string); ok {
						itemStr = v
					} else {
						return diag.FromErr(helpers.TypeAssertFailError("values element", item))
					}
					names = append(names, itemStr)
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

	// Sort environments in ascending or descending order
	if sortAscending {
		sortOrder = "asc"
	} else {
		sortOrder = "desc"
	}

	resp, err = client.ListTenants(&morpheus.Request{
		QueryParams: map[string]string{
			"max":       "100",
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

	var tenantIDs []string

	// store resource data
	var result *morpheus.ListTenantsResult
	if v, ok := resp.Result.(*morpheus.ListTenantsResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.Accounts == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Accounts"))
	}

	tenants := result.Accounts
	for _, tenant := range *tenants {
		if regexCheck(names, tenant.Name) {
			tenantIDs = append(tenantIDs, convert.Int64ToString(tenant.ID))
		}
	}
	d.SetId("1")
	d.Set("ids", tenantIDs)

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

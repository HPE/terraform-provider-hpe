// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package policy

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

func DataSourcePolicies() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus policies data source.",
		ReadContext: dataSourcePoliciesRead,
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
							Type: schema.TypeString,
							Description: "The name of the filter. Filter names are case-sensitive. " +
								"Valid names are (name, type)",
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"name", "type"}, false),
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

func dataSourcePoliciesRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
	var policyTypes []string
	var names []string

	var filterSet *schema.Set
	if v, ok := d.Get("filter").(*schema.Set); ok {
		filterSet = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("filter", d.Get("filter")))
	}

	if filterSet != nil && len(filterSet.List()) > 0 {
		filters := filterSet.List()
		for _, filter := range filters {
			var test map[string]any
			if v, ok := filter.(map[string]any); ok {
				test = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("filters element", filter))
			}

			var nameStr string
			if v, ok := test["name"].(string); ok {
				nameStr = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("name", test["name"]))
			}

			if nameStr == "type" {
				var valuesSet *schema.Set
				if v, ok := test["values"].(*schema.Set); ok {
					valuesSet = v
				} else {
					return diag.FromErr(helpers.TypeAssertFailError("values", test["values"]))
				}

				if valuesSet != nil {
					valuesList := valuesSet.List()
					for _, item := range valuesList {
						var itemStr string
						if v, ok := item.(string); ok {
							itemStr = v
						} else {
							return diag.FromErr(helpers.TypeAssertFailError("values element", item))
						}
						policyTypes = append(policyTypes, itemStr)
					}
				}
			}

			if nameStr == "name" {
				var valuesSet *schema.Set
				if v, ok := test["values"].(*schema.Set); ok {
					valuesSet = v
				} else {
					return diag.FromErr(helpers.TypeAssertFailError("values", test["values"]))
				}

				if valuesSet != nil {
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
	}

	if len(policyTypes) == 0 {
		policyTypes = append(policyTypes, "$")
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

	resp, err = client.ListPolicies(&morpheus.Request{
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

	var policyIDs []string

	// store resource data
	var result *morpheus.ListPoliciesResult
	if v, ok := resp.Result.(*morpheus.ListPoliciesResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("result", resp.Result))
	}

	if result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("ListPoliciesResult"))
	}

	if result.Policies == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Policies"))
	}

	policies := result.Policies
	if policies != nil {
		for _, policy := range *policies {
			if regexCheck(policyTypes, policy.PolicyType.Name) && regexCheck(names, policy.Name) {
				policyIDs = append(policyIDs, convert.Int64ToString(policy.ID))
			}
		}
	}

	d.SetId("1")
	d.Set("ids", policyIDs)

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

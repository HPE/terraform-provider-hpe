// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package image

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

func DataSourceImageVirtualImages() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus virtual images data source.",
		ReadContext: dataSourceImageVirtualImagesRead,
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
			"source": {
				Type:         schema.TypeString,
				Description:  "The source of the Morpheus virtual image (User, System, Synced) (Default: User)",
				Optional:     true,
				Default:      "User",
				ValidateFunc: validation.StringInSlice([]string{"User", "System", "Synced"}, false),
			},
			"filter": {
				Type:        schema.TypeSet,
				Description: "Custom filter block as described below.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "The name of the filter. Filter names are case-sensitive. Valid names are (name, type)",
							Required:    true,
							ValidateFunc: validation.StringInSlice(
								[]string{"name", "type"},
								false,
							),
						},
						"values": {
							Type: schema.TypeSet,
							Description: "The filter values. Filter values are case-sensitive. " +
								"Filters values support the use of Golang regex and can be tested " +
								"at https://regex101.com/",
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

func dataSourceImageVirtualImagesRead(
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

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	var sortOrder string
	var imageTypes []string
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
			var filterData map[string]any
			if v, ok := filter.(map[string]any); ok {
				filterData = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("filter", filter))
			}

			var filterName string
			if v, ok := filterData["name"].(string); ok {
				filterName = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("name", filterData["name"]))
			}

			if filterName == "type" {
				var valuesSet *schema.Set
				if v, ok := filterData["values"].(*schema.Set); ok {
					valuesSet = v
				} else {
					return diag.FromErr(helpers.TypeAssertFailError("values", filterData["values"]))
				}

				if valuesSet != nil {
					for _, item := range valuesSet.List() {
						var itemStr string
						if v, ok := item.(string); ok {
							itemStr = v
						} else {
							return diag.FromErr(helpers.TypeAssertFailError("values", item))
						}
						imageTypes = append(imageTypes, itemStr)
					}
				}
			}

			if filterName == "name" {
				var valuesSet *schema.Set
				if v, ok := filterData["values"].(*schema.Set); ok {
					valuesSet = v
				} else {
					return diag.FromErr(helpers.TypeAssertFailError("values", filterData["values"]))
				}

				if valuesSet != nil {
					for _, item := range valuesSet.List() {
						var itemStr string
						if v, ok := item.(string); ok {
							itemStr = v
						} else {
							return diag.FromErr(helpers.TypeAssertFailError("values", item))
						}
						names = append(names, itemStr)
					}
				}
			}
		}
	}

	if len(imageTypes) == 0 {
		imageTypes = append(imageTypes, "$")
	}

	if len(names) == 0 {
		names = append(names, "$")
	}

	// Sort virtual images in ascending or descending order
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

	var source string
	if v, ok := d.Get("source").(string); ok {
		source = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("source", d.Get("source")))
	}

	output, err := listAllVirtualImages(client, 200, sortOrder, source)
	if err != nil {
		return diag.FromErr(err)
	}

	var virtualImageIDs []string

	// store resource data
	for _, virtualImage := range output {
		if regexCheck(imageTypes, virtualImage.ImageType) &&
			regexCheck(names, virtualImage.Name) {
			virtualImageIDs = append(virtualImageIDs, convert.Int64ToString(virtualImage.ID))
		}
	}

	d.SetId("1")
	d.Set("ids", virtualImageIDs)

	return diags
}

func listAllVirtualImages(
	client *morpheus.Client,
	max int,
	sortOrder string,
	source string,
) ([]morpheus.VirtualImage, error) {
	var images []morpheus.VirtualImage

	// Fetch initial images
	params := make(map[string]string)
	params["max"] = convert.IntToString(max)
	params["sort"] = "id"
	params["direction"] = sortOrder
	params["filterType"] = source
	resp, err := client.ListVirtualImages(&morpheus.Request{
		QueryParams: params,
	})
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %v", resp, err)
		} else {
			log.Printf("API FAILURE: %s - %v", resp, err)
		}

		return nil, err
	}
	log.Printf("API RESPONSE: %s", resp)

	var result *morpheus.ListVirtualImagesResult
	if v, ok := resp.Result.(*morpheus.ListVirtualImagesResult); ok {
		result = v
	} else {
		return nil, helpers.TypeAssertFailError("Result", resp.Result)
	}

	if result == nil {
		return nil, helpers.NotFoundInResponseError("ListVirtualImagesResult")
	}

	if result.Meta == nil {
		return nil, helpers.NotFoundInResponseError("Meta")
	}

	pollIterations := result.Meta.Total / int64(max)

	// Add Page 1 virtual images
	if result.VirtualImages != nil {
		images = append(images, *result.VirtualImages...)
	}

	currentPage := 1
	for currentPage < int(pollIterations) {
		// Fetch initial images
		params := make(map[string]string)
		params["max"] = convert.IntToString(max)
		params["sort"] = "id"
		params["direction"] = sortOrder
		params["filterType"] = source

		params["offset"] = convert.IntToString(currentPage * max)
		resp, err := client.ListVirtualImages(&morpheus.Request{
			QueryParams: params,
		})
		if err != nil {
			if resp != nil && resp.StatusCode == 404 {
				log.Printf("API 404: %s - %v", resp, err)
			} else {
				log.Printf("API FAILURE: %s - %v", resp, err)
			}

			return nil, err
		}
		log.Printf("API RESPONSE: %s", resp)

		var result *morpheus.ListVirtualImagesResult
		if v, ok := resp.Result.(*morpheus.ListVirtualImagesResult); ok {
			result = v
		} else {
			return nil, helpers.TypeAssertFailError("Result", resp.Result)
		}

		if result == nil {
			return nil, helpers.NotFoundInResponseError("ListVirtualImagesResult")
		}

		if result.VirtualImages != nil {
			images = append(images, *result.VirtualImages...)
		}

		currentPage++
	}

	return images, nil
}

func regexCheck(patterns []string, str string) bool {
	var status int
	for _, v := range patterns {
		match, _ := regexp.MatchString(v, str)
		if match {
			status = status + 1
		}
	}

	return status > 0
}

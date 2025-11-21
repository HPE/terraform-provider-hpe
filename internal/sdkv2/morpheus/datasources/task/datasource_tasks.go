package task

import (
	"context"
	"log"
	"regexp"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"
)

func DataSourceTasks() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus tasks data source.",
		ReadContext: dataSourceTasksRead,
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
							Type:         schema.TypeString,
							Description:  "The name of the filter. Filter names are case-sensitive. Valid names are (name, type)",
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"name", "type"}, false),
						},
						"values": {
							Type: schema.TypeSet,
							//nolint:lll
							Description: "The filter values. Filter values are case-sensitive. Filters values support the use of Golang regex and can be tested at https://regex101.com/",
							Required:    true,
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

func dataSourceTasksRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
	var taskTypes []string
	var names []string

	var filterSet *schema.Set
	if v, ok := d.Get("filter").(*schema.Set); ok {
		filterSet = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("filter", d.Get("filter")))
	}

	filters := filterSet.List()
	//nolint:staticcheck
	if filters != nil && len(filters) > 0 {
		for _, filter := range filters {
			var test map[string]any
			if v, ok := filter.(map[string]any); ok {
				test = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("filter", filter))
			}

			var filterName string
			if v, ok := test["name"].(string); ok {
				filterName = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("name", test["name"]))
			}

			if filterName == "type" {
				var valuesSet *schema.Set
				if v, ok := test["values"].(*schema.Set); ok {
					valuesSet = v
				} else {
					return diag.FromErr(helpers.TypeAssertFailError("values", test["values"]))
				}
				values := valuesSet.List()
				if values == nil {
					continue
				}
				for _, item := range values {
					var itemStr string
					if v, ok := item.(string); ok {
						itemStr = v
					} else {
						return diag.FromErr(helpers.TypeAssertFailError("item", item))
					}
					taskTypes = append(taskTypes, itemStr)
				}
			}

			if filterName == "name" {
				var valuesSet *schema.Set
				if v, ok := test["values"].(*schema.Set); ok {
					valuesSet = v
				} else {
					return diag.FromErr(helpers.TypeAssertFailError("values", test["values"]))
				}
				values := valuesSet.List()
				if values == nil {
					continue
				}
				for _, item := range values {
					var itemStr string
					if v, ok := item.(string); ok {
						itemStr = v
					} else {
						return diag.FromErr(helpers.TypeAssertFailError("item", item))
					}
					names = append(names, itemStr)
				}
			}
		}
	}

	if len(taskTypes) == 0 {
		taskTypes = append(taskTypes, "$")
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

	resp, err = client.ListTasks(&morpheus.Request{
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
		} else {
			log.Printf("API FAILURE: %s - %v", resp, err)

			return diag.FromErr(err)
		}
	}
	log.Printf("API RESPONSE: %s", resp)

	var taskIDs []string

	// store resource data
	var result *morpheus.ListTasksResult
	if v, ok := resp.Result.(*morpheus.ListTasksResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("result", resp.Result))
	}
	if result.Tasks == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Tasks"))
	}
	tasks := result.Tasks
	for _, task := range *tasks {
		if regexCheck(taskTypes, task.TaskType.Name) && regexCheck(names, task.Name) {
			taskIDs = append(taskIDs, convert.Int64ToString(task.ID))
		}
	}
	d.SetId("1")
	d.Set("ids", taskIDs)

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

	if status > 0 {
		return true
	} else {
		return false
	}
}

// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package integration

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"
)

func DataSourceVroWorkflowRead() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus vRO workflow data source.",
		ReadContext: dataSourceVroWorkflowRead,
		Schema: map[string]*schema.Schema{
			"value": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Description:   "The name of the vRealize orchestrator workflow",
				Optional:      true,
				ConflictsWith: []string{"value"},
			},
		},
	}
}

func dataSourceVroWorkflowRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	var value int
	if v, ok := d.Get("value").(int); ok {
		value = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("value", d.Get("value")))
	}

	// lookup by name if we do not have an value yet
	var resp *morpheus.Response
	var err error

	resp, err = client.GetOptionSource("vroWorkflow", &morpheus.Request{})
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

	var result *morpheus.GetOptionSourceResult
	if v, ok := resp.Result.(*morpheus.GetOptionSourceResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.Data == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Data"))
	}

	allWorkflows := *result.Data
	if allWorkflows == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Data"))
	}

	var workflow morpheus.OptionSourceOption
	for i := range allWorkflows {
		if value == 0 && name != "" {
			if strings.EqualFold(allWorkflows[i].Name, name) {
				workflow = allWorkflows[i]

				break
			}
		} else if value != 0 {
			if value == allWorkflows[i].Value {
				workflow = allWorkflows[i]

				break
			}
		} else {
			return diag.Errorf("vRO workflow cannot be read without name or value")
		}
	}

	// store resource data
	var workflowValue float64
	if v, ok := workflow.Value.(float64); ok {
		workflowValue = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Value", workflow.Value))
	}

	d.SetId(fmt.Sprintf("%g", workflowValue))
	d.Set("value", workflow.Value)
	d.Set("name", workflow.Name)

	return diags
}

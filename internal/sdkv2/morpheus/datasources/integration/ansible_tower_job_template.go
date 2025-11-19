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

func DataSourceAnsibleTowerJobTemplate() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus ansible tower job template data source.",
		ReadContext: dataSourceAnsibleTowerJobTemplateRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Description:   "The name of the ansible tower job template",
				Optional:      true,
				ConflictsWith: []string{"id"},
			},
		},
	}
}

func dataSourceAnsibleTowerJobTemplateRead(
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

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}

	var value int
	if v, ok := d.Get("id").(int); ok {
		value = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("id", d.Get("id")))
	}

	var resp *morpheus.Response
	var err error

	resp, err = client.GetOptionSource("ansibleTowerJobTemplate", &morpheus.Request{})
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %v", resp, err)

			return nil
		}

		log.Printf("API FAILURE: %s - %v", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.GetOptionSourceResult
	if v, ok := resp.Result.(*morpheus.GetOptionSourceResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.Data == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Data"))
	}

	allTemplates := *result.Data

	var template morpheus.OptionSourceOption
	for i := range allTemplates {
		if value == 0 && name != "" {
			if strings.EqualFold(allTemplates[i].Name, name) {
				template = allTemplates[i]

				break
			}
		} else if value != 0 {
			if value == allTemplates[i].Value {
				template = allTemplates[i]

				break
			}
		} else {
			return diag.Errorf(
				"Ansible tower job template cannot be read without name or value",
			)
		}
	}

	var templateValue float64
	if v, ok := template.Value.(float64); ok {
		templateValue = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("template.Value", template.Value))
	}

	d.SetId(fmt.Sprintf("%g", templateValue))
	d.Set("id", template.Value)
	d.Set("name", template.Name)

	return diags
}

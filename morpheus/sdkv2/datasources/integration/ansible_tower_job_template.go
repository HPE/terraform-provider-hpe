// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package integration

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	morpheus "github.com/HPE/terraform-provider-hpe/internal/sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"
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
			"ansible_tower_integration_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the ansible tower integration",
				Required:    true,
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

	var ansibleTowerIntegrationId int
	if v, ok := d.Get("ansible_tower_integration_id").(int); ok {
		ansibleTowerIntegrationId = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("ansible_tower_integration_id",
			d.Get("ansible_tower_integration_id")))
	}

	resp, err = client.GetOptionSource("ansibleTowerJobTemplate", &morpheus.Request{
		QueryParams: map[string]string{
			"ansibleTowerIntegrationId": strconv.Itoa(ansibleTowerIntegrationId),
		},
	})
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %v", resp, err)

			// MORPH-17962: a 404 (e.g. a missing integration) must surface an
			// error rather than silently returning no state, consistent with
			// the other data sources.
			return diag.Errorf(
				"ansible tower job template not found "+
					"(ansible_tower_integration_id=%d)",
				ansibleTowerIntegrationId,
			)
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

	// Validate lookup inputs before iterating so empty input is rejected
	// cleanly even when the option source returns no templates (MORPH-17962).
	if value == 0 && name == "" {
		return diag.Errorf(
			"Ansible tower job template cannot be read without name or value",
		)
	}

	allTemplates := *result.Data

	var template morpheus.OptionSourceOption
	var found bool
	for i := range allTemplates {
		if value == 0 {
			if strings.EqualFold(allTemplates[i].Name, name) {
				template = allTemplates[i]
				found = true

				break
			}
		} else {
			if value == allTemplates[i].Value {
				template = allTemplates[i]
				found = true

				break
			}
		}
	}

	// MORPH-17962: a non-existent name/value (or a missing integration whose
	// option source is empty) must surface a clean diagnostic rather than a
	// raw "Type assertion failed for value: <nil>" from the assertion below.
	if !found {
		return diag.Errorf(
			"ansible tower job template not found "+
				"(name=%q, id=%d, ansible_tower_integration_id=%d)",
			name, value, ansibleTowerIntegrationId,
		)
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

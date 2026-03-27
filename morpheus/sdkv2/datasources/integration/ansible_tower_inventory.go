// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package integration

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/convert"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"
)

func DataSourceAnsibleTowerInventory() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus ansible tower inventory data source.",
		ReadContext: dataSourceAnsibleTowerInventoryRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
			},
			"ansible_tower_integration_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the ansible tower integration",
				Required:    true,
			},
			"name": {
				Type:          schema.TypeString,
				Description:   "The name of the ansible tower inventory",
				Optional:      true,
				ConflictsWith: []string{"id"},
			},
		},
	}
}

func dataSourceAnsibleTowerInventoryRead(
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

	resp, err = client.GetOptionSource("ansibleTowerInventory", &morpheus.Request{
		QueryParams: map[string]string{
			"ansibleTowerIntegrationId": strconv.Itoa(ansibleTowerIntegrationId),
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

	allInventories := *result.Data
	if allInventories == nil {
		return diag.FromErr(helpers.NilPointerError("allInventories"))
	}

	var inventory morpheus.OptionSourceOption
	for i := range allInventories {
		if value == 0 && name != "" {
			if strings.EqualFold(allInventories[i].Name, name) {
				inventory = allInventories[i]

				break
			}
		} else if value != 0 {
			if value == allInventories[i].Value {
				inventory = allInventories[i]

				break
			}
		} else {
			return diag.Errorf("Ansible tower inventory cannot be read without name or value")
		}
	}

	var inventoryValue int
	if v, ok := inventory.Value.(int); ok {
		inventoryValue = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("inventory.Value", inventory.Value))
	}

	d.SetId(convert.IntToString(inventoryValue))
	d.Set("id", inventory.Value)
	d.Set("name", inventory.Name)

	return diags
}

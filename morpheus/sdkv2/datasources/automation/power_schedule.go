// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package automation

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/convert"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func DataSourcePowerSchedule() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus power schedule data source.",
		ReadContext: dataSourcePowerScheduleRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeInt,
				Description:   "The ID of the power schedule",
				Optional:      true,
				ConflictsWith: []string{"name"},
				Computed:      true,
			},
			"name": {
				Type:          schema.TypeString,
				Description:   "The name of the power schedule",
				Optional:      true,
				ConflictsWith: []string{"id"},
			},
		},
	}
}

func dataSourcePowerScheduleRead(
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

	var id int
	if v, ok := d.Get("id").(int); ok {
		id = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("id", d.Get("id")))
	}

	var resp *morpheus.Response
	var err error
	if id == 0 && name != "" {
		resp, err = client.FindPowerScheduleByName(name)
	} else if id != 0 {
		resp, err = client.GetPowerSchedule(int64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Power schedule cannot be read without name or id")
	}

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

	var result *morpheus.GetPowerScheduleResult
	if v, ok := resp.Result.(*morpheus.GetPowerScheduleResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.PowerSchedule == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("PowerSchedule"))
	}

	powerSchedule := result.PowerSchedule
	d.SetId(convert.Int64ToString(powerSchedule.ID))
	d.Set("name", powerSchedule.Name)

	return diags
}

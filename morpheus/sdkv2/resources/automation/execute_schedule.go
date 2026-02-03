// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package automation

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/convert"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func ResourceExecuteSchedule() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides an execution schedule resource",
		CreateContext: resourceExecuteScheduleCreate,
		ReadContext:   resourceExecuteScheduleRead,
		UpdateContext: resourceExecuteScheduleUpdate,
		DeleteContext: resourceExecuteScheduleDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the execute schedule",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the execute schedule",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the execute schedule",
				Optional:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the execute schedule is enabled",
				Optional:    true,
				Default:     true,
			},
			"time_zone": {
				Type:        schema.TypeString,
				Description: "The time zone used for scheduling",
				Required:    true,
			},
			"schedule": {
				Type:        schema.TypeString,
				Description: "The cron style syntax for the scheduled execution",
				Required:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceExecuteScheduleCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	var description string
	if v, ok := d.Get("description").(string); ok {
		description = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("description", d.Get("description")))
	}

	var enabled bool
	if v, ok := d.Get("enabled").(bool); ok {
		enabled = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("enabled", d.Get("enabled")))
	}

	var timeZone string
	if v, ok := d.Get("time_zone").(string); ok {
		timeZone = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("time_zone", d.Get("time_zone")))
	}

	var scheduleStr string
	if v, ok := d.Get("schedule").(string); ok {
		scheduleStr = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("schedule", d.Get("schedule")))
	}

	schedule := make(map[string]any)
	schedule["name"] = name
	schedule["description"] = description
	schedule["enabled"] = enabled
	schedule["scheduleType"] = "execute"
	schedule["scheduleTimezone"] = timeZone
	schedule["cron"] = scheduleStr

	req := &morpheus.Request{
		Body: map[string]any{
			"schedule": schedule,
		},
	}

	resp, err := client.CreateExecuteSchedule(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.CreateExecuteScheduleResult
	if v, ok := resp.Result.(*morpheus.CreateExecuteScheduleResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.ExecuteSchedule == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("ExecuteSchedule"))
	}

	executeScheduleResult := result.ExecuteSchedule
	d.SetId(convert.Int64ToString(executeScheduleResult.ID))

	diags = append(diags, resourceExecuteScheduleRead(ctx, d, meta)...)

	return diags
}

func resourceExecuteScheduleRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	id := d.Id()

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}

	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindExecuteScheduleByName(name)
	} else if id != "" {
		resp, err = client.GetExecuteSchedule(convert.StringToInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Execute schedule cannot be read without name or id")
	}

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)
			log.Printf("Forcing recreation of resource")
			d.SetId("")

			return diags
		} else {
			log.Printf("API FAILURE: %s - %s", resp, err)

			return diag.FromErr(err)
		}
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Body == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Body"))
	}

	var executeSchedule ExecuteSchedule
	json.Unmarshal(resp.Body, &executeSchedule)

	d.SetId(convert.IntToString(executeSchedule.Schedule.ID))
	d.Set("name", executeSchedule.Schedule.Name)
	d.Set("description", executeSchedule.Schedule.Description)
	d.Set("enabled", executeSchedule.Schedule.Enabled)
	d.Set("time_zone", executeSchedule.Schedule.Scheduletimezone)
	d.Set("schedule", executeSchedule.Schedule.Cron)

	return diags
}

func resourceExecuteScheduleUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	id := d.Id()

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}

	var description string
	if v, ok := d.Get("description").(string); ok {
		description = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("description", d.Get("description")))
	}

	var enabled bool
	if v, ok := d.Get("enabled").(bool); ok {
		enabled = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("enabled", d.Get("enabled")))
	}

	var timeZone string
	if v, ok := d.Get("time_zone").(string); ok {
		timeZone = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("time_zone", d.Get("time_zone")))
	}

	var scheduleStr string
	if v, ok := d.Get("schedule").(string); ok {
		scheduleStr = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("schedule", d.Get("schedule")))
	}

	schedule := make(map[string]any)
	schedule["name"] = name
	schedule["description"] = description
	schedule["enabled"] = enabled
	schedule["scheduleType"] = "execute"
	schedule["scheduleTimezone"] = timeZone
	schedule["cron"] = scheduleStr

	req := &morpheus.Request{
		Body: map[string]any{
			"schedule": schedule,
		},
	}

	resp, err := client.UpdateExecuteSchedule(convert.StringToInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.UpdateExecuteScheduleResult
	if v, ok := resp.Result.(*morpheus.UpdateExecuteScheduleResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.ExecuteSchedule == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("ExecuteSchedule"))
	}

	executeSchedule := result.ExecuteSchedule
	d.SetId(convert.Int64ToString(executeSchedule.ID))

	return resourceExecuteScheduleRead(ctx, d, meta)
}

func resourceExecuteScheduleDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	id := d.Id()

	req := &morpheus.Request{}
	resp, err := client.DeleteExecuteSchedule(convert.StringToInt64(id), req)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)

			return diag.FromErr(err)
		} else {
			log.Printf("API FAILURE: %s - %s", resp, err)

			return diag.FromErr(err)
		}
	}
	log.Printf("API RESPONSE: %s", resp)
	d.SetId("")

	return diags
}

type ExecuteSchedule struct {
	Schedule struct {
		ID               int    `json:"id"`
		Name             string `json:"name"`
		Description      string `json:"description"`
		Enabled          bool   `json:"enabled"`
		Scheduletype     string `json:"scheduleType"`
		Scheduletimezone string `json:"scheduleTimezone"`
		Cron             string `json:"cron"`
		Datecreated      string `json:"dateCreated"`
		Lastupdated      string `json:"lastUpdated"`
	} `json:"schedule"`
}

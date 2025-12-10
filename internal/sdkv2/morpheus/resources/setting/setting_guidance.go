// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package setting

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"
)

func ResourceSettingGuidance() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus guidance setting resource.",
		CreateContext: resourceSettingGuidanceCreate,
		ReadContext:   resourceSettingGuidanceRead,
		UpdateContext: resourceSettingGuidanceUpdate,
		DeleteContext: resourceSettingGuidanceDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the guidance settings",
				Computed:    true,
			},
			"power_settings_average_cpu": {
				Type:        schema.TypeInt,
				Description: "Shutdown will be recommended if the average CPU usage is below this value",
				Optional:    true,
				Computed:    true,
			},
			"power_settings_maximum_cpu": {
				Type:        schema.TypeInt,
				Description: "Shutdown will be recommended if the CPU usage never exceeds this value",
				Optional:    true,
				Computed:    true,
			},
			"power_settings_network_threshold": {
				Type:        schema.TypeInt,
				Description: "Shutdown will be recommended if the average network bandwidth is below this value",
				Optional:    true,
				Computed:    true,
			},
			"cpu_upsize_average_cpu": {
				Type:        schema.TypeInt,
				Description: "CPU up-size is recommended if the average CPU percentage exceeds this value",
				Optional:    true,
				Computed:    true,
			},
			"cpu_upsize_maximum_cpu": {
				Type:        schema.TypeInt,
				Description: "CPU up-size is recommended if the maximum CPU percentage exceeds this value",
				Optional:    true,
				Computed:    true,
			},
			"memory_upsize_minimum_free_memory": {
				Type:        schema.TypeInt,
				Description: "Memory up-size will be recommended if free memory dips below this value",
				Optional:    true,
				Computed:    true,
			},
			"memory_downsize_average_free_memory": {
				Type:        schema.TypeInt,
				Description: "Memory down-size is recommended if the average free memory is above this value",
				Optional:    true,
				Computed:    true,
			},
			"memory_downsize_maximum_free_memory": {
				Type:        schema.TypeInt,
				Description: "Memory down-size is recommended if free memory has never dipped below this value",
				Optional:    true,
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceSettingGuidanceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	guidanceSettings := make(map[string]any)

	powerAverageCPU, powerAverageCPUok := d.GetOk("power_settings_average_cpu")
	if powerAverageCPUok {
		guidanceSettings["cpuAvgCutoffPower"] = powerAverageCPU
	}

	powerMaxCPU, powerMaxCPUok := d.GetOk("power_settings_maximum_cpu")
	if powerMaxCPUok {
		guidanceSettings["cpuMaxCutoffPower"] = powerMaxCPU
	}

	powerNetwork, powerNetworkok := d.GetOk("power_settings_network_threshold")
	if powerNetworkok {
		guidanceSettings["networkCutoffPower"] = powerNetwork
	}

	upsizeAverageCPU, upsizeAverageCPUok := d.GetOk("cpu_upsize_average_cpu")
	if upsizeAverageCPUok {
		guidanceSettings["cpuUpAvgStandardCutoffRightSize"] = upsizeAverageCPU
	}

	upsizeMaxCPU, upsizeMaxCPUok := d.GetOk("cpu_upsize_maximum_cpu")
	if upsizeMaxCPUok {
		guidanceSettings["cpuUpMaxStandardCutoffRightSize"] = upsizeMaxCPU
	}

	upsizeMaxMemory, upsizeMaxMemoryok := d.GetOk("memory_upsize_minimum_free_memory")
	if upsizeMaxMemoryok {
		guidanceSettings["memoryUpAvgStandardCutoffRightSize"] = upsizeMaxMemory
	}

	downsizeAverageMemory, downsizeAverageMemoryok := d.GetOk("memory_downsize_average_free_memory")
	if downsizeAverageMemoryok {
		guidanceSettings["memoryDownAvgStandardCutoffRightSize"] = downsizeAverageMemory
	}

	downsizeMaxMemory, downsizeMaxMemoryok := d.GetOk("memory_downsize_maximum_free_memory")
	if downsizeMaxMemoryok {
		guidanceSettings["memoryDownMaxStandardCutoffRightSize"] = downsizeMaxMemory
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"guidanceSettings": guidanceSettings,
		},
	}

	resp, err := client.UpdateGuidanceSettings(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	d.SetId(convert.Int64ToString(1))

	diags = append(diags, resourceSettingGuidanceRead(ctx, d, meta)...)

	return diags
}

func resourceSettingGuidanceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	var resp *morpheus.Response
	var err error

	resp, err = client.GetGuidanceSettings(&morpheus.Request{})
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)
			log.Printf("Forcing recreation of resource")
			d.SetId("")

			return diags
		}

		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.GetGuidanceSettingsResult
	if v, ok := resp.Result.(*morpheus.GetGuidanceSettingsResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.GuidanceSettings == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("GuidanceSettings"))
	}

	guidanceSetting := result.GuidanceSettings
	d.SetId(convert.Int64ToString(1))
	d.Set("power_settings_average_cpu", guidanceSetting.CpuAvgCutoffPower)
	d.Set("power_settings_maximum_cpu", guidanceSetting.CpuMaxCutoffPower)
	d.Set("power_settings_network_threshold", guidanceSetting.NetworkCutoffPower)
	d.Set("cpu_upsize_average_cpu", guidanceSetting.CpuUpAvgStandardCutoffRightSize)
	d.Set("cpu_upsize_maximum_cpu", guidanceSetting.CpuUpMaxStandardCutoffRightSize)
	d.Set("memory_upsize_minimum_free_memory", guidanceSetting.MemoryUpAvgStandardCutoffRightSize)
	d.Set("memory_downsize_average_free_memory", guidanceSetting.MemoryDownAvgStandardCutoffRightSize)
	d.Set("memory_downsize_maximum_free_memory", guidanceSetting.MemoryDownMaxStandardCutoffRightSize)

	return diags
}

func resourceSettingGuidanceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	guidanceSettings := make(map[string]any)

	if d.HasChange("power_settings_average_cpu") {
		guidanceSettings["cpuAvgCutoffPower"] = d.Get("power_settings_average_cpu")
	}

	if d.HasChange("power_settings_maximum_cpu") {
		guidanceSettings["cpuMaxCutoffPower"] = d.Get("power_settings_maximum_cpu")
	}

	if d.HasChange("power_settings_network_threshold") {
		guidanceSettings["networkCutoffPower"] = d.Get("power_settings_network_threshold")
	}

	if d.HasChange("cpu_upsize_average_cpu") {
		guidanceSettings["cpuUpAvgStandardCutoffRightSize"] = d.Get("cpu_upsize_average_cpu")
	}

	if d.HasChange("cpu_upsize_maximum_cpu") {
		guidanceSettings["cpuUpMaxStandardCutoffRightSize"] = d.Get("cpu_upsize_maximum_cpu")
	}

	if d.HasChange("memory_upsize_minimum_free_memory") {
		guidanceSettings["memoryUpAvgStandardCutoffRightSize"] = d.Get("memory_upsize_minimum_free_memory")
	}

	if d.HasChange("memory_downsize_average_free_memory") {
		guidanceSettings["memoryDownAvgStandardCutoffRightSize"] = d.Get("memory_downsize_average_free_memory")
	}

	if d.HasChange("memory_downsize_maximum_free_memory") {
		guidanceSettings["memoryDownMaxStandardCutoffRightSize"] = d.Get("memory_downsize_maximum_free_memory")
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"guidanceSettings": guidanceSettings,
		},
	}

	resp, err := client.UpdateGuidanceSettings(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.UpdateGuidanceSettingsResult
	if v, ok := resp.Result.(*morpheus.UpdateGuidanceSettingsResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.GuidanceSettings == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("GuidanceSettings"))
	}

	d.SetId(convert.Int64ToString(1))

	return resourceSettingGuidanceRead(ctx, d, meta)
}

func resourceSettingGuidanceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	d.SetId("")

	return diags
}

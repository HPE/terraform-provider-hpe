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

func ResourceScaleThreshold() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus scale threshold resource.",
		CreateContext: resourceScaleThresholdCreate,
		ReadContext:   resourceScaleThresholdRead,
		UpdateContext: resourceScaleThresholdUpdate,
		DeleteContext: resourceScaleThresholdDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the scale threshold",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the scale threshold",
				Required:    true,
			},
			"auto_upscale": {
				Type:        schema.TypeBool,
				Description: "Whether to scale up the number of instances",
				Required:    true,
			},
			"auto_downscale": {
				Type:        schema.TypeBool,
				Description: "Whether to scale down the number of instances",
				Required:    true,
			},
			"min_count": {
				Type:        schema.TypeInt,
				Description: "The minimum number of instances to scale down to",
				Required:    true,
			},
			"max_count": {
				Type:        schema.TypeInt,
				Description: "The maximum number of instances to scale up to",
				Required:    true,
			},
			"enable_cpu_threshold": {
				Type:        schema.TypeBool,
				Description: "Whether scaling operations based upon cpu usage is enabled or not",
				Optional:    true,
				Computed:    true,
			},
			"min_cpu_percentage": {
				Type:        schema.TypeFloat,
				Description: "The minimum cpu percentage for scaling",
				Optional:    true,
				Computed:    true,
			},
			"max_cpu_percentage": {
				Type:        schema.TypeFloat,
				Description: "The maximum memory percentage for scaling",
				Optional:    true,
				Computed:    true,
			},
			"enable_memory_threshold": {
				Type:        schema.TypeBool,
				Description: "Whether scaling operations based upon memory usage is enabled or not",
				Optional:    true,
				Computed:    true,
			},
			"min_memory_percentage": {
				Type:        schema.TypeFloat,
				Description: "The minimum memory percentage for scaling",
				Optional:    true,
				Computed:    true,
			},
			"max_memory_percentage": {
				Type:        schema.TypeFloat,
				Description: "The maximum memory percentage for scaling",
				Optional:    true,
				Computed:    true,
			},
			"enable_disk_threshold": {
				Type:        schema.TypeBool,
				Description: "Whether scaling operations based upon disk usage is enabled or not",
				Optional:    true,
				Computed:    true,
			},
			"min_disk_percentage": {
				Type:        schema.TypeFloat,
				Description: "The minimum disk percentage for scaling",
				Optional:    true,
				Computed:    true,
			},
			"max_disk_percentage": {
				Type:        schema.TypeFloat,
				Description: "The maximum disk percentage for scaling",
				Optional:    true,
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceScaleThresholdCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	var autoUpscale bool
	if v, ok := d.Get("auto_upscale").(bool); ok {
		autoUpscale = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("auto_upscale", d.Get("auto_upscale")))
	}

	var autoDownscale bool
	if v, ok := d.Get("auto_downscale").(bool); ok {
		autoDownscale = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("auto_downscale", d.Get("auto_downscale")))
	}

	var minCount int
	if v, ok := d.Get("min_count").(int); ok {
		minCount = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("min_count", d.Get("min_count")))
	}

	var maxCount int
	if v, ok := d.Get("max_count").(int); ok {
		maxCount = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("max_count", d.Get("max_count")))
	}

	var enableCPUThreshold bool
	if v, ok := d.Get("enable_cpu_threshold").(bool); ok {
		enableCPUThreshold = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("enable_cpu_threshold", d.Get("enable_cpu_threshold")))
	}

	var minCPUPercentage float64
	if v, ok := d.Get("min_cpu_percentage").(float64); ok {
		minCPUPercentage = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("min_cpu_percentage", d.Get("min_cpu_percentage")))
	}

	var maxCPUPercentage float64
	if v, ok := d.Get("max_cpu_percentage").(float64); ok {
		maxCPUPercentage = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("max_cpu_percentage", d.Get("max_cpu_percentage")))
	}

	var enableMemoryThreshold bool
	if v, ok := d.Get("enable_memory_threshold").(bool); ok {
		enableMemoryThreshold = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("enable_memory_threshold", d.Get("enable_memory_threshold")))
	}

	var minMemoryPercentage float64
	if v, ok := d.Get("min_memory_percentage").(float64); ok {
		minMemoryPercentage = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("min_memory_percentage", d.Get("min_memory_percentage")))
	}

	var maxMemoryPercentage float64
	if v, ok := d.Get("max_memory_percentage").(float64); ok {
		maxMemoryPercentage = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("max_memory_percentage", d.Get("max_memory_percentage")))
	}

	var enableDiskThreshold bool
	if v, ok := d.Get("enable_disk_threshold").(bool); ok {
		enableDiskThreshold = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("enable_disk_threshold", d.Get("enable_disk_threshold")))
	}

	var minDiskPercentage float64
	if v, ok := d.Get("min_disk_percentage").(float64); ok {
		minDiskPercentage = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("min_disk_percentage", d.Get("min_disk_percentage")))
	}

	var maxDiskPercentage float64
	if v, ok := d.Get("max_disk_percentage").(float64); ok {
		maxDiskPercentage = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("max_disk_percentage", d.Get("max_disk_percentage")))
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"scaleThreshold": map[string]any{
				"name":          name,
				"autoUp":        autoUpscale,
				"autoDown":      autoDownscale,
				"minCount":      minCount,
				"maxCount":      maxCount,
				"cpuEnabled":    enableCPUThreshold,
				"minCpu":        minCPUPercentage,
				"maxCpu":        maxCPUPercentage,
				"memoryEnabled": enableMemoryThreshold,
				"minMemory":     minMemoryPercentage,
				"maxMemory":     maxMemoryPercentage,
				"diskEnabled":   enableDiskThreshold,
				"minDisk":       minDiskPercentage,
				"maxDisk":       maxDiskPercentage,
			},
		},
	}

	resp, err := client.CreateScaleThreshold(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.CreateScaleThresholdResult
	if v, ok := resp.Result.(*morpheus.CreateScaleThresholdResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.ScaleThreshold == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("ScaleThreshold"))
	}

	scaleThreshold := result.ScaleThreshold
	d.SetId(convert.Int64ToString(scaleThreshold.ID))

	diags = append(diags, resourceScaleThresholdRead(ctx, d, meta)...)

	return diags
}

func resourceScaleThresholdRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
		resp, err = client.FindScaleThresholdByName(name)
	} else if id != "" {
		resp, err = client.GetScaleThreshold(convert.StringToInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("ScaleThreshold cannot be read without name or id")
	}

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

	var result *morpheus.GetScaleThresholdResult
	if v, ok := resp.Result.(*morpheus.GetScaleThresholdResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.ScaleThreshold == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("ScaleThreshold"))
	}

	scaleThreshold := result.ScaleThreshold
	d.SetId(convert.Int64ToString(scaleThreshold.ID))
	d.Set("name", scaleThreshold.Name)
	d.Set("auto_upscale", scaleThreshold.AutoUp)
	d.Set("auto_downscale", scaleThreshold.AutoDown)
	d.Set("min_count", scaleThreshold.MinCount)
	d.Set("max_count", scaleThreshold.MaxCount)
	d.Set("enable_cpu_threshold", scaleThreshold.CpuEnabled)
	d.Set("min_cpu_percentage", scaleThreshold.MinCpu)
	d.Set("max_cpu_percentage", scaleThreshold.MaxCpu)
	d.Set("enable_memory_threshold", scaleThreshold.MemoryEnabled)
	d.Set("min_memory_percentage", scaleThreshold.MinMemory)
	d.Set("max_memory_percentage", scaleThreshold.MaxMemory)
	d.Set("enable_disk_threshold", scaleThreshold.DiskEnabled)
	d.Set("min_disk_percentage", scaleThreshold.MinDisk)
	d.Set("max_disk_percentage", scaleThreshold.MaxDisk)

	return diags
}

func resourceScaleThresholdUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	var autoUpscale bool
	if v, ok := d.Get("auto_upscale").(bool); ok {
		autoUpscale = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("auto_upscale", d.Get("auto_upscale")))
	}

	var autoDownscale bool
	if v, ok := d.Get("auto_downscale").(bool); ok {
		autoDownscale = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("auto_downscale", d.Get("auto_downscale")))
	}

	var minCount int
	if v, ok := d.Get("min_count").(int); ok {
		minCount = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("min_count", d.Get("min_count")))
	}

	var maxCount int
	if v, ok := d.Get("max_count").(int); ok {
		maxCount = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("max_count", d.Get("max_count")))
	}

	var enableCPUThreshold bool
	if v, ok := d.Get("enable_cpu_threshold").(bool); ok {
		enableCPUThreshold = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("enable_cpu_threshold", d.Get("enable_cpu_threshold")))
	}

	var minCPUPercentage float64
	if v, ok := d.Get("min_cpu_percentage").(float64); ok {
		minCPUPercentage = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("min_cpu_percentage", d.Get("min_cpu_percentage")))
	}

	var maxCPUPercentage float64
	if v, ok := d.Get("max_cpu_percentage").(float64); ok {
		maxCPUPercentage = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("max_cpu_percentage", d.Get("max_cpu_percentage")))
	}

	var enableMemoryThreshold bool
	if v, ok := d.Get("enable_memory_threshold").(bool); ok {
		enableMemoryThreshold = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("enable_memory_threshold", d.Get("enable_memory_threshold")))
	}

	var minMemoryPercentage float64
	if v, ok := d.Get("min_memory_percentage").(float64); ok {
		minMemoryPercentage = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("min_memory_percentage", d.Get("min_memory_percentage")))
	}

	var maxMemoryPercentage float64
	if v, ok := d.Get("max_memory_percentage").(float64); ok {
		maxMemoryPercentage = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("max_memory_percentage", d.Get("max_memory_percentage")))
	}

	var enableDiskThreshold bool
	if v, ok := d.Get("enable_disk_threshold").(bool); ok {
		enableDiskThreshold = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("enable_disk_threshold", d.Get("enable_disk_threshold")))
	}

	var minDiskPercentage float64
	if v, ok := d.Get("min_disk_percentage").(float64); ok {
		minDiskPercentage = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("min_disk_percentage", d.Get("min_disk_percentage")))
	}

	var maxDiskPercentage float64
	if v, ok := d.Get("max_disk_percentage").(float64); ok {
		maxDiskPercentage = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("max_disk_percentage", d.Get("max_disk_percentage")))
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"scaleThreshold": map[string]any{
				"name":          name,
				"autoUp":        autoUpscale,
				"autoDown":      autoDownscale,
				"minCount":      minCount,
				"maxCount":      maxCount,
				"cpuEnabled":    enableCPUThreshold,
				"minCpu":        minCPUPercentage,
				"maxCpu":        maxCPUPercentage,
				"memoryEnabled": enableMemoryThreshold,
				"minMemory":     minMemoryPercentage,
				"maxMemory":     maxMemoryPercentage,
				"diskEnabled":   enableDiskThreshold,
				"minDisk":       minDiskPercentage,
				"maxDisk":       maxDiskPercentage,
			},
		},
	}

	resp, err := client.UpdateScaleThreshold(convert.StringToInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.UpdateScaleThresholdResult
	if v, ok := resp.Result.(*morpheus.UpdateScaleThresholdResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.ScaleThreshold == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("ScaleThreshold"))
	}

	scaleThreshold := result.ScaleThreshold
	d.SetId(convert.Int64ToString(scaleThreshold.ID))

	return resourceScaleThresholdRead(ctx, d, meta)
}

func resourceScaleThresholdDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteScaleThreshold(convert.StringToInt64(id), req)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)

			return diag.FromErr(err)
		}

		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	d.SetId("")

	return diags
}

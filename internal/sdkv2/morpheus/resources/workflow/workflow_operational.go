// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package workflow

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"
)

func ResourceWorkflowOperational() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus operational workflow resource.",
		CreateContext: resourceWorkflowOperationalCreate,
		ReadContext:   resourceWorkflowOperationalRead,
		UpdateContext: resourceWorkflowOperationalUpdate,
		DeleteContext: resourceWorkflowOperationalDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the operational workflow",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the operational workflow",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the operational workflow",
				Optional:    true,
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "The organization labels associated with the workflow (Only supported on Morpheus 5.5.3 or higher)",
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"option_types": {
				Type:        schema.TypeList,
				Description: "The option types associated with the operational workflow",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"platform": {
				Type:         schema.TypeString,
				Description:  "The operating system platforms the operational workflow is supported to run on",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"all", "linux", "macos", "windows", ""}, false),
			},
			"allow_custom_config": {
				Type:        schema.TypeBool,
				Description: "Allow a custom configuration to be supplied",
				Optional:    true,
				Default:     false,
			},
			"visibility": {
				Type:         schema.TypeString,
				Description:  "Whether the operational workflow is visible in sub-tenants or not",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"private", "public", ""}, false),
				Default:      "private",
			},
			"task_ids": {
				Type:        schema.TypeList,
				Description: "A list of tasks ids associated with the operational workflow",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceWorkflowOperationalCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	var tasks []map[string]any
	if d.Get("task_ids") != nil {
		var taskList []any
		if v, ok := d.Get("task_ids").([]any); ok {
			taskList = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("task_ids", d.Get("task_ids")))
		}

		if taskList != nil {
			for i := 0; i < len(taskList); i++ {
				row := make(map[string]any)
				row["taskId"] = taskList[i]
				row["taskPhase"] = "operation"
				tasks = append(tasks, row)
			}
		}
	}

	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		if labelSet, ok := attr.(*schema.Set); ok {
			for _, s := range labelSet.List() {
				if labelStr, ok := s.(string); ok {
					labelsPayload = append(labelsPayload, labelStr)
				}
			}
		}
	}

	var visibility string
	if v, ok := d.Get("visibility").(string); ok {
		visibility = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("visibility", d.Get("visibility")))
	}

	var platform string
	if v, ok := d.Get("platform").(string); ok {
		platform = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("platform", d.Get("platform")))
	}

	var allowCustomConfig bool
	if v, ok := d.Get("allow_custom_config").(bool); ok {
		allowCustomConfig = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("allow_custom_config", d.Get("allow_custom_config")))
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"taskSet": map[string]any{
				"name":              name,
				"description":       description,
				"labels":            labelsPayload,
				"type":              "operation",
				"optionTypes":       d.Get("option_types"),
				"visibility":        visibility,
				"platform":          platform,
				"allowCustomConfig": allowCustomConfig,
				"tasks":             tasks,
			},
		},
	}

	resp, err := client.CreateTaskSet(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.CreateTaskSetResult
	if v, ok := resp.Result.(*morpheus.CreateTaskSetResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.TaskSet == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("TaskSet"))
	}

	environment := result.TaskSet
	d.SetId(convert.Int64ToString(environment.ID))

	diags = append(diags, resourceWorkflowOperationalRead(ctx, d, meta)...)

	return diags
}

func resourceWorkflowOperationalRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
		resp, err = client.FindTaskSetByName(name)
	} else if id != "" {
		resp, err = client.GetTaskSet(convert.StringToInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("TaskSet cannot be read without name or id")
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

	var result *morpheus.GetTaskSetResult
	if v, ok := resp.Result.(*morpheus.GetTaskSetResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.TaskSet == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("TaskSet"))
	}

	workflow := result.TaskSet
	d.SetId(convert.Int64ToString(workflow.ID))
	d.Set("name", workflow.Name)
	d.Set("description", workflow.Description)

	if workflow.Labels != nil {
		d.Set("labels", workflow.Labels)
	}

	var optionTypes []int64
	if workflow.OptionTypes != nil {
		for i := 0; i < len(workflow.OptionTypes); i++ {
			var option map[string]any
			if v, ok := workflow.OptionTypes[i].(map[string]any); ok {
				option = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("label", workflow.OptionTypes[i]))
			}

			var optionID int64
			if v, ok := option["id"].(float64); ok {
				optionID = int64(v)
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("id", option["id"]))
			}

			optionTypes = append(optionTypes, optionID)
		}
	}

	d.Set("option_types", optionTypes)
	d.Set("task_ids", workflow.Tasks)
	d.Set("visibility", workflow.Visibility)
	d.Set("allow_custom_config", workflow.AllowCustomConfig)
	d.Set("platform", workflow.Platform)

	return diags
}

func resourceWorkflowOperationalUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	var tasks []map[string]any
	if d.Get("task_ids") != nil {
		var taskList []any
		if v, ok := d.Get("task_ids").([]any); ok {
			taskList = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("task_ids", d.Get("task_ids")))
		}

		if taskList != nil {
			for i := 0; i < len(taskList); i++ {
				row := make(map[string]any)
				row["taskId"] = taskList[i]
				row["taskPhase"] = "operation"
				tasks = append(tasks, row)
			}
		}
	}

	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		if labelSet, ok := attr.(*schema.Set); ok {
			for _, s := range labelSet.List() {
				if labelStr, ok := s.(string); ok {
					labelsPayload = append(labelsPayload, labelStr)
				}
			}
		}
	}

	var visibility string
	if v, ok := d.Get("visibility").(string); ok {
		visibility = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("visibility", d.Get("visibility")))
	}

	var platform string
	if v, ok := d.Get("platform").(string); ok {
		platform = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("platform", d.Get("platform")))
	}

	var allowCustomConfig bool
	if v, ok := d.Get("allow_custom_config").(bool); ok {
		allowCustomConfig = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("allow_custom_config", d.Get("allow_custom_config")))
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"taskSet": map[string]any{
				"name":              name,
				"description":       description,
				"labels":            labelsPayload,
				"type":              "operation",
				"optionTypes":       d.Get("option_types"),
				"visibility":        visibility,
				"platform":          platform,
				"allowCustomConfig": allowCustomConfig,
				"tasks":             tasks,
			},
		},
	}

	resp, err := client.UpdateTaskSet(convert.StringToInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.UpdateTaskSetResult
	if v, ok := resp.Result.(*morpheus.UpdateTaskSetResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.TaskSet == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("TaskSet"))
	}

	taskSet := result.TaskSet
	d.SetId(convert.Int64ToString(taskSet.ID))

	return resourceWorkflowOperationalRead(ctx, d, meta)
}

func resourceWorkflowOperationalDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteTaskSet(convert.StringToInt64(id), req)
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

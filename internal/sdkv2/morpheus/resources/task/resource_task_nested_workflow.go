package task

import (
	"context"
	"log"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTaskNestedWorkflow() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus nested workflow task resource",
		CreateContext: resourceTaskNestedWorkflowCreate,
		ReadContext:   resourceTaskNestedWorkflowRead,
		UpdateContext: resourceTaskNestedWorkflowUpdate,
		DeleteContext: resourceTaskNestedWorkflowDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the nested workflow task",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the nested workflow task",
				Required:    true,
			},
			"code": {
				Type:        schema.TypeString,
				Description: "The code of the nested workflow task",
				Optional:    true,
				Computed:    true,
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "The organization labels associated with the task (Only supported on Morpheus 5.5.3 or higher)",
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"operational_workflow_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the operational workflow",
				Optional:    true,
				Computed:    true,
			},
			"operational_workflow_name": {
				Type:        schema.TypeString,
				Description: "The name of the operational workflow",
				Optional:    true,
				Computed:    true,
			},
			"retryable": {
				Type:        schema.TypeBool,
				Description: "Whether to retry the task if there is a failure",
				Optional:    true,
				Computed:    true,
			},
			"retry_count": {
				Type:        schema.TypeInt,
				Description: "The number of times to retry the task if there is a failure",
				Optional:    true,
				Default:     5,
			},
			"retry_delay_seconds": {
				Type:        schema.TypeInt,
				Description: "The number of seconds to wait between retry attempts",
				Optional:    true,
				Default:     10,
			},
			"allow_custom_config": {
				Type:        schema.TypeBool,
				Description: "Custom configuration data to pass during the execution of the shell script",
				Optional:    true,
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTaskNestedWorkflowCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	var client *morpheus.Client
	if clientAssert, ok := meta.(*morpheus.Client); ok {
		client = clientAssert
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var name string
	if nameValue, ok := d.Get("name").(string); ok {
		name = nameValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}

	taskOptions := make(map[string]any)

	var operationalWorkflowId int
	if operationalWorkflowIdValue, ok := d.Get("operational_workflow_id").(int); ok {
		operationalWorkflowId = operationalWorkflowIdValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("operational_workflow_id", d.Get("operational_workflow_id")))
	}
	taskOptions["operationalWorkflowId"] = operationalWorkflowId

	var operationalWorkflowName string
	if operationalWorkflowNameValue, ok := d.Get("operational_workflow_name").(string); ok {
		operationalWorkflowName = operationalWorkflowNameValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("operational_workflow_name", d.Get("operational_workflow_name")))
	}
	taskOptions["operationalWorkflowName"] = operationalWorkflowName

	taskType := make(map[string]any)
	taskType["code"] = "nestedWorkflow"

	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		if labelSet, ok := attr.(*schema.Set); ok {
			for _, s := range labelSet.List() {
				if labelStr, ok := s.(string); ok {
					labelsPayload = append(labelsPayload, labelStr)
				} else {
					return diag.FromErr(helpers.TypeAssertFailError("label", s))
				}
			}
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("labels", attr))
		}
	}

	var code string
	if codeValue, ok := d.Get("code").(string); ok {
		code = codeValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("code", d.Get("code")))
	}

	var retryable bool
	if retryableValue, ok := d.Get("retryable").(bool); ok {
		retryable = retryableValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("retryable", d.Get("retryable")))
	}

	var allowCustomConfig bool
	if allowCustomConfigValue, ok := d.Get("allow_custom_config").(bool); ok {
		allowCustomConfig = allowCustomConfigValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("allow_custom_config", d.Get("allow_custom_config")))
	}

	var retryCount int
	if retryCountValue, ok := d.Get("retry_count").(int); ok {
		retryCount = retryCountValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("retry_count", d.Get("retry_count")))
	}

	var retryDelaySeconds int
	if retryDelaySecondsValue, ok := d.Get("retry_delay_seconds").(int); ok {
		retryDelaySeconds = retryDelaySecondsValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("retry_delay_seconds", d.Get("retry_delay_seconds")))
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"task": map[string]any{
				"name":              name,
				"code":              code,
				"labels":            labelsPayload,
				"taskType":          taskType,
				"taskOptions":       taskOptions,
				"executeTarget":     "local",
				"retryable":         retryable,
				"retryCount":        retryCount,
				"retryDelaySeconds": retryDelaySeconds,
				"allowCustomConfig": allowCustomConfig,
			},
		},
	}
	resp, err := client.CreateTask(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	var result *morpheus.CreateTaskResult
	if resultAssert, ok := resp.Result.(*morpheus.CreateTaskResult); ok {
		result = resultAssert
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("resp.Result", resp.Result))
	}

	if result.Task == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Task"))
	}
	task := result.Task
	// Successfully created resource, now set id
	d.SetId(convert.Int64ToString(task.ID))

	resourceTaskNestedWorkflowRead(ctx, d, meta)

	return diags
}

func resourceTaskNestedWorkflowRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	var client *morpheus.Client
	if clientAssert, ok := meta.(*morpheus.Client); ok {
		client = clientAssert
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	id := d.Id()

	var name string
	if nameValue, ok := d.Get("name").(string); ok {
		name = nameValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindTaskByName(name)
	} else if id != "" {
		resp, err = client.GetTask(convert.StringToInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Task cannot be read without name or id")
	}

	if err != nil {
		// 404 is ok?
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

	// store resource data
	var result *morpheus.GetTaskResult
	if resultAssert, ok := resp.Result.(*morpheus.GetTaskResult); ok {
		result = resultAssert
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("resp.Result", resp.Result))
	}

	if result.Task == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Task"))
	}
	nestedWorkflowTask := result.Task
	d.SetId(convert.Int64ToString(nestedWorkflowTask.ID))
	d.Set("name", nestedWorkflowTask.Name)
	d.Set("code", nestedWorkflowTask.Code)
	d.Set("labels", nestedWorkflowTask.Labels)
	d.Set("operational_workflow_id", nestedWorkflowTask.TaskOptions.OperationalWorkflowId)
	d.Set("operational_workflow_name", nestedWorkflowTask.TaskOptions.OperationalWorkflowName)
	d.Set("retryable", nestedWorkflowTask.Retryable)
	d.Set("retry_count", nestedWorkflowTask.RetryCount)
	d.Set("retry_delay_seconds", nestedWorkflowTask.RetryDelaySeconds)
	d.Set("allow_custom_config", nestedWorkflowTask.AllowCustomConfig)

	return diags
}

func resourceTaskNestedWorkflowUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if clientAssert, ok := meta.(*morpheus.Client); ok {
		client = clientAssert
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	id := d.Id()

	var name string
	if nameValue, ok := d.Get("name").(string); ok {
		name = nameValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}

	taskOptions := make(map[string]any)

	var operationalWorkflowId int
	if operationalWorkflowIdValue, ok := d.Get("operational_workflow_id").(int); ok {
		operationalWorkflowId = operationalWorkflowIdValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("operational_workflow_id", d.Get("operational_workflow_id")))
	}
	taskOptions["operationalWorkflowId"] = operationalWorkflowId

	var operationalWorkflowName string
	if operationalWorkflowNameValue, ok := d.Get("operational_workflow_name").(string); ok {
		operationalWorkflowName = operationalWorkflowNameValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("operational_workflow_name", d.Get("operational_workflow_name")))
	}
	taskOptions["operationalWorkflowName"] = operationalWorkflowName

	taskType := make(map[string]any)
	taskType["code"] = "nestedWorkflow"

	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		if labelSet, ok := attr.(*schema.Set); ok {
			for _, s := range labelSet.List() {
				if labelStr, ok := s.(string); ok {
					labelsPayload = append(labelsPayload, labelStr)
				} else {
					return diag.FromErr(helpers.TypeAssertFailError("label", s))
				}
			}
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("labels", attr))
		}
	}

	var code string
	if codeValue, ok := d.Get("code").(string); ok {
		code = codeValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("code", d.Get("code")))
	}

	var retryable bool
	if retryableValue, ok := d.Get("retryable").(bool); ok {
		retryable = retryableValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("retryable", d.Get("retryable")))
	}

	var allowCustomConfig bool
	if allowCustomConfigValue, ok := d.Get("allow_custom_config").(bool); ok {
		allowCustomConfig = allowCustomConfigValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("allow_custom_config", d.Get("allow_custom_config")))
	}

	var retryCount int
	if retryCountValue, ok := d.Get("retry_count").(int); ok {
		retryCount = retryCountValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("retry_count", d.Get("retry_count")))
	}

	var retryDelaySeconds int
	if retryDelaySecondsValue, ok := d.Get("retry_delay_seconds").(int); ok {
		retryDelaySeconds = retryDelaySecondsValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("retry_delay_seconds", d.Get("retry_delay_seconds")))
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"task": map[string]any{
				"name":              name,
				"code":              code,
				"labels":            labelsPayload,
				"taskType":          taskType,
				"taskOptions":       taskOptions,
				"executeTarget":     "local",
				"retryable":         retryable,
				"retryCount":        retryCount,
				"retryDelaySeconds": retryDelaySeconds,
				"allowCustomConfig": allowCustomConfig,
			},
		},
	}
	resp, err := client.UpdateTask(convert.StringToInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	var result *morpheus.UpdateTaskResult
	if resultAssert, ok := resp.Result.(*morpheus.UpdateTaskResult); ok {
		result = resultAssert
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("resp.Result", resp.Result))
	}

	if result.Task == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Task"))
	}
	nestedWorkflowTask := result.Task
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(convert.Int64ToString(nestedWorkflowTask.ID))

	return resourceTaskNestedWorkflowRead(ctx, d, meta)
}

func resourceTaskNestedWorkflowDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	var client *morpheus.Client
	if clientAssert, ok := meta.(*morpheus.Client); ok {
		client = clientAssert
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteTask(convert.StringToInt64(id), req)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)

			return nil
		} else {
			log.Printf("API FAILURE: %s - %s", resp, err)

			return diag.FromErr(err)
		}
	}
	log.Printf("API RESPONSE: %s", resp)
	d.SetId("")

	return diags
}

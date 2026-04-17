package task

import (
	"context"
	"log"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/convert"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceTaskVRO() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus vRealize Orchestrator (vRO) task resource",
		CreateContext: resourceTaskVROCreate,
		ReadContext:   resourceTaskVRORead,
		UpdateContext: resourceTaskVROUpdate,
		DeleteContext: resourceTaskVRODelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the vRO workflow task",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the vRO workflow task",
				Required:    true,
			},
			"code": {
				Type:        schema.TypeString,
				Description: "The code of the vRO workflow task",
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
			"result_type": {
				Type:         schema.TypeString,
				Description:  "The expected result type (value, keyValue, json)",
				ValidateFunc: validation.StringInSlice([]string{"value", "keyValue", "json"}, false),
				Optional:     true,
				Computed:     true,
			},
			"vro_integration_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the vRO integration",
				Required:    true,
			},
			"vro_workflow_value": {
				Type:        schema.TypeInt,
				Description: "The value of the vRO workflow",
				Required:    true,
			},
			"body": {
				Type:        schema.TypeString,
				Description: "The JSON body to send to vRO",
				Optional:    true,
			},
			"execute_target": {
				Type:        schema.TypeString,
				Description: "The target that the vRO workflow will be executed on",
				Optional:    true,
			},
			"retryable": {
				Type:        schema.TypeBool,
				Description: "Whether to retry the task if there is a failure",
				Optional:    true,
				Default:     false,
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
				Description: "Custom configuration data to pass during the execution of the vRO workflow task",
				Optional:    true,
				Default:     false,
			},
			"visibility": {
				Type:         schema.TypeString,
				Description:  "The visibility of the task (private or public)",
				ValidateFunc: validation.StringInSlice([]string{"private", "public"}, false),
				Optional:     true,
				Computed:     true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTaskVROCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	taskType := make(map[string]any)
	taskType["code"] = "vro"

	taskOptions := make(map[string]any)

	var vroIntegrationId int
	if vroIntegrationIdValue, ok := d.Get("vro_integration_id").(int); ok {
		vroIntegrationId = vroIntegrationIdValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("vro_integration_id", d.Get("vro_integration_id")))
	}
	taskOptions["vroIntegrationId"] = vroIntegrationId

	var vroWorkflowValue int
	if vroWorkflowValueValue, ok := d.Get("vro_workflow_value").(int); ok {
		vroWorkflowValue = vroWorkflowValueValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("vro_workflow_value", d.Get("vro_workflow_value")))
	}
	taskOptions["vroWorkflow"] = vroWorkflowValue

	var body string
	if bodyValue, ok := d.Get("body").(string); ok {
		body = bodyValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("body", d.Get("body")))
	}
	taskOptions["vroBody"] = body

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

	var resultType string
	if resultTypeValue, ok := d.Get("result_type").(string); ok {
		resultType = resultTypeValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("result_type", d.Get("result_type")))
	}

	var executeTarget string
	if executeTargetValue, ok := d.Get("execute_target").(string); ok {
		executeTarget = executeTargetValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("execute_target", d.Get("execute_target")))
	}

	var retryable bool
	if retryableValue, ok := d.Get("retryable").(bool); ok {
		retryable = retryableValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("retryable", d.Get("retryable")))
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

	var allowCustomConfig bool
	if allowCustomConfigValue, ok := d.Get("allow_custom_config").(bool); ok {
		allowCustomConfig = allowCustomConfigValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("allow_custom_config", d.Get("allow_custom_config")))
	}

	var visibility string
	if visibilityValue, ok := d.Get("visibility").(string); ok {
		visibility = visibilityValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("visibility", d.Get("visibility")))
	}

	if visibility != "" {
		taskOptions["visibility"] = visibility
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"task": map[string]any{
				"name":              name,
				"code":              code,
				"labels":            labelsPayload,
				"taskType":          taskType,
				"taskOptions":       taskOptions,
				"resultType":        resultType,
				"executeTarget":     executeTarget,
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

	diags = append(diags, resourceTaskVRORead(ctx, d, meta)...)

	return diags
}

func resourceTaskVRORead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
	workflowTask := result.Task

	d.SetId(convert.Int64ToString(workflowTask.ID))
	d.Set("name", workflowTask.Name)
	d.Set("code", workflowTask.Code)
	d.Set("labels", workflowTask.Labels)
	d.Set("result_type", workflowTask.ResultType)
	d.Set("vro_integration_id", workflowTask.TaskOptions.VroIntegrationId)
	d.Set("vro_workflow_value", workflowTask.TaskOptions.VroWorkflow)
	d.Set("body", workflowTask.TaskOptions.VroBody)
	d.Set("execute_target", workflowTask.ExecuteTarget)
	d.Set("retryable", workflowTask.Retryable)
	d.Set("retry_count", workflowTask.RetryCount)
	d.Set("retry_delay_seconds", workflowTask.RetryDelaySeconds)
	d.Set("allow_custom_config", workflowTask.AllowCustomConfig)
	d.Set("visibility", workflowTask.Visibility)

	return diags
}

func resourceTaskVROUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	taskType := make(map[string]any)
	taskType["code"] = "vro"

	taskOptions := make(map[string]any)

	var vroIntegrationId int
	if vroIntegrationIdValue, ok := d.Get("vro_integration_id").(int); ok {
		vroIntegrationId = vroIntegrationIdValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("vro_integration_id", d.Get("vro_integration_id")))
	}
	taskOptions["vroIntegrationId"] = vroIntegrationId

	var vroWorkflowValue int
	if vroWorkflowValueValue, ok := d.Get("vro_workflow_value").(int); ok {
		vroWorkflowValue = vroWorkflowValueValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("vro_workflow_value", d.Get("vro_workflow_value")))
	}
	taskOptions["vroWorkflow"] = vroWorkflowValue

	var body string
	if bodyValue, ok := d.Get("body").(string); ok {
		body = bodyValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("body", d.Get("body")))
	}
	taskOptions["vroBody"] = body

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

	var resultType string
	if resultTypeValue, ok := d.Get("result_type").(string); ok {
		resultType = resultTypeValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("result_type", d.Get("result_type")))
	}

	var executeTarget string
	if executeTargetValue, ok := d.Get("execute_target").(string); ok {
		executeTarget = executeTargetValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("execute_target", d.Get("execute_target")))
	}

	var retryable bool
	if retryableValue, ok := d.Get("retryable").(bool); ok {
		retryable = retryableValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("retryable", d.Get("retryable")))
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

	var allowCustomConfig bool
	if allowCustomConfigValue, ok := d.Get("allow_custom_config").(bool); ok {
		allowCustomConfig = allowCustomConfigValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("allow_custom_config", d.Get("allow_custom_config")))
	}

	var visibility string
	if visibilityValue, ok := d.Get("visibility").(string); ok {
		visibility = visibilityValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("visibility", d.Get("visibility")))
	}

	if visibility != "" {
		taskOptions["visibility"] = visibility
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"task": map[string]any{
				"name":              name,
				"code":              code,
				"labels":            labelsPayload,
				"taskType":          taskType,
				"taskOptions":       taskOptions,
				"resultType":        resultType,
				"executeTarget":     executeTarget,
				"retryable":         retryable,
				"retryCount":        retryCount,
				"retryDelaySeconds": retryDelaySeconds,
				"allowCustomConfig": allowCustomConfig,
			},
		},
	}
	log.Printf("API REQUEST: %s", req)
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
	vroTask := result.Task
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(convert.Int64ToString(vroTask.ID))

	return resourceTaskVRORead(ctx, d, meta)
}

func resourceTaskVRODelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

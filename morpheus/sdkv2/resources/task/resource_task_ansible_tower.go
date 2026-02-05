package task

import (
	"context"
	"log"
	"strconv"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/convert"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func ResourceTaskAnsibleTower() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus Ansible Tower task resource",
		CreateContext: resourceTaskAnsibleTowerCreate,
		ReadContext:   resourceTaskAnsibleTowerRead,
		UpdateContext: resourceTaskAnsibleTowerUpdate,
		DeleteContext: resourceTaskAnsibleTowerDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the ansible tower task",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the ansible tower task",
				Required:    true,
			},
			"code": {
				Type:        schema.TypeString,
				Description: "The code of the ansible tower task",
				Optional:    true,
				Computed:    true,
			},
			"labels": {
				Type: schema.TypeSet,
				Description: "The organization labels associated with the ansible tower task " +
					"(Only supported on Morpheus 5.5.3 or higher)",
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ansible_tower_integration_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the ansible tower integration",
				Required:    true,
			},
			"ansible_tower_inventory_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the ansible tower inventory",
				Required:    true,
			},
			"group": {
				Type:        schema.TypeString,
				Description: "The name of a new or existing group in the inventory",
				Optional:    true,
				Computed:    true,
			},
			"job_template_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the ansible tower job template",
				Required:    true,
			},
			"scm_override": {
				Type:        schema.TypeString,
				Description: "The git reference override",
				Optional:    true,
				Computed:    true,
			},
			"execute_mode": {
				Type:         schema.TypeString,
				Description:  "The ansible tower execution mode (executeHost, executeGroup, executeAll, off)",
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"executeHost", "executeGroup", "executeAll", "off"}, false),
			},
			"execute_target": {
				Type:         schema.TypeString,
				Description:  "The target that the ansible tower job will be executed on (local, remote, resource)",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"local", "remote", "resource"}, false),
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
				Description: "Custom configuration data to pass during the execution of the ansible tower job task",
				Optional:    true,
				Default:     false,
			},
			"visibility": {
				Type:         schema.TypeString,
				Description:  "The visibility of the ansible tower task (public or private)",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"public", "private"}, false),
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTaskAnsibleTowerCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
	taskType["code"] = "ansibleTowerTask"

	taskOptions := make(map[string]any)

	var ansibleTowerIntegrationId int
	if ansibleTowerIntegrationIdValue, ok := d.Get("ansible_tower_integration_id").(int); ok {
		ansibleTowerIntegrationId = ansibleTowerIntegrationIdValue
	} else {
		return diag.FromErr(
			helpers.TypeAssertFailError(
				"ansible_tower_integration_id",
				d.Get("ansible_tower_integration_id"),
			),
		)
	}
	taskOptions["ansibleTowerIntegrationId"] = ansibleTowerIntegrationId

	var ansibleTowerInventoryId int
	if ansibleTowerInventoryIdValue, ok := d.Get("ansible_tower_inventory_id").(int); ok {
		ansibleTowerInventoryId = ansibleTowerInventoryIdValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("ansible_tower_inventory_id", d.Get("ansible_tower_inventory_id")))
	}
	taskOptions["ansibleTowerInventoryId"] = ansibleTowerInventoryId

	var group string
	if groupValue, ok := d.Get("group").(string); ok {
		group = groupValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("group", d.Get("group")))
	}
	taskOptions["ansibleGroup"] = group

	var jobTemplateId int
	if jobTemplateIdValue, ok := d.Get("job_template_id").(int); ok {
		jobTemplateId = jobTemplateIdValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("job_template_id", d.Get("job_template_id")))
	}
	taskOptions["ansibleTowerJobTemplateId"] = jobTemplateId

	var executeMode string
	if executeModeValue, ok := d.Get("execute_mode").(string); ok {
		executeMode = executeModeValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("execute_mode", d.Get("execute_mode")))
	}
	taskOptions["ansibleTowerExecuteMode"] = executeMode

	var scmOverride string
	if scmOverrideValue, ok := d.Get("scm_override").(string); ok {
		scmOverride = scmOverrideValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("scm_override", d.Get("scm_override")))
	}
	taskOptions["ansibleTowerGitRef"] = scmOverride

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

	req := &morpheus.Request{
		Body: map[string]any{
			"task": map[string]any{
				"name":              name,
				"code":              code,
				"labels":            labelsPayload,
				"taskType":          taskType,
				"taskOptions":       taskOptions,
				"executeTarget":     executeTarget,
				"retryable":         retryable,
				"retryCount":        retryCount,
				"retryDelaySeconds": retryDelaySeconds,
				"allowCustomConfig": allowCustomConfig,
				"visibility":        visibility,
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

	diags = append(diags, resourceTaskAnsibleTowerRead(ctx, d, meta)...)

	return diags
}

func resourceTaskAnsibleTowerRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
	ansibleTowerTask := result.Task

	d.SetId(convert.Int64ToString(ansibleTowerTask.ID))
	d.Set("name", ansibleTowerTask.Name)
	d.Set("code", ansibleTowerTask.Code)
	d.Set("labels", ansibleTowerTask.Labels)
	integrationId, err := strconv.Atoi(ansibleTowerTask.TaskOptions.AnsibleTowerIntegrationId)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("ansible_tower_integration_id", integrationId)
	inventoryId, err := strconv.Atoi(ansibleTowerTask.TaskOptions.AnsibleTowerInventoryId)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("ansible_tower_inventory_id", inventoryId)
	d.Set("group", ansibleTowerTask.TaskOptions.AnsibleGroup)
	d.Set("scm_override", ansibleTowerTask.TaskOptions.AnsibleTowerGitRef)
	d.Set("execute_mode", ansibleTowerTask.TaskOptions.AnsibleTowerExecuteMode)
	jobTemplateId, err := strconv.Atoi(ansibleTowerTask.TaskOptions.AnsibleTowerJobTemplateId)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("job_template_id", jobTemplateId)
	d.Set("execute_target", ansibleTowerTask.ExecuteTarget)
	d.Set("retryable", ansibleTowerTask.Retryable)
	d.Set("retry_count", ansibleTowerTask.RetryCount)
	d.Set("retry_delay_seconds", ansibleTowerTask.RetryDelaySeconds)
	d.Set("allow_custom_config", ansibleTowerTask.AllowCustomConfig)
	d.Set("visibility", ansibleTowerTask.Visibility)

	return diags
}

func resourceTaskAnsibleTowerUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
	taskType["code"] = "ansibleTowerTask"

	taskOptions := make(map[string]any)
	if d.HasChange("ansible_tower_integration_id") {
		var ansibleTowerIntegrationId int
		if ansibleTowerIntegrationIdValue, ok := d.Get("ansible_tower_integration_id").(int); ok {
			ansibleTowerIntegrationId = ansibleTowerIntegrationIdValue
		} else {
			return diag.FromErr(
				helpers.TypeAssertFailError("ansible_tower_integration_id",
					d.Get("ansible_tower_integration_id"),
				),
			)
		}
		taskOptions["ansibleTowerIntegrationId"] = ansibleTowerIntegrationId
	}
	if d.HasChange("ansible_tower_inventory_id") {
		var ansibleTowerInventoryId int
		if ansibleTowerInventoryIdValue, ok := d.Get("ansible_tower_inventory_id").(int); ok {
			ansibleTowerInventoryId = ansibleTowerInventoryIdValue
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("ansible_tower_inventory_id", d.Get("ansible_tower_inventory_id")))
		}
		taskOptions["ansibleTowerInventoryId"] = ansibleTowerInventoryId
	}
	if d.HasChange("group") {
		var group string
		if groupValue, ok := d.Get("group").(string); ok {
			group = groupValue
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("group", d.Get("group")))
		}
		taskOptions["ansibleGroup"] = group
	}
	if d.HasChange("job_template_id") {
		var jobTemplateId int
		if jobTemplateIdValue, ok := d.Get("job_template_id").(int); ok {
			jobTemplateId = jobTemplateIdValue
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("job_template_id", d.Get("job_template_id")))
		}
		taskOptions["ansibleTowerJobTemplateId"] = jobTemplateId
	}
	if d.HasChange("execute_mode") {
		var executeMode string
		if executeModeValue, ok := d.Get("execute_mode").(string); ok {
			executeMode = executeModeValue
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("execute_mode", d.Get("execute_mode")))
		}
		taskOptions["ansibleTowerExecuteMode"] = executeMode
	}
	if d.HasChange("scm_override") {
		var scmOverride string
		if scmOverrideValue, ok := d.Get("scm_override").(string); ok {
			scmOverride = scmOverrideValue
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("scm_override", d.Get("scm_override")))
		}
		taskOptions["ansibleTowerGitRef"] = scmOverride
	}

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

	req := &morpheus.Request{
		Body: map[string]any{
			"task": map[string]any{
				"name":              name,
				"code":              code,
				"labels":            labelsPayload,
				"taskType":          taskType,
				"taskOptions":       taskOptions,
				"executeTarget":     executeTarget,
				"retryable":         retryable,
				"retryCount":        retryCount,
				"retryDelaySeconds": retryDelaySeconds,
				"allowCustomConfig": allowCustomConfig,
				"visibility":        visibility,
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
	ansibleTowerTask := result.Task
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(convert.Int64ToString(ansibleTowerTask.ID))

	return resourceTaskAnsibleTowerRead(ctx, d, meta)
}

func resourceTaskAnsibleTowerDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

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

func ResourceTaskAnsiblePlaybook() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus ansible playbook task resource",
		CreateContext: resourceTaskAnsiblePlaybookCreate,
		ReadContext:   resourceTaskAnsiblePlaybookRead,
		UpdateContext: resourceTaskAnsiblePlaybookUpdate,
		DeleteContext: resourceTaskAnsiblePlaybookDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the ansible playbook task",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the ansible playbook task",
				Required:    true,
			},
			"code": {
				Type:        schema.TypeString,
				Description: "The code of the ansible playbook task",
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
			"ansible_repo_id": {
				Type:        schema.TypeString,
				Description: "The id of the ansible repo",
				Optional:    true,
			},
			"git_ref": {
				Type:        schema.TypeString,
				Description: "The git reference of the ansible repo to pull (main, master, etc.)",
				Optional:    true,
			},
			"playbook": {
				Type:        schema.TypeString,
				Description: "The name of the ansible playbook to execute",
				Required:    true,
			},
			"tags": {
				Type:        schema.TypeString,
				Description: "The tags to specify during execution of the ansible playbook",
				Optional:    true,
				Computed:    true,
			},
			"skip_tags": {
				Type:        schema.TypeString,
				Description: "The tags to skip during execution of the ansible playbook",
				Optional:    true,
				Computed:    true,
			},
			"command_options": {
				Type:        schema.TypeString,
				Description: "Additional commands options to pass during the execution of the ansible playbook",
				Optional:    true,
				Computed:    true,
			},
			"execute_target": {
				Type:        schema.TypeString,
				Description: "The target that the ansible playbook will be executed on",
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
				Description: "Custom configuration data to pass during the execution of the ansible playbook",
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

func resourceTaskAnsiblePlaybookCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	var ansibleRepoId string
	if ansibleRepoIdValue, ok := d.Get("ansible_repo_id").(string); ok {
		ansibleRepoId = ansibleRepoIdValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("ansible_repo_id", d.Get("ansible_repo_id")))
	}
	taskOptions["ansibleGitId"] = ansibleRepoId

	var gitRef string
	if gitRefValue, ok := d.Get("git_ref").(string); ok {
		gitRef = gitRefValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("git_ref", d.Get("git_ref")))
	}
	taskOptions["ansibleGitRef"] = gitRef

	var playbook string
	if playbookValue, ok := d.Get("playbook").(string); ok {
		playbook = playbookValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("playbook", d.Get("playbook")))
	}
	taskOptions["ansiblePlaybook"] = playbook

	var tags string
	if tagsValue, ok := d.Get("tags").(string); ok {
		tags = tagsValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("tags", d.Get("tags")))
	}
	taskOptions["ansibleTags"] = tags

	var skipTags string
	if skipTagsValue, ok := d.Get("skip_tags").(string); ok {
		skipTags = skipTagsValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("skip_tags", d.Get("skip_tags")))
	}
	taskOptions["ansibleSkipTags"] = skipTags

	var commandOptions string
	if commandOptionsValue, ok := d.Get("command_options").(string); ok {
		commandOptions = commandOptionsValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("command_options", d.Get("command_options")))
	}
	taskOptions["ansibleOptions"] = commandOptions

	taskType := make(map[string]any)
	taskType["code"] = "ansibleTask"

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
				"executeTarget":     executeTarget,
				"visibility":        visibility,
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

	diags = append(diags, resourceTaskAnsiblePlaybookRead(ctx, d, meta)...)

	return diags
}

func resourceTaskAnsiblePlaybookRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
	ansiblePlaybookTask := result.Task
	d.SetId(convert.Int64ToString(ansiblePlaybookTask.ID))
	d.Set("name", ansiblePlaybookTask.Name)
	d.Set("code", ansiblePlaybookTask.Code)
	d.Set("labels", ansiblePlaybookTask.Labels)
	d.Set("ansible_repo_id", ansiblePlaybookTask.TaskOptions.AnsibleGitId)
	d.Set("git_ref", ansiblePlaybookTask.TaskOptions.AnsibleGitRef)
	d.Set("playbook", ansiblePlaybookTask.TaskOptions.AnsiblePlaybook)
	d.Set("tags", ansiblePlaybookTask.TaskOptions.AnsibleTags)
	d.Set("skip_tags", ansiblePlaybookTask.TaskOptions.AnsibleSkipTags)
	d.Set("command_options", ansiblePlaybookTask.TaskOptions.AnsibleOptions)
	d.Set("execute_target", ansiblePlaybookTask.ExecuteTarget)
	d.Set("retryable", ansiblePlaybookTask.Retryable)
	d.Set("retry_count", ansiblePlaybookTask.RetryCount)
	d.Set("retry_delay_seconds", ansiblePlaybookTask.RetryDelaySeconds)
	d.Set("allow_custom_config", ansiblePlaybookTask.AllowCustomConfig)
	d.Set("visibility", ansiblePlaybookTask.Visibility)

	return diags
}

func resourceTaskAnsiblePlaybookUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	var ansibleRepoId string
	if ansibleRepoIdValue, ok := d.Get("ansible_repo_id").(string); ok {
		ansibleRepoId = ansibleRepoIdValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("ansible_repo_id", d.Get("ansible_repo_id")))
	}
	taskOptions["ansibleGitId"] = ansibleRepoId

	var gitRef string
	if gitRefValue, ok := d.Get("git_ref").(string); ok {
		gitRef = gitRefValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("git_ref", d.Get("git_ref")))
	}
	taskOptions["ansibleGitRef"] = gitRef

	var playbook string
	if playbookValue, ok := d.Get("playbook").(string); ok {
		playbook = playbookValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("playbook", d.Get("playbook")))
	}
	taskOptions["ansiblePlaybook"] = playbook

	var tags string
	if tagsValue, ok := d.Get("tags").(string); ok {
		tags = tagsValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("tags", d.Get("tags")))
	}
	taskOptions["ansibleTags"] = tags

	var skipTags string
	if skipTagsValue, ok := d.Get("skip_tags").(string); ok {
		skipTags = skipTagsValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("skip_tags", d.Get("skip_tags")))
	}
	taskOptions["ansibleSkipTags"] = skipTags

	var commandOptions string
	if commandOptionsValue, ok := d.Get("command_options").(string); ok {
		commandOptions = commandOptionsValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("command_options", d.Get("command_options")))
	}
	taskOptions["ansibleOptions"] = commandOptions

	taskType := make(map[string]any)
	taskType["code"] = "ansibleTask"

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

	var visibility string
	if visibilityValue, ok := d.Get("visibility").(string); ok {
		visibility = visibilityValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("visibility", d.Get("visibility")))
	}

	var executeTarget string
	if executeTargetValue, ok := d.Get("execute_target").(string); ok {
		executeTarget = executeTargetValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("execute_target", d.Get("execute_target")))
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
				"executeTarget":     executeTarget,
				"visibility":        visibility,
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
	ansiblePlaybookTask := result.Task
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(convert.Int64ToString(ansiblePlaybookTask.ID))

	return resourceTaskAnsiblePlaybookRead(ctx, d, meta)
}

func resourceTaskAnsiblePlaybookDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

package task

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"strings"
	"time"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/convert"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceTaskPowerShellScript() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus powershell script task resource",
		CreateContext: resourceTaskPowerShellScriptCreate,
		ReadContext:   resourceTaskPowerShellScriptRead,
		UpdateContext: resourceTaskPowerShellScriptUpdate,
		DeleteContext: resourceTaskPowerShellScriptDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the powershell script task",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the powershell script task",
				Required:    true,
			},
			"code": {
				Type:        schema.TypeString,
				Description: "The code of the powershell script task",
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
			"elevated_shell": {
				Type:        schema.TypeBool,
				Description: "Run the powershell script with elevated permissions",
				Optional:    true,
				Default:     false,
			},
			"source_type": {
				Type:         schema.TypeString,
				Description:  "The source of the powershell script (local, url or repository)",
				ValidateFunc: validation.StringInSlice([]string{"local", "url", "repository"}, false),
				Required:     true,
			},
			"script_content": {
				Type:        schema.TypeString,
				Description: "The content of the powershell script. Used when the local source type is specified",
				Optional:    true,
				Computed:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					oldPayload := strings.TrimSpace(old)
					newPayload := strings.TrimSpace(new)

					return oldPayload == newPayload
				},
				StateFunc: func(val any) string {
					return strings.TrimSpace(val.(string))
				},
			},
			"script_path": {
				Type:        schema.TypeString,
				Description: "The path of the powershell script, either the url or the path in the repository",
				Optional:    true,
				Computed:    true,
			},
			"repository_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the git repository integration",
				Optional:    true,
				Computed:    true,
			},
			"version_ref": {
				Type:        schema.TypeString,
				Description: "The git reference of the repository to pull (main, master, etc.)",
				Optional:    true,
				Computed:    true,
			},
			"execute_target": {
				Type:         schema.TypeString,
				Description:  "The execute target for the powershell script (local, remote or resource)",
				ValidateFunc: validation.StringInSlice([]string{"local", "remote", "resource"}, false),
				Default:      "local",
				Optional:     true,
			},
			"visibility": {
				Type:         schema.TypeString,
				Description:  "The visibility of the task (private or public)",
				ValidateFunc: validation.StringInSlice([]string{"private", "public"}, false),
				Optional:     true,
				Computed:     true,
			},
			"remote_target_host": {
				Type:        schema.TypeString,
				Description: "The hostname or ip address of the remote target",
				Optional:    true,
				Computed:    true,
			},
			"remote_target_port": {
				Type:        schema.TypeString,
				Description: "The port used to connect to the remote target",
				Optional:    true,
				Computed:    true,
			},
			"remote_target_username": {
				Type:        schema.TypeString,
				Description: "The username of the user account used to authenticate to the remote target",
				Optional:    true,
				Computed:    true,
			},
			"remote_target_password": {
				Type:        schema.TypeString,
				Description: "The password of the user account used to authenticate to the remote target",
				Optional:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					h := sha256.New()
					h.Write([]byte(new))
					sha256Hash := hex.EncodeToString(h.Sum(nil))

					return strings.EqualFold(old, sha256Hash)
				},
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
				Description: "Custom configuration data to pass during the execution of the shell script",
				Optional:    true,
				Default:     false,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTaskPowerShellScriptCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	sourceOptions := make(map[string]any)

	var scriptContent string
	if scriptContentValue, ok := d.Get("script_content").(string); ok {
		scriptContent = scriptContentValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("script_content", d.Get("script_content")))
	}
	if scriptContent != "" {
		sourceOptions["content"] = scriptContent
	}

	var scriptPath string
	if scriptPathValue, ok := d.Get("script_path").(string); ok {
		scriptPath = scriptPathValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("script_path", d.Get("script_path")))
	}
	if scriptPath != "" {
		sourceOptions["contentPath"] = scriptPath
	}

	var versionRef string
	if versionRefValue, ok := d.Get("version_ref").(string); ok {
		versionRef = versionRefValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("version_ref", d.Get("version_ref")))
	}
	sourceOptions["contentRef"] = versionRef

	var repositoryId int
	if repositoryIdValue, ok := d.Get("repository_id").(int); ok {
		repositoryId = repositoryIdValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("repository_id", d.Get("repository_id")))
	}
	sourceOptions["repository"] = map[string]any{
		"id": repositoryId,
	}

	var sourceType string
	if sourceTypeValue, ok := d.Get("source_type").(string); ok {
		sourceType = sourceTypeValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("source_type", d.Get("source_type")))
	}
	sourceOptions["sourceType"] = sourceType

	taskType := make(map[string]any)
	taskType["code"] = "winrmTask"

	taskOptions := make(map[string]any)

	var elevatedShell bool
	if elevatedShellValue, ok := d.Get("elevated_shell").(bool); ok {
		elevatedShell = elevatedShellValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("elevated_shell", d.Get("elevated_shell")))
	}
	if elevatedShell {
		taskOptions["winrm.elevated"] = "on"
	} else {
		taskOptions["winrm.elevated"] = nil
	}

	var remoteTargetHost string
	if remoteTargetHostValue, ok := d.Get("remote_target_host").(string); ok {
		remoteTargetHost = remoteTargetHostValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("remote_target_host", d.Get("remote_target_host")))
	}
	if remoteTargetHost != "" {
		taskOptions["host"] = remoteTargetHost
	}

	var remoteTargetPort string
	if remoteTargetPortValue, ok := d.Get("remote_target_port").(string); ok {
		remoteTargetPort = remoteTargetPortValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("remote_target_port", d.Get("remote_target_port")))
	}
	if remoteTargetPort != "" {
		taskOptions["port"] = remoteTargetPort
	}

	var remoteTargetUsername string
	if remoteTargetUsernameValue, ok := d.Get("remote_target_username").(string); ok {
		remoteTargetUsername = remoteTargetUsernameValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("remote_target_username", d.Get("remote_target_username")))
	}
	if remoteTargetUsername != "" {
		taskOptions["username"] = remoteTargetUsername
	}

	var remoteTargetPassword string
	if remoteTargetPasswordValue, ok := d.Get("remote_target_password").(string); ok {
		remoteTargetPassword = remoteTargetPasswordValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("remote_target_password", d.Get("remote_target_password")))
	}
	if remoteTargetPassword != "" {
		taskOptions["password"] = remoteTargetPassword
	}

	var visibility string
	if visibilityValue, ok := d.Get("visibility").(string); ok {
		visibility = visibilityValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("visibility", d.Get("visibility")))
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

	req := &morpheus.Request{
		Body: map[string]any{
			"task": map[string]any{
				"name":              name,
				"code":              code,
				"labels":            labelsPayload,
				"file":              sourceOptions,
				"taskType":          taskType,
				"taskOptions":       taskOptions,
				"resultType":        resultType,
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
	log.Printf("Task ID: %s", convert.Int64ToString(task.ID))

	diags = append(diags, resourceTaskPowerShellScriptRead(ctx, d, meta)...)

	return diags
}

func resourceTaskPowerShellScriptRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
	powerShellScriptTask := result.Task
	d.SetId(convert.Int64ToString(powerShellScriptTask.ID))
	d.Set("name", powerShellScriptTask.Name)
	d.Set("code", powerShellScriptTask.Code)
	d.Set("labels", powerShellScriptTask.Labels)
	d.Set("result_type", powerShellScriptTask.ResultType)
	d.Set("source_type", powerShellScriptTask.File.SourceType)
	d.Set("script_content", powerShellScriptTask.File.Content)
	d.Set("script_path", powerShellScriptTask.File.ContentPath)
	d.Set("version_ref", powerShellScriptTask.File.ContentRef)
	d.Set("execute_target", powerShellScriptTask.ExecuteTarget)
	d.Set("repository_id", powerShellScriptTask.File.Repository.ID)
	if powerShellScriptTask.TaskOptions.WinrmElevated == "on" {
		d.Set("elevated_shell", true)
	} else {
		d.Set("elevated_shell", false)
	}
	d.Set("remote_target_host", powerShellScriptTask.TaskOptions.Host)
	d.Set("remote_target_port", powerShellScriptTask.TaskOptions.Port)
	d.Set("remote_target_username", powerShellScriptTask.TaskOptions.Username)
	d.Set("remote_target_password", powerShellScriptTask.TaskOptions.PasswordHash)
	d.Set("retryable", powerShellScriptTask.Retryable)
	d.Set("retry_count", powerShellScriptTask.RetryCount)
	d.Set("retry_delay_seconds", powerShellScriptTask.RetryDelaySeconds)
	d.Set("allow_custom_config", powerShellScriptTask.AllowCustomConfig)
	d.Set("visibility", powerShellScriptTask.Visibility)

	return diags
}

func resourceTaskPowerShellScriptUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	sourceOptions := make(map[string]any)

	var scriptContent string
	if scriptContentValue, ok := d.Get("script_content").(string); ok {
		scriptContent = scriptContentValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("script_content", d.Get("script_content")))
	}
	if scriptContent != "" {
		sourceOptions["content"] = scriptContent
	}

	var scriptPath string
	if scriptPathValue, ok := d.Get("script_path").(string); ok {
		scriptPath = scriptPathValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("script_path", d.Get("script_path")))
	}
	if scriptPath != "" {
		sourceOptions["contentPath"] = scriptPath
	}

	var versionRef string
	if versionRefValue, ok := d.Get("version_ref").(string); ok {
		versionRef = versionRefValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("version_ref", d.Get("version_ref")))
	}
	sourceOptions["contentRef"] = versionRef

	var repositoryId int
	if repositoryIdValue, ok := d.Get("repository_id").(int); ok {
		repositoryId = repositoryIdValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("repository_id", d.Get("repository_id")))
	}
	sourceOptions["repository"] = map[string]any{
		"id": repositoryId,
	}

	var sourceType string
	if sourceTypeValue, ok := d.Get("source_type").(string); ok {
		sourceType = sourceTypeValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("source_type", d.Get("source_type")))
	}
	sourceOptions["sourceType"] = sourceType

	taskType := make(map[string]any)
	taskType["code"] = "winrmTask"

	taskOptions := make(map[string]any)

	var elevatedShell bool
	if elevatedShellValue, ok := d.Get("elevated_shell").(bool); ok {
		elevatedShell = elevatedShellValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("elevated_shell", d.Get("elevated_shell")))
	}
	if elevatedShell {
		taskOptions["winrm.elevated"] = "on"
	} else {
		taskOptions["winrm.elevated"] = nil
	}

	if d.HasChange("remote_target_host") {
		if remoteTargetHost, ok := d.Get("remote_target_host").(string); ok {
			taskOptions["host"] = remoteTargetHost
		}
	}

	if d.HasChange("remote_target_port") {
		if remoteTargetPort, ok := d.Get("remote_target_port").(string); ok {
			taskOptions["port"] = remoteTargetPort
		}
	}
	if d.HasChange("remote_target_username") {
		if remoteTargetUsername, ok := d.Get("remote_target_username").(string); ok {
			taskOptions["username"] = remoteTargetUsername
		}
	}
	if d.HasChange("remote_target_password") {
		if remoteTargetPassword, ok := d.Get("remote_target_password").(string); ok {
			taskOptions["password"] = remoteTargetPassword
		}
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

	var visibility string
	if visibilityValue, ok := d.Get("visibility").(string); ok {
		visibility = visibilityValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("visibility", d.Get("visibility")))
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

	req := &morpheus.Request{
		Body: map[string]any{
			"task": map[string]any{
				"name":              name,
				"code":              code,
				"labels":            labelsPayload,
				"file":              sourceOptions,
				"taskType":          taskType,
				"taskOptions":       taskOptions,
				"resultType":        resultType,
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
	task := result.Task
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(convert.Int64ToString(task.ID))

	return resourceTaskPowerShellScriptRead(ctx, d, meta)
}

func resourceTaskPowerShellScriptDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

type PowerShellScript struct {
	Task struct {
		ID        int    `json:"id"`
		Accountid int    `json:"accountId"`
		Name      string `json:"name"`
		Code      string `json:"code"`
		Tasktype  struct {
			ID   int    `json:"id"`
			Code string `json:"code"`
			Name string `json:"name"`
		} `json:"taskType"`
		Taskoptions struct {
			Port              string `json:"port"`
			Host              string `json:"host"`
			Password          string `json:"password"`
			PasswordHash      string `json:"passwordHash"`
			Username          string `json:"username"`
			WinrmElevated     string `json:"winrm.elevated"`
			LocalScriptGitRef string `json:"localScriptGitRef"`
		}
		File struct {
			ID          int    `json:"id"`
			Sourcetype  string `json:"sourceType"`
			Contentref  string `json:"contentRef"`
			Contentpath string `json:"contentPath"`
			Repository  struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"repository"`
			Content any `json:"content"`
		} `json:"file"`
		Resulttype        string    `json:"resultType"`
		Executetarget     string    `json:"executeTarget"`
		Retryable         bool      `json:"retryable"`
		Retrycount        int       `json:"retryCount"`
		Retrydelayseconds int       `json:"retryDelaySeconds"`
		Allowcustomconfig bool      `json:"allowCustomConfig"`
		Datecreated       time.Time `json:"dateCreated"`
		Lastupdated       time.Time `json:"lastUpdated"`
	} `json:"task"`
}

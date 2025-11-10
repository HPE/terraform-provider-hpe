package task

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"strings"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceTaskShellScript() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus shell script task resource",
		CreateContext: resourceTaskShellScriptCreate,
		ReadContext:   resourceTaskShellScriptRead,
		UpdateContext: resourceTaskShellScriptUpdate,
		DeleteContext: resourceTaskShellScriptDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the shell script task",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the shell script task",
				Required:    true,
			},
			"code": {
				Type:        schema.TypeString,
				Description: "The code of the shell script task",
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
			"sudo": {
				Type:        schema.TypeBool,
				Description: "Whether to run the script as sudo",
				Optional:    true,
				Computed:    true,
			},
			"source_type": {
				Type:         schema.TypeString,
				Description:  "The source of the shell script (local, url or repository)",
				ValidateFunc: validation.StringInSlice([]string{"local", "url", "repository"}, false),
				Required:     true,
			},
			"script_content": {
				Type:        schema.TypeString,
				Description: "The content of the shell script. Used when the local source type is specified",
				Optional:    true,
				Computed:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					oldPayload := strings.TrimSuffix(old, "\n")
					newPayload := strings.TrimSuffix(new, "\n")

					return oldPayload == newPayload
				},
				StateFunc: func(val any) string {
					if str, ok := val.(string); ok {
						return strings.TrimSuffix(str, "\n")
					}

					return ""
				},
			},
			"script_path": {
				Type:        schema.TypeString,
				Description: "The path of the shell script, either the url or the path in the repository",
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
				Description:  "The execute target of the shell script (local, remote, resource)",
				ValidateFunc: validation.StringInSlice([]string{"local", "remote", "resource"}, false),
				Optional:     true,
				Computed:     true,
			},
			"local_repository_id": {
				Type:        schema.TypeString,
				Description: "The ID of the local git repository",
				Optional:    true,
				Computed:    true,
			},
			"local_repository_ref": {
				Type:        schema.TypeString,
				Description: "The git reference of the repository to pull (main, master, etc.)",
				Optional:    true,
				Computed:    true,
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
				Computed: true,
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

func resourceTaskShellScriptCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
	if scriptContent, ok := d.Get("script_content").(string); ok {
		if scriptContent != "" {
			sourceOptions["content"] = scriptContent
		}
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("script_content", d.Get("script_content")))
	}
	if scriptPath, ok := d.Get("script_path").(string); ok {
		if scriptPath != "" {
			sourceOptions["contentPath"] = scriptPath
		}
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("script_path", d.Get("script_path")))
	}
	if versionRef, ok := d.Get("version_ref").(string); ok {
		sourceOptions["contentRef"] = versionRef
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("version_ref", d.Get("version_ref")))
	}
	sourceOptions["repository"] = map[string]any{
		"id": d.Get("repository_id"),
	}
	if sourceType, ok := d.Get("source_type").(string); ok {
		sourceOptions["sourceType"] = sourceType
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("source_type", d.Get("source_type")))
	}

	taskType := make(map[string]any)
	taskType["code"] = "script"

	taskOptions := make(map[string]any)
	if sudo, ok := d.Get("sudo").(bool); ok {
		if sudo {
			taskOptions["shell.sudo"] = "on"
		} else {
			taskOptions["shell.sudo"] = nil
		}
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("sudo", d.Get("sudo")))
	}
	if remoteTargetHost, ok := d.Get("remote_target_host").(string); ok {
		if remoteTargetHost != "" {
			taskOptions["host"] = remoteTargetHost
		}
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("remote_target_host", d.Get("remote_target_host")))
	}
	if remoteTargetPort, ok := d.Get("remote_target_port").(string); ok {
		if remoteTargetPort != "" {
			taskOptions["port"] = remoteTargetPort
		}
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("remote_target_port", d.Get("remote_target_port")))
	}
	if remoteTargetUsername, ok := d.Get("remote_target_username").(string); ok {
		if remoteTargetUsername != "" {
			taskOptions["username"] = remoteTargetUsername
		}
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("remote_target_username", d.Get("remote_target_username")))
	}
	if remoteTargetPassword, ok := d.Get("remote_target_password").(string); ok {
		if remoteTargetPassword != "" {
			taskOptions["password"] = remoteTargetPassword
		}
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("remote_target_password", d.Get("remote_target_password")))
	}
	if localRepositoryId, ok := d.Get("local_repository_id").(string); ok {
		if localRepositoryId != "" {
			taskOptions["localScriptGitId"] = localRepositoryId
		}
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("local_repository_id", d.Get("local_repository_id")))
	}
	if localRepositoryRef, ok := d.Get("local_repository_ref").(string); ok {
		if localRepositoryRef != "" {
			taskOptions["localScriptGitRef"] = localRepositoryRef
		}
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("local_repository_ref", d.Get("local_repository_ref")))
	}
	if visibility, ok := d.Get("visibility").(string); ok {
		if visibility != "" {
			taskOptions["visibility"] = visibility
		}
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

	var localRepositoryRef string
	if localRepositoryRefValue, ok := d.Get("local_repository_ref").(string); ok {
		localRepositoryRef = localRepositoryRefValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("local_repository_ref", d.Get("local_repository_ref")))
	}

	var localRepositoryId string
	if localRepositoryIdValue, ok := d.Get("local_repository_id").(string); ok {
		localRepositoryId = localRepositoryIdValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("local_repository_id", d.Get("local_repository_id")))
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
				"localScriptGitRef": localRepositoryRef,
				"localScriptGitId":  localRepositoryId,
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

	return resourceTaskShellScriptRead(ctx, d, meta)
}

func resourceTaskShellScriptRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	shellScriptTask := result.Task

	d.SetId(convert.Int64ToString(shellScriptTask.ID))
	d.Set("name", shellScriptTask.Name)
	d.Set("code", shellScriptTask.Code)
	d.Set("labels", shellScriptTask.Labels)
	d.Set("result_type", shellScriptTask.ResultType)
	d.Set("source_type", shellScriptTask.File.SourceType)
	d.Set("script_content", shellScriptTask.File.Content)
	d.Set("script_path", shellScriptTask.File.ContentPath)
	d.Set("version_ref", shellScriptTask.File.ContentRef)
	d.Set("execute_target", shellScriptTask.ExecuteTarget)
	d.Set("local_repository_id", shellScriptTask.TaskOptions.LocalScriptGitId)
	d.Set("local_repository_ref", shellScriptTask.TaskOptions.LocalScriptGitRef)
	d.Set("repository_id", shellScriptTask.File.Repository.ID)
	if shellScriptTask.TaskOptions.ShellSudo == "on" {
		d.Set("sudo", true)
	} else {
		d.Set("sudo", false)
	}
	d.Set("remote_target_host", shellScriptTask.TaskOptions.Host)
	d.Set("remote_target_port", shellScriptTask.TaskOptions.Port)
	d.Set("remote_target_username", shellScriptTask.TaskOptions.Username)
	d.Set("remote_target_password", shellScriptTask.TaskOptions.PasswordHash)
	d.Set("retryable", shellScriptTask.Retryable)
	d.Set("retry_count", shellScriptTask.RetryCount)
	d.Set("retry_delay_seconds", shellScriptTask.RetryDelaySeconds)
	d.Set("allow_custom_config", shellScriptTask.AllowCustomConfig)
	d.Set("visibility", shellScriptTask.Visibility)

	return diags
}

func resourceTaskShellScriptUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
	if scriptContent, ok := d.Get("script_content").(string); ok {
		if scriptContent != "" {
			sourceOptions["content"] = scriptContent
		}
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("script_content", d.Get("script_content")))
	}
	if scriptPath, ok := d.Get("script_path").(string); ok {
		if scriptPath != "" {
			sourceOptions["contentPath"] = scriptPath
		}
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("script_path", d.Get("script_path")))
	}
	if versionRef, ok := d.Get("version_ref").(string); ok {
		sourceOptions["contentRef"] = versionRef
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("version_ref", d.Get("version_ref")))
	}
	sourceOptions["repository"] = map[string]any{
		"id": d.Get("repository_id"),
	}
	if sourceType, ok := d.Get("source_type").(string); ok {
		sourceOptions["sourceType"] = sourceType
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("source_type", d.Get("source_type")))
	}

	taskType := make(map[string]any)
	taskType["code"] = "script"

	taskOptions := make(map[string]any)
	if sudo, ok := d.Get("sudo").(bool); ok {
		if sudo {
			taskOptions["shell.sudo"] = "on"
		} else {
			taskOptions["shell.sudo"] = nil
		}
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("sudo", d.Get("sudo")))
	}
	if d.HasChange("remote_target_host") {
		if remoteTargetHost, ok := d.Get("remote_target_host").(string); ok {
			taskOptions["host"] = remoteTargetHost
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("remote_target_host", d.Get("remote_target_host")))
		}
	}
	if d.HasChange("remote_target_port") {
		if remoteTargetPort, ok := d.Get("remote_target_port").(string); ok {
			taskOptions["port"] = remoteTargetPort
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("remote_target_port", d.Get("remote_target_port")))
		}
	}
	if d.HasChange("remote_target_username") {
		if remoteTargetUsername, ok := d.Get("remote_target_username").(string); ok {
			taskOptions["username"] = remoteTargetUsername
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("remote_target_username", d.Get("remote_target_username")))
		}
	}
	if d.HasChange("remote_target_password") {
		if remoteTargetPassword, ok := d.Get("remote_target_password").(string); ok {
			taskOptions["password"] = remoteTargetPassword
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("remote_target_password", d.Get("remote_target_password")))
		}
	}
	if d.HasChange("local_repository_id") {
		if localRepositoryId, ok := d.Get("local_repository_id").(string); ok {
			taskOptions["localScriptGitId"] = localRepositoryId
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("local_repository_id", d.Get("local_repository_id")))
		}
	}
	if d.HasChange("local_repository_ref") {
		if localRepositoryRef, ok := d.Get("local_repository_ref").(string); ok {
			taskOptions["localScriptGitRef"] = localRepositoryRef
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("local_repository_ref", d.Get("local_repository_ref")))
		}
	}
	if d.HasChange("visibility") {
		if visibility, ok := d.Get("visibility").(string); ok {
			taskOptions["visibility"] = visibility
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("visibility", d.Get("visibility")))
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

	var localRepositoryRef string
	if localRepositoryRefValue, ok := d.Get("local_repository_ref").(string); ok {
		localRepositoryRef = localRepositoryRefValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("local_repository_ref", d.Get("local_repository_ref")))
	}

	var localRepositoryId string
	if localRepositoryIdValue, ok := d.Get("local_repository_id").(string); ok {
		localRepositoryId = localRepositoryIdValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("local_repository_id", d.Get("local_repository_id")))
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
				"localScriptGitRef": localRepositoryRef,
				"localScriptGitId":  localRepositoryId,
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

	shellScriptTask := result.Task
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(convert.Int64ToString(shellScriptTask.ID))

	return resourceTaskShellScriptRead(ctx, d, meta)
}

func resourceTaskShellScriptDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

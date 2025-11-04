package task

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceTaskGroovyScript() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus groovy script task resource",
		CreateContext: resourceTaskGroovyScriptCreate,
		ReadContext:   resourceTaskGroovyScriptRead,
		UpdateContext: resourceTaskGroovyScriptUpdate,
		DeleteContext: resourceTaskGroovyScriptDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the groovy script task",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the groovy script task",
				Required:    true,
			},
			"code": {
				Type:        schema.TypeString,
				Description: "The code of the groovy script task",
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
			"source_type": {
				Type:         schema.TypeString,
				Description:  "The source of the groovy script (local, url or repository)",
				ValidateFunc: validation.StringInSlice([]string{"local", "url", "repository"}, false),
				Required:     true,
			},
			"script_content": {
				Type:        schema.TypeString,
				Description: "The content of the groovy script. Used when the local source type is specified",
				Optional:    true,
				Computed:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					oldPayload := strings.TrimSuffix(old, "\n")
					newPayload := strings.TrimSuffix(new, "\n")

					return oldPayload == newPayload
				},
				StateFunc: func(val any) string {
					return strings.TrimSuffix(val.(string), "\n")
				},
			},
			"script_path": {
				Type:        schema.TypeString,
				Description: "The path of the groovy script, either the url or the path in the repository",
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
				Description: "Custom configuration data to pass during the execution of the groovy script",
				Optional:    true,
				Default:     false,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTaskGroovyScriptCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
	taskType["code"] = "groovyTask"

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
				"resultType":        resultType,
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

	task := result.Task
	// Successfully created resource, now set id
	d.SetId(convert.Int64ToString(task.ID))

	resourceTaskGroovyScriptRead(ctx, d, meta)

	return diags
}

func resourceTaskGroovyScriptRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	groovyScriptTask := result.Task
	d.SetId(convert.Int64ToString(groovyScriptTask.ID))
	d.Set("name", groovyScriptTask.Name)
	d.Set("code", groovyScriptTask.Code)
	d.Set("labels", groovyScriptTask.Labels)
	d.Set("result_type", groovyScriptTask.ResultType)
	d.Set("source_type", groovyScriptTask.File.SourceType)
	d.Set("script_content", groovyScriptTask.File.Content)
	d.Set("script_path", groovyScriptTask.File.ContentPath)
	d.Set("version_ref", groovyScriptTask.File.ContentRef)
	d.Set("repository_id", groovyScriptTask.File.Repository.ID)
	d.Set("retryable", groovyScriptTask.Retryable)
	d.Set("retry_count", groovyScriptTask.RetryCount)
	d.Set("retry_delay_seconds", groovyScriptTask.RetryDelaySeconds)
	d.Set("allow_custom_config", groovyScriptTask.AllowCustomConfig)

	return diags
}

func resourceTaskGroovyScriptUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
	taskType["code"] = "groovyTask"

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
				"resultType":        resultType,
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

	groovyScriptTask := result.Task
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(convert.Int64ToString(groovyScriptTask.ID))

	return resourceTaskGroovyScriptRead(ctx, d, meta)
}

func resourceTaskGroovyScriptDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

type GroovyScript struct {
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

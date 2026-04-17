package task

import (
	"context"
	"log"
	"strings"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/convert"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceTaskEmail() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus email task resource",
		CreateContext: resourceTaskEmailCreate,
		ReadContext:   resourceTaskEmailRead,
		UpdateContext: resourceTaskEmailUpdate,
		DeleteContext: resourceTaskEmailDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the email task",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the email task",
				Required:    true,
			},
			"code": {
				Type:        schema.TypeString,
				Description: "The code of the email task",
				Optional:    true,
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "The organization labels associated with the task (Only supported on Morpheus 5.5.3 or higher)",
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"email_address": {
				Type: schema.TypeString,
				//nolint:lll
				Description: "Email addresses can be entered literally or Morpheus automation variables can be injected, such as <%=instance.createdByEmail%>",
				Required:    true,
			},
			"subject": {
				Type:        schema.TypeString,
				Description: "The subject line of the email, Morpheus automation variables can be injected into the subject field",
				Required:    true,
			},
			"source": {
				Type: schema.TypeString,
				//nolint:lll
				Description: "Choose local to draft or paste the email directly into the Task. Choose Repository or URL to bring in a template from a Git repository or another outside source (local, repository, url)",
				Optional:    true,
				Default:     "local",
			},
			"content_url": {
				Type:        schema.TypeString,
				Description: "The URL of the template used for the email task, used with a source type of url",
				Optional:    true,
				Computed:    true,
			},
			"content_path": {
				Type:        schema.TypeString,
				Description: "The file path of the template used for the email task, used with a source type of repository",
				Optional:    true,
				Computed:    true,
			},
			"repository_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the git repository to fetch the email template",
				Optional:    true,
				Computed:    true,
			},
			"version_ref": {
				Type:        schema.TypeString,
				Description: "The git reference of the repository to pull (main, master, etc.)",
				Optional:    true,
				Computed:    true,
			},
			"content": {
				Type: schema.TypeString,
				//nolint:lll
				Description: "The body of the email is HTML. Morpheus automation variables can be injected into the email body when needed. Used with a source type of local",
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
			"skip_wrapped_email_template": {
				Type:        schema.TypeBool,
				Description: "Whether to ignore the Morpheus-styled email template",
				Optional:    true,
				Default:     false,
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
				Description: "Custom configuration data to pass during the execution of the email task",
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

func resourceTaskEmailCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	contentConfig := make(map[string]any)

	var source string
	if sourceValue, ok := d.Get("source").(string); ok {
		source = sourceValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("source", d.Get("source")))
	}

	switch source {
	//nolint:goconst
	case "local":
		contentConfig["sourceType"] = "local"
		var content string
		if contentValue, ok := d.Get("content").(string); ok {
			content = contentValue
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("content", d.Get("content")))
		}
		contentConfig["content"] = content
	//nolint:goconst
	case "url":
		contentConfig["sourceType"] = "url"
		var contentUrl string
		if contentUrlValue, ok := d.Get("content_url").(string); ok {
			contentUrl = contentUrlValue
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("content_url", d.Get("content_url")))
		}
		contentConfig["contentPath"] = contentUrl
	//nolint:goconst
	case "repository":
		contentConfig["sourceType"] = "repository"
		repository := make(map[string]any)

		var repositoryId int
		if repositoryIdValue, ok := d.Get("repository_id").(int); ok {
			repositoryId = repositoryIdValue
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("repository_id", d.Get("repository_id")))
		}
		repository["id"] = repositoryId

		var contentPath string
		if contentPathValue, ok := d.Get("content_path").(string); ok {
			contentPath = contentPathValue
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("content_path", d.Get("content_path")))
		}
		contentConfig["contentPath"] = contentPath

		var versionRef string
		if versionRefValue, ok := d.Get("version_ref").(string); ok {
			versionRef = versionRefValue
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("version_ref", d.Get("version_ref")))
		}
		if versionRef != "" {
			contentConfig["contentRef"] = versionRef
		}
		contentConfig["repository"] = repository
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

	var emailAddress string
	if emailAddressValue, ok := d.Get("email_address").(string); ok {
		emailAddress = emailAddressValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("email_address", d.Get("email_address")))
	}

	var subject string
	if subjectValue, ok := d.Get("subject").(string); ok {
		subject = subjectValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("subject", d.Get("subject")))
	}

	var skipWrappedEmailTemplate bool
	if skipWrappedEmailTemplateValue, ok := d.Get("skip_wrapped_email_template").(bool); ok {
		skipWrappedEmailTemplate = skipWrappedEmailTemplateValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("skip_wrapped_email_template", d.Get("skip_wrapped_email_template")))
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
				"name":   name,
				"code":   code,
				"labels": labelsPayload,
				"taskType": map[string]any{
					"code": "email",
				},
				"taskOptions": map[string]any{
					"emailAddress":      emailAddress,
					"emailSubject":      subject,
					"emailSkipTemplate": skipWrappedEmailTemplate,
				},
				"file":              contentConfig,
				"executeTarget":     "local",
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

	diags = append(diags, resourceTaskEmailRead(ctx, d, meta)...)

	return diags
}

func resourceTaskEmailRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
	emailTask := result.Task
	d.SetId(convert.Int64ToString(emailTask.ID))
	d.Set("name", emailTask.Name)
	d.Set("code", emailTask.Code)
	d.Set("labels", emailTask.Labels)
	d.Set("email_address", emailTask.TaskOptions.EmailAddress)
	d.Set("subject", emailTask.TaskOptions.EmailSubject)
	d.Set("source", emailTask.File.SourceType)
	if emailTask.File.SourceType == "url" {
		d.Set("content_url", emailTask.File.ContentPath)
	}
	if emailTask.File.SourceType == "repository" {
		d.Set("content_path", emailTask.File.ContentPath)
		d.Set("repository_id", emailTask.File.Repository.ID)
		d.Set("version_ref", emailTask.File.ContentRef)
	}
	if emailTask.File.SourceType == "local" {
		d.Set("content", emailTask.File.Content)
	}
	if emailTask.TaskOptions.EmailSkipTemplate == "on" {
		d.Set("skip_wrapped_email_template", true)
	} else {
		d.Set("skip_wrapped_email_template", false)
	}
	d.Set("retryable", emailTask.Retryable)
	d.Set("retry_count", emailTask.RetryCount)
	d.Set("retry_delay_seconds", emailTask.RetryDelaySeconds)
	d.Set("allow_custom_config", emailTask.AllowCustomConfig)
	d.Set("visibility", emailTask.Visibility)

	return diags
}

func resourceTaskEmailUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	contentConfig := make(map[string]any)

	var source string
	if sourceValue, ok := d.Get("source").(string); ok {
		source = sourceValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("source", d.Get("source")))
	}

	switch source {
	case "local":
		contentConfig["sourceType"] = "local"
		var content string
		if contentValue, ok := d.Get("content").(string); ok {
			content = contentValue
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("content", d.Get("content")))
		}
		contentConfig["content"] = content
	case "url":
		contentConfig["sourceType"] = "url"
		var contentUrl string
		if contentUrlValue, ok := d.Get("content_url").(string); ok {
			contentUrl = contentUrlValue
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("content_url", d.Get("content_url")))
		}
		contentConfig["contentPath"] = contentUrl
	case "repository":
		contentConfig["sourceType"] = "repository"
		repository := make(map[string]any)

		var repositoryId int
		if repositoryIdValue, ok := d.Get("repository_id").(int); ok {
			repositoryId = repositoryIdValue
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("repository_id", d.Get("repository_id")))
		}
		repository["id"] = repositoryId

		var contentPath string
		if contentPathValue, ok := d.Get("content_path").(string); ok {
			contentPath = contentPathValue
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("content_path", d.Get("content_path")))
		}
		contentConfig["contentPath"] = contentPath

		if d.HasChange("version_ref") {
			var versionRef string
			if versionRefValue, ok := d.Get("version_ref").(string); ok {
				versionRef = versionRefValue
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("version_ref", d.Get("version_ref")))
			}
			contentConfig["contentRef"] = versionRef
		}
		contentConfig["repository"] = repository
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

	var emailAddress string
	if emailAddressValue, ok := d.Get("email_address").(string); ok {
		emailAddress = emailAddressValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("email_address", d.Get("email_address")))
	}

	var subject string
	if subjectValue, ok := d.Get("subject").(string); ok {
		subject = subjectValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("subject", d.Get("subject")))
	}

	var skipWrappedEmailTemplate bool
	if skipWrappedEmailTemplateValue, ok := d.Get("skip_wrapped_email_template").(bool); ok {
		skipWrappedEmailTemplate = skipWrappedEmailTemplateValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("skip_wrapped_email_template", d.Get("skip_wrapped_email_template")))
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
				"name":   name,
				"code":   code,
				"labels": labelsPayload,
				"taskType": map[string]any{
					"code": "email",
				},
				"taskOptions": map[string]any{
					"emailAddress":      emailAddress,
					"emailSubject":      subject,
					"emailSkipTemplate": skipWrappedEmailTemplate,
				},
				"file":              contentConfig,
				"executeTarget":     "local",
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
	emailTask := result.Task
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(convert.Int64ToString(emailTask.ID))

	return resourceTaskEmailRead(ctx, d, meta)
}

func resourceTaskEmailDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

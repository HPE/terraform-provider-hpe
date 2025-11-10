package task

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"strconv"
	"strings"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceTaskChefBootstrap() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus chef bootstrap task resource",
		CreateContext: resourceTaskChefBootstrapCreate,
		ReadContext:   resourceTaskChefBootstrapRead,
		UpdateContext: resourceTaskChefBootstrapUpdate,
		DeleteContext: resourceTaskChefBootstrapDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the chef bootstrap task",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the chef bootstrap task",
				Required:    true,
			},
			"code": {
				Type:        schema.TypeString,
				Description: "The code of the chef bootstrap task",
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
			"chef_server_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the Chef Server integration",
				Optional:    true,
			},
			"environment": {
				Type:        schema.TypeString,
				Description: "The chef environment",
				Optional:    true,
			},
			"run_list": {
				Type:        schema.TypeString,
				Description: "The chef run list",
				Optional:    true,
			},
			"data_bag_key": {
				Type:        schema.TypeString,
				Description: "The chef databag key",
				Sensitive:   true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					h := sha256.New()
					h.Write([]byte(new))
					sha256Hash := hex.EncodeToString(h.Sum(nil))

					return strings.EqualFold(old, sha256Hash)
				},
				Optional: true,
			},
			"data_bag_key_path": {
				Type:        schema.TypeString,
				Description: "The chef databag key path",
				Optional:    true,
			},
			"node_name": {
				Type:        schema.TypeString,
				Description: "The chef node name",
				Optional:    true,
			},
			"node_attributes": {
				Type:             schema.TypeString,
				Description:      "The chef node attributes (JSON)",
				Optional:         true,
				DiffSuppressFunc: helpers.SuppressEquivalentJSONDiffs,
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
				Description: "Custom configuration data to pass during the execution of the chef bootstrap",
				Optional:    true,
				Default:     false,
			},
			"visibility": {
				Type:         schema.TypeString,
				Description:  "Whether the task is visible in sub-tenants or not",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"private", "public"}, false),
				Default:      "private",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTaskChefBootstrapCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
	var chefServerId int
	if chefServerIdValue, ok := d.Get("chef_server_id").(int); ok {
		chefServerId = chefServerIdValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("chef_server_id", d.Get("chef_server_id")))
	}
	taskOptions["chefServerId"] = chefServerId

	var environment string
	if environmentValue, ok := d.Get("environment").(string); ok {
		environment = environmentValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("environment", d.Get("environment")))
	}
	taskOptions["chefEnv"] = environment

	var runList string
	if runListValue, ok := d.Get("run_list").(string); ok {
		runList = runListValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("run_list", d.Get("run_list")))
	}
	taskOptions["chefRunList"] = runList

	var dataBagKey string
	if dataBagKeyValue, ok := d.Get("data_bag_key").(string); ok {
		dataBagKey = dataBagKeyValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("data_bag_key", d.Get("data_bag_key")))
	}
	taskOptions["chefDataKey"] = dataBagKey

	var dataBagKeyPath string
	if dataBagKeyPathValue, ok := d.Get("data_bag_key_path").(string); ok {
		dataBagKeyPath = dataBagKeyPathValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("data_bag_key_path", d.Get("data_bag_key_path")))
	}
	taskOptions["chefDataKeyPath"] = dataBagKeyPath

	var nodeName string
	if nodeNameValue, ok := d.Get("node_name").(string); ok {
		nodeName = nodeNameValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("node_name", d.Get("node_name")))
	}
	taskOptions["chefNodeName"] = nodeName

	var nodeAttributes string
	if nodeAttributesValue, ok := d.Get("node_attributes").(string); ok {
		nodeAttributes = nodeAttributesValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("node_attributes", d.Get("node_attributes")))
	}
	taskOptions["chefAttributes"] = nodeAttributes

	taskType := make(map[string]any)
	taskType["code"] = "chefTask"

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
				"executeTarget":     "resource",
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

	diags = append(diags, resourceTaskChefBootstrapRead(ctx, d, meta)...)

	return diags
}

func resourceTaskChefBootstrapRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
	chefBootstrapTask := result.Task
	d.SetId(convert.Int64ToString(chefBootstrapTask.ID))
	d.Set("name", chefBootstrapTask.Name)
	d.Set("code", chefBootstrapTask.Code)
	d.Set("labels", chefBootstrapTask.Labels)
	serverId, _ := strconv.Atoi(chefBootstrapTask.TaskOptions.ChefServerId)
	d.Set("chef_server_id", serverId)
	d.Set("environment", chefBootstrapTask.TaskOptions.ChefEnv)
	d.Set("run_list", chefBootstrapTask.TaskOptions.ChefRunList)
	d.Set("data_bag_key", chefBootstrapTask.TaskOptions.ChefDataKeyHash)
	d.Set("data_bag_key_path", chefBootstrapTask.TaskOptions.ChefDataKeyPath)
	d.Set("node_name", chefBootstrapTask.TaskOptions.ChefNodeName)
	d.Set("node_attributes", chefBootstrapTask.TaskOptions.ChefAttributes)
	d.Set("retryable", chefBootstrapTask.Retryable)
	d.Set("retry_count", chefBootstrapTask.RetryCount)
	d.Set("retry_delay_seconds", chefBootstrapTask.RetryDelaySeconds)
	d.Set("allow_custom_config", chefBootstrapTask.AllowCustomConfig)
	d.Set("visibility", chefBootstrapTask.Visibility)

	return diags
}

func resourceTaskChefBootstrapUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	var chefServerId int
	if chefServerIdValue, ok := d.Get("chef_server_id").(int); ok {
		chefServerId = chefServerIdValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("chef_server_id", d.Get("chef_server_id")))
	}
	taskOptions["chefServerId"] = chefServerId

	var environment string
	if environmentValue, ok := d.Get("environment").(string); ok {
		environment = environmentValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("environment", d.Get("environment")))
	}
	taskOptions["chefEnv"] = environment

	var runList string
	if runListValue, ok := d.Get("run_list").(string); ok {
		runList = runListValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("run_list", d.Get("run_list")))
	}
	taskOptions["chefRunList"] = runList

	if d.HasChange("data_bag_key") {
		var dataBagKey string
		if dataBagKeyValue, ok := d.Get("data_bag_key").(string); ok {
			dataBagKey = dataBagKeyValue
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("data_bag_key", d.Get("data_bag_key")))
		}
		taskOptions["chefDataKey"] = dataBagKey
	}

	var dataBagKeyPath string
	if dataBagKeyPathValue, ok := d.Get("data_bag_key_path").(string); ok {
		dataBagKeyPath = dataBagKeyPathValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("data_bag_key_path", d.Get("data_bag_key_path")))
	}
	taskOptions["chefDataKeyPath"] = dataBagKeyPath

	var nodeName string
	if nodeNameValue, ok := d.Get("node_name").(string); ok {
		nodeName = nodeNameValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("node_name", d.Get("node_name")))
	}
	taskOptions["chefNodeName"] = nodeName

	var nodeAttributes string
	if nodeAttributesValue, ok := d.Get("node_attributes").(string); ok {
		nodeAttributes = nodeAttributesValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("node_attributes", d.Get("node_attributes")))
	}
	taskOptions["chefAttributes"] = nodeAttributes

	taskType := make(map[string]any)
	taskType["code"] = "chefTask"

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
				"executeTarget":     "resource",
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
	chefBootstrapTask := result.Task
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(convert.Int64ToString(chefBootstrapTask.ID))

	return resourceTaskChefBootstrapRead(ctx, d, meta)
}

func resourceTaskChefBootstrapDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

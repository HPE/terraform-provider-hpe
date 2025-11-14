// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package workflow

import (
	"context"
	"log"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"
)

func ResourceWorkflowProvisioning() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus provisioning workflow resource.",
		CreateContext: resourceWorkflowProvisioningCreate,
		ReadContext:   resourceWorkflowProvisioningRead,
		UpdateContext: resourceWorkflowProvisioningUpdate,
		DeleteContext: resourceWorkflowProvisioningDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the provisioning workflow",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the provisioning workflow",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the provisioning workflow",
				Optional:    true,
				Computed:    true,
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "The organization labels associated with the workflow (Only supported on Morpheus 5.5.3 or higher)",
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"platform": {
				Type: schema.TypeString,
				Description: "The operating system platforms the workflow is supported on " +
					"(all, linux, macos, windows)",
				ValidateFunc: validation.StringInSlice([]string{"all", "linux", "macos", "windows"}, false),
				Optional:     true,
				Computed:     true,
			},
			"visibility": {
				Type:         schema.TypeString,
				Description:  "Whether the provisioning workflow is visible in sub-tenants or not",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"private", "public"}, false),
				Default:      "private",
			},
			"task": {
				Type:        schema.TypeList,
				Description: "A list of tasks associated with the provisioning workflow",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"task_id": {
							Type:        schema.TypeInt,
							Description: "The ID of the task to associate with the provisioning workflow",
							Required:    true,
						},
						"task_phase": {
							Type: schema.TypeString,
							Description: "The phase that the task is executed " +
								"(configure, price, preProvision, provision, postProvision, start, stop, " +
								"preDeploy, deploy, reconfigure, teardown, shutdown, startup)",
							Required: true,
							ValidateFunc: validation.StringInSlice(
								[]string{
									"configure", "price", "preProvision", "provision", "postProvision",
									"start", "stop", "preDeploy", "deploy", "reconfigure", "teardown",
									"shutdown", "startup",
								},
								false,
							),
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceWorkflowProvisioningCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

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

	// tasks
	var tasks []map[string]any
	if d.Get("task") != nil {
		var taskList []any
		if v, ok := d.Get("task").([]any); ok {
			taskList = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("task", d.Get("task")))
		}

		if taskList != nil {
			// iterate over the array of tasks
			for i := 0; i < len(taskList); i++ {
				row := make(map[string]any)
				var taskconfig map[string]any
				if v, ok := taskList[i].(map[string]any); ok {
					taskconfig = v
				} else {
					return diag.FromErr(helpers.TypeAssertFailError("task element", taskList[i]))
				}
				row["taskId"] = taskconfig["task_id"]
				row["taskPhase"] = taskconfig["task_phase"]
				tasks = append(tasks, row)
			}
		}
	}

	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		var labelSet *schema.Set
		if v, ok := attr.(*schema.Set); ok {
			labelSet = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("labels", attr))
		}

		for _, s := range labelSet.List() {
			var labelStr string
			if v, ok := s.(string); ok {
				labelStr = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("label", s))
			}
			labelsPayload = append(labelsPayload, labelStr)
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

	req := &morpheus.Request{
		Body: map[string]any{
			"taskSet": map[string]any{
				"name":        name,
				"description": description,
				"labels":      labelsPayload,
				"type":        "provision",
				"visibility":  visibility,
				"platform":    platform,
				"tasks":       tasks,
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

	diags = append(diags, resourceWorkflowProvisioningRead(ctx, d, meta)...)

	return diags
}

func resourceWorkflowProvisioningRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

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

	// lookup by name if we do not have an id yet
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

	// Tasks
	var tasks []map[string]any
	var taskOrderList []TaskOrder
	if len(workflow.TaskSetTasks) != 0 {
		for _, task := range workflow.TaskSetTasks {
			var data TaskOrder
			data.ID = task.Task.ID
			data.Phase = task.TaskPhase
			data.Order = task.TaskOrder
			taskOrderList = append(taskOrderList, data)
		}
		sort.Slice(taskOrderList, func(i, j int) bool { return taskOrderList[i].Order < taskOrderList[j].Order })
		for _, task := range taskOrderList {
			tag := make(map[string]any)
			tag["task_phase"] = task.Phase
			tag["task_id"] = task.ID
			tasks = append(tasks, tag)
		}
	}

	d.SetId(convert.Int64ToString(workflow.ID))
	d.Set("name", workflow.Name)
	d.Set("description", workflow.Description)
	d.Set("labels", workflow.Labels)
	d.Set("visibility", workflow.Visibility)
	if workflow.Platform == "" {
		d.Set("platform", "all")
	} else {
		d.Set("platform", workflow.Platform)
	}
	d.Set("task", tasks)

	return diags
}

func resourceWorkflowProvisioningUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	// tasks
	var tasks []map[string]any
	if d.Get("task") != nil {
		var taskList []any
		if v, ok := d.Get("task").([]any); ok {
			taskList = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("task", d.Get("task")))
		}

		if taskList != nil {
			// iterate over the array of tasks
			for i := 0; i < len(taskList); i++ {
				row := make(map[string]any)
				var taskconfig map[string]any
				if v, ok := taskList[i].(map[string]any); ok {
					taskconfig = v
				} else {
					return diag.FromErr(helpers.TypeAssertFailError("task element", taskList[i]))
				}
				row["taskId"] = taskconfig["task_id"]
				row["taskPhase"] = taskconfig["task_phase"]
				tasks = append(tasks, row)
			}
		}
	}

	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		var labelSet *schema.Set
		if v, ok := attr.(*schema.Set); ok {
			labelSet = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("labels", attr))
		}

		for _, s := range labelSet.List() {
			var labelStr string
			if v, ok := s.(string); ok {
				labelStr = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("label", s))
			}
			labelsPayload = append(labelsPayload, labelStr)
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

	req := &morpheus.Request{
		Body: map[string]any{
			"taskSet": map[string]any{
				"name":        name,
				"description": description,
				"labels":      labelsPayload,
				"visibility":  visibility,
				"platform":    platform,
				"tasks":       tasks,
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

	workflow := result.TaskSet
	d.SetId(convert.Int64ToString(workflow.ID))

	return resourceWorkflowProvisioningRead(ctx, d, meta)
}

func resourceWorkflowProvisioningDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

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

type TaskOrder struct {
	Order int64  `json:"order"`
	ID    int64  `json:"id"`
	Phase string `json:"phase"`
}

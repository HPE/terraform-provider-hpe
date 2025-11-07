package job

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"
)

// ResourceJobTask returns the task job resource
func ResourceJobTask() *schema.Resource {
	return resourceJobTask()
}

func resourceJobTask() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a task job resource",
		CreateContext: resourceJobTaskCreate,
		ReadContext:   resourceJobTaskRead,
		UpdateContext: resourceJobTaskUpdate,
		DeleteContext: resourceJobTaskDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the task job",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the task job",
				Required:    true,
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "The organization labels associated with the task job (Only supported on Morpheus 5.5.3 or higher)",
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the task job is enabled",
				Optional:    true,
				Default:     true,
			},
			"task_id": {
				Type:        schema.TypeInt,
				Description: "The id of the task associated with the job",
				Required:    true,
			},
			"schedule_mode": {
				Type:         schema.TypeString,
				Description:  "The job scheduling type (manual, date_and_time, scheduled)",
				ValidateFunc: validation.StringInSlice([]string{"manual", "date_and_time", "scheduled"}, false),
				Required:     true,
			},
			"scheduled_date_and_time": {
				Type:          schema.TypeString,
				Description:   "The date and time the job will be executed if schedule mode date_and_time is used",
				Optional:      true,
				ConflictsWith: []string{"execution_schedule_id"},
			},
			"execution_schedule_id": {
				Type:        schema.TypeInt,
				Description: "The id of the execution schedule associated with the job",
				Optional:    true,
			},
			"context_type": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"appliance", "server", "instance", "instance-label", "server-label"}, false),
				Description:  "The context that the job should run as (appliance, server, instance, instance-label, server-label)",
				Required:     true,
			},
			"server_ids": {
				Type:          schema.TypeList,
				Description:   "A list of server ids to associate with the job",
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeInt},
				ConflictsWith: []string{"instance_ids", "instance_label", "server_label"},
			},
			"server_label": {
				Type:          schema.TypeString,
				Description:   "The server label used for dynamic automation targeting",
				Optional:      true,
				ConflictsWith: []string{"instance_ids", "server_ids", "instance_label"},
			},
			"instance_ids": {
				Type:          schema.TypeList,
				Description:   "A list of instance ids to associate with the job",
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeInt},
				ConflictsWith: []string{"server_ids", "instance_label", "server_label"},
			},
			"instance_label": {
				Type:          schema.TypeString,
				Description:   "The instance label used for dynamic automation targeting",
				Optional:      true,
				ConflictsWith: []string{"instance_ids", "server_ids", "server_label"},
			},
			"custom_config": {
				Type:        schema.TypeString,
				Description: "The task custom configuration",
				Optional:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceJobTaskCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	job := make(map[string]any)

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}
	job["name"] = name
	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		var labelsSet *schema.Set
		if v, ok := attr.(*schema.Set); ok {
			labelsSet = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("labels", attr))
		}
		for _, s := range labelsSet.List() {
			var labelStr string
			if v, ok := s.(string); ok {
				labelStr = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("label", s))
			}
			labelsPayload = append(labelsPayload, labelStr)
		}
	}
	job["labels"] = labelsPayload
	var enabled bool
	if v, ok := d.Get("enabled").(bool); ok {
		enabled = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("enabled", d.Get("enabled")))
	}
	job["enabled"] = enabled
	var taskID int
	if v, ok := d.Get("task_id").(int); ok {
		taskID = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("task_id", d.Get("task_id")))
	}
	job["task"] = map[string]int{
		"id": taskID,
	}
	var contextType string
	if v, ok := d.Get("context_type").(string); ok {
		contextType = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("context_type", d.Get("context_type")))
	}
	job["targetType"] = contextType
	if contextType == "instance-label" {
		var instanceLabel string
		if v, ok := d.Get("instance_label").(string); ok {
			instanceLabel = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("instance_label", d.Get("instance_label")))
		}
		job["instanceLabel"] = instanceLabel
	}
	if contextType == "server-label" {
		var serverLabel string
		if v, ok := d.Get("server_label").(string); ok {
			serverLabel = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("server_label", d.Get("server_label")))
		}
		job["serverLabel"] = serverLabel
	}
	job["customConfig"] = d.Get("custom_config")

	// Evaluate different schedululing methods
	var scheduleMode string
	if v, ok := d.Get("schedule_mode").(string); ok {
		scheduleMode = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("schedule_mode", d.Get("schedule_mode")))
	}
	switch scheduleMode {
	case "manual":
		job["scheduleMode"] = "manual"
	case "date_and_time":
		job["scheduleMode"] = "dateTime"
		var scheduledDateTime string
		if v, ok := d.Get("scheduled_date_and_time").(string); ok {
			scheduledDateTime = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("scheduled_date_and_time", d.Get("scheduled_date_and_time")))
		}
		job["dateTime"] = scheduledDateTime
	case "scheduled":
		job["scheduleMode"] = d.Get("execution_schedule_id")
	}

	// instances
	var targets []map[string]any
	if contextType == "instance" {
		var instanceList []any
		if v, ok := d.Get("instance_ids").([]any); ok {
			instanceList = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("instance_ids", d.Get("instance_ids")))
		}
		// iterate over the array of instance ids
		for i := 0; i < len(instanceList); i++ {
			row := make(map[string]any)
			row["refId"] = instanceList[i]
			targets = append(targets, row)
		}
	}

	// servers
	if contextType == "server" {
		var serverList []any
		if v, ok := d.Get("server_ids").([]any); ok {
			serverList = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("server_ids", d.Get("server_ids")))
		}
		// iterate over the array of server ids
		for i := 0; i < len(serverList); i++ {
			row := make(map[string]any)
			row["refId"] = serverList[i]
			targets = append(targets, row)
		}
	}

	job["targets"] = targets

	req := &morpheus.Request{
		Body: map[string]any{
			"job": job,
		},
	}
	resp, err := client.CreateJob(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	var result map[string]any
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		log.Fatal(err)
	}

	jobId := fmt.Sprintf("%v", result["id"])
	// Successfully created resource, now set id
	d.SetId(jobId)

	resourceJobTaskRead(ctx, d, meta)

	return diags
}

func resourceJobTaskRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

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
		resp, err = client.FindJobByName(name)
	} else if id != "" {
		resp, err = client.GetJob(convert.StringToInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Execute schedule cannot be read without name or id")
	}

	if err != nil {
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
	var result *morpheus.GetJobResult
	if v, ok := resp.Result.(*morpheus.GetJobResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("result", resp.Result))
	}
	if result.Job == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Job"))
	}
	taskJob := result.Job

	d.SetId(convert.Int64ToString(taskJob.ID))
	d.Set("name", taskJob.Name)
	if taskJob.Labels != nil && len(taskJob.Labels) > 0 {
		d.Set("labels", taskJob.Labels)
	}
	d.Set("enabled", taskJob.Enabled)
	d.Set("task_id", taskJob.Task.ID)
	d.Set("context_type", taskJob.TargetType)
	switch taskJob.ScheduleMode {
	case "manual":
		d.Set("schedule_mode", "manual")
	case "dateTime":
		d.Set("schedule_mode", "date_and_time")
		d.Set("scheduled_date_and_time", taskJob.DateTime)
		// Execute schedule uses the schedule mode field for storing the execute schedule id
	default:
		d.Set("schedule_mode", "scheduled")
		intVar, err := strconv.Atoi(taskJob.ScheduleMode)
		if err != nil {
			log.Printf("String Conversion Failure: %s", err)
		}
		d.Set("execution_schedule_id", intVar)
	}
	if taskJob.CustomConfig != "" {
		d.Set("custom_config", taskJob.CustomConfig)
	}

	switch taskJob.TargetType {
	case "instance":
		// instances
		var instanceIds []int64
		if taskJob.Targets != nil {
			// iterate over the array of targets
			for i := 0; i < len(taskJob.Targets); i++ {
				instance := taskJob.Targets[i]
				instanceIds = append(instanceIds, int64(instance.RefId))
			}
		}
		d.Set("instance_ids", instanceIds)
	case "server":
		// servers
		var serverIds []int64
		if taskJob.Targets != nil {
			// iterate over the array of targets
			for i := 0; i < len(taskJob.Targets); i++ {
				server := taskJob.Targets[i]
				serverIds = append(serverIds, int64(server.RefId))
			}
		}
		d.Set("server_ids", serverIds)
	case "instance-label":
		if taskJob.Targets != nil && len(taskJob.Targets) > 0 {
			d.Set("instance_label", taskJob.Targets[0].Name)
		}
	case "server-label":
		if taskJob.Targets != nil && len(taskJob.Targets) > 0 {
			d.Set("server_label", taskJob.Targets[0].Name)
		}
	}

	return diags
}

func resourceJobTaskUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}
	id := d.Id()

	job := make(map[string]any)

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}
	job["name"] = name
	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		var labelsSet *schema.Set
		if v, ok := attr.(*schema.Set); ok {
			labelsSet = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("labels", attr))
		}
		for _, s := range labelsSet.List() {
			var labelStr string
			if v, ok := s.(string); ok {
				labelStr = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("label", s))
			}
			labelsPayload = append(labelsPayload, labelStr)
		}
	}
	job["labels"] = labelsPayload
	var enabled bool
	if v, ok := d.Get("enabled").(bool); ok {
		enabled = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("enabled", d.Get("enabled")))
	}
	job["enabled"] = enabled
	var taskID int
	if v, ok := d.Get("task_id").(int); ok {
		taskID = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("task_id", d.Get("task_id")))
	}
	job["task"] = map[string]int{
		"id": taskID,
	}
	var contextType string
	if v, ok := d.Get("context_type").(string); ok {
		contextType = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("context_type", d.Get("context_type")))
	}
	job["targetType"] = contextType
	if contextType == "instance-label" {
		var instanceLabel string
		if v, ok := d.Get("instance_label").(string); ok {
			instanceLabel = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("instance_label", d.Get("instance_label")))
		}
		job["instanceLabel"] = instanceLabel
	}
	if contextType == "server-label" {
		var serverLabel string
		if v, ok := d.Get("server_label").(string); ok {
			serverLabel = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("server_label", d.Get("server_label")))
		}
		job["serverLabel"] = serverLabel
	}
	job["customConfig"] = d.Get("custom_config")

	// Evaluate different schedululing methods
	var scheduleMode string
	if v, ok := d.Get("schedule_mode").(string); ok {
		scheduleMode = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("schedule_mode", d.Get("schedule_mode")))
	}
	switch scheduleMode {
	case "manual":
		job["scheduleMode"] = "manual"
	case "date_and_time":
		job["scheduleMode"] = "dateTime"
		var scheduledDateTime string
		if v, ok := d.Get("scheduled_date_and_time").(string); ok {
			scheduledDateTime = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("scheduled_date_and_time", d.Get("scheduled_date_and_time")))
		}
		job["dateTime"] = scheduledDateTime
	case "scheduled":
		job["scheduleMode"] = d.Get("execution_schedule_id")
	}

	// instances
	var targets []map[string]any
	if contextType == "instance" {
		var instanceList []any
		if v, ok := d.Get("instance_ids").([]any); ok {
			instanceList = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("instance_ids", d.Get("instance_ids")))
		}
		// iterate over the array of instance ids
		for i := 0; i < len(instanceList); i++ {
			row := make(map[string]any)
			row["refId"] = instanceList[i]
			targets = append(targets, row)
		}
	}

	// servers
	if contextType == "server" {
		var serverList []any
		if v, ok := d.Get("server_ids").([]any); ok {
			serverList = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("server_ids", d.Get("server_ids")))
		}
		// iterate over the array of instance ids
		for i := 0; i < len(serverList); i++ {
			row := make(map[string]any)
			row["refId"] = serverList[i]
			targets = append(targets, row)
		}
	}

	job["targets"] = targets

	req := &morpheus.Request{
		Body: map[string]any{
			"job": job,
		},
	}

	resp, err := client.UpdateJob(convert.StringToInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	var result map[string]any
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		log.Fatal(err)
	}

	// Successfully updated resource, now set id
	// the task API doesn't return the id so setting
	// to the original id
	d.SetId(id)

	return resourceJobTaskRead(ctx, d, meta)
}

func resourceJobTaskDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteJob(convert.StringToInt64(id), req)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)

			return diag.FromErr(err)
		} else {
			log.Printf("API FAILURE: %s - %s", resp, err)

			return diag.FromErr(err)
		}
	}
	log.Printf("API RESPONSE: %s", resp)
	d.SetId("")

	return diags
}

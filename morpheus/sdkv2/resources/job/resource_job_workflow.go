package job

import (
	"context"
	"log"
	"strconv"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/convert"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"
)

// ResourceJobWorkflow returns the workflow job resource
//
//nolint:lll
func ResourceJobWorkflow() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a workflow job resource",
		CreateContext: resourceJobWorkflowCreate,
		ReadContext:   resourceJobWorkflowRead,
		UpdateContext: resourceJobWorkflowUpdate,
		DeleteContext: resourceJobWorkflowDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the workflow job",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the workflow job",
				Required:    true,
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "The organization labels associated with the workflow job (Only supported on Morpheus 5.5.3 or higher)",
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the workflow job is enabled",
				Optional:    true,
				Default:     true,
			},
			"workflow_id": {
				Type:        schema.TypeInt,
				Description: "The id of the workflow associated with the job",
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
			"custom_options": {
				Description: "Custom options to pass to the workflow",
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

//nolint:goconst
func resourceJobWorkflowCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*morpheus.Client)

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
		for _, s := range attr.(*schema.Set).List() {
			var label string
			if v, ok := s.(string); ok {
				label = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("labels element", s))
			}
			labelsPayload = append(labelsPayload, label)
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

	var workflowID int
	if v, ok := d.Get("workflow_id").(int); ok {
		workflowID = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("workflow_id", d.Get("workflow_id")))
	}
	job["workflow"] = map[string]int{
		"id": workflowID,
	}

	// Evaluate different schedululing methods
	switch d.Get("schedule_mode") {
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

	// instances
	var targets []map[string]any
	if d.Get("context_type") == "instance" {
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
	if d.Get("context_type") == "server" {
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

	// Custom Options
	if d.Get("custom_options") != nil {
		var customOptionsInput map[string]any
		if v, ok := d.Get("custom_options").(map[string]any); ok {
			customOptionsInput = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("custom_options", d.Get("custom_options")))
		}
		customOptions := make(map[string]any)
		for key, value := range customOptionsInput {
			var strValue string
			if v, ok := value.(string); ok {
				strValue = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("custom_options value", value))
			}
			customOptions[key] = strValue
		}
		job["customOptions"] = customOptions
	}

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

	result := resp.Result.(*morpheus.CreateJobResult)
	jobResult := result.Job
	if jobResult == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Job"))
	}

	// Successfully created resource, now set id
	d.SetId(convert.Int64ToString(jobResult.ID))

	diags = append(diags, resourceJobWorkflowRead(ctx, d, meta)...)

	return diags
}

//nolint:goconst
func resourceJobWorkflowRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*morpheus.Client)
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
	result := resp.Result.(*morpheus.GetJobResult)
	workflowJob := result.Job
	if workflowJob == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Job"))
	}

	d.SetId(convert.Int64ToString(workflowJob.ID))
	d.Set("name", workflowJob.Name)
	if len(workflowJob.Labels) > 0 {
		d.Set("labels", workflowJob.Labels)
	}
	d.Set("enabled", workflowJob.Enabled)
	d.Set("workflow_id", workflowJob.Workflow.ID)
	d.Set("context_type", workflowJob.TargetType)
	switch workflowJob.ScheduleMode {
	case "manual":
		d.Set("schedule_mode", "manual")
	case "dateTime":
		d.Set("schedule_mode", "date_and_time")
		d.Set("scheduled_date_and_time", workflowJob.DateTime)
		// Execute schedule uses the schedule mode field for storing the execute schedule id
	default:
		d.Set("schedule_mode", "scheduled")
		intVar, err := strconv.Atoi(workflowJob.ScheduleMode)
		if err != nil {
			log.Printf("String Conversion Failure: %s", err)
		}
		d.Set("execution_schedule_id", intVar)
	}

	switch workflowJob.TargetType {
	case "instance":
		// instances
		var instanceIDs []int64
		if workflowJob.Targets != nil {
			// iterate over the array of targets
			for i := 0; i < len(workflowJob.Targets); i++ {
				instance := workflowJob.Targets[i]
				instanceIDs = append(instanceIDs, instance.RefId)
			}
		}
		d.Set("instance_ids", instanceIDs)
	case "server":
		// servers
		var serverIDs []int64
		if workflowJob.Targets != nil {
			// iterate over the array of targets
			for i := 0; i < len(workflowJob.Targets); i++ {
				server := workflowJob.Targets[i]
				serverIDs = append(serverIDs, server.RefId)
			}
		}
		d.Set("server_ids", serverIDs)
	case "instance-label":
		if len(workflowJob.Targets) > 0 {
			d.Set("instance_label", workflowJob.Targets[0].Name)
		}
	case "server-label":
		if len(workflowJob.Targets) > 0 {
			d.Set("server_label", workflowJob.Targets[0].Name)
		}
	}
	d.Set("custom_options", workflowJob.CustomOptions)

	return diags
}

func resourceJobWorkflowUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*morpheus.Client)
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
		for _, s := range attr.(*schema.Set).List() {
			var label string
			if v, ok := s.(string); ok {
				label = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("labels element", s))
			}
			labelsPayload = append(labelsPayload, label)
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

	var workflowID int
	if v, ok := d.Get("workflow_id").(int); ok {
		workflowID = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("workflow_id", d.Get("workflow_id")))
	}
	job["workflow"] = map[string]int{
		"id": workflowID,
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

	// Evaluate different schedululing methods
	switch d.Get("schedule_mode") {
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
	if d.Get("context_type") == "instance" {
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
	if d.Get("context_type") == "server" {
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

	// Custom Options
	if d.Get("custom_options") != nil {
		var customOptionsInput map[string]any
		if v, ok := d.Get("custom_options").(map[string]any); ok {
			customOptionsInput = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("custom_options", d.Get("custom_options")))
		}
		customOptions := make(map[string]any)
		for key, value := range customOptionsInput {
			var strValue string
			if v, ok := value.(string); ok {
				strValue = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("custom_options value", value))
			}
			customOptions[key] = strValue
		}
		job["customOptions"] = customOptions
	}

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
	result := resp.Result.(*morpheus.UpdateJobResult)
	jobResult := result.Job
	if jobResult == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Job"))
	}

	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(convert.Int64ToString(jobResult.ID))

	return resourceJobWorkflowRead(ctx, d, meta)
}

func resourceJobWorkflowDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*morpheus.Client)

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

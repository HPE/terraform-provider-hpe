// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package resources

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/HPE/terraform-provider-opsramp/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure implementation satisfies the expected interfaces
var _ resource.Resource = &ScheduleMaintenanceResource{}
var _ resource.ResourceWithImportState = &ScheduleMaintenanceResource{}
var _ resource.ResourceWithModifyPlan = &ScheduleMaintenanceResource{}

// ScheduleMaintenanceResource defines the resource implementation.
type ScheduleMaintenanceResource struct {
	BaseResource
}

// ScheduleMaintenanceModel maps Terraform schema attributes to the provider model.
type ScheduleMaintenanceModel struct {
	Id                types.String `tfsdk:"id"`
	Client            types.String `tfsdk:"client"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	RunRBA            types.Bool   `tfsdk:"run_rba"`
	InstallPatch      types.Bool   `tfsdk:"install_patch"`
	CorrelateAlerts   types.Bool   `tfsdk:"correlate_alerts"`
	RunEscalateAction types.Bool   `tfsdk:"run_escalate_action"`
	Status            types.String `tfsdk:"status"`

	// Schedule timing
	Schedule *ScheduleMaintenanceScheduleModel `tfsdk:"schedule"`

	// Resources
	DeviceIds      types.Set `tfsdk:"device_ids"`
	DeviceGroupIds types.Set `tfsdk:"device_group_ids"`
	SiteIds        types.Set `tfsdk:"site_ids"`

	// Alert conditions
	AlertConditions *ScheduleMaintenanceAlertConditionsModel `tfsdk:"alert_conditions"`

	UserIds               types.Set `tfsdk:"user_ids"`
	UserGroupIds          types.Set `tfsdk:"user_group_ids"`
	NotifyBeforeStartTime int64     `tfsdk:"notify_before_start_time"`
	NotifyBeforeEndTime   int64     `tfsdk:"notify_before_end_time"`
}

// ScheduleMaintenanceScheduleModel represents the schedule timing block
type ScheduleMaintenanceScheduleModel struct {
	Type      types.String                     `tfsdk:"type"`
	StartTime types.String                     `tfsdk:"start_time"`
	EndTime   types.String                     `tfsdk:"end_time"`
	EndBy     types.String                     `tfsdk:"end_by"`
	Timezone  types.String                     `tfsdk:"timezone"`
	Pattern   *ScheduleMaintenancePatternModel `tfsdk:"pattern"`
}

// ScheduleMaintenancePatternModel represents the recurrence pattern
type ScheduleMaintenancePatternModel struct {
	Type            types.String `tfsdk:"type"`
	WeekDays        types.String `tfsdk:"week_days"`
	DayOfWeek       types.String `tfsdk:"day_of_week"`
	WeekIndex       types.String `tfsdk:"week_index"`
	RepeatFrecuency types.Int64  `tfsdk:"repeat_frecuency"`
	DayFrequency    types.String `tfsdk:"day_frequency"`
	Months          types.String `tfsdk:"months"`
	DayOfMonth      types.String `tfsdk:"day_of_month"`
}

// ScheduleMaintenanceAlertConditionsModel represents alert conditions
type ScheduleMaintenanceAlertConditionsModel struct {
	MatchingType types.String                        `tfsdk:"matching_type"`
	Rules        []ScheduleMaintenanceAlertRuleModel `tfsdk:"rules"`
}

// ScheduleMaintenanceAlertRuleModel represents a single alert condition rule
type ScheduleMaintenanceAlertRuleModel struct {
	Key      types.String `tfsdk:"key"`
	Operator types.String `tfsdk:"operator"`
	Value    types.String `tfsdk:"value"`
}

// NewScheduleMaintenance creates a new instance of the resource.
func NewScheduleMaintenance() resource.Resource {
	return &ScheduleMaintenanceResource{}
}

// Metadata returns the resource type name.
func (r *ScheduleMaintenanceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_schedule_maintenance"
}

// Schema defines the schema for the resource.
func (r *ScheduleMaintenanceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an OpsRamp Scheduled Maintenance Window. Supports one-time and recurring schedules with device, device group, and location targeting.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the maintenance window (e.g. SM-...).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"client": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The unique ID of the client (sub-tenant). If not provided, uses the provider's tenant.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the scheduled maintenance window.",
			},
			"description": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "A description of the maintenance window.",
			},
			"run_rba": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Enable run book automation during maintenance.",
			},
			"install_patch": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Install patches during maintenance.",
			},
			"correlate_alerts": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Enable alert correlation during maintenance.",
			},
			"run_escalate_action": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Enable escalation policy during maintenance.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The current status of the maintenance window (Active, Pending, Completed, Suspended).",
			},
			"schedule": schema.SingleNestedAttribute{
				Required:            true,
				MarkdownDescription: "The schedule timing configuration.",
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The schedule type (`one-time` or `recurring`).",
						Validators: []validator.String{
							stringvalidator.OneOf("one-time", "recurring"),
						},
					},
					"start_time": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The start time in ISO 8601 format (e.g. `2026-06-10T17:00:00+0200`).",
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}[+-]\d{4}$`),
								"must be in ISO 8601 format: YYYY-MM-DDTHH:MM:SS+HHMM (e.g. 2026-06-10T17:00:00+0200)",
							),
						},
					},
					"end_time": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The end time in ISO 8601 format (e.g. `2026-06-10T18:00:00+0200`).",
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}[+-]\d{4}$`),
								"must be in ISO 8601 format: YYYY-MM-DDTHH:MM:SS+HHMM (e.g. 2026-06-10T18:00:00+0200)",
							),
						},
					},
					"end_by": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "When the recurring schedule ends. Use `Never` for no end date, or a datetime in ISO 8601 format (e.g. `2026-12-31T23:59:59+0000`).",
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								regexp.MustCompile(`^(Never|\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}[+-]\d{4})$`),
								"must be \"Never\" or a datetime in ISO 8601 format: YYYY-MM-DDTHH:MM:SS+HHMM",
							),
						},
					},
					"timezone": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The timezone for the schedule (e.g. `America/New_York`, `GMT`, `UTC`).",
					},
					"pattern": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "The recurrence pattern. Required when schedule type is `Recurring`.",
						Attributes: map[string]schema.Attribute{
							"type": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "The pattern type (`daily`, `weekly`, `monthly` or `yearly`).",
								Validators: []validator.String{
									stringvalidator.OneOf("daily", "weekly", "monthly", "yearly"),
								},
							},
							"repeat_frecuency": schema.Int64Attribute{
								Optional:            true,
								MarkdownDescription: "The repeat frequency for the pattern. For daily, how many days between occurrences; for weekly, how many weeks between occurrences; for monthly, how many months between occurrences.",
							},
							"day_frequency": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "For monthly pattern, the day of the month (e.g. `15`) or a descriptor like `everyday`.",
							},
							"week_days": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "Comma-separated days for weekly pattern (e.g. `Monday,Wednesday,Friday`).",
							},
							"months": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "For yearly pattern, comma-separated months (e.g. `January,June`).",
							},
							"day_of_month": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "For monthly pattern, the day of the month (e.g. `15`) or a descriptor like `LastDay`.",
							},
							"week_index": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "For monthly pattern, the week index (e.g. `First`, `Second`, `Third`, `Fourth`, `Last`).",
							},
							"day_of_week": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "For monthly pattern, the day of the week (e.g. `Monday`).",
							},
						},
					},
				},
			},
			"device_ids": schema.SetAttribute{
				Optional:            true,
				MarkdownDescription: "Device IDs to include in the maintenance window.",
				ElementType:         types.StringType,
			},
			"device_group_ids": schema.SetAttribute{
				Optional:            true,
				MarkdownDescription: "Device group IDs to include in the maintenance window.",
				ElementType:         types.StringType,
			},
			"site_ids": schema.SetAttribute{
				Optional:            true,
				MarkdownDescription: "Site IDs to include in the maintenance window.",
				ElementType:         types.StringType,
			},
			"alert_conditions": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Alert conditions to filter which alerts are suppressed during maintenance.",
				Attributes: map[string]schema.Attribute{
					"matching_type": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "How rules are combined (`ANY` or `ALL`).",
						Validators: []validator.String{
							stringvalidator.OneOf("ANY", "ALL"),
						},
					},
					"rules": schema.ListNestedAttribute{
						Required:            true,
						MarkdownDescription: "Alert condition rules.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"key": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "The alert field to match (e.g. `subject`, `description`, `serviceName`, `resourceName`).",
								},
								"operator": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "The comparison operator (e.g. `CONTAINS`, `EQUALS`).",
								},
								"value": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "The value to match against.",
								},
							},
						},
					},
				},
			},
			"user_ids": schema.SetAttribute{
				Optional:            true,
				MarkdownDescription: "User IDs to notify about the maintenance window.",
				ElementType:         types.StringType,
			},
			"user_group_ids": schema.SetAttribute{
				Optional:            true,
				MarkdownDescription: "User group IDs to notify about the maintenance window.",
				ElementType:         types.StringType,
			},
			"notify_before_start_time": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Number of minutes to notify users before maintenance starts.",
				Default:             int64default.StaticInt64(0),
			},
			"notify_before_end_time": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Number of minutes to notify users before maintenance ends.",
				Default:             int64default.StaticInt64(0),
			},
		},
	}
}

// getTenantId determines which tenant ID to use based on the optional client parameter
func (r *ScheduleMaintenanceResource) getTenantId(clientId types.String) string {
	if !clientId.IsNull() && clientId.ValueString() != "" {
		return clientId.ValueString()
	}
	return r.apiClient.TenantId
}

// mapScheduleMaintenanceResponse maps an API response onto a model in place.
// It is used by Create, Read, and Update to keep state consistent.
func mapScheduleMaintenanceResponse(existing *client.ScheduleMaintenanceResponse, model *ScheduleMaintenanceModel) {
	model.Name = types.StringValue(existing.Name)
	model.Description = types.StringValue(existing.Description)
	model.RunRBA = types.BoolValue(existing.RunRBA)
	model.InstallPatch = types.BoolValue(existing.InstallPatch)
	model.CorrelateAlerts = types.BoolValue(existing.CorrelateAlerts)
	model.RunEscalateAction = types.BoolValue(existing.RunEscalateAction)
	model.Status = types.StringValue(existing.Status)

	// Fully rebuild Schedule from API response (required for import and drift detection)
	sched := ScheduleMaintenanceScheduleModel{
		Type:      types.StringValue(existing.Schedule.Type),
		StartTime: types.StringValue(existing.Schedule.StartTime),
		EndTime:   types.StringValue(existing.Schedule.EndTime),
		Timezone:  types.StringValue(existing.Schedule.Timezone),
		EndBy:     types.StringNull(),
	}
	if existing.Schedule.EndBy != "" {
		sched.EndBy = types.StringValue(existing.Schedule.EndBy)
	}
	if existing.Schedule.Pattern != nil {
		pat := &ScheduleMaintenancePatternModel{
			Type:            types.StringValue(existing.Schedule.Pattern.Type),
			WeekDays:        types.StringNull(),
			DayOfWeek:       types.StringNull(),
			WeekIndex:       types.StringNull(),
			RepeatFrecuency: types.Int64Null(),
			DayFrequency:    types.StringNull(),
			Months:          types.StringNull(),
			DayOfMonth:      types.StringNull(),
		}
		if existing.Schedule.Pattern.WeekDays != "" {
			pat.WeekDays = types.StringValue(existing.Schedule.Pattern.WeekDays)
		}
		if existing.Schedule.Pattern.DayOfWeek != "" {
			pat.DayOfWeek = types.StringValue(existing.Schedule.Pattern.DayOfWeek)
		}
		if existing.Schedule.Pattern.WeekIndex != "" {
			pat.WeekIndex = types.StringValue(existing.Schedule.Pattern.WeekIndex)
		}
		if existing.Schedule.Pattern.RepeatFrequency != 0 {
			pat.RepeatFrecuency = types.Int64Value(int64(existing.Schedule.Pattern.RepeatFrequency))
		}
		if existing.Schedule.Pattern.DayFrequency != "" {
			pat.DayFrequency = types.StringValue(existing.Schedule.Pattern.DayFrequency)
		}
		if existing.Schedule.Pattern.Months != "" {
			pat.Months = types.StringValue(existing.Schedule.Pattern.Months)
		}
		if existing.Schedule.Pattern.DayOfMonth != "" {
			pat.DayOfMonth = types.StringValue(existing.Schedule.Pattern.DayOfMonth)
		}
		sched.Pattern = pat
	} else if model.Schedule != nil {
		// Preserve pattern from prior state so plan doesn't show a spurious diff
		sched.Pattern = model.Schedule.Pattern
	}
	model.Schedule = &sched

	// Rebuild alert conditions when present in the response
	if existing.AlertConditions != nil {
		ac := &ScheduleMaintenanceAlertConditionsModel{
			MatchingType: types.StringValue(existing.AlertConditions.MatchingType),
		}
		for _, rule := range existing.AlertConditions.Rules {
			mappedKey := rule.Key
			switch rule.Key {
			case "Alert : Subject":
				mappedKey = "subject"
			case "Alert : Description":
				mappedKey = "description"
			case "Alert : Metric":
				mappedKey = "serviceName"
			case "Resource : Resource Name":
				mappedKey = "resourceName"
			}

			ac.Rules = append(ac.Rules, ScheduleMaintenanceAlertRuleModel{
				Key:      types.StringValue(mappedKey),
				Operator: types.StringValue(rule.Operator),
				Value:    types.StringValue(rule.Value),
			})
		}
		model.AlertConditions = ac
	} else {
		model.AlertConditions = nil
	}
}

// Create handles the creation of the resource.
func (r *ScheduleMaintenanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ScheduleMaintenanceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.getTenantId(plan.Client)

	apiReq := r.buildRequest(plan)

	created, err := r.apiClient.CreateScheduleMaintenance(tenantId, apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Creation Error",
			fmt.Sprintf("Could not create schedule maintenance: %s", err.Error()))
		return
	}

	plan.Id = types.StringValue(created.UniqueId)

	// Fetch the freshly created resource to populate all computed fields (including Status).
	existing, err := r.apiClient.GetScheduleMaintenance(tenantId, created.UniqueId)
	if err != nil {
		// Non-fatal: store what we know and let the next Read reconcile.
		plan.Status = types.StringValue("Pending")
	} else {
		mapScheduleMaintenanceResponse(existing, &plan)
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Read handles reading the resource state.
func (r *ScheduleMaintenanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ScheduleMaintenanceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.getTenantId(state.Client)

	existing, err := r.apiClient.GetScheduleMaintenance(tenantId, state.Id.ValueString())
	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "404") || strings.Contains(errStr, "not found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read Error", err.Error())
		return
	}

	// Update scalar fields from response
	mapScheduleMaintenanceResponse(existing, &state)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update handles updating the resource.
func (r *ScheduleMaintenanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ScheduleMaintenanceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state ScheduleMaintenanceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.getTenantId(plan.Client)

	apiReq := r.buildRequest(plan)

	_, err := r.apiClient.UpdateScheduleMaintenance(tenantId, state.Id.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Update Error",
			fmt.Sprintf("Could not update schedule maintenance: %s", err.Error()))
		return
	}

	plan.Id = state.Id

	// Fetch the updated resource to populate all computed fields (including Status).
	existing, err := r.apiClient.GetScheduleMaintenance(tenantId, state.Id.ValueString())
	if err != nil {
		// Non-fatal: carry Status over from prior state and let the next Read reconcile.
		plan.Status = state.Status
	} else {
		mapScheduleMaintenanceResponse(existing, &plan)
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Delete handles deleting the resource.
func (r *ScheduleMaintenanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ScheduleMaintenanceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.getTenantId(state.Client)

	err := r.apiClient.DeleteScheduleMaintenance(tenantId, state.Id.ValueString())
	if err != nil {
		if !strings.Contains(err.Error(), "404") && !strings.Contains(err.Error(), "not found") {
			resp.Diagnostics.AddError("Delete Error", err.Error())
			return
		}
	}
}

// ImportState handles resource import.
// Use "<client_id>:<sm_id>" or just "<sm_id>" as the import ID.
func (r *ScheduleMaintenanceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importId := req.ID

	var state ScheduleMaintenanceModel

	if strings.Contains(importId, ":") {
		parts := strings.SplitN(importId, ":", 2)
		state.Client = types.StringValue(parts[0])
		state.Id = types.StringValue(parts[1])
	} else {
		state.Id = types.StringValue(importId)
		state.Client = types.StringNull()
	}

	state.Name = types.StringUnknown()
	state.Description = types.StringUnknown()
	state.Status = types.StringUnknown()
	state.RunRBA = types.BoolUnknown()
	state.InstallPatch = types.BoolUnknown()
	state.CorrelateAlerts = types.BoolUnknown()
	state.RunEscalateAction = types.BoolUnknown()

	// Initialize Schedule with Unknown values — Read will populate the real values.
	state.Schedule = &ScheduleMaintenanceScheduleModel{
		Type:      types.StringUnknown(),
		StartTime: types.StringUnknown(),
		EndTime:   types.StringUnknown(),
		Timezone:  types.StringUnknown(),
		EndBy:     types.StringNull(),
		Pattern:   nil,
	}

	// Optional list fields start empty; Read will not overwrite them
	// (devices/groups/locations come in different shape from GET, user re-adds them in config)
	state.DeviceIds = types.SetNull(types.StringType)
	state.DeviceGroupIds = types.SetNull(types.StringType)
	state.SiteIds = types.SetNull(types.StringType)
	state.AlertConditions = nil

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// buildRequest converts the Terraform model to the API request
func (r *ScheduleMaintenanceResource) buildRequest(plan ScheduleMaintenanceModel) client.ScheduleMaintenanceRequest {
	apiReq := client.ScheduleMaintenanceRequest{
		Name:              plan.Name.ValueString(),
		Description:       plan.Description.ValueString(),
		RunRBA:            plan.RunRBA.ValueBool(),
		InstallPatch:      plan.InstallPatch.ValueBool(),
		CorrelateAlerts:   plan.CorrelateAlerts.ValueBool(),
		RunEscalateAction: plan.RunEscalateAction.ValueBool(),
	}

	// Build schedule
	if plan.Schedule != nil {
		apiReq.Schedule = client.ScheduleMaintenanceTime{
			Type:      plan.Schedule.Type.ValueString(),
			StartTime: plan.Schedule.StartTime.ValueString(),
			EndTime:   plan.Schedule.EndTime.ValueString(),
			Timezone:  plan.Schedule.Timezone.ValueString(),
		}
		if !plan.Schedule.EndBy.IsNull() && plan.Schedule.EndBy.ValueString() != "" {
			apiReq.Schedule.EndBy = plan.Schedule.EndBy.ValueString()
		}
		if plan.Schedule.Pattern != nil {
			apiReq.Schedule.Pattern = &client.SMSchedulePattern{
				Type: plan.Schedule.Pattern.Type.ValueString(),
			}
			if plan.Schedule.Pattern.WeekDays.ValueString() != "" {
				apiReq.Schedule.Pattern.WeekDays = plan.Schedule.Pattern.WeekDays.ValueString()
			}
			if plan.Schedule.Pattern.DayOfWeek.ValueString() != "" {
				apiReq.Schedule.Pattern.DayOfWeek = plan.Schedule.Pattern.DayOfWeek.ValueString()
			}
			if plan.Schedule.Pattern.WeekIndex.ValueString() != "" {
				apiReq.Schedule.Pattern.WeekIndex = plan.Schedule.Pattern.WeekIndex.ValueString()
			}
			if plan.Schedule.Pattern.DayFrequency.ValueString() != "" {
				apiReq.Schedule.Pattern.DayFrequency = plan.Schedule.Pattern.DayFrequency.ValueString()
			}
			if plan.Schedule.Pattern.Months.ValueString() != "" {
				apiReq.Schedule.Pattern.Months = plan.Schedule.Pattern.Months.ValueString()
			}
			if plan.Schedule.Pattern.DayOfMonth.ValueString() != "" {
				apiReq.Schedule.Pattern.DayOfMonth = plan.Schedule.Pattern.DayOfMonth.ValueString()
			}
			if !plan.Schedule.Pattern.RepeatFrecuency.IsNull() {
				apiReq.Schedule.Pattern.RepeatFrequency = int(plan.Schedule.Pattern.RepeatFrecuency.ValueInt64())
			}
		}
	}

	// Build devices
	for _, d := range plan.DeviceIds.Elements() {
		dev := client.ScheduleDevice{UniqueId: d.(types.String).ValueString()}
		apiReq.Devices = append(apiReq.Devices, dev)
	}

	// Build device groups
	for _, dg := range plan.DeviceGroupIds.Elements() {
		group := client.ScheduleDeviceGroup{Id: dg.(types.String).ValueString()}
		apiReq.DeviceGroups = append(apiReq.DeviceGroups, group)
	}

	// Build sites
	for _, loc := range plan.SiteIds.Elements() {
		l := client.ScheduleLocation{Id: loc.(types.String).ValueString()}
		apiReq.Locations = append(apiReq.Locations, l)
	}

	// Build alert conditions
	if plan.AlertConditions != nil {
		apiReq.AlertConditions = &client.ScheduleAlertConditions{
			MatchingType: plan.AlertConditions.MatchingType.ValueString(),
		}
		for _, rule := range plan.AlertConditions.Rules {
			apiReq.AlertConditions.Rules = append(apiReq.AlertConditions.Rules, client.ScheduleAlertRule{
				Key:      rule.Key.ValueString(),
				Operator: rule.Operator.ValueString(),
				Value:    rule.Value.ValueString(),
			})
		}
	}

	return apiReq
}

// ModifyPlan validates cross-field constraints on the plan.
func (r *ScheduleMaintenanceResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	r.BaseResource.ModifyPlan(ctx, req, resp)

	// Skip on destroy
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan ScheduleMaintenanceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Schedule == nil {
		return
	}

	// Skip while schedule type is still unknown (during import plan)
	if plan.Schedule.Type.IsUnknown() {
		return
	}

	if plan.Schedule.Type.ValueString() == "Recurring" {
		if plan.Schedule.Pattern == nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("schedule").AtName("pattern"),
				"Missing required attribute",
				"`schedule.pattern` is required when `schedule.type` is \"Recurring\".",
			)
		}
	}

	if plan.Schedule.Type.ValueString() == "One-Time" {
		if plan.Schedule.Pattern != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("schedule").AtName("pattern"),
				"Invalid attribute",
				"`schedule.pattern` must not be set when `schedule.type` is \"One-Time\".",
			)
		}
	}
}

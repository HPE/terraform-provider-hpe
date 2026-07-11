// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/HPE/terraform-provider-hpe/opsramp/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &LogAlertDefinitionResource{}

type LogAlertDefinitionResource struct {
	BaseResource
}

type LogAlertDefinitionModel struct {
	Id          types.String `tfsdk:"id"`
	Client      types.String `tfsdk:"client"`
	Name        types.String `tfsdk:"name"`
	Query       types.String `tfsdk:"query"`
	AlertNoData types.String `tfsdk:"alert_no_data"`
	Status      types.String `tfsdk:"status"`
	HealQuery   types.String `tfsdk:"heal_query"`
	EntityType  types.String `tfsdk:"entity_type"`
	Component   types.String `tfsdk:"component"`
	Subject     types.String `tfsdk:"subject"`
	Description types.String `tfsdk:"description"`

	// Conditions
	Conditions []LogAlertConditionModel `tfsdk:"conditions"`

	// Schedule
	Schedule *LogAlertScheduleModel `tfsdk:"schedule"`

	// Resource attributes & labels
	ResourceAttributes types.Map `tfsdk:"resource_attributes"`
	Labels             types.Map `tfsdk:"labels"`
}

type LogAlertConditionModel struct {
	Severity types.String `tfsdk:"severity"`
	Operator types.String `tfsdk:"operator"`
	Value    types.Int64  `tfsdk:"value"`
}

type LogAlertScheduleModel struct {
	StartTime types.String         `tfsdk:"start_time"`
	EndTime   types.String         `tfsdk:"end_time"`
	Timezone  types.String         `tfsdk:"timezone"`
	Pattern   LogAlertPatternModel `tfsdk:"pattern"`
}

type LogAlertPatternModel struct {
	Type            types.String `tfsdk:"type"`
	RepeatFrequency types.Int64  `tfsdk:"repeat_frequency"`
	WeekDays        types.String `tfsdk:"week_days"`
	DayOfMonth      types.String `tfsdk:"day_of_month"`
	WeekIndex       types.String `tfsdk:"week_index"`
	DayOfWeek       types.String `tfsdk:"day_of_week"`
}

func NewLogAlertDefinition() resource.Resource {
	return &LogAlertDefinitionResource{}
}

func (r *LogAlertDefinitionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_log_alert_definition"
}

func (r *LogAlertDefinitionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an OpsRamp Log Alert Definition. Creates alert definitions based on log filter queries.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the log alert definition.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"client": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The unique ID of the client (sub-tenant). If not provided, the alert definition is created at the provider tenant level.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the log alert definition.",
			},
			"query": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The log filter query (LogQL syntax).",
			},
			"alert_no_data": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Action when no logs found (`noalert`, `critical`, `warning`).",
				Validators: []validator.String{
					stringvalidator.OneOf("noalert", "critical", "warning"),
				},
				Default: stringdefault.StaticString("noalert"),
			},
			"status": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Alert status (`enabled` or `disabled`).",
				Validators: []validator.String{
					stringvalidator.OneOf("enabled", "disabled"),
				},
				Default: stringdefault.StaticString("enabled"),
			},
			"heal_query": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Log filter query to auto-heal alerts. Use `noAutoHeal` to disable auto-healing or empty to heal on no data",
			},
			"entity_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Entity type for the alert (`RESOURCE` or `CLIENT`).",
				Validators: []validator.String{
					stringvalidator.OneOf("RESOURCE", "CLIENT"),
				},
			},
			"component": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Component identifier for the alert.",
			},
			"subject": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The alert subject. Supports `$variable` tokens.",
			},
			"description": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The alert description. Supports `$variable` tokens.",
			},
			"conditions": schema.ListNestedAttribute{
				Required:            true,
				MarkdownDescription: "Alert severity conditions.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"severity": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Severity level (`warning` or `critical`).",
							Validators: []validator.String{
								stringvalidator.OneOf("warning", "critical"),
							},
						},
						"operator": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Comparison operator (e.g. `>`, `>=`, `<`, `<=`, `!=`, `=`).",
						},
						"value": schema.Int64Attribute{
							Required:            true,
							MarkdownDescription: "Threshold value.",
						},
					},
				},
			},
			"schedule": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Schedule configuration for the log alert evaluation.",
				Attributes: map[string]schema.Attribute{
					"start_time": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Start time for the lookup period (mandatory for daily, weekly, monthly patterns). Format: `HH:mm:ss+0000`. Difference with end_time must not exceed 2 hours.",
					},
					"end_time": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "End time for the lookup period (mandatory for daily, weekly, monthly patterns). Format: `HH:mm:ss+0000`. Difference with start_time must not exceed 2 hours.",
					},
					"timezone": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Timezone for the schedule (e.g. `America/Los_Angeles`, `Pacific/Honolulu`).",
					},
					"pattern": schema.SingleNestedAttribute{
						Required:            true,
						MarkdownDescription: "Pattern configuration defining the alert schedule frequency.",
						Attributes: map[string]schema.Attribute{
							"type": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "The schedule pattern type.",
								Validators: []validator.String{
									stringvalidator.OneOf("second", "minute", "daily", "weekly", "monthly"),
								},
							},
							"repeat_frequency": schema.Int64Attribute{
								Optional:            true,
								MarkdownDescription: "Frequency for second, minute, or daily patterns. Valid minute values: 1, 2, 3, 4, 5, 6, 10, 15, 20, 30.",
							},
							"week_days": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "Comma-separated days for weekly pattern (e.g. `monday,thursday`). Valid: sunday, monday, tuesday, wednesday, thursday, friday, saturday.",
							},
							"day_of_month": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "Comma-separated days of month for monthly pattern (e.g. `4,6`). Valid: 1-31.",
							},
							"week_index": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "Week index for monthly pattern (e.g. `Second`). Valid: First, Second, Third, Fourth, Fifth, Last.",
								Validators: []validator.String{
									stringvalidator.OneOf("First", "Second", "Third", "Fourth", "Fifth", "Last"),
								},
							},
							"day_of_week": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "Day of week for monthly pattern with week_index (e.g. `tuesday`). Valid: sunday, monday, tuesday, wednesday, thursday, friday, saturday.",
								Validators: []validator.String{
									stringvalidator.OneOf("sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"),
								},
							},
						},
					},
				},
			},
			"resource_attributes": schema.MapAttribute{
				Optional:            true,
				MarkdownDescription: "Resource attributes for the alert (required when entity_type is RESOURCE). Keys like `hostName`, values can use `$variable` tokens.",
				ElementType:         types.StringType,
			},
			"labels": schema.MapAttribute{
				Optional:            true,
				MarkdownDescription: "Custom alert tags displayed in alert details.",
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *LogAlertDefinitionResource) getTenantId(clientId types.String) string {
	if !clientId.IsNull() && clientId.ValueString() != "" {
		return clientId.ValueString()
	}
	return r.apiClient.TenantId
}

func (r *LogAlertDefinitionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan LogAlertDefinitionModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.getTenantId(plan.Client)

	alert, err := r.buildAlert(&plan)
	if err != nil {
		resp.Diagnostics.AddError("Build Request Error", err.Error())
		return
	}

	apiReq := client.LogAlertDefinitionRequest{
		Alerts: []client.LogAlertWrapper{{Alert: *alert}},
	}

	result, err := r.apiClient.CreateLogAlertDefinition(tenantId, apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", err.Error())
		return
	}

	if len(result.Errors) > 0 {
		resp.Diagnostics.AddError("API Errors", fmt.Sprintf("%v", result.Errors))
		return
	}

	if len(result.Alerts) == 0 {
		resp.Diagnostics.AddError("Create Error", "No alert definition returned in response")
		return
	}

	plan.Id = types.StringValue(result.Alerts[0].Alert.AlertID)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *LogAlertDefinitionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state LogAlertDefinitionModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.getTenantId(state.Client)

	alert, err := r.apiClient.GetLogAlertDefinition(tenantId, state.Id.ValueString())
	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "404") || strings.Contains(errStr, "not found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read Error", err.Error())
		return
	}

	// Update state from API response
	state.Name = types.StringValue(alert.Name)
	state.Query = types.StringValue(alert.Query)
	state.AlertNoData = types.StringValue(alert.AlertNoData)
	state.Status = types.StringValue(alert.Status)
	state.EntityType = types.StringValue(alert.EntityType)
	state.Subject = types.StringValue(alert.Notification.Subject)
	state.Description = types.StringValue(alert.Notification.Description)

	if alert.HealQuery != "" {
		state.HealQuery = types.StringValue(alert.HealQuery)
	}
	if alert.Component != "" {
		state.Component = types.StringValue(alert.Component)
	}

	// Schedule
	if alert.Schedule != nil {
		schedule := &LogAlertScheduleModel{
			Pattern: LogAlertPatternModel{
				Type: types.StringValue(alert.Schedule.Pattern.Type),
			},
		}
		if alert.Schedule.StartTime != "" {
			schedule.StartTime = types.StringValue(alert.Schedule.StartTime)
		} else {
			schedule.StartTime = types.StringNull()
		}
		if alert.Schedule.EndTime != "" {
			schedule.EndTime = types.StringValue(alert.Schedule.EndTime)
		} else {
			schedule.EndTime = types.StringNull()
		}
		if alert.Schedule.Timezone != "" {
			schedule.Timezone = types.StringValue(alert.Schedule.Timezone)
		} else {
			schedule.Timezone = types.StringNull()
		}
		if alert.Schedule.Pattern.RepeatFrequency != 0 {
			schedule.Pattern.RepeatFrequency = types.Int64Value(int64(alert.Schedule.Pattern.RepeatFrequency))
		} else {
			schedule.Pattern.RepeatFrequency = types.Int64Null()
		}
		if alert.Schedule.Pattern.WeekDays != "" {
			schedule.Pattern.WeekDays = types.StringValue(alert.Schedule.Pattern.WeekDays)
		} else {
			schedule.Pattern.WeekDays = types.StringNull()
		}
		if alert.Schedule.Pattern.DayOfMonth != "" {
			schedule.Pattern.DayOfMonth = types.StringValue(alert.Schedule.Pattern.DayOfMonth)
		} else {
			schedule.Pattern.DayOfMonth = types.StringNull()
		}
		if alert.Schedule.Pattern.WeekIndex != "" {
			schedule.Pattern.WeekIndex = types.StringValue(alert.Schedule.Pattern.WeekIndex)
		} else {
			schedule.Pattern.WeekIndex = types.StringNull()
		}
		if alert.Schedule.Pattern.DayOfWeek != "" {
			schedule.Pattern.DayOfWeek = types.StringValue(alert.Schedule.Pattern.DayOfWeek)
		} else {
			schedule.Pattern.DayOfWeek = types.StringNull()
		}
		state.Schedule = schedule
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *LogAlertDefinitionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan LogAlertDefinitionModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state LogAlertDefinitionModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.getTenantId(plan.Client)

	alert, err := r.buildAlert(&plan)
	if err != nil {
		resp.Diagnostics.AddError("Build Request Error", err.Error())
		return
	}

	alert.AlertID = state.Id.ValueString()
	alert.TenantID = tenantId

	_, err = r.apiClient.UpdateLogAlertDefinition(tenantId, state.Id.ValueString(), *alert)
	if err != nil {
		resp.Diagnostics.AddError("Update Error", err.Error())
		return
	}

	plan.Id = state.Id

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *LogAlertDefinitionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state LogAlertDefinitionModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.getTenantId(state.Client)

	err := r.apiClient.DeleteLogAlertDefinition(tenantId, state.Id.ValueString())
	if err != nil {
		if !strings.Contains(err.Error(), "404") && !strings.Contains(err.Error(), "not found") {
			resp.Diagnostics.AddError("Delete Error", err.Error())
			return
		}
	}
}

// buildAlert converts the Terraform model to the API struct
func (r *LogAlertDefinitionResource) buildAlert(plan *LogAlertDefinitionModel) (*client.LogAlertDefinition, error) {
	alert := &client.LogAlertDefinition{
		Name:        plan.Name.ValueString(),
		Type:        "log",
		Query:       plan.Query.ValueString(),
		AlertNoData: plan.AlertNoData.ValueString(),
		Status:      plan.Status.ValueString(),
		EntityType:  plan.EntityType.ValueString(),
		Notification: client.LogAlertNotification{
			Subject:     plan.Subject.ValueString(),
			Description: plan.Description.ValueString(),
		},
	}

	if !plan.HealQuery.IsNull() {
		alert.HealQuery = plan.HealQuery.ValueString()
	}
	if !plan.Component.IsNull() {
		alert.Component = plan.Component.ValueString()
	}

	// Conditions
	for _, c := range plan.Conditions {
		alert.Conditions = append(alert.Conditions, client.LogAlertCondition{
			Severity: c.Severity.ValueString(),
			Operator: c.Operator.ValueString(),
			Value:    int(c.Value.ValueInt64()),
		})
	}

	// Resource attributes
	if !plan.ResourceAttributes.IsNull() && len(plan.ResourceAttributes.Elements()) > 0 {
		alert.ResourceAttributes = make(map[string]string)
		for k, v := range plan.ResourceAttributes.Elements() {
			alert.ResourceAttributes[k] = v.(types.String).ValueString()
		}
	}

	// Labels
	if !plan.Labels.IsNull() && len(plan.Labels.Elements()) > 0 {
		alert.Labels = make(map[string]string)
		for k, v := range plan.Labels.Elements() {
			alert.Labels[k] = v.(types.String).ValueString()
		}
	}

	// Schedule
	if plan.Schedule != nil {
		schedule := &client.LogAlertSchedule{
			Pattern: client.LogAlertPattern{
				Type: plan.Schedule.Pattern.Type.ValueString(),
			},
		}
		if !plan.Schedule.StartTime.IsNull() {
			schedule.StartTime = plan.Schedule.StartTime.ValueString()
		}
		if !plan.Schedule.EndTime.IsNull() {
			schedule.EndTime = plan.Schedule.EndTime.ValueString()
		}
		if !plan.Schedule.Timezone.IsNull() {
			schedule.Timezone = plan.Schedule.Timezone.ValueString()
		}
		if !plan.Schedule.Pattern.RepeatFrequency.IsNull() {
			schedule.Pattern.RepeatFrequency = int(plan.Schedule.Pattern.RepeatFrequency.ValueInt64())
		}
		if !plan.Schedule.Pattern.WeekDays.IsNull() {
			schedule.Pattern.WeekDays = plan.Schedule.Pattern.WeekDays.ValueString()
		}
		if !plan.Schedule.Pattern.DayOfMonth.IsNull() {
			schedule.Pattern.DayOfMonth = plan.Schedule.Pattern.DayOfMonth.ValueString()
		}
		if !plan.Schedule.Pattern.WeekIndex.IsNull() {
			schedule.Pattern.WeekIndex = plan.Schedule.Pattern.WeekIndex.ValueString()
		}
		if !plan.Schedule.Pattern.DayOfWeek.IsNull() {
			schedule.Pattern.DayOfWeek = plan.Schedule.Pattern.DayOfWeek.ValueString()
		}
		alert.Schedule = schedule
	}

	return alert, nil
}

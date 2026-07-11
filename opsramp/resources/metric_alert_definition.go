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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &MetricAlertDefinitionResource{}
var _ resource.ResourceWithModifyPlan = &MetricAlertDefinitionResource{}

type MetricAlertDefinitionResource struct {
	BaseResource
}

type MetricAlertDefinitionModel struct {
	Id                   types.String            `tfsdk:"id"`
	Client               types.String            `tfsdk:"client"`
	Name                 types.String            `tfsdk:"name"`
	Query                types.String            `tfsdk:"query"`
	AlertType            types.String            `tfsdk:"alert_type"`
	AlertThresholdType   types.String            `tfsdk:"alert_threshold_type"`
	AlertThresholdData   AlertThresholdDataModel `tfsdk:"alert_threshold_data"`
	NoDataCondition      types.String            `tfsdk:"no_data_condition"`
	AlertTriggerDuration types.String            `tfsdk:"alert_trigger_duration"`
	Subject              types.String            `tfsdk:"subject"`
	Description          types.String            `tfsdk:"description"`
	EntityType           types.List              `tfsdk:"entity_type"`
	Component            types.List              `tfsdk:"component"`
	Status               types.Bool              `tfsdk:"status"`
	IsObsolete           types.Bool              `tfsdk:"is_obsolete"`

	Labels     []NameValuePairModel `tfsdk:"labels"`
	Attributes []NameValuePairModel `tfsdk:"attributes"`
}

type AlertThresholdDataModel struct {
	WarningCondition  types.String `tfsdk:"warning_condition"`
	CriticalCondition types.String `tfsdk:"critical_condition"`
	Limit             types.Int64  `tfsdk:"limit"`
	Direction         types.String `tfsdk:"direction"`
	LearningPeriod    types.String `tfsdk:"learning_period"`
	StandardDeviation types.Int64  `tfsdk:"standard_deviation"`
	LowerLimit        types.String `tfsdk:"lower_limit"`
	UpperLimit        types.String `tfsdk:"upper_limit"`
}

type NameValuePairModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

func NewMetricAlertDefinition() resource.Resource {
	return &MetricAlertDefinitionResource{}
}

func (r *MetricAlertDefinitionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_metric_alert_definition"
}

func (r *MetricAlertDefinitionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an OpsRamp Metric Alert Definition. Creates alert definitions based on PromQL queries with static thresholds.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the alert definition.",
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
				MarkdownDescription: "The name of the alert definition. Must be unique across the client.",
			},
			"query": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The PromQL query for the alert definition.",
			},
			"alert_type": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The alert type (`METRICS`, `TRACE`).",
				Validators: []validator.String{
					stringvalidator.OneOf("METRICS", "TRACE"),
				},
				Default: stringdefault.StaticString("METRICS"),
			},
			"alert_threshold_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The threshold type (`STATIC_THRESHOLD`, `FORECAST`, `DYNAMIC_CHANGE_DETECTION`, `DYNAMIC_THRESHOLD`).",
				Validators: []validator.String{
					stringvalidator.OneOf("STATIC_THRESHOLD", "FORECAST", "DYNAMIC_CHANGE_DETECTION", "DYNAMIC_THRESHOLD"),
				},
			},
			"alert_threshold_data": schema.SingleNestedAttribute{
				Required:            true,
				MarkdownDescription: "The threshold data for the alert definition, including warning and critical conditions.",
				Attributes: map[string]schema.Attribute{
					"warning_condition": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The warning threshold condition. Comma-separated values with operators (e.g. `>0`, `>=12,<3,8-20`).",
					},
					"critical_condition": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The critical threshold condition. Comma-separated values with operators (e.g. `>0`, `>=12,<3,8-20`).",
					},
					"limit": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "The threshold limit for forecast or dynamic change detection types. Required if alert_threshold_type is not STATIC_THRESHOLD.",
					},
					"direction": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The direction for dynamic change detection (e.g. `increase`, `decrease`, `increaseordecrease`). Required if alert_threshold_type is DYNAMIC_CHANGE_DETECTION.",
						Validators: []validator.String{
							stringvalidator.OneOf("increase", "decrease", "increaseordecrease"),
						},
					},
					"learning_period": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The learning period for dynamic thresholds (e.g. `7d`).",
					},
					"standard_deviation": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "The number of standard deviations for dynamic thresholds.",
					},
					"lower_limit": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The lower limit for dynamic thresholds. Required if alert_threshold_type is DYNAMIC_THRESHOLD.",
					},
					"upper_limit": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The upper limit for dynamic thresholds. Required if alert_threshold_type is DYNAMIC_THRESHOLD.",
					},
				},
			},
			"no_data_condition": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Action when no data is received. Required when alert_threshold_type is STATIC_THRESHOLD or DYNAMIC_THRESHOLD. (e.g.: `NO_DATA_ALERT`, `WARNING_ALERT`, `CRITICAL_ALERT`).",
				Validators: []validator.String{
					stringvalidator.OneOf("NO_DATA_ALERT", "WARNING_ALERT", "CRITICAL_ALERT"),
				},
			},
			"alert_trigger_duration": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Duration the threshold must be breached before alerting (e.g. `0m`, `5m`, `30s`).",
			},
			"subject": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The alert subject. Supports tokens like `{{$host}}`.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The alert description. Supports tokens like `{{$host}}`.",
			},
			"entity_type": schema.ListAttribute{
				Required:            true,
				MarkdownDescription: "The entity type for the alert (e.g. `[\"RESOURCE\"]`).",
				ElementType:         types.StringType,
			},
			"component": schema.ListAttribute{
				Optional:            true,
				MarkdownDescription: "The alert component identifiers (e.g. `[\"$ip\"]`).",
				ElementType:         types.StringType,
			},
			"status": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Whether alerts are generated. Defaults to `true`.",
			},
			"is_obsolete": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "If true, all alerts generated on this definition become Obsolete. Only used on update.",
			},
			"labels": schema.ListNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Metric labels included in the alert.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Label name.",
						},
						"value": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Label value (can use `$variable` tokens).",
						},
					},
				},
			},
			"attributes": schema.ListNestedAttribute{
				Required:            true,
				MarkdownDescription: "Resource attributes for the alert. Each entry must have name as `name`, `host`, `ip`, or `uuid`.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Attribute name (`name`, `host`, `ip`, or `uuid`).",
							Validators: []validator.String{
								stringvalidator.OneOf("name", "host", "ip", "uuid"),
							},
						},
						"value": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Attribute value (can use `$variable` tokens).",
						},
					},
				},
			},
		},
	}
}

func (r *MetricAlertDefinitionResource) getTenantId(clientId types.String) string {
	if !clientId.IsNull() && clientId.ValueString() != "" {
		return clientId.ValueString()
	}
	return r.apiClient.TenantId
}

func (r *MetricAlertDefinitionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan MetricAlertDefinitionModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.getTenantId(plan.Client)

	apiReq, err := r.buildRequest(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError("Build Request Error", err.Error())
		return
	}

	result, err := r.apiClient.CreateMetricAlertDefinition(tenantId, *apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", err.Error())
		return
	}

	plan.Id = types.StringValue(result.GetID())

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *MetricAlertDefinitionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state MetricAlertDefinitionModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The metric-alerts API does not have a GET by ID endpoint.
	// We preserve state as-is. Changes are detected via Terraform plan diffs.
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *MetricAlertDefinitionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan MetricAlertDefinitionModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state MetricAlertDefinitionModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.getTenantId(plan.Client)

	apiReq, err := r.buildRequest(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError("Build Request Error", err.Error())
		return
	}

	// Set isObsolete if specified
	if !plan.IsObsolete.IsNull() {
		val := plan.IsObsolete.ValueBool()
		apiReq.IsObsolete = &val
	}

	_, err = r.apiClient.UpdateMetricAlertDefinition(tenantId, state.Id.ValueString(), *apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Update Error", err.Error())
		return
	}

	plan.Id = state.Id

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *MetricAlertDefinitionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state MetricAlertDefinitionModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.getTenantId(state.Client)

	err := r.apiClient.DeleteMetricAlertDefinition(tenantId, state.Id.ValueString())
	if err != nil {
		if !strings.Contains(err.Error(), "404") && !strings.Contains(err.Error(), "not found") {
			resp.Diagnostics.AddError("Delete Error", err.Error())
			return
		}
	}
}

// buildRequest converts the Terraform model to the API request
func (r *MetricAlertDefinitionResource) buildRequest(ctx context.Context, plan *MetricAlertDefinitionModel) (*client.MetricAlertDefinitionRequest, error) {
	apiReq := &client.MetricAlertDefinitionRequest{
		Name:                 plan.Name.ValueString(),
		Query:                plan.Query.ValueString(),
		AlertType:            plan.AlertType.ValueString(),
		AlertThresholdType:   plan.AlertThresholdType.ValueString(),
		AlertTriggerDuration: plan.AlertTriggerDuration.ValueString(),
		Subject:              plan.Subject.ValueString(),
		Description:          plan.Description.ValueString(),
		Status:               plan.Status.ValueBool(),
	}

	// Threshold data
	apiReq.AlertThresholdData = client.MetricAlertThresholdData{}

	if plan.AlertThresholdData.Direction.ValueString() != "" {
		apiReq.AlertThresholdData.Direction = plan.AlertThresholdData.Direction.ValueString()
	}
	if plan.AlertThresholdData.LearningPeriod.ValueString() != "" {
		apiReq.AlertThresholdData.LearningPeriod = plan.AlertThresholdData.LearningPeriod.ValueString()
	}
	if !plan.AlertThresholdData.StandardDeviation.IsNull() {
		apiReq.AlertThresholdData.StandardDeviation = plan.AlertThresholdData.StandardDeviation.ValueInt64()
	}
	if !plan.AlertThresholdData.Limit.IsNull() {
		apiReq.AlertThresholdData.Limit = plan.AlertThresholdData.Limit.ValueInt64()
	}
	if plan.AlertThresholdData.WarningCondition.ValueString() != "" {
		apiReq.AlertThresholdData.WarningCondition = plan.AlertThresholdData.WarningCondition.ValueString()
	}
	if plan.AlertThresholdData.CriticalCondition.ValueString() != "" {
		apiReq.AlertThresholdData.CriticalCondition = plan.AlertThresholdData.CriticalCondition.ValueString()
	}

	if plan.NoDataCondition.ValueString() != "" {
		apiReq.NoDataCondition = plan.NoDataCondition.ValueString()
	}

	// Entity type
	if !plan.EntityType.IsNull() {
		var entityTypes []string
		diags := plan.EntityType.ElementsAs(ctx, &entityTypes, false)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to parse entity_type")
		}
		apiReq.EntityType = entityTypes
	}

	// Component
	if !plan.Component.IsNull() {
		var components []string
		diags := plan.Component.ElementsAs(ctx, &components, false)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to parse component")
		}
		apiReq.Component = components
	}

	// Labels
	if plan.Labels != nil {
		for _, l := range plan.Labels {
			apiReq.Labels = append(apiReq.Labels, client.NameValuePair{
				Name:  l.Name.ValueString(),
				Value: l.Value.ValueString(),
			})
		}
	}

	// Attributes
	for _, a := range plan.Attributes {
		apiReq.Attributes = append(apiReq.Attributes, client.NameValuePair{
			Name:  a.Name.ValueString(),
			Value: a.Value.ValueString(),
		})
	}

	return apiReq, nil
}

// modify plan
func (r *MetricAlertDefinitionResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Don't validate during destroy — plan is null.
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan MetricAlertDefinitionModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if plan.NoDataCondition.ValueString() == "" {
		if plan.AlertThresholdType.ValueString() == "STATIC_THRESHOLD" || plan.AlertThresholdType.ValueString() == "DYNAMIC_THRESHOLD" {
			diags.AddError("No Data Condition Required", "no_data_condition must be specified when alert_threshold_type is STATIC_THRESHOLD or DYNAMIC_THRESHOLD")
		}
	}

	if plan.AlertThresholdType.ValueString() == "STATIC_THRESHOLD" {
		if plan.AlertThresholdData.WarningCondition.ValueString() == "" && plan.AlertThresholdData.CriticalCondition.ValueString() == "" {
			diags.AddError("Threshold Condition Required", "At least one of warning_condition or critical_condition must be specified when alert_threshold_type is STATIC_THRESHOLD or FORECAST")
		}
	}

	if plan.AlertThresholdType.ValueString() == "DYNAMIC_CHANGE_DETECTION" {
		if plan.AlertThresholdData.Direction.ValueString() == "" ||
			plan.AlertThresholdData.LearningPeriod.ValueString() == "" ||
			plan.AlertThresholdData.StandardDeviation.IsNull() {
			diags.AddError("Direction, Learning Period, and Standard Deviation Required", "direction, learning_period, and standard_deviation must be specified when alert_threshold_type is DYNAMIC_CHANGE_DETECTION")
		}

		if plan.NoDataCondition.ValueString() != "" {
			diags.AddError("No Data Condition Not Supported", "no_data_condition is not supported when alert_threshold_type is DYNAMIC_CHANGE_DETECTION")
		}
	}

	if plan.AlertThresholdType.ValueString() == "FORECAST" {
		if plan.AlertThresholdData.Limit.IsNull() {
			diags.AddError("Limit Required", "limit must be specified when alert_threshold_type is FORECAST")
		}

		if plan.AlertTriggerDuration.ValueString() != "" {
			diags.AddError("Alert Trigger Duration Not Supported", "alert_trigger_duration is not supported when alert_threshold_type is FORECAST")
		}

		if plan.NoDataCondition.ValueString() != "" {
			diags.AddError("No Data Condition Not Supported", "no_data_condition is not supported when alert_threshold_type is FORECAST")
		}
	} else {
		// opposite side-effect check
		if plan.AlertTriggerDuration.ValueString() == "" {
			diags.AddError("Alert Trigger Duration Required", "alert_trigger_duration must be specified when alert_threshold_type is not FORECAST")
		}
	}

	if plan.AlertThresholdType.ValueString() == "DYNAMIC_THRESHOLD" {
		if plan.AlertThresholdData.Limit.IsNull() {
			diags.AddError("Limit Required", "limit must be specified when alert_threshold_type is DYNAMIC_THRESHOLD")
		}
	}
}

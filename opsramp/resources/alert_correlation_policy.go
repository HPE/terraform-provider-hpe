// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package resources

import (
	"context"
	"strings"

	"github.com/HPE/terraform-provider-opsramp/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure implementation satisfies the expected interfaces
var _ resource.Resource = &AlertCorrelationPolicyResource{}
var _ resource.ResourceWithModifyPlan = &AlertCorrelationPolicyResource{}

// AlertCorrelationPolicyResource defines the resource implementation.
type AlertCorrelationPolicyResource struct {
	BaseResource
}

// AlertCorrelationPolicyModel maps Terraform schema attributes to the provider model.
type AlertCorrelationPolicyModel struct {
	Client               types.String               `tfsdk:"client"`
	Id                   types.String               `tfsdk:"id"`
	Name                 types.String               `tfsdk:"name"`
	EnabledMode          types.String               `tfsdk:"enabled_mode"`
	Precedence           types.Int64                `tfsdk:"precedence"`
	Type                 types.String               `tfsdk:"type"`
	FilterQuery          types.String               `tfsdk:"filter_query"`
	InferenceQuery       types.String               `tfsdk:"inference_query"`
	Review               types.Bool                 `tfsdk:"review"`
	InferenceSubject     types.String               `tfsdk:"inference_subject"`
	AlgorithmCorrelation *AlgorithmCorrelationModel `tfsdk:"algorithm_correlation"`
	MachineLearning      *MachineLearningModel      `tfsdk:"machine_learning"`
}

// AlgorithmCorrelationModel represents algorithm correlation settings
type AlgorithmCorrelationModel struct {
	MatchingConditions []MatchingConditionModel `tfsdk:"matching_conditions"`
	AlertsTimeWindow   types.String             `tfsdk:"alerts_time_window"`
	AlertTrigger       *AlertTriggerModel       `tfsdk:"alert_trigger"`
}

// AlertTriggerModel represents alert trigger settings
type AlertTriggerModel struct {
	Rules    []AlertTriggerRuleModel `tfsdk:"rules"`
	Duration types.Int64             `tfsdk:"duration"`
}

// AlertTriggerRuleModel represents a single trigger rule
type AlertTriggerRuleModel struct {
	EntityName  types.String `tfsdk:"entity_name"`
	Operator    types.String `tfsdk:"operator"`
	EntityValue types.String `tfsdk:"entity_value"`
}

// MachineLearningModel represents machine learning settings
type MachineLearningModel struct {
	ContinuousLearning types.Bool               `tfsdk:"continuous_learning"`
	Topology           types.Bool               `tfsdk:"topology"`
	TopologyDepth      types.Int64              `tfsdk:"topology_depth"`
	MatchingConditions []MatchingConditionModel `tfsdk:"matching_conditions"`
}

// MatchingConditionModel represents a matching condition
type MatchingConditionModel struct {
	Property  types.String `tfsdk:"property"`
	MatchType types.String `tfsdk:"match_type"`
}

// NewAlertCorrelationPolicy creates a new instance of the resource.
func NewAlertCorrelationPolicy() resource.Resource {
	return &AlertCorrelationPolicyResource{}
}

// Metadata returns the resource type name.
func (r *AlertCorrelationPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alert_correlation_policy"
}

// Schema defines the schema for the resource.
func (r *AlertCorrelationPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an OpsRamp Alert Correlation Policy. Correlates related alerts to reduce noise and identify root causes.",
		Attributes: map[string]schema.Attribute{
			"client": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The client (tenant) UUID where this policy should be created. If not specified, uses the provider's tenant.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the alert correlation policy.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the alert correlation policy.",
			},
			"enabled_mode": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The enabled mode of the policy. Valid values: `ON`, `OFF`, `OBSERVED`.",
				Validators: []validator.String{
					stringvalidator.OneOf("ON", "OFF", "OBSERVED"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"precedence": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The execution order of the policy. Lower values execute first.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The alert correlation policy type. Valid values: `DEPENDENCY`, `ALGORITHM`, `CO_OCCURRENCE`.",
				Validators: []validator.String{
					stringvalidator.OneOf("DEPENDENCY", "ALGORITHM", "CO_OCCURRENCE"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"inference_subject": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Template for the inference alert subject. Use `$subject` as a placeholder.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"filter_query": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Filter query to scope which alerts are evaluated by this policy.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"inference_query": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Inference query to scope which alerts are used for inference.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"review": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether the policy is in review mode.",
			},
			"algorithm_correlation": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Algorithm-based correlation settings. Used when type is `ALGORITHM`.",
				Attributes: map[string]schema.Attribute{
					"matching_conditions": schema.ListNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Conditions for matching correlated alerts.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"property": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "The alert property to match on.",
								},
								"match_type": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "The type of match.",
								},
							},
						},
					},
					"alerts_time_window": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Time window in minutes for correlating alerts.",
					},
					"alert_trigger": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Alert trigger conditions for algorithm correlation.",
						Attributes: map[string]schema.Attribute{
							"rules": schema.ListNestedAttribute{
								Optional:            true,
								MarkdownDescription: "List of trigger rules.",
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"entity_name": schema.StringAttribute{
											Required:            true,
											MarkdownDescription: "The entity name.",
										},
										"operator": schema.StringAttribute{
											Required:            true,
											MarkdownDescription: "The comparison operator.",
										},
										"entity_value": schema.StringAttribute{
											Required:            true,
											MarkdownDescription: "The value to compare against.",
										},
									},
								},
							},
							"duration": schema.Int64Attribute{
								Optional:            true,
								MarkdownDescription: "Duration in minutes for the trigger.",
							},
						},
					},
				},
			},
			"machine_learning": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Machine learning settings for alert correlation.",
				Attributes: map[string]schema.Attribute{
					"continuous_learning": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "Whether continuous learning is enabled.",
					},
					"topology": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "Whether topology-based correlation is enabled.",
					},
					"topology_depth": schema.Int64Attribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The depth of topology traversal for correlation.",
					},
					"matching_conditions": schema.ListNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Conditions for matching co-occurring alerts.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"property": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "The alert property to match on.",
								},
								"match_type": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "The type of match.",
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildAlertCorrelationPolicyRequest(plan AlertCorrelationPolicyModel) client.AlertCorrelationPolicy {
	policy := client.AlertCorrelationPolicy{
		Name:             plan.Name.ValueString(),
		EnabledMode:      plan.EnabledMode.ValueString(),
		Type:             plan.Type.ValueString(),
		FilterQuery:      plan.FilterQuery.ValueString(),
		InferenceQuery:   plan.InferenceQuery.ValueString(),
		Review:           plan.Review.ValueBool(),
		InferenceSubject: plan.InferenceSubject.ValueString(),
	}

	if !plan.Precedence.IsNull() && !plan.Precedence.IsUnknown() {
		policy.Precedence = int(plan.Precedence.ValueInt64())
	}

	if plan.AlgorithmCorrelation != nil {
		ac := &client.AlgorithmCorrelation{}
		if !plan.AlgorithmCorrelation.AlertsTimeWindow.IsNull() && !plan.AlgorithmCorrelation.AlertsTimeWindow.IsUnknown() {
			ac.AlertsTimeWindow = plan.AlgorithmCorrelation.AlertsTimeWindow.ValueString()
		}
		for _, mc := range plan.AlgorithmCorrelation.MatchingConditions {
			ac.MatchingConditions = append(ac.MatchingConditions, client.MatchingCondition{
				Property:  mc.Property.ValueString(),
				MatchType: mc.MatchType.ValueString(),
			})
		}
		if ac.MatchingConditions == nil {
			ac.MatchingConditions = []client.MatchingCondition{}
		}
		if plan.AlgorithmCorrelation.AlertTrigger != nil {
			at := &client.AlertTrigger{}
			if !plan.AlgorithmCorrelation.AlertTrigger.Duration.IsNull() && !plan.AlgorithmCorrelation.AlertTrigger.Duration.IsUnknown() {
				at.Duration = int(plan.AlgorithmCorrelation.AlertTrigger.Duration.ValueInt64())
			}
			for _, r := range plan.AlgorithmCorrelation.AlertTrigger.Rules {
				at.Rules = append(at.Rules, client.AlertTriggerRule{
					EntityName:  r.EntityName.ValueString(),
					Operator:    r.Operator.ValueString(),
					EntityValue: r.EntityValue.ValueString(),
				})
			}
			if at.Rules == nil {
				at.Rules = []client.AlertTriggerRule{}
			}
			ac.AlertTrigger = at
		}
		policy.AlgorithmCorrelation = ac
	}

	if plan.MachineLearning != nil {
		ml := &client.MachineLearning{
			ContinuousLearning: plan.MachineLearning.ContinuousLearning.ValueBool(),
			Topology:           plan.MachineLearning.Topology.ValueBool(),
		}
		if !plan.MachineLearning.TopologyDepth.IsNull() && !plan.MachineLearning.TopologyDepth.IsUnknown() {
			ml.TopologyDepth = int(plan.MachineLearning.TopologyDepth.ValueInt64())
		}
		for _, mc := range plan.MachineLearning.MatchingConditions {
			ml.MatchingConditions = append(ml.MatchingConditions, client.MatchingCondition{
				Property:  mc.Property.ValueString(),
				MatchType: mc.MatchType.ValueString(),
			})
		}
		if ml.MatchingConditions == nil {
			ml.MatchingConditions = []client.MatchingCondition{}
		}
		policy.MachineLearning = ml
	}

	return policy
}

func mapAlertCorrelationPolicyToState(resp *client.AlertCorrelationPolicy, state *AlertCorrelationPolicyModel) {
	state.Id = types.StringValue(resp.Id)
	state.Name = types.StringValue(resp.Name)
	state.EnabledMode = types.StringValue(resp.EnabledMode)
	state.Type = types.StringValue(resp.Type)
	state.FilterQuery = types.StringValue(resp.FilterQuery)
	state.InferenceQuery = types.StringValue(resp.InferenceQuery)
	state.Review = types.BoolValue(resp.Review)

	if resp.Precedence != 0 {
		state.Precedence = types.Int64Value(int64(resp.Precedence))
	}

	if resp.InferenceSubject != "" {
		state.InferenceSubject = types.StringValue(resp.InferenceSubject)
	}

	if resp.AlgorithmCorrelation != nil {
		ac := &AlgorithmCorrelationModel{}
		if resp.AlgorithmCorrelation.AlertsTimeWindow != "" {
			ac.AlertsTimeWindow = types.StringValue(resp.AlgorithmCorrelation.AlertsTimeWindow)
		}
		if len(resp.AlgorithmCorrelation.MatchingConditions) > 0 {
			ac.MatchingConditions = make([]MatchingConditionModel, len(resp.AlgorithmCorrelation.MatchingConditions))
			for i, mc := range resp.AlgorithmCorrelation.MatchingConditions {
				ac.MatchingConditions[i] = MatchingConditionModel{
					Property:  types.StringValue(mc.Property),
					MatchType: types.StringValue(mc.MatchType),
				}
			}
		}
		if resp.AlgorithmCorrelation.AlertTrigger != nil {
			at := &AlertTriggerModel{}
			if resp.AlgorithmCorrelation.AlertTrigger.Duration != 0 {
				at.Duration = types.Int64Value(int64(resp.AlgorithmCorrelation.AlertTrigger.Duration))
			}
			if len(resp.AlgorithmCorrelation.AlertTrigger.Rules) > 0 {
				at.Rules = make([]AlertTriggerRuleModel, len(resp.AlgorithmCorrelation.AlertTrigger.Rules))
				for i, r := range resp.AlgorithmCorrelation.AlertTrigger.Rules {
					at.Rules[i] = AlertTriggerRuleModel{
						EntityName:  types.StringValue(r.EntityName),
						Operator:    types.StringValue(r.Operator),
						EntityValue: types.StringValue(r.EntityValue),
					}
				}
			}
			ac.AlertTrigger = at
		}
		state.AlgorithmCorrelation = ac
	}

	if resp.MachineLearning != nil {
		ml := &MachineLearningModel{
			ContinuousLearning: types.BoolValue(resp.MachineLearning.ContinuousLearning),
			Topology:           types.BoolValue(resp.MachineLearning.Topology),
		}
		if resp.MachineLearning.TopologyDepth != 0 {
			ml.TopologyDepth = types.Int64Value(int64(resp.MachineLearning.TopologyDepth))
		}
		if len(resp.MachineLearning.MatchingConditions) > 0 {
			ml.MatchingConditions = make([]MatchingConditionModel, len(resp.MachineLearning.MatchingConditions))
			for i, mc := range resp.MachineLearning.MatchingConditions {
				ml.MatchingConditions[i] = MatchingConditionModel{
					Property:  types.StringValue(mc.Property),
					MatchType: types.StringValue(mc.MatchType),
				}
			}
		}
		state.MachineLearning = ml
	}
}

// Create handles the creation of the resource.
func (r *AlertCorrelationPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AlertCorrelationPolicyModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.apiClient.TenantId
	if !plan.Client.IsNull() && plan.Client.ValueString() != "" {
		tenantId = plan.Client.ValueString()
	}

	policy := buildAlertCorrelationPolicyRequest(plan)

	created, err := r.apiClient.CreateAlertCorrelationPolicy(tenantId, policy)
	if err != nil {
		resp.Diagnostics.AddError("Creation Error", err.Error())
		return
	}

	mapAlertCorrelationPolicyToState(created, &plan)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Read handles reading the resource state.
func (r *AlertCorrelationPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AlertCorrelationPolicyModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.apiClient.TenantId
	if !state.Client.IsNull() && state.Client.ValueString() != "" {
		tenantId = state.Client.ValueString()
	}

	existing, err := r.apiClient.GetAlertCorrelationPolicy(tenantId, state.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read Error", err.Error())
		return
	}

	mapAlertCorrelationPolicyToState(existing, &state)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update handles updating the resource.
func (r *AlertCorrelationPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state AlertCorrelationPolicyModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.apiClient.TenantId
	if !state.Client.IsNull() && state.Client.ValueString() != "" {
		tenantId = state.Client.ValueString()
	}

	policy := buildAlertCorrelationPolicyRequest(plan)

	updated, err := r.apiClient.UpdateAlertCorrelationPolicy(tenantId, state.Id.ValueString(), policy)
	if err != nil {
		resp.Diagnostics.AddError("Update Error", err.Error())
		return
	}

	plan.Client = state.Client
	mapAlertCorrelationPolicyToState(updated, &plan)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Delete handles deleting the resource.
func (r *AlertCorrelationPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AlertCorrelationPolicyModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.apiClient.TenantId
	if !state.Client.IsNull() && state.Client.ValueString() != "" {
		tenantId = state.Client.ValueString()
	}

	err := r.apiClient.DeleteAlertCorrelationPolicy(tenantId, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

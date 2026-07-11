// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package resources

import (
	"context"
	"strings"

	"github.com/HPE/terraform-provider-opsramp/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure implementation satisfies the expected interfaces
var _ resource.Resource = &FirstResponsePolicyResource{}
var _ resource.ResourceWithModifyPlan = &FirstResponsePolicyResource{}

// FirstResponsePolicyResource defines the resource implementation.
type FirstResponsePolicyResource struct {
	BaseResource
}

// FirstResponsePolicyModel maps Terraform schema attributes to the provider model.
type FirstResponsePolicyModel struct {
	Client           types.String `tfsdk:"client"`
	Id               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	EnabledMode      types.String `tfsdk:"enabled_mode"`
	FilterQuery      types.String `tfsdk:"filter_query"`
	AttributeActions types.Object `tfsdk:"attribute_actions"`
	PatternActions   types.Object `tfsdk:"pattern_actions"`
}

// FirstResponseAttrActionsModel represents attribute-based actions
type FirstResponseAttrActionsModel struct {
	ContinuousLearning types.Bool   `tfsdk:"continuous_learning"`
	Suppress           types.Object `tfsdk:"suppress"`
	Insights           types.Object `tfsdk:"insights"`
	RunProcess         types.Object `tfsdk:"run_process"`
}

// FirstResponsePatternActionsModel represents pattern-based actions
type FirstResponsePatternActionsModel struct {
	SeasonalityTimeFrame types.String `tfsdk:"seasonality_time_frame"`
	Suppress             types.Object `tfsdk:"suppress"`
}

var attrSuppressAttrTypes = map[string]attr.Type{
	"learned_configuration": types.BoolType,
	"suppress_duration":     types.Int64Type,
}

var attrRunProcessAttrTypes = map[string]attr.Type{
	"learned_configuration": types.BoolType,
	"run_immediately":       types.BoolType,
	"process_ids":           types.ListType{ElemType: types.StringType},
}

var insightsAttrTypes = map[string]attr.Type{
	"create_prc_insights": types.BoolType,
}

var patternSuppressAttrTypes = map[string]attr.Type{
	"seasonal_alerts": types.BoolType,
}

var attributeActionsAttrTypes = map[string]attr.Type{
	"continuous_learning": types.BoolType,
	"suppress":            types.ObjectType{AttrTypes: attrSuppressAttrTypes},
	"insights":            types.ObjectType{AttrTypes: insightsAttrTypes},
	"run_process":         types.ObjectType{AttrTypes: attrRunProcessAttrTypes},
}

var patternActionsAttrTypes = map[string]attr.Type{
	"seasonality_time_frame": types.StringType,
	"suppress":               types.ObjectType{AttrTypes: patternSuppressAttrTypes},
}

// NewFirstResponsePolicy creates a new instance of the resource.
func NewFirstResponsePolicy() resource.Resource {
	return &FirstResponsePolicyResource{}
}

// Metadata returns the resource type name.
func (r *FirstResponsePolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_first_response_policy"
}

// Schema defines the schema for the resource.
func (r *FirstResponsePolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an OpsRamp First Response Policy. Defines automated first response actions for incoming alerts.",
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
				MarkdownDescription: "The unique identifier of the first response policy.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the first response policy.",
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
			"filter_query": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "OpsQL filter query to scope which alerts are evaluated by this policy.",
			},
			"attribute_actions": schema.SingleNestedAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Attribute-based actions to apply when the policy matches.",
				Attributes: map[string]schema.Attribute{
					"continuous_learning": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
						MarkdownDescription: "Whether continuous learning is enabled.",
					},
					"suppress": schema.SingleNestedAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "Suppress settings for attribute actions.",
						Attributes: map[string]schema.Attribute{
							"learned_configuration": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								MarkdownDescription: "Whether to use learned configuration for suppression.",
							},
							"suppress_duration": schema.Int64Attribute{
								Optional:            true,
								Computed:            true,
								MarkdownDescription: "Duration in minutes to suppress alerts. Use -1 for indefinite.",
							},
						},
					},
					"insights": schema.SingleNestedAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "Insights settings for attribute actions.",
						Attributes: map[string]schema.Attribute{
							"create_prc_insights": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								MarkdownDescription: "Whether to create PRC insights.",
							},
						},
					},
					"run_process": schema.SingleNestedAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "Insights settings for attribute actions.",
						Attributes: map[string]schema.Attribute{
							"learned_configuration": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								MarkdownDescription: "Whether to create PRC insights.",
							},
							"run_immediately": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								MarkdownDescription: "Whether to create PRC insights.",
							},
							"process_ids": schema.ListAttribute{
								ElementType:         types.StringType,
								Optional:            true,
								Computed:            true,
								MarkdownDescription: "List of process IDs to run when the policy matches.",
							},
						},
					},
				},
			},
			"pattern_actions": schema.SingleNestedAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Pattern-based actions for the policy.",
				Attributes: map[string]schema.Attribute{
					"seasonality_time_frame": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The seasonality time frame (e.g., `10D`, `60D`).",
						Default:             stringdefault.StaticString("7D"),
					},
					"suppress": schema.SingleNestedAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "Suppress settings for pattern actions.",
						Attributes: map[string]schema.Attribute{
							"seasonal_alerts": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								MarkdownDescription: "Whether to suppress seasonal alerts.",
							},
						},
					},
				},
			},
		},
	}
}

func buildFirstResponsePolicyRequest(plan FirstResponsePolicyModel) client.FirstResponsePolicy {
	policy := client.FirstResponsePolicy{
		Name:        plan.Name.ValueString(),
		EnabledMode: plan.EnabledMode.ValueString(),
		FilterQuery: plan.FilterQuery.ValueString(),
	}

	if !plan.AttributeActions.IsNull() && !plan.AttributeActions.IsUnknown() {
		attrActionAttrs := plan.AttributeActions.Attributes()
		actions := &client.FirstResponseAttrActions{
			Suppress: &client.FirstResponseAttrSuppress{},
			Insights: &client.FirstResponseInsights{},
		}
		if v, ok := attrActionAttrs["continuous_learning"]; ok && !v.IsNull() && !v.IsUnknown() {
			actions.ContinuousLearning = v.(types.Bool).ValueBool()
		}

		if v, ok := attrActionAttrs["suppress"]; ok && !v.IsNull() && !v.IsUnknown() {
			attrs := v.(types.Object).Attributes()
			if v, ok := attrs["learned_configuration"]; ok && !v.IsNull() {
				actions.Suppress.LearnedConfiguration = v.(types.Bool).ValueBool()
			}
			if v, ok := attrs["suppress_duration"]; ok && !v.IsNull() {
				actions.Suppress.SuppressDuration = int(v.(types.Int64).ValueInt64())
			}
		}
		if v, ok := attrActionAttrs["run_process"]; ok && !v.IsNull() && !v.IsUnknown() {
			runProc := &client.FirstResponseRunProcess{}
			attrs := v.(types.Object).Attributes()
			if v, ok := attrs["learned_configuration"]; ok && !v.IsNull() {
				runProc.LearnedConfiguration = v.(types.Bool).ValueBool()
			}
			if v, ok := attrs["run_immediately"]; ok && !v.IsNull() {
				runProc.RunImmediately = v.(types.Bool).ValueBool()
			}
			if v, ok := attrs["process_ids"]; ok && !v.IsNull() {
				listVal := v.(types.List)
				for _, elem := range listVal.Elements() {
					runProc.ProcessIds = append(runProc.ProcessIds, elem.(types.String).ValueString())
				}
			}
			actions.RunProcess = runProc
		}
		if v, ok := attrActionAttrs["insights"]; ok && !v.IsNull() && !v.IsUnknown() {
			attrs := v.(types.Object).Attributes()
			if v, ok := attrs["create_prc_insights"]; ok && !v.IsNull() {
				actions.Insights.CreatePrcInsights = v.(types.Bool).ValueBool()
			}
		}
		policy.AttributeActions = actions
	}

	if !plan.PatternActions.IsNull() && !plan.PatternActions.IsUnknown() {
		patternActionAttrs := plan.PatternActions.Attributes()
		pa := &client.FirstResponsePatternActions{
			Suppress: &client.FirstResponsePatternSuppress{},
		}
		if v, ok := patternActionAttrs["seasonality_time_frame"]; ok && !v.IsNull() && !v.IsUnknown() {
			pa.SeasonalityTimeFrame = v.(types.String).ValueString()
		}

		if v, ok := patternActionAttrs["suppress"]; ok && !v.IsNull() && !v.IsUnknown() {
			attrs := v.(types.Object).Attributes()
			if v, ok := attrs["seasonal_alerts"]; ok && !v.IsNull() {
				pa.Suppress.SeasonalAlerts = v.(types.Bool).ValueBool()
			}
		}
		policy.PatternActions = pa
	}

	return policy
}

func mapFirstResponsePolicyToState(resp *client.FirstResponsePolicy, state *FirstResponsePolicyModel) {
	state.Id = types.StringValue(resp.Id)
	state.Name = types.StringValue(resp.Name)
	state.EnabledMode = types.StringValue(resp.EnabledMode)
	state.FilterQuery = types.StringValue(resp.FilterQuery)
	state.AttributeActions = types.ObjectNull(attributeActionsAttrTypes)
	state.PatternActions = types.ObjectNull(patternActionsAttrTypes)

	if resp.AttributeActions != nil {
		// Suppress - always present
		suppress := resp.AttributeActions.Suppress
		if suppress == nil {
			suppress = &client.FirstResponseAttrSuppress{}
		}
		suppressObj, _ := types.ObjectValue(attrSuppressAttrTypes, map[string]attr.Value{
			"learned_configuration": types.BoolValue(suppress.LearnedConfiguration),
			"suppress_duration":     types.Int64Value(int64(suppress.SuppressDuration)),
		})

		// RunProcess - only present when configured
		var runProcessObj types.Object
		if resp.AttributeActions.RunProcess != nil {
			runProcess := resp.AttributeActions.RunProcess
			processIdValues := make([]attr.Value, len(runProcess.ProcessIds))
			for i, pid := range runProcess.ProcessIds {
				processIdValues[i] = types.StringValue(pid)
			}
			runProcessObj, _ = types.ObjectValue(attrRunProcessAttrTypes, map[string]attr.Value{
				"learned_configuration": types.BoolValue(runProcess.LearnedConfiguration),
				"run_immediately":       types.BoolValue(runProcess.RunImmediately),
				"process_ids":           types.ListValueMust(types.StringType, processIdValues),
			})
		} else {
			runProcessObj = types.ObjectNull(attrRunProcessAttrTypes)
		}

		// Insights - always present
		insights := resp.AttributeActions.Insights
		if insights == nil {
			insights = &client.FirstResponseInsights{}
		}
		insightsObj, _ := types.ObjectValue(insightsAttrTypes, map[string]attr.Value{
			"create_prc_insights": types.BoolValue(insights.CreatePrcInsights),
		})

		actionsObj, _ := types.ObjectValue(attributeActionsAttrTypes, map[string]attr.Value{
			"continuous_learning": types.BoolValue(resp.AttributeActions.ContinuousLearning),
			"suppress":            suppressObj,
			"insights":            insightsObj,
			"run_process":         runProcessObj,
		})

		state.AttributeActions = actionsObj
	}

	if resp.PatternActions != nil {
		seasonalityTimeFrame := types.StringNull()
		if resp.PatternActions.SeasonalityTimeFrame != "" {
			seasonalityTimeFrame = types.StringValue(resp.PatternActions.SeasonalityTimeFrame)
		}

		// Suppress - always present
		patternSuppress := resp.PatternActions.Suppress
		if patternSuppress == nil {
			patternSuppress = &client.FirstResponsePatternSuppress{}
		}
		suppressObj, _ := types.ObjectValue(patternSuppressAttrTypes, map[string]attr.Value{
			"seasonal_alerts": types.BoolValue(patternSuppress.SeasonalAlerts),
		})

		patternActionsObj, _ := types.ObjectValue(patternActionsAttrTypes, map[string]attr.Value{
			"seasonality_time_frame": seasonalityTimeFrame,
			"suppress":               suppressObj,
		})

		state.PatternActions = patternActionsObj
	}
}

// Create handles the creation of the resource.
func (r *FirstResponsePolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan FirstResponsePolicyModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.apiClient.TenantId
	if !plan.Client.IsNull() && plan.Client.ValueString() != "" {
		tenantId = plan.Client.ValueString()
	}

	policy := buildFirstResponsePolicyRequest(plan)

	created, err := r.apiClient.CreateFirstResponsePolicy(tenantId, policy)
	if err != nil {
		resp.Diagnostics.AddError("Creation Error", err.Error())
		return
	}

	mapFirstResponsePolicyToState(created, &plan)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Read handles reading the resource state.
func (r *FirstResponsePolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state FirstResponsePolicyModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.apiClient.TenantId
	if !state.Client.IsNull() && state.Client.ValueString() != "" {
		tenantId = state.Client.ValueString()
	}

	existing, err := r.apiClient.GetFirstResponsePolicy(tenantId, state.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read Error", err.Error())
		return
	}

	mapFirstResponsePolicyToState(existing, &state)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update handles updating the resource.
func (r *FirstResponsePolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state FirstResponsePolicyModel
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

	policy := buildFirstResponsePolicyRequest(plan)

	updated, err := r.apiClient.UpdateFirstResponsePolicy(tenantId, state.Id.ValueString(), policy)
	if err != nil {
		resp.Diagnostics.AddError("Update Error", err.Error())
		return
	}

	plan.Client = state.Client
	mapFirstResponsePolicyToState(updated, &plan)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Delete handles deleting the resource.
func (r *FirstResponsePolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state FirstResponsePolicyModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.apiClient.TenantId
	if !state.Client.IsNull() && state.Client.ValueString() != "" {
		tenantId = state.Client.ValueString()
	}

	err := r.apiClient.DeleteFirstResponsePolicy(tenantId, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

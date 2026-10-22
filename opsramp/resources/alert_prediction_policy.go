// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"
	"strings"

	"github.com/HPE/terraform-provider-hpe/opsramp/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure implementation satisfies the expected interfaces
var (
	_ resource.Resource               = &AlertPredictionPolicyResource{}
	_ resource.ResourceWithModifyPlan = &AlertPredictionPolicyResource{}
)

// AlertPredictionPolicyResource defines the resource implementation.
type AlertPredictionPolicyResource struct {
	BaseResource
}

// AlertPredictionPolicyModel maps Terraform schema attributes to the provider model.
type AlertPredictionPolicyModel struct {
	Client                   types.String `tfsdk:"client"`
	Id                       types.String `tfsdk:"id"`
	Name                     types.String `tfsdk:"name"`
	OrganizationMatchingType types.String `tfsdk:"organization_matching_type"`
	IncludedClients          types.List   `tfsdk:"included_clients"`
	EnabledMode              types.String `tfsdk:"enabled_mode"`
	SeasonalityTimeFrame     types.String `tfsdk:"seasonality_time_frame"`
	GeneratePredictionAlert  types.Bool   `tfsdk:"generate_prediction_alert"`
	FilterQuery              types.String `tfsdk:"filter_query"`
}

// NewAlertPredictionPolicy creates a new instance of the resource.
func NewAlertPredictionPolicy() resource.Resource {
	return &AlertPredictionPolicyResource{}
}

// Metadata returns the resource type name.
func (r *AlertPredictionPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alert_prediction_policy"
}

// ModifyPlan enforces the shared "client" attribute check plus organization_matching_type /
// included_clients scope rules.
func (r *AlertPredictionPolicyResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	r.BaseResource.ModifyPlan(ctx, req, resp)
	r.ValidateOrganizationMatching(ctx, req, resp)
}

// Schema defines the schema for the resource.
func (r *AlertPredictionPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an OpsRamp Alert Prediction Policy. Uses machine learning to predict future alerts based on historical data.",
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
				MarkdownDescription: "The unique identifier of the alert prediction policy.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the alert prediction policy.",
			},
			"enabled_mode": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The enabled mode of the policy. Valid values: `ON`, `OFF`.",
				Validators: []validator.String{
					stringvalidator.OneOf("ON", "OFF"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"seasonality_time_frame": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Learning based on the data for last N days. Valid values: `7D`, `10D`, `30D`, `60D`, `90D`.",
				Validators: []validator.String{
					stringvalidator.OneOf("7D", "10D", "30D", "60D", "90D"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"generate_prediction_alert": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether to generate predictions based on the policy.",
			},
			"filter_query": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Filter query to scope which alerts are evaluated by this policy.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"organization_matching_type": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The organization matching type for the policy. Applies only when the operation is performed at the MSP level; must not be set at the client level. Defaults to `ALL` when omitted at the MSP level. Valid values: `ALL`, `INCLUDE`.",
				Validators: []validator.String{
					stringvalidator.OneOf("ALL", "INCLUDE"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"included_clients": schema.ListAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "List of client names or UUIDs to include when `organization_matching_type` is `INCLUDE`. Required and must be non-empty in that case; must not be set at the client level.",
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func buildAlertPredictionPolicyRequest(plan AlertPredictionPolicyModel, apiClient *client.OpsRampClient) client.AlertPredictionPolicy {
	policy := client.AlertPredictionPolicy{
		Name:                 plan.Name.ValueString(),
		EnabledMode:          plan.EnabledMode.ValueString(),
		SeasonalityTimeFrame: plan.SeasonalityTimeFrame.ValueString(),
		FilterQuery:          plan.FilterQuery.ValueString(),
	}

	isMSP := apiClient.Scope == "MSP" && (plan.Client.IsNull() || plan.Client.ValueString() == "")
	if isMSP {
		if !plan.OrganizationMatchingType.IsNull() && !plan.OrganizationMatchingType.IsUnknown() {
			policy.OrganizationMatchingType = plan.OrganizationMatchingType.ValueString()
		}

		if !plan.IncludedClients.IsNull() && !plan.IncludedClients.IsUnknown() {
			var includedClients []types.String
			_ = plan.IncludedClients.ElementsAs(context.Background(), &includedClients, false)
			for _, c := range includedClients {
				policy.IncludedClients = append(policy.IncludedClients, c.ValueString())
			}
		}
	}

	return policy
}

func mapAlertPredictionPolicyToState(
	resp *client.AlertPredictionPolicy, state *AlertPredictionPolicyModel, isMSP bool, fromRead bool,
) {
	state.Id = types.StringValue(resp.Id)
	state.Name = types.StringValue(resp.Name)
	state.EnabledMode = types.StringValue(resp.EnabledMode)
	state.SeasonalityTimeFrame = types.StringValue(resp.SeasonalityTimeFrame)
	state.FilterQuery = types.StringValue(resp.FilterQuery)

	// The API reports MSP scope via either a populated "includedClients" list (INCLUDE) or
	// "clientsIncluded":"ALL" (ALL); it never echoes back "organizationMatchingType".
	applyIncludedClientsTargeting(
		resp.IncludedClients, resp.ClientsIncluded, isMSP, fromRead, &state.OrganizationMatchingType, &state.IncludedClients,
	)
}

// Create handles the creation of the resource.
func (r *AlertPredictionPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AlertPredictionPolicyModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.apiClient.TenantId
	if !plan.Client.IsNull() && plan.Client.ValueString() != "" {
		tenantId = plan.Client.ValueString()
	}

	policy := buildAlertPredictionPolicyRequest(plan, r.apiClient)

	created, err := r.apiClient.CreateAlertPredictionPolicy(tenantId, policy)
	if err != nil {
		resp.Diagnostics.AddError("Creation Error", err.Error())

		return
	}

	isMSP := r.apiClient.Scope == "MSP" && (plan.Client.IsNull() || plan.Client.ValueString() == "")
	mapAlertPredictionPolicyToState(created, &plan, isMSP, false)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Read handles reading the resource state.
func (r *AlertPredictionPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AlertPredictionPolicyModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.apiClient.TenantId
	if !state.Client.IsNull() && state.Client.ValueString() != "" {
		tenantId = state.Client.ValueString()
	}

	existing, err := r.apiClient.GetAlertPredictionPolicy(tenantId, state.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)

			return
		}
		resp.Diagnostics.AddError("Read Error", err.Error())

		return
	}

	mapAlertPredictionPolicyToState(
		existing, &state, r.apiClient.Scope == "MSP" && (state.Client.IsNull() || state.Client.ValueString() == ""), true,
	)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update handles updating the resource.
func (r *AlertPredictionPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state AlertPredictionPolicyModel
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

	policy := buildAlertPredictionPolicyRequest(plan, r.apiClient)

	updated, err := r.apiClient.UpdateAlertPredictionPolicy(tenantId, state.Id.ValueString(), policy)
	if err != nil {
		resp.Diagnostics.AddError("Update Error", err.Error())

		return
	}

	plan.Client = state.Client
	isMSP := r.apiClient.Scope == "MSP" && (plan.Client.IsNull() || plan.Client.ValueString() == "")
	mapAlertPredictionPolicyToState(updated, &plan, isMSP, false)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Delete handles deleting the resource.
func (r *AlertPredictionPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AlertPredictionPolicyModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.apiClient.TenantId
	if !state.Client.IsNull() && state.Client.ValueString() != "" {
		tenantId = state.Client.ValueString()
	}

	err := r.apiClient.DeleteAlertPredictionPolicy(tenantId, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", err.Error())

		return
	}

	resp.State.RemoveResource(ctx)
}

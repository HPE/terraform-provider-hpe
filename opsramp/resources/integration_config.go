// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/HPE/terraform-provider-hpe/opsramp/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure implementation satisfies the expected interfaces
var _ resource.Resource = &IntegrationConfigResource{}
var _ resource.ResourceWithImportState = &IntegrationConfigResource{}

// IntegrationConfigResource defines the resource implementation.
type IntegrationConfigResource struct {
	BaseResource
}

// IntegrationConfigModel maps Terraform schema attributes to the provider model.
type IntegrationConfigModel struct {
	Id             types.String   `tfsdk:"id"`
	IntegrationId  types.String   `tfsdk:"integration_id"`
	Client         types.String   `tfsdk:"client"`
	Name           types.String   `tfsdk:"name"`
	Config         types.String   `tfsdk:"config"`
	State          types.String   `tfsdk:"state"`
	AllResources   types.Bool     `tfsdk:"all_resources"`
	Schedule       *ScheduleModel `tfsdk:"schedule"`
	ProxyProfileId types.String   `tfsdk:"proxy_profile_id"`
}

type ScheduleModel struct {
	PatternType types.String `tfsdk:"pattern_type"`
	Pattern     types.Int64  `tfsdk:"pattern"`
	StartTime   types.String `tfsdk:"start_time"`
}

// NewIntegrationConfig creates a new instance of the resource.
func NewIntegrationConfig() resource.Resource {
	return &IntegrationConfigResource{}
}

// Metadata returns the resource type name.
func (r *IntegrationConfigResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integration_config"
}

// Schema defines the schema for the resource.
func (r *IntegrationConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an OpsRamp Integration Config. This is a sub-resource of an installed integration that represents a specific configuration (e.g. a Kubernetes cluster config under a K8s integration).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the config (e.g. ADAPTER-MANIFEST-...).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"integration_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the installed integration this config belongs to (e.g. INTG-...).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"client": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The unique ID of the client (sub-tenant) where the integration is installed. If not provided, uses the provider tenant.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the configuration.",
			},
			"config": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The configuration as a JSON-encoded string.",
			},
			"state": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The current state of the configuration (e.g. ENABLED).",
			},
			"all_resources": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether this config applies to all resources discovered by the integration.",
			},
			"proxy_profile_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The UUID of the Gateway profile to use for proxying collected data.",
			},
			"schedule": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Schedule for integration execution. Applicable for configuration-based integrations that run on a schedule (e.g. `PROMETHEUSREMOTEWRITE`). Not applicable for event-based integrations.",
				Attributes: map[string]schema.Attribute{
					"pattern_type": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Operator type (e.g. `HOURLY`, `DAILY`, `WEEKLY`, `MONTHLY`).",
						Validators: []validator.String{
							stringvalidator.OneOf("HOURLY", "DAILY", "WEEKLY", "MONTHLY"),
						},
					},
					"pattern": schema.Int64Attribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The schedule pattern. Interpretation depends on pattern_type (e.g. for HOURLY, 2 means every 2 hours; for WEEKLY, 1 means every Monday, 2 means every Tuesday, etc.).",
						Validators: []validator.Int64{
							int64validator.AtLeast(1),
						},
						Default: int64default.StaticInt64(1),
					},
					"start_time": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The schedule start time in RFC3339 format. If not specified, defaults to the time of integration creation.",
						Default:             stringdefault.StaticString(""),
					},
				},
			},
		},
	}
}

// getTenantId determines which tenant ID to use based on the optional client parameter
func (r *IntegrationConfigResource) getTenantId(clientId types.String) string {
	if !clientId.IsNull() && clientId.ValueString() != "" {
		return clientId.ValueString()
	}
	return r.apiClient.TenantId
}

// Create handles the creation of the resource.
func (r *IntegrationConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan IntegrationConfigModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.getTenantId(plan.Client)
	integrationId := plan.IntegrationId.ValueString()

	// Parse config JSON
	var configMap map[string]any
	if err := json.Unmarshal([]byte(plan.Config.ValueString()), &configMap); err != nil {
		resp.Diagnostics.AddError("Invalid Config JSON", fmt.Sprintf("config must be valid JSON: %s", err.Error()))
		return
	}

	var schedule *client.Schedule
	scheduleNone := true

	if plan.Schedule != nil {
		scheduleNone = false
		schedule = &client.Schedule{
			PatternType: plan.Schedule.PatternType.ValueString(),
			Pattern:     plan.Schedule.Pattern.ValueInt64(),
			StartTime:   plan.Schedule.StartTime.ValueString(),
		}
	}

	createReq := client.IntegrationConfigRequest{
		Name:         plan.Name.ValueString(),
		Config:       configMap,
		ScheduleNone: scheduleNone,
		Schedule:     schedule,
		AllResources: plan.AllResources.ValueBool(),
	}

	if plan.ProxyProfileId.ValueString() != "" {
		createReq.InfoMap = client.InfoMap{KubernetesProxy: &client.ProxyRef{Uuid: plan.ProxyProfileId.ValueString()}}
	}

	created, err := r.apiClient.CreateIntegrationConfig(tenantId, integrationId, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Integration Config Create Error",
			fmt.Sprintf("Could not create config for integration '%s': %s", integrationId, err.Error()))
		return
	}

	plan.Id = types.StringValue(created.ID)
	plan.State = types.StringValue(created.State)

	// Normalize config from response
	if created.Config != nil {
		configBytes, _ := json.Marshal(created.Config)
		plan.Config = types.StringValue(string(configBytes))
	}

	if created.InfoMap != nil && created.InfoMap.KubernetesProxy != nil {
		plan.ProxyProfileId = types.StringValue(created.InfoMap.KubernetesProxy.Uuid)
	} else {
		plan.ProxyProfileId = types.StringNull()
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Read handles reading the resource state.
func (r *IntegrationConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state IntegrationConfigModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.getTenantId(state.Client)

	existing, err := r.apiClient.GetIntegrationConfig(tenantId, state.IntegrationId.ValueString(), state.Id.ValueString())
	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "404") || strings.Contains(errStr, "not found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read Error", err.Error())
		return
	}

	state.Name = types.StringValue(existing.Name)
	state.State = types.StringValue(existing.State)

	if existing.Config != nil {
		configBytes, _ := json.Marshal(existing.Config)
		state.Config = types.StringValue(string(configBytes))
	}

	if existing.InfoMap != nil && existing.InfoMap.KubernetesProxy != nil {
		state.ProxyProfileId = types.StringValue(existing.InfoMap.KubernetesProxy.Uuid)
	}

	if existing.InfoMap != nil && existing.InfoMap.KubernetesProxy != nil {
		state.ProxyProfileId = types.StringValue(existing.InfoMap.KubernetesProxy.Uuid)
	} else {
		state.ProxyProfileId = types.StringNull()
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update handles updating the resource.
func (r *IntegrationConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan IntegrationConfigModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state IntegrationConfigModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.getTenantId(plan.Client)
	integrationId := plan.IntegrationId.ValueString()
	configId := state.Id.ValueString()

	// Parse config JSON
	var configMap map[string]any
	if err := json.Unmarshal([]byte(plan.Config.ValueString()), &configMap); err != nil {
		resp.Diagnostics.AddError("Invalid Config JSON", fmt.Sprintf("config must be valid JSON: %s", err.Error()))
		return
	}

	updateReq := client.IntegrationConfigRequest{
		Name:         plan.Name.ValueString(),
		Config:       configMap,
		AllResources: plan.AllResources.ValueBool(),
	}

	if plan.Schedule != nil {
		updateReq.ScheduleNone = true
		updateReq.Schedule = &client.Schedule{
			PatternType: plan.Schedule.PatternType.ValueString(),
			Pattern:     plan.Schedule.Pattern.ValueInt64(),
			StartTime:   plan.Schedule.StartTime.ValueString(),
		}
	}

	if plan.ProxyProfileId.ValueString() != "" {
		updateReq.InfoMap = client.InfoMap{KubernetesProxy: &client.ProxyRef{Uuid: plan.ProxyProfileId.ValueString()}}
	}

	updated, err := r.apiClient.UpdateIntegrationConfig(tenantId, integrationId, configId, updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Integration Config Update Error", err.Error())
		return
	}

	plan.Id = state.Id
	plan.State = types.StringValue(updated.State)

	if updated.Config != nil {
		configBytes, _ := json.Marshal(updated.Config)
		plan.Config = types.StringValue(string(configBytes))
	}

	if updated.InfoMap != nil && updated.InfoMap.KubernetesProxy != nil {
		plan.ProxyProfileId = types.StringValue(updated.InfoMap.KubernetesProxy.Uuid)
	} else {
		plan.ProxyProfileId = types.StringNull()
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Delete handles deleting the resource.
func (r *IntegrationConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state IntegrationConfigModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.getTenantId(state.Client)

	err := r.apiClient.DeleteIntegrationConfig(tenantId, state.IntegrationId.ValueString(), state.Id.ValueString())
	if err != nil {
		if !strings.Contains(err.Error(), "404") && !strings.Contains(err.Error(), "not found") {
			resp.Diagnostics.AddError("Delete Error", err.Error())
			return
		}
	}
}

// ImportState handles resource import.
// Import ID format: integration_id:config_id or client_id:integration_id:config_id
func (r *IntegrationConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importId := req.ID

	var state IntegrationConfigModel

	parts := strings.Split(importId, ":")
	switch len(parts) {
	case 2:
		// Format: integration_id:config_id
		state.IntegrationId = types.StringValue(parts[0])
		state.Id = types.StringValue(parts[1])
		state.Client = types.StringNull()
	case 3:
		// Format: client_id:integration_id:config_id
		state.Client = types.StringValue(parts[0])
		state.IntegrationId = types.StringValue(parts[1])
		state.Id = types.StringValue(parts[2])
	default:
		resp.Diagnostics.AddError("Invalid Import ID",
			"Import ID must be in the format 'integration_id:config_id' or 'client_id:integration_id:config_id'.")
		return
	}

	state.Name = types.StringUnknown()
	state.Config = types.StringUnknown()
	state.State = types.StringUnknown()

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

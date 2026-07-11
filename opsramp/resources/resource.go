// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package resources

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/HPE/terraform-provider-hpe/opsramp/client" // Adjust import path to match your client package

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure implementation satisfies the expected interfaces
var _ resource.Resource = &Resource{}
var _ resource.ResourceWithImportState = &Resource{}
var _ resource.ResourceWithModifyPlan = &Resource{}

// Resource defines the resource implementation.
type Resource struct {
	BaseResource
}

// New creates a new instance of the resource.
func NewResource() resource.Resource {
	return &Resource{}
}

// Metadata returns the resource type name.
func (r *Resource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource"
}

// Schema defines the schema for the resource.
func (r *Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"client": schema.StringAttribute{
				Optional:    true,
				Description: "The tenant/client ID to use for this resource. Defaults to the provider's tenant ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"uuid": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The UUID used to identify an existing resource. Provide either uuid or hostname / resource_name together with resource_type.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"resource_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The resource name. When uuid and hostname are not set, resource_name must be provided together with resource_type.",
				Default:     stringdefault.StaticString(""),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"resource_type": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("Server"),
				MarkdownDescription: "The type of the resource. Required with resource_name and hostname when uuid is not set.",
			},
			"hostname": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The hostname. When uuid and resource_name are not set, resource_name must be provided together with resource_type.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"alias_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The resource alias.",
				Default:     stringdefault.StaticString(""),
			},
			"agent_installed": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "If the resource has an agent installed.",
			},
		},
	}
}

// Create handles the creation of the resource.
func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !hasValidResourceIdentifier(plan, ResourceModel{}) {
		resp.Diagnostics.AddError(
			"Missing required identifier",
			"Provide either uuid or hostname / resource_name together with resource_type.",
		)
		return
	}

	tenantId := r.apiClient.TenantId
	if !plan.Client.IsNull() && plan.Client.ValueString() != "" {
		tenantId = plan.Client.ValueString()
	}

	var uuid string

	// Check if UUID is provided
	if plan.Uuid.ValueString() != "" {

		uuid = plan.Uuid.ValueString()
		existing, err := r.apiClient.GetResource(tenantId, uuid)
		if err != nil {
			resp.Diagnostics.AddError("Backend query error", err.Error())
			return
		}

		_, err = r.PerformUpdate(tenantId, plan, uuid)

		if err != nil {
			resp.Diagnostics.AddError("Update Error", err.Error())
			return
		}

		state := ResourceModel{
			Client:         plan.Client,
			Uuid:           types.StringValue(existing.Uuid),
			ResourceName:   types.StringValue(existing.GeneralInfo.ResourceName),
			HostName:       types.StringValue(existing.GeneralInfo.HostName),
			ResourceType:   types.StringValue(plan.ResourceType.ValueString()),
			AliasName:      types.StringValue(existing.GeneralInfo.AliasName),
			AgentInstalled: types.BoolValue(existing.AgentInstalled),
		}

		// Save state
		resp.State.Set(ctx, &state)

	} else {
		// Create the resource
		createResource := client.CreateResource{
			ResourceName: plan.ResourceName.ValueString(),
			ResourceType: plan.ResourceType.ValueString(),
			HostName:     plan.HostName.ValueString(),
			AliasName:    plan.AliasName.ValueString(),
		}

		created, err := r.apiClient.CreateResource(tenantId, createResource)
		if err != nil {
			resp.Diagnostics.AddError("Creation Error", err.Error())
			return
		}
		uuid = created.ResourceUuid

		plan.Uuid = types.StringValue(uuid)
		plan.AgentInstalled = types.BoolValue(false)
		plan.ResourceName = types.StringValue(plan.ResourceName.ValueString())
		plan.HostName = types.StringValue(plan.HostName.ValueString())
		plan.AliasName = types.StringValue(plan.AliasName.ValueString())

		// Save state
		resp.State.Set(ctx, &plan)
	}
}

// --- Read ---
func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.apiClient.TenantId
	if !state.Client.IsNull() && state.Client.ValueString() != "" {
		tenantId = state.Client.ValueString()
	}

	existing, err := r.apiClient.GetResource(tenantId, state.Uuid.ValueString())
	// TODO: handle http error codes
	if err != nil && !strings.Contains(err.Error(), "No Resource found") {
		resp.Diagnostics.AddError("Read Error", err.Error())
		return
	}

	if existing == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state = ResourceModel{
		Client:         state.Client,
		Uuid:           types.StringValue(existing.Uuid),
		ResourceName:   types.StringValue(existing.GeneralInfo.ResourceName),
		ResourceType:   types.StringValue(existing.GeneralInfo.ResourceType),
		HostName:       types.StringValue(existing.GeneralInfo.HostName),
		AliasName:      types.StringValue(existing.GeneralInfo.AliasName),
		AgentInstalled: types.BoolValue(existing.AgentInstalled),
	}

	resp.State.Set(ctx, &state)
}

func (r *Resource) PerformUpdate(tenantId string, plan ResourceModel, uuid string) (interface{}, error) {
	var updateResource client.UpdateResource

	if plan.ResourceType.ValueString() != "" {
		updateResource.ResourceType = plan.ResourceType.ValueString()
	}
	if plan.AliasName.ValueString() != "" {
		updateResource.AliasName = plan.AliasName.ValueString()
	}

	result, err := r.apiClient.UpdateResource(tenantId, uuid, updateResource)
	return result, err
}

// Update modifies an existing resource.
func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// UUID cannot be changed; ensure UUID stays the same.
	if plan.Uuid.ValueString() != "" && plan.Uuid.ValueString() != state.Uuid.ValueString() {
		resp.Diagnostics.AddError("UUID Immutable", "UUID cannot be changed after creation.")
		return
	}

	tenantId := r.apiClient.TenantId
	if !state.Client.IsNull() && state.Client.ValueString() != "" {
		tenantId = state.Client.ValueString()
	}

	_, err := r.PerformUpdate(tenantId, plan, state.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Resource",
			fmt.Sprintf("Could not update resource: %s", err),
		)
		return
	}

	// Computed/immutable attributes
	plan.Client = state.Client
	plan.Uuid = state.Uuid
	plan.AgentInstalled = state.AgentInstalled

	// Save
	resp.State.Set(ctx, &plan)
}

// Delete removes the resource from the API.
func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.apiClient.TenantId
	if !state.Client.IsNull() && state.Client.ValueString() != "" {
		tenantId = state.Client.ValueString()
	}

	_, err := r.apiClient.DeleteResource(tenantId, state.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

// ImportState handles importing an existing resource.
func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	uuid := req.ID
	tenantId := r.apiClient.TenantId
	res, err := r.apiClient.GetResource(tenantId, uuid)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Importing Resource",
			fmt.Sprintf("Could not import resource: %s", err),
		)
		return
	}

	state := ResourceModel{
		Client:       types.StringNull(),
		Uuid:         types.StringValue(res.Uuid),
		ResourceName: types.StringValue(res.GeneralInfo.ResourceName),
		ResourceType: types.StringValue(res.GeneralInfo.ResourceType),
		HostName:     types.StringValue(res.GeneralInfo.HostName),
		AliasName:    types.StringValue(res.GeneralInfo.AliasName),
	}

	resp.State.Set(ctx, &state)
}

// ResourceModel maps Terraform schema attributes to the provider model.
type ResourceModel struct {
	Client         types.String `tfsdk:"client"`
	Uuid           types.String `tfsdk:"uuid"`
	ResourceName   types.String `tfsdk:"resource_name"`
	HostName       types.String `tfsdk:"hostname"`
	ResourceType   types.String `tfsdk:"resource_type"`
	AliasName      types.String `tfsdk:"alias_name"`
	AgentInstalled types.Bool   `tfsdk:"agent_installed"`
}

func (r *Resource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Call base implementation
	r.BaseResource.ModifyPlan(ctx, req, resp)

	// Don't modify plan during destroy
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan, state ResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	req.State.Get(ctx, &state)
	if resp.Diagnostics.HasError() {
		return
	}

	hasClientOverride := !plan.Client.IsNull() && (plan.Client.IsUnknown() || strings.TrimSpace(plan.Client.ValueString()) != "")
	if strings.ToUpper(r.apiClient.Scope) != "CLIENT" && !hasClientOverride {
		resp.Diagnostics.AddError("Resources can only be created at Client level", "Use a client-scoped provider configuration or specify the client using unique ID.")
		return
	}

	// Use client parameter if set, otherwise provider's tenant ID
	tenantId := r.apiClient.TenantId
	if !plan.Client.IsNull() && plan.Client.ValueString() != "" {
		tenantId = plan.Client.ValueString()
	}

	// Validate Device Type
	if plan.ResourceType.ValueString() != "" {
		deviceTypes, err := r.apiClient.GetResourceTypes(tenantId)
		if err != nil {
			resp.Diagnostics.AddError("Error fetching GetDeviceTypes", err.Error())
			return
		}

		deviceType := plan.ResourceType.ValueString()
		if !slices.Contains(deviceTypes, deviceType) {
			resp.Diagnostics.AddError("Error Device Type does not macth API device types", fmt.Sprintf("Allowed types: %s", strings.Join(deviceTypes, ",")))
			return
		}
	}

	if !hasValidResourceIdentifier(plan, state) {
		resp.Diagnostics.AddError("Missing required identifier", "Provide either uuid, hostname, or resource_name together with resource_type.")
	}

}

func hasValidResourceIdentifier(plan ResourceModel, state ResourceModel) bool {
	uuid := firstNonEmpty(plan.Uuid.ValueString(), state.Uuid.ValueString())
	hostname := firstNonEmpty(plan.HostName.ValueString(), state.HostName.ValueString())
	resourceName := firstNonEmpty(plan.ResourceName.ValueString(), state.ResourceName.ValueString())
	resourceType := firstNonEmpty(plan.ResourceType.ValueString(), state.ResourceType.ValueString())

	return uuid != "" || ((resourceName != "" || hostname != "") && resourceType != "")
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}

	return ""
}

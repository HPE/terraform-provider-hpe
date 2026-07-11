// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/HPE/terraform-provider-hpe/opsramp/client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &ManagementProfileResource{}
var _ resource.ResourceWithModifyPlan = &ManagementProfileResource{}

// ManagementProfileResource defines the resource implementation.
type ManagementProfileResource struct {
	BaseResource
}

// ManagementProfileModel maps Terraform schema attributes to the provider model.
type ManagementProfileModel struct {
	Client      types.String `tfsdk:"client"`
	Uuid        types.String `tfsdk:"uuid"`
	Id          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Type        types.String `tfsdk:"type"`
}

// NewManagementProfile creates a new instance of the resource.
func NewManagementProfile() resource.Resource {
	return &ManagementProfileResource{}
}

func (r *ManagementProfileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_management_profile"
}

func (r *ManagementProfileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an OpsRamp Management Profile (Gateway type). Management profiles define the gateway agent configuration used when installing managed resources.",
		Attributes: map[string]schema.Attribute{
			"client": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The client (sub-tenant) UUID where the profile should be created. If omitted, the provider's configured tenant is used.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"uuid": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The numeric identifier of the management profile assigned by OpsRamp.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the management profile assigned by OpsRamp.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The management profile name.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "Summary describing the management profile.",
			},
			"type": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("Gateway"),
				MarkdownDescription: "The management profile type. Always `Gateway`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *ManagementProfileResource) resolveTenantId(clientAttr types.String) string {
	if !clientAttr.IsNull() && clientAttr.ValueString() != "" {
		return clientAttr.ValueString()
	}
	return r.apiClient.TenantId
}

// Create handles the creation of the resource.
func (r *ManagementProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ManagementProfileModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.resolveTenantId(plan.Client)

	createReq := client.ManagementProfile{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Type:        plan.Type.ValueString(),
	}

	created, err := r.apiClient.CreateManagementProfile(tenantId, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Management Profile Create Error",
			fmt.Sprintf("Could not create management profile: %s", err))
		return
	}

	plan.Uuid = types.StringValue(created.Uuid)
	plan.Id = types.Int64Value(int64(created.Id))
	plan.Name = types.StringValue(created.Name)
	plan.Description = types.StringValue(created.Description)
	plan.Type = types.StringValue(created.Type)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Read handles reading the resource state.
func (r *ManagementProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ManagementProfileModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.resolveTenantId(state.Client)

	profile, err := r.apiClient.GetManagementProfile(tenantId, int(state.Id.ValueInt64()))
	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "No collector profile") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Management Profile Read Error",
			fmt.Sprintf("Could not read management profile %s: %s", state.Uuid.ValueString(), err))
		return
	}

	state.Name = types.StringValue(profile.Name)
	state.Description = types.StringValue(profile.Description)
	state.Type = types.StringValue(profile.Type)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update handles updating the resource.
func (r *ManagementProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ManagementProfileModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state ManagementProfileModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.resolveTenantId(plan.Client)

	updateReq := client.ManagementProfile{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Type:        plan.Type.ValueString(),
	}

	updated, err := r.apiClient.UpdateManagementProfile(tenantId, int(state.Id.ValueInt64()), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Management Profile Update Error",
			fmt.Sprintf("Could not update management profile %s: %s", state.Uuid.ValueString(), err))
		return
	}

	plan.Uuid = state.Uuid
	plan.Id = state.Id
	plan.Name = types.StringValue(updated.Name)
	plan.Description = types.StringValue(updated.Description)
	plan.Type = types.StringValue(updated.Type)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Delete handles deleting the resource.
func (r *ManagementProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ManagementProfileModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.resolveTenantId(state.Client)

	if err := r.apiClient.DeleteManagementProfile(tenantId, int(state.Id.ValueInt64())); err != nil {
		resp.Diagnostics.AddError("Management Profile Delete Error",
			fmt.Sprintf("Could not delete management profile %s: %s", state.Uuid.ValueString(), err))
	}
}

func (r *ManagementProfileResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Call base implementation
	r.BaseResource.ModifyPlan(ctx, req, resp)

	// Don't modify plan during destroy
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan, state ManagementProfileModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	req.State.Get(ctx, &state)
	if resp.Diagnostics.HasError() {
		return
	}

	hasClientOverride := !plan.Client.IsNull() && (plan.Client.IsUnknown() || strings.TrimSpace(plan.Client.ValueString()) != "")
	if strings.ToUpper(r.apiClient.Scope) != "CLIENT" && !hasClientOverride {
		resp.Diagnostics.AddError("Management Profiles can only be created at Client level", "Use a client-scoped provider configuration or specify the client using unique ID.")
		return
	}
}

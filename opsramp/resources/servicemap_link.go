// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package resources

import (
	"context"
	"strings"

	"github.com/HPE/terraform-provider-hpe/opsramp/client" // Adjust import path to match your client package

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure implementation satisfies the expected interfaces
var _ resource.Resource = &ServicemapResourceLink{}
var _ resource.ResourceWithImportState = &ServicemapResourceLink{}
var _ resource.ResourceWithModifyPlan = &ServicemapResourceLink{}

// ServicemapResourceLink defines the resource implementation.
type ServicemapResourceLink struct {
	BaseResource
}

// NewServicemapLink creates a new instance of the resource.
func NewServicemapLink() resource.Resource {
	return &ServicemapResourceLink{}
}

// Metadata returns the resource type name.
func (r *ServicemapResourceLink) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_servicemap_link"
}

// Recursive schema definition example
func (r *ServicemapResourceLink) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"client": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			// parent attribute
			"parent": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			// link attribute
			"link": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func translateLinkPlanToModel(plan ServicemapLinkModel) client.CreateServicemapLink {
	var parent *client.Parent

	if plan.Parent.ValueString() != "" {
		parent = &client.Parent{
			Id: plan.Parent.ValueString(),
		}
	}

	return client.CreateServicemapLink{
		Parent: parent,
		Id:     plan.Link.ValueString(),
	}
}

func (r *ServicemapResourceLink) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ServicemapLinkModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.apiClient.TenantId
	if !plan.Client.IsNull() && plan.Client.ValueString() != "" {
		tenantId = plan.Client.ValueString()
	}

	// Use the JSON from the plan
	CreateServicemapLink := translateLinkPlanToModel(plan)

	// Create the Servicemap in the backend
	_, err := r.apiClient.CreateServicemapLink(tenantId, CreateServicemapLink)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", err.Error())
		return
	}

	// Save new state back to Terraform
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// --- Read ---
func (r *ServicemapResourceLink) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ServicemapLinkModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.apiClient.TenantId
	if !state.Client.IsNull() && state.Client.ValueString() != "" {
		tenantId = state.Client.ValueString()
	}

	serviceMapLink := client.CreateServicemapLink{
		Id: state.Link.ValueString(),
		Parent: &client.Parent{
			Id: state.Parent.ValueString(),
		},
	}

	backendServicemapLink, err := r.apiClient.GetServicemapLink(tenantId, serviceMapLink)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", err.Error())
		return
	}

	if backendServicemapLink == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update modifies an existing resource.
func (r *ServicemapResourceLink) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ServicemapLinkModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save
	resp.State.Set(ctx, &state)
}

// Delete removes the resource from the API.
func (r *ServicemapResourceLink) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ServicemapLinkModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.apiClient.TenantId
	if !state.Client.IsNull() && state.Client.ValueString() != "" {
		tenantId = state.Client.ValueString()
	}

	servicemapLink := translateLinkPlanToModel(state)

	err := r.apiClient.DeleteServicemapLink(tenantId, servicemapLink)
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

// ImportState handles importing an existing servicemap link.
// Import ID format: "{parent_id}/{link_id}"
func (r *ServicemapResourceLink) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			"Expected format: '{parent_id}/{link_id}'",
		)
		return
	}

	state := ServicemapLinkModel{
		Client: types.StringNull(),
		Parent: types.StringValue(parts[0]),
		Link:   types.StringValue(parts[1]),
	}

	resp.State.Set(ctx, &state)
}

type ServicemapLinkModel struct {
	Client types.String `tfsdk:"client"`
	Parent types.String `tfsdk:"parent"`
	Link   types.String `tfsdk:"link"`
}

func (r *ServicemapResourceLink) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Call base implementation
	r.BaseResource.ModifyPlan(ctx, req, resp)

	// Don't modify plan during destroy
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan, state ServicemapLinkModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	req.State.Get(ctx, &state)
	if resp.Diagnostics.HasError() {
		return
	}

	hasClientOverride := !plan.Client.IsNull() && (plan.Client.IsUnknown() || strings.TrimSpace(plan.Client.ValueString()) != "")
	if strings.ToUpper(r.apiClient.Scope) != "CLIENT" && !hasClientOverride {
		resp.Diagnostics.AddError("ServiceMap Links can only be created at Client level", "Use a client-scoped provider configuration or specify the client using unique ID.")
		return
	}
}

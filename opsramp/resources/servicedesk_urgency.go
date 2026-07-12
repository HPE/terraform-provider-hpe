// (C) Copyright 2026 Hewlett Packard Enterprise Development LP
package resources

import (
	"context"

	"github.com/HPE/terraform-provider-hpe/opsramp/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ServiceDeskUrgencyModel struct {
	Client      types.String `tfsdk:"client"`
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	State       types.Bool   `tfsdk:"state"`
}

// ServiceDeskUrgency defines the resource implementation.
type ServiceDeskUrgency struct {
	BaseResource
}

// Ensure implementation satisfies the expected interfaces
var _ resource.Resource = &ServiceDeskUrgency{}
var _ resource.ResourceWithModifyPlan = &ServiceDeskUrgency{}

//var _ resource.ResourceWithImportState = &ServiceDeskUrgency{}

// New creates a new instance of the resource.
func NewServiceDeskUrgency() resource.Resource {
	return &ServiceDeskUrgency{}
}

// Metadata returns the resource type name.
func (r *ServiceDeskUrgency) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_servicedesk_urgency"
}

// Schema defines the schema for the resource.
func (r *ServiceDeskUrgency) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"client": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The client (tenant) UUID where this resource should be managed. If not specified, uses the provider tenant.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The ID of the ServiceDeskUrgency. May be retrieved from the backend.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the ServiceDeskUrgency.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The description of the ServiceDeskUrgency.",
			},
			"state": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
		},
	}
}

func (r *ServiceDeskUrgency) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ServiceDeskUrgencyModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	urgency := client.ServiceDeskUrgency{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		State:       plan.State.ValueBool(),
	}

	created, err := r.apiClient.CreateServiceDeskUrgency(urgency)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", err.Error())
		return
	}

	plan.Id = types.StringValue(created.Id)
	plan.State = types.BoolValue(created.State)
	plan.Name = types.StringValue(created.Name)
	plan.Description = types.StringValue(created.Description)
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ServiceDeskUrgency) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ServiceDeskUrgencyModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use ticket type from state (required for API)
	urgency, err := r.apiClient.GetServiceDeskUrgency(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read Error", err.Error())
		return
	}

	// If the urgency is not found, remove the resource from state
	// This is a common pattern in Terraform providers to handle resources that may not exist anymore.
	if urgency == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(urgency.Name)
	state.Description = types.StringValue(urgency.Description)
	state.State = types.BoolValue(urgency.State)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ServiceDeskUrgency) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ServiceDeskUrgencyModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	urgency := client.ServiceDeskUrgency{
		Id:          plan.Id.ValueString(),
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		State:       plan.State.ValueBool(),
	}

	// The API may not support update directly; if not, implement as needed.
	// For now, try create (upsert) pattern:
	updated, err := r.apiClient.UpdateServiceDeskUrgency(plan.Id.ValueString(), urgency)
	if err != nil {
		resp.Diagnostics.AddError("Update Error", err.Error())
		return
	}

	plan.Id = types.StringValue(updated.Id)
	plan.Name = types.StringValue(updated.Name)
	plan.Description = types.StringValue(updated.Description)
	plan.State = types.BoolValue(updated.State)
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ServiceDeskUrgency) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ServiceDeskUrgencyModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.apiClient.DeleteServiceDeskUrgency(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", err.Error())
	}

	resp.State.RemoveResource(ctx)
}

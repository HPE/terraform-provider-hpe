// (C) Copyright 2026 Hewlett Packard Enterprise Development LP
package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/opsramp/client"
)

type ServiceDeskBusinessImpactModel struct {
	Client      types.String `tfsdk:"client"`
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	State       types.Bool   `tfsdk:"state"`
}

// ServiceDeskBusinessImpact defines the resource implementation.
type ServiceDeskBusinessImpact struct {
	BaseResource
}

// Ensure implementation satisfies the expected interfaces
var (
	_ resource.Resource               = &ServiceDeskBusinessImpact{}
	_ resource.ResourceWithModifyPlan = &ServiceDeskBusinessImpact{}
)

// var _ resource.ResourceWithImportState = &ServiceDeskBusinessImpact{}

// New creates a new instance of the resource.
func NewServiceDeskBusinessImpact() resource.Resource {
	return &ServiceDeskBusinessImpact{}
}

// Metadata returns the resource type name.
func (r *ServiceDeskBusinessImpact) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_servicedesk_business_impact"
}

// Schema defines the schema for the resource.
func (r *ServiceDeskBusinessImpact) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
				Description: "The ID of the ServiceDeskBusinessImpact. May be retrieved from the backend.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the ServiceDeskBusinessImpact.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The description of the ServiceDeskBusinessImpact.",
			},
			"state": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The state of the ServiceDeskBusinessImpact. Defaults to 'enabled'.",
			},
		},
	}
}

func (r *ServiceDeskBusinessImpact) resolveTenantId(clientAttr types.String) string {
	if !clientAttr.IsNull() && clientAttr.ValueString() != "" {
		return clientAttr.ValueString()
	}

	return r.apiClient.TenantId
}

func (r *ServiceDeskBusinessImpact) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ServiceDeskBusinessImpactModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.resolveTenantId(plan.Client)

	businessImpact := client.ServiceDeskBusinessImpact{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		State:       plan.State.ValueBool(),
	}

	created, err := r.apiClient.CreateServiceDeskBusinessImpact(tenantId, businessImpact)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", err.Error())

		return
	}

	plan.Id = types.StringValue(created.Id)
	plan.State = types.BoolValue(created.State)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ServiceDeskBusinessImpact) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ServiceDeskBusinessImpactModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.resolveTenantId(state.Client)

	businessImpact, err := r.apiClient.GetServiceDeskBusinessImpact(tenantId, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read Error", err.Error())

		return
	}

	// If the businessImpact is not found, remove the resource from state
	// This is a common pattern in Terraform providers to handle resources that may not exist anymore.
	if businessImpact == nil {
		resp.State.RemoveResource(ctx)

		return
	}

	state.Name = types.StringValue(businessImpact.Name)
	state.Description = types.StringValue(businessImpact.Description)
	state.State = types.BoolValue(businessImpact.State)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ServiceDeskBusinessImpact) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ServiceDeskBusinessImpactModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.resolveTenantId(plan.Client)

	businessImpact := client.ServiceDeskBusinessImpact{
		Id:          plan.Id.ValueString(),
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		State:       plan.State.ValueBool(),
	}

	// The API may not support update directly; if not, implement as needed.
	// For now, try create (upsert) pattern:
	updated, err := r.apiClient.UpdateServiceDeskBusinessImpact(tenantId, plan.Id.ValueString(), businessImpact)
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

func (r *ServiceDeskBusinessImpact) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ServiceDeskBusinessImpactModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.resolveTenantId(state.Client)

	err := r.apiClient.DeleteServiceDeskBusinessImpact(tenantId, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", err.Error())

		return
	}

	resp.State.RemoveResource(ctx)
}

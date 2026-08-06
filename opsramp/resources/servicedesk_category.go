// (C) Copyright 2026 Hewlett Packard Enterprise Development LP
package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/opsramp/client"
)

type ServiceDeskCategoryModel struct {
	Client      types.String `tfsdk:"client"`
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	TicketType  types.String `tfsdk:"ticket_type"`
}

// ServiceDeskCategory defines the resource implementation.
type ServiceDeskCategory struct {
	BaseResource
}

// Ensure implementation satisfies the expected interfaces
var (
	_ resource.Resource               = &ServiceDeskCategory{}
	_ resource.ResourceWithModifyPlan = &ServiceDeskCategory{}
)

// var _ resource.ResourceWithImportState = &ServiceDeskCategory{}

// New creates a new instance of the resource.
func NewServiceDeskCategory() resource.Resource {
	return &ServiceDeskCategory{}
}

// Metadata returns the resource type name.
func (r *ServiceDeskCategory) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_servicedesk_category"
}

// Schema defines the schema for the resource.
func (r *ServiceDeskCategory) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
				Description: "The ID of the ServiceDeskCategory. May be retrieved from the backend.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the ServiceDeskCategory.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The description of the ServiceDeskCategory.",
			},
			"ticket_type": schema.StringAttribute{
				Required:    true,
				Description: "The ticket type of the ServiceDeskCategory (e.g. incidents, problems, serviceRequests).",
				Validators: []validator.String{
					stringvalidator.OneOf("incidents", "problems", "serviceRequests"),
				},
			},
		},
	}
}

func (r *ServiceDeskCategory) resolveTenantId(clientAttr types.String) string {
	if !clientAttr.IsNull() && clientAttr.ValueString() != "" {
		return clientAttr.ValueString()
	}

	return r.apiClient.TenantId
}

func (r *ServiceDeskCategory) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ServiceDeskCategoryModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.resolveTenantId(plan.Client)

	category := client.ServiceDeskCategory{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		TicketType:  plan.TicketType.ValueString(),
	}

	created, err := r.apiClient.CreateServiceDeskCategory(tenantId, category)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", err.Error())

		return
	}

	plan.Id = types.StringValue(created.Id)
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ServiceDeskCategory) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ServiceDeskCategoryModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.resolveTenantId(state.Client)

	category, err := r.apiClient.GetServiceDeskCategory(tenantId, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read Error", err.Error())

		return
	}

	// If the category is not found, remove the resource from state
	// This is a common pattern in Terraform providers to handle resources that may not exist anymore.
	if category == nil {
		resp.State.RemoveResource(ctx)

		return
	}

	switch category.TicketType {
	case "Incident":
		state.TicketType = types.StringValue("incidents")
	case "Problem":
		state.TicketType = types.StringValue("problems")
	case "Service Request":
		state.TicketType = types.StringValue("serviceRequests")
	}

	state.Name = types.StringValue(category.Name)
	state.Description = types.StringValue(category.Description)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ServiceDeskCategory) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ServiceDeskCategoryModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	category := client.ServiceDeskCategory{
		Id:          plan.Id.ValueString(),
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		TicketType:  plan.TicketType.ValueString(),
	}

	// The API may not support update directly; if not, implement as needed.
	// For now, try create (upsert) pattern:
	updated, err := r.apiClient.UpdateServiceDeskCategory(r.resolveTenantId(plan.Client), plan.Id.ValueString(), category)
	if err != nil {
		resp.Diagnostics.AddError("Update Error", err.Error())

		return
	}

	plan.Id = types.StringValue(updated.Id)
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ServiceDeskCategory) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ServiceDeskCategoryModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.apiClient.DeleteServiceDeskCategory(r.resolveTenantId(state.Client), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", err.Error())
	}

	resp.State.RemoveResource(ctx)
}

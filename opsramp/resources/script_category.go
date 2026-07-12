// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/HPE/terraform-provider-hpe/opsramp/client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure implementation satisfies the expected interfaces
var _ resource.Resource = &ScriptCategoryResource{}
var _ resource.ResourceWithModifyPlan = &ScriptCategoryResource{}

// ScriptCategoryResource defines the resource implementation.
type ScriptCategoryResource struct {
	BaseResource
}

// ScriptCategoryModel maps Terraform schema attributes to the provider model.
type ScriptCategoryModel struct {
	Client   types.String `tfsdk:"client"`
	Name     types.String `tfsdk:"name"`
	ParentId types.String `tfsdk:"parent_id"`
	Uuid     types.String `tfsdk:"uuid"`
}

// NewScriptCategory creates a new instance of the resource.
func NewScriptCategory() resource.Resource {
	return &ScriptCategoryResource{}
}

// Metadata returns the resource type name.
func (r *ScriptCategoryResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_script_category"
}

// Schema defines the schema for the resource.
func (r *ScriptCategoryResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an OpsRamp RBA Script Category. Categories organize scripts and can be nested using `parent_id`.",
		Attributes: map[string]schema.Attribute{
			"client": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The client (tenant) UUID where this category should be created. If not specified, uses the provider's tenant.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"uuid": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The UUID of the script category, assigned by OpsRamp.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the script category.",
			},
			"parent_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Uuid of the parent category.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *ScriptCategoryResource) resolveTenantId(clientAttr types.String) string {
	if !clientAttr.IsNull() && clientAttr.ValueString() != "" {
		return clientAttr.ValueString()
	}
	return r.apiClient.TenantId
}

// buildCategoryRequest converts the Terraform model into an API request body.
func buildCategoryRequest(plan ScriptCategoryModel) client.ScriptCategory {
	cat := client.ScriptCategory{
		Name: plan.Name.ValueString(),
		Uuid: plan.Uuid.ValueString(),
	}
	if !plan.ParentId.IsNull() && !plan.ParentId.IsUnknown() && plan.ParentId.ValueString() != "" {
		cat.Parent = &client.ScriptCategoryParentRef{Uuid: plan.ParentId.ValueString()}
	}
	return cat
}

// Create handles the creation of the resource.
func (r *ScriptCategoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ScriptCategoryModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.resolveTenantId(plan.Client)
	catReq := buildCategoryRequest(plan)

	created, err := r.apiClient.CreateScriptCategory(tenantId, catReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Creating Script Category", fmt.Sprintf("Could not create script category: %s", err))
		return
	}

	// Update plan with values from the created category
	plan.Uuid = types.StringValue(created.Uuid)
	if created.Parent != nil && created.Parent.Uuid != "" {
		plan.ParentId = types.StringValue(created.Parent.Uuid)
	} else {
		plan.ParentId = types.StringNull()
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Read handles reading the resource state.
func (r *ScriptCategoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ScriptCategoryModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.resolveTenantId(state.Client)

	category, err := r.apiClient.GetScriptCategory(tenantId, state.Uuid.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "No task category found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error Reading Script Category", fmt.Sprintf("Could not read script category %s: %s", state.Uuid.ValueString(), err))
		return
	}

	state.Name = types.StringValue(category.Name)
	if category.Parent != nil && category.Parent.Uuid != "" {
		state.ParentId = types.StringValue(category.Parent.Uuid)
	} else {
		state.ParentId = types.StringNull()
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update handles updating the resource.
func (r *ScriptCategoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ScriptCategoryModel
	var state ScriptCategoryModel
	req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)

	tenantId := r.resolveTenantId(plan.Client)

	catReq := buildCategoryRequest(plan)

	updated, err := r.apiClient.UpdateScriptCategory(tenantId, catReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Updating Script Category", fmt.Sprintf("Could not update script category %s: %s", state.Uuid.ValueString(), err))
		return
	}

	plan.Uuid = types.StringValue(updated.Uuid)
	plan.Name = types.StringValue(updated.Name)
	if updated.Parent != nil && updated.Parent.Uuid != "" {
		plan.ParentId = types.StringValue(updated.Parent.Uuid)
	} else {
		plan.ParentId = types.StringNull()
	}

	diags := resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Delete handles deleting the resource.
func (r *ScriptCategoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ScriptCategoryModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.resolveTenantId(state.Client)

	err := r.apiClient.DeleteScriptCategory(tenantId, state.Uuid.ValueString())
	if err != nil && !strings.Contains(err.Error(), "status: 204") {
		resp.Diagnostics.AddError("Error Deleting Script Category", fmt.Sprintf("Could not delete script category %s: %s", state.Uuid.ValueString(), err))
		return
	}

	resp.State.RemoveResource(ctx)
}

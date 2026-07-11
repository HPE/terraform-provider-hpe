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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure implementation satisfies the expected interfaces
var _ resource.Resource = &KBCategoryResource{}
var _ resource.ResourceWithImportState = &KBCategoryResource{}
var _ resource.ResourceWithModifyPlan = &KBCategoryResource{}

// KBCategoryResource defines the resource implementation.
type KBCategoryResource struct {
	BaseResource
}

// KBCategoryModel maps Terraform schema attributes to the provider model.
type KBCategoryModel struct {
	Client           types.String `tfsdk:"client"`
	Id               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	ParentCategoryId types.String `tfsdk:"parent_category_id"`
	State            types.String `tfsdk:"state"`
}

// NewKBCategory creates a new instance of the resource.
func NewKBCategory() resource.Resource {
	return &KBCategoryResource{}
}

// Metadata returns the resource type name.
func (r *KBCategoryResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kb_category"
}

// Schema defines the schema for the resource.
func (r *KBCategoryResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an OpsRamp Knowledge Base Category. Categories organise KB articles and can optionally reference a parent category.",
		Attributes: map[string]schema.Attribute{
			"client": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The client (tenant) UUID where this category should be created. If not specified, uses the provider's tenant.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The string ID of the KB category, assigned by the API.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the KB category.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The description of the KB category.",
			},
			"parent_category_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The string ID of the parent KB category. Omit for a root-level category.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"state": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The state of the KB category (e.g., `ACTIVE`, `TRASH`).",
			},
		},
	}
}

func (r *KBCategoryResource) resolveTenantId(clientAttr types.String) string {
	if !clientAttr.IsNull() && clientAttr.ValueString() != "" {
		return clientAttr.ValueString()
	}
	return r.apiClient.TenantId
}

func buildKBCategoryRequest(plan KBCategoryModel, apiClient *client.OpsRampClient) client.KBCategory {

	scope := "PARTNER"
	if apiClient.Scope != "MSP" || (!plan.Client.IsNull() && plan.Client.ValueString() != "") {
		scope = "CLIENT"
	}

	cat := client.KBCategory{
		Name:        plan.Name.ValueString(),
		Scope:       scope,
		Description: plan.Description.ValueString(),
	}
	if !plan.ParentCategoryId.IsNull() && !plan.ParentCategoryId.IsUnknown() {
		cat.ParentCategory = &client.KBCategoryRef{Id: plan.ParentCategoryId.ValueString()}
	}
	return cat
}

func populateKBCategoryModel(model *KBCategoryModel, cat *client.KBCategory) {
	model.Id = types.StringValue(cat.Id)
	model.Name = types.StringValue(cat.Name)
	model.Description = types.StringValue(cat.Description)
	model.State = types.StringValue(cat.State)
	if cat.ParentCategory != nil {
		model.ParentCategoryId = types.StringValue(cat.ParentCategory.Id)
	} else {
		model.ParentCategoryId = types.StringNull()
	}
}

// Create handles the creation of the resource.
func (r *KBCategoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan KBCategoryModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.resolveTenantId(plan.Client)
	catReq := buildKBCategoryRequest(plan, r.apiClient)

	created, err := r.apiClient.CreateKBCategory(tenantId, catReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Creating KB Category", fmt.Sprintf("Could not create KB category: %s", err))
		return
	}

	populateKBCategoryModel(&plan, created)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Read handles reading the resource state.
func (r *KBCategoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state KBCategoryModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.resolveTenantId(state.Client)
	categoryId := state.Id.ValueString()

	existing, err := r.apiClient.GetKBCategory(tenantId, categoryId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error Reading KB Category", fmt.Sprintf("Could not read KB category %s: %s", categoryId, err))
		return
	}

	populateKBCategoryModel(&state, existing)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update handles updating the resource.
func (r *KBCategoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state KBCategoryModel
	req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)

	tenantId := r.resolveTenantId(plan.Client)
	categoryId := state.Id.ValueString()
	catReq := buildKBCategoryRequest(plan, r.apiClient)

	updated, err := r.apiClient.UpdateKBCategory(tenantId, categoryId, catReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Updating KB Category", fmt.Sprintf("Could not update KB category %s: %s", categoryId, err))
		return
	}

	plan.Id = state.Id
	populateKBCategoryModel(&plan, updated)

	diags := resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Delete handles deleting the resource.
func (r *KBCategoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state KBCategoryModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.resolveTenantId(state.Client)
	categoryId := state.Id.ValueString()

	err := r.apiClient.DeleteKBCategory(tenantId, categoryId)
	if err != nil {
		resp.Diagnostics.AddError("Error Deleting KB Category", fmt.Sprintf("Could not delete KB category %s: %s", categoryId, err))
		return
	}

	resp.State.RemoveResource(ctx)
}

// ImportState handles importing an existing resource.
// Import ID format: "{categoryId}" or "{tenantId}/{categoryId}"
func (r *KBCategoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")

	var tenantId, categoryId string
	switch len(parts) {
	case 1:
		tenantId = r.apiClient.TenantId
		categoryId = parts[0]
	case 2:
		tenantId = parts[0]
		categoryId = parts[1]
	default:
		resp.Diagnostics.AddError("Invalid Import ID", "Expected format: '{categoryId}' or '{tenantId}/{categoryId}'")
		return
	}

	existing, err := r.apiClient.GetKBCategory(tenantId, categoryId)
	if err != nil {
		resp.Diagnostics.AddError("Error Importing KB Category", fmt.Sprintf("Could not import KB category: %s", err))
		return
	}

	var state KBCategoryModel
	if len(parts) == 2 {
		state.Client = types.StringValue(tenantId)
	} else {
		state.Client = types.StringNull()
	}
	populateKBCategoryModel(&state, existing)

	resp.State.Set(ctx, &state)
}

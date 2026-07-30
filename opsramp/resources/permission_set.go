// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/HPE/terraform-provider-hpe/opsramp/client"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure implementation satisfies the expected interfaces
var (
	_ resource.Resource                = &PermissionSetResource{}
	_ resource.ResourceWithImportState = &PermissionSetResource{}
	_ resource.ResourceWithModifyPlan  = &PermissionSetResource{}
)

// PermissionSetResource defines the resource implementation.
type PermissionSetResource struct {
	BaseResource
}

// PermissionSetModel maps Terraform schema attributes to the provider model.
type PermissionSetModel struct {
	Client      types.String `tfsdk:"client"`
	UniqueId    types.String `tfsdk:"unique_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Permissions types.Set    `tfsdk:"permissions"`
}

// PermissionModel represents a single permission in the set.
type PermissionModel struct {
	Name types.String `tfsdk:"name"`
	Type types.String `tfsdk:"type"`
}

// NewPermissionSet creates a new instance of the resource.
func NewPermissionSet() resource.Resource {
	return &PermissionSetResource{}
}

// Metadata returns the resource type name.
func (r *PermissionSetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_permission_set"
}

// Schema defines the schema for the resource.
func (r *PermissionSetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an OpsRamp Permission Set resource. Permission sets define specific permissions that can be assigned to roles.",
		Attributes: map[string]schema.Attribute{
			"client": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The client (tenant) UUID where this permission set should be created. If not specified, uses the provider's tenant.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"unique_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the permission set.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the permission set.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The description of the permission set.",
			},
			"permissions": schema.SetNestedAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Set of permissions in the set. Order does not matter.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The name of the permission.",
						},
						"type": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The Type of the permission.",
						},
					},
				},
			},
		},
	}
}

// ModifyPlan filters permissions with empty type from the plan before comparison
func (r *PermissionSetResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Call base implementation
	r.BaseResource.ModifyPlan(ctx, req, resp)

	// Don't modify plan during destroy
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan PermissionSetModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If permissions is null or unknown, nothing to filter
	if plan.Permissions.IsNull() || plan.Permissions.IsUnknown() {
		return
	}

	// Extract permissions and filter out empty types
	var permModels []PermissionModel
	diags = plan.Permissions.ElementsAs(ctx, &permModels, false)
	if diags.HasError() {
		return
	}

	// Build filtered permissions set
	var permAttrs []attr.Value
	for _, p := range permModels {
		// Skip permissions with empty type
		if p.Type.IsNull() || p.Type.ValueString() == "" {
			continue
		}
		permObj, _ := types.ObjectValue(
			map[string]attr.Type{
				"name": types.StringType,
				"type": types.StringType,
			},
			map[string]attr.Value{
				"name": p.Name,
				"type": p.Type,
			},
		)
		permAttrs = append(permAttrs, permObj)
	}

	permsSet, _ := types.SetValue(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"name": types.StringType,
				"type": types.StringType,
			},
		},
		permAttrs,
	)
	plan.Permissions = permsSet

	diags = resp.Plan.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *PermissionSetResource) getPermissionsFromPlan(ctx context.Context, plan PermissionSetModel) ([]client.Permission, error) {
	var permissions []client.Permission

	if plan.Permissions.IsNull() || plan.Permissions.IsUnknown() {
		return permissions, nil
	}

	var permModels []PermissionModel
	diags := plan.Permissions.ElementsAs(ctx, &permModels, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to parse permissions")
	}

	for _, p := range permModels {
		// Skip permissions with empty type (they're considered deleted)
		if p.Type.ValueString() == "" {
			continue
		}
		permissions = append(permissions, client.Permission{
			Name: p.Name.ValueString(),
			Type: p.Type.ValueString(),
		})
	}

	return permissions, nil
}

// getPermissionsFromState extracts permissions from current state
func (r *PermissionSetResource) getPermissionsFromState(
	ctx context.Context,
	state PermissionSetModel,
) ([]client.Permission, error) {
	var permissions []client.Permission

	if state.Permissions.IsNull() || state.Permissions.IsUnknown() {
		return permissions, nil
	}

	var permModels []PermissionModel
	diags := state.Permissions.ElementsAs(ctx, &permModels, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to parse permissions")
	}

	for _, p := range permModels {
		if p.Type.ValueString() == "" {
			continue
		}
		permissions = append(permissions, client.Permission{
			Name: p.Name.ValueString(),
			Type: p.Type.ValueString(),
		})
	}

	return permissions, nil
}

// Create handles the creation of the resource.
func (r *PermissionSetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan PermissionSetModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build permissions set (filters out empty types)
	permissions, err := r.getPermissionsFromPlan(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError("Error parsing permissions", err.Error())

		return
	}

	// Use client parameter if set, otherwise provider's tenant ID + scope
	tenantId := r.apiClient.TenantId
	scope := r.apiClient.Scope

	if !plan.Client.IsNull() && plan.Client.ValueString() != "" {
		scope = "CLIENT"
		tenantId = plan.Client.ValueString()
	}

	createPermSet := client.CreatePermissionSet{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Permissions: permissions,
		Scope:       scope,
	}

	created, err := r.apiClient.CreatePermissionSet(tenantId, createPermSet)
	if err != nil {
		resp.Diagnostics.AddError("Creation Error", err.Error())

		return
	}

	plan.UniqueId = types.StringValue(created.UniqueId)

	// Build permissions set for state - use returned permissions which have IDs
	var permAttrs []attr.Value
	for _, p := range created.Permissions {
		if p.Type == "" {
			continue
		}
		permObj, _ := types.ObjectValue(
			map[string]attr.Type{
				"name": types.StringType,
				"type": types.StringType,
			},
			map[string]attr.Value{
				"name": types.StringValue(p.Name),
				"type": types.StringValue(p.Type),
			},
		)
		permAttrs = append(permAttrs, permObj)
	}

	permsSet, _ := types.SetValue(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"name": types.StringType,
				"type": types.StringType,
			},
		},
		permAttrs,
	)
	plan.Permissions = permsSet

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Read handles reading the resource state.
func (r *PermissionSetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state PermissionSetModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use client parameter if set, otherwise provider's tenant ID
	tenantId := r.apiClient.TenantId
	if !state.Client.IsNull() && state.Client.ValueString() != "" {
		tenantId = state.Client.ValueString()
	}

	existing, err := r.apiClient.GetPermissionSet(tenantId, state.UniqueId.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)

			return
		}
		resp.Diagnostics.AddError("Read Error", err.Error())

		return
	}

	// Update state from API response
	state.UniqueId = types.StringValue(existing.UniqueId)
	state.Name = types.StringValue(existing.Name)
	state.Description = types.StringValue(existing.Description)

	// Update permissions set - only include permissions with non-empty type
	var permAttrs []attr.Value
	for _, p := range existing.Permissions {
		// Skip permissions with empty type (they're deleted/inactive)
		if p.Type == "" {
			continue
		}
		permObj, _ := types.ObjectValue(
			map[string]attr.Type{
				"name": types.StringType,
				"type": types.StringType,
			},
			map[string]attr.Value{
				"name": types.StringValue(p.Name),
				"type": types.StringValue(p.Type),
			},
		)
		permAttrs = append(permAttrs, permObj)
	}

	permsSet, _ := types.SetValue(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"name": types.StringType,
				"type": types.StringType,
			},
		},
		permAttrs,
	)
	state.Permissions = permsSet

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update handles updating the resource.
func (r *PermissionSetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state PermissionSetModel
	req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)

	// Use client parameter if set, otherwise provider's tenant ID
	tenantId := r.apiClient.TenantId
	if !plan.Client.IsNull() && plan.Client.ValueString() != "" {
		tenantId = plan.Client.ValueString()
	}

	// Build permissions from plan
	planPermissions, err := r.getPermissionsFromPlan(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError("Error parsing plan permissions", err.Error())

		return
	}

	// Build permissions from current state
	statePermissions, err := r.getPermissionsFromState(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError("Error parsing state permissions", err.Error())

		return
	}

	// Create a map of plan permission names for quick lookup
	planPermMap := make(map[string]bool)
	for _, p := range planPermissions {
		planPermMap[p.Name] = true
	}

	// Find permissions in state but not in plan - these need to be deleted (type="")
	var permissions []client.Permission
	permissions = append(permissions, planPermissions...)
	for _, sp := range statePermissions {
		if !planPermMap[sp.Name] {
			// Permission was removed, send with empty type to delete
			permissions = append(permissions, client.Permission{
				Name: sp.Name,
				Type: "",
			})
		}
	}

	updatePermSet := client.UpdatePermissionSet{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Permissions: permissions,
	}

	updated, err := r.apiClient.UpdatePermissionSet(tenantId, state.UniqueId.ValueString(), updatePermSet)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Permission Set",
			fmt.Sprintf("Could not update permission set: %s", err),
		)

		return
	}

	plan.UniqueId = state.UniqueId

	// Rebuild permissions from response
	var updatedPermAttrs []attr.Value
	for _, p := range updated.Permissions {
		if p.Type == "" {
			continue
		}
		permObj, _ := types.ObjectValue(
			map[string]attr.Type{
				"name": types.StringType,
				"type": types.StringType,
			},
			map[string]attr.Value{
				"name": types.StringValue(p.Name),
				"type": types.StringValue(p.Type),
			},
		)
		updatedPermAttrs = append(updatedPermAttrs, permObj)
	}
	plan.Permissions, _ = types.SetValue(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"name": types.StringType,
				"type": types.StringType,
			},
		},
		updatedPermAttrs,
	)

	diags := resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Delete handles deleting the resource.
func (r *PermissionSetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state PermissionSetModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use client parameter if set, otherwise provider's tenant ID
	tenantId := r.apiClient.TenantId
	if !state.Client.IsNull() && state.Client.ValueString() != "" {
		tenantId = state.Client.ValueString()
	}

	err := r.apiClient.DeletePermissionSet(tenantId, state.UniqueId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", err.Error())

		return
	}

	resp.State.RemoveResource(ctx)
}

// ImportState handles importing an existing resource.
func (r *PermissionSetResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// Import format: permSetId (uses provider's tenant)
	permSetId := req.ID

	existing, err := r.apiClient.GetPermissionSet(r.apiClient.TenantId, permSetId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Importing Permission Set",
			fmt.Sprintf("Could not import permission set: %s", err),
		)

		return
	}

	// Build permissions list - only include permissions with non-empty type
	var permAttrs []attr.Value
	for _, p := range existing.Permissions {
		// Skip permissions with empty type (they're deleted/inactive)
		if p.Type == "" {
			continue
		}
		permObj, _ := types.ObjectValue(
			map[string]attr.Type{
				"name": types.StringType,
				"type": types.StringType,
			},
			map[string]attr.Value{
				"name": types.StringValue(p.Name),
				"type": types.StringValue(p.Type),
			},
		)
		permAttrs = append(permAttrs, permObj)
	}

	permsSet, _ := types.SetValue(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"name": types.StringType,
				"type": types.StringType,
			},
		},
		permAttrs,
	)

	importState := PermissionSetModel{
		UniqueId:    types.StringValue(existing.UniqueId),
		Name:        types.StringValue(existing.Name),
		Description: types.StringValue(existing.Description),
		Permissions: permsSet,
	}

	resp.State.Set(ctx, &importState)
}

// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/HPE/terraform-provider-opsramp/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure implementation satisfies the expected interfaces
var _ resource.Resource = &UserGroupResource{}
var _ resource.ResourceWithImportState = &UserGroupResource{}
var _ resource.ResourceWithModifyPlan = &UserGroupResource{}

// UserGroupResource defines the resource implementation.
type UserGroupResource struct {
	BaseResource
}

// UserGroupModel maps Terraform schema attributes to the provider model.
type UserGroupModel struct {
	Client      types.String   `tfsdk:"client"`
	UniqueId    types.String   `tfsdk:"unique_id"`
	Name        types.String   `tfsdk:"name"`
	Description types.String   `tfsdk:"description"`
	Users       []types.String `tfsdk:"users"`
	Roles       []types.String `tfsdk:"roles"`
}

// NewUserGroup creates a new instance of the resource.
func NewUserGroup() resource.Resource {
	return &UserGroupResource{}
}

// Metadata returns the resource type name.
func (r *UserGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_group"
}

// Schema defines the schema for the resource.
func (r *UserGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an OpsRamp User Group resource. User groups organize users and can be assigned roles.",
		Attributes: map[string]schema.Attribute{
			"client": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The client (tenant) UUID where this user should be created. If not specified, uses the provider's tenant.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"unique_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "A unique identifier for the user group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the user group.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The description of the user group.",
			},
			"users": schema.SetAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "List of user IDs in the group.",
			},
			"roles": schema.SetAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "List of role IDs assigned to the group.",
			},
		},
	}
}

// Create handles the creation of the resource.
func (r *UserGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan UserGroupModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use client parameter if set, otherwise provider's tenant ID
	tenantId := r.apiClient.TenantId
	if !plan.Client.IsNull() && plan.Client.ValueString() != "" {
		tenantId = plan.Client.ValueString()
	}

	// Build roles list
	var roles []client.UserRoleRef
	for _, roleId := range plan.Roles {
		roles = append(roles, client.UserRoleRef{UniqueId: roleId.ValueString()})
	}

	createGroup := client.CreateUserGroup{
		Name:        plan.Name.ValueString(),
		UniqueId:    plan.UniqueId.ValueString(),
		Description: plan.Description.ValueString(),
		Roles:       roles,
	}

	created, err := r.apiClient.CreateUserGroup(tenantId, createGroup)
	if err != nil {
		resp.Diagnostics.AddError("Creation Error", err.Error())
		return
	}

	addedUserIds := make([]types.String, 0)

	if len(plan.Users) > 0 {

		// Build users list
		var users []client.UserRef
		for _, userId := range plan.Users {
			users = append(users, client.UserRef{Id: userId.ValueString()})
		}

		err = r.apiClient.AddUserToUserGroup(tenantId, created.UniqueId, users)
		if err != nil {
			resp.Diagnostics.AddError("Error Adding Users to Group", err.Error())
			return
		}

		// Store the successfully added users in state
		for _, u := range users {
			addedUserIds = append(addedUserIds, types.StringValue(u.Id))
		}
	}

	plan.UniqueId = types.StringValue(created.UniqueId)
	if plan.Users != nil {
		plan.Users = addedUserIds
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Read handles reading the resource state.
func (r *UserGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state UserGroupModel
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

	existing, err := r.apiClient.GetUserGroup(tenantId, state.UniqueId.ValueString())
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

	// Fetch users from separate endpoint
	existingUsers, err := r.apiClient.GetUserGroupUsers(tenantId, state.UniqueId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error Reading User Group Users", err.Error())
		return
	}

	// Update users list while preserving null for optional unset attribute.
	if state.Users != nil {
		users := make([]types.String, 0, len(existingUsers))
		for _, user := range existingUsers {
			users = append(users, types.StringValue(user.Id))
		}
		state.Users = users
	}

	// Update roles list while preserving null for optional unset attribute.
	if state.Roles != nil {
		roles := make([]types.String, 0, len(existing.Roles))
		for _, role := range existing.Roles {
			roles = append(roles, types.StringValue(role.UniqueId))
		}
		state.Roles = roles
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update handles updating the resource.
func (r *UserGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state UserGroupModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use client parameter if set, otherwise provider's tenant ID
	tenantId := r.apiClient.TenantId
	if !plan.Client.IsNull() && plan.Client.ValueString() != "" {
		tenantId = plan.Client.ValueString()
	}

	// Build roles list. If roles are omitted in config (nil), preserve existing roles.
	desiredRoles := plan.Roles
	if desiredRoles == nil {
		desiredRoles = state.Roles
	}

	var roles []client.UserRoleRef
	for _, roleId := range desiredRoles {
		roles = append(roles, client.UserRoleRef{UniqueId: roleId.ValueString()})
	}

	updateGroup := client.UpdateUserGroup{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Roles:       roles,
	}

	_, err := r.apiClient.UpdateUserGroup(tenantId, state.UniqueId.ValueString(), updateGroup)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating User Group",
			fmt.Sprintf("Could not update user group: %s", err),
		)
		return
	}

	// Manage users: compute diff between current state and plan
	currentUsers := make(map[string]bool)
	for _, u := range state.Users {
		currentUsers[u.ValueString()] = true
	}
	plannedUsers := make(map[string]bool)
	for _, u := range plan.Users {
		plannedUsers[u.ValueString()] = true
	}

	// Users to add (in plan but not in state)
	var usersToAdd []client.UserRef
	for _, u := range plan.Users {
		if !currentUsers[u.ValueString()] {
			usersToAdd = append(usersToAdd, client.UserRef{Id: u.ValueString()})
		}
	}

	// Users to remove (in state but not in plan)
	var usersToRemove []client.UserRef
	for _, u := range state.Users {
		if !plannedUsers[u.ValueString()] {
			usersToRemove = append(usersToRemove, client.UserRef{Id: u.ValueString()})
		}
	}

	if len(usersToAdd) > 0 {
		err := r.apiClient.AddUserToUserGroup(tenantId, state.UniqueId.ValueString(), usersToAdd)
		if err != nil {
			resp.Diagnostics.AddError("Error Adding Users to Group", err.Error())
			return
		}
	}

	if len(usersToRemove) > 0 {
		err := r.apiClient.RemoveUsersFromUserGroup(tenantId, state.UniqueId.ValueString(), usersToRemove)
		if err != nil {
			resp.Diagnostics.AddError("Error Removing Users from Group", err.Error())
			return
		}
	}

	plan.UniqueId = state.UniqueId
	plan.Roles = desiredRoles

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Delete handles deleting the resource.
func (r *UserGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state UserGroupModel
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

	err := r.apiClient.DeleteUserGroup(tenantId, state.UniqueId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

// ImportState handles importing an existing resource.
func (r *UserGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import format: groupId (uses provider's tenant)
	groupId := req.ID

	existingGroup, err := r.apiClient.GetUserGroup(r.apiClient.TenantId, groupId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Importing User Group",
			fmt.Sprintf("Could not import user group: %s", err),
		)
		return
	}

	// Build users list from separate endpoint
	existingUsers, err2 := r.apiClient.GetUserGroupUsers(r.apiClient.TenantId, groupId)
	if err2 != nil {
		resp.Diagnostics.AddError("Error Fetching User Group Users", err2.Error())
		return
	}

	users := make([]types.String, 0, len(existingUsers))
	for _, user := range existingUsers {
		users = append(users, types.StringValue(user.Id))
	}

	// Build roles list
	roles := make([]types.String, 0, len(existingGroup.Roles))
	for _, role := range existingGroup.Roles {
		roles = append(roles, types.StringValue(role.UniqueId))
	}

	state := UserGroupModel{
		UniqueId:    types.StringValue(existingGroup.UniqueId),
		Name:        types.StringValue(existingGroup.Name),
		Description: types.StringValue(existingGroup.Description),
		Users:       users,
		Roles:       roles,
	}

	resp.State.Set(ctx, &state)
}

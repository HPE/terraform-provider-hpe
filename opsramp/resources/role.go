// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/HPE/terraform-provider-hpe/opsramp/client"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure implementation satisfies the expected interfaces
var _ resource.Resource = &RoleResource{}
var _ resource.ResourceWithImportState = &RoleResource{}
var _ resource.ResourceWithModifyPlan = &RoleResource{}

// RoleResource defines the resource implementation.
type RoleResource struct {
	BaseResource
}

// RoleModel maps Terraform schema attributes to the provider model.
type RoleModel struct {
	Client         types.String   `tfsdk:"client"`
	Id             types.String   `tfsdk:"id"`
	Name           types.String   `tfsdk:"name"`
	Description    types.String   `tfsdk:"description"`
	DefaultRole    types.Bool     `tfsdk:"default_role"`
	AllClients     types.Bool     `tfsdk:"all_clients"`
	AllDevices     types.Bool     `tfsdk:"all_devices"`
	AllCredentials types.Bool     `tfsdk:"all_credentials"`
	AllAuthzTags   types.Bool     `tfsdk:"all_authz_tags"`
	Clients        []types.String `tfsdk:"clients"`
	Users          types.Set      `tfsdk:"users"`
	UserGroups     types.Set      `tfsdk:"user_groups"`
	Devices        []types.String `tfsdk:"devices"`
	DeviceGroups   []types.String `tfsdk:"device_groups"`
	CredentialSets []types.String `tfsdk:"credential_sets"`
	Permissions    []types.String `tfsdk:"permissions"`
}

// NewRole creates a new instance of the resource.
func NewRole() resource.Resource {
	return &RoleResource{}
}

// Metadata returns the resource type name.
func (r *RoleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role"
}

// Schema defines the schema for the resource.
func (r *RoleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an OpsRamp Role resource. Roles define access control with permission sets, users, devices, and more.",
		Attributes: map[string]schema.Attribute{
			"client": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The client (tenant) UUID where this role should be created. If not specified, uses the provider's tenant.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the role.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the role.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The description of the role.",
			},
			"default_role": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether this is a default role.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"all_clients": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether the role applies to all clients. Computed based on whether clients list is empty.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"all_devices": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether the role applies to all devices. Computed based on whether devices/device_groups lists are empty.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"all_credentials": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether the role applies to all credentials. Computed based on whether credential_sets list is empty.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"all_authz_tags": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether the role applies to all authorization tags. Always true.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"clients": schema.SetAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "List of client unique IDs associated with the role.",
			},
			"users": schema.SetAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "List of user IDs associated with the role. Read-only - manage role assignments via the user resource.",
			},
			"user_groups": schema.SetAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "List of user group unique IDs associated with the role. Read-only - manage role assignments via the user_group resource.",
			},
			"devices": schema.SetAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "List of device IDs associated with the role.",
			},
			"device_groups": schema.SetAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "List of device group IDs associated with the role.",
			},
			"credential_sets": schema.SetAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "List of credential set unique IDs associated with the role.",
			},
			"permissions": schema.SetAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "List of permission set unique IDs assigned to the role.",
			},
		},
	}
}

// Create handles the creation of the resource.
func (r *RoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan RoleModel
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

	// Build clients list
	var clients []client.RoleClientRef
	for _, clientId := range plan.Clients {
		clients = append(clients, client.RoleClientRef{UniqueId: clientId.ValueString()})
	}

	// Build devices list
	var devices []client.RoleDeviceRef
	for _, deviceId := range plan.Devices {
		devices = append(devices, client.RoleDeviceRef{Id: deviceId.ValueString()})
	}

	// Build device groups list
	var deviceGroups []client.RoleDeviceGroupRef
	for _, groupId := range plan.DeviceGroups {
		deviceGroups = append(deviceGroups, client.RoleDeviceGroupRef{Id: groupId.ValueString()})
	}

	// Build credential sets list
	var credentialSets []client.RoleCredentialRef
	for _, credId := range plan.CredentialSets {
		credentialSets = append(credentialSets, client.RoleCredentialRef{UniqueId: credId.ValueString()})
	}

	// Build permissions list - look up numeric IDs from unique IDs
	var permissions []client.RolePermissionRef
	for _, permUniqueId := range plan.Permissions {
		numericId, err := r.apiClient.GetPermissionSetIdByUniqueId(tenantId, permUniqueId.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error looking up permission set",
				fmt.Sprintf("Could not find permission set with uniqueId '%s': %s", permUniqueId.ValueString(), err),
			)
			return
		}
		permissions = append(permissions, client.RolePermissionRef{Id: numericId})
	}

	scope := "MSP"
	if r.apiClient.Scope != "MSP" || (!plan.Client.IsNull() && plan.Client.ValueString() != "") {
		scope = "CLIENT"
	}

	// Compute boolean flags based on list contents
	allClients := len(clients) == 0
	allDevices := len(devices) == 0 && len(deviceGroups) == 0
	allCredentials := len(credentialSets) == 0
	allAuthzTags := true // Always true

	createRole := client.Role{
		Name:           plan.Name.ValueString(),
		Description:    plan.Description.ValueString(),
		Scope:          scope,
		Devices:        devices,
		DeviceGroups:   deviceGroups,
		CredentialSets: credentialSets,
		Permissions:    permissions,
		AllDevices:     allDevices,
		AllCredentials: allCredentials,
		AllAuthzTags:   allAuthzTags,
	}

	// For CLIENT scope, do not send allClients or clients
	if scope != "CLIENT" {
		createRole.Clients = clients
		createRole.AllClients = allClients
	}

	created, err := r.apiClient.CreateRole(tenantId, createRole)
	if err != nil {
		resp.Diagnostics.AddError("Creation Error", err.Error())
		return
	}

	plan.Id = types.StringValue(created.UniqueId)
	plan.DefaultRole = types.BoolValue(created.DefaultRole)
	plan.AllClients = types.BoolValue(created.AllClients)
	plan.AllDevices = types.BoolValue(created.AllDevices)
	plan.AllCredentials = types.BoolValue(created.AllCredentials)
	plan.AllAuthzTags = types.BoolValue(created.AllAuthzTags)

	// Populate computed users from API response
	userElems := make([]attr.Value, len(created.Users))
	for i, u := range created.Users {
		userElems[i] = types.StringValue(u.Id)
	}
	plan.Users, diags = types.SetValue(types.StringType, userElems)
	resp.Diagnostics.Append(diags...)

	// Populate computed user groups from API response
	ugElems := make([]attr.Value, len(created.UserGroups))
	for i, ug := range created.UserGroups {
		ugElems[i] = types.StringValue(ug.UniqueId)
	}
	plan.UserGroups, diags = types.SetValue(types.StringType, ugElems)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Read handles reading the resource state.
func (r *RoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state RoleModel
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

	existing, err := r.apiClient.GetRole(tenantId, state.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read Error", err.Error())
		return
	}

	// Update state from API response
	state.Id = types.StringValue(existing.UniqueId)
	state.Name = types.StringValue(existing.Name)
	state.Description = types.StringValue(existing.Description)
	// Update boolean flags
	state.DefaultRole = types.BoolValue(existing.DefaultRole)
	state.AllClients = types.BoolValue(existing.AllClients)
	state.AllDevices = types.BoolValue(existing.AllDevices)
	state.AllCredentials = types.BoolValue(existing.AllCredentials)
	state.AllAuthzTags = types.BoolValue(existing.AllAuthzTags)

	// Update clients - only update if API returned clients, preserve state otherwise
	if len(existing.Clients) > 0 {
		var clients []types.String
		for _, c := range existing.Clients {
			clients = append(clients, types.StringValue(c.UniqueId))
		}
		state.Clients = clients
	}

	// Update users - always update from API since this is computed
	userElems := make([]attr.Value, len(existing.Users))
	for i, u := range existing.Users {
		userElems[i] = types.StringValue(u.Id)
	}
	state.Users, diags = types.SetValue(types.StringType, userElems)
	resp.Diagnostics.Append(diags...)

	// Update user groups - always update from API since this is computed
	ugElems := make([]attr.Value, len(existing.UserGroups))
	for i, ug := range existing.UserGroups {
		ugElems[i] = types.StringValue(ug.UniqueId)
	}
	state.UserGroups, diags = types.SetValue(types.StringType, ugElems)
	resp.Diagnostics.Append(diags...)

	// Update devices - only update if API returned devices
	// Update devices - only update if API returned devices
	if len(existing.Devices) > 0 {
		var devices []types.String
		for _, d := range existing.Devices {
			devices = append(devices, types.StringValue(d.Id))
		}
		state.Devices = devices
	}

	// Update device groups - only update if API returned device groups
	if len(existing.DeviceGroups) > 0 {
		var deviceGroups []types.String
		for _, dg := range existing.DeviceGroups {
			deviceGroups = append(deviceGroups, types.StringValue(dg.Id))
		}
		state.DeviceGroups = deviceGroups
	}

	// Update credential sets - only update if API returned credential sets
	if len(existing.CredentialSets) > 0 {
		var credentialSets []types.String
		for _, cs := range existing.CredentialSets {
			credentialSets = append(credentialSets, types.StringValue(cs.UniqueId))
		}
		state.CredentialSets = credentialSets
	}

	// Permissions: preserve the prior state since API returns numeric IDs
	// and we store unique IDs - no conversion needed

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update handles updating the resource.
func (r *RoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state RoleModel
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

	// Build clients list
	var clients []client.RoleClientRef
	for _, clientId := range plan.Clients {
		clients = append(clients, client.RoleClientRef{UniqueId: clientId.ValueString()})
	}

	// Build devices list
	var devices []client.RoleDeviceRef
	for _, deviceId := range plan.Devices {
		devices = append(devices, client.RoleDeviceRef{Id: deviceId.ValueString()})
	}

	// Build device groups list
	var deviceGroups []client.RoleDeviceGroupRef
	for _, groupId := range plan.DeviceGroups {
		deviceGroups = append(deviceGroups, client.RoleDeviceGroupRef{Id: groupId.ValueString()})
	}

	// Build credential sets list
	var credentialSets []client.RoleCredentialRef
	for _, credId := range plan.CredentialSets {
		credentialSets = append(credentialSets, client.RoleCredentialRef{UniqueId: credId.ValueString()})
	}

	// Build permissions list - look up numeric IDs from unique IDs
	var permissions []client.RolePermissionRef
	for _, permUniqueId := range plan.Permissions {
		numericId, err := r.apiClient.GetPermissionSetIdByUniqueId(tenantId, permUniqueId.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error looking up permission set",
				fmt.Sprintf("Could not find permission set with uniqueId '%s': %s", permUniqueId.ValueString(), err),
			)
			return
		}
		permissions = append(permissions, client.RolePermissionRef{Id: numericId})
	}

	scope := "MSP"
	if r.apiClient.Scope != "MSP" || (!plan.Client.IsNull() && plan.Client.ValueString() != "") {
		scope = "CLIENT"
	}

	// Compute boolean flags based on list contents
	allClients := len(clients) == 0
	allDevices := len(devices) == 0 && len(deviceGroups) == 0
	allCredentials := len(credentialSets) == 0
	allAuthzTags := true // Always true

	updateRole := client.Role{
		Name:           plan.Name.ValueString(),
		Description:    plan.Description.ValueString(),
		Scope:          scope,
		Devices:        devices,
		DeviceGroups:   deviceGroups,
		CredentialSets: credentialSets,
		Permissions:    permissions,
		AllDevices:     allDevices,
		AllCredentials: allCredentials,
		AllAuthzTags:   allAuthzTags,
	}

	// For CLIENT scope, do not send allClients or clients
	if scope != "CLIENT" {
		updateRole.Clients = clients
		updateRole.AllClients = allClients
	}

	updated, err := r.apiClient.UpdateRole(tenantId, state.Id.ValueString(), updateRole)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Role",
			fmt.Sprintf("Could not update role: %s", err),
		)
		return
	}

	plan.Id = state.Id
	plan.DefaultRole = types.BoolValue(updated.DefaultRole)
	plan.AllClients = types.BoolValue(updated.AllClients)
	plan.AllDevices = types.BoolValue(updated.AllDevices)
	plan.AllCredentials = types.BoolValue(updated.AllCredentials)
	plan.AllAuthzTags = types.BoolValue(updated.AllAuthzTags)

	// Populate computed users from API response
	userElems := make([]attr.Value, len(updated.Users))
	for i, u := range updated.Users {
		userElems[i] = types.StringValue(u.Id)
	}
	plan.Users, diags = types.SetValue(types.StringType, userElems)
	resp.Diagnostics.Append(diags...)

	// Populate computed user groups from API response
	ugElems := make([]attr.Value, len(updated.UserGroups))
	for i, ug := range updated.UserGroups {
		ugElems[i] = types.StringValue(ug.UniqueId)
	}
	plan.UserGroups, diags = types.SetValue(types.StringType, ugElems)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Delete handles deleting the resource.
func (r *RoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state RoleModel
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

	err := r.apiClient.DeleteRole(tenantId, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

// ImportState handles importing an existing resource.
func (r *RoleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import format: roleId (uses provider's tenant)
	roleId := req.ID

	existing, err := r.apiClient.GetRole(r.apiClient.TenantId, roleId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Importing Role",
			fmt.Sprintf("Could not import role: %s", err),
		)
		return
	}

	// Build clients list
	var clients []types.String
	for _, c := range existing.Clients {
		clients = append(clients, types.StringValue(c.UniqueId))
	}

	// Build users list
	userElems := make([]attr.Value, len(existing.Users))
	for i, u := range existing.Users {
		userElems[i] = types.StringValue(u.Id)
	}
	users, userDiags := types.SetValue(types.StringType, userElems)
	resp.Diagnostics.Append(userDiags...)

	// Build user groups list
	ugElems := make([]attr.Value, len(existing.UserGroups))
	for i, ug := range existing.UserGroups {
		ugElems[i] = types.StringValue(ug.UniqueId)
	}
	userGroups, ugDiags := types.SetValue(types.StringType, ugElems)
	resp.Diagnostics.Append(ugDiags...)

	// Build devices list
	var devices []types.String
	for _, d := range existing.Devices {
		devices = append(devices, types.StringValue(d.Id))
	}

	// Build device groups list
	var deviceGroups []types.String
	for _, dg := range existing.DeviceGroups {
		deviceGroups = append(deviceGroups, types.StringValue(dg.Id))
	}

	// Build credential sets list
	var credentialSets []types.String
	for _, cs := range existing.CredentialSets {
		credentialSets = append(credentialSets, types.StringValue(cs.UniqueId))
	}

	// Build permissions list - look up unique IDs from numeric IDs
	var permissions []types.String
	permSetList, err := r.apiClient.SearchPermissionSets(r.apiClient.TenantId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error fetching permission sets",
			fmt.Sprintf("Could not fetch permission sets for import: %s", err),
		)
		return
	}
	permIdToUniqueId := make(map[int]string)
	for _, ps := range permSetList.Results {
		permIdToUniqueId[ps.Id] = ps.UniqueId
	}
	for _, p := range existing.Permissions {
		if uniqueId, ok := permIdToUniqueId[p.Id]; ok {
			permissions = append(permissions, types.StringValue(uniqueId))
		}
	}

	state := RoleModel{
		Id:             types.StringValue(existing.UniqueId),
		Name:           types.StringValue(existing.Name),
		Description:    types.StringValue(existing.Description),
		DefaultRole:    types.BoolValue(existing.DefaultRole),
		AllClients:     types.BoolValue(existing.AllClients),
		AllDevices:     types.BoolValue(existing.AllDevices),
		AllCredentials: types.BoolValue(existing.AllCredentials),
		AllAuthzTags:   types.BoolValue(existing.AllAuthzTags),
		Clients:        clients,
		Users:          users,
		UserGroups:     userGroups,
		Devices:        devices,
		DeviceGroups:   deviceGroups,
		CredentialSets: credentialSets,
		Permissions:    permissions,
	}

	resp.State.Set(ctx, &state)
}

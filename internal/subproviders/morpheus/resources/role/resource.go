// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package role

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/compare"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/errors"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource = &Resource{}
)

func NewResource() resource.Resource {
	return &Resource{}
}

// Resource defines the resource implementation.
type Resource struct {
	configure.ResourceWithMorpheusConfigure
	resource.Resource
}

func (r *Resource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_role"
}

func (r *Resource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = RoleResourceSchema(ctx)
}

// populate role resource model with current API values
func getRoleAsState(
	ctx context.Context,
	id int64,
	client *sdk.APIClient,
) (RoleModel, diag.Diagnostics) {
	var state RoleModel
	var diags diag.Diagnostics

	r, hresp, err := client.RolesAPI.GetRole(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		diags.AddError(
			"populate role resource",
			fmt.Sprintf("role %d GET failed: ", id)+errors.ErrMsg(err, hresp),
		)

		return state, diags
	}

	state.Id = convert.Int64ToType(r.Role.Id)
	state.Name = convert.StrToType(r.Role.Name)
	state.Description = convert.StrToType(r.Role.Description)
	state.LandingUrl = convert.StrToType(r.Role.LandingUrl)
	state.Multitenant = convert.BoolToType(r.Role.Multitenant)
	state.MultitenantLocked = convert.BoolToType(r.Role.MultitenantLocked)
	state.RoleType = convert.StrToType(r.Role.RoleType)

	// for sorting the permission keys and storing to state, we don't want the Role properties,
	// only the permission related ones
	r.Role = nil

	sortedPermissions, err := json.Marshal(r)
	if err != nil {
		diags.AddError(
			"get role (read permissions)",
			fmt.Sprintf("role %d: failed to marshal permissions: "+err.Error(), id),
		)

		return state, diags
	}

	sortedPermissionsStr := string(sortedPermissions)

	state.Permissions = convert.StrToType(&sortedPermissionsStr)

	return state, diags
}

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan RoleModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRole := sdk.NewAddRolesRequestRoleWithDefaults()
	name := plan.Name.ValueString()

	// Only add to create request if user has set permissions explicitly.
	// Also, set permissions first so that it doesn't override the other
	// addRole fields when we unmarshal.
	if !plan.Permissions.IsNull() && !plan.Permissions.IsUnknown() {

		data := []byte(plan.Permissions.ValueString())

		// populate the addRole request with user-provided config data
		err := json.Unmarshal(data, &addRole)
		if err != nil {
			resp.Diagnostics.AddError(
				"create role resource",
				"role "+name+": failed to unmarshal permissions to request: "+err.Error(),
			)

			return

		}
	}

	// required
	addRole.SetAuthority(name)

	// optional
	if !plan.Description.IsUnknown() {
		addRole.SetDescription(plan.Description.ValueString())
	}

	if !plan.LandingUrl.IsUnknown() {
		addRole.SetLandingUrl(plan.LandingUrl.ValueString())
	}

	// optional_computed
	if !plan.Multitenant.IsUnknown() {
		// default: false
		addRole.SetMultitenant(plan.Multitenant.ValueBool())
	}
	if !plan.MultitenantLocked.IsUnknown() {
		// default: false
		addRole.SetMultitenantLocked(plan.MultitenantLocked.ValueBool())
	}

	if !plan.RoleType.IsUnknown() {
		// default: user
		addRole.SetRoleType(plan.RoleType.ValueString())
	}

	addRoleReq := sdk.NewAddRolesRequest(*addRole)

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"create role resource",
			"role "+name+": failed to create client: "+err.Error(),
		)

		return
	}

	role, hresp, err := client.RolesAPI.AddRoles(ctx).
		AddRolesRequest(*addRoleReq).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"create role resource",
			"role "+name+" POST failed: "+errors.ErrMsg(err, hresp),
		)

		return
	}

	if role.GetRole().Id == nil {
		resp.Diagnostics.AddError(
			"create role resource",
			"role "+name+": id is nil",
		)

		return
	}

	id := *role.GetRole().Id
	plan.Id = types.Int64Value(id)

	// write id as soon as possible
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, pdiags := getRoleAsState(ctx, id, client)
	if pdiags.HasError() {
		resp.Diagnostics.Append(pdiags...)
		resp.Diagnostics.AddError(
			"create role resource",
			fmt.Sprintf("role %d: failed to read from api", id),
		)

		return
	}

	// If the user provided a config as part of the create,
	// then set the state to what was in the plan (optional).
	// Otherwise, in the case of the user providing NO config,
	// set it to what was read from the API (computed, set in getRoleAsState).
	if !plan.Permissions.IsNull() && !plan.Permissions.IsUnknown() {
		state.Permissions = plan.Permissions
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state RoleModel

	diags := req.State.Get(ctx, &state)
	if diags.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"read role resource",
			"new client call failed with "+err.Error(),
		)

		return
	}

	id := state.Id.ValueInt64()
	apiState, pdiags := getRoleAsState(ctx, id, client)
	if pdiags.HasError() {
		resp.Diagnostics.Append(pdiags...)
		resp.Diagnostics.AddError(
			"read role resource",
			fmt.Sprintf("role %d: failed to read from api", id),
		)

		return
	}

	var statePermissionData, apiPermissionData sdk.GetRole200Response

	// On import, or when the user does not set the permissions attribute,
	// the permissions attribute will be null or unknown, so we need to ignore the subset check
	// and just set it to the API Permissions - i.e. fully computed
	if !state.Permissions.IsNull() && !state.Permissions.IsUnknown() {

		statePermissionStr := state.Permissions.ValueString()
		err = json.Unmarshal([]byte(statePermissionStr), &statePermissionData)
		if err != nil {
			resp.Diagnostics.Append(pdiags...)
			resp.Diagnostics.AddError(
				"read role resource",
				fmt.Sprintf("role %d: failed to unmarshal permissions from state; permissions: %s",
					id, statePermissionStr),
			)

			return

		}

		apiPermissionStr := apiState.Permissions.ValueString()
		err = json.Unmarshal([]byte(apiPermissionStr), &apiPermissionData)
		if err != nil {
			resp.Diagnostics.Append(pdiags...)
			resp.Diagnostics.AddError(
				"read role resource",
				fmt.Sprintf("role %d: failed to unmarshal permissions from api; permissions: %s",
					id, apiPermissionStr),
			)

			return

		}

		// If the existing state is a subset of the response from the API,
		// then we're safe to keep using the existing state as the new state.
		// Otherwise, it'll attempt to set the state to the API permissions which will show a mismatch.
		if eq, err := compare.ContainsSubset(apiPermissionData, statePermissionData); eq && err == nil {
			apiState.Permissions = state.Permissions
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &apiState)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *Resource) Update(
	_ context.Context,
	_ resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	resp.Diagnostics.AddError(
		"update role resource",
		"update of 'role' resources has not been implemented",
	)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var data RoleModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id.ValueInt64()
	client, _ := r.NewClient(ctx)
	_, hresp, err := client.RolesAPI.DeleteRole(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"delete role resource",
			fmt.Sprintf("role %d: DELETE failed ", id)+errors.ErrMsg(err, hresp),
		)

		return
	}
}

func (r *Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"import role resource",
			"provided import ID '"+req.ID+"' is invalid (non-number)",
		)

		return
	}

	diags := resp.State.SetAttribute(ctx, path.Root("id"), id)

	resp.Diagnostics.Append(diags...)
}

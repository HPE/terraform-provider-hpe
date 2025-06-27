// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package user

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/errors"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &Resource{}
	_ resource.ResourceWithImportState = &Resource{}
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
	resp.TypeName = req.ProviderTypeName + "_morpheus_user"
}

func (r *Resource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = UserResourceSchema(ctx)
}

// populate user resource model with current API values
func getUserAsState(
	ctx context.Context,
	id int64,
	client *sdk.APIClient,
) (UserModel, diag.Diagnostics) {
	var state UserModel
	var diags diag.Diagnostics

	u, hresp, err := client.UsersAPI.GetUser(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		diags.AddError(
			"populate user resource",
			fmt.Sprintf("user %d GET failed: ", id)+errors.ErrMsg(err, hresp),
		)

		return state, diags
	}

	roleIDValues := []attr.Value{}
	for _, role := range u.GetUser().Roles {
		roleIDValues = append(roleIDValues, convert.Int64ToType(role.Id))
	}

	roleIDSet, d := types.SetValue(types.Int64Type, roleIDValues)
	diags.Append(d...)
	if diags.HasError() {
		return state, diags
	}

	state.Id = convert.Int64ToType(u.User.Id)
	state.TenantId = convert.Int64ToType(u.User.AccountId)
	state.Username = convert.StrToType(u.User.Username)
	state.Email = convert.StrToType(u.User.Email)
	state.FirstName = convert.StrToType(u.User.FirstName)
	state.LastName = convert.StrToType(u.User.LastName)
	state.LinuxUsername = convert.StrToType(u.User.LinuxUsername.Get())
	state.WindowsUsername = convert.StrToType(u.User.WindowsUsername.Get())
	state.LinuxKeyPairId = convert.Int64ToType(u.User.LinuxKeyPairId.Get())
	state.PasswordExpired = convert.BoolToType(u.User.PasswordExpired)
	state.ReceiveNotifications = convert.BoolToType(u.User.ReceiveNotifications)
	state.RoleIds = roleIDSet

	return state, diags
}

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan UserModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var roleIDs []int64
	if !plan.RoleIds.IsNull() && !plan.RoleIds.IsUnknown() {
		diags := plan.RoleIds.ElementsAs(ctx, &roleIDs, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	var roles []sdk.GetAlerts200ResponseAllOfChecksInnerAccount
	for _, roleID := range roleIDs {
		rolevalue := sdk.GetAlerts200ResponseAllOfChecksInnerAccount{
			Id: &roleID,
		}
		roles = append(roles, rolevalue)
	}

	addUser := sdk.NewAddUserTenantRequestUserWithDefaults()

	var config UserModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// required
	username := plan.Username.ValueString()
	addUser.SetUsername(username)
	addUser.SetEmail(plan.Email.ValueString())
	addUser.SetRoles(roles)
	addUser.SetPassword(config.PasswordWo.ValueString())

	// optional
	if !plan.FirstName.IsUnknown() {
		addUser.SetFirstName(plan.FirstName.ValueString())
	}
	if !plan.LastName.IsUnknown() {
		addUser.SetLastName(plan.LastName.ValueString())
	}
	if !plan.LinuxUsername.IsUnknown() {
		addUser.SetLinuxUsername(plan.LinuxUsername.ValueString())
	}
	if !plan.LinuxPasswordWo.IsUnknown() {
		addUser.SetLinuxPassword(plan.LinuxPasswordWo.ValueString())
	}
	if !plan.WindowsUsername.IsUnknown() {
		addUser.SetWindowsUsername(plan.WindowsUsername.ValueString())
	}
	if !plan.WindowsPasswordWo.IsUnknown() {
		addUser.SetWindowsPassword(plan.WindowsPasswordWo.ValueString())
	}
	if !plan.LinuxKeyPairId.IsUnknown() {
		addUser.SetLinuxKeyPairId(plan.LinuxKeyPairId.ValueInt64())
	}
	if !plan.ReceiveNotifications.IsUnknown() {
		addUser.SetReceiveNotifications(plan.ReceiveNotifications.ValueBool())
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"create user resource",
			"user "+username+": failed to create client: "+err.Error(),
		)

		return
	}

	apiAddUserReq := client.UsersAPI.AddUser(ctx)
	if !plan.TenantId.IsUnknown() {
		apiAddUserReq = apiAddUserReq.AccountId(plan.TenantId.ValueInt64())
	}

	addUserReq := sdk.NewAddUserTenantRequest(*addUser)
	user, hresp, err := apiAddUserReq.AddUserTenantRequest(*addUserReq).Execute()

	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"create user resource",
			"user "+username+" POST failed: "+errors.ErrMsg(err, hresp),
		)

		return
	}

	if user.GetUser().Id == nil {
		resp.Diagnostics.AddError(
			"create user resource",
			"user "+username+": id is nil",
		)

		return
	}

	id := *user.GetUser().Id
	plan.Id = types.Int64Value(id)

	// write id as soon as possible
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, pdiags := getUserAsState(ctx, id, client)
	if pdiags.HasError() {
		resp.Diagnostics.Append(pdiags...)
		resp.Diagnostics.AddError(
			"create user resource",
			fmt.Sprintf("user %d: failed to read from api", id),
		)

		return
	}

	// special case - can't read from API
	state.PasswordWoVersion = plan.PasswordWoVersion
	state.WindowsPasswordWoVersion = plan.WindowsPasswordWoVersion
	state.LinuxPasswordWoVersion = plan.LinuxPasswordWoVersion

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Note that the following are not updateable via the API:
// LinuxUsername
// WindowsUsername
// LinuxKeyPairId
// ReceiveNotifications
// TenantId
func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state, config UserModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var roleIDs []int64
	if !plan.RoleIds.IsNull() && !plan.RoleIds.IsUnknown() {
		diags := plan.RoleIds.ElementsAs(ctx, &roleIDs, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	var roles []sdk.UpdateUserRequestUserRolesInner
	for _, roleID := range roleIDs {
		rolevalue := sdk.UpdateUserRequestUserRolesInner{
			Id: roleID,
		}
		roles = append(roles, rolevalue)
	}

	updateUser := sdk.NewUpdateUserRequestUser()

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	username := plan.Username.ValueString()

	// non-nullable
	updateUser.SetUsername(username)
	updateUser.SetEmail(plan.Email.ValueString())
	updateUser.SetRoles(roles)

	if !plan.PasswordWoVersion.Equal(state.PasswordWoVersion) {
		if config.PasswordWo.IsUnknown() {
			resp.Diagnostics.AddError(
				"update user resource",
				fmt.Sprintf("user %s: 'password_wo_version' changed, "+
					"but 'password_wo' is not set", username),
			)

			return
		}
		updateUser.SetPassword(config.PasswordWo.ValueString())
	}

	// nullable
	if plan.FirstName.IsNull() {
		updateUser.SetFirstNameNil()
	} else {
		updateUser.SetFirstName(plan.FirstName.ValueString())
	}

	if plan.LastName.IsNull() {
		updateUser.SetLastNameNil()
	} else {
		updateUser.SetLastName(plan.LastName.ValueString())
	}

	if plan.LinuxKeyPairId.IsNull() {
		updateUser.SetLinuxKeyPairIdNil()
	} else {
		updateUser.SetLinuxKeyPairId(plan.LinuxKeyPairId.ValueInt64())
	}

	if plan.LinuxUsername.IsNull() {
		updateUser.SetLinuxUsernameNil()
	} else {
		updateUser.SetLinuxUsername(plan.LinuxUsername.ValueString())
	}

	if plan.WindowsUsername.IsNull() {
		updateUser.SetWindowsUsernameNil()
	} else {
		updateUser.SetWindowsUsername(plan.WindowsUsername.ValueString())
	}

	if !plan.LinuxPasswordWoVersion.Equal(state.LinuxPasswordWoVersion) {
		if config.LinuxPasswordWo.IsUnknown() {
			resp.Diagnostics.AddError(
				"update user resource",
				fmt.Sprintf("user %s: 'linux_password_wo_version' changed, "+
					"but 'linux_password_wo' is not set", username),
			)

			return
		}
		if plan.LinuxPasswordWo.IsNull() {
			updateUser.SetLinuxPasswordNil()
		} else {
			updateUser.SetLinuxPassword(plan.LinuxPasswordWo.ValueString())
		}
	}

	if !plan.WindowsPasswordWoVersion.Equal(state.WindowsPasswordWoVersion) {
		if config.WindowsPasswordWo.IsUnknown() {
			resp.Diagnostics.AddError(
				"update user resource",
				fmt.Sprintf("user %s: 'windows_password_wo_version' changed, "+
					"but 'windows_password_wo' is not set", username),
			)

			return
		}
		if plan.WindowsPasswordWo.IsNull() {
			updateUser.SetWindowsPasswordNil()
		} else {
			updateUser.SetWindowsPassword(plan.WindowsPasswordWo.ValueString())
		}
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"update user resource",
			"user "+username+": failed to create client: "+err.Error(),
		)

		return
	}

	id := plan.Id.ValueInt64()
	apiUpdateUserReq := client.UsersAPI.UpdateUser(ctx, id)

	updateUserReq := sdk.NewUpdateUserRequest(*updateUser)
	user, hresp, err := apiUpdateUserReq.UpdateUserRequest(*updateUserReq).Execute()

	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"update user resource",
			"user "+username+" PUT failed: "+errors.ErrMsg(err, hresp),
		)

		return
	}

	if user.GetUser().Id == nil {
		resp.Diagnostics.AddError(
			"update user resource",
			"user "+username+": id is nil",
		)

		return
	}

	newid := *user.GetUser().Id
	if newid != id {
		resp.Diagnostics.AddError(
			"update user resource",
			"user "+username+": id mismatch "+fmt.Sprintf("%d != %d", id, newid),
		)

		return
	}

	state, pdiags := getUserAsState(ctx, newid, client)
	if pdiags.HasError() {
		resp.Diagnostics.Append(pdiags...)
		resp.Diagnostics.AddError(
			"update user resource",
			fmt.Sprintf("user %d: failed to read from api", id),
		)

		return
	}

	// special cases - can't read from API
	state.PasswordWoVersion = plan.PasswordWoVersion
	state.WindowsPasswordWoVersion = plan.WindowsPasswordWoVersion
	state.LinuxPasswordWoVersion = plan.LinuxPasswordWoVersion

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var plan UserModel

	diags := req.State.Get(ctx, &plan)
	if diags.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"read user resource",
			"new client call failed with "+err.Error(),
		)

		return
	}

	id := plan.Id.ValueInt64()
	state, pdiags := getUserAsState(ctx, id, client)
	if pdiags.HasError() {
		resp.Diagnostics.Append(pdiags...)
		resp.Diagnostics.AddError(
			"read user resource",
			fmt.Sprintf("user %d: failed to read from api", id),
		)

		return
	}

	// special cases - can't read from API
	state.PasswordWoVersion = plan.PasswordWoVersion
	state.WindowsPasswordWoVersion = plan.WindowsPasswordWoVersion
	state.LinuxPasswordWoVersion = plan.LinuxPasswordWoVersion

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var data UserModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id.ValueInt64()
	client, _ := r.NewClient(ctx)
	_, hresp, err := client.UsersAPI.DeleteUser(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"delete user resource",
			fmt.Sprintf("user %d: DELETE failed ", id)+errors.ErrMsg(err, hresp),
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
			"import user resource",
			"provided import ID '"+req.ID+"' is invalid (non-number)",
		)

		return
	}

	diags := resp.State.SetAttribute(
		ctx, path.Root("id"), id,
	)
	resp.Diagnostics.Append(diags...)
}

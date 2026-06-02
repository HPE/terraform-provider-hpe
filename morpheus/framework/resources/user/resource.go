// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package user

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
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

func getUser(
	ctx context.Context,
	id int64,
	client *sdk.APIClient,
) (*sdk.GetUser200Response, error) {
	u, hresp, err := client.UsersAPI.GetUser(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user %d GET failed: %s", id, errfmt.ErrMsg(err, hresp))
	}

	return u, nil
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
			fmt.Sprintf("user %d GET failed: ", id)+errfmt.ErrMsg(err, hresp),
		)

		return state, diags
	}

	roleIDValues := []attr.Value{}
	for _, role := range u.User.Roles {
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

	var roles []sdk.AddUserRequestUserRolesInner
	for _, roleID := range roleIDs {
		rolevalue := sdk.AddUserRequestUserRolesInner{
			Id: &roleID,
		}
		roles = append(roles, rolevalue)
	}

	addUser := sdk.NewAddUserRequestUserWithDefaults()

	var config UserModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// required
	username := plan.Username.ValueString()
	addUser.Username = username
	addUser.Email = plan.Email.ValueString()
	addUser.Roles = roles
	addUser.Password = config.PasswordWo.ValueString()

	// optional
	if !plan.FirstName.IsUnknown() {
		addUser.FirstName = sdk.PtrString(plan.FirstName.ValueString())
	}
	if !plan.LastName.IsUnknown() {
		addUser.LastName = sdk.PtrString(plan.LastName.ValueString())
	}
	if !plan.LinuxUsername.IsUnknown() {
		addUser.LinuxUsername = sdk.PtrString(plan.LinuxUsername.ValueString())
	}
	if !plan.LinuxPasswordWo.IsUnknown() {
		addUser.LinuxPassword = sdk.PtrString(plan.LinuxPasswordWo.ValueString())
	}
	if !plan.WindowsUsername.IsUnknown() {
		addUser.WindowsUsername = sdk.PtrString(plan.WindowsUsername.ValueString())
	}
	if !plan.WindowsPasswordWo.IsUnknown() {
		addUser.WindowsPassword = sdk.PtrString(plan.WindowsPasswordWo.ValueString())
	}
	if !plan.LinuxKeyPairId.IsUnknown() {
		addUser.LinuxKeyPairId = sdk.PtrInt64(plan.LinuxKeyPairId.ValueInt64())
	}
	if !plan.ReceiveNotifications.IsUnknown() {
		addUser.ReceiveNotifications = sdk.PtrBool(plan.ReceiveNotifications.ValueBool())
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

	addUserReq := &sdk.AddUserRequest{User: *addUser}
	user, hresp, err := apiAddUserReq.AddUserRequest(*addUserReq).Execute()

	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"create user resource",
			"user "+username+" POST failed: "+errfmt.ErrMsg(err, hresp),
		)

		return
	}

	if user.User == nil || user.User.Id == nil {
		resp.Diagnostics.AddError(
			"create user resource",
			"user "+username+": id is nil",
		)

		return
	}

	id := *user.User.Id
	plan.Id = types.Int64Value(id)

	// Helper to taint the resource state on an error after the POST request
	taintResourceState := func(id int64) {
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "user",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})
	}

	state, pdiags := getUserAsState(ctx, id, client)
	if pdiags.HasError() {
		resp.Diagnostics.Append(pdiags...)
		resp.Diagnostics.AddError(
			"create user resource",
			fmt.Sprintf("user %d: failed to read from api", id),
		)
		taintResourceState(id)

		return
	}

	// special case - can't read from API
	state.PasswordWoVersion = plan.PasswordWoVersion
	state.WindowsPasswordWoVersion = plan.WindowsPasswordWoVersion
	state.LinuxPasswordWoVersion = plan.LinuxPasswordWoVersion

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError(
			"failed to set user state",
			fmt.Sprintf("User %d was created but state could not be saved", id),
		)
		taintResourceState(id)

		return
	}
}

// Note that the following are not updateable via the API:
// LinuxUsername
// WindowsUsername
// LinuxKeyPairId
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

	updateUser := sdk.NewUpdateUserRequestUserWithDefaults()

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	username := plan.Username.ValueString()

	// non-nullable
	updateUser.Username = sdk.PtrString(username)
	updateUser.Email = sdk.PtrString(plan.Email.ValueString())
	updateUser.Roles = roles

	if !plan.PasswordWoVersion.Equal(state.PasswordWoVersion) {
		if config.PasswordWo.IsUnknown() {
			resp.Diagnostics.AddError(
				"update user resource",
				fmt.Sprintf("user %s: 'password_wo_version' changed, "+
					"but 'password_wo' is not set", username),
			)

			return
		}
		updateUser.Password = sdk.PtrString(config.PasswordWo.ValueString())
	}

	// nullable
	if plan.FirstName.IsNull() {
		updateUser.FirstName.Set(nil)
	} else {
		updateUser.FirstName.Set(sdk.PtrString(plan.FirstName.ValueString()))
	}

	if plan.LastName.IsNull() {
		updateUser.LastName.Set(nil)
	} else {
		updateUser.LastName.Set(sdk.PtrString(plan.LastName.ValueString()))
	}

	if plan.LinuxKeyPairId.IsUnknown() || plan.LinuxKeyPairId.IsNull() {
		updateUser.LinuxKeyPairId.Set(nil)
	} else {
		updateUser.LinuxKeyPairId.Set(sdk.PtrInt64(plan.LinuxKeyPairId.ValueInt64()))
	}

	if plan.LinuxUsername.IsNull() {
		updateUser.LinuxUsername.Set(nil)
	} else {
		updateUser.LinuxUsername.Set(sdk.PtrString(plan.LinuxUsername.ValueString()))
	}

	if plan.WindowsUsername.IsNull() {
		updateUser.WindowsUsername.Set(nil)
	} else {
		updateUser.WindowsUsername.Set(sdk.PtrString(plan.WindowsUsername.ValueString()))
	}

	updateUser.ReceiveNotifications = sdk.PtrBool(plan.ReceiveNotifications.ValueBool())

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
			updateUser.LinuxPassword.Set(nil)
		} else {
			updateUser.LinuxPassword.Set(sdk.PtrString(plan.LinuxPasswordWo.ValueString()))
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
			updateUser.WindowsPassword.Set(nil)
		} else {
			updateUser.WindowsPassword.Set(sdk.PtrString(plan.WindowsPasswordWo.ValueString()))
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

	originalUserState, err := getUser(ctx, id, client)
	if err != nil {
		resp.Diagnostics.AddError(
			"update user resource",
			fmt.Sprintf("user %d: failed to read from api", id),
		)

		return
	}

	apiUpdateUserReq := client.UsersAPI.UpdateUser(ctx, id)

	updateUserReq := &sdk.UpdateUserRequest{User: *updateUser}
	user, hresp, err := apiUpdateUserReq.UpdateUserRequest(*updateUserReq).Execute()

	if err != nil || hresp.StatusCode != http.StatusOK {
		if hresp.StatusCode != http.StatusInternalServerError {
			resp.Diagnostics.AddError(
				"update user resource",
				"user "+username+" PUT failed: "+errfmt.ErrMsg(err, hresp),
			)

			return
		}

		// Work around an API bug
		newUserState, err := getUser(ctx, id, client)
		if err != nil {
			resp.Diagnostics.AddError(
				"update user resource",
				fmt.Sprintf("user %d: failed to read from api", id),
			)

			return
		}

		if !newUserState.User.LastUpdated.After(*originalUserState.User.LastUpdated) {
			resp.Diagnostics.AddError(
				"update user resource",
				fmt.Sprintf("user %d: resource was not updated", id),
			)
		}

		innerUser := sdk.NewAddUser200ResponseAllOfUserWithDefaults()
		innerUser.Id = sdk.PtrInt64(id)
		user = &sdk.UpdateUser200Response{}
		user.User = innerUser
	}

	if user.User == nil || user.User.Id == nil {
		resp.Diagnostics.AddError(
			"update user resource",
			"user "+username+": id is nil",
		)

		return
	}

	newid := *user.User.Id
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
			fmt.Sprintf("user %d: DELETE failed ", id)+errfmt.ErrMsg(err, hresp),
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

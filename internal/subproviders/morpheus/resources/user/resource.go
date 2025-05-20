// (C) Copyright 2024 Hewlett Packard Enterprise Development LP

package user

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/configure"
	"github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	//"github.com/hashicorp/terraform-plugin-log/tflog"
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

func errMsg(err error, resp *http.Response) string {
	var msg string

	if err != nil {
		msg = err.Error()
	}

	if resp != nil {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return msg
		}
		code := http.StatusText(resp.StatusCode)
		msg = fmt.Sprintf("%s (%s): %s", msg, code, string(bodyBytes))
	}

	return msg
}

func strToType(s *string) types.String {
	if s == nil {
		return types.StringNull()
	}

	return types.StringValue(*s)
}

func boolToType(b *bool) types.Bool {
	if b == nil {
		return types.BoolNull()
	}

	return types.BoolValue(*b)
}

func int64ToType(i *int64) types.Int64 {
	if i == nil {
		return types.Int64Null()
	}

	return types.Int64Value(*i)
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
			fmt.Sprintf("user %d GET failed: ", id)+errMsg(err, hresp),
		)

		return state, diags
	}

	roleIDValues := []attr.Value{}
	for _, role := range u.GetUser().Roles {
		roleIDValues = append(roleIDValues, int64ToType(role.Id))
	}

	roleIDSet, d := types.SetValue(types.Int64Type, roleIDValues)
	diags.Append(d...)
	if diags.HasError() {
		return state, diags
	}

	state.Id = int64ToType(u.User.Id)
	state.Username = strToType(u.User.Username)
	state.Email = strToType(u.User.Email)
	state.FirstName = strToType(u.User.FirstName)
	state.LastName = strToType(u.User.LastName)
	state.LinuxUsername = strToType(u.User.LinuxUsername)
	state.WindowsUsername = strToType(u.User.WindowsUsername)
	state.LinuxKeyPairId = int64ToType(u.User.LinuxKeyPairId)
	state.PasswordExpired = boolToType(u.User.PasswordExpired)
	state.ReceiveNotifications = boolToType(u.User.ReceiveNotifications)
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

	// required
	username := plan.Username.ValueString()
	addUser.SetUsername(username)
	addUser.SetEmail(plan.Email.ValueString())
	addUser.SetPassword(plan.Password.ValueString())
	addUser.SetRoles(roles)

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
	if !plan.WindowsUsername.IsUnknown() {
		addUser.SetWindowsUsername(plan.WindowsUsername.ValueString())
	}
	if !plan.LinuxKeyPairId.IsUnknown() {
		addUser.SetLinuxKeyPairId(plan.LinuxKeyPairId.ValueInt64())
	}
	if !plan.ReceiveNotifications.IsUnknown() {
		addUser.SetReceiveNotifications(plan.ReceiveNotifications.ValueBool())
	}

	addUserReq := sdk.NewAddUserTenantRequest(*addUser)

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"create user resource",
			"user "+username+": failed to create client: "+err.Error(),
		)

		return
	}

	user, hresp, err := client.UsersAPI.AddUser(ctx).
		AddUserTenantRequest(*addUserReq).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"create user resource",
			"user "+username+" POST failed: "+errMsg(err, hresp),
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

	// special case (for now)
	state.Password, _ = plan.Password.ToStringValue(ctx)

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

	// special case (for now)
	state.Password, _ = plan.Password.ToStringValue(ctx)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
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
		"update user resource",
		"update of 'user' resources has not been implemented",
	)
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
			fmt.Sprintf("user %d: DELETE failed ", id)+errMsg(err, hresp),
		)

		return
	}
}

func (r *Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	parts := strings.SplitN(req.ID, ",", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"import user resource",
			"expected import format: <id>,<password>",
		)
	}
	password := parts[1]
	id, err := strconv.Atoi(parts[0])
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
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.SetAttribute(
		ctx, path.Root("password"), password,
	)
	resp.Diagnostics.Append(diags...)
}

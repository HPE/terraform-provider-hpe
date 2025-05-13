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
	sdk "github.com/HewlettPackard/hpe-morpheus-client/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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

	username := plan.Username.ValueString()
	email := plan.Email.ValueString()
	password := plan.Password.ValueString()

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

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"create user resource",
			"user "+username+": failed to create client: "+err.Error(),
		)

		return
	}

	addUser := sdk.NewAddUserTenantRequestUserWithDefaults()
	addUser.SetEmail(email)
	addUser.SetUsername(username)
	addUser.SetPassword(password)
	addUser.SetRoles(roles)

	addUserReq := sdk.NewAddUserTenantRequest(*addUser)

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
	u, hresp, err := client.UsersAPI.GetUser(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"create user resource",
			"user "+username+" GET failed: "+errMsg(err, hresp),
		)

		return
	}

	roleIDValues := []attr.Value{}
	for _, role := range u.GetUser().Roles {
		roleIDValues = append(roleIDValues, types.Int64Value(*role.Id))
	}

	roleIDSet, diags := types.SetValue(types.Int64Type, roleIDValues)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state UserModel

	state.Id = types.Int64Value(id)
	state.Username = types.StringValue(username)
	state.Email = types.StringValue(*u.GetUser().Email)
	state.Password = types.StringValue(plan.Password.ValueString())
	state.RoleIds = roleIDSet

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

	id := plan.Id.ValueInt64()

	var roleIDs []int64
	if !plan.RoleIds.IsNull() && !plan.RoleIds.IsUnknown() {
		diags := plan.RoleIds.ElementsAs(ctx, &roleIDs, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			tflog.Error(ctx, fmt.Sprintf("Could not parse roles for user %d", id))

			return
		}
	}

	client, _ := r.NewClient(ctx)
	u, hresp, err := client.UsersAPI.GetUser(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"read user resource",
			fmt.Sprintf("user %d GET failed: ", id)+errMsg(err, hresp),
		)

		return
	}

	username := u.GetUser().Username
	if username == nil {
		resp.Diagnostics.AddError(
			"read user resource",
			fmt.Sprintf("user %d has nil name: ", id),
		)

		return
	}

	roleIDValues := []attr.Value{}
	for _, role := range u.GetUser().Roles {
		roleIDValues = append(roleIDValues, types.Int64Value(*role.Id))
	}

	roleIDSet, diags := types.SetValue(types.Int64Type, roleIDValues)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state UserModel

	state.Id = types.Int64Value(id)
	state.Username = types.StringValue(*username)
	state.Email = types.StringValue(*u.GetUser().Email)
	state.Password = types.StringValue(plan.Password.ValueString())
	state.RoleIds = roleIDSet

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *Resource) Update(
	ctx context.Context,
	_ resource.UpdateRequest,
	_ *resource.UpdateResponse,
) {
	tflog.Error(ctx, "update 'user' is not implemented")
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

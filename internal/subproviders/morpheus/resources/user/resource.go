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
	msg := err.Error()
	if resp != nil {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return msg
		}
		msg = msg + " http response: " + string(bodyBytes)

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

	var rolesList []RolesValue
	resp.Diagnostics.Append(plan.Roles.ElementsAs(ctx, &rolesList, false)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Could not parse roles for user "+username)

		return
	}

	var roles []sdk.GetAlerts200ResponseAllOfChecksInnerAccount
	for _, role := range rolesList {
		if role.Id.IsUnknown() {
			resp.Diagnostics.AddError(
				"create user error",
				username+": missing role id",
			)

			return
		}

		rolevalue := sdk.GetAlerts200ResponseAllOfChecksInnerAccount{
			Id: role.Id.ValueInt64Pointer(),
		}
		roles = append(roles, rolevalue)
	}
	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"create user error",
			username+": failed to create client",
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
	if err != nil {
		resp.Diagnostics.AddError(
			"create user error",
			username+" post failed: "+errMsg(err, hresp),
		)

		return
	}

	if user.GetUser().Id == nil {
		resp.Diagnostics.AddError(
			"create user error",
			username+": id is nil",
		)

		return

	}

	id := *user.GetUser().Id
	u, hresp, err := client.UsersAPI.GetUser(ctx, id).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"create user error",
			username+" get failed: "+errMsg(err, hresp),
		)

		return
	}
	roleValues := []attr.Value{}
	uroles := u.GetUser().Roles
	for _, role := range uroles {
		value := RolesValue{
			Authority:   types.StringValue(role.GetAuthority()),
			Description: types.StringValue(role.GetDescription()),
			Id:          types.Int64Value(role.GetId()),
			Name:        types.StringValue(role.GetName()),
			state:       attr.ValueStateKnown,
		}
		roleValues = append(roleValues, value)
	}

	var state UserModel

	state.Id = types.Int64Value(id)
	state.Username = types.StringValue(*u.GetUser().Username)
	state.Email = types.StringValue(*u.GetUser().Email)
	state.Roles = types.SetValueMust(RolesValue{}.Type(ctx), roleValues)
	state.Password = types.StringValue(password)

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
	// password := plan.Password.ValueString()

	id := plan.Id.ValueInt64()

	var rolesList []RolesValue
	resp.Diagnostics.Append(plan.Roles.ElementsAs(ctx, &rolesList, false)...)

	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Could not parse roles for user %d", id))

		return
	}
	client, _ := r.NewClient(ctx)
	u, hresp, err := client.UsersAPI.GetUser(ctx, id).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"read user error",
			fmt.Sprintf("%d get failed: ", id)+errMsg(err, hresp),
		)

		return
	}

	username := u.GetUser().Username
	if username == nil {
		resp.Diagnostics.AddError(
			"read user error",
			fmt.Sprintf("%d has nil name: ", id),
		)

		return

	}

	roleValues := []attr.Value{}
	for _, apirole := range u.GetUser().Roles {
		if apirole.Id == nil {
			resp.Diagnostics.AddError(
				"read user error",
				*username+" has missing role id",
			)

			return

		}
		value := RolesValue{
			Authority:   types.StringValue(apirole.GetAuthority()),
			Description: types.StringValue(apirole.GetDescription()),
			Id:          types.Int64Value(apirole.GetId()),
			Name:        types.StringValue(apirole.GetName()),
			state:       attr.ValueStateKnown,
		}
		roleValues = append(roleValues, value)

	}
	var state UserModel

	state.Id = types.Int64Value(id)
	state.Username = types.StringValue(*username)
	state.Email = types.StringValue(*u.GetUser().Email)
	state.Roles = types.SetValueMust(RolesValue{}.Type(ctx), roleValues)
	state.Password = types.StringValue(plan.Password.ValueString())

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
	_, _, err := client.UsersAPI.DeleteUser(ctx, id).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"delete user error",
			fmt.Sprintf("%d: delete failed ", id)+err.Error(),
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
			"import user error",
			"expected import format: <id>,<password>",
		)
	}
	password := parts[1]
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		resp.Diagnostics.AddError(
			"import error",
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

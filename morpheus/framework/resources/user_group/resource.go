package user_group

import (
	"context"
	"fmt"
	"strconv"

	sdk "github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

var (
	_ resource.Resource                = &userGroupResource{}
	_ resource.ResourceWithConfigure   = &userGroupResource{}
	_ resource.ResourceWithImportState = &userGroupResource{}
)

type userGroupResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &userGroupResource{}
}

func (r *userGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_user_group"
}

func (r *userGroupResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = UserGroupSchema(ctx)
}

func (r *userGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan userGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := sdk.AddUserGroupRequestUserGroup{
		Name: plan.Name.ValueStringPointer(),
	}
	if !plan.Description.IsNull() {
		body.Description = *sdk.NewNullableString(plan.Description.ValueStringPointer())
	}
	if !plan.SudoUser.IsNull() {
		body.SudoUser = plan.SudoUser.ValueBoolPointer()
	}
	if !plan.ServerGroup.IsNull() {
		body.ServerGroup = *sdk.NewNullableString(plan.ServerGroup.ValueStringPointer())
	}

	result, httpResp, err := client.UsersAPI.AddUserGroup(ctx).AddUserGroupRequest(sdk.AddUserGroupRequest{
		UserGroup: body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "user_group", plan.Name.ValueString(), err, httpResp)
		return
	}

	ug := result.GetUserGroup()
	mapAddResponseToModel(&plan, &ug)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *userGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state userGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	result, httpResp, err := client.UsersAPI.GetUserGroup(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "user_group", "", err, httpResp)
		return
	}

	ug := result.GetUserGroup()
	mapGetResponseToModel(&state, &ug)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *userGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan userGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueInt64()

	body := sdk.UpdateUserGroupRequestUserGroup{
		Name: plan.Name.ValueStringPointer(),
	}
	if !plan.Description.IsNull() {
		body.Description = *sdk.NewNullableString(plan.Description.ValueStringPointer())
	}
	if !plan.SudoUser.IsNull() {
		body.SudoUser = plan.SudoUser.ValueBoolPointer()
	}
	if !plan.ServerGroup.IsNull() {
		body.ServerGroup = *sdk.NewNullableString(plan.ServerGroup.ValueStringPointer())
	}

	result, httpResp, err := client.UsersAPI.UpdateUserGroup(ctx, id).UpdateUserGroupRequest(sdk.UpdateUserGroupRequest{
		UserGroup: body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "user_group", plan.Name.ValueString(), err, httpResp)
		return
	}

	ug := result.GetUserGroup()
	mapUpdateResponseToModel(&plan, &ug)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *userGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state userGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	_, httpResp, err := client.UsersAPI.DeleteUserGroup(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "user_group", "", err, httpResp)
		return
	}
}

func (r *userGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse ID %q as integer: %s", req.ID, err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func mapAddResponseToModel(model *userGroupModel, ug *sdk.AddUserGroup200ResponseAllOfUserGroup) {
	if ug.Id != nil {
		model.ID = types.Int64Value(*ug.Id)
	}
	if ug.Name != nil {
		model.Name = types.StringValue(*ug.Name)
	}
	if v := ug.Description.Get(); v != nil {
		model.Description = types.StringValue(*v)
	} else {
		model.Description = types.StringNull()
	}
	if ug.SudoUser != nil {
		model.SudoUser = types.BoolValue(*ug.SudoUser)
	}
	if v := ug.ServerGroup.Get(); v != nil {
		model.ServerGroup = types.StringValue(*v)
	} else {
		model.ServerGroup = types.StringNull()
	}
}

func mapGetResponseToModel(model *userGroupModel, ug *sdk.GetUserGroup200ResponseUserGroup) {
	if ug.Id != nil {
		model.ID = types.Int64Value(*ug.Id)
	}
	if ug.Name != nil {
		model.Name = types.StringValue(*ug.Name)
	}
	if v := ug.Description.Get(); v != nil {
		model.Description = types.StringValue(*v)
	} else {
		model.Description = types.StringNull()
	}
	if ug.SudoUser != nil {
		model.SudoUser = types.BoolValue(*ug.SudoUser)
	}
	if v := ug.ServerGroup.Get(); v != nil {
		model.ServerGroup = types.StringValue(*v)
	} else {
		model.ServerGroup = types.StringNull()
	}
}

func mapUpdateResponseToModel(model *userGroupModel, ug *sdk.UpdateUserGroup200ResponseAllOfUserGroup) {
	if ug.Id != nil {
		model.ID = types.Int64Value(*ug.Id)
	}
	if ug.Name != nil {
		model.Name = types.StringValue(*ug.Name)
	}
	if v := ug.Description.Get(); v != nil {
		model.Description = types.StringValue(*v)
	} else {
		model.Description = types.StringNull()
	}
	if ug.SudoUser != nil {
		model.SudoUser = types.BoolValue(*ug.SudoUser)
	}
	if v := ug.ServerGroup.Get(); v != nil {
		model.ServerGroup = types.StringValue(*v)
	} else {
		model.ServerGroup = types.StringNull()
	}
}

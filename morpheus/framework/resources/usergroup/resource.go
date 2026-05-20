package usergroup

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/polling"
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

	// TODO: Build request body from plan and call API
	// result, httpResp, err := client.UsersAPI.AddUserGroup(ctx).Body(body).Execute()
	// if err := errfmt.CheckResponse(err, httpResp); err != nil {
	//     errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "user_group", plan.Name.ValueString(), err, httpResp)
	//     return
	// }
	_ = client
	_ = polling.ForCreate

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

	// TODO: Call API to read resource
	// result, httpResp, err := client.UsersAPI.GetUserGroup(ctx, id).Execute()
	var httpResp *http.Response
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)
		return
	}
	// if err := errfmt.CheckResponse(err, httpResp); err != nil {
	//     errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "user_group", "", err, httpResp)
	//     return
	// }
	_ = client
	_ = id

	// TODO: Map response to state model
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

	// TODO: Build request body from plan and call API
	// _, httpResp, err := client.UsersAPI.UpdateUserGroup(ctx, id).Body(body).Execute()
	// if err := errfmt.CheckResponse(err, httpResp); err != nil {
	//     errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "user_group", plan.Name.ValueString(), err, httpResp)
	//     return
	// }
	_ = client
	_ = id

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

	// TODO: Call API to delete resource
	// httpResp, err := client.UsersAPI.DeleteUserGroup(ctx, id).Execute()
	// if err := errfmt.CheckResponse(err, httpResp); err != nil {
	//     errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "user_group", "", err, httpResp)
	//     return
	// }
	_ = client
	_ = id
	_ = fmt.Sprintf
	_ = polling.ForDelete
}

func (r *userGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse ID %q as integer: %s", req.ID, err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

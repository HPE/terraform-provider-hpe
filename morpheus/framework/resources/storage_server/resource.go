package storage_server

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
	_ resource.Resource                = &storageServerResource{}
	_ resource.ResourceWithConfigure   = &storageServerResource{}
	_ resource.ResourceWithImportState = &storageServerResource{}
)

type storageServerResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &storageServerResource{}
}

func (r *storageServerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_storage_server"
}

func (r *storageServerResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = StorageServerSchema(ctx)
}

func (r *storageServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan storageServerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config := map[string]interface{}{}
	if !plan.ServiceUrl.IsNull() {
		config["serviceUrl"] = plan.ServiceUrl.ValueString()
	}
	if !plan.ServiceUsername.IsNull() {
		config["serviceUsername"] = plan.ServiceUsername.ValueString()
	}
	if !plan.ServicePassword.IsNull() {
		config["servicePassword"] = plan.ServicePassword.ValueString()
	}

	body := sdk.AddStorageServersRequestStorageServer{
		Name:   plan.Name.ValueString(),
		Type:   plan.Type.ValueString(),
		Config: config,
	}
	if !plan.Description.IsNull() {
		body.Description = plan.Description.ValueStringPointer()
	}
	if !plan.Enabled.IsNull() {
		body.Enabled = plan.Enabled.ValueBoolPointer()
	}

	result, httpResp, err := client.StorageAPI.AddStorageServers(ctx).AddStorageServersRequest(sdk.AddStorageServersRequest{
		StorageServer: body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "storage_server", plan.Name.ValueString(), err, httpResp)
		return
	}

	ss := result.GetStorageServer()
	mapCreateResponseToModel(&plan, &ss)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *storageServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state storageServerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	result, httpResp, err := client.StorageAPI.GetStorageServers(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "storage_server", "", err, httpResp)
		return
	}

	ss := result.GetStorageServer()
	mapGetResponseToModel(&state, &ss)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *storageServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan storageServerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueInt64()

	config := map[string]interface{}{}
	if !plan.ServiceUrl.IsNull() {
		config["serviceUrl"] = plan.ServiceUrl.ValueString()
	}
	if !plan.ServiceUsername.IsNull() {
		config["serviceUsername"] = plan.ServiceUsername.ValueString()
	}
	if !plan.ServicePassword.IsNull() {
		config["servicePassword"] = plan.ServicePassword.ValueString()
	}

	body := sdk.UpdateStorageServersRequestStorageServer{
		Name:   plan.Name.ValueStringPointer(),
		Type:   plan.Type.ValueStringPointer(),
		Config: config,
	}
	if !plan.Description.IsNull() {
		body.Description = plan.Description.ValueStringPointer()
	}
	if !plan.Enabled.IsNull() {
		body.Enabled = plan.Enabled.ValueBoolPointer()
	}

	_, httpResp, err := client.StorageAPI.UpdateStorageServers(ctx, id).UpdateStorageServersRequest(sdk.UpdateStorageServersRequest{
		StorageServer: body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "storage_server", plan.Name.ValueString(), err, httpResp)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *storageServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state storageServerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	_, httpResp, err := client.StorageAPI.RemoveStorageServers(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "storage_server", "", err, httpResp)
		return
	}
}

func (r *storageServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse ID %q as integer: %s", req.ID, err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func mapCreateResponseToModel(model *storageServerModel, ss *sdk.AddStorageServers200ResponseAllOfStorageServer) {
	if ss.Id != nil {
		model.ID = types.Int64Value(*ss.Id)
	}
	if ss.Name != nil {
		model.Name = types.StringValue(*ss.Name)
	}
	if ss.Enabled != nil {
		model.Enabled = types.BoolValue(*ss.Enabled)
	}
	if v := ss.Description.Get(); v != nil {
		model.Description = types.StringValue(*v)
	} else {
		model.Description = types.StringNull()
	}
	if v := ss.ServiceUrl.Get(); v != nil {
		model.ServiceUrl = types.StringValue(*v)
	} else {
		model.ServiceUrl = types.StringNull()
	}
	// Sensitive fields: preserve plan values (API won't return them)
}

func mapGetResponseToModel(model *storageServerModel, ss *sdk.GetStorageServers200ResponseStorageServer) {
	if ss.Id != nil {
		model.ID = types.Int64Value(*ss.Id)
	}
	if ss.Name != nil {
		model.Name = types.StringValue(*ss.Name)
	}
	if ss.Enabled != nil {
		model.Enabled = types.BoolValue(*ss.Enabled)
	}
	if v := ss.Description.Get(); v != nil {
		model.Description = types.StringValue(*v)
	} else {
		model.Description = types.StringNull()
	}
	if v := ss.ServiceUrl.Get(); v != nil {
		model.ServiceUrl = types.StringValue(*v)
	} else {
		model.ServiceUrl = types.StringNull()
	}
	// Sensitive fields (service_username, service_password): preserve state values
}

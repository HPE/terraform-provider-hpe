package network_pool_server

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
	_ resource.Resource                = &networkPoolServerResource{}
	_ resource.ResourceWithConfigure   = &networkPoolServerResource{}
	_ resource.ResourceWithImportState = &networkPoolServerResource{}
)

type networkPoolServerResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &networkPoolServerResource{}
}

func (r *networkPoolServerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_network_pool_server"
}

func (r *networkPoolServerResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = NetworkPoolServerSchema(ctx)
}

func (r *networkPoolServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan networkPoolServerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The SDK uses a oneOf union for this request. We use InfobloxNetworkPoolServer as the
	// concrete type since all pool server types share the same common fields.
	// TODO: Support multiple pool server types based on TypeID.
	infoblox := sdk.InfobloxNetworkPoolServer{
		Type: "infoblox",
		Name: plan.Name.ValueString(),
	}
	if !plan.ServiceUrl.IsNull() {
		infoblox.ServiceUrl = *sdk.NewNullableString(plan.ServiceUrl.ValueStringPointer())
	}
	if !plan.ServiceUsername.IsNull() {
		infoblox.ServiceUsername = *sdk.NewNullableString(plan.ServiceUsername.ValueStringPointer())
	}
	if !plan.ServicePassword.IsNull() {
		infoblox.ServicePassword = *sdk.NewNullableString(plan.ServicePassword.ValueStringPointer())
	}
	if !plan.IgnoreSsl.IsNull() {
		infoblox.IgnoreSsl = plan.IgnoreSsl.ValueBoolPointer()
	}
	if !plan.Enabled.IsNull() {
		infoblox.Enabled = plan.Enabled.ValueBoolPointer()
	}

	serverReq := sdk.CreateNetworkPoolServerRequestNetworkPoolServer{
		InfobloxNetworkPoolServer: &infoblox,
	}

	// TODO: The SDK uses a oneOf union for this request. If the above struct does not compile,
	// the request may need to be constructed differently depending on the IPAM type.
	result, httpResp, err := client.NetworksAPI.CreateNetworkPoolServer(ctx).CreateNetworkPoolServerRequest(sdk.CreateNetworkPoolServerRequest{
		NetworkPoolServer: &serverReq,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "network_pool_server", plan.Name.ValueString(), err, httpResp)
		return
	}

	server := result.GetNetworkPoolServer()
	mapCreateResponseToModel(&plan, &server)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *networkPoolServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state networkPoolServerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	result, httpResp, err := client.NetworksAPI.GetNetworkPoolServer(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "network_pool_server", "", err, httpResp)
		return
	}

	server := result.GetNetworkPoolServer()
	mapReadResponseToModel(&state, &server)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *networkPoolServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan networkPoolServerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueInt64()

	infobloxUpdate := sdk.InfobloxNetworkPoolServerUpdate{
		Name: plan.Name.ValueStringPointer(),
	}
	if !plan.ServiceUrl.IsNull() {
		infobloxUpdate.ServiceUrl = *sdk.NewNullableString(plan.ServiceUrl.ValueStringPointer())
	}
	if !plan.ServiceUsername.IsNull() {
		infobloxUpdate.ServiceUsername = *sdk.NewNullableString(plan.ServiceUsername.ValueStringPointer())
	}
	if !plan.ServicePassword.IsNull() {
		infobloxUpdate.ServicePassword = *sdk.NewNullableString(plan.ServicePassword.ValueStringPointer())
	}
	if !plan.IgnoreSsl.IsNull() {
		infobloxUpdate.IgnoreSsl = plan.IgnoreSsl.ValueBoolPointer()
	}
	if !plan.Enabled.IsNull() {
		infobloxUpdate.Enabled = plan.Enabled.ValueBoolPointer()
	}

	serverReq := sdk.InfobloxNetworkPoolServerUpdateAsUpdateNetworkPoolServerRequestNetworkPoolServer(&infobloxUpdate)

	_, httpResp, err := client.NetworksAPI.UpdateNetworkPoolServer(ctx, id).UpdateNetworkPoolServerRequest(sdk.UpdateNetworkPoolServerRequest{
		NetworkPoolServer: &serverReq,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "network_pool_server", plan.Name.ValueString(), err, httpResp)
		return
	}

	// Re-read to get current state
	readResult, httpResp, err := client.NetworksAPI.GetNetworkPoolServer(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "network_pool_server", "", err, httpResp)
		return
	}

	server := readResult.GetNetworkPoolServer()
	mapReadResponseToModel(&plan, &server)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *networkPoolServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state networkPoolServerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	_, httpResp, err := client.NetworksAPI.DeleteNetworkPoolServer(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "network_pool_server", "", err, httpResp)
		return
	}
}

func (r *networkPoolServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse ID %q as integer: %s", req.ID, err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func mapCreateResponseToModel(model *networkPoolServerModel, server *sdk.CreateNetworkPoolServer200ResponseAllOfNetworkPoolServer) {
	if server.Id != nil {
		model.ID = types.Int64Value(*server.Id)
	}
	if server.Name != nil {
		model.Name = types.StringValue(*server.Name)
	}
	if t := server.Type; t != nil && t.Id != nil {
		model.TypeID = types.Int64Value(*t.Id)
	}
	if server.ServiceUrl.IsSet() && server.ServiceUrl.Get() != nil {
		model.ServiceUrl = types.StringValue(*server.ServiceUrl.Get())
	}
	if server.IgnoreSsl.IsSet() && server.IgnoreSsl.Get() != nil {
		model.IgnoreSsl = types.BoolValue(*server.IgnoreSsl.Get())
	}
	if server.Enabled != nil {
		model.Enabled = types.BoolValue(*server.Enabled)
	}
	if server.Status != nil {
		model.Status = types.StringValue(*server.Status)
	} else {
		model.Status = types.StringNull()
	}
}

func mapReadResponseToModel(model *networkPoolServerModel, server *sdk.GetNetworkPoolServer200ResponseNetworkPoolServer) {
	if server.Id != nil {
		model.ID = types.Int64Value(*server.Id)
	}
	if server.Name != nil {
		model.Name = types.StringValue(*server.Name)
	}
	if t := server.Type; t != nil && t.Id != nil {
		model.TypeID = types.Int64Value(*t.Id)
	}
	if server.ServiceUrl.IsSet() && server.ServiceUrl.Get() != nil {
		model.ServiceUrl = types.StringValue(*server.ServiceUrl.Get())
	} else {
		model.ServiceUrl = types.StringNull()
	}
	if server.IgnoreSsl.IsSet() && server.IgnoreSsl.Get() != nil {
		model.IgnoreSsl = types.BoolValue(*server.IgnoreSsl.Get())
	}
	if server.Enabled != nil {
		model.Enabled = types.BoolValue(*server.Enabled)
	}
	if server.Status != nil {
		model.Status = types.StringValue(*server.Status)
	} else {
		model.Status = types.StringNull()
	}
	// Note: service_username and service_password are write-only and not returned by the API.
	// We preserve the plan values for those fields by not overwriting them here.
}

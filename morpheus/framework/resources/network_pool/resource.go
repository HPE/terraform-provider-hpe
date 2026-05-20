package network_pool

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
	_ resource.Resource                = &networkPoolResource{}
	_ resource.ResourceWithConfigure   = &networkPoolResource{}
	_ resource.ResourceWithImportState = &networkPoolResource{}
)

type networkPoolResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &networkPoolResource{}
}

func (r *networkPoolResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_network_pool"
}

func (r *networkPoolResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = NetworkPoolSchema(ctx)
}

func (r *networkPoolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan networkPoolModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	poolReq := sdk.CreateNetworkPoolRequestNetworkPool{
		Name: plan.Name.ValueStringPointer(),
		Type: &sdk.CreateNetworkPoolRequestNetworkPoolType{
			AdditionalProperties: map[string]interface{}{
				"id": strconv.FormatInt(plan.TypeID.ValueInt64(), 10),
			},
		},
	}

	result, httpResp, err := client.NetworksAPI.CreateNetworkPool(ctx).CreateNetworkPoolRequest(sdk.CreateNetworkPoolRequest{
		NetworkPool: &poolReq,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "network_pool", plan.Name.ValueString(), err, httpResp)
		return
	}

	pool := result.GetNetworkPool()
	mapCreateResponseToModel(&plan, &pool)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *networkPoolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state networkPoolModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	result, httpResp, err := client.NetworksAPI.GetNetworkPool(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "network_pool", "", err, httpResp)
		return
	}

	pool := result.GetNetworkPool()
	mapReadResponseToModel(&state, &pool)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *networkPoolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan networkPoolModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueInt64()

	poolReq := sdk.UpdateNetworkPoolRequestNetworkPool{
		Name: plan.Name.ValueStringPointer(),
	}

	_, httpResp, err := client.NetworksAPI.UpdateNetworkPool(ctx, id).UpdateNetworkPoolRequest(sdk.UpdateNetworkPoolRequest{
		NetworkPool: &poolReq,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "network_pool", plan.Name.ValueString(), err, httpResp)
		return
	}

	// Re-read to get current state
	readResult, httpResp, err := client.NetworksAPI.GetNetworkPool(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "network_pool", "", err, httpResp)
		return
	}

	pool := readResult.GetNetworkPool()
	mapReadResponseToModel(&plan, &pool)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *networkPoolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state networkPoolModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	_, httpResp, err := client.NetworksAPI.DeleteNetworkPool(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "network_pool", "", err, httpResp)
		return
	}
}

func (r *networkPoolResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse ID %q as integer: %s", req.ID, err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func mapCreateResponseToModel(model *networkPoolModel, pool *sdk.CreateNetworkPool200ResponseNetworkPool) {
	if pool.Id != nil {
		model.ID = types.Int64Value(*pool.Id)
	}
	if pool.Name != nil {
		model.Name = types.StringValue(*pool.Name)
	}
	if pool.IpCount != nil {
		model.IpCount = types.Int64Value(*pool.IpCount)
	}
	if pool.FreeCount != nil {
		model.FreeCount = types.Int64Value(*pool.FreeCount)
	}
	if pool.PoolEnabled != nil {
		model.PoolEnabled = types.BoolValue(*pool.PoolEnabled)
	}
	if pool.DhcpServer != nil {
		model.DhcpServer = types.BoolValue(*pool.DhcpServer)
	}
	if pool.DnsDomain.IsSet() && pool.DnsDomain.Get() != nil {
		model.DnsDomain = types.StringValue(*pool.DnsDomain.Get())
	} else {
		model.DnsDomain = types.StringNull()
	}
	if pool.Gateway.IsSet() && pool.Gateway.Get() != nil {
		model.Gateway = types.StringValue(*pool.Gateway.Get())
	} else {
		model.Gateway = types.StringNull()
	}
	if pool.Netmask.IsSet() && pool.Netmask.Get() != nil {
		model.Netmask = types.StringValue(*pool.Netmask.Get())
	} else {
		model.Netmask = types.StringNull()
	}
	if pool.SubnetAddress.IsSet() && pool.SubnetAddress.Get() != nil {
		model.SubnetAddress = types.StringValue(*pool.SubnetAddress.Get())
	} else {
		model.SubnetAddress = types.StringNull()
	}
}

func mapReadResponseToModel(model *networkPoolModel, pool *sdk.GetNetworkPool200ResponseNetworkPool) {
	if pool.Id != nil {
		model.ID = types.Int64Value(*pool.Id)
	}
	if pool.Name != nil {
		model.Name = types.StringValue(*pool.Name)
	}
	if t := pool.Type; t != nil && t.Id != nil {
		model.TypeID = types.Int64Value(*t.Id)
	}
	if pool.IpCount != nil {
		model.IpCount = types.Int64Value(*pool.IpCount)
	}
	if pool.FreeCount != nil {
		model.FreeCount = types.Int64Value(*pool.FreeCount)
	}
	if pool.PoolEnabled != nil {
		model.PoolEnabled = types.BoolValue(*pool.PoolEnabled)
	}
	if v, ok := pool.GetDnsDomainOk(); ok && v != nil {
		model.DnsDomain = types.StringValue(*v)
	} else {
		model.DnsDomain = types.StringNull()
	}
	if pool.DhcpServer != nil {
		model.DhcpServer = types.BoolValue(*pool.DhcpServer)
	}
	if v, ok := pool.GetGatewayOk(); ok && v != nil {
		model.Gateway = types.StringValue(*v)
	} else {
		model.Gateway = types.StringNull()
	}
	if v, ok := pool.GetNetmaskOk(); ok && v != nil {
		model.Netmask = types.StringValue(*v)
	} else {
		model.Netmask = types.StringNull()
	}
	if v, ok := pool.GetSubnetAddressOk(); ok && v != nil {
		model.SubnetAddress = types.StringValue(*v)
	} else {
		model.SubnetAddress = types.StringNull()
	}
}

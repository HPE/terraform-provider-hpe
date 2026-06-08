package network_pool

import (
	"context"
	"fmt"
	"strconv"

	sdk "github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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

func (r *networkPoolResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
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
				"code": plan.TypeCode.ValueString(),
			},
		},
	}

	additionalProps := map[string]interface{}{}
	if !plan.PoolEnabled.IsNull() && !plan.PoolEnabled.IsUnknown() {
		additionalProps["poolEnabled"] = plan.PoolEnabled.ValueBool()
	}
	if !plan.DNSDomain.IsNull() && !plan.DNSDomain.IsUnknown() {
		additionalProps["dnsDomain"] = plan.DNSDomain.ValueString()
	}
	if !plan.DhcpServer.IsNull() && !plan.DhcpServer.IsUnknown() {
		additionalProps["dhcpServer"] = plan.DhcpServer.ValueBool()
	}
	if !plan.Gateway.IsNull() && !plan.Gateway.IsUnknown() {
		additionalProps["gateway"] = plan.Gateway.ValueString()
	}
	if !plan.Netmask.IsNull() && !plan.Netmask.IsUnknown() {
		additionalProps["netmask"] = plan.Netmask.ValueString()
	}
	if !plan.SubnetAddress.IsNull() && !plan.SubnetAddress.IsUnknown() {
		additionalProps["subnetAddress"] = plan.SubnetAddress.ValueString()
	}
	if !plan.IpRanges.IsNull() && !plan.IpRanges.IsUnknown() {
		attrs := plan.IpRanges.Attributes()
		startAddr := attrs["starting_address"].(types.String).ValueString()
		endAddr := attrs["ending_address"].(types.String).ValueString()
		additionalProps["ipRanges"] = []map[string]interface{}{
			{"startAddress": startAddr, "endAddress": endAddr},
		}
	}
	if len(additionalProps) > 0 {
		poolReq.AdditionalProperties = additionalProps
	}

	result, httpResp, err := client.NetworksAPI.CreateNetworkPool(ctx).
		CreateNetworkPoolRequest(sdk.CreateNetworkPoolRequest{
			NetworkPool: &poolReq,
		}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "network_pool", plan.Name.ValueString(), err, httpResp)

		return
	}

	pool := result.NetworkPool
	if pool == nil {
		resp.Diagnostics.AddError("API returned nil", "NetworkPool is nil in the response")

		return
	}
	mapCreateResponseToModel(&plan, pool)

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

	pool := result.NetworkPool
	if pool == nil {
		resp.Diagnostics.AddError("API returned nil", "NetworkPool is nil in the response")

		return
	}
	mapReadResponseToModel(&state, pool)

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

	additionalProps := map[string]interface{}{}
	if !plan.PoolEnabled.IsNull() && !plan.PoolEnabled.IsUnknown() {
		additionalProps["poolEnabled"] = plan.PoolEnabled.ValueBool()
	}
	if !plan.DNSDomain.IsNull() && !plan.DNSDomain.IsUnknown() {
		additionalProps["dnsDomain"] = plan.DNSDomain.ValueString()
	}
	if !plan.DhcpServer.IsNull() && !plan.DhcpServer.IsUnknown() {
		additionalProps["dhcpServer"] = plan.DhcpServer.ValueBool()
	}
	if !plan.Gateway.IsNull() && !plan.Gateway.IsUnknown() {
		additionalProps["gateway"] = plan.Gateway.ValueString()
	}
	if !plan.Netmask.IsNull() && !plan.Netmask.IsUnknown() {
		additionalProps["netmask"] = plan.Netmask.ValueString()
	}
	if !plan.SubnetAddress.IsNull() && !plan.SubnetAddress.IsUnknown() {
		additionalProps["subnetAddress"] = plan.SubnetAddress.ValueString()
	}
	if !plan.IpRanges.IsNull() && !plan.IpRanges.IsUnknown() {
		attrs := plan.IpRanges.Attributes()
		startAddr := attrs["starting_address"].(types.String).ValueString()
		endAddr := attrs["ending_address"].(types.String).ValueString()
		additionalProps["ipRanges"] = []map[string]interface{}{
			{"startAddress": startAddr, "endAddress": endAddr},
		}
	}
	if len(additionalProps) > 0 {
		poolReq.AdditionalProperties = additionalProps
	}

	_, httpResp, err := client.NetworksAPI.UpdateNetworkPool(ctx, id).
		UpdateNetworkPoolRequest(sdk.UpdateNetworkPoolRequest{
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

	pool := readResult.NetworkPool
	if pool == nil {
		resp.Diagnostics.AddError("API returned nil", "NetworkPool is nil in the response")

		return
	}
	mapReadResponseToModel(&plan, pool)

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

func (r *networkPoolResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
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
		model.DNSDomain = types.StringValue(*pool.DnsDomain.Get())
	} else {
		model.DNSDomain = types.StringNull()
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
	if t := pool.Type; t != nil && t.Code != nil {
		model.TypeCode = types.StringValue(*t.Code)
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
	if pool.DnsDomain.IsSet() && pool.DnsDomain.Get() != nil {
		model.DNSDomain = types.StringValue(*pool.DnsDomain.Get())
	} else {
		model.DNSDomain = types.StringNull()
	}
	if pool.DhcpServer != nil {
		model.DhcpServer = types.BoolValue(*pool.DhcpServer)
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
	if len(pool.IpRanges) > 0 {
		r := pool.IpRanges[0]
		startAddr := types.StringNull()
		endAddr := types.StringNull()
		if r.StartAddress.IsSet() && r.StartAddress.Get() != nil {
			startAddr = types.StringValue(*r.StartAddress.Get())
		}
		if r.EndAddress.IsSet() && r.EndAddress.Get() != nil {
			endAddr = types.StringValue(*r.EndAddress.Get())
		}
		model.IpRanges, _ = types.ObjectValue(
			map[string]attr.Type{
				"starting_address": types.StringType,
				"ending_address":   types.StringType,
			},
			map[string]attr.Value{
				"starting_address": startAddr,
				"ending_address":   endAddr,
			},
		)
	} else {
		model.IpRanges = types.ObjectNull(map[string]attr.Type{
			"starting_address": types.StringType,
			"ending_address":   types.StringType,
		})
	}
}

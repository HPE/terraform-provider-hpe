// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package networkdhcpserver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

func getNetworkDhcpServerAsState(
	ctx context.Context,
	id int64,
	serverId float32,
	client *sdk.APIClient,
	plan NetworkDhcpServerModel,
) (NetworkDhcpServerModel, diag.Diagnostics) {
	var state NetworkDhcpServerModel
	var diags diag.Diagnostics

	dhcpResp, hresp, err := client.NetworksAPI.
		GetNetworkDhcpServer(ctx, id, serverId).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		diags.AddError(
			"populate network dhcp server resource",
			fmt.Sprintf("network dhcp server %d GET failed: ", id)+
				errfmt.ErrMsg(err, hresp),
		)

		return state, diags
	}

	dhcpServer := dhcpResp.GetNetworkDhcpServer()

	if dhcpServer.Id != nil {
		state.Id = types.Int64Value(int64(*dhcpServer.Id))
	}

	state.Name = convert.StrToType(dhcpServer.Name)
	state.ServerIpAddress = convert.StrToType(dhcpServer.ServerIpAddress)

	// leaseTime may not be returned by the GET response;
	// preserve from plan / prior state when absent.
	if dhcpServer.LeaseTime != nil {
		state.LeaseTime = types.Int64Value(int64(*dhcpServer.LeaseTime))
	} else {
		state.LeaseTime = plan.LeaseTime
	}

	// network_server_id is a path parameter not returned by the GET
	// response; preserve it from the plan / prior state.
	state.NetworkServerId = plan.NetworkServerId

	state.Config = types.DynamicNull()

	// success is not part of the GET response; set to null.
	state.Success = types.BoolNull()

	return state, diags
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var plan NetworkDhcpServerModel

	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"read network dhcp server resource",
			"failed to create client: "+err.Error(),
		)

		return
	}

	id := plan.Id.ValueInt64()
	serverId := float32(plan.NetworkServerId.ValueInt64())

	state, diags := getNetworkDhcpServerAsState(
		ctx, id, serverId, client, plan,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.Config.IsNull() && !plan.Config.IsUnknown() {
		state.Config = plan.Config
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

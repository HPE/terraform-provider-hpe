// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package networkdhcpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

func getNetworkDhcpServerAsState(
	ctx context.Context,
	id int64,
	serverID int64,
	client *sdk.APIClient,
	plan NetworkDhcpServerModel,
) (NetworkDhcpServerModel, diag.Diagnostics) {
	var state NetworkDhcpServerModel
	var diags diag.Diagnostics

	dhcpResp, hresp, err := client.NetworksAPI.
		GetNetworkDhcpServer(ctx, id, float32(serverID)).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		diags.AddError(
			"populate network dhcp server resource",
			fmt.Sprintf("network dhcp server %d GET failed: ", id)+
				errfmt.ErrMsg(err, hresp),
		)

		return state, diags
	}

	dhcpServer := dhcpResp.GetNetworkDhcpServer()

	if dhcpServer.Id == nil {
		diags.AddError(
			"populate network dhcp server resource",
			fmt.Sprintf("network dhcp server %d GET response missing id", id),
		)

		return state, diags
	}

	state.Id = types.Int64Value(int64(*dhcpServer.Id))

	state.Name = convert.StrToType(dhcpServer.Name)
	state.ServerIpAddress = convert.StrToType(dhcpServer.ServerIpAddress)

	// leaseTime may not be returned by the GET response;
	// preserve from plan / prior state when absent.
	if dhcpServer.LeaseTime != nil {
		state.LeaseTime = types.Int64Value(int64(*dhcpServer.LeaseTime))
	} else {
		state.LeaseTime = plan.LeaseTime
	}

	if networkServerID, ok := getNetworkServerIDFromAPI(dhcpServer); ok {
		state.NetworkServerId = types.Int64Value(networkServerID)
	} else {
		// Some API versions do not include networkServer in this response.
		state.NetworkServerId = plan.NetworkServerId
	}

	state.Config = types.DynamicNull()
	state.ConfigNsx = NewConfigNsxValueNull()

	cfg := dhcpServer.GetConfig()

	switch {
	case !plan.ConfigNsx.IsNull() && !plan.ConfigNsx.IsUnknown():
		edgeCluster := types.StringNull()
		activeEdgeNode := types.StringNull()
		standbyEdgeNode := types.StringNull()

		if nsxCfg := cfg.NSXDHCPServerConfiguration; nsxCfg != nil {
			edgeCluster = convert.StrToType(nsxCfg.EdgeCluster.Get())
			activeEdgeNode = convert.StrToType(nsxCfg.PreferredEdgeNode1.Get())
			standbyEdgeNode = convert.StrToType(nsxCfg.PreferredEdgeNode2.Get())
		}

		nsxValue, nsxDiags := NewConfigNsxValue(
			ConfigNsxValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"edge_cluster":      edgeCluster,
				"active_edge_node":  activeEdgeNode,
				"standby_edge_node": standbyEdgeNode,
			},
		)
		if nsxDiags.HasError() {
			diags.Append(nsxDiags...)

			return state, diags
		}

		state.ConfigNsx = nsxValue
	case !plan.Config.IsNull() && !plan.Config.IsUnknown():
		state.Config = plan.Config
	}

	// success is not part of the GET response; set to null.
	state.Success = types.BoolNull()

	return state, diags
}

func getNetworkServerIDFromAPI(dhcpServer any) (int64, bool) {
	encoded, err := json.Marshal(dhcpServer)
	if err != nil {
		return 0, false
	}

	var payload map[string]any
	if err := json.Unmarshal(encoded, &payload); err != nil {
		return 0, false
	}

	networkServerRaw, ok := payload["networkServer"]
	if !ok {
		return 0, false
	}

	networkServerMap, ok := networkServerRaw.(map[string]any)
	if !ok {
		return 0, false
	}

	idRaw, ok := networkServerMap["id"]
	if !ok {
		return 0, false
	}

	idFloat, ok := idRaw.(float64)
	if !ok {
		return 0, false
	}

	return int64(idFloat), true
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
	serverID := plan.NetworkServerId.ValueInt64()

	state, diags := getNetworkDhcpServerAsState(
		ctx, id, serverID, client, plan,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

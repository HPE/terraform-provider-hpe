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

// dhcpServerReadPayload is a local struct for decoding the untyped GET
// response returned by the SDK (interface{}).
type dhcpServerReadPayload struct {
	Id              *int64                      `json:"id"`
	Name            *string                     `json:"name"`
	ServerIpAddress *string                     `json:"serverIpAddress"`
	LeaseTime       *int64                      `json:"leaseTime"`
	Config          json.RawMessage             `json:"config"`
	NetworkServer   *dhcpServerNetworkServerRef `json:"networkServer"`
}

type dhcpServerNetworkServerRef struct {
	Id *float64 `json:"id"`
}

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
		GetNetworkDhcpServer(ctx, id, serverID).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		diags.AddError(
			"populate network dhcp server resource",
			fmt.Sprintf("network dhcp server %d GET failed: ", id)+
				errfmt.ErrMsg(err, hresp),
		)

		return state, diags
	}

	raw := dhcpResp.GetNetworkDhcpServer()

	encoded, err := json.Marshal(raw)
	if err != nil {
		diags.AddError(
			"populate network dhcp server resource",
			fmt.Sprintf(
				"network dhcp server %d: failed to marshal response: %s",
				id, err,
			),
		)

		return state, diags
	}

	var payload dhcpServerReadPayload
	if err := json.Unmarshal(encoded, &payload); err != nil {
		diags.AddError(
			"populate network dhcp server resource",
			fmt.Sprintf(
				"network dhcp server %d: failed to unmarshal response: %s",
				id, err,
			),
		)

		return state, diags
	}

	if payload.Id == nil {
		diags.AddError(
			"populate network dhcp server resource",
			fmt.Sprintf("network dhcp server %d GET response missing id", id),
		)

		return state, diags
	}

	state.Id = types.Int64Value(*payload.Id)
	state.Name = convert.StrToType(payload.Name)
	state.ServerIpAddress = convert.StrToType(payload.ServerIpAddress)

	// leaseTime may not be returned by the GET response;
	// preserve from plan / prior state when absent.
	if payload.LeaseTime != nil {
		state.LeaseTime = types.Int64Value(*payload.LeaseTime)
	} else {
		state.LeaseTime = plan.LeaseTime
	}

	if payload.NetworkServer != nil && payload.NetworkServer.Id != nil {
		state.NetworkServerId = types.Int64Value(int64(*payload.NetworkServer.Id))
	} else {
		// Some API versions do not include networkServer in this response.
		state.NetworkServerId = plan.NetworkServerId
	}

	state.Config = types.DynamicNull()
	state.ConfigNsx = NewConfigNsxValueNull()

	switch {
	case !plan.ConfigNsx.IsNull() && !plan.ConfigNsx.IsUnknown():
		edgeCluster := types.StringNull()
		activeEdgeNode := types.StringNull()
		standbyEdgeNode := types.StringNull()

		if len(payload.Config) > 0 {
			var nsxCfg sdk.NetworkDhcpServerConfigNSX
			if err := json.Unmarshal(payload.Config, &nsxCfg); err == nil {
				edgeCluster = convert.StrToType(nsxCfg.EdgeCluster.Get())
				activeEdgeNode = convert.StrToType(
					nsxCfg.PreferredEdgeNode1.Get(),
				)
				standbyEdgeNode = convert.StrToType(
					nsxCfg.PreferredEdgeNode2.Get(),
				)
			}
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

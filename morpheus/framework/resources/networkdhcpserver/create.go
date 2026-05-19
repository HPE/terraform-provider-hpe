// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package networkdhcpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan NetworkDhcpServerModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"create network dhcp server resource",
			"failed to create client: "+err.Error(),
		)

		return
	}

	name := plan.Name.ValueString()

	dhcpServer := sdk.
		NewCreateNetworkDhcpServerRequestNetworkDhcpServerWithDefaults()
	dhcpServer.SetName(name)
	dhcpServer.SetServerIpAddress(plan.ServerIpAddress.ValueString())

	if !plan.LeaseTime.IsNull() && !plan.LeaseTime.IsUnknown() {
		dhcpServer.SetLeaseTime(plan.LeaseTime.ValueInt64())
	}

	switch {
	case !plan.ConfigNsxt.IsNull() && !plan.ConfigNsxt.IsUnknown():
		nsxConfig := sdk.NewNSXDHCPServerConfiguration1()

		if !plan.ConfigNsxt.EdgeCluster.IsNull() &&
			!plan.ConfigNsxt.EdgeCluster.IsUnknown() {
			nsxConfig.SetEdgeCluster(plan.ConfigNsxt.EdgeCluster.ValueString())
		}

		if !plan.ConfigNsxt.ActiveEdgeNode.IsNull() &&
			!plan.ConfigNsxt.ActiveEdgeNode.IsUnknown() {
			nsxConfig.SetPreferredEdgeNode1(
				plan.ConfigNsxt.ActiveEdgeNode.ValueString(),
			)
		}

		if !plan.ConfigNsxt.StandbyEdgeNode.IsNull() &&
			!plan.ConfigNsxt.StandbyEdgeNode.IsUnknown() {
			nsxConfig.SetPreferredEdgeNode2(
				plan.ConfigNsxt.StandbyEdgeNode.ValueString(),
			)
		}

		dhcpServer.SetConfig(sdk.CreateNetworkDhcpServerRequestNetworkDhcpServerConfig{
			NSXDHCPServerConfiguration1: nsxConfig,
		})

	case !plan.Config.IsNull() && !plan.Config.IsUnknown():
		configValue := plan.Config.UnderlyingValue()
		configMap, err := convert.ValueToAny(ctx, configValue)
		if err != nil {
			resp.Diagnostics.AddError(
				"create network dhcp server resource",
				"network dhcp server "+name+
					": failed to convert config: "+err.Error(),
			)

			return
		}

		configDataMap, ok := configMap.(map[string]any)
		if !ok {
			resp.Diagnostics.AddError(
				"create network dhcp server resource",
				"network dhcp server "+name+
					": config must be a valid object/map",
			)

			return
		}

		configJSON, err := json.Marshal(configDataMap)
		if err != nil {
			resp.Diagnostics.AddError(
				"create network dhcp server resource",
				"network dhcp server "+name+
					": failed to marshal config: "+err.Error(),
			)

			return
		}

		var config sdk.CreateNetworkDhcpServerRequestNetworkDhcpServerConfig
		if err := json.Unmarshal(configJSON, &config); err != nil {
			resp.Diagnostics.AddError(
				"create network dhcp server resource",
				"network dhcp server "+name+
					": failed to unmarshal config: "+err.Error(),
			)

			return
		}

		dhcpServer.SetConfig(config)
	}

	createReq := sdk.NewCreateNetworkDhcpServerRequestWithDefaults()
	createReq.SetNetworkDhcpServer(*dhcpServer)
	serverID := plan.NetworkIntegrationId.ValueInt64()

	createResp, hresp, err := client.NetworksAPI.
		CreateNetworkDhcpServer(ctx, serverID).
		CreateNetworkDhcpServerRequest(*createReq).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"create network dhcp server resource",
			fmt.Sprintf("network dhcp server %s POST failed: %s",
				name, errfmt.ErrMsg(err, hresp)),
		)

		return
	}

	idPtr, ok := createResp.GetIdOk()
	if !ok || idPtr == nil {
		resp.Diagnostics.AddError(
			"create network dhcp server resource",
			"network dhcp server "+name+": id is nil",
		)

		return
	}

	dhcpServerID := *idPtr
	plan.Id = types.Int64Value(dhcpServerID)

	taintResourceState := func(id int64) {
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "network_dhcp_server",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})
	}

	state, pdiags := getNetworkDhcpServerAsState(
		ctx, dhcpServerID, serverID, client, plan,
	)
	if pdiags.HasError() {
		resp.Diagnostics.Append(pdiags...)
		resp.Diagnostics.AddError(
			"failed to read network dhcp server state",
			fmt.Sprintf(
				"Network DHCP Server %d was created but could not be read",
				dhcpServerID,
			),
		)
		taintResourceState(dhcpServerID)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError(
			"failed to set network dhcp server state",
			fmt.Sprintf(
				"Network DHCP Server %d was created but state could not be saved",
				dhcpServerID,
			),
		)
		taintResourceState(dhcpServerID)

		return
	}
}

// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package networkdhcpserver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state NetworkDhcpServerModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"update network dhcp server resource",
			"failed to create client: "+err.Error(),
		)

		return
	}

	id := state.Id.ValueInt64()
	serverId := float32(plan.NetworkServerId.ValueInt64())

	dhcpServerMap := map[string]interface{}{}

	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		dhcpServerMap["name"] = plan.Name.ValueString()
	}

	if !plan.ServerIpAddress.IsNull() &&
		!plan.ServerIpAddress.IsUnknown() {
		dhcpServerMap["serverIpAddress"] = plan.
			ServerIpAddress.ValueString()
	}

	if !plan.LeaseTime.IsNull() && !plan.LeaseTime.IsUnknown() {
		dhcpServerMap["leaseTime"] = plan.LeaseTime.ValueInt64()
	}

	switch {
	case !plan.ConfigNsx.IsNull() && !plan.ConfigNsx.IsUnknown():
		configMap := map[string]interface{}{}

		if !plan.ConfigNsx.EdgeCluster.IsNull() &&
			!plan.ConfigNsx.EdgeCluster.IsUnknown() {
			configMap["edgeCluster"] = plan.ConfigNsx.EdgeCluster.ValueString()
		}

		if !plan.ConfigNsx.ActiveEdgeNode.IsNull() &&
			!plan.ConfigNsx.ActiveEdgeNode.IsUnknown() {
			configMap["preferredEdgeNode1"] = plan.ConfigNsx.
				ActiveEdgeNode.ValueString()
		}

		if !plan.ConfigNsx.StandbyEdgeNode.IsNull() &&
			!plan.ConfigNsx.StandbyEdgeNode.IsUnknown() {
			configMap["preferredEdgeNode2"] = plan.ConfigNsx.
				StandbyEdgeNode.ValueString()
		}

		dhcpServerMap["config"] = configMap

	case !plan.Config.IsNull() && !plan.Config.IsUnknown():
		configValue := plan.Config.UnderlyingValue()
		configMap, err := convert.ValueToAny(ctx, configValue)
		if err != nil {
			resp.Diagnostics.AddError(
				"update network dhcp server resource",
				fmt.Sprintf(
					"network dhcp server %d: failed to convert config: %s",
					id, err.Error(),
				),
			)

			return
		}

		configDataMap, ok := configMap.(map[string]any)
		if ok {
			dhcpServerMap["config"] = configDataMap
		} else {
			resp.Diagnostics.AddError(
				"update network dhcp server resource",
				fmt.Sprintf(
					"network dhcp server %d: config must be a valid object/map",
					id,
				),
			)

			return
		}
	}

	updateReq := sdk.NewUpdateNetworkDhcpServerRequestWithDefaults()
	updateReq.SetNetworkDhcpServer(dhcpServerMap)

	_, hresp, err := client.NetworksAPI.
		UpdateNetworkDhcpServer(ctx, id, serverId).
		UpdateNetworkDhcpServerRequest(*updateReq).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"update network dhcp server resource",
			fmt.Sprintf(
				"network dhcp server %d UPDATE failed: %s",
				id, errfmt.ErrMsg(err, hresp),
			),
		)

		return
	}

	newState, diags := getNetworkDhcpServerAsState(
		ctx, id, serverId, client, plan,
	)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		resp.Diagnostics.AddError(
			"update network dhcp server resource",
			fmt.Sprintf(
				"network dhcp server %d: failed to read from api",
				id,
			),
		)

		return
	}

	if !plan.Config.IsNull() && !plan.Config.IsUnknown() {
		newState.Config = plan.Config
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

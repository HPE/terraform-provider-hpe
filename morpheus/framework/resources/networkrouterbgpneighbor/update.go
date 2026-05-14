// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterbgpneighbor

import (
	"context"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const updateOperation = "update network router bgp neighbor resource"

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state NetworkRouterBgpNeighborModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			updateOperation,
			"failed to create client: "+err.Error(),
		)

		return
	}

	id := state.Id.ValueInt64()
	routerID := state.RouterId.ValueInt64()

	neighbor := sdk.NewUpdateNetworkRouterBgpNeighborRequestNetworkRouterBgpNeighbor()

	if !plan.IpAddress.IsNull() && !plan.IpAddress.IsUnknown() {
		neighbor.SetIpAddress(plan.IpAddress.ValueString())
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		neighbor.SetDescription(plan.Description.ValueString())
	}

	if !plan.ForwardingAddress.IsNull() && !plan.ForwardingAddress.IsUnknown() {
		neighbor.SetForwardingAddress(plan.ForwardingAddress.ValueString())
	}

	if !plan.ProtocolAddress.IsNull() && !plan.ProtocolAddress.IsUnknown() {
		neighbor.SetProtocolAddress(plan.ProtocolAddress.ValueString())
	}

	if !plan.RemoteAs.IsNull() && !plan.RemoteAs.IsUnknown() {
		neighbor.SetRemoteAs(plan.RemoteAs.ValueString())
	}

	if !plan.Weight.IsNull() && !plan.Weight.IsUnknown() {
		neighbor.SetWeight(plan.Weight.ValueInt64())
	}

	if !plan.KeepAlive.IsNull() && !plan.KeepAlive.IsUnknown() {
		neighbor.SetKeepAlive(plan.KeepAlive.ValueInt64())
	}

	if !plan.HoldDown.IsNull() && !plan.HoldDown.IsUnknown() {
		neighbor.SetHoldDown(plan.HoldDown.ValueInt64())
	}

	if !plan.PasswordWo.IsNull() && !plan.PasswordWo.IsUnknown() {
		neighbor.SetPassword(plan.PasswordWo.ValueString())
	}

	if !plan.RouteFilteringType.IsNull() && !plan.RouteFilteringType.IsUnknown() {
		neighbor.SetRouteFilteringType(plan.RouteFilteringType.ValueString())
	}

	if !plan.RouteFilteringIn.IsNull() && !plan.RouteFilteringIn.IsUnknown() {
		neighbor.SetRouteFilteringIn(plan.RouteFilteringIn.ValueString())
	}

	if !plan.RouteFilteringOut.IsNull() && !plan.RouteFilteringOut.IsUnknown() {
		neighbor.SetRouteFilteringOut(plan.RouteFilteringOut.ValueString())
	}

	if !plan.BfdEnabled.IsNull() && !plan.BfdEnabled.IsUnknown() {
		bfdVal := sdk.UpdateNetworkRouterBgpNeighborRequestNetworkRouterBgpNeighborBfdEnabled{}
		boolVal := plan.BfdEnabled.ValueBool()
		bfdVal.Bool = &boolVal
		neighbor.SetBfdEnabled(bfdVal)
	}

	if !plan.BfdInterval.IsNull() && !plan.BfdInterval.IsUnknown() {
		neighbor.SetBfdInterval(plan.BfdInterval.ValueInt64())
	}

	if !plan.BfdMultiple.IsNull() && !plan.BfdMultiple.IsUnknown() {
		neighbor.SetBfdMultiple(plan.BfdMultiple.ValueInt64())
	}

	if !plan.AllowAsIn.IsNull() && !plan.AllowAsIn.IsUnknown() {
		allowVal := sdk.UpdateNetworkRouterBgpNeighborRequestNetworkRouterBgpNeighborAllowAsIn{}
		boolVal := plan.AllowAsIn.ValueBool()
		allowVal.Bool = &boolVal
		neighbor.SetAllowAsIn(allowVal)
	}

	if !plan.HopLimit.IsNull() && !plan.HopLimit.IsUnknown() {
		neighbor.SetHopLimit(plan.HopLimit.ValueInt64())
	}

	if !plan.RestartMode.IsNull() && !plan.RestartMode.IsUnknown() {
		neighbor.SetRestartMode(plan.RestartMode.ValueString())
	}

	updateConfig := buildUpdateConfig(ctx, plan, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	if updateConfig != nil {
		neighbor.SetConfig(*updateConfig)
	}

	updateReq := sdk.NewUpdateNetworkRouterBgpNeighborRequest()
	updateReq.SetNetworkRouterBgpNeighbor(*neighbor)

	_, hresp, err := client.NetworksAPI.
		UpdateNetworkRouterBgpNeighbor(ctx, id, routerID).
		UpdateNetworkRouterBgpNeighborRequest(*updateReq).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			updateOperation,
			fmt.Sprintf("bgp neighbor %d UPDATE failed: %s",
				id, errfmt.ErrMsg(err, hresp)),
		)

		return
	}

	newState, diags := getNetworkRouterBgpNeighborAsState(ctx, id, routerID, client, plan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		resp.Diagnostics.AddError(
			updateOperation,
			fmt.Sprintf("bgp neighbor %d: failed to read from api", id),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func buildUpdateConfig(
	ctx context.Context,
	plan NetworkRouterBgpNeighborModel,
	resp *resource.UpdateResponse,
) *sdk.UpdateNetworkRouterBgpNeighborRequestNetworkRouterBgpNeighborConfig {
	config := &sdk.UpdateNetworkRouterBgpNeighborRequestNetworkRouterBgpNeighborConfig{}

	switch {
	case !plan.ConfigNsxt.IsNull() && !plan.ConfigNsxt.IsUnknown():
		nsxtConfig := sdk.NewNSXTBGPNeighborConfig3()

		if !plan.ConfigNsxt.SourceAddresses.IsNull() &&
			!plan.ConfigNsxt.SourceAddresses.IsUnknown() {
			var addresses []types.String
			diags := plan.ConfigNsxt.SourceAddresses.ElementsAs(ctx, &addresses, false)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return nil
			}

			var strAddrs []string
			for _, addr := range addresses {
				if !addr.IsNull() {
					strAddrs = append(strAddrs, addr.ValueString())
				}
			}

			nsxtConfig.SetSourceAddresses(strAddrs)
		}

		config.NSXTBGPNeighborConfig3 = nsxtConfig

	case !plan.ConfigNsxv.IsNull() && !plan.ConfigNsxv.IsUnknown():
		nsxvConfig := sdk.NewNSXVBGPNeighborConfig3()

		if !plan.ConfigNsxv.RouterId.IsNull() &&
			!plan.ConfigNsxv.RouterId.IsUnknown() {
			nsxvConfig.SetRouterId(plan.ConfigNsxv.RouterId.ValueString())
		}

		if !plan.ConfigNsxv.Interface.IsNull() &&
			!plan.ConfigNsxv.Interface.IsUnknown() {
			nsxvConfig.SetInterface(plan.ConfigNsxv.Interface.ValueString())
		}

		config.NSXVBGPNeighborConfig3 = nsxvConfig

	case !plan.Config.IsNull() && !plan.Config.IsUnknown():
		configValue := plan.Config.UnderlyingValue()
		configMap, err := convert.ValueToAny(ctx, configValue)
		if err != nil {
			resp.Diagnostics.AddError(
				updateOperation,
				"failed to convert config: "+err.Error(),
			)

			return nil
		}

		configDataMap, ok := configMap.(map[string]any)
		if !ok {
			resp.Diagnostics.AddError(
				updateOperation,
				"config must be a valid object/map",
			)

			return nil
		}

		config.MapmapOfStringAny = &configDataMap

	default:
		return nil
	}

	return config
}

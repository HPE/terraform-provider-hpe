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

	neighbor := &sdk.UpdateNetworkRouterBgpNeighborRequestNetworkRouterBgpNeighbor{}

	if !plan.IpAddress.IsNull() && !plan.IpAddress.IsUnknown() {
		neighbor.IpAddress = plan.IpAddress.ValueStringPointer()
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		neighbor.Description = plan.Description.ValueStringPointer()
	}

	if !plan.ForwardingAddress.IsNull() && !plan.ForwardingAddress.IsUnknown() {
		neighbor.ForwardingAddress = plan.ForwardingAddress.ValueStringPointer()
	}

	if !plan.ProtocolAddress.IsNull() && !plan.ProtocolAddress.IsUnknown() {
		neighbor.ProtocolAddress = plan.ProtocolAddress.ValueStringPointer()
	}

	if !plan.RemoteAs.IsNull() && !plan.RemoteAs.IsUnknown() {
		neighbor.RemoteAs = plan.RemoteAs.ValueStringPointer()
	}

	if !plan.Weight.IsNull() && !plan.Weight.IsUnknown() {
		neighbor.Weight = plan.Weight.ValueInt64Pointer()
	}

	if !plan.KeepAlive.IsNull() && !plan.KeepAlive.IsUnknown() {
		neighbor.KeepAlive = plan.KeepAlive.ValueInt64Pointer()
	}

	if !plan.HoldDown.IsNull() && !plan.HoldDown.IsUnknown() {
		neighbor.HoldDown = plan.HoldDown.ValueInt64Pointer()
	}

	if !plan.PasswordWo.IsNull() && !plan.PasswordWo.IsUnknown() {
		neighbor.Password = plan.PasswordWo.ValueStringPointer()
	}

	if !plan.RouteFilteringType.IsNull() && !plan.RouteFilteringType.IsUnknown() {
		neighbor.RouteFilteringType = plan.RouteFilteringType.ValueStringPointer()
	}

	if !plan.RouteFilteringIn.IsNull() && !plan.RouteFilteringIn.IsUnknown() {
		neighbor.RouteFilteringIn = plan.RouteFilteringIn.ValueStringPointer()
	}

	if !plan.RouteFilteringOut.IsNull() && !plan.RouteFilteringOut.IsUnknown() {
		neighbor.RouteFilteringOut = plan.RouteFilteringOut.ValueStringPointer()
	}

	if !plan.BfdEnabled.IsNull() && !plan.BfdEnabled.IsUnknown() {
		bfdVal := sdk.UpdateNetworkRouterBgpNeighborRequestNetworkRouterBgpNeighborBfdEnabled{}
		boolVal := plan.BfdEnabled.ValueBool()
		bfdVal.Bool = &boolVal
		neighbor.BfdEnabled = &bfdVal
	}

	if !plan.BfdInterval.IsNull() && !plan.BfdInterval.IsUnknown() {
		neighbor.BfdInterval = plan.BfdInterval.ValueInt64Pointer()
	}

	if !plan.BfdMultiple.IsNull() && !plan.BfdMultiple.IsUnknown() {
		neighbor.BfdMultiple = plan.BfdMultiple.ValueInt64Pointer()
	}

	if !plan.AllowAsIn.IsNull() && !plan.AllowAsIn.IsUnknown() {
		allowVal := sdk.UpdateNetworkRouterBgpNeighborRequestNetworkRouterBgpNeighborAllowAsIn{}
		boolVal := plan.AllowAsIn.ValueBool()
		allowVal.Bool = &boolVal
		neighbor.AllowAsIn = &allowVal
	}

	if !plan.HopLimit.IsNull() && !plan.HopLimit.IsUnknown() {
		neighbor.HopLimit = plan.HopLimit.ValueInt64Pointer()
	}

	if !plan.RestartMode.IsNull() && !plan.RestartMode.IsUnknown() {
		neighbor.RestartMode = plan.RestartMode.ValueStringPointer()
	}

	updateConfig := buildUpdateConfig(ctx, plan, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	if updateConfig != nil {
		neighbor.Config = updateConfig
	}

	updateReq := &sdk.UpdateNetworkRouterBgpNeighborRequest{}
	updateReq.NetworkRouterBgpNeighbor = neighbor

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
		nsxtConfig := &sdk.NSXTBGPNeighborConfig3{}

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

			nsxtConfig.SourceAddresses = strAddrs
		}

		config.NSXTBGPNeighborConfig3 = nsxtConfig

	case !plan.ConfigNsxv.IsNull() && !plan.ConfigNsxv.IsUnknown():
		nsxvConfig := &sdk.NSXVBGPNeighborConfig3{}

		if !plan.ConfigNsxv.RouterId.IsNull() &&
			!plan.ConfigNsxv.RouterId.IsUnknown() {
			nsxvConfig.RouterId = plan.ConfigNsxv.RouterId.ValueStringPointer()
		}

		if !plan.ConfigNsxv.Interface.IsNull() &&
			!plan.ConfigNsxv.Interface.IsUnknown() {
			nsxvConfig.Interface = plan.ConfigNsxv.Interface.ValueStringPointer()
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

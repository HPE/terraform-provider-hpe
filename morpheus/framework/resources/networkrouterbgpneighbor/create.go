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
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const createOperation = "create network router bgp neighbor resource"

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan NetworkRouterBgpNeighborModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			createOperation,
			"failed to create client: "+err.Error(),
		)

		return
	}

	routerID := plan.RouterId.ValueInt64()

	neighbor := &sdk.CreateNetworkRouterBgpNeighborRequestNetworkRouterBgpNeighbor{
		IpAddress: plan.IpAddress.ValueString(),
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
		bfdVal := sdk.CreateNetworkRouterBgpNeighborRequestNetworkRouterBgpNeighborBfdEnabled{}
		bfdVal.Bool = plan.BfdEnabled.ValueBoolPointer()
		neighbor.BfdEnabled = &bfdVal
	}

	if !plan.BfdInterval.IsNull() && !plan.BfdInterval.IsUnknown() {
		neighbor.BfdInterval = plan.BfdInterval.ValueInt64Pointer()
	}

	if !plan.BfdMultiple.IsNull() && !plan.BfdMultiple.IsUnknown() {
		neighbor.BfdMultiple = plan.BfdMultiple.ValueInt64Pointer()
	}

	if !plan.AllowAsIn.IsNull() && !plan.AllowAsIn.IsUnknown() {
		allowVal := sdk.CreateNetworkRouterBgpNeighborRequestNetworkRouterBgpNeighborAllowAsIn{}
		allowVal.Bool = plan.AllowAsIn.ValueBoolPointer()
		neighbor.AllowAsIn = &allowVal
	}

	if !plan.HopLimit.IsNull() && !plan.HopLimit.IsUnknown() {
		neighbor.HopLimit = plan.HopLimit.ValueInt64Pointer()
	}

	if !plan.RestartMode.IsNull() && !plan.RestartMode.IsUnknown() {
		neighbor.RestartMode = plan.RestartMode.ValueStringPointer()
	}

	config := buildCreateConfig(ctx, plan, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	if config != nil {
		neighbor.Config = config
	}

	createReq := &sdk.CreateNetworkRouterBgpNeighborRequest{}
	createReq.NetworkRouterBgpNeighbor = neighbor

	result, hresp, err := client.NetworksAPI.
		CreateNetworkRouterBgpNeighbor(ctx, routerID).
		CreateNetworkRouterBgpNeighborRequest(*createReq).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			createOperation,
			fmt.Sprintf("router %d bgp neighbor POST failed: %s",
				routerID, errfmt.ErrMsg(err, hresp)),
		)

		return
	}

	idPtr := result.Id.Get()
	if idPtr == nil {
		resp.Diagnostics.AddError(
			createOperation,
			fmt.Sprintf("router %d bgp neighbor POST succeeded but response id was missing", routerID),
		)

		return
	}

	id := *idPtr
	plan.Id = types.Int64Value(id)

	taintResourceState := func(id int64) {
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "network_router_bgp_neighbor",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})
	}

	state, pdiags := getNetworkRouterBgpNeighborAsState(ctx, id, routerID, client, plan)
	if pdiags.HasError() {
		resp.Diagnostics.Append(pdiags...)
		resp.Diagnostics.AddError(
			createOperation,
			fmt.Sprintf("BGP neighbor %d was created but could not be read", id),
		)
		taintResourceState(id)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError(
			createOperation,
			fmt.Sprintf("BGP neighbor %d was created but state could not be saved", id),
		)
		taintResourceState(id)

		return
	}
}

func buildCreateConfig(
	ctx context.Context,
	plan NetworkRouterBgpNeighborModel,
	resp *resource.CreateResponse,
) *sdk.CreateNetworkRouterBgpNeighborRequestNetworkRouterBgpNeighborConfig {
	config := &sdk.CreateNetworkRouterBgpNeighborRequestNetworkRouterBgpNeighborConfig{}

	switch {
	case !plan.ConfigNsxt.IsNull() && !plan.ConfigNsxt.IsUnknown():
		nsxtConfig := &sdk.NSXTBGPNeighborConfig1{}

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

		config.NSXTBGPNeighborConfig1 = nsxtConfig

	case !plan.ConfigNsxv.IsNull() && !plan.ConfigNsxv.IsUnknown():
		nsxvConfig := &sdk.NSXVBGPNeighborConfig1{}

		if !plan.ConfigNsxv.RouterId.IsNull() &&
			!plan.ConfigNsxv.RouterId.IsUnknown() {
			nsxvConfig.RouterId = plan.ConfigNsxv.RouterId.ValueStringPointer()
		}

		if !plan.ConfigNsxv.Interface.IsNull() &&
			!plan.ConfigNsxv.Interface.IsUnknown() {
			nsxvConfig.Interface = plan.ConfigNsxv.Interface.ValueStringPointer()
		}

		config.NSXVBGPNeighborConfig1 = nsxvConfig

	case !plan.Config.IsNull() && !plan.Config.IsUnknown():
		configValue := plan.Config.UnderlyingValue()
		configMap, err := convert.ValueToAny(ctx, configValue)
		if err != nil {
			resp.Diagnostics.AddError(
				createOperation,
				"failed to convert config: "+err.Error(),
			)

			return nil
		}

		configDataMap, ok := configMap.(map[string]any)
		if !ok {
			resp.Diagnostics.AddError(
				createOperation,
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

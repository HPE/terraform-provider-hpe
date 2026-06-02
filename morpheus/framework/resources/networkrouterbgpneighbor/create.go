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
		neighbor.Description = sdk.PtrString(plan.Description.ValueString())
	}

	if !plan.ForwardingAddress.IsNull() && !plan.ForwardingAddress.IsUnknown() {
		neighbor.ForwardingAddress = sdk.PtrString(plan.ForwardingAddress.ValueString())
	}

	if !plan.ProtocolAddress.IsNull() && !plan.ProtocolAddress.IsUnknown() {
		neighbor.ProtocolAddress = sdk.PtrString(plan.ProtocolAddress.ValueString())
	}

	if !plan.RemoteAs.IsNull() && !plan.RemoteAs.IsUnknown() {
		neighbor.RemoteAs = sdk.PtrString(plan.RemoteAs.ValueString())
	}

	if !plan.Weight.IsNull() && !plan.Weight.IsUnknown() {
		neighbor.Weight = sdk.PtrInt64(plan.Weight.ValueInt64())
	}

	if !plan.KeepAlive.IsNull() && !plan.KeepAlive.IsUnknown() {
		neighbor.KeepAlive = sdk.PtrInt64(plan.KeepAlive.ValueInt64())
	}

	if !plan.HoldDown.IsNull() && !plan.HoldDown.IsUnknown() {
		neighbor.HoldDown = sdk.PtrInt64(plan.HoldDown.ValueInt64())
	}

	if !plan.PasswordWo.IsNull() && !plan.PasswordWo.IsUnknown() {
		neighbor.Password = sdk.PtrString(plan.PasswordWo.ValueString())
	}

	if !plan.RouteFilteringType.IsNull() && !plan.RouteFilteringType.IsUnknown() {
		neighbor.RouteFilteringType = sdk.PtrString(plan.RouteFilteringType.ValueString())
	}

	if !plan.RouteFilteringIn.IsNull() && !plan.RouteFilteringIn.IsUnknown() {
		neighbor.RouteFilteringIn = sdk.PtrString(plan.RouteFilteringIn.ValueString())
	}

	if !plan.RouteFilteringOut.IsNull() && !plan.RouteFilteringOut.IsUnknown() {
		neighbor.RouteFilteringOut = sdk.PtrString(plan.RouteFilteringOut.ValueString())
	}

	if !plan.BfdEnabled.IsNull() && !plan.BfdEnabled.IsUnknown() {
		bfdVal := sdk.CreateNetworkRouterBgpNeighborRequestNetworkRouterBgpNeighborBfdEnabled{}
		boolVal := plan.BfdEnabled.ValueBool()
		bfdVal.Bool = &boolVal
		neighbor.BfdEnabled = &bfdVal
	}

	if !plan.BfdInterval.IsNull() && !plan.BfdInterval.IsUnknown() {
		neighbor.BfdInterval = sdk.PtrInt64(plan.BfdInterval.ValueInt64())
	}

	if !plan.BfdMultiple.IsNull() && !plan.BfdMultiple.IsUnknown() {
		neighbor.BfdMultiple = sdk.PtrInt64(plan.BfdMultiple.ValueInt64())
	}

	if !plan.AllowAsIn.IsNull() && !plan.AllowAsIn.IsUnknown() {
		allowVal := sdk.CreateNetworkRouterBgpNeighborRequestNetworkRouterBgpNeighborAllowAsIn{}
		boolVal := plan.AllowAsIn.ValueBool()
		allowVal.Bool = &boolVal
		neighbor.AllowAsIn = &allowVal
	}

	if !plan.HopLimit.IsNull() && !plan.HopLimit.IsUnknown() {
		neighbor.HopLimit = sdk.PtrInt64(plan.HopLimit.ValueInt64())
	}

	if !plan.RestartMode.IsNull() && !plan.RestartMode.IsUnknown() {
		neighbor.RestartMode = sdk.PtrString(plan.RestartMode.ValueString())
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
			nsxvConfig.RouterId = sdk.PtrString(plan.ConfigNsxv.RouterId.ValueString())
		}

		if !plan.ConfigNsxv.Interface.IsNull() &&
			!plan.ConfigNsxv.Interface.IsUnknown() {
			nsxvConfig.Interface = sdk.PtrString(plan.ConfigNsxv.Interface.ValueString())
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

// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterbgpneighbor

import (
	"context"
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

const readOperation = "read network router bgp neighbor resource"

func getNetworkRouterBgpNeighborAsState(
	ctx context.Context,
	id int64,
	routerID int64,
	client *sdk.APIClient,
	plan NetworkRouterBgpNeighborModel,
) (NetworkRouterBgpNeighborModel, diag.Diagnostics) {
	var state NetworkRouterBgpNeighborModel
	var diags diag.Diagnostics

	resp, hresp, err := client.NetworksAPI.
		GetNetworkRouterBgpNeighbor(ctx, id, routerID).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		diags.AddError(
			readOperation,
			fmt.Sprintf("bgp neighbor %d GET failed: ", id)+
				errfmt.ErrMsg(err, hresp),
		)

		return state, diags
	}

	neighbor := resp.NetworkRouterBgpNeighbor
	if neighbor == nil {
		diags.AddError(
			readOperation,
			fmt.Sprintf("bgp neighbor %d GET succeeded but response payload was missing", id),
		)

		return state, diags
	}

	state.Id = convert.Int64ToType(neighbor.Id)
	state.RouterId = plan.RouterId
	state.IpAddress = convert.StrToType(neighbor.IpAddress)

	if neighbor.Description.IsSet() {
		state.Description = convert.StrToType(neighbor.Description.Get())
	} else {
		state.Description = types.StringNull()
	}

	if neighbor.ForwardingAddress.IsSet() {
		state.ForwardingAddress = convert.StrToType(neighbor.ForwardingAddress.Get())
	} else {
		state.ForwardingAddress = types.StringNull()
	}

	if neighbor.ProtocolAddress.IsSet() {
		state.ProtocolAddress = convert.StrToType(neighbor.ProtocolAddress.Get())
	} else {
		state.ProtocolAddress = types.StringNull()
	}

	state.RemoteAs = convert.StrToType(neighbor.RemoteAs)
	state.Weight = convert.Int64ToType(neighbor.Weight)
	state.KeepAlive = convert.Int64ToType(neighbor.KeepAlive)
	state.HoldDown = convert.Int64ToType(neighbor.HoldDown)

	if neighbor.RouteFilteringType.IsSet() {
		state.RouteFilteringType = convert.StrToType(neighbor.RouteFilteringType.Get())
	} else {
		state.RouteFilteringType = types.StringNull()
	}

	if neighbor.RouteFilteringIn.IsSet() {
		state.RouteFilteringIn = convert.StrToType(neighbor.RouteFilteringIn.Get())
	} else {
		state.RouteFilteringIn = types.StringNull()
	}

	if neighbor.RouteFilteringOut.IsSet() {
		state.RouteFilteringOut = convert.StrToType(neighbor.RouteFilteringOut.Get())
	} else {
		state.RouteFilteringOut = types.StringNull()
	}

	if neighbor.BfdEnabled.IsSet() {
		state.BfdEnabled = convert.BoolToType(neighbor.BfdEnabled.Get())
	} else {
		state.BfdEnabled = types.BoolNull()
	}

	if neighbor.BfdInterval.IsSet() {
		state.BfdInterval = convert.Int64ToType(neighbor.BfdInterval.Get())
	} else {
		state.BfdInterval = types.Int64Null()
	}

	if neighbor.BfdMultiple.IsSet() {
		state.BfdMultiple = convert.Int64ToType(neighbor.BfdMultiple.Get())
	} else {
		state.BfdMultiple = types.Int64Null()
	}

	if neighbor.AllowAsIn.IsSet() {
		state.AllowAsIn = convert.BoolToType(neighbor.AllowAsIn.Get())
	} else {
		state.AllowAsIn = types.BoolNull()
	}

	if neighbor.HopLimit.IsSet() {
		state.HopLimit = convert.Int64ToType(neighbor.HopLimit.Get())
	} else {
		state.HopLimit = types.Int64Null()
	}

	if neighbor.RestartMode.IsSet() {
		state.RestartMode = convert.StrToType(neighbor.RestartMode.Get())
	} else {
		state.RestartMode = types.StringNull()
	}

	// Password is write-only — never read back from API
	state.PasswordWo = types.StringNull()
	state.PasswordWoVersion = plan.PasswordWoVersion

	// Handle config
	state.Config = types.DynamicNull()
	state.ConfigNsxt = NewConfigNsxtValueNull()
	state.ConfigNsxv = NewConfigNsxvValueNull()

	if neighbor.Config != nil {
		cfg := neighbor.Config

		if cfg.NSXTBGPNeighborConfig2 != nil && !plan.ConfigNsxt.IsNull() {
			nsxt := cfg.NSXTBGPNeighborConfig2
			sourceAddrs := nsxt.SourceAddresses

			var addrValues []attr.Value
			for _, addr := range sourceAddrs {
				addrValues = append(addrValues, types.StringValue(addr))
			}

			sourceAddressesSet, setDiags := types.SetValue(types.StringType, addrValues)
			diags.Append(setDiags...)

			if !diags.HasError() {
				configNsxt, d := NewConfigNsxtValue(
					ConfigNsxtValue{}.AttributeTypes(ctx),
					map[string]attr.Value{
						"source_addresses": sourceAddressesSet,
					},
				)
				diags.Append(d...)

				if !diags.HasError() {
					state.ConfigNsxt = configNsxt
				}
			}
		}

		if cfg.NSXVBGPNeighborConfig2 != nil && !plan.ConfigNsxv.IsNull() {
			nsxv := cfg.NSXVBGPNeighborConfig2

			configNsxv, d := NewConfigNsxvValue(
				ConfigNsxvValue{}.AttributeTypes(ctx),
				map[string]attr.Value{
					"router_id": convert.StrToType(nsxv.RouterId),
					"interface": convert.StrToType(nsxv.Interface),
				},
			)
			diags.Append(d...)
			if !diags.HasError() {
				state.ConfigNsxv = configNsxv
			}
		}

		if !plan.Config.IsNull() && !plan.Config.IsUnknown() {
			state.Config = plan.Config
		}
	}

	return state, diags
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var plan NetworkRouterBgpNeighborModel

	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			readOperation,
			"failed to create client: "+err.Error(),
		)

		return
	}

	id := plan.Id.ValueInt64()
	routerID := plan.RouterId.ValueInt64()

	state, diags := getNetworkRouterBgpNeighborAsState(ctx, id, routerID, client, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

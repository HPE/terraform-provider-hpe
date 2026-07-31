// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterbgpneighbor

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

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
	} else if !plan.Description.IsUnknown() {
		// The API accepts description on create/update but never returns it, so
		// preserve the configured (or prior) value. Without this a configured
		// description is nulled out and the apply fails with "Provider produced
		// inconsistent result after apply". The IsUnknown guard matters because
		// description is Optional+Computed: an unset value arrives unknown in
		// the plan and must never be copied into state.
		state.Description = plan.Description
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

	// Handle config.
	//
	// Which typed config block to populate is driven by the content of the API
	// response first, and only then by the plan. Read has to be self-sufficient:
	// on import the plan is empty (ImportState seeds only router_id and id), so
	// a plan-gated branch silently drops the block and ImportStateVerify fails.
	//
	// Content-based discrimination is also required because the config schema is
	// an anyOf and every variant carries a `,remain` catch-all for additional
	// properties. An NSX-T payload such as {"sourceAddresses":[...]} therefore
	// decodes into a non-nil NSXVBGPNeighborConfig2 as well (the unknown key
	// lands in AdditionalProperties, so the struct is non-zero and survives the
	// SDK's IsEmpty check). Testing for the variant's own fields is the only
	// reliable way to tell the two apart.
	//
	// The plan is still consulted as a secondary trigger so that a configured
	// but empty block (both attributes are Optional+Computed) resolves to a
	// non-null object rather than null, which would fail the "inconsistent
	// result after apply" check. Note this must test IsNull *and* IsUnknown: an
	// unset Optional+Computed block arrives unknown, not null, so an IsNull-only
	// guard treats "user did not set it" as "user set it".
	state.Config = types.DynamicNull()
	state.ConfigNsxt = NewConfigNsxtValueNull()
	state.ConfigNsxv = NewConfigNsxvValueNull()

	if neighbor.Config != nil {
		cfg := neighbor.Config

		nsxtInPlan := !plan.ConfigNsxt.IsNull() && !plan.ConfigNsxt.IsUnknown()
		nsxvInPlan := !plan.ConfigNsxv.IsNull() && !plan.ConfigNsxv.IsUnknown()

		if nsxt := cfg.NSXTBGPNeighborConfig2; nsxt != nil &&
			(len(nsxt.SourceAddresses) > 0 || nsxtInPlan) {
			var addrValues []attr.Value
			for _, addr := range nsxt.SourceAddresses {
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

		if nsxv := cfg.NSXVBGPNeighborConfig2; nsxv != nil &&
			(nsxv.RouterId != nil || nsxv.Interface != nil || nsxvInPlan) {
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

// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterroute

import (
	"context"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const readOperation = "read network router route resource"

func getNetworkRouterRouteAsState(
	ctx context.Context,
	id int64,
	routerID int64,
	client *sdk.APIClient,
	plan NetworkRouterRouteModel,
) (NetworkRouterRouteModel, diag.Diagnostics) {
	var state NetworkRouterRouteModel
	var diags diag.Diagnostics

	resp, hresp, err := client.NetworksAPI.
		GetNetworkRouterRoute(ctx, id, routerID).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		diags.AddError(
			readOperation,
			fmt.Sprintf("route %d GET failed: ", id)+
				errfmt.ErrMsg(err, hresp),
		)

		return state, diags
	}

	route := resp.GetNetworkRoute()

	state.Id = convert.Int64ToType(route.Id)
	state.RouterId = plan.RouterId
	state.Name = convert.StrToType(route.Name)

	if route.IsSetCode() {
		state.Code = types.StringValue(route.GetCode())
	} else {
		state.Code = types.StringNull()
	}

	if route.IsSetDescription() {
		state.Description = types.StringValue(route.GetDescription())
	} else {
		state.Description = types.StringNull()
	}

	state.Enabled = convert.BoolToType(route.Enabled)
	state.DefaultRoute = convert.BoolToType(route.DefaultRoute)
	state.Network = convert.StrToType(route.Source)
	state.NextHop = convert.StrToType(route.Destination)
	state.RouteType = convert.StrToType(route.RouteType)
	state.SourceType = convert.StrToType(route.SourceType)
	state.ExternalId = convert.StrToType(route.ExternalId)
	state.ProviderId = convert.StrToType(route.ProviderId)

	// TODO: make this not a float32 in sdk....
	if route.IsSetNetworkMtu() {
		mtuVal := route.GetNetworkMtu()
		state.Mtu = types.Int64Value(int64(mtuVal))
	} else {
		state.Mtu = types.Int64Null()
	}

	return state, diags
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var plan NetworkRouterRouteModel

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

	state, diags := getNetworkRouterRouteAsState(ctx, id, routerID, client, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

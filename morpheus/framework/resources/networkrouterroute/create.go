// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterroute

import (
	"context"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
)

const createOperation = "create network router route resource"

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan NetworkRouterRouteModel

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

	route := sdk.NewCreateNetworkRouterRouteRequestNetworkRouteWithDefaults()

	routerID := plan.RouterId.ValueInt64()

	if !plan.Mtu.IsNull() && !plan.Mtu.IsUnknown() {
		route.SetNetworkMtu(float32(plan.Mtu.ValueInt64()))
	}

	if !plan.Network.IsNull() && !plan.Network.IsUnknown() {
		route.SetSource(plan.Network.ValueString())
	}

	if !plan.NextHop.IsNull() && !plan.NextHop.IsUnknown() {
		route.SetDestination(plan.NextHop.ValueString())
	}

	route.SetName(plan.Name.ValueString())

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		route.SetDescription(plan.Description.ValueString())
	}

	if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() {
		route.SetEnabled(plan.Enabled.ValueBool())
	}

	if !plan.DefaultRoute.IsNull() && !plan.DefaultRoute.IsUnknown() {
		route.SetDefaultRoute(plan.DefaultRoute.ValueBool())
	}

	createReq := sdk.NewCreateNetworkRouterRouteRequest()
	createReq.SetNetworkRoute(*route)

	result, hresp, err := client.NetworksAPI.
		CreateNetworkRouterRoute(ctx, routerID).
		CreateNetworkRouterRouteRequest(*createReq).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			createOperation,
			fmt.Sprintf("router %d route POST failed: %s",
				routerID, errfmt.ErrMsg(err, hresp)),
		)

		return
	}

	id := result.GetId()
	plan.Id = types.Int64Value(id)

	taintResourceState := func(id int64) {
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "network_router_route",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})
	}

	state, pdiags := getNetworkRouterRouteAsState(ctx, id, routerID, client, plan)
	if pdiags.HasError() {
		resp.Diagnostics.Append(pdiags...)
		resp.Diagnostics.AddError(
			createOperation,
			fmt.Sprintf("route %d was created but could not be read", id),
		)
		taintResourceState(id)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError(
			createOperation,
			fmt.Sprintf("route %d was created but state could not be saved", id),
		)
		taintResourceState(id)

		return
	}
}

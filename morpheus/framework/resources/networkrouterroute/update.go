// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterroute

import (
	"context"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

const updateOperation = "update network router route resource"

// Update deletes and recreates the route because the API does not support
// an update endpoint for network router routes.
func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state NetworkRouterRouteModel

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

	oldID := state.Id.ValueInt64()
	routerID := state.RouterId.ValueInt64()

	// Delete the existing route
	tflog.Debug(ctx, fmt.Sprintf("Deleting route %d on router %d for replacement", oldID, routerID))

	_, hresp, err := client.NetworksAPI.
		DeleteNetworkRouterRoute(ctx, oldID, routerID).Execute()
	if err != nil {
		if hresp == nil || hresp.StatusCode != http.StatusNotFound {
			resp.Diagnostics.AddError(
				updateOperation,
				fmt.Sprintf("route %d DELETE failed: %s",
					oldID, errfmt.ErrMsg(err, hresp)),
			)

			return
		}

		tflog.Debug(ctx, fmt.Sprintf("Route %d already deleted (404)", oldID))
	}

	// Recreate with new values
	route := sdk.NewCreateNetworkRouterRouteRequestNetworkRouteWithDefaults()

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

	newRouterID := plan.RouterId.ValueInt64()

	result, hresp, err := client.NetworksAPI.
		CreateNetworkRouterRoute(ctx, newRouterID).
		CreateNetworkRouterRouteRequest(*createReq).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			updateOperation,
			fmt.Sprintf("router %d route POST (recreate) failed: %s",
				newRouterID, errfmt.ErrMsg(err, hresp)),
		)

		return
	}

	newID := result.GetId()
	plan.Id = types.Int64Value(newID)

	newState, diags := getNetworkRouterRouteAsState(ctx, newID, newRouterID, client, plan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		resp.Diagnostics.AddError(
			updateOperation,
			fmt.Sprintf("route %d: failed to read from api after recreate", newID),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

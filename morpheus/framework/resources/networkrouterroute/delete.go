// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterroute

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state NetworkRouterRouteModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"delete network router route resource",
			"failed to create client: "+err.Error(),
		)

		return
	}

	id := state.Id.ValueInt64()
	routerID := state.RouterId.ValueInt64()

	tflog.Debug(ctx, fmt.Sprintf("Deleting route %d on router %d", id, routerID))

	_, hresp, err := client.NetworksAPI.
		DeleteNetworkRouterRoute(ctx, id, routerID).Execute()
	if err != nil {
		if hresp != nil && hresp.StatusCode == http.StatusNotFound {
			tflog.Debug(ctx, fmt.Sprintf("Route %d already deleted (404)", id))

			return
		}

		resp.Diagnostics.AddError(
			"delete network router route resource",
			fmt.Sprintf("route %d DELETE failed: %s",
				id, errfmt.ErrMsg(err, hresp)),
		)
	}
}

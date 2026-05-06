// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterbgpneighbor

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
	var state NetworkRouterBgpNeighborModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"delete network router bgp neighbor resource",
			"failed to create client: "+err.Error(),
		)

		return
	}

	id := state.Id.ValueInt64()
	routerID := state.RouterId.ValueInt64()

	tflog.Debug(ctx, fmt.Sprintf("Deleting BGP neighbor %d on router %d", id, routerID))

	_, hresp, err := client.NetworksAPI.
		DeleteNetworkRouterBgpNeighbor(ctx, id, routerID).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"delete network router bgp neighbor resource",
			fmt.Sprintf("bgp neighbor %d DELETE failed: %s",
				id, errfmt.ErrMsg(err, hresp)),
		)
	}
}

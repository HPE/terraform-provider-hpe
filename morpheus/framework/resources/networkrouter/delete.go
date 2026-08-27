// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouter

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state NetworkRouterModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"delete network router resource",
			"failed to create client: "+err.Error(),
		)

		return
	}

	id := state.Id.ValueInt64()

	_, hresp, err := client.NetworksAPI.DeleteNetworkRouter(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"delete network router resource",
			fmt.Sprintf("network router %d DELETE failed: %s",
				id, errfmt.ErrMsg(err, hresp)),
		)
	}
}

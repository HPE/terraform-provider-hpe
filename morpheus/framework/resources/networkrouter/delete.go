// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouter

import (
	"context"
	"fmt"

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

	result, hresp, err := client.NetworksAPI.DeleteNetworkRouter(ctx, id).Execute()
	if err := errfmt.CheckDeleted(result, err, hresp); err != nil {
		resp.Diagnostics.AddError(
			"delete network router resource",
			fmt.Sprintf("network router %d DELETE failed: %s",
				id, err.Error()),
		)
	}
}

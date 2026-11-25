// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkfloatingip

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

func (r *networkFloatingIpResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state NetworkFloatingIpModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueInt64()

	// Releasing a floating IP is a PUT .../release. A 404 means it is already
	// gone. A 400 is returned while a VM is still attached ("Unable to release
	// floating IP because it is still assigned."); surface that message so the
	// user knows to detach first.
	// CheckDeleted treats a 404 as already released and a 200 carrying
	// success=false (the plugin reporting a refusal) as a failure with its msg.
	result, httpResp, err := client.NetworksAPI.ReleaseNetworkFloatingIp(ctx, id).Execute()
	if err := errfmt.CheckDeleted(result, err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "network_floating_ip", "", err, httpResp)

		return
	}
}

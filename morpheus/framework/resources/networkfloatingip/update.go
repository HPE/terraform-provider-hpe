// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkfloatingip

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

// Update sends nothing: every user-settable attribute (network_server_id,
// floating_ip_pool_id) is RequiresReplace, so a change to any of them destroys
// and recreates the resource rather than reaching Update. It still reads the
// floating IP back and maps it, as Read does, so the computed attributes
// (ip_status, server_id, …) come from the appliance rather than from the plan
// — which matters the day an in-place attribute is added.
func (r *networkFloatingIpResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan NetworkFloatingIpModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	id := plan.Id.ValueInt64()

	result, httpResp, err := client.NetworksAPI.GetNetworkFloatingIp(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "network_floating_ip", "", err, httpResp)

		return
	}

	if result.NetworkFloatingIp == nil {
		resp.Diagnostics.AddError("API returned nil", "NetworkFloatingIp is nil in the response")

		return
	}

	// network_server_id and floating_ip_pool_id are not returned by the API;
	// the plan carries them.
	state := floatingIpAsState(result.NetworkFloatingIp, &plan)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

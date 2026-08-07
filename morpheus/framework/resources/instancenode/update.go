// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instancenode

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/containerip"
)

const defaultUpdateTimeout = 45 * time.Minute

// Update implements resource.Resource.
// Only wait_for_ip_address can change without replacement. When it flips
// from false to true, run the IP wait.
func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan instanceNodeModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateTimeout, diags := plan.Timeouts.Update(ctx, defaultUpdateTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state instanceNodeModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Carry forward computed values.
	plan.ID = state.ID
	plan.ContainerID = state.ContainerID
	plan.ServerID = state.ServerID
	plan.IPAddress = state.IPAddress
	plan.Name = state.Name
	plan.UUID = state.UUID

	if plan.WaitForIPAddress.ValueBool() &&
		(!state.IPAddress.IsNull() && containerip.Ready(state.IPAddress.ValueString())) {
		// Already have a valid IP — no wait needed.
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

		return
	}

	if plan.WaitForIPAddress.ValueBool() {
		client, err := r.NewClient(ctx)
		if err != nil {
			resp.Diagnostics.AddError("failed to create API client", err.Error())

			return
		}

		ip, warned, waitErr := containerip.Wait(
			ctx, client,
			plan.InstanceID.ValueInt64(),
			plan.ContainerID.ValueInt64(),
			updateTimeout,
		)
		if waitErr != nil {
			resp.Diagnostics.AddError(
				"IP address wait failed",
				fmt.Sprintf("container %d: %s",
					plan.ContainerID.ValueInt64(), waitErr.Error()),
			)

			return
		}

		if !warned && ip != "" {
			plan.IPAddress = types.StringValue(ip)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

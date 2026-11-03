// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instancenode

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/containerip"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
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
	plan.ContainerUUID = state.ContainerUUID

	// ssh_password is create-only (a change replaces the node); an update
	// planned through write_only_revision only records the digests.
	plan.WriteOnlyRevision = state.WriteOnlyRevision
	if wodigest.Resend(ctx, req.Plan) {
		plan.WriteOnlyRevision = wodigest.NewRevision()
	}

	// Record the configured-value digests and re-baseline the appliance-hash
	// drift detection off a fresh read of the node's server (see
	// write_only.go). The client is best-effort: on a construction failure the
	// baseline is simply not refreshed, which never blocks the update.
	var serverResp *sdk.GetHost200ResponseServer
	if client, cerr := r.NewClient(ctx); cerr == nil {
		serverResp = serverForDrift(ctx, client, plan.ServerID.ValueInt64())
	}
	resp.Diagnostics.Append(wo.StoreObserved(ctx, resp.Private, req.Config, serverResp)...)

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

// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkfloatingip

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
)

func (r *networkFloatingIpResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan NetworkFloatingIpModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	allocReq := sdk.AllocateNetworkFloatingIpRequest{
		NetworkServerId:  plan.NetworkServerId.ValueInt64Pointer(),
		FloatingIpPoolId: plan.FloatingIpPoolId.ValueInt64Pointer(),
	}

	result, httpResp, err := client.NetworksAPI.AllocateNetworkFloatingIp(ctx).
		AllocateNetworkFloatingIpRequest(allocReq).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "network_floating_ip", "", err, httpResp)

		return
	}

	// The endpoint returns HTTP 200 even for a business-rule failure, signalled
	// by success=false and a human-readable msg. Surface that message.
	if result.Success != nil && !*result.Success {
		errfmt.DiagErrorf(
			&resp.Diagnostics, errfmt.OpCreate, "network_floating_ip",
			"the appliance refused to allocate the floating IP: %s", successMsg(result.AdditionalProperties),
		)

		return
	}

	if result.NetworkFloatingIp == nil || result.NetworkFloatingIp.Id == nil {
		resp.Diagnostics.AddError("API returned nil ID", "NetworkFloatingIp ID is nil in the allocate response")

		return
	}

	id := *result.NetworkFloatingIp.Id

	// The allocate response omits some fields the pool fills in later, so do a
	// full read by id for the authoritative state.
	readResult, httpResp, err := client.NetworksAPI.GetNetworkFloatingIp(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "network_floating_ip", "", err, httpResp)
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "network_floating_ip",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	if readResult.NetworkFloatingIp == nil {
		resp.Diagnostics.AddError("API returned nil", "NetworkFloatingIp is nil in the response")
		// The create succeeded; keep the id so the next apply can finish or
		// destroy it instead of orphaning it on the appliance.
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "network_floating_ip",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	state := floatingIpAsState(readResult.NetworkFloatingIp, &plan)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// successMsg pulls the human-readable "msg" the appliance returns alongside a
// success=false envelope; it falls back to a generic string when absent.
func successMsg(props map[string]interface{}) string {
	if props == nil {
		return "no message returned"
	}

	if msg, ok := props["msg"].(string); ok && msg != "" {
		return msg
	}

	return fmt.Sprintf("%v", props)
}

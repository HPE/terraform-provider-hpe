// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancerprofile

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// Import ID format: "loadBalancerId/profileId"
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"import load_balancer_profile resource",
			fmt.Sprintf(
				"provided import ID %q is invalid: expected format 'loadBalancerId/profileId'",
				req.ID,
			),
		)

		return
	}

	loadBalancerID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"import load_balancer_profile resource",
			fmt.Sprintf("invalid load balancer ID %q: %s", parts[0], err),
		)

		return
	}

	profileID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"import load_balancer_profile resource",
			fmt.Sprintf("invalid profile ID %q: %s", parts[1], err),
		)

		return
	}

	resp.Diagnostics.Append(
		resp.State.SetAttribute(ctx, path.Root("load_balancer_id"), loadBalancerID)...,
	)
	resp.Diagnostics.Append(
		resp.State.SetAttribute(ctx, path.Root("id"), profileID)...,
	)
}

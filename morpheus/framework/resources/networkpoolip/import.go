// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkpoolip

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
	// Import ID format: "<network_pool_id>:<id>".
	parts := strings.SplitN(req.ID, ":", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"import network_pool_ip resource",
			fmt.Sprintf(
				"provided import ID %q is invalid: expected format '<network_pool_id>:<id>'",
				req.ID,
			),
		)

		return
	}

	networkPoolID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"import network_pool_ip resource",
			fmt.Sprintf("invalid network pool ID %q: %s", parts[0], err),
		)

		return
	}

	id, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"import network_pool_ip resource",
			fmt.Sprintf("invalid reservation ID %q: %s", parts[1], err),
		)

		return
	}

	resp.Diagnostics.Append(
		resp.State.SetAttribute(ctx, path.Root("network_pool_id"), networkPoolID)...,
	)
	resp.Diagnostics.Append(
		resp.State.SetAttribute(ctx, path.Root("id"), id)...,
	)
}

// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancervirtualserver

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ImportState supports importing with the format "loadBalancerId.id".
func (r *Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	parts := strings.SplitN(req.ID, ".", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"import load balancer virtual server",
			fmt.Sprintf(
				"provided import ID %q is invalid; expected format: load_balancer_id.id",
				req.ID,
			),
		)

		return
	}

	lbID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"import load balancer virtual server",
			fmt.Sprintf("load_balancer_id %q is not a valid integer", parts[0]),
		)

		return
	}

	vsID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"import load balancer virtual server",
			fmt.Sprintf("id %q is not a valid integer", parts[1]),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("load_balancer_id"),
		types.Int64Value(lbID),
	)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("id"), vsID,
	)...)
}

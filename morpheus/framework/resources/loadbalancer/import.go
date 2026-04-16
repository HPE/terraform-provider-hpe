// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancer

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"import load balancer resource",
			"provided import ID '"+req.ID+"' is invalid (non-number)",
		)

		return
	}

	diags := resp.State.SetAttribute(ctx, path.Root("id"), id)
	resp.Diagnostics.Append(diags...)
}

// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterroute

import (
	"context"
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
	parts := strings.SplitN(req.ID, ".", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"import network router route resource",
			"provided import ID '"+req.ID+
				"' is invalid, expected format 'router_id.id'",
		)

		return
	}

	routerID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"import network router route resource",
			"provided router_id '"+parts[0]+"' is invalid (non-number)",
		)

		return
	}

	id, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"import network router route resource",
			"provided id '"+parts[1]+"' is invalid (non-number)",
		)

		return
	}

	resp.Diagnostics.Append(
		resp.State.SetAttribute(ctx, path.Root("router_id"), routerID)...,
	)
	resp.Diagnostics.Append(
		resp.State.SetAttribute(ctx, path.Root("id"), id)...,
	)
}

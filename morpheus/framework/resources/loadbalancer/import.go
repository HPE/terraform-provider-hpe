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
	// "<id>" or "<id>#group_id=<sha256 hex>&network_server_id=<sha256 hex>":
	// see wodigest.
	rawID, ok := wo.Import(ctx, req, resp)
	if !ok {
		return
	}

	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"import load balancer resource",
			"provided import ID '"+rawID+"' is invalid (non-number)",
		)

		return
	}

	diags := resp.State.SetAttribute(ctx, path.Root("id"), id)
	resp.Diagnostics.Append(diags...)
}

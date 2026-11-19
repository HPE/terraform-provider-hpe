// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouter

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
	// "<id>" or "<id>#config_juniper_bgp.md5_key_wo=<sha256 hex>": see wodigest.
	rawID, ok := wo.Import(ctx, req, resp)
	if !ok {
		return
	}

	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"import network router resource",
			"provided import ID '"+rawID+"' is invalid (non-number)",
		)

		return
	}

	diags := resp.State.SetAttribute(ctx, path.Root("id"), id)
	resp.Diagnostics.Append(diags...)
}

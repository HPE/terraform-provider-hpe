// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package image

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
	// "<id>" or "<id>#ssh_password_wo=<sha256 hex>": see wodigest.
	rawID, ok := wo.Import(ctx, req, resp)
	if !ok {
		return
	}

	id, err := strconv.Atoi(rawID)
	if err != nil {
		resp.Diagnostics.AddError(
			"import image resource",
			"provided import ID '"+rawID+"' is invalid (non-number)",
		)

		return
	}

	diags := resp.State.SetAttribute(
		ctx, path.Root("id"), id,
	)
	resp.Diagnostics.Append(diags...)
}

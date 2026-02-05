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
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"import image resource",
			"provided import ID '"+req.ID+"' is invalid (non-number)",
		)

		return
	}

	diags := resp.State.SetAttribute(
		ctx, path.Root("id"), id,
	)
	resp.Diagnostics.Append(diags...)
}

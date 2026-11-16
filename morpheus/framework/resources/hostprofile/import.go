// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package hostprofile

import (
	"context"
	"fmt"
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
			"import host profile resource",
			fmt.Sprintf("invalid host profile id %q: %s", req.ID, err),
		)

		return
	}

	resp.Diagnostics.Append(
		resp.State.SetAttribute(ctx, path.Root("id"), id)...,
	)
}

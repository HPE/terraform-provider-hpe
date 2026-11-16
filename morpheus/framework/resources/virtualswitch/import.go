// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package virtualswitch

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ImportState imports a virtual switch by its numeric id. cluster_id and the
// segment map keys are populated by the Read that follows (keys synthesised by
// carrySegmentKeys).
func (r *Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"import virtual_switch resource",
			"provided import ID '"+req.ID+"' is invalid (non-number)",
		)

		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

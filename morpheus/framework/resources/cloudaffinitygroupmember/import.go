// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cloudaffinitygroupmember

import (
	"context"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ImportState accepts cloud_id/affinity_group_id/server_id.
//
// A membership has no identifier of its own -- it is the pairing that exists,
// not an object -- so the import id has to carry all three parts.
func (r *Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	parts := strings.Split(req.ID, "/")

	const wantParts = 3

	if len(parts) != wantParts {
		resp.Diagnostics.AddAttributeError(
			importIDPath,
			"unexpected import identifier",
			"Expected cloud_id/affinity_group_id/server_id, for example 2/6045/572255, "+
				"but got: "+req.ID,
		)

		return
	}

	names := []string{"cloud_id", "affinity_group_id", "server_id"}
	values := make([]int64, wantParts)

	for i, part := range parts {
		v, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				importIDPath,
				"unexpected import identifier",
				"Expected "+names[i]+" to be a number in cloud_id/affinity_group_id/server_id, "+
					"but got: "+part,
			)

			return
		}

		values[i] = v
	}

	state := memberModel{
		ID:              types.StringValue(membershipID(values[0], values[1], values[2])),
		CloudID:         types.Int64Value(values[0]),
		AffinityGroupID: types.Int64Value(values[1]),
		ServerID:        types.Int64Value(values[2]),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

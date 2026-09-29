// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package tenant

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

const importOperation = "import tenant resource"

func (r *Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			importOperation,
			fmt.Sprintf("provided import ID %q is invalid, expected a numeric tenant ID", req.ID),
		)

		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	// Seed prior with an explicit null for remove_resources, the only attribute
	// the GET does not return. It is a delete-time flag unknown to the API, so
	// it is null after import; re-apply to set it. parent_id is
	// computed_optional and comes from the read response.
	prior := TenantModel{
		RemoveResources: types.BoolNull(),
	}

	state, notFound, diags := getTenantAsState(ctx, id, client, prior)
	if notFound {
		resp.Diagnostics.AddError(
			importOperation,
			fmt.Sprintf("tenant %d not found", id),
		)

		return
	}

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

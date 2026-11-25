// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package hostprofile

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state HostProfileModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The seeded HVM default profile cannot be deleted on the API.
	if state.Code.ValueString() == seedCode {
		resp.Diagnostics.AddError(
			"delete host profile resource",
			"the seeded default host profile '"+seedCode+
				"' is immutable and cannot be deleted",
		)

		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"delete host profile resource",
			"failed to create client: "+err.Error(),
		)

		return
	}

	id := state.Id.ValueInt64()

	result, hresp, err := client.HostProfilesAPI.
		DeleteHostProfile(ctx, id).Execute()
	if err := errfmt.CheckDeleted(result, err, hresp); err != nil {
		resp.Diagnostics.AddError(
			"delete host profile resource",
			"host profile DELETE failed: "+err.Error(),
		)

		return
	}
}

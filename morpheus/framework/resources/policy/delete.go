// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package policy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var data PolicyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id.ValueInt64()

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"delete policy resource",
			fmt.Sprintf("policy %d: failed to create client: %s", id, err.Error()),
		)

		return
	}

	result, hresp, err := client.PoliciesAPI.RemovePolicies(ctx, id).Execute()
	if err := errfmt.CheckDeleted(result, err, hresp); err != nil {
		resp.Diagnostics.AddError(
			"delete policy resource",
			fmt.Sprintf("policy %d: DELETE failed ", id)+err.Error(),
		)

		return
	}
}

// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package cloud

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
	var data CloudModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id.ValueInt64()

	client, _ := r.NewClient(ctx)

	result, hresp, err := client.CloudsAPI.RemoveClouds(ctx, id).Force(true).Execute()
	if err := errfmt.CheckDeleted(result, err, hresp); err != nil {
		resp.Diagnostics.AddError(
			"delete cloud resource",
			fmt.Sprintf("cloud %d: DELETE failed ", id)+err.Error(),
		)

		return
	}
}

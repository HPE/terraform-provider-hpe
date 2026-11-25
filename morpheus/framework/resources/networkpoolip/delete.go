// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkpoolip

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
	var data NetworkPoolIpModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	networkPoolID := data.NetworkPoolId.ValueInt64()
	id := data.Id.ValueInt64()

	// CheckDeleted judges the DELETE by status and body: a 404 is success
	// (already gone), and a 200 with success:false is a refusal that must be
	// an error so state is not cleared while the reservation remains
	// (MORPH-17987).
	result, httpResp, err := client.NetworksAPI.
		DeleteNetworkPoolIp(ctx, networkPoolID, id).Execute()
	if err := errfmt.CheckDeleted(result, err, httpResp); err != nil {
		resp.Diagnostics.AddError(
			"error deleting network pool IP",
			fmt.Sprintf("network pool IP %d DELETE failed: ", id)+err.Error(),
		)
	}
}

// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package network

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/constants"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/errors"
)

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state NetworkModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"delete network resource",
			"failed to create client: "+err.Error(),
		)

		return
	}

	id := state.Id.ValueInt64()

	// Create a custom timeout for the delete operation
	// Mainly needed for GCP network delete, which is
	// synchronous, and can take some time
	deleteCtx, cancel := context.WithTimeout(ctx, constants.NetworkDeleteTimeout)
	defer cancel()

	tflog.Debug(ctx, fmt.Sprintf("Deleting network %d", id))
	_, hresp, err := client.NetworksAPI.DeleteNetwork(deleteCtx, id).
		Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"delete network resource",
			fmt.Sprintf("network %d DELETE failed: %s",
				id, errors.ErrMsg(err, hresp)),
		)
	}
}

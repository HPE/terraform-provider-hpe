// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package network

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/constants"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
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

	// Create a custom timeout for the delete operation
	// Mainly needed for GCP network delete, which is
	// synchronous, and can take some time
	// TODO: redo this, we should use the morpheus client
	// with a context to set timeouts, rather than digging
	// into the http client like this
	if httpClient := client.GetConfig().HTTPClient; httpClient != nil {
		httpClient.Timeout = constants.NetworkDeleteTimeout
	}

	id := state.Id.ValueInt64()

	tflog.Debug(ctx, fmt.Sprintf("Deleting network %d", id))
	_, hresp, err := client.NetworksAPI.DeleteNetwork(ctx, id).
		Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"delete network resource",
			fmt.Sprintf("network %d DELETE failed: %s",
				id, errfmt.ErrMsg(err, hresp)),
		)
	}
}

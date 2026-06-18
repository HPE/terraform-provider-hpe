// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instanceclone

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state InstanceCloneModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteTimeout, diags := state.Timeouts.Delete(ctx, 10*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"delete instance clone",
			"failed to create client: "+err.Error(),
		)

		return
	}

	instanceID := state.Id.ValueInt64()

	_, hresp, err := client.InstancesAPI.DeleteInstance(ctx, instanceID).Execute()
	if err != nil {
		if hresp != nil && hresp.StatusCode == http.StatusNotFound {
			tflog.Warn(ctx, fmt.Sprintf("Instance clone %d already gone", instanceID))

			return
		}

		resp.Diagnostics.AddError(
			"delete instance clone",
			fmt.Sprintf("instance %d DELETE failed: %s",
				instanceID, errfmt.ErrMsg(err, hresp)),
		)

		return
	}
}

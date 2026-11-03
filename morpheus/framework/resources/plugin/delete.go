// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package plugin

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state PluginModel
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
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	id := state.Id.ValueInt64()

	// CheckDeleted treats a 404 as already gone and also inspects the
	// {"success": false, ...} body Morpheus returns with a 200 on failure.
	result, httpResp, err := client.PluginsAPI.RemovePlugin(ctx, id).Execute()
	if err := errfmt.CheckDeleted(result, err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, resourceName, fmt.Sprint(id), err, httpResp)

		return
	}

	// Wait for the GET to turn 404.
	what := fmt.Sprintf("plugin %d", id)
	waitForDeleted := func() (struct{}, error) {
		_, httpResp, err := client.PluginsAPI.GetPlugin(ctx, id).Execute()

		return struct{}{}, errfmt.PollDeleted(what, err, httpResp)
	}

	if _, err := backoff.Retry(
		ctx,
		waitForDeleted,
		backoff.WithBackOff(backoff.NewConstantBackOff(3*time.Second)),
		backoff.WithMaxElapsedTime(deleteTimeout),
	); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, resourceName, fmt.Sprint(id), err, nil)
	}
}

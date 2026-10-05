// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cluster

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ClusterModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get timeout from HCL if set, the default is 45 minutes
	deleteTimeout, diags := data.Timeouts.Delete(ctx, 45*time.Minute)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

	id := data.Id.ValueInt64()
	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	// We'll do a force delete so that we can delete
	// clusters that are in a provisioning state.
	deleteReq := client.ClustersAPI.DeleteCluster(ctx, id).Force("on")
	result, hresp, err := deleteReq.Execute()
	if err := errfmt.CheckDeleted(result, err, hresp); err != nil {
		resp.Diagnostics.AddError(
			"delete cluster resource",
			fmt.Sprintf("cluster %d: DELETE failed ", id)+err.Error(),
		)

		return
	}

	// Morpheus accepts the delete and removes the cluster in a background
	// job; its server group stays readable (status deprovisioning) until that
	// job finishes. Wait for the GET to turn 404 -- a create that reuses the
	// name races the job otherwise ("Cluster name must be unique",
	// MORPH-18379). Morpheus is left to delete the cluster's servers.
	what := fmt.Sprintf("cluster %d", id)
	waitForDeleted := func() (struct{}, error) {
		_, httpResp, err := client.ClustersAPI.GetCluster(ctx, id).Execute()

		return struct{}{}, errfmt.PollDeleted(what, err, httpResp)
	}

	if _, err := backoff.Retry(
		ctx,
		waitForDeleted,
		backoff.WithBackOff(backoff.NewConstantBackOff(5*time.Second)),
		backoff.WithMaxElapsedTime(deleteTimeout),
	); err != nil {
		resp.Diagnostics.AddError(
			"delete cluster resource",
			fmt.Sprintf("cluster %d: DELETE failed ", id)+err.Error(),
		)
	}
}

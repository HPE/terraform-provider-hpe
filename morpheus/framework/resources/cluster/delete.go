// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cluster

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

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
	_, hresp, err := deleteReq.Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"delete cluster resource",
			fmt.Sprintf("cluster %d: DELETE failed ", id)+errfmt.ErrMsg(err, hresp),
		)

		return
	}

	// We'll just check that the cluster itself has been deleted for now
	// and assume that Morpheus took care of deleting its servers.
	waitForDeleted := func() (*sdk.GetCluster200Response, error) {
		_, httpResp, err := client.ClustersAPI.GetCluster(ctx, id).Execute()
		// 404 status code counts as a successful delete
		if err != nil {
			if httpResp == nil || httpResp.StatusCode != http.StatusNotFound {
				return nil, backoff.Permanent(err)
			}
		}

		return nil, nil
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

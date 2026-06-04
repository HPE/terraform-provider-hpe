// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancerpool

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var data LoadBalancerPoolModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	loadBalancerID := data.LoadBalancerId.ValueInt64()
	id := data.Id.ValueInt64()

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	_, httpResp, err := client.LoadBalancersAPI.
		DeleteLoadBalancerPool(ctx, loadBalancerID, id).Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"error deleting load balancer pool",
			fmt.Sprintf("load balancer pool %d DELETE failed: ", id)+
				errfmt.ErrMsg(err, httpResp),
		)

		return
	}

	waitForDeleted := func() (struct{}, error) {
		_, httpResp, err := client.LoadBalancersAPI.
			GetLoadBalancerPool(ctx, loadBalancerID, id).Execute()
		if err != nil {
			if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
				return struct{}{}, nil
			}

			return struct{}{}, backoff.Permanent(err)
		}

		return struct{}{}, fmt.Errorf("load balancer pool %d still exists", id)
	}

	if _, err := backoff.Retry(
		ctx,
		waitForDeleted,
		backoff.WithBackOff(backoff.NewConstantBackOff(5*time.Second)),
		backoff.WithMaxElapsedTime(45*time.Minute),
	); err != nil {
		resp.Diagnostics.AddError(
			"error deleting load balancer pool",
			fmt.Sprintf("load balancer pool %d: DELETE confirmation failed: ", id)+
				err.Error(),
		)
	}
}

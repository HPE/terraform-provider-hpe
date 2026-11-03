// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancervirtualserver

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
	var data LoadBalancerVirtualServerModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("delete load balancer virtual server", "failed to create client: "+err.Error())

		return
	}

	lbID, err := loadBalancerIDFromInt64(data.LoadBalancerId)
	if err != nil {
		resp.Diagnostics.AddError("delete load balancer virtual server", err.Error())

		return
	}

	id := data.Id.ValueInt64()

	result, hresp, err := client.LoadBalancersAPI.
		DeleteLoadBalancerVirtualServer(ctx, lbID, id).
		Execute()
	if err := errfmt.CheckDeleted(result, err, hresp); err != nil {
		resp.Diagnostics.AddError(
			"error deleting load balancer virtual server",
			fmt.Sprintf("load balancer %d virtual server %d DELETE failed: %s",
				lbID, id, err.Error()),
		)
	}
}

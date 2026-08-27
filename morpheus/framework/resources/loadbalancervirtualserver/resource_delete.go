// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancervirtualserver

import (
	"context"
	"fmt"
	"net/http"

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

	_, hresp, err := client.LoadBalancersAPI.
		DeleteLoadBalancerVirtualServer(ctx, lbID, id).
		Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		// 404 means the resource is already gone — treat as success.
		if hresp != nil && hresp.StatusCode == http.StatusNotFound {
			return
		}

		resp.Diagnostics.AddError(
			"error deleting load balancer virtual server",
			fmt.Sprintf("load balancer %d virtual server %d DELETE failed: %s",
				lbID, id, errfmt.ErrMsg(err, hresp)),
		)
	}
}

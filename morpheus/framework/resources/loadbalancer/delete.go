// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancer

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state LoadBalancerModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"delete load balancer resource",
			"failed to create client: "+err.Error(),
		)

		return
	}

	id := state.Id.ValueInt64()

	tflog.Debug(ctx, fmt.Sprintf("Deleting load balancer %d", id))
	_, hresp, err := client.LoadBalancersAPI.DeleteLoadBalancer(ctx, id).
		Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		// 404 means the resource is already gone — treat as success
		if hresp != nil && hresp.StatusCode == http.StatusNotFound {
			return
		}

		resp.Diagnostics.AddError(
			"delete load balancer resource",
			fmt.Sprintf("load balancer %d DELETE failed: %s",
				id, errfmt.ErrMsg(err, hresp)),
		)
	}
}

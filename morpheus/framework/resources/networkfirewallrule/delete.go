// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkfirewallrule

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
	var data NetworkFirewallRuleModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id.ValueInt64()
	serverId := data.NetworkServerId.ValueInt64()

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	_, httpResp, err := client.NetworksAPI.
		DeleteNetworkFirewallRule(ctx, id, serverId).Execute()
	if err != nil || (httpResp != nil && httpResp.StatusCode != http.StatusOK) {
		resp.Diagnostics.AddError(
			"error deleting network firewall rule",
			fmt.Sprintf("network firewall rule %d DELETE failed: ", id)+errfmt.ErrMsg(err, httpResp),
		)

		return
	}

	waitForDeleted := func() (struct{}, error) {
		_, httpResp, err := client.NetworksAPI.
			GetNetworkFirewallRule(ctx, id, serverId).Execute()
		if err != nil {
			if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
				return struct{}{}, nil
			}

			return struct{}{}, backoff.Permanent(err)
		}

		// Resource still exists — keep retrying
		return struct{}{}, fmt.Errorf("network firewall rule %d still exists", id)
	}

	if _, err := backoff.Retry(
		ctx,
		waitForDeleted,
		backoff.WithBackOff(backoff.NewConstantBackOff(5*time.Second)),
		backoff.WithMaxElapsedTime(45*time.Minute),
	); err != nil {
		resp.Diagnostics.AddError(
			"error deleting network firewall rule",
			fmt.Sprintf("network firewall rule %d: DELETE confirmation failed: ", id)+err.Error(),
		)
	}
}

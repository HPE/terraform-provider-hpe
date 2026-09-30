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
	networkIntegrationId := data.NetworkIntegrationId.ValueInt64()

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	result, httpResp, err := client.NetworksAPI.
		DeleteNetworkFirewallRule(ctx, id, networkIntegrationId).Execute()
	if err := errfmt.CheckDeleted(result, err, httpResp); err != nil {
		resp.Diagnostics.AddError(
			"error deleting network firewall rule",
			fmt.Sprintf("network firewall rule %d DELETE failed: ", id)+err.Error(),
		)

		return
	}

	waitForDeleted := func() (struct{}, error) {
		_, httpResp, err := client.NetworksAPI.
			GetNetworkFirewallRule(ctx, id, networkIntegrationId).Execute()
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
		backoff.WithMaxElapsedTime(5*time.Minute),
	); err != nil {
		resp.Diagnostics.AddError(
			"error deleting network firewall rule",
			fmt.Sprintf("network firewall rule %d: DELETE confirmation failed: ", id)+err.Error(),
		)
	}
}

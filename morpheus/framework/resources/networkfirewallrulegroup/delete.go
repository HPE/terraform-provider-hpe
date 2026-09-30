// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkfirewallrulegroup

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var data NetworkFirewallRuleGroupModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	id := data.Id.ValueInt64()
	serverID := data.NetworkIntegrationId.ValueInt64()

	tflog.Debug(ctx, fmt.Sprintf(
		"Deleting network firewall rule group %d", id,
	))

	result, httpResp, err := client.NetworksAPI.
		DeleteNetworkFirewallRuleGroup(ctx, id, serverID).Execute()
	if err := errfmt.CheckDeleted(result, err, httpResp); err != nil {
		resp.Diagnostics.AddError(
			"error deleting network firewall rule group",
			fmt.Sprintf("network firewall rule group %d DELETE failed: ", id)+
				err.Error(),
		)
	}
}

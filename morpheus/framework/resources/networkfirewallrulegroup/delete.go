// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package networkfirewallrulegroup

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

	_, httpResp, err := client.NetworksAPI.
		DeleteNetworkFirewallRuleGroup(ctx, id, serverID).Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		// 404 means the resource is already gone — treat as success
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			return
		}

		resp.Diagnostics.AddError(
			"error deleting network firewall rule group",
			fmt.Sprintf("network firewall rule group %d DELETE failed: ", id)+
				errfmt.ErrMsg(err, httpResp),
		)
	}
}

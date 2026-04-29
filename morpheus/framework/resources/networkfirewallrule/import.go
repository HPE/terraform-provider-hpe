// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkfirewallrule

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	parts := strings.SplitN(req.ID, ":", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"import network_firewall_rule resource",
			fmt.Sprintf(
				"provided import ID %q is invalid; expected format 'network_integration_id:rule_id'",
				req.ID,
			),
		)

		return
	}

	serverId, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"import network_firewall_rule resource",
			"network_integration_id '"+parts[0]+"' is invalid (non-number)",
		)

		return
	}

	ruleId, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"import network_firewall_rule resource",
			"rule_id '"+parts[1]+"' is invalid (non-number)",
		)

		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_integration_id"), serverId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), ruleId)...)
}

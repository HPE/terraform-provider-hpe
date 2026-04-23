// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package networkfirewallrulegroup

import (
	"context"
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
			"import network firewall rule group resource",
			"provided import ID '"+req.ID+
				"' is invalid, expected format 'network_server_id:id'",
		)

		return
	}

	serverID, err := strconv.Atoi(parts[0])
	if err != nil {
		resp.Diagnostics.AddError(
			"import network firewall rule group resource",
			"provided network_server_id '"+parts[0]+"' is invalid (non-number)",
		)

		return
	}

	id, err := strconv.Atoi(parts[1])
	if err != nil {
		resp.Diagnostics.AddError(
			"import network firewall rule group resource",
			"provided id '"+parts[1]+"' is invalid (non-number)",
		)

		return
	}

	resp.Diagnostics.Append(
		resp.State.SetAttribute(
			ctx, path.Root("network_server_id"), serverID,
		)...,
	)
	resp.Diagnostics.Append(
		resp.State.SetAttribute(ctx, path.Root("id"), id)...,
	)
}

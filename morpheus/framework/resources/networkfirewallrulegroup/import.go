// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkfirewallrulegroup

import (
	"context"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	parts := strings.SplitN(req.ID, ".", 3)
	if len(parts) != 3 {
		resp.Diagnostics.AddError(
			"import network firewall rule group resource",
			"provided import ID '"+req.ID+
				"' is invalid, expected format 'network_integration_id.id.external_type'",
		)

		return
	}

	serverID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"import network firewall rule group resource",
			"provided network_integration_id '"+parts[0]+"' is invalid (non-number)",
		)

		return
	}

	id, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"import network firewall rule group resource",
			"provided id '"+parts[1]+"' is invalid (non-number)",
		)

		return
	}

	externalType := parts[2]
	if externalType == "" {
		resp.Diagnostics.AddError(
			"import network firewall rule group resource",
			"provided external_type is empty, expected a non-empty value (e.g. 'SecurityPolicy')",
		)

		return
	}

	resp.Diagnostics.Append(
		resp.State.SetAttribute(
			ctx, path.Root("network_integration_id"), types.Int64Value(serverID),
		)...,
	)
	resp.Diagnostics.Append(
		resp.State.SetAttribute(ctx, path.Root("id"), types.Int64Value(id))...,
	)
	resp.Diagnostics.Append(
		resp.State.SetAttribute(
			ctx, path.Root("external_type"), types.StringValue(externalType),
		)...,
	)
}

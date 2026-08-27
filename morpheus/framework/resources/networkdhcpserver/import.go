// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package networkdhcpserver

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
	parts := strings.SplitN(req.ID, ".", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"import network dhcp server resource",
			"provided import ID '"+req.ID+
				"' is invalid, expected format 'network_integration_id.id'",
		)

		return
	}

	integrationID, err := strconv.Atoi(parts[0])
	if err != nil {
		resp.Diagnostics.AddError(
			"import network dhcp server resource",
			"provided network_integration_id '"+parts[0]+"' is invalid (non-number)",
		)

		return
	}

	id, err := strconv.Atoi(parts[1])
	if err != nil {
		resp.Diagnostics.AddError(
			"import network dhcp server resource",
			"provided id '"+parts[1]+"' is invalid (non-number)",
		)

		return
	}

	resp.Diagnostics.Append(
		resp.State.SetAttribute(ctx, path.Root("network_integration_id"), integrationID)...,
	)
	resp.Diagnostics.Append(
		resp.State.SetAttribute(ctx, path.Root("id"), id)...,
	)
}

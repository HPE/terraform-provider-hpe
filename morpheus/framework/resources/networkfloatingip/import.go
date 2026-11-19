// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkfloatingip

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// importIDs is the parsed import identifier. networkServerID and
// floatingIpPoolID are nil for the short form.
type importIDs struct {
	id               int64
	networkServerID  *int64
	floatingIpPoolID *int64
}

// parseImportID parses "<id>" or
// "<id>:<network_server_id>:<floating_ip_pool_id>". It is a pure function so the
// parsing is unit-testable without HTTP.
func parseImportID(raw string) (importIDs, error) {
	parts := strings.Split(raw, ":")

	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return importIDs{}, fmt.Errorf("could not parse floating IP ID %q as integer: %w", parts[0], err)
	}

	switch len(parts) {
	case 1:
		return importIDs{id: id}, nil
	case 3:
		networkServerID, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return importIDs{}, fmt.Errorf("could not parse network_server_id %q as integer: %w", parts[1], err)
		}

		floatingIpPoolID, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			return importIDs{}, fmt.Errorf("could not parse floating_ip_pool_id %q as integer: %w", parts[2], err)
		}

		return importIDs{id: id, networkServerID: &networkServerID, floatingIpPoolID: &floatingIpPoolID}, nil
	default:
		return importIDs{}, fmt.Errorf(
			"expected \"<id>\" or \"<id>:<network_server_id>:<floating_ip_pool_id>\", got %q", raw,
		)
	}
}

// ImportState accepts either the short or long form (see parseImportID). The
// short form leaves the two inputs null, which makes the next plan show a
// replacement; a warning is added so the user is not surprised.
func (r *networkFloatingIpResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	ids, err := parseImportID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Invalid import ID", err.Error())

		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), ids.id)...)
	resp.Diagnostics.Append(
		resp.State.SetAttribute(ctx, path.Root("network_server_id"), ids.networkServerID)...,
	)
	resp.Diagnostics.Append(
		resp.State.SetAttribute(ctx, path.Root("floating_ip_pool_id"), ids.floatingIpPoolID)...,
	)

	if ids.networkServerID == nil {
		resp.Diagnostics.AddWarning(
			"Floating IP imported without its source",
			"network_server_id and floating_ip_pool_id are not returned by the API and were not "+
				"supplied in the import ID, so they are null in state. The next plan will show a "+
				"replacement. Import as \"<id>:<network_server_id>:<floating_ip_pool_id>\" to avoid this.",
		)
	}
}

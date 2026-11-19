// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package host

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ImportState implements resource.ResourceWithImportState. The host is imported
// by its numeric ID, optionally with the sha256 digests of its write-only
// values ("<id>#ssh_password_wo=<hex>"). Create-only inputs that the API does
// not return — the typed/free-form config blocks, network_interfaces and
// volumes — are not populated and plan a replacement if configured, so an
// imported host's configuration should leave them out; delete_options may be
// set and ssh_password_wo is updated in place (see the docs).
func (r *Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// "<id>" or "<id>#ssh_password_wo=<sha256 hex>": see wodigest.
	rawID, ok := wo.Import(ctx, req, resp)
	if !ok {
		return
	}

	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"invalid import ID",
			fmt.Sprintf("expected numeric host ID, got %q: %s", rawID, err.Error()),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkdomain

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// Import ID format: "<id>", optionally followed by
	// "#domain_password_wo=<sha256 hex>" (see wodigest).
	rawID, ok := wo.Import(ctx, req, resp)
	if !ok {
		return
	}

	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"import network_domain resource",
			fmt.Sprintf("invalid network domain ID %q: %s", rawID, err),
		)

		return
	}

	resp.Diagnostics.Append(
		resp.State.SetAttribute(ctx, path.Root("id"), id)...,
	)
}

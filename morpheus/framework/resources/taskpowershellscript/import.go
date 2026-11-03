// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package taskpowershellscript

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// Import ID format: "<id>", optionally followed by
	// "#remote_target_password_wo=<sha256 hex>" (see wodigest).
	kind.Import(ctx, req, resp)
}

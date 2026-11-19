// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkserver

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *networkServerResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// "<id>" or "<id>#service_password_wo=<hex>&service_token_wo=<hex>" (either
	// digest optional): see wodigest.
	rawID, ok := wo.Import(ctx, req, resp)
	if !ok {
		return
	}

	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse ID %q as integer: %s", rawID, err))

		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

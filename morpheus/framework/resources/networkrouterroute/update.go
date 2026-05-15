// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterroute

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

const updateOperation = "update network router route resource"

// Update deletes and recreates the route because the API does not support
// an update endpoint for network router routes.
func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	resp.Diagnostics.Append(
		diag.NewErrorDiagnostic(updateOperation, "Update is not supported for this resource type."),
	)
}

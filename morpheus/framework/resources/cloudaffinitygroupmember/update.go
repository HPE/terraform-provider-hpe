// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cloudaffinitygroupmember

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Update exists only to satisfy the interface.
//
// Every attribute forces replacement, so Terraform never calls this: a
// membership cannot be changed, only created or destroyed.
func (r *Resource) Update(
	_ context.Context,
	_ resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	resp.Diagnostics.AddError(
		"update not supported",
		"Every attribute of this resource forces replacement, so it should never "+
			"be updated in place. Please report this as a bug.",
	)
}

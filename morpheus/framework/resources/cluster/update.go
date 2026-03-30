// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cluster

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

const (
	updateOperation = "update cluster resource"
)

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	resp.Diagnostics.AddError(updateOperation, "Cluster update not implemented yet")
}

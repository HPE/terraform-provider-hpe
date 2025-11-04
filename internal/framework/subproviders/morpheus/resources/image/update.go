// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

//go:build experimental

package image

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Update implements resource.Resource.
func (r *Resource) Update(context.Context, resource.UpdateRequest, *resource.UpdateResponse) {
	panic("unimplemented")
}

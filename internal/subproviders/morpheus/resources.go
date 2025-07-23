// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

//go:build !experimental

// This file is used to include the Morpheus subprovider resources in the
// release build. It is not used in the experimental build. It is used to
// include only the stable resources in the release build. When building the
// experimental version, use the `-tags experimental` flag to exclude this
// file.

// When resources are ready for production use, they should be moved to this
// file

package morpheus

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/resources/group"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/resources/user"
)

func (s SubProvider) GetResources(
	_ context.Context,
) []func() resource.Resource {
	resources := []func() resource.Resource{
		group.NewResource,
		user.NewResource,
	}

	return resources
}

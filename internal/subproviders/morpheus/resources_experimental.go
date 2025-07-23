// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

//go:build experimental

// This file is used to include experimental resources in the Morpheus
// subprovider. It is not included in the release build. It is used to test new
// resources before they are included in the release build. It is not intended
// for production use and may contain unstable or incomplete features.

// When building the provider, use the `-tags experimental` flag to include
// this file.

// When resources are ready for production use, they should be moved to the
// `resources.go` file.

package morpheus

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/resources/group"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/resources/role"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/resources/user"
)

func (s SubProvider) GetResources(
	_ context.Context,
) []func() resource.Resource {
	resources := []func() resource.Resource{
		group.NewResource,
		user.NewResource,
		role.NewResource,
	}

	return resources
}

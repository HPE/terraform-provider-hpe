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

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/resources/cloud"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/resources/datastore"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/resources/group"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/resources/instance"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/resources/network"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/resources/policy"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/resources/role"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/resources/serviceplan"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/resources/user"
)

func (s SubProvider) GetResources(
	_ context.Context,
) []func() resource.Resource {
	resources := []func() resource.Resource{
		cloud.NewResource,
		datastore.NewResource,
		group.NewResource,
		network.NewResource,
		user.NewResource,
		role.NewResource,
		serviceplan.NewResource,
		instance.NewResource,
		policy.NewResource,
	}

	return resources
}

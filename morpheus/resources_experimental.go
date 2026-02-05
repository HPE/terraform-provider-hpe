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

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/cloud"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/datastore"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/group"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/image"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/instance"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/network"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/policy"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/role"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/serviceplan"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/task"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/user"
)

func (s SubProvider) GetResources(
	_ context.Context,
) []func() resource.Resource {
	resources := []func() resource.Resource{
		cloud.NewResource,
		datastore.NewResource,
		group.NewResource,
		image.NewResource,
		instance.NewResource,
		network.NewResource,
		policy.NewResource,
		role.NewResource,
		serviceplan.NewResource,
		task.NewResource,
		user.NewResource,
	}

	return resources
}

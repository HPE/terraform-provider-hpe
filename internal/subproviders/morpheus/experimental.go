// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

//go:build experimental

// This file is used to include experimental features in the Morpheus subprovider.
// It is not included in the release build.
// It is used to test new features before they are included in the release build.
// It is not intended for production use and may contain unstable or incomplete features.

// When building the provider, use the `-tags experimental` flag to include this file.

// When datasources or resources are ready for production use, they should be moved to the `release.go` file.
package morpheus

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/datasources/cloud"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/datasources/environment"
	dsgroup "github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/datasources/group"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/datasources/instancetypelayout"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/datasources/network"
	dsrole "github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/datasources/role"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/datasources/rolepermissions"
	dsserviceplan "github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/datasources/serviceplan"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/resources/group"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/resources/role"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/resources/user"
)

func (SubProvider) GetDataSources(
	_ context.Context,
) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		cloud.NewDataSource,
		environment.NewDataSource,
		dsgroup.NewDataSource,
		instancetypelayout.NewDataSource,
		network.NewDataSource,
		dsrole.NewDataSource,
		rolepermissions.NewDataSource,
		dsserviceplan.NewDataSource,
	}
}

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

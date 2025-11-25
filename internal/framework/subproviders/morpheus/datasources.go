// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

//go:build !experimental

// This file is used to include the Morpheus subprovider datasources in the
// release build. It is not used in the experimental build. It is used to
// include only the stable datasources in the release build. When building the
// experimental version, use the `-tags experimental` flag to exclude this
// file.

// When datasources are ready for production use, they should be moved to this
// file

package morpheus

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/datasources/cloud"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/datasources/datastore"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/datasources/environment"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/datasources/group"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/datasources/image"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/datasources/instancetypelayout"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/datasources/network"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/datasources/role"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/datasources/serviceplan"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/datasources/user"
)

func (SubProvider) GetDataSources(
	_ context.Context,
) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		cloud.NewDataSource,
		datastore.NewDataSource,
		environment.NewDataSource,
		group.NewDataSource,
		instancetypelayout.NewDataSource,
		image.NewDataSource,
		network.NewDataSource,
		role.NewDataSource,
		serviceplan.NewDataSource,
		user.NewDataSource,
	}
}

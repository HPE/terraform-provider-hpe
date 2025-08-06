// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

//go:build experimental

// This file is used to include experimental datasources in the Morpheus
// subprovider. It is not included in the release build. It is used to test new
// datasources before they are included in the release build. It is not
// intended for production use and may contain unstable or incomplete features.

// When building the provider, use the `-tags experimental` flag to include
// this file.

// When datasources are ready for production use, they should be moved to the
// `datasources.go` file.

package morpheus

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/datasources/cloud"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/datasources/environment"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/datasources/group"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/datasources/instancetypelayout"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/datasources/network"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/datasources/role"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/datasources/rolepermissions"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/datasources/serviceplan"
)

func (SubProvider) GetDataSources(
	_ context.Context,
) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		cloud.NewDataSource,
		environment.NewDataSource,
		group.NewDataSource,
		instancetypelayout.NewDataSource,
		network.NewDataSource,
		role.NewDataSource,
		rolepermissions.NewDataSource,
		serviceplan.NewDataSource,
	}
}

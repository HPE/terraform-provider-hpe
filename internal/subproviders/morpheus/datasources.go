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

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/datasources/cloud"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/datasources/environment"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/datasources/group"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/datasources/instancetypelayout"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/datasources/network"
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
	}
}

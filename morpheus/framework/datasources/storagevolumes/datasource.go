// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storagevolumes

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
)

const (
	summary = "read storage volumes data source"
	// listMax bounds the number of records fetched from the API in one call.
	listMax = 250
)

// Ensure the implementation satisfies the expected interfaces.
var _ datasource.DataSource = &DataSource{}

// NewDataSource is a helper function to simplify the provider implementation.
func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

// DataSource is the data source implementation.
type DataSource struct {
	configure.DataSourceWithMorpheusConfigure
	datasource.DataSource
}

// Metadata returns the data source type name.
func (d *DataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + "storage_volumes"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = StorageVolumesDataSourceSchema(ctx)
	resp.Schema.Description = "Retrieves a list of Morpheus storage volumes, optionally " +
		"filtered using one or more filter blocks."
	resp.Schema.MarkdownDescription = "Retrieves a list of Morpheus storage volumes, optionally " +
		"filtered using one or more filter blocks."
}

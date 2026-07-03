// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storagevolume

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
)

const (
	summary                     = "read storage volume data source"
	ErrorNoValidSearchTerms     = `no valid search terms - an id or name is required`
	ErrorNoStorageVolumeFound   = `no storage volume found`
	ErrorMultipleStorageVolumes = `multiple storage volumes were returned`
	// listMax bounds the number of records fetched when resolving a volume by
	// name via the list endpoint (mirrors the storage_volumes data source).
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
	resp.TypeName = req.ProviderTypeName + "_" + "storage_volume"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = StorageVolumeDataSourceSchema(ctx)
	resp.Schema.Description = "Retrieves information about a single Morpheus storage volume by id or name."
	resp.Schema.MarkdownDescription = "Retrieves information about a single Morpheus storage volume by id or name."
}

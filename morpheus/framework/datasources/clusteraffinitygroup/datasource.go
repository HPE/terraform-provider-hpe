// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package clusteraffinitygroup implements a data source for cluster_affinity_group
package clusteraffinitygroup

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
)

const (
	summary                            = "read cluster affinity group data source"
	ErrorNoValidSearchTerms            = `no valid search terms - an id or name is required`
	ErrorNoClusterAffinityGroupFound   = `no cluster affinity group found`
	ErrorMultipleClusterAffinityGroups = `multiple cluster affinity groups were returned`
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
	resp.TypeName = req.ProviderTypeName + "_" + "cluster_affinity_group"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = ClusterAffinityGroupDataSourceSchema(ctx)
}

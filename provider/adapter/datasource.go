// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package adapter

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
)

// DataSourceAdapter wraps a Terraform Plugin Framework data source and adapts
// it to work as a child data source within a parent provider architecture. It
// allows child provider data sources to be properly namespaced and integrated
// into a parent provider.
//
// The adapter implements all optional data source interfaces
// (DataSourceWithConfigure, DataSourceWithConfigValidators, and
// DataSourceWithValidateConfig) by delegating to the wrapped data source if it
// implements them. For optional interfaces not implemented by the wrapped data
// source, the adapter returns nil or no-op responses.
//
// Key transformations performed by the adapter:
//
// 1. Metadata: Transforms the data source TypeName by prepending the child
// provider's TypeName. For example, a "network" data source from a "morpheus"
// provider becomes "morpheus_network", and when used in a parent "hpe"
// provider becomes "hpe_morpheus_network".
//
// 2. Configure: Extracts the child provider's configuration data from the
// parent provider's ConfigureRequest.ProviderData map, ensuring the data
// source receives only its own provider's data.
//
// DataSourceAdapter has fewer optional interfaces than ResourceAdapter because
// data sources do not support state mutations (create, update, delete) and
// thus don't need interfaces like ResourceWithImportState,
// ResourceWithModifyPlan, ResourceWithMoveState, or ResourceWithUpgradeState.
type DataSourceAdapter struct {
	in       datasource.DataSource
	provider provider.Provider

	withConfigure        datasource.DataSourceWithConfigure
	withConfigValidators datasource.DataSourceWithConfigValidators
	withValidateConfig   datasource.DataSourceWithValidateConfig
}

var _ datasource.DataSource = &DataSourceAdapter{}
var _ datasource.DataSourceWithConfigure = &DataSourceAdapter{}
var _ datasource.DataSourceWithConfigValidators = &DataSourceAdapter{}
var _ datasource.DataSourceWithValidateConfig = &DataSourceAdapter{}

func NewDataSourceAdapter(in datasource.DataSource, p provider.Provider) *DataSourceAdapter {
	d := &DataSourceAdapter{in: in, provider: p}

	d.withConfigure, _ = in.(datasource.DataSourceWithConfigure)
	d.withConfigValidators, _ = in.(datasource.DataSourceWithConfigValidators)
	d.withValidateConfig, _ = in.(datasource.DataSourceWithValidateConfig)

	return d
}

func NewAdaptedDataSource(in datasource.DataSource, p provider.Provider) datasource.DataSource {
	return NewDataSourceAdapter(in, p)
}

// Metadata is the only method implementation that varies from `in`
// We use the Provider Adapter's name to the Metadata request.
// This will transform the data source name from e.g.:
// datasource -> {child_provider}_datasource
// When a parent provider is introduced, the data source name will
// then be registered as e.g.:
// {parent_provider}_{child_provider}_datasource
func (d *DataSourceAdapter) Metadata(
	ctx context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	// Get the provider's name from its metadata
	providerMetaResp := &provider.MetadataResponse{}
	d.provider.Metadata(ctx, provider.MetadataRequest{}, providerMetaResp)

	req.ProviderTypeName = req.ProviderTypeName + "_" + providerMetaResp.TypeName
	d.in.Metadata(ctx, req, resp)
}

func (d *DataSourceAdapter) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	d.in.Schema(ctx, req, resp)
}

func (d *DataSourceAdapter) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	d.in.Read(ctx, req, resp)
}

func (d *DataSourceAdapter) Configure(
	ctx context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	if d.withConfigure == nil {
		return
	}

	// Extract child provider configure data for ConfigureRequest.ProviderData
	if providerData, ok := req.ProviderData.(map[string]any); ok {
		metaResp := &provider.MetadataResponse{}
		d.provider.Metadata(ctx, provider.MetadataRequest{}, metaResp)

		req.ProviderData = providerData[metaResp.TypeName]
	}

	d.withConfigure.Configure(ctx, req, resp)
}

func (d *DataSourceAdapter) ConfigValidators(
	ctx context.Context,
) []datasource.ConfigValidator {
	if d.withConfigValidators == nil {
		return nil
	}

	return d.withConfigValidators.ConfigValidators(ctx)
}

func (d *DataSourceAdapter) ValidateConfig(
	ctx context.Context,
	req datasource.ValidateConfigRequest,
	resp *datasource.ValidateConfigResponse,
) {
	if d.withValidateConfig == nil {
		return
	}

	d.withValidateConfig.ValidateConfig(ctx, req, resp)
}

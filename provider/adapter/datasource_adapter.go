package adapter

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
)

type AdapterDataSource struct {
	in       datasource.DataSource
	provider provider.Provider

	withConfigure        datasource.DataSourceWithConfigure
	withConfigValidators datasource.DataSourceWithConfigValidators
	withValidateConfig   datasource.DataSourceWithValidateConfig
}

var _ datasource.DataSource = &AdapterDataSource{}
var _ datasource.DataSourceWithConfigure = &AdapterDataSource{}
var _ datasource.DataSourceWithConfigValidators = &AdapterDataSource{}
var _ datasource.DataSourceWithValidateConfig = &AdapterDataSource{}

func NewAdapterDataSource(in datasource.DataSource, p provider.Provider) *AdapterDataSource {
	d := &AdapterDataSource{in: in, provider: p}

	d.withConfigure, _ = in.(datasource.DataSourceWithConfigure)
	d.withConfigValidators, _ = in.(datasource.DataSourceWithConfigValidators)
	d.withValidateConfig, _ = in.(datasource.DataSourceWithValidateConfig)

	return d
}

func NewAdaptedDataSource(in datasource.DataSource, p provider.Provider) datasource.DataSource {
	return NewAdapterDataSource(in, p)
}

// Metadata is the only method implementation that varies from `in`
// We use the Provider Adapter's name to the Metadata request.
// This will transform the data source name from e.g.:
// datasource -> {child_provider}_datasource
// When a parent provier is introduced, the data source name will
// then be registered as e.g.:
// {parent_provider}_{child_provider}_datasource
func (d *AdapterDataSource) Metadata(
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

func (d *AdapterDataSource) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	d.in.Schema(ctx, req, resp)
}

func (d *AdapterDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	d.in.Read(ctx, req, resp)
}

func (d *AdapterDataSource) Configure(
	ctx context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	if d.withConfigure == nil {
		return
	}

	d.withConfigure.Configure(ctx, req, resp)
}

func (d *AdapterDataSource) ConfigValidators(
	ctx context.Context,
) []datasource.ConfigValidator {
	if d.withConfigValidators == nil {
		return nil
	}

	return d.withConfigValidators.ConfigValidators(ctx)
}

func (d *AdapterDataSource) ValidateConfig(
	ctx context.Context,
	req datasource.ValidateConfigRequest,
	resp *datasource.ValidateConfigResponse,
) {
	if d.withValidateConfig == nil {
		return
	}

	d.withValidateConfig.ValidateConfig(ctx, req, resp)
}

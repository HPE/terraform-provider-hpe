// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package adapter

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/statestore"
)

// ProviderAdapter wraps a Terraform Plugin Framework provider and adapts it to
// work as a child provider within a parent provider architecture. It allows
// child providers to be embedded as nested blocks in the parent provider's
// configuration.
//
// The adapter implements the standard provider.Provider interface and multiple
// optional provider interfaces, delegating method calls to the wrapped
// provider while optionally supporting advanced features like actions, config
// validators, functions, ephemeral resources, list resources, meta schema,
// state stores, and config validation if the wrapped provider implements them.
//
// The adapted provider's schema is transformed so that the child provider's
// attributes appear as a single nested block (identified by the child's
// TypeName), allowing configuration like:
//
//	provider "provider" {
//	  child_provider {
//	   first_attribute  = ""
//	   second_attribute = ""
//	  }
//	}
//
// Resources and data sources from the child provider are wrapped with
// AdaptedResource and AdaptedDataSource respectively to ensure they work
// correctly within the parent provider context.
type ProviderAdapter struct {
	in provider.Provider

	withActions          provider.ProviderWithActions
	withConfigValidators provider.ProviderWithConfigValidators
	withEphemeral        provider.ProviderWithEphemeralResources
	withFunctions        provider.ProviderWithFunctions
	withListResources    provider.ProviderWithListResources
	withMetaSchema       provider.ProviderWithMetaSchema
	withStateStores      provider.ProviderWithStateStores
	withValidateConfig   provider.ProviderWithValidateConfig
}

var (
	_ provider.Provider                       = &ProviderAdapter{}
	_ provider.ProviderWithActions            = &ProviderAdapter{}
	_ provider.ProviderWithConfigValidators   = &ProviderAdapter{}
	_ provider.ProviderWithEphemeralResources = &ProviderAdapter{}
	_ provider.ProviderWithFunctions          = &ProviderAdapter{}
	_ provider.ProviderWithListResources      = &ProviderAdapter{}
	_ provider.ProviderWithMetaSchema         = &ProviderAdapter{}
	_ provider.ProviderWithStateStores        = &ProviderAdapter{}
	_ provider.ProviderWithValidateConfig     = &ProviderAdapter{}
)

// Constructs a new Provider Adapter from any standard Terraform Plugin Framework Provider.
func NewProviderAdapter(in provider.Provider) *ProviderAdapter {
	p := &ProviderAdapter{in: in}

	p.withActions, _ = in.(provider.ProviderWithActions)
	p.withConfigValidators, _ = in.(provider.ProviderWithConfigValidators)
	p.withEphemeral, _ = in.(provider.ProviderWithEphemeralResources)
	p.withFunctions, _ = in.(provider.ProviderWithFunctions)
	p.withListResources, _ = in.(provider.ProviderWithListResources)
	p.withMetaSchema, _ = in.(provider.ProviderWithMetaSchema)
	p.withStateStores, _ = in.(provider.ProviderWithStateStores)
	p.withValidateConfig, _ = in.(provider.ProviderWithValidateConfig)

	return p
}

// The constructor that should be used for passing the child into the parent Provider.
func NewAdaptedProvider(in provider.Provider) provider.Provider {
	return NewProviderAdapter(in)
}

// Maintains metadata from `in`.
// Effectively, the adapted Provider will maintain the metadata of its input.
// We want to be able to pass the child Provider's name to the parent.
// Typically, Provider versions are set to "dev" in source code, and tagged at build time.
// The parent Provider should set its own Version string.
func (p *ProviderAdapter) Metadata(
	ctx context.Context,
	_ provider.MetadataRequest,
	resp *provider.MetadataResponse,
) {
	inMetaResp := &provider.MetadataResponse{}
	p.in.Metadata(ctx, provider.MetadataRequest{}, inMetaResp)
	resp.TypeName = inMetaResp.TypeName
	resp.Version = inMetaResp.Version
}

// Transforms the schema of `in` into a ListNestedBlock.
// ListNestedBlock is needed to be compatible with SDKV2 Providers if
// muxing the parent provider with an SDKV2 Provider using terraform-plugin-mux.
func (p *ProviderAdapter) Schema(
	ctx context.Context,
	_ provider.SchemaRequest,
	resp *provider.SchemaResponse,
) {
	// Fetch metadata from Adapter's `in`
	inMetaResp := &provider.MetadataResponse{}
	p.in.Metadata(ctx, provider.MetadataRequest{}, inMetaResp)

	inSchemaResp := &provider.SchemaResponse{}

	p.in.Schema(ctx, provider.SchemaRequest{}, inSchemaResp)

	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{
			inMetaResp.TypeName: schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: inSchemaResp.Schema.Attributes,
				},
				Validators: []validator.List{
					listvalidator.SizeBetween(0, 1),
				},
			},
		},
		Description:         inSchemaResp.Schema.GetDescription(),
		DeprecationMessage:  inSchemaResp.Schema.GetDeprecationMessage(),
		MarkdownDescription: inSchemaResp.Schema.GetMarkdownDescription(),
	}
}

func (p *ProviderAdapter) Configure(
	ctx context.Context,
	req provider.ConfigureRequest,
	resp *provider.ConfigureResponse,
) {
	p.in.Configure(ctx, req, resp)
}

func (p *ProviderAdapter) Resources(
	ctx context.Context,
) []func() resource.Resource {
	var adaptedResources []func() resource.Resource
	for _, f := range p.in.Resources(ctx) {
		adaptedResources = append(
			adaptedResources,
			func() resource.Resource {
				return NewAdaptedResource(f(), p)
			},
		)
	}

	return adaptedResources
}

func (p *ProviderAdapter) DataSources(
	ctx context.Context,
) []func() datasource.DataSource {
	var adaptedDataSources []func() datasource.DataSource
	for _, f := range p.in.DataSources(ctx) {
		adaptedDataSources = append(
			adaptedDataSources,
			func() datasource.DataSource {
				return NewAdaptedDataSource(f(), p)
			},
		)
	}

	return adaptedDataSources
}

// Additional functionality beyond standard provider.Provider interface
func (p *ProviderAdapter) Actions(
	ctx context.Context,
) []func() action.Action {
	if p.withActions == nil {
		return nil
	}

	return p.withActions.Actions(ctx)
}

func (p *ProviderAdapter) ConfigValidators(
	ctx context.Context,
) []provider.ConfigValidator {
	if p.withConfigValidators == nil {
		return nil
	}

	return p.withConfigValidators.ConfigValidators(ctx)
}

func (p *ProviderAdapter) Functions(
	ctx context.Context,
) []func() function.Function {
	if p.withFunctions == nil {
		return nil
	}

	return p.withFunctions.Functions(ctx)
}

func (p *ProviderAdapter) EphemeralResources(
	ctx context.Context,
) []func() ephemeral.EphemeralResource {
	if p.withEphemeral == nil {
		return nil
	}

	return p.withEphemeral.EphemeralResources(ctx)
}

func (p *ProviderAdapter) MetaSchema(
	ctx context.Context,
	req provider.MetaSchemaRequest,
	resp *provider.MetaSchemaResponse,
) {
	if p.withMetaSchema == nil {
		return
	}

	p.withMetaSchema.MetaSchema(ctx, req, resp)
}

func (p *ProviderAdapter) ListResources(
	ctx context.Context,
) []func() list.ListResource {
	if p.withListResources == nil {
		return nil
	}

	return p.withListResources.ListResources(ctx)
}

func (p *ProviderAdapter) StateStores(
	ctx context.Context,
) []func() statestore.StateStore {
	if p.withStateStores == nil {
		return nil
	}

	return p.withStateStores.StateStores(ctx)
}

func (p *ProviderAdapter) ValidateConfig(
	ctx context.Context,
	req provider.ValidateConfigRequest,
	resp *provider.ValidateConfigResponse,
) {
	if p.withValidateConfig == nil {
		return
	}

	p.withValidateConfig.ValidateConfig(ctx, req, resp)
}

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

type AdapterProvider struct {
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

var _ provider.Provider = &AdapterProvider{}
var _ provider.ProviderWithActions = &AdapterProvider{}
var _ provider.ProviderWithConfigValidators = &AdapterProvider{}
var _ provider.ProviderWithEphemeralResources = &AdapterProvider{}
var _ provider.ProviderWithFunctions = &AdapterProvider{}
var _ provider.ProviderWithListResources = &AdapterProvider{}
var _ provider.ProviderWithMetaSchema = &AdapterProvider{}
var _ provider.ProviderWithStateStores = &AdapterProvider{}
var _ provider.ProviderWithValidateConfig = &AdapterProvider{}

// Constructs a new Adapter Provider from any standard Terraform Plugin Framework Provider
func NewAdapterProvider(in provider.Provider) *AdapterProvider {
	p := &AdapterProvider{in: in}

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

// Helper constructor to construct a list of Adapted providers
func NewAdaptedChildProviders(in ...provider.Provider) []provider.Provider {
	var adaptedProviders []provider.Provider
	for _, childProvider := range in {
		adaptedProviders = append(adaptedProviders, NewAdapterProvider(childProvider))
	}
	return adaptedProviders
}

// TODO: allow for setting optional or required
// TODO: Figure this out, if it's worth having Optional or Required as settable
// type OptionalRequired = int
//
// const (
// 	Optional OptionalRequired = iota
// 	Required OptionalRequired = 1
// )

// Maintains metadata from `in`.
// Effectively, the adapted Provider will maintain the metadata of its input.
// We want to be able to pass the child Provider's name to the parent.
// Typically, Provider versions are set to "dev" in source code, and tagged at build time.
// The parent Provider should set its own Version string.
func (p *AdapterProvider) Metadata(
	ctx context.Context,
	_ provider.MetadataRequest,
	resp *provider.MetadataResponse,
) {
	inMetaResp := &provider.MetadataResponse{}
	p.in.Metadata(ctx, provider.MetadataRequest{}, inMetaResp)
	resp.TypeName = inMetaResp.TypeName
	resp.Version = inMetaResp.Version
}

// Transforms the schema of p.in into a SingleNestedAttribute
// This will create a provider schema that look like:
//
//	child_provider {
//	  first_attribute  = ""
//	  second_attribute = ""
//	}

func (p *AdapterProvider) Schema(
	ctx context.Context,
	_ provider.SchemaRequest,
	resp *provider.SchemaResponse,
) {
	// Fetch metadata from Adapter's input
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

func (p *AdapterProvider) Configure(
	ctx context.Context,
	req provider.ConfigureRequest,
	resp *provider.ConfigureResponse,
) {
	p.in.Configure(ctx, req, resp)
}

func (p *AdapterProvider) Resources(
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

func (p *AdapterProvider) DataSources(
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

func (p *AdapterProvider) Actions(
	ctx context.Context,
) []func() action.Action {
	if p.withActions == nil {
		return nil
	}

	return p.withActions.Actions(ctx)
}

func (p *AdapterProvider) ConfigValidators(
	ctx context.Context,
) []provider.ConfigValidator {
	if p.withConfigValidators == nil {
		return nil
	}

	return p.withConfigValidators.ConfigValidators(ctx)
}

func (p *AdapterProvider) Functions(
	ctx context.Context,
) []func() function.Function {
	if p.withFunctions == nil {
		return nil
	}

	return p.withFunctions.Functions(ctx)
}

func (p *AdapterProvider) EphemeralResources(
	ctx context.Context,
) []func() ephemeral.EphemeralResource {
	if p.withEphemeral == nil {
		return nil
	}

	return p.withEphemeral.EphemeralResources(ctx)
}

func (p *AdapterProvider) MetaSchema(
	ctx context.Context,
	req provider.MetaSchemaRequest,
	resp *provider.MetaSchemaResponse,
) {
	if p.withMetaSchema == nil {
		return
	}

	p.withMetaSchema.MetaSchema(ctx, req, resp)
}

func (p *AdapterProvider) ListResources(
	ctx context.Context,
) []func() list.ListResource {
	if p.withListResources == nil {
		return nil
	}

	return p.withListResources.ListResources(ctx)
}

func (p *AdapterProvider) StateStores(
	ctx context.Context,
) []func() statestore.StateStore {
	if p.withStateStores == nil {
		return nil
	}

	return p.withStateStores.StateStores(ctx)
}

func (p *AdapterProvider) ValidateConfig(
	ctx context.Context,
	req provider.ValidateConfigRequest,
	resp *provider.ValidateConfigResponse,
) {
	if p.withValidateConfig == nil {
		return
	}

	p.withValidateConfig.ValidateConfig(ctx, req, resp)
}

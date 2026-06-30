// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package provider

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/HPE/terraform-provider-hpe/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/notify"

	version "github.com/hashicorp/go-version"
)

var _ provider.Provider = &HpeProvider{}
var _ provider.ProviderWithEphemeralResources = &HpeProvider{}

func New(
	version string,
	providers ...provider.Provider,
) func() provider.Provider {
	return func() provider.Provider {
		return &HpeProvider{
			version:        version,
			childProviders: providers,
		}
	}
}

type HpeProvider struct {
	version        string
	childProviders []provider.Provider
}

func (p *HpeProvider) Metadata(
	_ context.Context,
	_ provider.MetadataRequest,
	resp *provider.MetadataResponse,
) {
	resp.TypeName = "hpe"
	resp.Version = p.version
}

func (p *HpeProvider) Schema(
	ctx context.Context,
	_ provider.SchemaRequest,
	resp *provider.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Blocks: make(map[string]schema.Block),
	}

	for _, s := range p.childProviders {
		metaResp := &provider.MetadataResponse{}
		schemaResp := &provider.SchemaResponse{}

		s.Metadata(ctx, provider.MetadataRequest{}, metaResp)
		s.Schema(ctx, provider.SchemaRequest{}, schemaResp)

		blockName := metaResp.TypeName

		// Prevent a panic if one of the child providers was passed
		// without using the adapter layer, and surface the error as a diagnostic.
		block, ok := schemaResp.Schema.Blocks[blockName]
		if !ok || block == nil {
			resp.Diagnostics.AddError(errfmt.ChildProviderSchemaErr(metaResp.TypeName))

			return
		}

		resp.Schema.Blocks[blockName] = block
	}
}

func (p *HpeProvider) Configure(
	ctx context.Context,
	req provider.ConfigureRequest,
	resp *provider.ConfigureResponse,
) {
	var wg sync.WaitGroup
	// check version if not running on dev or test (accepance test builds)
	if p.version != "dev" && p.version != "test" && notify.IsEnabled() && notify.TryDial() == nil {
		// Do this in a separate goroutine while the rest of the Configure method runs.
		wg.Add(1)
		go func() {
			defer wg.Done()
			localVer, err := version.NewVersion(strings.ToLower(strings.TrimPrefix(p.version, "v")))
			if err != nil {
				// Installed provider version is messed up, so surface an Error Diagnostic.
				// Terraform should catch this first, though, as it doesn't allow malformed versions.
				resp.Diagnostics.Append(
					diag.NewErrorDiagnostic(
						"failed to convert local version string",
						err.Error(),
					),
				)

				return
			}

			retry := func() (*version.Version, error) {
				return notify.GetProviderVersion(notify.RegistryUrl)
			}
			// If we passed the TryDial check, we'll probably hit the API relatively quickly.
			remoteVer, err := backoff.Retry(
				ctx,
				retry,
				backoff.WithMaxElapsedTime(30*time.Second),
			)
			if err != nil {
				// Continue provider execution if this fails.
				resp.Diagnostics.Append(
					diag.NewWarningDiagnostic(
						"failed to fetch latest remote version",
						err.Error(),
					),
				)

				return
			}

			err = notify.CompareProviderVersion(localVer, remoteVer)
			if err != nil {
				// i.e. if localVer < remoteVer
				resp.Diagnostics.Append(diag.NewWarningDiagnostic("Outdated provider version", err.Error()))
			}
		}()

	}

	resourceData := map[string]any{}
	dataSourceData := map[string]any{}
	ephemeralResourceData := map[string]any{}

	// Parse the parent raw config value as a map of block names → tftypes.Value.
	// This is because we need to provide tftypes.Value for the 'Raw' part of
	// constructing the childProverConfigReq later in this function.
	parentAttrs := map[string]tftypes.Value{}
	if err := req.Config.Raw.As(&parentAttrs); err != nil {
		resp.Diagnostics.AddError("Failed to read provider config", err.Error())
		wg.Wait()

		return
	}

	for _, s := range p.childProviders {
		childMetaResp := &provider.MetadataResponse{}
		s.Metadata(ctx, provider.MetadataRequest{}, childMetaResp)

		// Since the "hpe" provider is using ListNestedBlock for its configs,
		// we need to pass the 0th ListNestedBlock to the child provider
		// so that it can parse its config as a flat map[string]Attribute.
		blockName := childMetaResp.TypeName
		blocks := req.Config.Schema.GetBlocks()
		block := blocks[blockName]
		fwAttrs := block.GetNestedObject().GetAttributes()

		schemaAttrs := make(map[string]schema.Attribute)
		// assert fwschema UnderlyingAttributes to schema Attribute.
		for k, v := range fwAttrs {
			if schemaAttr, ok := v.(schema.Attribute); ok {
				schemaAttrs[k] = schemaAttr
			}
		}

		// Extract the list tftypes.Value for this child's block.
		listVal, ok := parentAttrs[blockName]
		// The blocks are optional.
		if !ok || listVal.IsNull() || !listVal.IsKnown() {
			continue
		}

		// Unwrap the list to its elements.
		var elems []tftypes.Value
		if err := listVal.As(&elems); err != nil || len(elems) == 0 {
			if err != nil {
				// Only fail on error.
				resp.Diagnostics.AddError(
					"failed to configure hpe provider",
					err.Error(),
				)

				return
			}

			continue
		}

		childProviderConfigReq := provider.ConfigureRequest{
			TerraformVersion:   req.TerraformVersion,
			ClientCapabilities: req.ClientCapabilities,
			Config: tfsdk.Config{
				Schema: schema.Schema{
					Attributes: schemaAttrs,
				},
				Raw: elems[0], // flat object: {url, username, ...}
			},
		}

		childConfigResp := &provider.ConfigureResponse{}
		s.Configure(ctx, childProviderConfigReq, childConfigResp)

		resp.Diagnostics.Append(childConfigResp.Diagnostics...)

		if childConfigResp.ResourceData != nil {
			resourceData[childMetaResp.TypeName] = childConfigResp.ResourceData
		}

		if childConfigResp.DataSourceData != nil {
			dataSourceData[childMetaResp.TypeName] = childConfigResp.DataSourceData
		}

		if childConfigResp.EphemeralResourceData != nil {
			ephemeralResourceData[childMetaResp.TypeName] = childConfigResp.EphemeralResourceData
		}
	}

	resp.ResourceData = resourceData
	resp.DataSourceData = dataSourceData
	resp.EphemeralResourceData = ephemeralResourceData

	wg.Wait()
}

func (p *HpeProvider) Resources(
	ctx context.Context,
) []func() resource.Resource {
	var resources []func() resource.Resource
	for _, s := range p.childProviders {
		resources = append(resources, s.Resources(ctx)...)
	}

	return resources
}

func (p *HpeProvider) DataSources(
	ctx context.Context,
) []func() datasource.DataSource {
	var datasources []func() datasource.DataSource
	for _, s := range p.childProviders {
		datasources = append(datasources, s.DataSources(ctx)...)
	}

	return datasources
}

func (p *HpeProvider) EphemeralResources(
	ctx context.Context,
) []func() ephemeral.EphemeralResource {
	var ephemerals []func() ephemeral.EphemeralResource

	for _, s := range p.childProviders {
		if ep, ok := s.(provider.ProviderWithEphemeralResources); ok {
			ephemerals = append(ephemerals, ep.EphemeralResources(ctx)...)
		}
	}

	return ephemerals
}

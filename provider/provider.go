// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package provider

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/HPE/terraform-provider-hpe/provider/adapter"
	"github.com/HPE/terraform-provider-hpe/provider/subprovider"
	"github.com/HPE/terraform-provider-hpe/utils/notify"

	version "github.com/hashicorp/go-version"
)

var _ provider.Provider = &HpeProvider{}

func New(
	version string,
	b ...subprovider.SubProvider,
) func() provider.Provider {
	return func() provider.Provider {
		return &HpeProvider{
			version:      version,
			subproviders: b,
		}
	}
}

func New2(
	version string,
	providers ...provider.Provider,
) func() provider.Provider {
	return func() provider.Provider {
		return &HpeProvider2{
			version: version,
			childProviders: adapter.NewAdaptedChildProviders(
				providers...,
			),
		}
	}
}

type HpeProvider struct {
	version      string
	subproviders []subprovider.SubProvider
}

type HpeProvider2 struct {
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

func (p *HpeProvider2) Metadata(
	_ context.Context,
	_ provider.MetadataRequest,
	resp *provider.MetadataResponse,
) {
	resp.TypeName = "hpe"
	resp.Version = p.version
}

type AttrMap struct {
	name       string
	attributes map[string]schema.Attribute
}

func createListNestedBlock(attrmaps []AttrMap) map[string]schema.Block {
	blockmap := map[string]schema.Block{}
	for _, attrmap := range attrmaps {
		block := schema.ListNestedBlock{
			NestedObject: schema.NestedBlockObject{
				Attributes: attrmap.attributes,
			},
			Validators: []validator.List{
				listvalidator.SizeBetween(0, 1),
			},
		}
		blockmap[attrmap.name] = block
	}

	return blockmap
}

// func createListNestedAttributes(attrmaps []AttrMap) map[string]schema.Attribute {
// 	blockmap := map[string]schema.Block{}
// 	for _, attrmap := range attrmaps {
// 		block := schema.SingleNestedAttribute{
// 			NestedObject: schema.NestedBlockObject{
// 				Attributes: attrmap.attributes,
// 			},
// 			Validators: []validator.List{
// 				listvalidator.SizeBetween(0, 1),
// 			},
// 		}
// 		blockmap[attrmap.name] = block
// 	}
//
// 	return blockmap
// }

func (p *HpeProvider) Schema(
	ctx context.Context,
	_ provider.SchemaRequest,
	resp *provider.SchemaResponse,
) {
	var a []AttrMap
	for _, s := range p.subproviders {
		a = append(a, AttrMap{
			name:       s.GetName(ctx),
			attributes: s.GetSchema(ctx),
		})
	}

	blocks := createListNestedBlock(a)

	resp.Schema = schema.Schema{
		Blocks: blocks,
	}
}

func (p *HpeProvider2) Schema(
	ctx context.Context,
	_ provider.SchemaRequest,
	resp *provider.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		// Attributes: make(map[string]schema.Attribute),
		Blocks: make(map[string]schema.Block),
	}

	for _, s := range p.childProviders {
		metaResp := &provider.MetadataResponse{}
		schemaResp := &provider.SchemaResponse{}

		s.Metadata(ctx, provider.MetadataRequest{}, metaResp)
		s.Schema(ctx, provider.SchemaRequest{}, schemaResp)

		// resp.Schema.Attributes[metaResp.TypeName] = schemaResp.Schema.Attributes[metaResp.TypeName]
		resp.Schema.Blocks[metaResp.TypeName] = schemaResp.Schema.Blocks[metaResp.TypeName]
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

	f := func(ctx context.Context, c tfsdk.Config, name string) func(any) {
		return func(target any) {
			c.GetAttribute(ctx, path.Root(name), target)
		}
	}

	d := map[string]any{}
	for _, s := range p.subproviders {
		v, err := s.Configure(ctx, f(ctx, req.Config, s.GetName(ctx)))
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to configure "+s.GetName(ctx),
				err.Error(),
			)
		}
		d[s.GetName(ctx)] = v
	}

	resp.ResourceData = d
	resp.DataSourceData = d
	wg.Wait()
}

func (p *HpeProvider2) Configure(
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

	// Navigate the parent raw value as a map of block names → tftypes.Value.
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
		// so that it can parse its config as a flat map[string]Attribute
		blocks := req.Config.Schema.GetBlocks()
		block := blocks[childMetaResp.TypeName]
		fwAttrs := block.GetNestedObject().GetAttributes()

		schemaAttrs := make(map[string]schema.Attribute)
		// assert fwschema UnderlyingAttributes to schema Attribute
		for k, v := range fwAttrs {
			if schemaAttr, ok := v.(schema.Attribute); ok {
				schemaAttrs[k] = schemaAttr
			}
		}

		// Extract the list tftypes.Value for this child's block.
		listVal, ok := parentAttrs[childMetaResp.TypeName]
		if !ok || listVal.IsNull() || !listVal.IsKnown() {
			continue
		}

		// Unwrap the list to its elements.
		var elems []tftypes.Value
		if err := listVal.As(&elems); err != nil || len(elems) == 0 {
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
	}

	resp.ResourceData = resourceData
	resp.DataSourceData = dataSourceData

	wg.Wait()
}

func (p *HpeProvider) Resources(
	ctx context.Context,
) []func() resource.Resource {
	var resources []func() resource.Resource
	for _, s := range p.subproviders {
		resources = append(resources, s.GetResources(ctx)...)
	}

	return resources
}

func (p *HpeProvider2) Resources(
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
	for _, s := range p.subproviders {
		datasources = append(datasources, s.GetDataSources(ctx)...)
	}

	return datasources
}

func (p *HpeProvider2) DataSources(
	ctx context.Context,
) []func() datasource.DataSource {
	var datasources []func() datasource.DataSource
	for _, s := range p.childProviders {
		datasources = append(datasources, s.DataSources(ctx)...)
	}

	return datasources
}

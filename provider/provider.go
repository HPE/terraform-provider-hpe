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

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/notify"
	"github.com/HPE/terraform-provider-hpe/provider/subprovider"

	version "github.com/hashicorp/go-version"
)

var _ provider.Provider = &hpeProvider{}

func New(
	version string,
	b ...subprovider.SubProvider,
) func() provider.Provider {
	return func() provider.Provider {
		return &hpeProvider{
			version:      version,
			subproviders: b,
		}
	}
}

type hpeProvider struct {
	version      string
	subproviders []subprovider.SubProvider
}

func (p *hpeProvider) Metadata(
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

func (p *hpeProvider) Schema(
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

func (p *hpeProvider) Configure(
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

func (p *hpeProvider) Resources(
	ctx context.Context,
) []func() resource.Resource {
	var resources []func() resource.Resource
	for _, s := range p.subproviders {
		resources = append(resources, s.GetResources(ctx)...)
	}

	return resources
}

func (p *hpeProvider) DataSources(
	ctx context.Context,
) []func() datasource.DataSource {
	var datasources []func() datasource.DataSource
	for _, s := range p.subproviders {
		datasources = append(datasources, s.GetDataSources(ctx)...)
	}

	return datasources
}

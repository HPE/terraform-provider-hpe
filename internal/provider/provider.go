// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package provider

import (
	"context"

	"github.com/HPE/terraform-provider-hpe/subprovider"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

var _ provider.Provider = &HPEProvider{}

func New(
	version string,
	b ...subprovider.SubProvider,
) func() provider.Provider {
	return func() provider.Provider {
		return &HPEProvider{
			version:      version,
			subproviders: b,
		}
	}
}

type HPEProvider struct {
	version      string
	subproviders []subprovider.SubProvider
}

func (p *HPEProvider) Metadata(
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

func (p *HPEProvider) Schema(
	_ context.Context,
	_ provider.SchemaRequest,
	resp *provider.SchemaResponse,
) {
	var a []AttrMap
	for _, s := range p.subproviders {
		a = append(a, AttrMap{
			name:       s.GetName(),
			attributes: s.GetSchema(),
		})
	}

	blocks := createListNestedBlock(a)

	resp.Schema = schema.Schema{
		Blocks: blocks,
	}
}

func (p *HPEProvider) Configure(
	ctx context.Context,
	req provider.ConfigureRequest,
	_ *provider.ConfigureResponse,
) {
	f := func(ctx context.Context, c tfsdk.Config, name string) func(any) {
		return func(target any) {
			c.GetAttribute(ctx, path.Root(name), target)
		}
	}

	for _, s := range p.subproviders {
		s.Configure(f(ctx, req.Config, s.GetName()))
	}
}

func (p *HPEProvider) Resources(
	_ context.Context,
) []func() resource.Resource {
	var resources []func() resource.Resource
	for _, s := range p.subproviders {
		resources = append(resources, s.GetResources()...)
	}

	return resources
}

func (p *HPEProvider) DataSources(
	_ context.Context,
) []func() datasource.DataSource {
	var datasources []func() datasource.DataSource
	for _, s := range p.subproviders {
		datasources = append(datasources, s.GetDataSources()...)
	}

	return datasources
}

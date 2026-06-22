// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	frameworkschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	sdkv2schema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func FwToSdkv2SchemaViaProto(in frameworkschema.Schema) *sdkv2schema.Schema {
	ctx := context.Background()

	server := providerserver.NewProtocol5(staticSchemaProvider{schema: in})()
	resp, err := server.GetProviderSchema(ctx, &tfprotov5.GetProviderSchemaRequest{})
	if err != nil {
		panic(err)
	}

	for _, diag := range resp.Diagnostics {
		if diag != nil && diag.Severity == tfprotov5.DiagnosticSeverityError {
			panic(diag.Summary + ": " + diag.Detail)
		}
	}

	return tfprotov5SchemaToSdkv2(resp.Provider)
}

type staticSchemaProvider struct {
	schema frameworkschema.Schema
}

func (p staticSchemaProvider) Metadata(context.Context, provider.MetadataRequest, *provider.MetadataResponse) {
}

func (p staticSchemaProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = p.schema
}

func (p staticSchemaProvider) Configure(context.Context, provider.ConfigureRequest, *provider.ConfigureResponse) {
}

func (p staticSchemaProvider) DataSources(context.Context) []func() datasource.DataSource {
	return nil
}

func (p staticSchemaProvider) Resources(context.Context) []func() resource.Resource {
	return nil
}

func tfprotov5SchemaToSdkv2(in *tfprotov5.Schema) *sdkv2schema.Schema {
	return &sdkv2schema.Schema{
		Type:     sdkv2schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &sdkv2schema.Resource{
			Schema: tfprotov5BlockToSdkv2Map(in.Block),
		},
	}
}

func tfprotov5BlockToSdkv2Map(block *tfprotov5.SchemaBlock) map[string]*sdkv2schema.Schema {
	if block == nil {
		return nil
	}

	out := make(map[string]*sdkv2schema.Schema, len(block.Attributes)+len(block.BlockTypes))

	for _, attr := range block.Attributes {
		s := &sdkv2schema.Schema{
			Required:    attr.Required,
			Optional:    attr.Optional,
			Computed:    attr.Computed,
			Sensitive:   attr.Sensitive,
			WriteOnly:   attr.WriteOnly,
			Description: attr.Description,
			Deprecated:  attr.DeprecationMessage,
		}

		applyTFTypeToSdkSchema(context.Background(), s, attr.Type)
		out[attr.Name] = s
	}

	for _, nested := range block.BlockTypes {
		out[nested.TypeName] = tfprotov5NestedBlockToSdkv2(nested)
	}

	return out
}

func tfprotov5NestedBlockToSdkv2(block *tfprotov5.SchemaNestedBlock) *sdkv2schema.Schema {
	out := &sdkv2schema.Schema{
		Optional: true,
		MinItems: int(block.MinItems),
		MaxItems: int(block.MaxItems),
		Elem: &sdkv2schema.Resource{
			Schema: tfprotov5BlockToSdkv2Map(block.Block),
		},
	}

	switch block.Nesting {
	case tfprotov5.SchemaNestedBlockNestingModeSingle,
		tfprotov5.SchemaNestedBlockNestingModeGroup:
		out.Type = sdkv2schema.TypeList
		out.MaxItems = 1
	case tfprotov5.SchemaNestedBlockNestingModeList:
		out.Type = sdkv2schema.TypeList
	case tfprotov5.SchemaNestedBlockNestingModeSet:
		out.Type = sdkv2schema.TypeSet
	default:
		panic(fmt.Sprintf("unsupported tfprotov5 nested block mode %v", block.Nesting))
	}

	return out
}

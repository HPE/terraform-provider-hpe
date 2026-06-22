// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package provider

import (
	"context"
	"fmt"

	frameworkschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	sdkv2schema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// FwToSdkv2Schema converts a Terraform Plugin Framework provider schema into a
// lossy SDKv2 schema wrapper. Prefer FwToSdkv2SchemaMap when assigning to
// (*schema.Provider).Schema.
func FwToSdkv2Schema(in frameworkschema.Schema) *sdkv2schema.Schema {
	return &sdkv2schema.Schema{
		Type:     sdkv2schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &sdkv2schema.Resource{
			Schema: FwToSdkv2SchemaMap(in),
		},
	}
}

func FwToSdkv2SchemaMap(in frameworkschema.Schema) map[string]*sdkv2schema.Schema {
	return fwObjectToSdkv2SchemaMap(context.Background(), in.Attributes, in.Blocks)
}

func fwObjectToSdkv2SchemaMap(ctx context.Context, attrs map[string]frameworkschema.Attribute, blocks map[string]frameworkschema.Block) map[string]*sdkv2schema.Schema {
	out := make(map[string]*sdkv2schema.Schema, len(attrs)+len(blocks))

	for name, attr := range attrs {
		out[name] = fwAttributeToSdkv2(ctx, attr)
	}

	for name, block := range blocks {
		out[name] = fwBlockToSdkv2(ctx, block)
	}

	return out
}

func fwAttributeToSdkv2(ctx context.Context, attr frameworkschema.Attribute) *sdkv2schema.Schema {
	out := &sdkv2schema.Schema{
		Required:    attr.IsRequired(),
		Optional:    attr.IsOptional(),
		Computed:    attr.IsComputed(),
		Sensitive:   attr.IsSensitive(),
		WriteOnly:   attr.IsWriteOnly(),
		Description: attr.GetDescription(),
		Deprecated:  attr.GetDeprecationMessage(),
	}

	// if attr.GetMarkdownDescription() != "" {
	// 	out.Description = attr.GetMarkdownDescription()
	// }

	switch a := attr.(type) {
	case frameworkschema.StringAttribute:
		out.Type = sdkv2schema.TypeString
	case frameworkschema.BoolAttribute:
		out.Type = sdkv2schema.TypeBool
	case frameworkschema.Int32Attribute, frameworkschema.Int64Attribute:
		out.Type = sdkv2schema.TypeInt
	case frameworkschema.Float32Attribute, frameworkschema.Float64Attribute, frameworkschema.NumberAttribute:
		out.Type = sdkv2schema.TypeFloat
	case frameworkschema.ListAttribute:
		out.Type = sdkv2schema.TypeList
		out.Elem = sdkElemFromTFType(ctx, a.ElementType.TerraformType(ctx))
	case frameworkschema.SetAttribute:
		out.Type = sdkv2schema.TypeSet
		out.Elem = sdkElemFromTFType(ctx, a.ElementType.TerraformType(ctx))
	case frameworkschema.MapAttribute:
		out.Type = sdkv2schema.TypeMap
		out.Elem = sdkElemFromTFType(ctx, a.ElementType.TerraformType(ctx))
	case frameworkschema.SingleNestedAttribute:
		out.Type = sdkv2schema.TypeList
		out.MaxItems = 1
		out.Elem = &sdkv2schema.Resource{
			Schema: fwObjectToSdkv2SchemaMap(ctx, a.Attributes, nil),
		}
	case frameworkschema.ListNestedAttribute:
		out.Type = sdkv2schema.TypeList
		out.Elem = &sdkv2schema.Resource{
			Schema: fwObjectToSdkv2SchemaMap(ctx, a.NestedObject.Attributes, nil),
		}
	case frameworkschema.SetNestedAttribute:
		out.Type = sdkv2schema.TypeSet
		out.Elem = &sdkv2schema.Resource{
			Schema: fwObjectToSdkv2SchemaMap(ctx, a.NestedObject.Attributes, nil),
		}
	default:
		applyTFTypeToSdkSchema(ctx, out, attr.GetType().TerraformType(ctx))
	}

	return out
}

func fwBlockToSdkv2(ctx context.Context, block frameworkschema.Block) *sdkv2schema.Schema {
	out := &sdkv2schema.Schema{
		Optional:    true,
		Description: block.GetDescription(),
		Deprecated:  block.GetDeprecationMessage(),
	}

	if block.GetMarkdownDescription() != "" {
		out.Description = block.GetMarkdownDescription()
	}

	switch b := block.(type) {
	case frameworkschema.SingleNestedBlock:
		out.Type = sdkv2schema.TypeList
		out.MaxItems = 1
		out.Elem = &sdkv2schema.Resource{
			Schema: fwObjectToSdkv2SchemaMap(ctx, b.Attributes, b.Blocks),
		}
	case frameworkschema.ListNestedBlock:
		out.Type = sdkv2schema.TypeList
		out.Elem = &sdkv2schema.Resource{
			Schema: fwObjectToSdkv2SchemaMap(ctx, b.NestedObject.Attributes, b.NestedObject.Blocks),
		}
	case frameworkschema.SetNestedBlock:
		out.Type = sdkv2schema.TypeSet
		out.Elem = &sdkv2schema.Resource{
			Schema: fwObjectToSdkv2SchemaMap(ctx, b.NestedObject.Attributes, b.NestedObject.Blocks),
		}
	default:
		panic(fmt.Sprintf("unsupported framework block type %T", block))
	}

	return out
}

func sdkElemFromTFType(ctx context.Context, typ tftypes.Type) any {
	s := &sdkv2schema.Schema{}
	applyTFTypeToSdkSchema(ctx, s, typ)

	return s
}

func applyTFTypeToSdkSchema(ctx context.Context, out *sdkv2schema.Schema, typ tftypes.Type) {
	switch {
	case typ.Is(tftypes.String):
		out.Type = sdkv2schema.TypeString
	case typ.Is(tftypes.Bool):
		out.Type = sdkv2schema.TypeBool
	case typ.Is(tftypes.Number):
		out.Type = sdkv2schema.TypeFloat
	case typ.Is(tftypes.List{}):
		t := typ.(tftypes.List)
		out.Type = sdkv2schema.TypeList
		out.Elem = sdkElemFromTFType(ctx, t.ElementType)
	case typ.Is(tftypes.Set{}):
		t := typ.(tftypes.Set)
		out.Type = sdkv2schema.TypeSet
		out.Elem = sdkElemFromTFType(ctx, t.ElementType)
	case typ.Is(tftypes.Map{}):
		t := typ.(tftypes.Map)
		out.Type = sdkv2schema.TypeMap
		out.Elem = sdkElemFromTFType(ctx, t.ElementType)
	default:
		panic(fmt.Sprintf("unsupported Terraform type %T", typ))
	}
}

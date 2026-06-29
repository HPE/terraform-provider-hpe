// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package convert

import (
	"context"
	"fmt"

	fwschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	sdkv2schema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// FwToSdkv2Schema converts a Terraform Plugin Framework provider schema into a
// SDKv2 schema. The conversion is inherently lossy as not all Framework
// concepts can be represented equivalently in SDKV2.
// However, this conversion is good enough for our usecase of injecting the
// Framework Provider Schema into an SDKv2 provider.
func FwToSdkv2Schema(in fwschema.Schema) *sdkv2schema.Schema {
	return &sdkv2schema.Schema{
		Type:     sdkv2schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &sdkv2schema.Resource{
			Schema: FwToSdkv2SchemaMap(in),
		},
	}
}

func FwToSdkv2SchemaMap(in fwschema.Schema) map[string]*sdkv2schema.Schema {
	return fwObjectToSdkv2SchemaMap(context.Background(), in.Attributes, in.Blocks)
}

func fwObjectToSdkv2SchemaMap(ctx context.Context, attrs map[string]fwschema.Attribute, blocks map[string]fwschema.Block) map[string]*sdkv2schema.Schema {
	out := make(map[string]*sdkv2schema.Schema, len(attrs)+len(blocks))

	for name, attr := range attrs {
		out[name] = fwAttributeToSdkv2(ctx, attr)
	}

	for name, block := range blocks {
		out[name] = fwBlockToSdkv2(ctx, block)
	}

	return out
}

func fwAttributeToSdkv2(ctx context.Context, attr fwschema.Attribute) *sdkv2schema.Schema {
	out := &sdkv2schema.Schema{
		Required:    attr.IsRequired(),
		Optional:    attr.IsOptional(),
		Computed:    attr.IsComputed(),
		Sensitive:   attr.IsSensitive(),
		WriteOnly:   attr.IsWriteOnly(),
		Description: attr.GetDescription(),
		Deprecated:  attr.GetDeprecationMessage(),
	}

	switch a := attr.(type) {
	case fwschema.StringAttribute:
		out.Type = sdkv2schema.TypeString
	case fwschema.BoolAttribute:
		out.Type = sdkv2schema.TypeBool
	case fwschema.Int32Attribute, fwschema.Int64Attribute:
		out.Type = sdkv2schema.TypeInt
	case fwschema.Float32Attribute, fwschema.Float64Attribute, fwschema.NumberAttribute:
		out.Type = sdkv2schema.TypeFloat
	case fwschema.ListAttribute:
		out.Type = sdkv2schema.TypeList
		out.Elem = sdkElemFromTFType(ctx, a.ElementType.TerraformType(ctx))
	case fwschema.SetAttribute:
		out.Type = sdkv2schema.TypeSet
		out.Elem = sdkElemFromTFType(ctx, a.ElementType.TerraformType(ctx))
	case fwschema.MapAttribute:
		out.Type = sdkv2schema.TypeMap
		out.Elem = sdkElemFromTFType(ctx, a.ElementType.TerraformType(ctx))
	case fwschema.SingleNestedAttribute:
		out.Type = sdkv2schema.TypeList
		out.MaxItems = 1
		out.Elem = &sdkv2schema.Resource{
			Schema: fwObjectToSdkv2SchemaMap(ctx, a.Attributes, nil),
		}
	case fwschema.ListNestedAttribute:
		out.Type = sdkv2schema.TypeList
		out.Elem = &sdkv2schema.Resource{
			Schema: fwObjectToSdkv2SchemaMap(ctx, a.NestedObject.Attributes, nil),
		}
	case fwschema.SetNestedAttribute:
		out.Type = sdkv2schema.TypeSet
		out.Elem = &sdkv2schema.Resource{
			Schema: fwObjectToSdkv2SchemaMap(ctx, a.NestedObject.Attributes, nil),
		}
	default:
		applyTFTypeToSdkSchema(ctx, out, attr.GetType().TerraformType(ctx))
	}

	return out
}

func fwBlockToSdkv2(ctx context.Context, block fwschema.Block) *sdkv2schema.Schema {
	out := &sdkv2schema.Schema{
		Optional:    true,
		Description: block.GetDescription(),
		Deprecated:  block.GetDeprecationMessage(),
	}

	if block.GetMarkdownDescription() != "" {
		out.Description = block.GetMarkdownDescription()
	}

	switch b := block.(type) {
	case fwschema.SingleNestedBlock:
		out.Type = sdkv2schema.TypeList
		out.MaxItems = 1
		out.Elem = &sdkv2schema.Resource{
			Schema: fwObjectToSdkv2SchemaMap(ctx, b.Attributes, b.Blocks),
		}
	case fwschema.ListNestedBlock:
		out.Type = sdkv2schema.TypeList
		out.Elem = &sdkv2schema.Resource{
			Schema: fwObjectToSdkv2SchemaMap(ctx, b.NestedObject.Attributes, b.NestedObject.Blocks),
		}
	case fwschema.SetNestedBlock:
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

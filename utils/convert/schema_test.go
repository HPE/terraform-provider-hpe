// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package convert

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	fwschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	sdkv2schema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestFwToSdkv2Schema(t *testing.T) {
	t.Parallel()

	in := fwschema.Schema{
		Attributes: map[string]fwschema.Attribute{
			"url": fwschema.StringAttribute{
				Required:    true,
				Description: "The URL",
			},
		},
	}

	got := FwToSdkv2Schema(in)

	if got.Type != sdkv2schema.TypeList {
		t.Errorf("expected TypeList, got %v", got.Type)
	}

	if !got.Optional {
		t.Error("expected Optional to be true")
	}

	if got.MaxItems != 1 {
		t.Errorf("expected MaxItems 1, got %d", got.MaxItems)
	}

	elem, ok := got.Elem.(*sdkv2schema.Resource)
	if !ok {
		t.Fatalf("expected Elem to be *schema.Resource, got %T", got.Elem)
	}

	urlSchema, ok := elem.Schema["url"]
	if !ok {
		t.Fatal("expected url attribute in schema")
	}

	if urlSchema.Type != sdkv2schema.TypeString {
		t.Errorf("expected url TypeString, got %v", urlSchema.Type)
	}

	if !urlSchema.Required {
		t.Error("expected url Required to be true")
	}
}

func TestFwToSdkv2SchemaMap_Attributes(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		input    fwschema.Schema
		expected map[string]*sdkv2schema.Schema
	}{
		"string": {
			input: fwschema.Schema{
				Attributes: map[string]fwschema.Attribute{
					"name": fwschema.StringAttribute{
						Required:    true,
						Description: "A name",
					},
				},
			},
			expected: map[string]*sdkv2schema.Schema{
				"name": {
					Type:        sdkv2schema.TypeString,
					Required:    true,
					Description: "A name",
				},
			},
		},
		"bool": {
			input: fwschema.Schema{
				Attributes: map[string]fwschema.Attribute{
					"enabled": fwschema.BoolAttribute{
						Optional: true,
					},
				},
			},
			expected: map[string]*sdkv2schema.Schema{
				"enabled": {
					Type:     sdkv2schema.TypeBool,
					Optional: true,
				},
			},
		},
		"int64": {
			input: fwschema.Schema{
				Attributes: map[string]fwschema.Attribute{
					"count": fwschema.Int64Attribute{
						Optional: true,
					},
				},
			},
			expected: map[string]*sdkv2schema.Schema{
				"count": {
					Type:     sdkv2schema.TypeInt,
					Optional: true,
				},
			},
		},
		"int32": {
			input: fwschema.Schema{
				Attributes: map[string]fwschema.Attribute{
					"port": fwschema.Int32Attribute{
						Required: true,
					},
				},
			},
			expected: map[string]*sdkv2schema.Schema{
				"port": {
					Type:     sdkv2schema.TypeInt,
					Required: true,
				},
			},
		},
		"float64": {
			input: fwschema.Schema{
				Attributes: map[string]fwschema.Attribute{
					"ratio": fwschema.Float64Attribute{
						Optional: true,
					},
				},
			},
			expected: map[string]*sdkv2schema.Schema{
				"ratio": {
					Type:     sdkv2schema.TypeFloat,
					Optional: true,
				},
			},
		},
		"float32": {
			input: fwschema.Schema{
				Attributes: map[string]fwschema.Attribute{
					"ratio": fwschema.Float32Attribute{
						Optional: true,
					},
				},
			},
			expected: map[string]*sdkv2schema.Schema{
				"ratio": {
					Type:     sdkv2schema.TypeFloat,
					Optional: true,
				},
			},
		},
		"number": {
			input: fwschema.Schema{
				Attributes: map[string]fwschema.Attribute{
					"value": fwschema.NumberAttribute{
						Optional: true,
					},
				},
			},
			expected: map[string]*sdkv2schema.Schema{
				"value": {
					Type:     sdkv2schema.TypeFloat,
					Optional: true,
				},
			},
		},
		"sensitive": {
			input: fwschema.Schema{
				Attributes: map[string]fwschema.Attribute{
					"password": fwschema.StringAttribute{
						Required:  true,
						Sensitive: true,
					},
				},
			},
			expected: map[string]*sdkv2schema.Schema{
				"password": {
					Type:      sdkv2schema.TypeString,
					Required:  true,
					Sensitive: true,
				},
			},
		},
		"deprecated": {
			input: fwschema.Schema{
				Attributes: map[string]fwschema.Attribute{
					"old_field": fwschema.StringAttribute{
						Optional:           true,
						DeprecationMessage: "use new_field instead",
					},
				},
			},
			expected: map[string]*sdkv2schema.Schema{
				"old_field": {
					Type:       sdkv2schema.TypeString,
					Optional:   true,
					Deprecated: "use new_field instead",
				},
			},
		},
		"list-of-strings": {
			input: fwschema.Schema{
				Attributes: map[string]fwschema.Attribute{
					"tags": fwschema.ListAttribute{
						Optional:    true,
						ElementType: types.StringType,
					},
				},
			},
			expected: map[string]*sdkv2schema.Schema{
				"tags": {
					Type:     sdkv2schema.TypeList,
					Optional: true,
					Elem:     &sdkv2schema.Schema{Type: sdkv2schema.TypeString},
				},
			},
		},
		"set-of-strings": {
			input: fwschema.Schema{
				Attributes: map[string]fwschema.Attribute{
					"ids": fwschema.SetAttribute{
						Optional:    true,
						ElementType: types.StringType,
					},
				},
			},
			expected: map[string]*sdkv2schema.Schema{
				"ids": {
					Type:     sdkv2schema.TypeSet,
					Optional: true,
					Elem:     &sdkv2schema.Schema{Type: sdkv2schema.TypeString},
				},
			},
		},
		"map-of-strings": {
			input: fwschema.Schema{
				Attributes: map[string]fwschema.Attribute{
					"labels": fwschema.MapAttribute{
						Optional:    true,
						ElementType: types.StringType,
					},
				},
			},
			expected: map[string]*sdkv2schema.Schema{
				"labels": {
					Type:     sdkv2schema.TypeMap,
					Optional: true,
					Elem:     &sdkv2schema.Schema{Type: sdkv2schema.TypeString},
				},
			},
		},
		"single-nested": {
			input: fwschema.Schema{
				Attributes: map[string]fwschema.Attribute{
					"config": fwschema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]fwschema.Attribute{
							"host": fwschema.StringAttribute{Required: true},
						},
					},
				},
			},
			expected: map[string]*sdkv2schema.Schema{
				"config": {
					Type:     sdkv2schema.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &sdkv2schema.Resource{
						Schema: map[string]*sdkv2schema.Schema{
							"host": {Type: sdkv2schema.TypeString, Required: true},
						},
					},
				},
			},
		},
		"list-nested": {
			input: fwschema.Schema{
				Attributes: map[string]fwschema.Attribute{
					"rules": fwschema.ListNestedAttribute{
						Optional: true,
						NestedObject: fwschema.NestedAttributeObject{
							Attributes: map[string]fwschema.Attribute{
								"priority": fwschema.Int64Attribute{Required: true},
							},
						},
					},
				},
			},
			expected: map[string]*sdkv2schema.Schema{
				"rules": {
					Type:     sdkv2schema.TypeList,
					Optional: true,
					Elem: &sdkv2schema.Resource{
						Schema: map[string]*sdkv2schema.Schema{
							"priority": {Type: sdkv2schema.TypeInt, Required: true},
						},
					},
				},
			},
		},
		"set-nested": {
			input: fwschema.Schema{
				Attributes: map[string]fwschema.Attribute{
					"endpoints": fwschema.SetNestedAttribute{
						Optional: true,
						NestedObject: fwschema.NestedAttributeObject{
							Attributes: map[string]fwschema.Attribute{
								"url": fwschema.StringAttribute{Required: true},
							},
						},
					},
				},
			},
			expected: map[string]*sdkv2schema.Schema{
				"endpoints": {
					Type:     sdkv2schema.TypeSet,
					Optional: true,
					Elem: &sdkv2schema.Resource{
						Schema: map[string]*sdkv2schema.Schema{
							"url": {Type: sdkv2schema.TypeString, Required: true},
						},
					},
				},
			},
		},
		"multiple-attributes": {
			input: fwschema.Schema{
				Attributes: map[string]fwschema.Attribute{
					"url":      fwschema.StringAttribute{Required: true},
					"insecure": fwschema.BoolAttribute{Optional: true},
					"port":     fwschema.Int64Attribute{Optional: true},
				},
			},
			expected: map[string]*sdkv2schema.Schema{
				"url":      {Type: sdkv2schema.TypeString, Required: true},
				"insecure": {Type: sdkv2schema.TypeBool, Optional: true},
				"port":     {Type: sdkv2schema.TypeInt, Optional: true},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := FwToSdkv2SchemaMap(tc.input)
			if diff := cmp.Diff(tc.expected, got); diff != "" {
				t.Errorf("unexpected result (-want +got):\n%s", diff)
			}
		})
	}
}

func TestFwToSdkv2SchemaMap_Blocks(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		input    fwschema.Schema
		expected map[string]*sdkv2schema.Schema
	}{
		"single-nested-block": {
			input: fwschema.Schema{
				Blocks: map[string]fwschema.Block{
					"settings": fwschema.SingleNestedBlock{
						Description: "Provider settings",
						Attributes: map[string]fwschema.Attribute{
							"debug": fwschema.BoolAttribute{Optional: true},
						},
					},
				},
			},
			expected: map[string]*sdkv2schema.Schema{
				"settings": {
					Type:        sdkv2schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: "Provider settings",
					Elem: &sdkv2schema.Resource{
						Schema: map[string]*sdkv2schema.Schema{
							"debug": {Type: sdkv2schema.TypeBool, Optional: true},
						},
					},
				},
			},
		},
		"list-nested-block": {
			input: fwschema.Schema{
				Blocks: map[string]fwschema.Block{
					"backend": fwschema.ListNestedBlock{
						Description: "Backend configs",
						NestedObject: fwschema.NestedBlockObject{
							Attributes: map[string]fwschema.Attribute{
								"address": fwschema.StringAttribute{Required: true},
							},
						},
					},
				},
			},
			expected: map[string]*sdkv2schema.Schema{
				"backend": {
					Type:        sdkv2schema.TypeList,
					Optional:    true,
					Description: "Backend configs",
					Elem: &sdkv2schema.Resource{
						Schema: map[string]*sdkv2schema.Schema{
							"address": {Type: sdkv2schema.TypeString, Required: true},
						},
					},
				},
			},
		},
		"list-nested-block-with-size-between-validator": {
			input: fwschema.Schema{
				Blocks: map[string]fwschema.Block{
					"morpheus": fwschema.ListNestedBlock{
						NestedObject: fwschema.NestedBlockObject{
							Attributes: map[string]fwschema.Attribute{
								"url": fwschema.StringAttribute{Required: true},
							},
						},
						Validators: []validator.List{
							listvalidator.SizeBetween(0, 1),
						},
					},
				},
			},
			expected: map[string]*sdkv2schema.Schema{
				"morpheus": {
					Type:     sdkv2schema.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &sdkv2schema.Resource{
						Schema: map[string]*sdkv2schema.Schema{
							"url": {Type: sdkv2schema.TypeString, Required: true},
						},
					},
				},
			},
		},
		"list-nested-block-with-size-at-most-validator": {
			input: fwschema.Schema{
				Blocks: map[string]fwschema.Block{
					"config": fwschema.ListNestedBlock{
						NestedObject: fwschema.NestedBlockObject{
							Attributes: map[string]fwschema.Attribute{
								"host": fwschema.StringAttribute{Required: true},
							},
						},
						Validators: []validator.List{
							listvalidator.SizeAtMost(3),
						},
					},
				},
			},
			expected: map[string]*sdkv2schema.Schema{
				"config": {
					Type:     sdkv2schema.TypeList,
					Optional: true,
					MaxItems: 3,
					Elem: &sdkv2schema.Resource{
						Schema: map[string]*sdkv2schema.Schema{
							"host": {Type: sdkv2schema.TypeString, Required: true},
						},
					},
				},
			},
		},
		"list-nested-block-with-size-at-least-validator": {
			input: fwschema.Schema{
				Blocks: map[string]fwschema.Block{
					"backend": fwschema.ListNestedBlock{
						NestedObject: fwschema.NestedBlockObject{
							Attributes: map[string]fwschema.Attribute{
								"addr": fwschema.StringAttribute{Required: true},
							},
						},
						Validators: []validator.List{
							listvalidator.SizeAtLeast(1),
						},
					},
				},
			},
			expected: map[string]*sdkv2schema.Schema{
				"backend": {
					Type:     sdkv2schema.TypeList,
					Optional: true,
					MinItems: 1,
					Elem: &sdkv2schema.Resource{
						Schema: map[string]*sdkv2schema.Schema{
							"addr": {Type: sdkv2schema.TypeString, Required: true},
						},
					},
				},
			},
		},
		"set-nested-block": {
			input: fwschema.Schema{
				Blocks: map[string]fwschema.Block{
					"listener": fwschema.SetNestedBlock{
						NestedObject: fwschema.NestedBlockObject{
							Attributes: map[string]fwschema.Attribute{
								"port": fwschema.Int64Attribute{Required: true},
							},
						},
					},
				},
			},
			expected: map[string]*sdkv2schema.Schema{
				"listener": {
					Type:     sdkv2schema.TypeSet,
					Optional: true,
					Elem: &sdkv2schema.Resource{
						Schema: map[string]*sdkv2schema.Schema{
							"port": {Type: sdkv2schema.TypeInt, Required: true},
						},
					},
				},
			},
		},
		"block-markdown-description-preferred": {
			input: fwschema.Schema{
				Blocks: map[string]fwschema.Block{
					"tls": fwschema.SingleNestedBlock{
						Description:         "plain description",
						MarkdownDescription: "**markdown** description",
						Attributes: map[string]fwschema.Attribute{
							"cert": fwschema.StringAttribute{Required: true},
						},
					},
				},
			},
			expected: map[string]*sdkv2schema.Schema{
				"tls": {
					Type:        sdkv2schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: "**markdown** description",
					Elem: &sdkv2schema.Resource{
						Schema: map[string]*sdkv2schema.Schema{
							"cert": {Type: sdkv2schema.TypeString, Required: true},
						},
					},
				},
			},
		},
		"nested-blocks-within-block": {
			input: fwschema.Schema{
				Blocks: map[string]fwschema.Block{
					"outer": fwschema.SingleNestedBlock{
						Attributes: map[string]fwschema.Attribute{
							"name": fwschema.StringAttribute{Required: true},
						},
						Blocks: map[string]fwschema.Block{
							"inner": fwschema.SingleNestedBlock{
								Attributes: map[string]fwschema.Attribute{
									"value": fwschema.StringAttribute{Optional: true},
								},
							},
						},
					},
				},
			},
			expected: map[string]*sdkv2schema.Schema{
				"outer": {
					Type:     sdkv2schema.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &sdkv2schema.Resource{
						Schema: map[string]*sdkv2schema.Schema{
							"name": {Type: sdkv2schema.TypeString, Required: true},
							"inner": {
								Type:     sdkv2schema.TypeList,
								Optional: true,
								MaxItems: 1,
								Elem: &sdkv2schema.Resource{
									Schema: map[string]*sdkv2schema.Schema{
										"value": {Type: sdkv2schema.TypeString, Optional: true},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := FwToSdkv2SchemaMap(tc.input)
			if diff := cmp.Diff(tc.expected, got); diff != "" {
				t.Errorf("unexpected result (-want +got):\n%s", diff)
			}
		})
	}
}

func TestFwToSdkv2SchemaMap_Empty(t *testing.T) {
	t.Parallel()

	got := FwToSdkv2SchemaMap(fwschema.Schema{})
	if len(got) != 0 {
		t.Errorf("expected empty map, got %d entries", len(got))
	}
}

func TestFwToSdkv2SchemaMap_AttributesAndBlocks(t *testing.T) {
	t.Parallel()

	in := fwschema.Schema{
		Attributes: map[string]fwschema.Attribute{
			"name": fwschema.StringAttribute{Required: true},
		},
		Blocks: map[string]fwschema.Block{
			"config": fwschema.SingleNestedBlock{
				Attributes: map[string]fwschema.Attribute{
					"debug": fwschema.BoolAttribute{Optional: true},
				},
			},
		},
	}

	got := FwToSdkv2SchemaMap(in)

	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}

	if got["name"].Type != sdkv2schema.TypeString {
		t.Errorf("expected name TypeString, got %v", got["name"].Type)
	}

	if got["config"].Type != sdkv2schema.TypeList {
		t.Errorf("expected config TypeList, got %v", got["config"].Type)
	}
}

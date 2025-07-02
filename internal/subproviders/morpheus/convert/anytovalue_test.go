// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package convert

import (
	"context"
	"math/big"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestAnyToValue(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		input      any
		targetType attr.Type
		expected   attr.Value
		expectErr  bool
	}{
		"string": {
			input:      "test",
			targetType: types.StringType,
			expected:   types.StringValue("test"),
		},
		"bool": {
			input:      true,
			targetType: types.BoolType,
			expected:   types.BoolValue(true),
		},
		"number-int": {
			input:      42,
			targetType: types.NumberType,
			expected:   types.NumberValue(new(big.Float).SetInt64(42)),
		},
		"number-float64": {
			input:      42.5,
			targetType: types.NumberType,
			expected:   types.NumberValue(new(big.Float).SetFloat64(42.5)),
		},
		"float64": {
			input:      42.5,
			targetType: types.Float64Type,
			expected:   types.Float64Value(42.5),
		},
		"list-strings": {
			input: []any{"a", "b", "c"},
			targetType: types.ListType{
				ElemType: types.StringType,
			},
			expected: types.ListValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("a"),
					types.StringValue("b"),
					types.StringValue("c"),
				},
			),
		},
		"set-numbers": {
			input: []any{1, 2, 3},
			targetType: types.SetType{
				ElemType: types.NumberType,
			},
			expected: types.SetValueMust(
				types.NumberType,
				[]attr.Value{
					types.NumberValue(new(big.Float).SetInt64(1)),
					types.NumberValue(new(big.Float).SetInt64(2)),
					types.NumberValue(new(big.Float).SetInt64(3)),
				},
			),
		},
		"object": {
			input: map[string]any{
				"name":    "test",
				"enabled": true,
				"count":   42,
			},
			targetType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"name":    types.StringType,
					"enabled": types.BoolType,
					"count":   types.NumberType,
				},
			},
			expected: types.ObjectValueMust(
				map[string]attr.Type{
					"name":    types.StringType,
					"enabled": types.BoolType,
					"count":   types.NumberType,
				},
				map[string]attr.Value{
					"name":    types.StringValue("test"),
					"enabled": types.BoolValue(true),
					"count":   types.NumberValue(new(big.Float).SetInt64(42)),
				},
			),
		},
		"tuple": {
			input: []any{"test", true, 42},
			targetType: types.TupleType{
				ElemTypes: []attr.Type{
					types.StringType,
					types.BoolType,
					types.NumberType,
				},
			},
			expected: types.TupleValueMust(
				[]attr.Type{
					types.StringType,
					types.BoolType,
					types.NumberType,
				},
				[]attr.Value{
					types.StringValue("test"),
					types.BoolValue(true),
					types.NumberValue(new(big.Float).SetInt64(42)),
				},
			),
		},
		"list-empty": {
			input: []any{},
			targetType: types.ListType{
				ElemType: types.StringType,
			},
			expected: types.ListValueMust(
				types.StringType,
				[]attr.Value{},
			),
		},
		"set-empty": {
			input: []any{},
			targetType: types.SetType{
				ElemType: types.NumberType,
			},
			expected: types.SetValueMust(
				types.NumberType,
				[]attr.Value{},
			),
		},
		"object-empty": {
			input: map[string]any{},
			targetType: types.ObjectType{
				AttrTypes: map[string]attr.Type{},
			},
			expected: types.ObjectValueMust(
				map[string]attr.Type{},
				map[string]attr.Value{},
			),
		},
		"tuple-empty": {
			input: []any{},
			targetType: types.TupleType{
				ElemTypes: []attr.Type{},
			},
			expected: types.TupleValueMust(
				[]attr.Type{},
				[]attr.Value{},
			),
		},
		"complex-object": {
			input: map[string]any{
				"name": "test",
				"metadata": map[string]any{
					"labels": []any{"prod", "us-west"},
					"annotations": map[string]any{
						"created_by": "admin",
						"timestamp":  1625148107,
					},
				},
				"spec": map[string]any{
					"replicas": 3,
					"ports":    []any{80, 443},
					"config": map[string]any{
						"enabled": true,
						"timeout": 30.5,
					},
					"tags": []any{"web", "api"},
				},
				"status": []any{"running", true, 42.5},
			},
			targetType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"name": types.StringType,
					"metadata": types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"labels": types.SetType{
								ElemType: types.StringType,
							},
							"annotations": types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"created_by": types.StringType,
									"timestamp":  types.NumberType,
								},
							},
						},
					},
					"spec": types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"replicas": types.NumberType,
							"ports": types.ListType{
								ElemType: types.NumberType,
							},
							"config": types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"enabled": types.BoolType,
									"timeout": types.Float64Type,
								},
							},
							"tags": types.SetType{
								ElemType: types.StringType,
							},
						},
					},
					"status": types.TupleType{
						ElemTypes: []attr.Type{
							types.StringType,
							types.BoolType,
							types.Float64Type,
						},
					},
				},
			},
			expected: types.ObjectValueMust(
				map[string]attr.Type{
					"name": types.StringType,
					"metadata": types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"labels": types.SetType{
								ElemType: types.StringType,
							},
							"annotations": types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"created_by": types.StringType,
									"timestamp":  types.NumberType,
								},
							},
						},
					},
					"spec": types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"replicas": types.NumberType,
							"ports": types.ListType{
								ElemType: types.NumberType,
							},
							"config": types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"enabled": types.BoolType,
									"timeout": types.Float64Type,
								},
							},
							"tags": types.SetType{
								ElemType: types.StringType,
							},
						},
					},
					"status": types.TupleType{
						ElemTypes: []attr.Type{
							types.StringType,
							types.BoolType,
							types.Float64Type,
						},
					},
				},
				map[string]attr.Value{
					"name": types.StringValue("test"),
					"metadata": types.ObjectValueMust(
						map[string]attr.Type{
							"labels": types.SetType{
								ElemType: types.StringType,
							},
							"annotations": types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"created_by": types.StringType,
									"timestamp":  types.NumberType,
								},
							},
						},
						map[string]attr.Value{
							"labels": types.SetValueMust(
								types.StringType,
								[]attr.Value{
									types.StringValue("prod"),
									types.StringValue("us-west"),
								},
							),
							"annotations": types.ObjectValueMust(
								map[string]attr.Type{
									"created_by": types.StringType,
									"timestamp":  types.NumberType,
								},
								map[string]attr.Value{
									"created_by": types.StringValue("admin"),
									"timestamp":  types.NumberValue(new(big.Float).SetInt64(1625148107)),
								},
							),
						},
					),
					"spec": types.ObjectValueMust(
						map[string]attr.Type{
							"replicas": types.NumberType,
							"ports": types.ListType{
								ElemType: types.NumberType,
							},
							"config": types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"enabled": types.BoolType,
									"timeout": types.Float64Type,
								},
							},
							"tags": types.SetType{
								ElemType: types.StringType,
							},
						},
						map[string]attr.Value{
							"replicas": types.NumberValue(new(big.Float).SetInt64(3)),
							"ports": types.ListValueMust(
								types.NumberType,
								[]attr.Value{
									types.NumberValue(new(big.Float).SetInt64(80)),
									types.NumberValue(new(big.Float).SetInt64(443)),
								},
							),
							"config": types.ObjectValueMust(
								map[string]attr.Type{
									"enabled": types.BoolType,
									"timeout": types.Float64Type,
								},
								map[string]attr.Value{
									"enabled": types.BoolValue(true),
									"timeout": types.Float64Value(30.5),
								},
							),
							"tags": types.SetValueMust(
								types.StringType,
								[]attr.Value{
									types.StringValue("web"),
									types.StringValue("api"),
								},
							),
						},
					),
					"status": types.TupleValueMust(
						[]attr.Type{
							types.StringType,
							types.BoolType,
							types.Float64Type,
						},
						[]attr.Value{
							types.StringValue("running"),
							types.BoolValue(true),
							types.Float64Value(42.5),
						},
					),
				},
			),
		},
		"invalid-string": {
			input:      42,
			targetType: types.StringType,
			expectErr:  true,
		},
		"invalid-bool": {
			input:      "true",
			targetType: types.BoolType,
			expectErr:  true,
		},
		"invalid-number": {
			input:      "42",
			targetType: types.NumberType,
			expectErr:  true,
		},
		"invalid-list": {
			input: map[string]any{},
			targetType: types.ListType{
				ElemType: types.StringType,
			},
			expectErr: true,
		},
		"invalid-set": {
			input: map[string]any{},
			targetType: types.SetType{
				ElemType: types.StringType,
			},
			expectErr: true,
		},
		"invalid-tuple": {
			input: map[string]any{},
			targetType: types.TupleType{
				ElemTypes: []attr.Type{types.StringType},
			},
			expectErr: true,
		},
		"invalid-object": {
			input: []any{},
			targetType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"name": types.StringType,
				},
			},
			expectErr: true,
		},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := AnyToValue(context.Background(), tc.input, tc.targetType)
			if err != nil {
				if !tc.expectErr {
					t.Errorf("unexpected error: %s", err)
				}

				return
			}

			if tc.expectErr {
				t.Error("expected error, got none")

				return
			}

			if diff := cmp.Diff(got, tc.expected); diff != "" {
				t.Errorf("unexpected diff: %s", diff)
			}
		})
	}
}

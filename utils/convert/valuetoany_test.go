package convert

import (
	"context"
	"math/big"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestValueToAnyObject(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		input    attr.Value
		expected any
		wantErr  bool
	}{
		"simple object": {
			input: types.ObjectValueMust(
				map[string]attr.Type{
					"name": types.StringType,
					"age":  types.NumberType,
				},
				map[string]attr.Value{
					"name": types.StringValue("John Doe"),
					"age":  types.NumberValue(big.NewFloat(42)),
				},
			),
			expected: map[string]any{
				"name": "John Doe",
				"age":  float64(42),
			},
			wantErr: false,
		},
		"object with multiple types": {
			input: types.ObjectValueMust(
				map[string]attr.Type{
					"string_val": types.StringType,
					"number_val": types.NumberType,
					"bool_val":   types.BoolType,
				},
				map[string]attr.Value{
					"string_val": types.StringValue("hello"),
					"number_val": types.NumberValue(big.NewFloat(123.45)),
					"bool_val":   types.BoolValue(true),
				},
			),
			expected: map[string]any{
				"string_val": "hello",
				"number_val": float64(123.45),
				"bool_val":   true,
			},
			wantErr: false,
		},
		"nested object": {
			input: types.ObjectValueMust(
				map[string]attr.Type{
					"person": types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"name":  types.StringType,
							"email": types.StringType,
						},
					},
				},
				map[string]attr.Value{
					"person": types.ObjectValueMust(
						map[string]attr.Type{
							"name":  types.StringType,
							"email": types.StringType,
						},
						map[string]attr.Value{
							"name":  types.StringValue("Jane Smith"),
							"email": types.StringValue("jane@example.com"),
						},
					),
				},
			),
			expected: map[string]any{
				"person": map[string]any{
					"name":  "Jane Smith",
					"email": "jane@example.com",
				},
			},
			wantErr: false,
		},
		"object with list": {
			input: types.ObjectValueMust(
				map[string]attr.Type{
					"name": types.StringType,
					"tags": types.ListType{
						ElemType: types.StringType,
					},
				},
				map[string]attr.Value{
					"name": types.StringValue("resource"),
					"tags": types.ListValueMust(
						types.StringType,
						[]attr.Value{
							types.StringValue("production"),
							types.StringValue("web"),
							types.StringValue("2025"),
						},
					),
				},
			),
			expected: map[string]any{
				"name": "resource",
				"tags": []any{"production", "web", "2025"},
			},
			wantErr: false,
		},
		"empty object": {
			input: types.ObjectValueMust(
				map[string]attr.Type{},
				map[string]attr.Value{},
			),
			expected: map[string]any{},
			wantErr:  false,
		},
		"null object": {
			input: types.ObjectNull(
				map[string]attr.Type{
					"name": types.StringType,
					"age":  types.NumberType,
				},
			),
			expected: nil,
			wantErr:  false,
		},
		"object with null attribute": {
			input: types.ObjectValueMust(
				map[string]attr.Type{
					"name":        types.StringType,
					"description": types.StringType,
				},
				map[string]attr.Value{
					"name":        types.StringValue("test"),
					"description": types.StringNull(),
				},
			),
			expected: map[string]any{
				"name":        "test",
				"description": nil,
			},
			wantErr: false,
		},
		"deeply nested object": {
			input: types.ObjectValueMust(
				map[string]attr.Type{
					"level1": types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"level2": types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"level3": types.StringType,
								},
							},
						},
					},
				},
				map[string]attr.Value{
					"level1": types.ObjectValueMust(
						map[string]attr.Type{
							"level2": types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"level3": types.StringType,
								},
							},
						},
						map[string]attr.Value{
							"level2": types.ObjectValueMust(
								map[string]attr.Type{
									"level3": types.StringType,
								},
								map[string]attr.Value{
									"level3": types.StringValue("deep value"),
								},
							),
						},
					),
				},
			),
			expected: map[string]any{
				"level1": map[string]any{
					"level2": map[string]any{
						"level3": "deep value",
					},
				},
			},
			wantErr: false,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result, err := ValueToAny(context.Background(), tc.input)

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				if diff := cmp.Diff(tc.expected, result); diff != "" {
					t.Errorf("unexpected result: %s", diff)
				}
			}
		})
	}
}

func TestValueToAnyPrimitives(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		input    attr.Value
		expected any
		wantErr  bool
	}{
		"string value": {
			input:    types.StringValue("test string"),
			expected: "test string",
			wantErr:  false,
		},
		"bool value true": {
			input:    types.BoolValue(true),
			expected: true,
			wantErr:  false,
		},
		"bool value false": {
			input:    types.BoolValue(false),
			expected: false,
			wantErr:  false,
		},
		"number value integer": {
			input:    types.NumberValue(big.NewFloat(42)),
			expected: float64(42),
			wantErr:  false,
		},
		"number value float": {
			input:    types.NumberValue(big.NewFloat(123.45)),
			expected: float64(123.45),
			wantErr:  false,
		},
		"float64 value": {
			input:    types.Float64Value(987.65),
			expected: float64(987.65),
			wantErr:  false,
		},
		"null string": {
			input:    types.StringNull(),
			expected: nil,
			wantErr:  false,
		},
		"null bool": {
			input:    types.BoolNull(),
			expected: nil,
			wantErr:  false,
		},
		"null number": {
			input:    types.NumberNull(),
			expected: nil,
			wantErr:  false,
		},
		"null float64": {
			input:    types.Float64Null(),
			expected: nil,
			wantErr:  false,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result, err := ValueToAny(context.Background(), tc.input)

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				if diff := cmp.Diff(tc.expected, result); diff != "" {
					t.Errorf("unexpected result: %s", diff)
				}
			}
		})
	}
}

func TestValueToAnyList(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		input    attr.Value
		expected any
		wantErr  bool
	}{
		"string list": {
			input: types.ListValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("a"),
					types.StringValue("b"),
					types.StringValue("c"),
				},
			),
			expected: []any{"a", "b", "c"},
			wantErr:  false,
		},
		"number list": {
			input: types.ListValueMust(
				types.NumberType,
				[]attr.Value{
					types.NumberValue(big.NewFloat(1)),
					types.NumberValue(big.NewFloat(2)),
					types.NumberValue(big.NewFloat(3.5)),
				},
			),
			expected: []any{float64(1), float64(2), float64(3.5)},
			wantErr:  false,
		},
		"bool list": {
			input: types.ListValueMust(
				types.BoolType,
				[]attr.Value{
					types.BoolValue(true),
					types.BoolValue(false),
					types.BoolValue(true),
				},
			),
			expected: []any{true, false, true},
			wantErr:  false,
		},
		"empty list": {
			input: types.ListValueMust(
				types.StringType,
				[]attr.Value{},
			),
			expected: []any{},
			wantErr:  false,
		},
		"null list": {
			input:    types.ListNull(types.StringType),
			expected: nil,
			wantErr:  false,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result, err := ValueToAny(context.Background(), tc.input)

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				if diff := cmp.Diff(tc.expected, result); diff != "" {
					t.Errorf("unexpected result: %s", diff)
				}
			}
		})
	}
}

func TestValueToAnySet(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		input    attr.Value
		expected any
		wantErr  bool
	}{
		"string set": {
			input: types.SetValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("a"),
					types.StringValue("b"),
					types.StringValue("c"),
				},
			),
			expected: []any{"a", "b", "c"},
			wantErr:  false,
		},
		"number set": {
			input: types.SetValueMust(
				types.NumberType,
				[]attr.Value{
					types.NumberValue(big.NewFloat(1)),
					types.NumberValue(big.NewFloat(2)),
					types.NumberValue(big.NewFloat(3.5)),
				},
			),
			expected: []any{float64(1), float64(2), float64(3.5)},
			wantErr:  false,
		},
		"bool set": {
			input: types.SetValueMust(
				types.BoolType,
				[]attr.Value{
					types.BoolValue(true),
					types.BoolValue(false),
				},
			),
			expected: []any{true, false},
			wantErr:  false,
		},
		"empty set": {
			input: types.SetValueMust(
				types.StringType,
				[]attr.Value{},
			),
			expected: []any{},
			wantErr:  false,
		},
		"null set": {
			input:    types.SetNull(types.StringType),
			expected: nil,
			wantErr:  false,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result, err := ValueToAny(context.Background(), tc.input)

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				// Note: Set elements don't have guaranteed order, so we'd need a more sophisticated
				// comparison for real-world testing. For this example, we assume the order is preserved.
				if diff := cmp.Diff(tc.expected, result); diff != "" {
					t.Errorf("unexpected result: %s", diff)
				}
			}
		})
	}
}

func TestValueToAnyTuple(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		input    attr.Value
		expected any
		wantErr  bool
	}{
		"mixed tuple": {
			input: types.TupleValueMust(
				[]attr.Type{types.StringType, types.NumberType, types.BoolType},
				[]attr.Value{
					types.StringValue("test"),
					types.NumberValue(big.NewFloat(42)),
					types.BoolValue(true),
				},
			),
			expected: []any{"test", float64(42), true},
			wantErr:  false,
		},
		"string tuple": {
			input: types.TupleValueMust(
				[]attr.Type{types.StringType, types.StringType},
				[]attr.Value{
					types.StringValue("a"),
					types.StringValue("b"),
				},
			),
			expected: []any{"a", "b"},
			wantErr:  false,
		},
		"empty tuple": {
			input: types.TupleValueMust(
				[]attr.Type{},
				[]attr.Value{},
			),
			expected: []any{},
			wantErr:  false,
		},
		"null tuple": {
			input:    types.TupleNull([]attr.Type{types.StringType, types.NumberType}),
			expected: nil,
			wantErr:  false,
		},
		"tuple with null element": {
			input: types.TupleValueMust(
				[]attr.Type{types.StringType, types.NumberType},
				[]attr.Value{
					types.StringValue("test"),
					types.NumberNull(),
				},
			),
			expected: []any{"test", nil},
			wantErr:  false,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result, err := ValueToAny(context.Background(), tc.input)

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				if diff := cmp.Diff(tc.expected, result); diff != "" {
					t.Errorf("unexpected result: %s", diff)
				}
			}
		})
	}
}

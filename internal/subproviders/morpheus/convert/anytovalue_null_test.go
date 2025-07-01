// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package convert

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
)

// mockUnsupportedType implements attr.Type for testing unsupported types
type mockUnsupportedType struct{}

func (t mockUnsupportedType) ApplyTerraform5AttributePathStep(
	step tftypes.AttributePathStep,
) (interface{}, error) {
	return nil, fmt.Errorf(
		"cannot apply AttributePathStep %T",
		step,
	)
}

func (t mockUnsupportedType) String() string {
	return "mock"
}

func (t mockUnsupportedType) Equal(o attr.Type) bool {
	_, ok := o.(mockUnsupportedType)

	return ok
}

func (t mockUnsupportedType) Type() attr.Type {
	return mockUnsupportedType{}
}

func (t mockUnsupportedType) ValueType(_ context.Context) attr.Value {
	return nil
}

func (t mockUnsupportedType) ValueFromTerraform(
	_ context.Context,
	_ tftypes.Value,
) (attr.Value, error) {
	return nil, fmt.Errorf("cannot convert value")
}

func (t mockUnsupportedType) TerraformType(_ context.Context) tftypes.Type {
	return tftypes.DynamicPseudoType
}

func (t mockUnsupportedType) Validate(
	_ context.Context,
	_ tftypes.Value,
	_ path.Path,
) diag.Diagnostics {
	return nil
}

func TestAnyToValueNullCases(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		input      any
		targetType attr.Type
		expected   attr.Value
	}{
		"string-null": {
			input:      nil,
			targetType: types.StringType,
			expected:   types.StringNull(),
		},
		"bool-null": {
			input:      nil,
			targetType: types.BoolType,
			expected:   types.BoolNull(),
		},
		"int64-null": {
			input:      nil,
			targetType: types.Int64Type,
			expected:   types.Int64Null(),
		},
		"float64-null": {
			input:      nil,
			targetType: types.Float64Type,
			expected:   types.Float64Null(),
		},
		"number-null": {
			input:      nil,
			targetType: types.NumberType,
			expected:   types.NumberNull(),
		},
		"list-null": {
			input:      nil,
			targetType: types.ListType{ElemType: types.StringType},
			expected:   types.ListNull(types.StringType),
		},
		"set-null": {
			input:      nil,
			targetType: types.SetType{ElemType: types.StringType},
			expected:   types.SetNull(types.StringType),
		},
		"map-null": {
			input:      nil,
			targetType: types.MapType{ElemType: types.StringType},
			expected:   types.MapNull(types.StringType),
		},
		"object-null": {
			input: nil,
			targetType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"test": types.StringType,
				},
			},
			expected: types.ObjectNull(map[string]attr.Type{
				"test": types.StringType,
			}),
		},
		"tuple-null": {
			input: nil,
			targetType: types.TupleType{
				ElemTypes: []attr.Type{types.StringType},
			},
			expected: types.TupleNull([]attr.Type{
				types.StringType,
			}),
		},
	}

	ctx := context.Background()

	for name, tc := range cases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := AnyToValue(ctx, tc.input, tc.targetType)
			if err != nil {
				t.Fatalf(
					"unexpected error: %s",
					err,
				)
			}

			assert.Equal(
				t,
				tc.expected,
				got,
				"values should match",
			)
		})
	}
}

func TestAnyToValueNilTargetType(t *testing.T) {
	t.Parallel()

	_, err := AnyToValue(
		context.Background(),
		"test",
		nil,
	)

	assert.EqualError(
		t,
		err,
		"target type is required but was nil",
		"should error when target type is nil",
	)
}

func TestAnyToValueUnsupportedType(t *testing.T) {
	t.Parallel()

	_, err := AnyToValue(
		context.Background(),
		nil,
		mockUnsupportedType{},
	)

	assert.EqualError(
		t,
		err,
		"unsupported type: convert.mockUnsupportedType",
		"should error with unsupported type message",
	)
}

func TestNullToValueUnsupportedType(t *testing.T) {
	t.Parallel()

	_, err := NullToValue(mockUnsupportedType{})

	assert.EqualError(
		t,
		err,
		"unsupported type: convert.mockUnsupportedType",
		"should error with unsupported type message",
	)
}

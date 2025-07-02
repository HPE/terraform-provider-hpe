// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package convert

import (
	"context"
	"fmt"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// AnyToValue allows converting from an arbitrary go struct (eg map[string]any)
// to a framework value of the specified target type. The target type is used to
// resolve ambiguity: for example it allows us to know if a go slice should be
// converted to a framework list or set.
func AnyToValue(
	ctx context.Context,
	a any,
	targetType attr.Type,
) (
	attr.Value,
	error,
) {
	if targetType == nil {
		return nil, fmt.Errorf(
			"target type is required but was nil",
		)
	}

	if a == nil {
		return NullToValue(targetType)
	}

	switch t := targetType.(type) {
	case basetypes.StringType:
		s, ok := a.(string)
		if !ok {
			return nil, fmt.Errorf(
				"expected string, got %T",
				a,
			)
		}

		return types.StringValue(s), nil

	case basetypes.BoolType:
		b, ok := a.(bool)
		if !ok {
			return nil, fmt.Errorf(
				"expected bool, got %T",
				a,
			)
		}

		return types.BoolValue(b), nil

	case basetypes.NumberType:
		switch v := a.(type) {
		case int:
			return types.NumberValue(new(big.Float).SetInt64(int64(v))), nil
		case int64:
			return types.NumberValue(new(big.Float).SetInt64(v)), nil
		case float64:
			return types.NumberValue(new(big.Float).SetFloat64(v)), nil
		default:
			return nil, fmt.Errorf(
				"expected number, got %T",
				a,
			)
		}

	case basetypes.Float64Type:
		switch v := a.(type) {
		case float64:
			return types.Float64Value(v), nil
		case int:
			return types.Float64Value(float64(v)), nil
		case int64:
			return types.Float64Value(float64(v)), nil
		default:
			return nil, fmt.Errorf(
				"expected float64, got %T",
				a,
			)
		}

	case types.ListType:
		l, ok := a.([]any)
		if !ok {
			return nil, fmt.Errorf(
				"expected slice, got %T",
				a,
			)
		}

		return ListToValue(ctx, l, t)

	case types.SetType:
		s, ok := a.([]any)
		if !ok {
			return nil, fmt.Errorf(
				"expected slice, got %T",
				a,
			)
		}

		return SetToValue(ctx, s, t)

	case types.TupleType:
		tup, ok := a.([]any)
		if !ok {
			return nil, fmt.Errorf(
				"expected slice, got %T",
				a,
			)
		}

		return TupleToValue(ctx, tup, t)

	case types.ObjectType:
		m, ok := a.(map[string]any)
		if !ok {
			return nil, fmt.Errorf(
				"expected map, got %T",
				a,
			)
		}

		return MapToValue(ctx, m, t)

	default:
		return nil, fmt.Errorf(
			"unsupported type: %T",
			targetType,
		)
	}
}

func NullToValue(targetType attr.Type) (attr.Value, error) {
	// primitives require "basetypes" while collections require "types"
	switch t := targetType.(type) {
	case basetypes.StringType:
		return types.StringNull(), nil
	case basetypes.BoolType:
		return types.BoolNull(), nil
	case basetypes.Int64Type:
		return types.Int64Null(), nil
	case basetypes.Float64Type:
		return types.Float64Null(), nil
	case basetypes.NumberType:
		return types.NumberNull(), nil
	case types.ListType:
		return types.ListNull(t.ElemType), nil
	case types.SetType:
		return types.SetNull(t.ElemType), nil
	case types.MapType:
		return types.MapNull(t.ElemType), nil
	case types.ObjectType:
		return types.ObjectNull(t.AttrTypes), nil
	case types.TupleType:
		return types.TupleNull(t.ElemTypes), nil
	case attr.TypeWithElementTypes:
		return types.TupleNull(t.ElementTypes()), nil
	default:
		return nil, fmt.Errorf(
			"unsupported type: %T",
			targetType,
		)
	}
}

func MapToValue(
	ctx context.Context,
	m map[string]any,
	targetType attr.Type,
) (
	attr.Value,
	error,
) {
	typ, ok := targetType.(types.ObjectType)
	if !ok {
		return nil, fmt.Errorf(
			"expected Object type for map, got %T",
			targetType,
		)
	}

	if len(m) == 0 {
		return types.ObjectValueMust(
			typ.AttrTypes,
			map[string]attr.Value{},
		), nil
	}

	vm := make(map[string]attr.Value, len(m))
	tm := make(map[string]attr.Type, len(typ.AttrTypes))

	// Only process keys that have target types
	for k, v := range m {
		if targetValue, exists := typ.AttrTypes[k]; exists {
			vv, err := AnyToValue(ctx, v, targetValue)
			if err != nil {
				return types.ObjectNull(typ.AttrTypes),
					fmt.Errorf(
						"error decoding key %q: %w",
						k,
						err,
					)
			}
			vm[k] = vv
			tm[k] = targetValue
		} else {
			tflog.Trace(
				ctx,
				"skipping map key with no target type",
				map[string]any{
					"key":   k,
					"value": fmt.Sprintf("%v", v),
				},
			)
		}
	}

	return types.ObjectValueMust(tm, vm), nil
}

func CollectionToValue(
	ctx context.Context,
	c []any,
	targetType attr.Type,
) (
	attr.Value,
	error,
) {
	switch t := targetType.(type) {
	case types.ListType:
		return ListToValue(ctx, c, t)
	case types.SetType:
		return SetToValue(ctx, c, t)
	case types.TupleType:
		return TupleToValue(ctx, c, t)
	default:
		return nil, fmt.Errorf(
			"expected List, Set, or Tuple type for array input, "+
				"got %T",
			targetType,
		)
	}
}

func ListToValue(
	ctx context.Context,
	l []any,
	listType types.ListType,
) (
	attr.Value,
	error,
) {
	if len(l) == 0 {
		return types.ListValueMust(
			listType.ElemType,
			[]attr.Value{},
		), nil
	}

	vl := make([]attr.Value, len(l))
	for i, v := range l {
		val, err := AnyToValue(ctx, v, listType.ElemType)
		if err != nil {
			return nil, fmt.Errorf(
				"error decoding list element %d: %w",
				i,
				err,
			)
		}
		vl[i] = val
	}

	return types.ListValueMust(listType.ElemType, vl), nil
}

func SetToValue(
	ctx context.Context,
	s []any,
	setType types.SetType,
) (
	attr.Value,
	error,
) {
	if len(s) == 0 {
		return types.SetValueMust(
			setType.ElemType,
			[]attr.Value{},
		), nil
	}

	vl := make([]attr.Value, len(s))
	for i, v := range s {
		val, err := AnyToValue(ctx, v, setType.ElemType)
		if err != nil {
			return nil, fmt.Errorf(
				"error decoding set element %d: %w",
				i,
				err,
			)
		}
		vl[i] = val
	}

	return types.SetValueMust(setType.ElemType, vl), nil
}

func TupleToValue(
	ctx context.Context,
	t []any,
	tupleType types.TupleType,
) (
	attr.Value,
	error,
) {
	if len(t) == 0 {
		if len(tupleType.ElemTypes) > 0 {
			return nil, fmt.Errorf(
				"invalid tuple length: expected %d elements, "+
					"got 0",
				len(tupleType.ElemTypes),
			)
		}

		return types.TupleValueMust(
			tupleType.ElemTypes,
			[]attr.Value{},
		), nil
	}

	if len(t) != len(tupleType.ElemTypes) {
		return nil, fmt.Errorf(
			"invalid tuple length: expected %d elements, got %d",
			len(tupleType.ElemTypes),
			len(t),
		)
	}

	vl := make([]attr.Value, len(t))
	for i, v := range t {
		val, err := AnyToValue(ctx, v, tupleType.ElemTypes[i])
		if err != nil {
			return nil, fmt.Errorf(
				"error decoding tuple element %d: %w",
				i,
				err,
			)
		}
		vl[i] = val
	}

	return types.TupleValueMust(tupleType.ElemTypes, vl), nil
}

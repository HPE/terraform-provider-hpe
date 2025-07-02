// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package convert

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// ValueToAny converts a Terraform Plugin Framework Value into a Go type.
// This function handles null values, primitive types, and complex types
// like objects, lists, sets, tuples and maps.
func ValueToAny(ctx context.Context, v attr.Value) (any, error) {
	if v == nil || v.IsNull() {
		return nil, nil
	}

	switch val := v.(type) {
	case basetypes.StringValue:
		return val.ValueString(), nil
	case basetypes.BoolValue:
		return val.ValueBool(), nil
	case basetypes.NumberValue:
		f, acc := val.ValueBigFloat().Float64()
		if acc != 0 {
			return nil, fmt.Errorf(
				"loss of precision converting number value %v to float64",
				val.ValueBigFloat())
		}

		return f, nil
	case basetypes.Float64Value:
		return val.ValueFloat64(), nil
	case basetypes.ListValue:
		return ListToAny(ctx, val)
	case basetypes.SetValue:
		return SetToAny(ctx, val)
	case basetypes.MapValue:
		return MapToAny(ctx, val)
	case basetypes.ObjectValue:
		return ObjectToAny(ctx, val)
	case basetypes.TupleValue:
		return TupleToAny(ctx, val)
	default:
		return nil, fmt.Errorf("unsupported type for ValueToAny conversion: %T", v)
	}
}

func ListToAny(ctx context.Context, l basetypes.ListValue) ([]any, error) {
	if l.IsNull() {
		return nil, nil
	}

	elems := l.Elements()
	result := make([]any, len(elems))

	for i, elem := range elems {
		var err error
		result[i], err = ValueToAny(ctx, elem)
		if err != nil {
			return nil, fmt.Errorf("error converting list element %d: %w", i, err)
		}
	}

	return result, nil
}

func SetToAny(ctx context.Context, s basetypes.SetValue) ([]any, error) {
	if s.IsNull() {
		return nil, nil
	}

	elems := s.Elements()
	result := make([]any, len(elems))

	for i, elem := range elems {
		var err error
		result[i], err = ValueToAny(ctx, elem)
		if err != nil {
			return nil, fmt.Errorf("error converting set element %d: %w", i, err)
		}
	}

	return result, nil
}

func MapToAny(ctx context.Context, m basetypes.MapValue) (map[string]any, error) {
	if m.IsNull() {
		return nil, nil
	}

	elems := m.Elements()
	result := make(map[string]any, len(elems))

	for k, v := range elems {
		var err error
		result[k], err = ValueToAny(ctx, v)
		if err != nil {
			return nil, fmt.Errorf("error converting map value for key %q: %w", k, err)
		}
	}

	return result, nil
}

func ObjectToAny(ctx context.Context, o basetypes.ObjectValue) (map[string]any, error) {
	if o.IsNull() {
		return nil, nil
	}

	attrs := o.Attributes()
	result := make(map[string]any, len(attrs))

	for k, v := range attrs {
		var err error
		result[k], err = ValueToAny(ctx, v)
		if err != nil {
			return nil, fmt.Errorf("error converting object attribute %q: %w", k, err)
		}
	}

	return result, nil
}

func TupleToAny(ctx context.Context, t basetypes.TupleValue) ([]any, error) {
	if t.IsNull() {
		return nil, nil
	}

	elems := t.Elements()
	result := make([]any, len(elems))

	for i, elem := range elems {
		var err error
		result[i], err = ValueToAny(ctx, elem)
		if err != nil {
			return nil, fmt.Errorf("error converting tuple element %d: %w", i, err)
		}
	}

	return result, nil
}

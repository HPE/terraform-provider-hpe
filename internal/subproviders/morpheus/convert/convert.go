// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package convert

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func StrToType(s *string) types.String {
	if s == nil {
		return types.StringNull()
	}

	return types.StringValue(*s)
}

func StrSliceToSet(items []string) types.Set {
	if len(items) == 0 {
		return types.SetNull(types.StringType)
	}

	var vals []attr.Value
	for _, i := range items {
		vals = append(vals, types.StringValue(i))
	}

	set, diags := types.SetValue(types.StringType, vals)
	if diags.HasError() {
		return types.SetNull(types.StringType)
	}

	return set
}

func SetToStrSlice(set types.Set) ([]string, error) {
	var items []string

	for _, elem := range set.Elements() {
		switch val := elem.(type) {
		case basetypes.StringValue:
			items = append(items, val.ValueString())
		default:
			return nil, fmt.Errorf("value %v is not a string", val)
		}
	}

	return items, nil
}

func BoolToType(b *bool) types.Bool {
	if b == nil {
		return types.BoolNull()
	}

	return types.BoolValue(*b)
}

func Int64ToType(i *int64) types.Int64 {
	if i == nil {
		return types.Int64Null()
	}

	return types.Int64Value(*i)
}

func Int64SliceToSet(items []int64) types.Set {
	if len(items) == 0 {
		return types.SetNull(types.Int64Type)
	}

	var vals []attr.Value
	for _, i := range items {
		vals = append(vals, types.Int64Value(i))
	}

	set, diags := types.SetValue(types.Int64Type, vals)
	if diags.HasError() {
		return types.SetNull(types.Int64Type)
	}

	return set
}

type MappingFunc[I any, O any] func(in I) O

// Map objects in a slice into a Terraform Set Type according to the mapping
// function.
// Mapping function must take in an arbitrary Go value (e.g. API response
// struct) and return the mapped Terraform value.
func ToSetType[S any, O basetypes.ObjectValuable](
	ctx context.Context,
	slice []S,
	mapper MappingFunc[S, O],
) (basetypes.SetValue, diag.Diagnostics) {
	values := []attr.Value{}
	var obj O
	v, _ := obj.ToObjectValue(ctx)

	if len(slice) == 0 {
		return basetypes.NewSetNull(v.Type(ctx)), nil
	}

	for _, i := range slice {
		v := mapper(i)

		obj, d := v.ToObjectValue(ctx)
		if d.HasError() {
			return types.SetUnknown(basetypes.ObjectType{}), d
		}

		values = append(values, obj)
	}

	return types.SetValue(v.Type(ctx), values)
}

// Map objects in a slice into a Terraform List Type according to the mapping
// function.
// Mapping function must take in an arbitrary Go value (e.g. API response
// struct) and return the mapped Terraform value.
func ToListType[S any, O basetypes.ObjectValuable](
	ctx context.Context,
	slice []S,
	mapper MappingFunc[S, O],
) (basetypes.ListValue, diag.Diagnostics) {
	values := []attr.Value{}
	var obj O
	v, _ := obj.ToObjectValue(ctx)

	if len(slice) == 0 {
		return basetypes.NewListNull(v.Type(ctx)), nil
	}

	for _, i := range slice {
		v := mapper(i)

		obj, d := v.ToObjectValue(ctx)
		if d.HasError() {
			return types.ListUnknown(basetypes.ObjectType{}), d
		}

		values = append(values, obj)
	}

	return types.ListValue(v.Type(ctx), values)
}

// Map List objects into a slice of objects according to the mapping function.
// Mapping function must take in a Terraform value and return the mapped object
// (e.g. fill out an API request struct)
func FromSetType[S attr.Value, O any](
	ctx context.Context,
	set types.Set,
	mapper MappingFunc[S, O],
) ([]O, diag.Diagnostics) {
	var out []O
	var elems []S

	if diags := set.ElementsAs(ctx, &elems, false); diags.HasError() {
		return nil, diags
	}

	for _, el := range elems {
		out = append(out, mapper(el))
	}

	return out, nil
}

// Map List objects into a slice of objects according to the mapping function.
// Mapping function must take in a Terraform value and return the mapped object
// (e.g. fill out an API request struct)
func FromListType[S attr.Value, O any](
	ctx context.Context,
	list types.List,
	mapper MappingFunc[S, O],
) ([]O, diag.Diagnostics) {
	var out []O
	var elems []S

	if diags := list.ElementsAs(ctx, &elems, false); diags.HasError() {
		return nil, diags
	}

	for _, el := range elems {
		out = append(out, mapper(el))
	}

	return out, nil
}

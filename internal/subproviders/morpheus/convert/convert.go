// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package convert

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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

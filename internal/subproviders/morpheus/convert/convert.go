// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package convert

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

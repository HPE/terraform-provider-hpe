// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package convert

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func StrToType(s *string) types.String {
	if s == nil {
		return types.StringNull()
	}

	return types.StringValue(*s)
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

// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package convert_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

func strPtr(s string) *string { return &s }
func i64Ptr(i int64) *int64   { return &i }
func boolPtr(b bool) *bool    { return &b }

func TestStrOrPlan(t *testing.T) {
	tests := []struct {
		name string
		api  *string
		plan types.String
		want types.String
	}{
		{"api non-nil overrides plan", strPtr("api"), types.StringValue("plan"), types.StringValue("api")},
		{"api nil, plan known", nil, types.StringValue("plan"), types.StringValue("plan")},
		{"api nil, plan known-null", nil, types.StringNull(), types.StringNull()},
		{"api nil, plan unknown", nil, types.StringUnknown(), types.StringNull()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convert.StrOrPlan(tt.api, tt.plan)
			if !got.Equal(tt.want) {
				t.Errorf("StrOrPlan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64OrPlan(t *testing.T) {
	tests := []struct {
		name string
		api  *int64
		plan types.Int64
		want types.Int64
	}{
		{"api non-nil overrides plan", i64Ptr(2), types.Int64Value(1), types.Int64Value(2)},
		{"api nil, plan known", nil, types.Int64Value(1), types.Int64Value(1)},
		{"api nil, plan known-null", nil, types.Int64Null(), types.Int64Null()},
		{"api nil, plan unknown", nil, types.Int64Unknown(), types.Int64Null()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convert.Int64OrPlan(tt.api, tt.plan)
			if !got.Equal(tt.want) {
				t.Errorf("Int64OrPlan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBoolOrPlan(t *testing.T) {
	tests := []struct {
		name string
		api  *bool
		plan types.Bool
		want types.Bool
	}{
		{"api non-nil overrides plan", boolPtr(false), types.BoolValue(true), types.BoolValue(false)},
		{"api nil, plan known", nil, types.BoolValue(true), types.BoolValue(true)},
		{"api nil, plan known-null", nil, types.BoolNull(), types.BoolNull()},
		{"api nil, plan unknown", nil, types.BoolUnknown(), types.BoolNull()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convert.BoolOrPlan(tt.api, tt.plan)
			if !got.Equal(tt.want) {
				t.Errorf("BoolOrPlan() = %v, want %v", got, tt.want)
			}
		})
	}
}

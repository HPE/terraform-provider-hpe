// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package validators_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/HPE/terraform-provider-hpe/utils/validators"
)

// TestUnitGreaterThanOrEqual exercises the GreaterThanOrEqual validator
// directly, covering the greater/equal cases that must pass, the less-than
// case that must fail, and the null/unknown values that must be skipped.
func TestUnitGreaterThanOrEqual(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// The validated attribute ("max_idle") is compared against a sibling
	// root-level attribute ("min_idle").
	testSchema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"min_idle": schema.Int64Attribute{Optional: true},
			"max_idle": schema.Int64Attribute{Optional: true},
		},
	}

	objectType := tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"min_idle": tftypes.Number,
			"max_idle": tftypes.Number,
		},
	}

	cases := []struct {
		name      string
		configVal types.Int64
		minIdle   *int64
		wantError bool
	}{
		{"greater than", types.Int64Value(5), ptr(2), false},
		{"equal", types.Int64Value(2), ptr(2), false},
		{"less than", types.Int64Value(1), ptr(2), true},
		{"this null skipped", types.Int64Null(), ptr(2), false},
		{"this unknown skipped", types.Int64Unknown(), ptr(2), false},
		{"ref null skipped", types.Int64Value(1), nil, false},
	}

	v := validators.GreaterThanOrEqual("min_idle")

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			minRaw := tftypes.NewValue(tftypes.Number, nil)
			if tc.minIdle != nil {
				minRaw = tftypes.NewValue(tftypes.Number, *tc.minIdle)
			}

			maxRaw := tftypes.NewValue(tftypes.Number, nil)
			switch {
			case tc.configVal.IsUnknown():
				maxRaw = tftypes.NewValue(tftypes.Number, tftypes.UnknownValue)
			case !tc.configVal.IsNull():
				maxRaw = tftypes.NewValue(tftypes.Number, tc.configVal.ValueInt64())
			}

			config := tfsdk.Config{
				Schema: testSchema,
				Raw: tftypes.NewValue(objectType, map[string]tftypes.Value{
					"min_idle": minRaw,
					"max_idle": maxRaw,
				}),
			}

			req := validator.Int64Request{
				Path:        path.Root("max_idle"),
				ConfigValue: tc.configVal,
				Config:      config,
			}
			resp := &validator.Int64Response{}
			v.ValidateInt64(ctx, req, resp)

			if got := resp.Diagnostics.HasError(); got != tc.wantError {
				t.Fatalf("value %v (min_idle %v): got error=%v, want error=%v (diagnostics: %v)",
					tc.configVal, tc.minIdle, got, tc.wantError, resp.Diagnostics)
			}
		})
	}
}

func ptr(v int64) *int64 {
	return &v
}

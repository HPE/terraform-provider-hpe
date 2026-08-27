// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package validators_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/utils/validators"
)

// TestHexColor exercises the HexColor validator directly, covering valid hex
// forms, skipped null/unknown values, and invalid input that must be rejected
// with an attribute error.
func TestUnitHexColor(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	cases := []struct {
		name      string
		value     types.String
		wantError bool
	}{
		{"valid six digit", types.StringValue("#1a73e8"), false},
		{"valid three digit", types.StringValue("#fff"), false},
		{"valid four digit alpha", types.StringValue("#abcd"), false},
		{"valid eight digit alpha", types.StringValue("#1a73e8ff"), false},
		{"valid uppercase", types.StringValue("#ABCDEF"), false},
		{"null skipped", types.StringNull(), false},
		{"unknown skipped", types.StringUnknown(), false},
		{"invalid word", types.StringValue("not-a-color"), true},
		{"invalid named", types.StringValue("red"), true},
		{"invalid missing hash", types.StringValue("1a73e8"), true},
		{"invalid too short", types.StringValue("#12"), true},
		{"invalid non hex", types.StringValue("#12345g"), true},
		{"invalid empty", types.StringValue(""), true},
		{"invalid rgb", types.StringValue("rgb(0,0,0)"), true},
	}

	v := validators.HexColor()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := validator.StringRequest{
				Path:        path.Root("primary_color"),
				ConfigValue: tc.value,
			}
			resp := &validator.StringResponse{}
			v.ValidateString(ctx, req, resp)

			if got := resp.Diagnostics.HasError(); got != tc.wantError {
				t.Fatalf("value %v: got error=%v, want error=%v (diagnostics: %v)",
					tc.value, got, tc.wantError, resp.Diagnostics)
			}
		})
	}
}

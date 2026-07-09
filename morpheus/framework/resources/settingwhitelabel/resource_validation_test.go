// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package settingwhitelabel_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/settingwhitelabel"
)

// TestSettingWhitelabelColorValidators exercises the schema-level hex-color
// validators on primary_color and secondary_color directly, without a live
// Morpheus backend. It guards the regression reported where an invalid color
// (e.g. "not-a-color") was accepted instead of rejected.
func TestSettingWhitelabelColorValidators(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	s := settingwhitelabel.SettingWhitelabelResourceSchema(ctx)

	cases := []struct {
		name      string
		attribute string
		value     string
		wantError bool
	}{
		{"primary valid six digit", "primary_color", "#1a73e8", false},
		{"primary valid three digit", "primary_color", "#fff", false},
		{"primary valid eight digit alpha", "primary_color", "#1a73e8ff", false},
		{"primary valid uppercase", "primary_color", "#ABCDEF", false},
		{"primary invalid word", "primary_color", "not-a-color", true},
		{"primary invalid named", "primary_color", "red", true},
		{"primary invalid missing hash", "primary_color", "1a73e8", true},
		{"primary invalid too short", "primary_color", "#12", true},
		{"primary invalid non hex", "primary_color", "#12345g", true},
		{"primary invalid empty", "primary_color", "", true},
		{"secondary valid six digit", "secondary_color", "#ffffff", false},
		{"secondary valid three digit", "secondary_color", "#161", false},
		{"secondary invalid word", "secondary_color", "not-a-color", true},
		{"secondary invalid rgb", "secondary_color", "rgb(0,0,0)", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			attr, ok := s.Attributes[tc.attribute].(schema.StringAttribute)
			if !ok {
				t.Fatalf("attribute %q is not a schema.StringAttribute", tc.attribute)
			}
			if len(attr.Validators) == 0 {
				t.Fatalf("attribute %q has no validators; expected a hex-color validator", tc.attribute)
			}

			req := validator.StringRequest{
				Path:        path.Root(tc.attribute),
				ConfigValue: types.StringValue(tc.value),
			}
			resp := &validator.StringResponse{}
			for _, v := range attr.Validators {
				v.ValidateString(ctx, req, resp)
			}

			if got := resp.Diagnostics.HasError(); got != tc.wantError {
				t.Fatalf("value %q on %s: got error=%v, want error=%v (diagnostics: %v)",
					tc.value, tc.attribute, got, tc.wantError, resp.Diagnostics)
			}
		})
	}
}

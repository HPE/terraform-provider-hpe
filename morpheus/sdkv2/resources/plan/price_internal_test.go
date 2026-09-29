// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package plan

import (
	"strings"
	"testing"

	morpheus "github.com/HPE/terraform-provider-hpe/internal/sdk/legacy"
)

// TestPriceResultGuard verifies the HTTP-200-with-success:false handling that
// fixes MORPH-15913: the API answers price validation failures with a 2xx
// status and {"success":false,"errors":{...}}, which must become an error
// diagnostic (carrying the field messages) rather than being treated as a
// successful create with id 0. The pre-flight lookup in resourcePriceCreate
// normally intercepts duplicate codes before any POST, so this path cannot be
// reached deterministically by an acceptance test.
func TestPriceResultGuard(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		result     morpheus.CreatePriceResult
		action     string
		wantErr    bool
		wantInText []string
	}{
		{
			name:   "create success with id",
			result: morpheus.CreatePriceResult{Success: true, ID: 42},
			action: "create",
		},
		{
			name:       "create validation failure on HTTP 200",
			result:     morpheus.CreatePriceResult{Success: false, Errors: map[string]string{"code": "must be unique"}},
			action:     "create",
			wantErr:    true,
			wantInText: []string{"code: must be unique"},
		},
		{
			name:       "create success:true without id",
			result:     morpheus.CreatePriceResult{Success: true},
			action:     "create",
			wantErr:    true,
			wantInText: []string{"failed to create price"},
		},
		{
			name:       "message preferred over generic text",
			result:     morpheus.CreatePriceResult{Success: false, Message: "Enter your password"},
			action:     "create",
			wantErr:    true,
			wantInText: []string{"Enter your password"},
		},
		{
			name:   "update success without id is fine",
			result: morpheus.CreatePriceResult{Success: true},
			action: "update",
		},
		{
			name:       "update with field errors",
			result:     morpheus.CreatePriceResult{Success: true, Errors: map[string]string{"cost": "is required"}},
			action:     "update",
			wantErr:    true,
			wantInText: []string{"failed to update price", "cost: is required"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			diags := priceResultGuard(tc.result, tc.action)

			if !tc.wantErr {
				if diags != nil {
					t.Fatalf("expected no diagnostics, got %v", diags)
				}

				return
			}

			if !diags.HasError() {
				t.Fatalf("expected an error diagnostic, got %v", diags)
			}
			text := diags[0].Summary + " " + diags[0].Detail
			for _, want := range tc.wantInText {
				if !strings.Contains(text, want) {
					t.Errorf("diagnostic %q does not contain %q", text, want)
				}
			}
		})
	}
}

// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package tenant

import "testing"

func boolPtr(b bool) *bool { return &b }

func TestApiFailure(t *testing.T) {
	cases := []struct {
		name       string
		success    *bool
		props      map[string]any
		wantFailed bool
		wantDetail string
	}{
		{
			name:    "typed success true",
			success: boolPtr(true),
			props:   map[string]any{"msg": "ignored"},
		},
		{
			name:  "no success anywhere",
			props: map[string]any{"msg": "ignored"},
		},
		{
			name:  "untyped success true",
			props: map[string]any{"success": true},
		},
		{
			// The pre-8.1.0 "Invalid role id" shape from POST /api/accounts.
			name:    "typed success false with msg and errors",
			success: boolPtr(false),
			props: map[string]any{
				"msg":    "Error saving account",
				"errors": map[string]any{"role": "Invalid role id 99"},
			},
			wantFailed: true,
			wantDetail: "Error saving account; role: Invalid role id 99",
		},
		{
			// UpdateTenant's model has no Success field, so the envelope flag
			// arrives untyped in AdditionalProperties.
			name: "untyped success false",
			props: map[string]any{
				"success": false,
				"msg":     "Error saving account",
			},
			wantFailed: true,
			wantDetail: "Error saving account",
		},
		{
			name:    "errors are sorted by key for a stable message",
			success: boolPtr(false),
			props: map[string]any{
				"errors": map[string]any{
					"subdomain": "must be unique",
					"name":      "must be unique",
				},
			},
			wantFailed: true,
			wantDetail: "name: must be unique; subdomain: must be unique",
		},
		{
			name:       "success false with nothing else",
			success:    boolPtr(false),
			props:      map[string]any{},
			wantFailed: true,
			wantDetail: "the API reported success=false without a message",
		},
		{
			name:       "nil props",
			success:    boolPtr(false),
			wantFailed: true,
			wantDetail: "the API reported success=false without a message",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			detail, failed := apiFailure(tc.success, tc.props)

			if failed != tc.wantFailed {
				t.Fatalf("failed = %v, want %v (detail %q)", failed, tc.wantFailed, detail)
			}

			if detail != tc.wantDetail {
				t.Errorf("detail = %q, want %q", detail, tc.wantDetail)
			}
		})
	}
}

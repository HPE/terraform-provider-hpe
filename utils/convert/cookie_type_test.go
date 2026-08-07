// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package convert

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestCookieTypeToAPI(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    types.String
		expected *string
	}{
		{
			name:     "session",
			input:    types.StringValue("session"),
			expected: PtrString("LBSessionCookieTime"),
		},
		{
			name:     "persistence",
			input:    types.StringValue("persistence"),
			expected: PtrString("LBPersistenceCookieTime"),
		},
		{
			name:     "null returns nil",
			input:    types.StringNull(),
			expected: nil,
		},
		{
			name:     "unknown returns nil",
			input:    types.StringUnknown(),
			expected: nil,
		},
		{
			name:     "invalid returns nil",
			input:    types.StringValue("invalid"),
			expected: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := CookieTypeToAPI(tt.input)
			if tt.expected == nil {
				if result != nil {
					t.Errorf("expected nil, got %q", *result)
				}

				return
			}

			if result == nil {
				t.Fatalf("expected %q, got nil", *tt.expected)

				return
			}

			if *result != *tt.expected {
				t.Errorf("expected %q, got %q", *tt.expected, *result)
			}
		})
	}
}

func TestCookieTypeFromAPI(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    *string
		expected types.String
	}{
		{
			name:     "LBSessionCookieTime",
			input:    PtrString("LBSessionCookieTime"),
			expected: types.StringValue("session"),
		},
		{
			name:     "LBPersistenceCookieTime",
			input:    PtrString("LBPersistenceCookieTime"),
			expected: types.StringValue("persistence"),
		},
		{
			name:     "nil returns null",
			input:    nil,
			expected: types.StringNull(),
		},
		{
			name:     "unknown value returns null",
			input:    PtrString("unknown"),
			expected: types.StringNull(),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := CookieTypeFromAPI(tt.input)
			if !result.Equal(tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package convert

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ptr(s string) *string {
	return &s
}

func TestStrToType(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input    *string
		expected types.String
	}{
		"nil": {
			input:    nil,
			expected: types.StringNull(),
		},
		"empty": {
			input:    ptr(""),
			expected: types.StringValue(""),
		},
		"zero": {
			input:    ptr("0"),
			expected: types.StringValue("0"),
		},
		"value": {
			input:    ptr("hello"),
			expected: types.StringValue("hello"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := StrToType(tc.input)
			if !result.Equal(tc.expected) {
				t.Errorf("StrToType(%v) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestStrToTypeEmptyAsNull(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input    *string
		expected types.String
	}{
		"nil": {
			input:    nil,
			expected: types.StringNull(),
		},
		"empty string is null": {
			input:    ptr(""),
			expected: types.StringNull(),
		},
		"zero is not null": {
			input:    ptr("0"),
			expected: types.StringValue("0"),
		},
		"normal value": {
			input:    ptr("hello"),
			expected: types.StringValue("hello"),
		},
		"whitespace is not empty": {
			input:    ptr(" "),
			expected: types.StringValue(" "),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := StrToTypeEmptyAsNull(tc.input)
			if !result.Equal(tc.expected) {
				t.Errorf("StrToTypeEmptyAsNull(%v) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestStrToTypeZeroIDAsNull(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input    *string
		expected types.String
	}{
		"nil": {
			input:    nil,
			expected: types.StringNull(),
		},
		"empty string is null": {
			input:    ptr(""),
			expected: types.StringNull(),
		},
		"zero is null": {
			input:    ptr("0"),
			expected: types.StringNull(),
		},
		"valid ID": {
			input:    ptr("42"),
			expected: types.StringValue("42"),
		},
		"large ID": {
			input:    ptr("99999"),
			expected: types.StringValue("99999"),
		},
		"non-numeric string": {
			input:    ptr("abc"),
			expected: types.StringValue("abc"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := StrToTypeZeroIDAsNull(tc.input)
			if !result.Equal(tc.expected) {
				t.Errorf("StrToTypeZeroIDAsNull(%v) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}

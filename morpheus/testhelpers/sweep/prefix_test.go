// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package sweep

import "testing"

func TestResolveSweepPrefix(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"empty falls back to default", "", defaultTestResourcePrefix},
		{"whitespace only falls back to default", "   ", defaultTestResourcePrefix},
		{"tab and newline fall back to default", "\t\n", defaultTestResourcePrefix},
		{"override used verbatim", "Foo", "Foo"},
		{"override is trimmed", "  Foo  ", "Foo"},
		{"default passes through", defaultTestResourcePrefix, defaultTestResourcePrefix},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolveSweepPrefix(tt.in); got != tt.want {
				t.Errorf("resolveSweepPrefix(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

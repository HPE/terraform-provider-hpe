// (C) Copyright 2026 Hewlett Packard Enterprise Development LP
package auth_test

import (
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/auth"
)

func TestNormalizeBaseURL(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"no trailing slash", "https://morpheus.example.com", "https://morpheus.example.com"},
		{"single trailing slash", "https://morpheus.example.com/", "https://morpheus.example.com"},
		{"multiple trailing slashes", "https://morpheus.example.com///", "https://morpheus.example.com"},
		{"trailing slash with path", "https://morpheus.example.com/sub/", "https://morpheus.example.com/sub"},
		{"trailing whitespace after slash", "https://morpheus.example.com/ ", "https://morpheus.example.com"},
		{"leading and trailing whitespace", "  https://morpheus.example.com/  ", "https://morpheus.example.com"},
		{"whitespace only", "   ", ""},
		{"empty string", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := auth.NormalizeBaseURL(tt.in); got != tt.want {
				t.Errorf("NormalizeBaseURL(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

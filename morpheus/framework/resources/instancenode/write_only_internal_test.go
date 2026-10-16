// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instancenode

import (
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// TestUnitServerSSHPasswordHashExtractor covers the instance_node hash
// extractor: a present hash, a present-but-null hash, an unset field, a nil
// server and a wrong-typed response.
func TestUnitServerSSHPasswordHashExtractor(t *testing.T) {
	t.Parallel()

	withHash := func(s string) *sdk.GetHost200ResponseServer {
		return &sdk.GetHost200ResponseServer{SshPasswordHash: *sdk.NewNullableString(&s)}
	}

	nullHash := func() *sdk.GetHost200ResponseServer {
		return &sdk.GetHost200ResponseServer{SshPasswordHash: *sdk.NewNullableString(nil)}
	}

	cases := []struct {
		name        string
		resp        any
		wantHash    string
		wantPresent bool
	}{
		{"present with value", withHash("abc123"), "abc123", true},
		{"present but null is a comparable empty value", nullHash(), "", true},
		{"unset field: not present", &sdk.GetHost200ResponseServer{}, "", false},
		{"nil server: not present", (*sdk.GetHost200ResponseServer)(nil), "", false},
		{"wrong type: not present", "not a server", "", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			hash, present := serverSSHPasswordHashExtractor(tc.resp)
			if hash != tc.wantHash || present != tc.wantPresent {
				t.Fatalf("serverSSHPasswordHashExtractor() = (%q, %v), want (%q, %v)",
					hash, present, tc.wantHash, tc.wantPresent)
			}
		})
	}
}

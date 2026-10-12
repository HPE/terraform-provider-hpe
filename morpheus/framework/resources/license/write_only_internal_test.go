// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package license

import (
	"regexp"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
)

func TestLicenseKeyID(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"trims whitespace", "  ABCDEFGHIJ  ", "ABCDEFGH"},
		{"longer than 8 truncated", "ABCDEFGHIJKL", "ABCDEFGH"},
		{"exactly 8", "ABCDEFGH", "ABCDEFGH"},
		{"shorter than 8", "ABC", "ABC"},
		{"newline trimmed then truncated", "ABCDEFGHIJ\n", "ABCDEFGH"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := licenseKeyID(tc.in); got != tc.want {
				t.Fatalf("licenseKeyID(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestKeyDigest(t *testing.T) {
	// sha256("ABCDEFGHIJ") computed independently.
	const want = "261305762671a58cae5b74990bcfc236c2336fb04a0fbac626166d9491d2884c"

	if got := keyDigest("ABCDEFGHIJ"); got != want {
		t.Fatalf("keyDigest = %q, want %q", got, want)
	}

	// Whitespace-insensitive: surrounding whitespace is trimmed before hashing.
	if got := keyDigest("  ABCDEFGHIJ\n"); got != want {
		t.Fatalf("keyDigest with whitespace = %q, want %q", got, want)
	}

	// 64 lowercase hex characters.
	if !regexp.MustCompile(`^[0-9a-f]{64}$`).MatchString(keyDigest("anything")) {
		t.Fatalf("keyDigest is not 64 hex characters")
	}
}

func TestPlanDecision(t *testing.T) {
	const digest = "261305762671a58cae5b74990bcfc236c2336fb04a0fbac626166d9491d2884c"
	upper := "261305762671A58CAE5B74990BCFC236C2336FB04A0FBAC626166D9491D2884C"

	cases := []struct {
		name       string
		stored     string
		ok         bool
		configured string
		want       wodigest.Decision
	}{
		{"no stored digest seeds", "", false, digest, wodigest.Seed},
		{"equal is unchanged", digest, true, digest, wodigest.Unchanged},
		{"equal different case is unchanged", upper, true, digest, wodigest.Unchanged},
		{"different is update", digest, true, "deadbeef", wodigest.Update},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := planDecision(tc.stored, tc.ok, tc.configured); got != tc.want {
				t.Fatalf("planDecision(%q, %v, %q) = %s, want %s",
					tc.stored, tc.ok, tc.configured, got, tc.want)
			}
		})
	}
}

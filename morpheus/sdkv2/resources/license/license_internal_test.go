// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package license

import (
	"testing"

	morpheus "github.com/HPE/terraform-provider-hpe/internal/sdk/legacy"
)

// The identifiers Morpheus reports for an installed licence are derived from
// the key: hash = SHA-256 of the trimmed key, keyId = its first eight
// characters (ApplianceLicenseService). Both were checked against a live
// appliance (MORPH-18310).
func TestUnitLicenseKeyIdentifiers(t *testing.T) {
	t.Parallel()

	const key = "4a407f6dEXAMPLEKEYBODY"

	if got := licenseKeyID(key); got != "4a407f6d" {
		t.Errorf("licenseKeyID = %q, want 4a407f6d", got)
	}

	if got := licenseKeyID("short"); got != "short" {
		t.Errorf("licenseKeyID(short) = %q, want short", got)
	}

	plain := licenseKeyHash(key)
	if len(plain) != 64 {
		t.Fatalf("licenseKeyHash length = %d, want 64 hex characters", len(plain))
	}

	// Morpheus trims the key before hashing; so must the resource, or a key
	// pasted with a trailing newline would never match the installed one.
	if got := licenseKeyHash("  " + key + "\n"); got != plain {
		t.Errorf("licenseKeyHash is not whitespace-insensitive: %q != %q", got, plain)
	}

	if got := licenseKeyHash("other"); got == plain {
		t.Errorf("licenseKeyHash: different keys hashed alike")
	}
}

func TestUnitFindInstalledLicense(t *testing.T) {
	t.Parallel()

	active := morpheus.ApplianceLicense{KeyID: "aaaaaaaa", Hash: "HASH-A", Producttier: "capacity"}
	stacked := morpheus.ApplianceLicense{KeyID: "bbbbbbbb", Hash: "hash-b", Producttier: "enterprise"}
	result := &morpheus.GetLicenseResult{
		License:           &active,
		InstalledLicenses: []morpheus.ApplianceLicense{active, stacked},
	}

	cases := []struct {
		name     string
		hash     string
		keyID    string
		wantTier string // "" = not found
	}{
		{"active licence by hash (case-insensitive)", "hash-a", "", "capacity"},
		{"stacked licence by hash", "hash-b", "", "enterprise"},
		{"import: by key id when there is no key", "", "bbbbbbbb", "enterprise"},
		{"hash wins over key id", "hash-b", "aaaaaaaa", "enterprise"},
		{"unknown hash", "nope", "aaaaaaaa", ""},
		{"nothing to match on", "", "", ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := findInstalledLicense(result, tc.hash, tc.keyID)
			switch {
			case tc.wantTier == "" && got != nil:
				t.Fatalf("found %q, want none", got.KeyID)
			case tc.wantTier != "" && got == nil:
				t.Fatalf("found none, want %s", tc.wantTier)
			case got != nil && got.Producttier != tc.wantTier:
				t.Fatalf("found %s, want %s", got.Producttier, tc.wantTier)
			}
		})
	}

	// No active licence (appliance unlicensed): License is nil, the list may
	// still carry inactive entries or be empty.
	if got := findInstalledLicense(&morpheus.GetLicenseResult{}, "hash-a", "aaaaaaaa"); got != nil {
		t.Errorf("empty result: found %q, want none", got.KeyID)
	}

	if got := findInstalledLicense(nil, "hash-a", ""); got != nil {
		t.Errorf("nil result: found %q, want none", got.KeyID)
	}
}

// TestUnitKeyDiffSuppress: the only difference key suppresses is surrounding
// whitespace. An empty prior value -- what Read leaves when the configured key
// is not installed -- and a different key must both plan a change, or the
// resource could neither re-install a missing key nor install a new one.
func TestUnitKeyDiffSuppress(t *testing.T) {
	t.Parallel()

	const (
		k1 = "4a407f6dFIRSTKEY"
		k2 = "9b2c1d0eSECONDKEY"
	)

	cases := []struct {
		name     string
		old, new string
		suppress bool
	}{
		{"same key", k1, k1, true},
		{"same key, whitespace around the new value", k1, "  " + k1 + "\n", true},
		{"same key, whitespace around the old value", k1 + "\n", k1, true},
		{"different key", k1, k2, false},
		{"not installed (Read cleared the key): re-install", "", k1, false},
		{"not installed and a new key", "", k2, false},
		{"new resource (no prior value)", "", "", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := keyDiffSuppress("key", tc.old, tc.new, nil); got != tc.suppress {
				t.Errorf("keyDiffSuppress(old=%q, new=%q) = %v, want %v", tc.old, tc.new, got, tc.suppress)
			}
		})
	}
}

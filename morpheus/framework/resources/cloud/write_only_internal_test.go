// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cloud

import (
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// zoneWithConfig builds a cloud whose config map carries the given keys, using
// the generic map variant of the config anyOf (which MarshalJSON renders as the
// flat object the extractor reads).
func zoneWithConfig(m map[string]any) *sdk.GetClouds200ResponseZone {
	cfg := &sdk.GetClouds200ResponseZoneConfig{MapmapOfStringAny: &m}

	return &sdk.GetClouds200ResponseZone{Config: cfg}
}

// TestUnitZoneConfigHash covers the cloud hash extractor: a present key with a
// value, a present key that is null, an absent key, a nil config and a nil
// zone.
func TestUnitZoneConfigHash(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		zone        *sdk.GetClouds200ResponseZone
		key         string
		wantHash    string
		wantPresent bool
	}{
		{
			"present with value",
			zoneWithConfig(map[string]any{"passwordHash": "abc123"}),
			"passwordHash", "abc123", true,
		},
		{
			"present but null is a comparable empty value",
			zoneWithConfig(map[string]any{"passwordHash": nil}),
			"passwordHash", "", true,
		},
		{
			"absent key: nothing can be concluded",
			zoneWithConfig(map[string]any{"clientSecretHash": "xyz"}),
			"passwordHash", "", false,
		},
		{
			"azure client secret hash present",
			zoneWithConfig(map[string]any{"clientSecretHash": "secret-hash"}),
			"clientSecretHash", "secret-hash", true,
		},
		{
			"nil config: not present",
			&sdk.GetClouds200ResponseZone{}, "passwordHash", "", false,
		},
		{
			"nil zone: not present",
			nil, "passwordHash", "", false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			hash, present := zoneConfigHash(tc.zone, tc.key)
			if hash != tc.wantHash || present != tc.wantPresent {
				t.Fatalf("zoneConfigHash() = (%q, %v), want (%q, %v)",
					hash, present, tc.wantHash, tc.wantPresent)
			}
		})
	}
}

// TestUnitZoneConfigHashExtractorWrongType asserts the extractor returns
// not-present for a response that is not a *GetClouds200ResponseZone.
func TestUnitZoneConfigHashExtractorWrongType(t *testing.T) {
	t.Parallel()

	extract := zoneConfigHashExtractor("passwordHash")

	if _, present := extract("not a zone"); present {
		t.Fatal("a wrong-typed response must be not present")
	}

	if _, present := extract(zoneWithConfig(map[string]any{"passwordHash": "h"})); !present {
		t.Fatal("a matching zone must be present")
	}
}

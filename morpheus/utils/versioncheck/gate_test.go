// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package versioncheck_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/constants"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/versioncheck"
)

const feature = "Cloud affinity groups"

// mustParse is a test helper for building a Version from a build string.
func mustParse(t *testing.T, build string) *versioncheck.Version {
	t.Helper()

	v, err := versioncheck.Parse(build)
	if err != nil {
		t.Fatalf("Parse(%q) unexpected error: %v", build, err)
	}

	return v
}

// TestDecideBlocksBelowConstraint pins the whole point of the gate: an
// appliance below the constraint must be refused, and the diagnostic must name
// both the required version and the version the appliance actually reports.
func TestDecideBlocksBelowConstraint(t *testing.T) {
	ctx := context.Background()

	for _, build := range []string{"8.0.8", "8.0.9", "8.0.9.4", "7.0.4", "8.0.1"} {
		t.Run(build, func(t *testing.T) {
			diags := versioncheck.Decide(
				ctx, feature, constants.AffinityGroupMinVersion,
				mustParse(t, build), nil,
			)

			if !diags.HasError() {
				t.Fatalf("Decide(%q, %q) allowed the operation, want an error",
					build, constants.AffinityGroupMinVersion)
			}

			if len(diags) != 1 {
				t.Fatalf("Decide(%q) produced %d diagnostics, want 1", build, len(diags))
			}

			detail := diags[0].Detail()

			if !strings.Contains(detail, constants.AffinityGroupMinVersion) {
				t.Errorf("detail %q does not name the required version %q",
					detail, constants.AffinityGroupMinVersion)
			}

			if !strings.Contains(detail, build) {
				t.Errorf("detail %q does not name the appliance version %q", detail, build)
			}

			if !strings.Contains(diags[0].Summary(), feature) {
				t.Errorf("summary %q does not name the feature %q", diags[0].Summary(), feature)
			}
		})
	}
}

// TestDecideAllowsAtOrAboveConstraint covers the boundary — 8.0.10 itself must
// pass — plus the four-segment builds Morpheus actually reports, which a naive
// major.minor.patch comparison gets wrong.
func TestDecideAllowsAtOrAboveConstraint(t *testing.T) {
	ctx := context.Background()

	builds := []string{
		"8.0.10",     // exactly the boundary
		"8.0.10.1",   // four-segment build of the boundary release
		"8.0.11",     // later patch
		"8.1.0",      // later minor
		"9.0.2.18",   // much later, four segments
		"8.0.10-rc1", // prerelease suffix is stripped before comparing
		" 8.0.10 ",   // whitespace is trimmed
		"8.0.100",    // 100 > 10, not a string comparison
	}

	for _, build := range builds {
		t.Run(build, func(t *testing.T) {
			diags := versioncheck.Decide(
				ctx, feature, constants.AffinityGroupMinVersion,
				mustParse(t, build), nil,
			)

			if diags.HasError() {
				t.Errorf("Decide(%q, %q) blocked the operation: %v",
					build, constants.AffinityGroupMinVersion, diags)
			}
		})
	}
}

// TestDecideFailsOpenWhenVersionUnknown is the deliberate policy choice: if the
// appliance version cannot be read — most plausibly because the API token lacks
// the admin-health permission that GET /api/health requires — the gate steps
// aside rather than blocking a possibly healthy appliance.
func TestDecideFailsOpenWhenVersionUnknown(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		name      string
		version   *versioncheck.Version
		lookupErr error
	}{
		{
			name:      "lookup failed",
			version:   nil,
			lookupErr: errors.New("query appliance version: 403 Forbidden"),
		},
		{
			name:      "no version and no error",
			version:   nil,
			lookupErr: nil,
		},
		{
			name:      "error wins over a stale version value",
			version:   mustParse(t, "8.0.9"),
			lookupErr: errors.New("query appliance version: connection refused"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			diags := versioncheck.Decide(
				ctx, feature, constants.AffinityGroupMinVersion, tc.version, tc.lookupErr,
			)

			if len(diags) != 0 {
				t.Errorf("Decide fail-open case produced diagnostics, want none: %v", diags)
			}
		})
	}
}

// TestDecideFailsOpenOnMalformedConstraint: a bad constraint is a provider bug.
// Blocking the practitioner for our own typo helps nobody.
func TestDecideFailsOpenOnMalformedConstraint(t *testing.T) {
	diags := versioncheck.Decide(
		context.Background(), feature, "not-a-constraint", mustParse(t, "8.0.9"), nil,
	)

	if len(diags) != 0 {
		t.Errorf("Decide with a malformed constraint produced diagnostics, want none: %v", diags)
	}
}

// TestAffinityGroupMinVersionIsDocumentedValue guards the number quoted in the
// resource and data source doc templates ("Requires Morpheus 8.0.10 or later").
// If the gate moves, this fails and the templates must move with it.
func TestAffinityGroupMinVersionIsDocumentedValue(t *testing.T) {
	if constants.AffinityGroupMinVersion != ">= 8.0.10" {
		t.Errorf(
			"AffinityGroupMinVersion = %q, want %q; the affinity group doc "+
				"templates state 8.0.10 and must be updated together with this",
			constants.AffinityGroupMinVersion, ">= 8.0.10",
		)
	}
}

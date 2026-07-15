// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package versioncheck provides helpers for checking the Morpheus appliance
// version against version constraints. It wraps hashicorp/go-version and, unlike
// a naive major.minor.patch trim, preserves the full build version (for example
// 9.0.2.18) so that four-segment constraints compare correctly.
package versioncheck

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	goversion "github.com/hashicorp/go-version"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

// Version is a parsed, comparable version. It aliases hashicorp/go-version's
// Version so callers do not need to import go-version directly.
type Version = goversion.Version

// Appliance fetches the target Morpheus appliance's build version via the health
// API and parses it. An error is returned if the version cannot be retrieved or
// parsed; callers that treat the version as best-effort (for example plan-time
// checks) should skip their logic on error rather than fail.
func Appliance(ctx context.Context, client *sdk.APIClient) (*Version, error) {
	resp, hresp, err := client.HealthAPI.ListHealth(ctx).Execute()
	if err != nil || hresp == nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("query appliance version: %s", errfmt.ErrMsg(err, hresp))
	}
	if resp.Health == nil || resp.Health.BuildVersion == nil {
		return nil, fmt.Errorf("appliance health response did not include a build version")
	}

	return Parse(*resp.Health.BuildVersion)
}

// Parse parses a Morpheus build version (for example "9.0.2.18") into a
// comparable Version. It preserves the full dotted numeric portion — unlike a
// major.minor.patch trim — and drops any trailing prerelease/build/whitespace
// suffix that would otherwise skew comparisons.
func Parse(buildVersion string) (*Version, error) {
	v, err := goversion.NewVersion(clean(buildVersion))
	if err != nil {
		return nil, fmt.Errorf("parse appliance version %q: %w", buildVersion, err)
	}

	return v, nil
}

// Satisfies reports whether v satisfies the given constraint (for example
// ">= 9.0.2").
func Satisfies(v *Version, constraint string) (bool, error) {
	c, err := goversion.NewConstraint(constraint)
	if err != nil {
		return false, fmt.Errorf("parse version constraint %q: %w", constraint, err)
	}

	return c.Check(v), nil
}

// AtLeast fetches the appliance version and reports whether it satisfies the
// constraint. It is a convenience wrapper over Appliance and Satisfies.
func AtLeast(ctx context.Context, client *sdk.APIClient, constraint string) (bool, error) {
	v, err := Appliance(ctx, client)
	if err != nil {
		return false, err
	}

	return Satisfies(v, constraint)
}

// clean keeps the numeric dotted portion of a Morpheus build version and drops
// any trailing prerelease/build/whitespace suffix so it can be parsed and
// compared.
func clean(v string) string {
	v = strings.TrimSpace(v)
	if i := strings.IndexAny(v, "-+ "); i >= 0 {
		v = v[:i]
	}

	return v
}

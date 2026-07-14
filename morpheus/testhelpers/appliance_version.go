// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package testhelpers

import (
	"context"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/versioncheck"
)

// ApplianceVersionAtLeast reports whether the target Morpheus appliance's build
// version satisfies the given constraint (for example ">= 9.0.2"). It delegates
// to the shared morpheus/utils/version library and fails the test if the
// appliance version cannot be determined.
func ApplianceVersionAtLeast(ctx context.Context, t *testing.T, constraint string) bool {
	t.Helper()

	ok, err := versioncheck.AtLeast(ctx, newClient(ctx, t), constraint)
	if err != nil {
		t.Fatalf("failed to check appliance version: %v", err)
	}

	return ok
}

// SkipUnlessApplianceVersionAtLeast skips the test unless the target appliance's
// build version satisfies the constraint. Use it to gate assertions for fields
// that only exist on newer Morpheus versions.
func SkipUnlessApplianceVersionAtLeast(
	ctx context.Context,
	t *testing.T,
	constraint string,
) {
	t.Helper()

	if !ApplianceVersionAtLeast(ctx, t, constraint) {
		t.Skipf("appliance build version does not satisfy %q", constraint)
	}
}

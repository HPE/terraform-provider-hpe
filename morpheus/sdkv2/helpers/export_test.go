// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package helpers

import morpheus "github.com/HPE/terraform-provider-hpe/internal/sdk/legacy"

// ResetTenancyCache forgets every cached tenancy determination so a test can
// make the same client resolve again (for example against a whoami server that
// now answers differently). Test-only: the export_test mechanism keeps it out
// of the production build.
//
// Because the cache is a package-level global, tests that reset it must run
// serially. A test that calls t.Parallel() while touching CallerIsMaster would
// race this reset against another test's determination.
func ResetTenancyCache() {
	cacheMu.Lock()
	defer cacheMu.Unlock()

	callers = map[*morpheus.Client]caller{}
}

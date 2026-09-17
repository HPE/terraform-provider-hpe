// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package helpers

// ResetTenancyCache clears the process-wide tenancy determination so tests can
// exercise differently-credentialed callers in one process. Test-only: the
// export_test mechanism keeps it out of the production build.
//
// Because the cache is a package-level global, tests that reset it must run
// serially. A test that calls t.Parallel() while touching CallerIsMaster would
// race this reset against another test's determination.
func ResetTenancyCache() {
	cacheMu.Lock()
	defer cacheMu.Unlock()

	resolved = false
	isMaster = false
}

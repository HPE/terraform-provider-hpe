// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

//go:build !unix

package testhelpers

import "testing"

// LockFixture is a no-op where flock is unavailable (Windows): the acceptance
// tests that need the fixture lock are run from Unix hosts, and this stub
// only keeps the testhelpers package compiling everywhere.
func LockFixture(t *testing.T, name string) {
	t.Helper()
	t.Logf("fixture lock %q not supported on this platform; tests sharing the fixture must run serially", name)
}

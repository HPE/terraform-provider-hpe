// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

//go:build unix

package testhelpers

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"
)

// LockFixture serialises tests -- across packages, so across `go test`
// processes -- that share an appliance fixture which cannot be used by two
// tests at once. The lock is a flock on a file in the OS temp directory named
// after the fixture and is released when the test ends. Tests holding it must
// not also call t.Parallel against each other; the lock makes that moot.
//
// Used by the cloud affinity group tests: a VMware DRS cluster refuses the
// same VM in two KEEP_TOGETHER rules, so every group seeded with the fixture VM
// must be created and destroyed before the next.
func LockFixture(t *testing.T, name string) {
	t.Helper()

	path := filepath.Join(os.TempDir(), "tf-acc-fixture-"+name+".lock")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		t.Fatalf("fixture lock %s: %v", path, err)
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		_ = f.Close()
		t.Fatalf("fixture lock %s: %v", path, err)
	}
	t.Cleanup(func() {
		_ = syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
		_ = f.Close()
	})
}

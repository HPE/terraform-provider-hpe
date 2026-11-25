// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package virtualswitch

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

// TestUnitSweepPrefixMatches pins the sweeper's copy of the test-name prefix
// to TestNamePrefix. The sweeper is build-tagged and cannot import this
// package (cycle), so the literal is duplicated and checked here by reading
// the source.
func TestUnitSweepPrefixMatches(t *testing.T) {
	t.Parallel()

	src, err := os.ReadFile(filepath.Join("sweep", "sweep.go"))
	if err != nil {
		t.Fatal(err)
	}

	m := regexp.MustCompile(`const TestNamePrefix = "([^"]+)"`).FindSubmatch(src)
	if m == nil {
		t.Fatal("sweep/sweep.go no longer declares TestNamePrefix")
	}

	if got := string(m[1]); got != TestNamePrefix {
		t.Fatalf("sweeper prefix %q differs from TestNamePrefix %q", got, TestNamePrefix)
	}
}

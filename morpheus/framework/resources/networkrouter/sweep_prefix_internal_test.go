// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouter

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

// TestUnitSweepPrefixMatches pins the sweeper's copy of the SDN test-name
// prefix to SdnTestNamePrefix. The sweeper is build-tagged and cannot import
// this package (cycle), so the literal is duplicated and checked here by
// reading the source.
func TestUnitSweepPrefixMatches(t *testing.T) {
	t.Parallel()

	src, err := os.ReadFile(filepath.Join("sweep", "sweep.go"))
	if err != nil {
		t.Fatal(err)
	}

	m := regexp.MustCompile(`const SdnTestNamePrefix = "([^"]+)"`).FindSubmatch(src)
	if m == nil {
		t.Fatal("sweep/sweep.go no longer declares SdnTestNamePrefix")
	}

	if got := string(m[1]); got != SdnTestNamePrefix {
		t.Fatalf("sweeper prefix %q differs from SdnTestNamePrefix %q", got, SdnTestNamePrefix)
	}
}

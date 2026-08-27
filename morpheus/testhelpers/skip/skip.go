package skip

import (
	"os"
	"testing"
)

const (
	envRunSkippedByDefault = "RUN_SKIPPED_BY_DEFAULT"
)

// SkipByDefault returns true if the test should not run (i.e., the
// RUN_SKIPPED_BY_DEFAULT env var is not set). Callers should return
// early when true. The test will not be marked as skipped.
func SkipByDefault(t *testing.T) bool {
	t.Helper()

	_, run := os.LookupEnv(envRunSkippedByDefault)

	return !run
}

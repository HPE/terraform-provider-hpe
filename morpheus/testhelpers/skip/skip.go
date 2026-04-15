package skip

import (
	"os"
	"testing"
)

const (
	envRunSkippedByDefault = "RUN_SKIPPED_BY_DEFAULT"
)

// SkipByDefault will skip tests by default
// unless RUN_SKIPPED_BY_DEFAULT was set.
func SkipByDefault(t *testing.T) {
	if _, run := os.LookupEnv(envRunSkippedByDefault); !run {
		t.Skip()
	}
}

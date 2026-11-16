// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package virtualswitch

import (
	"os"
	"time"
)

// The fake appliance answers at once; the real pauses between refused
// attempts and sync polls would only slow the unit run. Set before any test
// starts (init runs before TestMain) -- but not for an acceptance run, which
// shares this test binary and must exercise the shipped cadence against the
// real appliance.
func init() {
	if os.Getenv("TF_ACC") != "" {
		return
	}

	refusalPause = 50 * time.Millisecond
	syncPollInterval = 50 * time.Millisecond
}

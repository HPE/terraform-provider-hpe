// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package instance

import (
	"errors"
	"testing"

	"github.com/cenkalti/backoff/v5"
)

// TestUnitCheckStatusDoneUpdate pins how the instance resize/update poller
// classifies each reported status against the real
// UpdateTargetStatuses/UpdateErrorStatuses sets.
//
// As with create, a resize can settle in "stopped" rather than "running" when
// the instance's cloud has "Automatically Power On VMs" (autoRecoverPowerState)
// disabled, or when resizing an already-stopped instance. The resize has still
// completed, so update polling must treat "stopped" as done (nil) rather than
// polling until the update timeout elapses.
func TestUnitCheckStatusDoneUpdate(t *testing.T) {
	t.Parallel()

	const (
		outcomeDone  = iota // resize finished successfully (nil error)
		outcomeError        // terminal failure (permanent error, stop polling)
		outcomeRetry        // not finished yet (retryable, keep polling)
	)

	tests := []struct {
		name    string
		status  string
		outcome int
	}{
		// Target statuses -> done.
		{"running is done", "running", outcomeDone},
		{"stopped is done", "stopped", outcomeDone},

		// Error statuses -> permanent failure.
		{"denied is error", "denied", outcomeError},
		{"cancelled is error", "cancelled", outcomeError},
		{"failed is error", "failed", outcomeError},
		{"suspended is error", "suspended", outcomeError},
		{"removing is error", "removing", outcomeError},
		{"pendingRemoval is error", "pendingRemoval", outcomeError},

		// In-progress statuses -> keep polling.
		{"resizing keeps polling", "resizing", outcomeRetry},
		{"starting keeps polling", "starting", outcomeRetry},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := checkStatusDone(tc.status, UpdateTargetStatuses, UpdateErrorStatuses)

			var permanent *backoff.PermanentError
			isPermanent := errors.As(err, &permanent)

			switch tc.outcome {
			case outcomeDone:
				if err != nil {
					t.Fatalf("status %q: expected done (nil error), got %v", tc.status, err)
				}
			case outcomeError:
				if err == nil {
					t.Fatalf("status %q: expected a permanent error, got nil", tc.status)
				}
				if !isPermanent {
					t.Fatalf("status %q: expected a permanent error, got retryable %v", tc.status, err)
				}
			case outcomeRetry:
				if err == nil {
					t.Fatalf("status %q: expected a retryable error, got nil", tc.status)
				}
				if isPermanent {
					t.Fatalf("status %q: expected a retryable error, got permanent %v", tc.status, err)
				}
			}
		})
	}
}

// TestUnitUpdateStatusSetsDisjoint guards the invariant that no status is in both
// the target and error sets. checkStatusDone evaluates error statuses first, so a
// status appearing in both would be silently treated as a failure.
func TestUnitUpdateStatusSetsDisjoint(t *testing.T) {
	t.Parallel()

	errorSet := make(map[string]struct{}, len(UpdateErrorStatuses))
	for _, s := range UpdateErrorStatuses {
		errorSet[s] = struct{}{}
	}

	for _, s := range UpdateTargetStatuses {
		if _, clash := errorSet[s]; clash {
			t.Errorf("status %q is in both UpdateTargetStatuses and UpdateErrorStatuses", s)
		}
	}
}

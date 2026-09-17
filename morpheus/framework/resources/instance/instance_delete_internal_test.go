// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instance

import (
	"errors"
	"testing"

	"github.com/cenkalti/backoff/v5"
)

// TestUnitCheckDeleteStatusDone pins which instance statuses end the wait after
// a DELETE, and which keep it polling. "stopped" and "suspended" must keep
// polling: Morpheus's status sync writes them onto a removing instance while
// its containers are stopped, so they are transient during a normal teardown
// (MORPH-4733). "warning" must end it: that is what Morpheus sets when the
// removal job fails.
func TestUnitCheckDeleteStatusDone(t *testing.T) {
	t.Parallel()

	const (
		outcomeError = iota // terminal failure (permanent error, stop polling)
		outcomeRetry        // not gone yet (retryable, keep polling)
	)

	tests := []struct {
		name    string
		status  string
		outcome int
	}{
		// Error statuses -> permanent failure.
		{"denied is error", "denied", outcomeError},
		{"cancelled is error", "cancelled", outcomeError},
		{"failed is error", "failed", outcomeError},
		{"warning is error", "warning", outcomeError},
		{"restoring is error", "restoring", outcomeError},

		// Transient during teardown -> keep polling.
		{"removing keeps polling", "removing", outcomeRetry},
		{"stopping keeps polling", "stopping", outcomeRetry},
		{"stopped keeps polling", "stopped", outcomeRetry},
		{"suspended keeps polling", "suspended", outcomeRetry},
		{"running keeps polling", "running", outcomeRetry},
		{"unknown keeps polling", "unknown", outcomeRetry},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := checkDeleteStatusDone(tc.status, nil)

			var permanent *backoff.PermanentError
			isPermanent := errors.As(err, &permanent)

			switch tc.outcome {
			case outcomeError:
				if !isPermanent {
					t.Errorf("status %q: want a permanent error, got %v", tc.status, err)
				}
			case outcomeRetry:
				if err == nil || isPermanent {
					t.Errorf("status %q: want a retryable error, got %v", tc.status, err)
				}
			}
		})
	}
}

// TestUnitCheckDeleteStatusDoneCarriesStatusMessage pins the failure text: the
// status, plus Morpheus's statusMessage when there is one, since for a failed
// removal that message is the only record of why.
func TestUnitCheckDeleteStatusDoneCarriesStatusMessage(t *testing.T) {
	t.Parallel()

	str := func(s string) *string { return &s }

	tests := []struct {
		name          string
		status        string
		statusMessage *string
		want          string
	}{
		{
			"with message",
			"warning", str("Unable to remove instance: Failed to delete VM"),
			"reached error status: warning (Unable to remove instance: Failed to delete VM)",
		},
		{"nil message", "failed", nil, "reached error status: failed"},
		{"empty message", "failed", str(""), "reached error status: failed"},
		{"blank message", "failed", str("   "), "reached error status: failed"},
		{"message is trimmed", "denied", str("  denied by policy  "), "reached error status: denied (denied by policy)"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := checkDeleteStatusDone(tc.status, tc.statusMessage)
			if err == nil {
				t.Fatalf("status %q: want an error", tc.status)
			}

			var permanent *backoff.PermanentError
			if !errors.As(err, &permanent) {
				t.Fatalf("status %q: want a permanent error, got %v", tc.status, err)
			}

			if got := permanent.Unwrap().Error(); got != tc.want {
				t.Errorf("message\n want: %s\n  got: %s", tc.want, got)
			}
		})
	}
}

// TestUnitDeleteStatusSetsDisjoint guards the two invariants the delete wait
// depends on: the transient teardown statuses are never classed as errors, and
// the removal-failure status always is.
func TestUnitDeleteStatusSetsDisjoint(t *testing.T) {
	t.Parallel()

	for _, transient := range []string{"removing", "stopping", "stopped", "suspended"} {
		for _, e := range DeleteErrorStatuses {
			if e == transient {
				t.Errorf("%q is transient during teardown and must not be in DeleteErrorStatuses", transient)
			}
		}
	}

	found := false

	for _, e := range DeleteErrorStatuses {
		if e == "warning" {
			found = true
		}
	}

	if !found {
		t.Error(`"warning" is the status Morpheus sets when a removal fails and must be in DeleteErrorStatuses`)
	}
}

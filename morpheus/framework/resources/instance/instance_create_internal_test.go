// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package instance

import (
	"errors"
	"strings"
	"testing"

	"github.com/cenkalti/backoff/v5"
)

// TestUnitCheckStatusDoneCreate pins how the instance create poller classifies
// each reported status against the real CreateTargetStatuses/CreateErrorStatuses
// sets.
//
// The key case is "stopped": an instance can settle in "stopped" rather than
// "running" when its cloud has "Automatically Power On VMs"
// (autoRecoverPowerState) disabled - which is the API default. Provisioning has
// still succeeded, so create polling must treat "stopped" as done (nil), not as
// a permanent error. Genuine provisioning failures must remain permanent, and
// in-progress statuses must keep polling (a retryable, non-permanent error).
func TestUnitCheckStatusDoneCreate(t *testing.T) {
	t.Parallel()

	const (
		outcomeDone  = iota // provisioning finished successfully (nil error)
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
		{"provisioning keeps polling", "provisioning", outcomeRetry},
		{"starting keeps polling", "starting", outcomeRetry},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := checkStatusDone(tc.status, CreateTargetStatuses, CreateErrorStatuses)

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

// TestUnitCreateStatusSetsDisjoint guards the invariant that no status is in both
// the target and error sets. checkStatusDone evaluates error statuses first, so a
// status appearing in both would be silently treated as a failure - which is
// exactly the "stopped" regression this fix corrects.
func TestUnitCreateStatusSetsDisjoint(t *testing.T) {
	t.Parallel()

	errorSet := make(map[string]struct{}, len(CreateErrorStatuses))
	for _, s := range CreateErrorStatuses {
		errorSet[s] = struct{}{}
	}

	for _, s := range CreateTargetStatuses {
		if _, clash := errorSet[s]; clash {
			t.Errorf("status %q is in both CreateTargetStatuses and CreateErrorStatuses", s)
		}
	}
}

// TestUnitServerUUIDApplied verifies no diagnostic when the requested UUID
// is found among the assigned ones.
func TestUnitServerUUIDApplied(t *testing.T) {
	t.Parallel()

	assigned := map[string]struct{}{
		"uuid-a":     {},
		"uuid-extra": {},
	}

	diags := validateServerUUIDLogic("uuid-a", assigned, 42)
	if diags.HasError() {
		t.Fatalf("expected no diagnostics, got: %v", diags)
	}
}

// TestUnitServerUUIDNotApplied verifies that a requested UUID not found
// produces an error diagnostic naming that UUID.
func TestUnitServerUUIDNotApplied(t *testing.T) {
	t.Parallel()

	assigned := map[string]struct{}{
		"uuid-other": {},
	}

	diags := validateServerUUIDLogic("uuid-missing", assigned, 99)
	if !diags.HasError() {
		t.Fatal("expected an error diagnostic when a UUID is missing")
	}

	detail := diags.Errors()[0].Detail()
	if !strings.Contains(detail, "uuid-missing") {
		t.Errorf("expected error to name 'uuid-missing', got: %s", detail)
	}
}

// TestUnitCheckServerUUIDsAllPresent verifies that when all requested UUIDs
// are found among the assigned ones, no diagnostic is emitted.
func TestUnitCheckServerUUIDsAllPresent(t *testing.T) {
	t.Parallel()

	requested := []string{"uuid-a", "uuid-b"}
	assigned := map[string]struct{}{
		"uuid-a": {},
		"uuid-b": {},
		"uuid-c": {}, // extra assigned UUID (from another container) is fine
	}

	diags := checkServerUUIDs(requested, assigned, 42)
	if diags.HasError() {
		t.Fatalf("expected no diagnostics, got: %v", diags)
	}
}

// TestUnitCheckServerUUIDsOneMissing verifies that a requested UUID not found
// among the assigned ones produces an error diagnostic naming that UUID.
func TestUnitCheckServerUUIDsOneMissing(t *testing.T) {
	t.Parallel()

	requested := []string{"uuid-a", "uuid-missing"}
	assigned := map[string]struct{}{
		"uuid-a":     {},
		"uuid-other": {},
	}

	diags := checkServerUUIDs(requested, assigned, 99)
	if !diags.HasError() {
		t.Fatal("expected an error diagnostic when a UUID is missing")
	}

	detail := diags.Errors()[0].Detail()
	if !strings.Contains(detail, "uuid-missing") {
		t.Errorf("expected error to name 'uuid-missing', got: %s", detail)
	}
	if strings.Contains(detail, "uuid-a") {
		t.Errorf("error should not name 'uuid-a' which was applied, got: %s", detail)
	}
}

// TestUnitCheckServerUUIDsExcessRequested verifies that requesting more UUIDs
// than servers provisioned produces an error for the excess UUIDs.
func TestUnitCheckServerUUIDsExcessRequested(t *testing.T) {
	t.Parallel()

	// Only one server was provisioned (layout_size=1), but two UUIDs requested.
	requested := []string{"uuid-a", "uuid-excess"}
	assigned := map[string]struct{}{
		"uuid-a": {},
	}

	diags := checkServerUUIDs(requested, assigned, 10)
	if !diags.HasError() {
		t.Fatal("expected an error diagnostic for excess UUIDs")
	}

	detail := diags.Errors()[0].Detail()
	if !strings.Contains(detail, "uuid-excess") {
		t.Errorf("expected error to name 'uuid-excess', got: %s", detail)
	}
}

// TestUnitCheckServerUUIDsEmptyRequested verifies that an empty requested list
// produces no diagnostics (edge case: set is non-null but empty).
func TestUnitCheckServerUUIDsEmptyRequested(t *testing.T) {
	t.Parallel()

	diags := checkServerUUIDs(nil, map[string]struct{}{"uuid-a": {}}, 1)
	if diags.HasError() {
		t.Fatalf("expected no diagnostics for empty requested list, got: %v", diags)
	}
}

// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package testhelpers

import (
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

// TestRecordResultRecordsCapabilitySkip protects the contract that an
// acceptance test skipped because a required capability is missing is recorded
// as "Skipped". It exercises the canonical pattern: register RecordResult
// first, then gate on capabilities.MustHaveOrSkip (which skips the test via t.Skip
// when the capability is absent). The deferred RecordResult must observe the
// skip and record it.
//
// This test is sequential; it briefly enables result recording and clears the
// capability registry, restoring both on cleanup. Parallel tests in this
// package only resume after the sequential phase, so the toggles do not race.
func TestRecordResultRecordsCapabilitySkip(t *testing.T) {
	t.Setenv(capabilities.EnvCapabilities, "")
	capabilities.ResetForTesting()
	t.Cleanup(capabilities.ResetForTesting)

	prevRecording := recordingEnabled
	recordingEnabled = true
	t.Cleanup(func() { recordingEnabled = prevRecording })

	const sub = "missing_capability"
	t.Run(sub, func(st *testing.T) {
		defer RecordResult(st)

		capabilities.MustHaveOrSkip(st, capabilities.Capability("nonexistent_capability"))

		st.Fatal("test body must not run when a required capability is missing")
	})

	name := t.Name() + "/" + sub
	m.Lock()
	got, ok := testResults[name]
	m.Unlock()

	if !ok {
		t.Fatalf("expected a recorded result for %q, found none", name)
	}
	if got.Status != "Skipped" {
		t.Fatalf("expected status %q, got %q", "Skipped", got.Status)
	}
}

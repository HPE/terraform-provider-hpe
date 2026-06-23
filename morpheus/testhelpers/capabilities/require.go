// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package capabilities

import (
	"fmt"
	"testing"
)

// Missing reports whether ANY of the required capabilities are missing.
// Use it at the top of a test function. When a required capability is
// unavailable the test is marked skipped (via t.Skip); register
// "defer testhelpers.RecordResult(t)" beforehand so the skip is recorded.
// The bool result is retained for call-site readability, but when a
// capability is missing the call does not return: the test is skipped.
//
// When TF_ACC_CAPABILITIES_VERBOSE is set, logs which capability is missing.
//
// Example:
//
//	func TestAccMorpheusVMwareCloud(t *testing.T) {
//	    defer testhelpers.RecordResult(t)
//	    if capabilities.Missing(t, capabilities.VMware) {
//	        t.Log("Skipping test due to missing capabilities")
//	    }
//	    t.Parallel()
//	    // ... rest of test
//	}
func Missing(t *testing.T, required ...Capability) bool {
	t.Helper()
	for _, cap := range required {
		if !Has(cap) {
			if IsVerbose() {
				t.Logf("capability %q not available, test not running", cap)
			}

			t.Skipf("required capability %q not available", cap)
		}
	}

	return false
}

// MissingAll returns true if ALL of the specified capabilities are missing.
// Use when a test can run with any one of several alternative capabilities.
// When all are missing the test is marked skipped (via t.Skip); register
// "defer testhelpers.RecordResult(t)" beforehand so the skip is recorded.
//
// Example:
//
//	func TestAccMorpheusRouterGeneric(t *testing.T) {
//	    defer testhelpers.RecordResult(t)
//	    // Runs if NSXT OR NSXV is available
//	    if capabilities.MissingAll(t, capabilities.NSXT, capabilities.NSXV) {
//	        t.Log("Skipping test due to missing capabilities")
//	    }
//	    // ... rest of test
//	}
func MissingAll(t *testing.T, anyOf ...Capability) bool {
	t.Helper()
	if HasAny(anyOf...) {
		return false
	}
	if IsVerbose() {
		t.Logf("none of capabilities %v available, test not running", capabilityNames(anyOf))
	}

	t.Skipf("none of required capabilities %v available", capabilityNames(anyOf))

	return true
}

// capabilityNames returns a formatted string of capability names for logging.
func capabilityNames(caps []Capability) string {
	if len(caps) == 0 {
		return "[]"
	}
	names := make([]string, len(caps))
	for i, cap := range caps {
		names[i] = string(cap)
	}

	return fmt.Sprintf("%v", names)
}

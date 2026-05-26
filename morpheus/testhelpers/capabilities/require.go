// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package capabilities

import (
	"fmt"
	"testing"
)

// Missing returns true if ANY of the required capabilities are missing.
// Use at the top of a test function; if true, return immediately.
// The test will silently return without being marked as skipped.
//
// When TF_ACC_CAPABILITIES_VERBOSE is set, logs which capability is missing.
//
// Example:
//
//	func TestAccMorpheusVMwareCloud(t *testing.T) {
//	    if capabilities.Missing(t, capabilities.MorpheusCore, capabilities.VMware) {
//	        return
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

			return true
		}
	}

	return false
}

// MissingAll returns true if ALL of the specified capabilities are missing.
// Use when a test can run with any one of several alternative capabilities.
//
// Example:
//
//	func TestAccMorpheusRouterGeneric(t *testing.T) {
//	    // Runs if NSXT OR NSXV is available (in addition to MorpheusCore)
//	    if capabilities.Missing(t, capabilities.MorpheusCore) {
//	        return
//	    }
//	    if capabilities.MissingAll(t, capabilities.NSXT, capabilities.NSXV) {
//	        return
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

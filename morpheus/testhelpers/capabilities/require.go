// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package capabilities

import (
	"strings"
	"testing"
)

// MustHave skips the test unless ALL of the required capabilities are
// available. Use it as the first capability gate in a test function (after
// "defer testhelpers.RecordResult(t)"). When one or more required
// capabilities are missing, the test is skipped via t.Skip with a message
// naming only the missing capabilities; the call does not return.
//
// The skip reason is emitted into "go test -json" as an output event for the
// test, in the stable, human-readable form:
//
//	missing required capabilities: a, b, c
//
// When TF_ACC_CAPABILITIES_VERBOSE is set, the same reason is also logged.
//
// Example:
//
//	func TestAccMorpheusVMwareCloud(t *testing.T) {
//	    defer testhelpers.RecordResult(t)
//	    capabilities.MustHave(t, capabilities.VMware)
//	    t.Parallel()
//	    // ... rest of test
//	}
func MustHave(t *testing.T, required ...Capability) {
	t.Helper()

	missing := missingCapabilities(required...)
	if len(missing) == 0 {
		return
	}

	msg := "missing required capabilities: " + capabilityNames(missing)
	if IsVerbose() {
		t.Logf("%s; test not running", msg)
	}

	// t.Skip (not Skipf): msg is pre-built and must not be treated as a format
	// string. This emits msg as an output event in "go test -json".
	t.Skip(msg)
}

// missingCapabilities returns the subset of required capabilities that are not
// available, preserving the order in which they were requested. It is pure
// (no *testing.T side effects), which makes it straightforward to unit test.
func missingCapabilities(required ...Capability) []Capability {
	var missing []Capability
	for _, capability := range required {
		if !Has(capability) {
			missing = append(missing, capability)
		}
	}

	return missing
}

// capabilityNames returns a human-readable, comma-separated list of capability
// names (e.g. "azure, network_dhcp") for use in skip messages and logs.
func capabilityNames(caps []Capability) string {
	names := make([]string, len(caps))
	for i, capability := range caps {
		names[i] = string(capability)
	}

	return strings.Join(names, ", ")
}

// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package capabilities

import (
	"os"
	"strings"
	"sync"
)

const (
	// EnvCapabilities is the environment variable that specifies available capabilities.
	// Format: comma-separated list of capability names.
	// Example: TF_ACC_CAPABILITIES="morpheus_core,vmware,nsxt"
	EnvCapabilities = "TF_ACC_CAPABILITIES"

	// EnvCapabilitiesVerbose enables verbose logging when tests don't run
	// due to missing capabilities.
	EnvCapabilitiesVerbose = "TF_ACC_CAPABILITIES_VERBOSE"
)

var (
	registry     map[Capability]bool
	registryOnce sync.Once
	verbose      bool
)

// loadRegistry loads the capability registry from environment variables.
// This is called once on first access via sync.Once.
func loadRegistry() {
	registry = make(map[Capability]bool)
	verbose = os.Getenv(EnvCapabilitiesVerbose) != ""

	env := os.Getenv(EnvCapabilities)
	if env == "" {
		// Safe default: if TF_ACC_CAPABILITIES is not set,
		// no capabilities are available and no tests run.
		// Tests can be destructive, so explicit opt-in is required.
		return
	}

	// Parse comma-separated capabilities
	for _, cap := range strings.Split(env, ",") {
		cap = strings.TrimSpace(cap)
		if cap != "" {
			registry[Capability(cap)] = true
		}
	}
}

// Has returns true if the specified capability is available.
// If TF_ACC_CAPABILITIES is not set, no capabilities are available.
func Has(cap Capability) bool {
	registryOnce.Do(loadRegistry)

	return registry[cap]
}

// HasAll returns true if all specified capabilities are available.
func HasAll(caps ...Capability) bool {
	for _, cap := range caps {
		if !Has(cap) {
			return false
		}
	}

	return true
}

// HasAny returns true if at least one of the specified capabilities is available.
func HasAny(caps ...Capability) bool {
	for _, cap := range caps {
		if Has(cap) {
			return true
		}
	}

	return len(caps) == 0 // empty list = true
}

// IsVerbose returns true if verbose capability logging is enabled.
func IsVerbose() bool {
	registryOnce.Do(loadRegistry)

	return verbose
}

// ResetForTesting resets the registry for testing purposes.
// This should only be used in unit tests.
func ResetForTesting() {
	registry = nil
	registryOnce = sync.Once{}
	verbose = false
}

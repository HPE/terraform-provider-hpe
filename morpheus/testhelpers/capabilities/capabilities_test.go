// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package capabilities

import (
	"os"
	"slices"
	"testing"
)

func TestHas_WithEnvSet(t *testing.T) {
	// Reset and set up test environment
	ResetForTesting()
	os.Setenv(EnvCapabilities, "vmware,nsxt,network_dhcp")
	defer os.Unsetenv(EnvCapabilities)
	defer ResetForTesting()

	tests := []struct {
		cap      Capability
		expected bool
	}{
		{VMware, true},
		{NSXT, true},
		{NetworkDHCP, true},
		{AWS, false},
		{Azure, false},
	}

	for _, tc := range tests {
		t.Run(string(tc.cap), func(t *testing.T) {
			got := Has(tc.cap)
			if got != tc.expected {
				t.Errorf("Has(%q) = %v, want %v", tc.cap, got, tc.expected)
			}
		})
	}
}

func TestHas_SafeDefault(t *testing.T) {
	// Reset and ensure env is not set
	ResetForTesting()
	os.Unsetenv(EnvCapabilities)
	defer ResetForTesting()

	// Safe default: when env not set, no capabilities are available
	caps := []Capability{
		VMware,
		AWS,
		Azure,
		NSXT,
		NSXV,
		Kubernetes,
	}

	for _, cap := range caps {
		t.Run(string(cap), func(t *testing.T) {
			if Has(cap) {
				t.Errorf("Has(%q) = true with env unset, want false (safe default)", cap)
			}
		})
	}
}

func TestHas_EmptyEnv(t *testing.T) {
	// Reset and set empty env
	ResetForTesting()
	os.Setenv(EnvCapabilities, "")
	defer os.Unsetenv(EnvCapabilities)
	defer ResetForTesting()

	// Empty string means no capabilities available (safe default)
	if Has(VMware) {
		t.Error("Has(VMware) = true with empty env, want false (safe default)")
	}
}

func TestHas_WhitespaceHandling(t *testing.T) {
	// Reset and set env with whitespace
	ResetForTesting()
	os.Setenv(EnvCapabilities, "  vmware , nsxt  ,  aws  ")
	defer os.Unsetenv(EnvCapabilities)
	defer ResetForTesting()

	tests := []struct {
		cap      Capability
		expected bool
	}{
		{VMware, true},
		{NSXT, true},
		{AWS, true},
		{Azure, false},
	}

	for _, tc := range tests {
		t.Run(string(tc.cap), func(t *testing.T) {
			got := Has(tc.cap)
			if got != tc.expected {
				t.Errorf("Has(%q) = %v, want %v", tc.cap, got, tc.expected)
			}
		})
	}
}

func TestHasAll(t *testing.T) {
	ResetForTesting()
	os.Setenv(EnvCapabilities, "vmware,nsxt")
	defer os.Unsetenv(EnvCapabilities)
	defer ResetForTesting()

	tests := []struct {
		name     string
		caps     []Capability
		expected bool
	}{
		{"all present", []Capability{VMware, NSXT}, true},
		{"one missing", []Capability{VMware, AWS}, false},
		{"all missing", []Capability{AWS, Azure}, false},
		{"empty list", []Capability{}, true},
		{"single present", []Capability{NSXT}, true},
		{"single missing", []Capability{Azure}, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := HasAll(tc.caps...)
			if got != tc.expected {
				t.Errorf("HasAll(%v) = %v, want %v", tc.caps, got, tc.expected)
			}
		})
	}
}

func TestHasAny(t *testing.T) {
	ResetForTesting()
	os.Setenv(EnvCapabilities, "nsxt,vmware")
	defer os.Unsetenv(EnvCapabilities)
	defer ResetForTesting()

	tests := []struct {
		name     string
		caps     []Capability
		expected bool
	}{
		{"one present", []Capability{NSXT, AWS}, true},
		{"all present", []Capability{VMware, NSXT}, true},
		{"none present", []Capability{AWS, Azure}, false},
		{"empty list", []Capability{}, true},
		{"single present", []Capability{NSXT}, true},
		{"single missing", []Capability{GCP}, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := HasAny(tc.caps...)
			if got != tc.expected {
				t.Errorf("HasAny(%v) = %v, want %v", tc.caps, got, tc.expected)
			}
		})
	}
}

func TestMissingCapabilities(t *testing.T) {
	ResetForTesting()
	os.Setenv(EnvCapabilities, "vmware,nsxt")
	defer os.Unsetenv(EnvCapabilities)
	defer ResetForTesting()

	tests := []struct {
		name     string
		caps     []Capability
		expected []Capability
	}{
		{"all present", []Capability{VMware, NSXT}, nil},
		{"one missing", []Capability{VMware, AWS}, []Capability{AWS}},
		{"all missing", []Capability{AWS, Azure}, []Capability{AWS, Azure}},
		{"empty list", []Capability{}, nil},
		{"single present", []Capability{VMware}, nil},
		{"single missing", []Capability{GCP}, []Capability{GCP}},
		{
			"order preserved, only missing reported",
			[]Capability{Azure, VMware, AWS},
			[]Capability{Azure, AWS},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := missingCapabilities(tc.caps...)
			if !slices.Equal(got, tc.expected) {
				t.Errorf("missingCapabilities(%v) = %v, want %v", tc.caps, got, tc.expected)
			}
		})
	}
}

func TestCapabilityNames(t *testing.T) {
	tests := []struct {
		name     string
		caps     []Capability
		expected string
	}{
		{"empty", nil, ""},
		{"single", []Capability{Azure}, "azure"},
		{"multiple", []Capability{Azure, NetworkDHCP}, "azure, network_dhcp"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := capabilityNames(tc.caps); got != tc.expected {
				t.Errorf("capabilityNames(%v) = %q, want %q", tc.caps, got, tc.expected)
			}
		})
	}
}

func TestMustHaveOrSkip(t *testing.T) {
	ResetForTesting()
	os.Setenv(EnvCapabilities, "vmware,nsxt")
	defer os.Unsetenv(EnvCapabilities)
	defer ResetForTesting()

	t.Run("all present runs the test body", func(st *testing.T) {
		var bodyRan bool
		st.Cleanup(func() {
			if st.Skipped() {
				t.Error("MustHaveOrSkip skipped when all required capabilities were present")
			}
			if !bodyRan {
				t.Error("test body did not run when all required capabilities were present")
			}
		})

		MustHaveOrSkip(st, VMware, NSXT)
		bodyRan = true
	})

	t.Run("missing capability skips before the body", func(st *testing.T) {
		var skipped bool
		// Registered first => runs last (LIFO), so it observes the value set
		// by the cleanup below after MustHaveOrSkip skips via runtime.Goexit.
		st.Cleanup(func() {
			if !skipped {
				t.Error("expected MustHaveOrSkip to mark the test as skipped")
			}
		})
		st.Cleanup(func() { skipped = st.Skipped() })

		MustHaveOrSkip(st, AWS) // not in the registry -> skips here
		t.Error("code after MustHaveOrSkip ran despite a missing capability")
	})
}

func TestIsVerbose(t *testing.T) {
	// Test verbose mode off
	ResetForTesting()
	os.Unsetenv(EnvCapabilitiesVerbose)
	os.Setenv(EnvCapabilities, "vmware")
	defer os.Unsetenv(EnvCapabilities)

	if IsVerbose() {
		t.Error("IsVerbose() = true without env set, want false")
	}

	// Test verbose mode on
	ResetForTesting()
	os.Setenv(EnvCapabilitiesVerbose, "1")
	os.Setenv(EnvCapabilities, "vmware")
	defer os.Unsetenv(EnvCapabilitiesVerbose)

	if !IsVerbose() {
		t.Error("IsVerbose() = false with env set, want true")
	}
}

func TestCapabilityString(t *testing.T) {
	tests := []struct {
		cap      Capability
		expected string
	}{
		{VMware, "vmware"},
		{NSXT, "nsxt"},
		{NetworkDHCP, "network_dhcp"},
	}

	for _, tc := range tests {
		t.Run(tc.expected, func(t *testing.T) {
			got := tc.cap.String()
			if got != tc.expected {
				t.Errorf("Capability.String() = %q, want %q", got, tc.expected)
			}
		})
	}
}

// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package capabilities

import (
	"os"
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

func TestMissing(t *testing.T) {
	ResetForTesting()
	os.Setenv(EnvCapabilities, "vmware,nsxt")
	defer os.Unsetenv(EnvCapabilities)
	defer ResetForTesting()

	tests := []struct {
		name     string
		caps     []Capability
		expected bool
	}{
		{"all present", []Capability{VMware, NSXT}, false},
		{"one missing", []Capability{VMware, AWS}, true},
		{"all missing", []Capability{AWS, Azure}, true},
		{"empty list", []Capability{}, false},
		{"single present", []Capability{VMware}, false},
		{"single missing", []Capability{GCP}, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Missing(t, tc.caps...)
			if got != tc.expected {
				t.Errorf("Missing(%v) = %v, want %v", tc.caps, got, tc.expected)
			}
		})
	}
}

func TestMissingAll(t *testing.T) {
	ResetForTesting()
	os.Setenv(EnvCapabilities, "nsxt,vmware")
	defer os.Unsetenv(EnvCapabilities)
	defer ResetForTesting()

	tests := []struct {
		name     string
		caps     []Capability
		expected bool
	}{
		{"one present (NSXT)", []Capability{NSXT, NSXV}, false},
		{"none present", []Capability{AWS, Azure}, true},
		{"all present", []Capability{VMware, NSXT}, false},
		{"empty list", []Capability{}, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := MissingAll(t, tc.caps...)
			if got != tc.expected {
				t.Errorf("MissingAll(%v) = %v, want %v", tc.caps, got, tc.expected)
			}
		})
	}
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

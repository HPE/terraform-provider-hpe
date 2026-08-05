// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instancenode

import (
	"testing"
)

func TestIPReady_EmptyString(t *testing.T) {
	t.Parallel()

	if IPReady("") {
		t.Error("expected IPReady=false for empty string")
	}
}

func TestIPReady_ZeroAddress(t *testing.T) {
	t.Parallel()

	if IPReady("0.0.0.0") {
		t.Error("expected IPReady=false for 0.0.0.0")
	}
}

func TestIPReady_ValidIPv4(t *testing.T) {
	t.Parallel()

	if !IPReady("10.0.1.5") {
		t.Error("expected IPReady=true for 10.0.1.5")
	}
}

func TestIPReady_ValidIPv6(t *testing.T) {
	t.Parallel()

	if !IPReady("fe80::1") {
		t.Error("expected IPReady=true for fe80::1")
	}
}

func TestIPReady_WhitespaceTrimmed(t *testing.T) {
	t.Parallel()

	if IPReady("  ") {
		t.Error("expected IPReady=false for whitespace-only string")
	}
}

func TestIPReady_WhitespaceAroundValid(t *testing.T) {
	t.Parallel()

	if !IPReady("  10.0.1.5  ") {
		t.Error("expected IPReady=true for whitespace-padded valid IP")
	}
}

func TestIPReady_WhitespaceAroundZero(t *testing.T) {
	t.Parallel()

	if IPReady("  0.0.0.0  ") {
		t.Error("expected IPReady=false for whitespace-padded 0.0.0.0")
	}
}

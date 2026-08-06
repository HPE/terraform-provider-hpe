// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package containerip

import (
	"testing"
)

func TestUnitReady_EmptyString(t *testing.T) {
	t.Parallel()

	if Ready("") {
		t.Error("expected Ready=false for empty string")
	}
}

func TestUnitReady_ZeroAddress(t *testing.T) {
	t.Parallel()

	if Ready("0.0.0.0") {
		t.Error("expected Ready=false for 0.0.0.0")
	}
}

func TestUnitReady_ValidIPv4(t *testing.T) {
	t.Parallel()

	if !Ready("10.0.1.5") {
		t.Error("expected Ready=true for 10.0.1.5")
	}
}

func TestUnitReady_ValidIPv6(t *testing.T) {
	t.Parallel()

	if !Ready("fe80::1") {
		t.Error("expected Ready=true for fe80::1")
	}
}

func TestUnitReady_WhitespaceTrimmed(t *testing.T) {
	t.Parallel()

	if Ready("  ") {
		t.Error("expected Ready=false for whitespace-only string")
	}
}

func TestUnitReady_WhitespaceAroundValid(t *testing.T) {
	t.Parallel()

	if !Ready("  10.0.1.5  ") {
		t.Error("expected Ready=true for whitespace-padded valid IP")
	}
}

func TestUnitReady_WhitespaceAroundZero(t *testing.T) {
	t.Parallel()

	if Ready("  0.0.0.0  ") {
		t.Error("expected Ready=false for whitespace-padded 0.0.0.0")
	}
}

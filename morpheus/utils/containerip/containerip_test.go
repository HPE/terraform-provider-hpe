// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package containerip

import (
	"testing"
)

func TestReady_EmptyString(t *testing.T) {
	t.Parallel()

	if Ready("") {
		t.Error("expected Ready=false for empty string")
	}
}

func TestReady_ZeroAddress(t *testing.T) {
	t.Parallel()

	if Ready("0.0.0.0") {
		t.Error("expected Ready=false for 0.0.0.0")
	}
}

func TestReady_ValidIPv4(t *testing.T) {
	t.Parallel()

	if !Ready("10.0.1.5") {
		t.Error("expected Ready=true for 10.0.1.5")
	}
}

func TestReady_ValidIPv6(t *testing.T) {
	t.Parallel()

	if !Ready("fe80::1") {
		t.Error("expected Ready=true for fe80::1")
	}
}

func TestReady_WhitespaceTrimmed(t *testing.T) {
	t.Parallel()

	if Ready("  ") {
		t.Error("expected Ready=false for whitespace-only string")
	}
}

func TestReady_WhitespaceAroundValid(t *testing.T) {
	t.Parallel()

	if !Ready("  10.0.1.5  ") {
		t.Error("expected Ready=true for whitespace-padded valid IP")
	}
}

func TestReady_WhitespaceAroundZero(t *testing.T) {
	t.Parallel()

	if Ready("  0.0.0.0  ") {
		t.Error("expected Ready=false for whitespace-padded 0.0.0.0")
	}
}

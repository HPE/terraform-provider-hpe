// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cluster

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func strptr(s string) *string { return &s }

// TestDecideHvm2LayoutNoFieldsSet: when no 2.0 fields are set, the gate never
// refuses, regardless of the layout's clusterVersion (including nil).
func TestDecideHvm2LayoutNoFieldsSet(t *testing.T) {
	ctx := context.Background()

	for _, cv := range []*string{nil, strptr("1.0"), strptr("2.0"), strptr("garbage")} {
		refuse, _ := decideHvm2Layout(ctx, false, cv)
		if refuse {
			t.Errorf("decideHvm2Layout(anySet=false, %v) refused, want allow", cv)
		}
	}
}

// TestDecideHvm2LayoutAllows: fields set against a 2.0+ layout must be allowed.
func TestDecideHvm2LayoutAllows(t *testing.T) {
	ctx := context.Background()

	for _, cv := range []string{"2.0", "2.1", "2.0.1", "3.0", "10.0"} {
		t.Run(cv, func(t *testing.T) {
			refuse, _ := decideHvm2Layout(ctx, true, strptr(cv))
			if refuse {
				t.Errorf("decideHvm2Layout(anySet=true, %q) refused, want allow", cv)
			}
		})
	}
}

// TestDecideHvm2LayoutRefuses: fields set against a pre-2.0 layout must be
// refused, and the detail must name the actual clusterVersion.
func TestDecideHvm2LayoutRefuses(t *testing.T) {
	ctx := context.Background()

	for _, cv := range []string{"1.4", "1.3", "1.0", "0.9"} {
		t.Run(cv, func(t *testing.T) {
			refuse, detail := decideHvm2Layout(ctx, true, strptr(cv))
			if !refuse {
				t.Fatalf("decideHvm2Layout(anySet=true, %q) allowed, want refuse", cv)
			}

			if !strings.Contains(detail, cv) {
				t.Errorf("detail %q does not name the layout clusterVersion %q", detail, cv)
			}

			if !strings.Contains(detail, "2.0") {
				t.Errorf("detail %q does not name the required version 2.0", detail)
			}
		})
	}
}

// TestDecideHvm2LayoutFailsOpenNilVersion: fields set but the layout reports no
// clusterVersion must fail open (allow) rather than block.
func TestDecideHvm2LayoutFailsOpenNilVersion(t *testing.T) {
	refuse, _ := decideHvm2Layout(context.Background(), true, nil)
	if refuse {
		t.Errorf("decideHvm2Layout(anySet=true, nil) refused, want fail-open allow")
	}
}

// TestDecideHvm2LayoutFailsOpenUnparseable: an unparseable clusterVersion must
// fail open rather than block a possibly healthy operation.
func TestDecideHvm2LayoutFailsOpenUnparseable(t *testing.T) {
	for _, cv := range []string{"not-a-version", "", "abc.def"} {
		refuse, _ := decideHvm2Layout(context.Background(), true, strptr(cv))
		if refuse {
			t.Errorf("decideHvm2Layout(anySet=true, %q) refused, want fail-open allow", cv)
		}
	}
}

// TestHvm2FieldsSet verifies the shared trigger used by both the create and
// update gates: each of the nine HVM 2.0 network fields, set on its own, must
// trip the gate; an all-null config must not; and an unknown (not-yet-known)
// value must not count as set.
func TestHvm2FieldsSet(t *testing.T) {
	if hvm2FieldsSet(ConfigHvmValue{}) {
		t.Error("hvm2FieldsSet(all-null) = true, want false")
	}

	set := basetypes.NewStringValue("x")
	cases := map[string]ConfigHvmValue{
		"virt_switch_name": {VirtSwitchName: set},
		"dns1":             {Dns1: set},
		"dns2":             {Dns2: set},
		"dns3":             {Dns3: set},
		"search_domain":    {SearchDomain: set},
		"ntp1":             {Ntp1: set},
		"ntp2":             {Ntp2: set},
		"proxy":            {Proxy: set},
		"no_proxy":         {NoProxy: set},
	}

	for name, cfg := range cases {
		t.Run(name, func(t *testing.T) {
			if !hvm2FieldsSet(cfg) {
				t.Errorf("hvm2FieldsSet(%s set) = false, want true", name)
			}
		})
	}

	// An unknown value is not yet known and must not count as set.
	if hvm2FieldsSet(ConfigHvmValue{Dns1: basetypes.NewStringUnknown()}) {
		t.Error("hvm2FieldsSet(dns1 unknown) = true, want false")
	}
}

// TestClusterHvmVirtSwitchMinLayoutVersionIsDocumentedValue pins the constant
// quoted in the cluster resource docs ("clusterVersion 2.0 or greater").
func TestClusterHvmVirtSwitchMinLayoutVersionIsDocumentedValue(t *testing.T) {
	if ClusterHvmVirtSwitchMinLayoutVersion != ">= 2.0" {
		t.Errorf(
			"ClusterHvmVirtSwitchMinLayoutVersion = %q, want %q; the cluster "+
				"HVM 2.0 doc/example and gate detail state 2.0 and must be "+
				"updated together with this",
			ClusterHvmVirtSwitchMinLayoutVersion, ">= 2.0",
		)
	}
}

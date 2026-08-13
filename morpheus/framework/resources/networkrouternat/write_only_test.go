// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouternat

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// These tests drive the production builders directly. An earlier version
// re-implemented the branching inside the test body and asserted on its own
// copy, which passed even when the production default was deliberately broken.
// If a test here can pass while resource.go is wrong, it is not doing its job.

func TestUnitBuildCreateNatConfig_FirewallDefaultApplied(t *testing.T) {
	t.Parallel()

	got := buildCreateNatConfig(&NetworkRouterNatModel{
		Action:   types.StringValue("SNAT"),
		Firewall: types.StringNull(),
		Service:  types.StringNull(),
	})

	if got.Firewall == nil {
		t.Fatal("firewall should be defaulted on create, got nil")
	}

	if *got.Firewall != defaultFirewallMatch {
		t.Fatalf("firewall = %q, want %q", *got.Firewall, defaultFirewallMatch)
	}

	// The default exists because the Morpheus OptionType is required. Guard the
	// literal too: silently changing it would send an unaccepted value.
	if defaultFirewallMatch != "MATCH_INTERNAL_ADDRESS" {
		t.Fatalf("default firewall literal changed to %q", defaultFirewallMatch)
	}

	if got.Service != nil {
		t.Fatalf("service should be omitted when unset, got %q", *got.Service)
	}
}

func TestUnitBuildCreateNatConfig_ExplicitFirewallWins(t *testing.T) {
	t.Parallel()

	got := buildCreateNatConfig(&NetworkRouterNatModel{
		Action:   types.StringValue("DNAT"),
		Firewall: types.StringValue("BYPASS"),
		Service:  types.StringValue("HTTPS"),
	})

	if got.Firewall == nil || *got.Firewall != "BYPASS" {
		t.Fatalf("firewall = %v, want BYPASS", got.Firewall)
	}

	if got.Service == nil || *got.Service != "HTTPS" {
		t.Fatalf("service = %v, want HTTPS", got.Service)
	}

	if got.Action != "DNAT" {
		t.Fatalf("action = %q, want DNAT", got.Action)
	}
}

// The T2-2 fix. Update must omit an unset firewall rather than defaulting it:
// the controller merges the payload over the current config, so sending a
// default would overwrite the real NSX-T value with a guess.
func TestUnitBuildUpdateNatConfig_UnsetFirewallOmitted(t *testing.T) {
	t.Parallel()

	got := buildUpdateNatConfig(&NetworkRouterNatModel{
		Action:   types.StringValue("SNAT"),
		Firewall: types.StringNull(),
		Service:  types.StringNull(),
	})

	if got.Firewall != nil {
		t.Fatalf("firewall must be omitted on update, got %q — this would "+
			"clobber the server-side value", *got.Firewall)
	}

	if got.Service != nil {
		t.Fatalf("service must be omitted on update, got %q", *got.Service)
	}
}

func TestUnitBuildUpdateNatConfig_ExplicitValuesSent(t *testing.T) {
	t.Parallel()

	got := buildUpdateNatConfig(&NetworkRouterNatModel{
		Action:   types.StringValue("DNAT"),
		Firewall: types.StringValue("MATCH_EXTERNAL_ADDRESS"),
		Service:  types.StringValue("SSH"),
	})

	if got.Firewall == nil || *got.Firewall != "MATCH_EXTERNAL_ADDRESS" {
		t.Fatalf("firewall = %v, want MATCH_EXTERNAL_ADDRESS", got.Firewall)
	}

	if got.Service == nil || *got.Service != "SSH" {
		t.Fatalf("service = %v, want SSH", got.Service)
	}
}

// Guards the shape the builders depend on: both write-only fields are optional
// pointers, so "unset" is expressible. If the SDK ever makes them plain
// strings, omission becomes impossible and the update fix silently regresses.
func TestUnitNatConfigFieldsArePointers(t *testing.T) {
	t.Parallel()

	var upd sdk.UpdateNetworkRouterNatRequestNetworkRouterNATConfig
	if upd.Firewall != nil || upd.Service != nil {
		t.Fatal("expected nil-able firewall/service on the update config")
	}
}

// Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package networkrouter

import (
	"context"
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// setInterfacesState builds each InterfacesValue from an explicit attribute map,
// so it has to supply every attribute the generated schema declares. When the
// schema gained "config" during a regeneration the map was not updated, and
// NewInterfacesValue answered with diagnostics rather than a value — which this
// function turns into an error, failing the read for any router that has
// interfaces at all.
//
// Nothing exercised this path, so the break was invisible: the empty-interfaces
// case returns early and every existing test took it. This covers the populated
// case, and will fail again if the schema grows another attribute that the
// mapping does not supply.
func TestSetInterfacesStatePopulatesEveryDeclaredAttribute(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	id := int64(7)
	cidr := "10.0.0.0/24"
	ip := "10.0.0.1"

	router := &sdk.GetNetworkRouter200ResponseNetworkRouter{
		Interfaces: []sdk.GetNetworkRouter200ResponseNetworkRouterInterfacesInner{
			{
				Id:        &id,
				Cidr:      &cidr,
				IpAddress: &ip,
				// The API returns config as a free-form map. The generated
				// ConfigValue has no attributes to hold it, so the mapping can
				// only represent null — but it must still supply the key.
				Config: map[string]interface{}{"anything": "here"},
			},
		},
	}

	var state NetworkRouterModel

	if err := setInterfacesState(ctx, &state, router); err != nil {
		t.Fatalf("setInterfacesState returned an error for a router with interfaces: %v", err)
	}

	if state.Interfaces.IsNull() {
		t.Fatal("expected interfaces to be populated, got null")
	}

	if n := len(state.Interfaces.Elements()); n != 1 {
		t.Fatalf("expected 1 interface, got %d", n)
	}
}

// The early return is the case every previous test took. Kept so a change to
// the null handling is caught too.
func TestSetInterfacesStateNullWhenAbsent(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	var state NetworkRouterModel

	if err := setInterfacesState(ctx, &state, &sdk.GetNetworkRouter200ResponseNetworkRouter{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !state.Interfaces.IsNull() {
		t.Fatal("expected null interfaces when the router has none")
	}
}

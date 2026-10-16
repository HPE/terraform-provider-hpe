// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterbgpneighbor

import (
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// TestUnitPasswordHash checks the appliance-hash extractor: a set hash is
// returned present; a present-but-null field is present with an empty hash (the
// appliance holds no secret); an unset (API-omitted) field is not present (no
// information); a wrong type is not present.
func TestUnitPasswordHash(t *testing.T) {
	t.Parallel()

	hash := "deadbeef"

	set := &sdk.GetNetworkRouterBgpNeighbor200ResponseNetworkRouterBgpNeighbor{
		PasswordHash: *sdk.NewNullableString(&hash),
	}
	if h, ok := passwordHash(set); !ok || h != hash {
		t.Fatalf("set hash: got (%q,%v), want (%q,true)", h, ok, hash)
	}

	null := &sdk.GetNetworkRouterBgpNeighbor200ResponseNetworkRouterBgpNeighbor{
		PasswordHash: *sdk.NewNullableString(nil),
	}
	if h, ok := passwordHash(null); !ok || h != "" {
		t.Fatalf("null hash: got (%q,%v), want (\"\",true)", h, ok)
	}

	unset := &sdk.GetNetworkRouterBgpNeighbor200ResponseNetworkRouterBgpNeighbor{}
	if h, ok := passwordHash(unset); ok || h != "" {
		t.Fatalf("unset hash: got (%q,%v), want (\"\",false)", h, ok)
	}

	if h, ok := passwordHash("not the object"); ok || h != "" {
		t.Fatalf("wrong type: got (%q,%v), want (\"\",false)", h, ok)
	}

	if h, ok := passwordHash((*sdk.GetNetworkRouterBgpNeighbor200ResponseNetworkRouterBgpNeighbor)(nil)); ok || h != "" {
		t.Fatalf("nil object: got (%q,%v), want (\"\",false)", h, ok)
	}
}

// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkpoolserver

import (
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// TestUnitServicePasswordHash checks the appliance-hash extractor: a set hash
// is returned present; a nil pointer is present with an empty hash (Morpheus
// always renders the field, so nil means null -- no secret); a wrong type is
// not present.
func TestUnitServicePasswordHash(t *testing.T) {
	t.Parallel()

	hash := "deadbeef"

	set := &sdk.GetNetworkPoolServer200ResponseNetworkPoolServer{
		ServicePasswordHash: &hash,
	}
	if h, ok := servicePasswordHash(set); !ok || h != hash {
		t.Fatalf("set hash: got (%q,%v), want (%q,true)", h, ok, hash)
	}

	// A plain *string cannot tell omitted from null: nil is not-present.
	omitted := &sdk.GetNetworkPoolServer200ResponseNetworkPoolServer{}
	if h, ok := servicePasswordHash(omitted); ok || h != "" {
		t.Fatalf("nil hash: got (%q,%v), want (\"\",false)", h, ok)
	}

	if h, ok := servicePasswordHash("not the object"); ok || h != "" {
		t.Fatalf("wrong type: got (%q,%v), want (\"\",false)", h, ok)
	}

	if h, ok := servicePasswordHash((*sdk.GetNetworkPoolServer200ResponseNetworkPoolServer)(nil)); ok || h != "" {
		t.Fatalf("nil object: got (%q,%v), want (\"\",false)", h, ok)
	}
}

// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storageserver

import (
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// TestUnitServicePasswordHash checks the appliance-hash extractor: a set hash
// is returned present; a present-but-null field is present with an empty hash
// (the appliance holds no secret); an unset (API-omitted) field is not present
// (no information); a wrong type is not present.
func TestUnitServicePasswordHash(t *testing.T) {
	t.Parallel()

	hash := "deadbeef"

	set := &sdk.GetStorageServers200ResponseStorageServer{
		ServicePasswordHash: *sdk.NewNullableString(&hash),
	}
	if h, ok := servicePasswordHash(set); !ok || h != hash {
		t.Fatalf("set hash: got (%q,%v), want (%q,true)", h, ok, hash)
	}

	null := &sdk.GetStorageServers200ResponseStorageServer{
		ServicePasswordHash: *sdk.NewNullableString(nil),
	}
	if h, ok := servicePasswordHash(null); !ok || h != "" {
		t.Fatalf("null hash: got (%q,%v), want (\"\",true)", h, ok)
	}

	unset := &sdk.GetStorageServers200ResponseStorageServer{}
	if h, ok := servicePasswordHash(unset); ok || h != "" {
		t.Fatalf("unset hash: got (%q,%v), want (\"\",false)", h, ok)
	}

	if h, ok := servicePasswordHash("not the object"); ok || h != "" {
		t.Fatalf("wrong type: got (%q,%v), want (\"\",false)", h, ok)
	}

	if h, ok := servicePasswordHash((*sdk.GetStorageServers200ResponseStorageServer)(nil)); ok || h != "" {
		t.Fatalf("nil object: got (%q,%v), want (\"\",false)", h, ok)
	}
}

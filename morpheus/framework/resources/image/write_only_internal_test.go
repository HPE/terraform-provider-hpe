// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package image

import (
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// TestUnitSshPasswordHash checks the appliance-hash extractor: a set hash is
// returned present; a present-but-null field is present with an empty hash (the
// appliance holds no secret); an unset (API-omitted) field is not present (no
// information); a wrong type is not present.
func TestUnitSshPasswordHash(t *testing.T) {
	t.Parallel()

	hash := "deadbeef"

	set := &sdk.GetVirtualImage200ResponseVirtualImage{
		SshPasswordHash: *sdk.NewNullableString(&hash),
	}
	if h, ok := sshPasswordHash(set); !ok || h != hash {
		t.Fatalf("set hash: got (%q,%v), want (%q,true)", h, ok, hash)
	}

	null := &sdk.GetVirtualImage200ResponseVirtualImage{
		SshPasswordHash: *sdk.NewNullableString(nil),
	}
	if h, ok := sshPasswordHash(null); !ok || h != "" {
		t.Fatalf("null hash: got (%q,%v), want (\"\",true)", h, ok)
	}

	unset := &sdk.GetVirtualImage200ResponseVirtualImage{}
	if h, ok := sshPasswordHash(unset); ok || h != "" {
		t.Fatalf("unset hash: got (%q,%v), want (\"\",false)", h, ok)
	}

	if h, ok := sshPasswordHash("not the object"); ok || h != "" {
		t.Fatalf("wrong type: got (%q,%v), want (\"\",false)", h, ok)
	}

	if h, ok := sshPasswordHash((*sdk.GetVirtualImage200ResponseVirtualImage)(nil)); ok || h != "" {
		t.Fatalf("nil object: got (%q,%v), want (\"\",false)", h, ok)
	}
}

// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package hostprofile

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// ntpServersFromAPI decodes a JSON ntp block through the SDK's response model
// (the same mapstructure-backed decoder the live client uses) and returns the
// servers set that flattenSystem produces for it.
func ntpServersFromAPI(t *testing.T, ntpJSON string) types.Set {
	t.Helper()

	var system sdk.GetHostProfile200ResponseHostProfileAllOfConfigSystem
	if err := json.Unmarshal([]byte(`{"ntp":`+ntpJSON+`}`), &system); err != nil {
		t.Fatalf("decode: %v", err)
	}

	obj, diags := flattenSystem(context.Background(), &system)
	if diags.HasError() {
		t.Fatalf("flattenSystem: %v", diags.Errors())
	}

	ntp, ok := obj.Attributes()["ntp"].(types.Object)
	if !ok || ntp.IsNull() {
		t.Fatalf("ntp object missing or null: %v", obj.Attributes()["ntp"])
	}

	servers, ok := ntp.Attributes()["servers"].(types.Set)
	if !ok {
		t.Fatalf("servers is %T, want types.Set", ntp.Attributes()["servers"])
	}

	return servers
}

// TestUnitFlattenNtpServersPreservesEmptySet pins that an explicitly empty
// server list reads back as an empty set, not null. config.system.ntp.servers
// is plain Optional, so a configured `servers = []` is sent as "servers":[]
// and the read-back must return the same known empty set, or the apply fails
// with "provider produced inconsistent result". An absent servers key still
// reads back as null, which is what an omitted attribute plans as.
func TestUnitFlattenNtpServersPreservesEmptySet(t *testing.T) {
	t.Parallel()

	got := ntpServersFromAPI(t, `{"enabled":false,"servers":[]}`)
	if got.IsNull() || got.IsUnknown() {
		t.Fatalf(`"servers":[] read back as %v, want a known empty set`, got)
	}
	if n := len(got.Elements()); n != 0 {
		t.Fatalf(`"servers":[] read back with %d elements`, n)
	}

	want := types.SetValueMust(types.StringType, []attr.Value{})
	if !got.Equal(want) {
		t.Fatalf("empty set: got %v, want %v", got, want)
	}
}

func TestUnitFlattenNtpServersAbsentIsNull(t *testing.T) {
	t.Parallel()

	got := ntpServersFromAPI(t, `{"enabled":false}`)
	if !got.IsNull() {
		t.Fatalf("absent servers read back as %v, want null", got)
	}
}

func TestUnitFlattenNtpServersPopulated(t *testing.T) {
	t.Parallel()

	got := ntpServersFromAPI(t, `{"enabled":true,"servers":["a.ntp","b.ntp"]}`)
	want := types.SetValueMust(types.StringType, []attr.Value{
		types.StringValue("a.ntp"), types.StringValue("b.ntp"),
	})
	if !got.Equal(want) {
		t.Fatalf("populated set: got %v, want %v", got, want)
	}
}

// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package clusteraffinitygroup

import (
	"context"
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// TestUnitAffinityGroupAsStateResourcePermissionGroups guards the same class of
// bug as the plural data sources: resource_permissions declares its groups set
// with the custom GroupsType as the element type, so the elements have to be
// GroupsValue. Putting bare objects in the set makes types.SetValue reject
// them and hand back an unknown set, which then reaches state as an unknown
// value.
func TestUnitAffinityGroupAsStateResourcePermissionGroups(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ag := &sdk.GetClusterAffinityGroup200ResponseAffinityGroup{
		Id:   sdk.PtrInt64(1),
		Name: sdk.PtrString("example"),
		ResourcePermissions: &sdk.GetClusterAffinityGroup200ResponseAffinityGroupResourcePermissions{
			All: sdk.PtrBool(false),
			Sites: []sdk.GetClusterAffinityGroup200ResponseAffinityGroupResourcePermissionsSitesInner{
				{Id: sdk.PtrInt64(7), Default: sdk.PtrBool(true)},
			},
		},
	}

	state := affinityGroupAsState(ctx, ag, 1)

	groups := state.ResourcePermissions.Groups
	if groups.IsUnknown() {
		t.Fatal("resource_permissions.groups is unknown; the set rejected its elements")
	}

	if got := len(groups.Elements()); got != 1 {
		t.Fatalf("resource_permissions.groups has %d elements, want 1", got)
	}

	want := GroupsValue{}.Type(ctx)
	if got := groups.ElementType(ctx); !want.Equal(got) {
		t.Fatalf("groups element type = %s, want %s", got, want)
	}
}

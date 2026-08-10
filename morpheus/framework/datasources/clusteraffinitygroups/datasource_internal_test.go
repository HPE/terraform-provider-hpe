// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package clusteraffinitygroups

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// schemaElementType returns the element type the generated schema declares for
// the affinity_groups set. Reading it from the schema rather than restating it
// is the point: the Read implementation and the schema have to agree, and
// asserting against the schema is what catches them drifting apart.
func schemaElementType(t *testing.T, ctx context.Context) attr.Type {
	t.Helper()

	a, ok := ClusterAffinityGroupsDataSourceSchema(ctx).Attributes["affinity_groups"]
	if !ok {
		t.Fatal("schema has no affinity_groups attribute")
	}

	withElem, ok := a.GetType().(attr.TypeWithElementType)
	if !ok {
		t.Fatalf("affinity_groups is not a collection type: %T", a.GetType())
	}

	return withElem.ElementType()
}

// TestUnitAffinityGroupValueMatchesSchemaElementType asserts that the value the
// Read implementation puts in the affinity_groups set has exactly the element
// type the schema declares. Building a bare types.Object here instead produces
// "Invalid Set Element Type" at plan time against a real appliance, which is a
// failure mode that should never need an appliance to surface.
func TestUnitAffinityGroupValueMatchesSchemaElementType(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	elemType := schemaElementType(t, ctx)

	for name, ag := range affinityGroupFixtures() {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			v, diags := affinityGroupToValue(ctx, ag)
			if diags.HasError() {
				t.Fatalf("affinityGroupToValue: %v", diags.Errors())
			}

			if got := v.Type(ctx); !elemType.Equal(got) {
				t.Fatalf("element type mismatch:\n schema: %s\n  value: %s", elemType, got)
			}

			// The set itself is the thing that rejected the old bare objects,
			// so exercise it rather than trusting the type comparison alone.
			if _, d := types.SetValue(elemType, []attr.Value{v}); d.HasError() {
				t.Fatalf("types.SetValue rejected the element: %v", d.Errors())
			}
		})
	}
}

// affinityGroupFixtures covers both the minimal payload and one with every
// nested structure populated, because the nested groups set has its own custom
// element type and is only built when resource permissions come back.
func affinityGroupFixtures() map[string]*sdk.ListClusterAffinityGroups200ResponseAllOfAffinityGroupsInner {
	minimal := &sdk.ListClusterAffinityGroups200ResponseAllOfAffinityGroupsInner{
		Id:   sdk.PtrInt64(1),
		Name: sdk.PtrString("minimal"),
	}

	full := &sdk.ListClusterAffinityGroups200ResponseAllOfAffinityGroupsInner{
		Id:           sdk.PtrInt64(2),
		Name:         sdk.PtrString("full"),
		AffinityType: sdk.PtrString("KEEP_TOGETHER"),
		Source:       sdk.PtrString("morpheus"),
		RefType:      sdk.PtrString("ComputeServerGroup"),
		RefId:        sdk.PtrInt64(3),
		Active:       sdk.PtrBool(true),
		Visibility:   sdk.PtrString("private"),
		Pool: &sdk.ListClusterAffinityGroups200ResponseAllOfAffinityGroupsInnerPool{
			Id: sdk.PtrInt64(4),
		},
		Servers: []sdk.ListClusterAffinityGroups200ResponseAllOfAffinityGroupsInnerServersInner{
			{Id: sdk.PtrInt64(5)},
		},
		Tenants: []sdk.ListClusterAffinityGroups200ResponseAllOfAffinityGroupsInnerTenantsInner{
			{Id: sdk.PtrInt64(6)},
		},
		ResourcePermissions: &sdk.ListClusterAffinityGroups200ResponseAllOfAffinityGroupsInnerResourcePermissions{
			All: sdk.PtrBool(false),
			Sites: []sdk.ListClusterAffinityGroups200ResponseAllOfAffinityGroupsInnerResourcePermissionsSitesInner{
				{Id: sdk.PtrInt64(7), Default: sdk.PtrBool(true)},
			},
		},
	}

	return map[string]*sdk.ListClusterAffinityGroups200ResponseAllOfAffinityGroupsInner{
		"minimal": minimal,
		"full":    full,
	}
}

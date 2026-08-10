// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package computeservers

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// schemaElementType returns the element type the generated schema declares for
// the servers set. Reading it from the schema rather than restating it is the
// point: the Read implementation and the schema have to agree, and asserting
// against the schema is what catches them drifting apart.
func schemaElementType(t *testing.T, ctx context.Context) attr.Type {
	t.Helper()

	a, ok := ComputeServersDataSourceSchema(ctx).Attributes["servers"]
	if !ok {
		t.Fatal("schema has no servers attribute")
	}

	withElem, ok := a.GetType().(attr.TypeWithElementType)
	if !ok {
		t.Fatalf("servers is not a collection type: %T", a.GetType())
	}

	return withElem.ElementType()
}

// TestUnitServerValueMatchesSchemaElementType asserts that the value the Read
// implementation puts in the servers set has exactly the element type the
// schema declares. Building a bare types.Object here instead produces "Invalid
// Set Element Type" at plan time against a real appliance, which is a failure
// mode that should never need an appliance to surface.
func TestUnitServerValueMatchesSchemaElementType(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	elemType := schemaElementType(t, ctx)

	for name, srv := range serverFixtures() {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			v, diags := serverToValue(ctx, srv)
			if diags.HasError() {
				t.Fatalf("serverToValue: %v", diags.Errors())
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

// serverFixtures covers both the minimal payload and one with every nested
// structure populated, so the mapping is exercised on both branches of each
// nil check.
func serverFixtures() map[string]*sdk.ListHosts200ResponseAllOfServersInner {
	minimal := &sdk.ListHosts200ResponseAllOfServersInner{
		Id:   sdk.PtrInt64(1),
		Name: sdk.PtrString("minimal"),
	}

	full := &sdk.ListHosts200ResponseAllOfServersInner{
		Id:             sdk.PtrInt64(2),
		Name:           sdk.PtrString("full"),
		Hostname:       sdk.PtrString("full.example.com"),
		Uuid:           sdk.PtrString("6f0b1f8e-0000-0000-0000-000000000000"),
		Status:         sdk.PtrString("provisioned"),
		PowerState:     sdk.PtrString("on"),
		Visibility:     sdk.PtrString("private"),
		AgentInstalled: sdk.PtrBool(true),
		ZoneId:         sdk.PtrInt64(3),
		SiteId:         sdk.PtrInt64(4),
		MaxMemory:      sdk.PtrInt64(1024),
		MaxStorage:     sdk.PtrInt64(2048),
		Labels:         []string{"one", "two"},
		Description:    *sdk.NewNullableString(sdk.PtrString("a host")),
		ExternalId:     *sdk.NewNullableString(sdk.PtrString("ext-1")),
		InternalId:     *sdk.NewNullableString(sdk.PtrString("int-1")),
		ExternalIp:     *sdk.NewNullableString(sdk.PtrString("203.0.113.1")),
		InternalIp:     *sdk.NewNullableString(sdk.PtrString("10.0.0.1")),
		Platform:       *sdk.NewNullableString(sdk.PtrString("linux")),
		Zone: *sdk.NewNullableListHosts200ResponseAllOfServersInnerZone(
			&sdk.ListHosts200ResponseAllOfServersInnerZone{
				Id:   sdk.PtrInt64(3),
				Name: sdk.PtrString("a cloud"),
			},
		),
		ComputeServerType: &sdk.ListHosts200ResponseAllOfServersInnerComputeServerType{
			Id:      sdk.PtrInt64(5),
			Code:    sdk.PtrString("vmwareLinux"),
			Name:    sdk.PtrString("VMware Linux"),
			Managed: sdk.PtrBool(true),
		},
		Plan: &sdk.ListHosts200ResponseAllOfServersInnerPlan{
			Id:   *sdk.NewNullableInt64(sdk.PtrInt64(6)),
			Code: *sdk.NewNullableString(sdk.PtrString("plan-code")),
			Name: *sdk.NewNullableString(sdk.PtrString("plan name")),
		},
		Instance: &sdk.ListHosts200ResponseAllOfServersInnerInstance{
			Id: sdk.PtrInt64(7),
		},
	}

	return map[string]*sdk.ListHosts200ResponseAllOfServersInner{
		"minimal": minimal,
		"full":    full,
	}
}

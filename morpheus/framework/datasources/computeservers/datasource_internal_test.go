// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package computeservers

import (
	"context"
	"encoding/json"
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

// TestUnitListDecodesStatsTimestampWithNumericOffset pins the server list
// against an appliance that mixes two timestamp forms in one response. Morpheus
// writes a server's statistics sample time with a "Z" offset for most servers
// and with "+0000" for the node of an external Kubernetes cluster. The second
// form is not RFC 3339, and while the field was typed as a date-time one such
// server made the whole list undecodable, failing every lookup on the
// appliance.
func TestUnitListDecodesStatsTimestampWithNumericOffset(t *testing.T) {
	t.Parallel()

	const body = `{"servers":[
		{"id":1,"name":"hvm-host","stats":{"ts":"2026-11-02T15:39:25Z","usedMemory":1024}},
		{"id":2,"name":"k8s-node","stats":{"ts":"2026-11-02T15:39:25+0000","usedMemory":2048}},
		{"id":3,"name":"no-stats-yet","stats":{}}
	]}`

	var list sdk.ListHosts200Response
	if err := json.Unmarshal([]byte(body), &list); err != nil {
		t.Fatalf("server list did not decode: %v", err)
	}

	if got := len(list.Servers); got != 3 {
		t.Fatalf("decoded %d servers, want 3", got)
	}

	want := []string{"2026-11-02T15:39:25Z", "2026-11-02T15:39:25+0000", ""}
	for i, srv := range list.Servers {
		got := ""
		if srv.Stats != nil && srv.Stats.Ts != nil {
			got = *srv.Stats.Ts
		}

		if got != want[i] {
			t.Errorf("server %d: stats timestamp %q, want %q", i+1, got, want[i])
		}
	}

	// The lookup by id reads one server through a different response model,
	// with a statistics model of its own.
	const numericOffset = "2026-11-02T15:39:25+0000"

	var one sdk.GetHost200Response
	if err := json.Unmarshal([]byte(
		`{"server":{"id":2,"name":"k8s-node","stats":{"ts":"`+numericOffset+`"}}}`,
	), &one); err != nil {
		t.Fatalf("single server did not decode: %v", err)
	}

	if one.Server == nil || one.Server.Stats == nil || one.Server.Stats.Ts == nil {
		t.Fatalf("single server: stats timestamp missing after decode: %+v", one.Server)
	}

	if got := *one.Server.Stats.Ts; got != numericOffset {
		t.Errorf("single server: stats timestamp %q, want %q", got, numericOffset)
	}
}

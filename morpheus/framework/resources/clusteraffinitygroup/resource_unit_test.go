// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package clusteraffinitygroup

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// int64Set builds a known set of Int64 values. Passing no ids yields a known EMPTY
// set, which is deliberately distinct from a null or unknown one.
func int64Set(ids ...int64) types.Set {
	vals := make([]attr.Value, 0, len(ids))
	for _, id := range ids {
		vals = append(vals, types.Int64Value(id))
	}

	return types.SetValueMust(types.Int64Type, vals)
}

// TestResolveUpdateServers guards against a destructive data-loss bug.
//
// servers is Optional+Computed with no UseStateForUnknown, so Terraform marks it
// UNKNOWN in the plan whenever the practitioner has not configured it and some other
// attribute changes. The Morpheus API treats servers on update as a wholesale replace
// in which an omitted key and an empty array both REMOVE EVERY MEMBER, so resolving
// that unknown to an empty array silently destroyed the group's membership on
// something as innocuous as a rename.
//
// The unknown and null rows below are the regression guard: they must resolve to the
// membership held in STATE, never to an empty array.
// TestResolveUpdateServersSerialisation pins the on-the-wire consequence of the
// resolution, since that is what actually reaches the API.
//
// The SDK omits servers only when the slice is nil (ToMap guards with IsNil, which for
// a slice is reflect nil-ness, not emptiness). A known-empty set must therefore
// serialise as "servers": [] so an explicit `servers = []` really does clear the group,
// while an unknown plan must serialise the members held in state.
// group builds a known GroupsValue.
func group(id int64, dflt types.Bool) GroupsValue {
	return GroupsValue{
		Id:      types.Int64Value(id),
		Default: dflt,
		state:   attr.ValueStateKnown,
	}
}

// groupSet builds a known set of GroupsValue.
func groupSet(t *testing.T, groups ...GroupsValue) types.Set {
	t.Helper()

	set, diags := types.SetValueFrom(
		context.Background(), GroupsValue{}.Type(context.Background()), groups,
	)
	if diags.HasError() {
		t.Fatalf("building groups set: %v", diags.Errors())
	}

	return set
}

// rp builds a known ResourcePermissionsValue.
func rp(all types.Bool, groups types.Set) ResourcePermissionsValue {
	return ResourcePermissionsValue{
		All:    all,
		Groups: groups,
		state:  attr.ValueStateKnown,
	}
}

// assertNoUnknown fails if any part of the resource_permissions value is unknown.
// That is the whole point of the resolution: Terraform rejects a post-apply state
// containing an unknown at ANY depth, not just at the top level.
func assertNoUnknown(t *testing.T, got ResourcePermissionsValue) {
	t.Helper()

	if got.IsUnknown() {
		t.Fatalf("resource_permissions is unknown")
	}

	if got.IsNull() {
		return
	}

	if got.All.IsUnknown() {
		t.Fatalf("resource_permissions.all is unknown")
	}

	if got.Groups.IsUnknown() {
		t.Fatalf("resource_permissions.groups is unknown")
	}

	if got.Groups.IsNull() {
		return
	}

	for i, elem := range got.Groups.Elements() {
		g, ok := elem.(GroupsValue)
		if !ok {
			t.Fatalf("resource_permissions.groups[%d] is %T, want GroupsValue", i, elem)
		}
		if g.IsUnknown() {
			t.Fatalf("resource_permissions.groups[%d] is unknown", i)
		}
		if g.Id.IsUnknown() {
			t.Fatalf("resource_permissions.groups[%d].id is unknown", i)
		}
		if g.Default.IsUnknown() {
			t.Fatalf("resource_permissions.groups[%d].default is unknown", i)
		}
	}
}

// TestResolveTenantIds guards the "Provider returned invalid result object after
// apply" failure reported against a live Morpheus 9.0.2 appliance.
//
// tenant_ids is Optional+Computed, so leaving it out of the configuration makes
// Terraform mark it UNKNOWN in the plan. The resource preserves the planned value
// rather than the API echo (to avoid a perpetual diff when the API normalises the
// list), and preserving that unknown put an unknown into post-apply state. Every row
// below must yield a KNOWN value; the unknown-plan rows are the regression guard.
func TestResolveTenantIds(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		planned types.Set
		fromAPI types.Set
		want    types.Set
	}{
		// UNKNOWN plan: the practitioner configured nothing, so the API is the
		// only source of truth and a typed null is the fallback.
		"unknown plan takes the api tenants": {
			planned: types.SetUnknown(types.Int64Type),
			fromAPI: int64Set(4, 5),
			want:    int64Set(4, 5),
		},
		"unknown plan with no api tenants becomes a typed null": {
			// What the appliance actually returns for a fresh group: "tenants": [].
			planned: types.SetUnknown(types.Int64Type),
			fromAPI: int64Set(),
			want:    types.SetNull(types.Int64Type),
		},
		"unknown plan with null api becomes a typed null": {
			planned: types.SetUnknown(types.Int64Type),
			fromAPI: types.SetNull(types.Int64Type),
			want:    types.SetNull(types.Int64Type),
		},
		"unknown plan with unknown api becomes a typed null": {
			planned: types.SetUnknown(types.Int64Type),
			fromAPI: types.SetUnknown(types.Int64Type),
			want:    types.SetNull(types.Int64Type),
		},

		// KNOWN plan: preserved verbatim, so a normalised or reordered API echo
		// cannot produce a perpetual diff.
		"known plan is preserved over a reordered api echo": {
			planned: int64Set(7, 8),
			fromAPI: int64Set(8, 7, 9),
			want:    int64Set(7, 8),
		},
		"known empty plan is preserved": {
			planned: int64Set(),
			fromAPI: int64Set(9),
			want:    int64Set(),
		},
		"null plan is preserved": {
			planned: types.SetNull(types.Int64Type),
			fromAPI: int64Set(9),
			want:    types.SetNull(types.Int64Type),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := resolveTenantIds(tc.planned, tc.fromAPI)

			// The regression guard. Terraform rejects an unknown here outright.
			if got.IsUnknown() {
				t.Fatalf("resolveTenantIds returned an UNKNOWN value")
			}

			// A null must be a correctly typed null, not a zero value: the
			// element type has to match the schema or the state write fails.
			wantType := types.SetType{ElemType: types.Int64Type}
			if gotType := got.Type(context.Background()); !gotType.Equal(wantType) {
				t.Fatalf("resolveTenantIds type = %s, want %s", gotType, wantType)
			}

			if !got.Equal(tc.want) {
				t.Fatalf("resolveTenantIds = %s, want %s", got, tc.want)
			}
		})
	}
}

// TestResolveResourcePermissions guards the second half of the same failure.
//
// resource_permissions is Optional+Computed, and so are its `all` and `groups`
// members and the `id`/`default` inside each group, so a partial configuration leaves
// unknowns nested inside an otherwise known object. Every row must come back with no
// unknown at any depth.
func TestResolveResourcePermissions(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	nullGroups := types.SetNull(GroupsValue{}.Type(ctx))

	tests := map[string]struct {
		planned ResourcePermissionsValue
		fromAPI ResourcePermissionsValue
		want    ResourcePermissionsValue
	}{
		// UNKNOWN attribute: not configured, so take the API value.
		"unknown plan with null api becomes null": {
			// What the appliance actually returns: "resourcePermissions": null.
			planned: NewResourcePermissionsValueUnknown(),
			fromAPI: NewResourcePermissionsValueNull(),
			want:    NewResourcePermissionsValueNull(),
		},
		"unknown plan with unknown api becomes null": {
			planned: NewResourcePermissionsValueUnknown(),
			fromAPI: NewResourcePermissionsValueUnknown(),
			want:    NewResourcePermissionsValueNull(),
		},
		"unknown plan takes the api value": {
			planned: NewResourcePermissionsValueUnknown(),
			fromAPI: rp(types.BoolValue(true), groupSet(t, group(1, types.BoolValue(true)))),
			want:    rp(types.BoolValue(true), groupSet(t, group(1, types.BoolValue(true)))),
		},

		// KNOWN attribute: preserved, so the API echo cannot churn the diff.
		"known plan is preserved": {
			planned: rp(types.BoolValue(false), groupSet(t, group(1, types.BoolValue(false)))),
			fromAPI: rp(types.BoolValue(true), groupSet(t, group(2, types.BoolValue(true)))),
			want:    rp(types.BoolValue(false), groupSet(t, group(1, types.BoolValue(false)))),
		},
		"null plan is preserved": {
			planned: NewResourcePermissionsValueNull(),
			fromAPI: rp(types.BoolValue(true), groupSet(t, group(1, types.BoolValue(true)))),
			want:    NewResourcePermissionsValueNull(),
		},

		// Partial configuration: the object is known but its Optional+Computed
		// members are not. Each one is resolved individually.
		"unknown all is filled from the api": {
			planned: rp(types.BoolUnknown(), groupSet(t, group(1, types.BoolValue(false)))),
			fromAPI: rp(types.BoolValue(true), groupSet(t, group(1, types.BoolValue(false)))),
			want:    rp(types.BoolValue(true), groupSet(t, group(1, types.BoolValue(false)))),
		},
		"unknown all with null api falls back to null": {
			planned: rp(types.BoolUnknown(), groupSet(t, group(1, types.BoolValue(false)))),
			fromAPI: NewResourcePermissionsValueNull(),
			want:    rp(types.BoolNull(), groupSet(t, group(1, types.BoolValue(false)))),
		},
		"unknown groups is filled from the api": {
			planned: rp(types.BoolValue(true), types.SetUnknown(GroupsValue{}.Type(ctx))),
			fromAPI: rp(types.BoolValue(true), groupSet(t, group(3, types.BoolValue(true)))),
			want:    rp(types.BoolValue(true), groupSet(t, group(3, types.BoolValue(true)))),
		},
		"unknown groups with null api falls back to a typed null set": {
			planned: rp(types.BoolValue(true), types.SetUnknown(GroupsValue{}.Type(ctx))),
			fromAPI: NewResourcePermissionsValueNull(),
			want:    rp(types.BoolValue(true), nullGroups),
		},
		"unknown nested default is filled from the matching api group": {
			planned: rp(types.BoolValue(false), groupSet(t, group(1, types.BoolUnknown()))),
			fromAPI: rp(
				types.BoolValue(false),
				groupSet(t, group(2, types.BoolValue(false)), group(1, types.BoolValue(true))),
			),
			want: rp(types.BoolValue(false), groupSet(t, group(1, types.BoolValue(true)))),
		},
		"unknown nested default with no api match falls back to null": {
			planned: rp(types.BoolValue(false), groupSet(t, group(1, types.BoolUnknown()))),
			fromAPI: rp(types.BoolValue(false), groupSet(t, group(2, types.BoolValue(true)))),
			want:    rp(types.BoolValue(false), groupSet(t, group(1, types.BoolNull()))),
		},
		"unknown nested default with null api falls back to null": {
			planned: rp(types.BoolValue(false), groupSet(t, group(1, types.BoolUnknown()))),
			fromAPI: NewResourcePermissionsValueNull(),
			want:    rp(types.BoolValue(false), groupSet(t, group(1, types.BoolNull()))),
		},
		// id is Optional+Computed on the CLUSTER resource, so it too can arrive
		// unknown. It cannot be correlated, so it resolves to null rather than
		// staying unknown.
		"unknown nested id falls back to null": {
			planned: rp(types.BoolValue(false), groupSet(t, GroupsValue{
				Id:      types.Int64Unknown(),
				Default: types.BoolUnknown(),
				state:   attr.ValueStateKnown,
			})),
			fromAPI: rp(types.BoolValue(false), groupSet(t, group(2, types.BoolValue(true)))),
			want: rp(types.BoolValue(false), groupSet(t, GroupsValue{
				Id:      types.Int64Null(),
				Default: types.BoolNull(),
				state:   attr.ValueStateKnown,
			})),
		},
		"known empty groups is preserved": {
			planned: rp(types.BoolValue(true), groupSet(t)),
			fromAPI: rp(types.BoolValue(true), groupSet(t, group(1, types.BoolValue(true)))),
			want:    rp(types.BoolValue(true), groupSet(t)),
		},
		"null groups is preserved": {
			planned: rp(types.BoolValue(true), nullGroups),
			fromAPI: rp(types.BoolValue(true), groupSet(t, group(1, types.BoolValue(true)))),
			want:    rp(types.BoolValue(true), nullGroups),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, diags := resolveResourcePermissions(context.Background(), tc.planned, tc.fromAPI)
			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags.Errors())
			}

			// The regression guard, applied at every depth.
			assertNoUnknown(t, got)

			if !got.Equal(tc.want) {
				t.Fatalf("resolveResourcePermissions = %s, want %s", got, tc.want)
			}
		})
	}
}

// createModelState reproduces exactly what Create writes to state: map the read-back,
// resolve the two preserved attributes, then backfill the ID from the create response.
func createModelState(
	t *testing.T,
	plan ClusterAffinityGroupModel,
	ag *sdk.GetClusterAffinityGroup200ResponseAffinityGroup,
	clusterID int64,
	id int64,
) tfsdk.State {
	t.Helper()

	ctx := context.Background()

	if diags := mapAndResolveResponse(ctx, &plan, ag, clusterID); diags.HasError() {
		t.Fatalf("mapAndResolveResponse: %v", diags.Errors())
	}

	if plan.Id.IsUnknown() {
		plan.Id = types.Int64Value(id)
	}

	state := tfsdk.State{Schema: ClusterAffinityGroupResourceSchema(ctx)}
	if diags := state.Set(ctx, &plan); diags.HasError() {
		t.Fatalf("state.Set: %v", diags.Errors())
	}

	return state
}

// unconfiguredPlan is the plan Terraform hands Create when the practitioner sets only
// the required attributes: every Optional+Computed attribute is UNKNOWN.
func unconfiguredPlan(clusterID int64, name string) ClusterAffinityGroupModel {
	return ClusterAffinityGroupModel{
		Active:              types.BoolUnknown(),
		AffinityType:        types.StringUnknown(),
		ClusterId:           types.Int64Value(clusterID),
		Id:                  types.Int64Unknown(),
		Name:                types.StringValue(name),
		PoolId:              types.Int64Unknown(),
		ResourcePermissions: NewResourcePermissionsValueUnknown(),
		Servers:             types.SetUnknown(types.Int64Type),
		Source:              types.StringUnknown(),
		TenantIds:           types.SetUnknown(types.Int64Type),
		Visibility:          types.StringUnknown(),
	}
}

// TestCreateStateIsFullyKnown is the end-to-end regression test for the acceptance
// failure: "After the apply operation, the provider still indicated an unknown value
// for hpe_morpheus_cluster_affinity_group.example.<attr>. All values must be known
// after apply, so this is always a bug in the provider."
//
// It runs the plan the framework produces for an unconfigured resource through the
// exact sequence Create uses, writes the result to a real tfsdk.State backed by the
// generated schema, and asserts that Terraform's own IsFullyKnown check passes. That
// covers every attribute at once rather than the two the acceptance run happened to
// name, so a regression in any of them fails here.
func TestCreateStateIsFullyKnown(t *testing.T) {
	t.Parallel()

	tests := map[string]*sdk.GetClusterAffinityGroup200ResponseAffinityGroup{
		// Mirrors what a Morpheus 9.0.2 appliance returns for a freshly created
		// group: resourcePermissions is null and tenants is empty, which is
		// precisely what broke the preserve-the-plan path.
		"appliance response for a freshly created group": {
			Id:           ptr(int64(6126)),
			Name:         ptr("example"),
			AffinityType: ptr("KEEP_TOGETHER"),
			Source:       ptr("user"),
			Active:       ptr(true),
			Visibility:   ptr("private"),
			Pool: &sdk.GetClusterAffinityGroup200ResponseAffinityGroupPool{
				Id: ptr(int64(1)),
			},
			Servers:             []sdk.GetClusterAffinityGroup200ResponseAffinityGroupServersInner{},
			Tenants:             []sdk.GetClusterAffinityGroup200ResponseAffinityGroupTenantsInner{},
			ResourcePermissions: nil,
		},
		// Nothing populated at all. No attribute may fall through its nil guard and
		// leave the planned unknown in place — the defect that hid in `active`.
		"empty response": {},
		// Fully populated, including permissions, to exercise the mapped branches.
		"fully populated response": {
			Id:           ptr(int64(7)),
			Name:         ptr("example"),
			AffinityType: ptr("KEEP_SEPARATE"),
			Source:       ptr("user"),
			Active:       ptr(false),
			Visibility:   ptr("public"),
			Pool: &sdk.GetClusterAffinityGroup200ResponseAffinityGroupPool{
				Id: ptr(int64(2)),
			},
			Servers: []sdk.GetClusterAffinityGroup200ResponseAffinityGroupServersInner{
				{Id: ptr(int64(11))},
			},
			Tenants: []sdk.GetClusterAffinityGroup200ResponseAffinityGroupTenantsInner{
				{Id: ptr(int64(3))},
			},
			ResourcePermissions: &sdk.GetClusterAffinityGroup200ResponseAffinityGroupResourcePermissions{
				All: ptr(true),
				Sites: []sdk.GetClusterAffinityGroup200ResponseAffinityGroupResourcePermissionsSitesInner{
					{Id: ptr(int64(4)), Default: ptr(true)},
				},
			},
		},
	}

	for name, ag := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			state := createModelState(t, unconfiguredPlan(3, "example"), ag, 3, 6126)

			if !state.Raw.IsFullyKnown() {
				t.Fatalf("post-apply state still contains unknown values: %s", state.Raw)
			}
		})
	}
}

// TestCreateStateIsFullyKnownWithPartialResourcePermissions covers the same defect one
// level down: a configured resource_permissions whose Optional+Computed members were
// left out of the configuration, so Terraform marks them unknown inside an otherwise
// known object.
func TestCreateStateIsFullyKnownWithPartialResourcePermissions(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	plans := map[string]ResourcePermissionsValue{
		// resource_permissions = { all = true }
		"groups omitted": rp(types.BoolValue(true), types.SetUnknown(GroupsValue{}.Type(ctx))),
		// resource_permissions = { groups = [{ id = 4 }] }
		"all and nested default omitted": rp(
			types.BoolUnknown(), groupSet(t, group(4, types.BoolUnknown())),
		),
	}

	for name, plannedRP := range plans {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			plan := unconfiguredPlan(3, "example")
			plan.ResourcePermissions = plannedRP

			state := createModelState(t, plan, &sdk.GetClusterAffinityGroup200ResponseAffinityGroup{
				Id:                  ptr(int64(6126)),
				Name:                ptr("example"),
				Active:              ptr(true),
				ResourcePermissions: nil,
			}, 3, 6126)

			if !state.Raw.IsFullyKnown() {
				t.Fatalf("post-apply state still contains unknown values: %s", state.Raw)
			}
		})
	}
}

// ptr returns a pointer to v, for building SDK response structs inline.
func ptr[T any](v T) *T {
	return &v
}

// importedModelState reproduces what Read writes to state on import: the prior state
// holds only the two IDs that ImportState set, so the response is mapped straight in
// with no preservation.
func importedModelState(
	t *testing.T,
	ag *sdk.GetClusterAffinityGroup200ResponseAffinityGroup,
	clusterID int64,
	id int64,
) tfsdk.State {
	t.Helper()

	ctx := context.Background()

	imported := ClusterAffinityGroupModel{
		Active:              types.BoolNull(),
		AffinityType:        types.StringNull(),
		ClusterId:           types.Int64Value(clusterID),
		Id:                  types.Int64Value(id),
		Name:                types.StringNull(),
		PoolId:              types.Int64Null(),
		ResourcePermissions: NewResourcePermissionsValueNull(),
		Servers:             types.SetNull(types.Int64Type),
		Source:              types.StringNull(),
		TenantIds:           types.SetNull(types.Int64Type),
		Visibility:          types.StringNull(),
	}

	if diags := mapGetResponseToModel(ctx, &imported, ag, clusterID); diags.HasError() {
		t.Fatalf("mapGetResponseToModel: %v", diags.Errors())
	}

	state := tfsdk.State{Schema: ClusterAffinityGroupResourceSchema(ctx)}
	if diags := state.Set(ctx, &imported); diags.HasError() {
		t.Fatalf("state.Set: %v", diags.Errors())
	}

	return state
}

// TestImportStateMatchesCreateState guards ImportStateVerify.
//
// Create resolves an unconfigured (UNKNOWN) tenant_ids / resource_permissions to a
// typed null, while import maps the same response straight in. If the mapping
// represented "no tenants" as an empty set the two states would disagree — the created
// resource would carry a null and the imported one a `tenant_ids = []` — and the
// import step of the acceptance test would fail on a difference that is pure
// representation. Both paths must land on the same value.
func TestImportStateMatchesCreateState(t *testing.T) {
	t.Parallel()

	tests := map[string]*sdk.GetClusterAffinityGroup200ResponseAffinityGroup{
		"appliance response with no tenants and no permissions": {
			Id:           ptr(int64(6126)),
			Name:         ptr("example"),
			AffinityType: ptr("KEEP_TOGETHER"),
			Source:       ptr("user"),
			Active:       ptr(true),
			Visibility:   ptr("private"),
			Pool: &sdk.GetClusterAffinityGroup200ResponseAffinityGroupPool{
				Id: ptr(int64(1)),
			},
			Servers:             []sdk.GetClusterAffinityGroup200ResponseAffinityGroupServersInner{},
			Tenants:             []sdk.GetClusterAffinityGroup200ResponseAffinityGroupTenantsInner{},
			ResourcePermissions: nil,
		},
		"response with tenants and permissions": {
			Id:         ptr(int64(6126)),
			Name:       ptr("example"),
			Active:     ptr(true),
			Visibility: ptr("private"),
			Tenants: []sdk.GetClusterAffinityGroup200ResponseAffinityGroupTenantsInner{
				{Id: ptr(int64(3))},
			},
			ResourcePermissions: &sdk.GetClusterAffinityGroup200ResponseAffinityGroupResourcePermissions{
				All: ptr(true),
				Sites: []sdk.GetClusterAffinityGroup200ResponseAffinityGroupResourcePermissionsSitesInner{
					{Id: ptr(int64(4)), Default: ptr(true)},
				},
			},
		},
	}

	for name, ag := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			created := createModelState(t, unconfiguredPlan(3, "example"), ag, 3, 6126)
			imported := importedModelState(t, ag, 3, 6126)

			if !created.Raw.Equal(imported.Raw) {
				t.Fatalf("import state differs from create state:\ncreate: %s\nimport: %s",
					created.Raw, imported.Raw)
			}
		})
	}
}

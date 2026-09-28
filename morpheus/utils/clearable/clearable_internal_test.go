// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package clearable

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// testSchema has one attribute of each kind the classification has to deal
// with: Optional+Computed strings, an integer and sets, an Optional-only
// string, and the kinds that are never clearable — required, computed-only,
// write-only and boolean.
func testSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":          schema.StringAttribute{Computed: true},
			"name":        schema.StringAttribute{Required: true},
			"description": schema.StringAttribute{Optional: true, Computed: true},
			"gateway":     schema.StringAttribute{Optional: true, Computed: true},
			"cidr":        schema.StringAttribute{Optional: true},
			"vlan_id":     schema.Int64Attribute{Optional: true, Computed: true},
			"labels":      schema.SetAttribute{ElementType: types.StringType, Optional: true, Computed: true},
			"tenant_ids":  schema.SetAttribute{ElementType: types.Int64Type, Optional: true, Computed: true},
			"enabled":     schema.BoolAttribute{Optional: true, Computed: true},
			"secret":      schema.StringAttribute{Optional: true, WriteOnly: true},
		},
	}
}

// testAttrs is the classification of testSchema's clearable attributes.
var testAttrs = Attributes{
	{Name: "description", Mode: Clear},
	{Name: "gateway", Mode: Release},
	{Name: "cidr", Mode: Clear},
	{Name: "vlan_id", Mode: Clear},
	{Name: "labels", Mode: Clear},
	{Name: "tenant_ids", Mode: Clear},
}

// testNot is the rest of testSchema's optional attributes, with reasons.
var testNot = NotClearable{
	"enabled": "a boolean has no cleared state",
}

// object builds a fully populated object value for testSchema: every
// attribute the caller does not name is null.
func object(ctx context.Context, t *testing.T, s schema.Schema, set map[string]tftypes.Value) tftypes.Value {
	t.Helper()

	objType, ok := s.Type().TerraformType(ctx).(tftypes.Object)
	if !ok {
		t.Fatalf("schema type is %T, want tftypes.Object", s.Type().TerraformType(ctx))
	}

	vals := make(map[string]tftypes.Value, len(objType.AttributeTypes))

	for name, typ := range objType.AttributeTypes {
		if v, ok := set[name]; ok {
			vals[name] = v

			continue
		}

		vals[name] = tftypes.NewValue(typ, nil)
	}

	return tftypes.NewValue(objType, vals)
}

func str(v string) tftypes.Value {
	return tftypes.NewValue(tftypes.String, v)
}

func num(v int64) tftypes.Value {
	return tftypes.NewValue(tftypes.Number, v)
}

func strSet(vs ...string) tftypes.Value {
	elems := make([]tftypes.Value, 0, len(vs))
	for _, v := range vs {
		elems = append(elems, str(v))
	}

	return tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, elems)
}

func numSet(vs ...int64) tftypes.Value {
	elems := make([]tftypes.Value, 0, len(vs))
	for _, v := range vs {
		elems = append(elems, num(v))
	}

	return tftypes.NewValue(tftypes.Set{ElementType: tftypes.Number}, elems)
}

func config(ctx context.Context, t *testing.T, set map[string]tftypes.Value) tfsdk.Config {
	t.Helper()

	s := testSchema()

	return tfsdk.Config{Schema: s, Raw: object(ctx, t, s, set)}
}

func state(ctx context.Context, t *testing.T, set map[string]tftypes.Value) tfsdk.State {
	t.Helper()

	s := testSchema()

	return tfsdk.State{Schema: s, Raw: object(ctx, t, s, set)}
}

func plan(ctx context.Context, t *testing.T, set map[string]tftypes.Value) tfsdk.Plan {
	t.Helper()

	s := testSchema()

	return tfsdk.Plan{Schema: s, Raw: object(ctx, t, s, set)}
}

// fakePrivate is the framework's private state, as far as the package uses
// it.
type fakePrivate struct {
	data map[string][]byte
}

func (f *fakePrivate) GetKey(_ context.Context, key string) ([]byte, diag.Diagnostics) {
	return f.data[key], nil
}

func (f *fakePrivate) SetKey(_ context.Context, key string, value []byte) diag.Diagnostics {
	if f.data == nil {
		f.data = map[string][]byte{}
	}

	f.data[key] = value

	return nil
}

func mustNoDiags(t *testing.T, diags diag.Diagnostics) {
	t.Helper()

	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
}

func TestRecordWritesTheConfiguredClearableNamesSorted(t *testing.T) {
	ctx := context.Background()

	// name is set but is not clearable; gateway and vlan_id are clearable but
	// not set; the rest are set and must be recorded, in name order.
	cfg := config(ctx, t, map[string]tftypes.Value{
		"name":        str("n"),
		"labels":      strSet("a"),
		"description": str("d"),
		"cidr":        str("10.0.0.0/24"),
		"tenant_ids":  numSet(1, 2),
	})

	private := &fakePrivate{}

	mustNoDiags(t, Record(ctx, cfg, private, testAttrs))

	got := string(private.data[PrivateKey])
	want := `{"v":1,"attrs":["cidr","description","labels","tenant_ids"]}`

	if got != want {
		t.Fatalf("record = %s, want %s", got, want)
	}

	if !json.Valid(private.data[PrivateKey]) {
		t.Fatalf("record is not valid JSON: %s", got)
	}
}

func TestRecordWithNothingConfiguredWritesAnEmptyRecord(t *testing.T) {
	ctx := context.Background()

	private := &fakePrivate{}

	mustNoDiags(t, Record(ctx, config(ctx, t, map[string]tftypes.Value{"name": str("n")}), private, testAttrs))

	if got := string(private.data[PrivateKey]); got != `{"v":1,"attrs":[]}` {
		t.Fatalf("record = %s, want an empty attrs list", got)
	}
}

func TestRecordWithoutPrivateStateIsAnError(t *testing.T) {
	ctx := context.Background()

	diags := Record(ctx, config(ctx, t, nil), nil, testAttrs)
	if !diags.HasError() {
		t.Fatal("expected an error when private state is unavailable")
	}
}

func TestComputeWithoutARecordPlansNothing(t *testing.T) {
	ctx := context.Background()

	// An imported or upgraded resource: description is in state and not in
	// config, but there is no record saying the user ever set it.
	pending, diags := Compute(ctx,
		config(ctx, t, map[string]tftypes.Value{"name": str("n")}),
		state(ctx, t, map[string]tftypes.Value{"name": str("n"), "description": str("server-set")}),
		&fakePrivate{},
		testAttrs,
	)

	mustNoDiags(t, diags)

	if len(pending) != 0 {
		t.Fatalf("pending = %v, want none", pending.Names())
	}
}

func TestComputeOptionalOnlyAttributeNeedsNoRecord(t *testing.T) {
	ctx := context.Background()

	// cidr is Optional but not Computed: its plan is its configuration, so a
	// null with a value in state is a removal even with no record (an
	// upgraded or imported resource) — Terraform will hold the apply to it.
	// description is Optional+Computed and, unrecorded, is left alone.
	pending, diags := Compute(ctx,
		config(ctx, t, map[string]tftypes.Value{"name": str("n")}),
		state(ctx, t, map[string]tftypes.Value{"name": str("n"), "cidr": str("10.0.0.0/24"), "description": str("d")}),
		&fakePrivate{},
		testAttrs,
	)

	mustNoDiags(t, diags)

	if got := strings.Join(pending.Names(), ","); got != "cidr" {
		t.Fatalf("pending = %q, want %q", got, "cidr")
	}
}

func TestComputeRejectsAnAttributeMissingFromTheSchema(t *testing.T) {
	ctx := context.Background()

	_, diags := Compute(ctx,
		config(ctx, t, map[string]tftypes.Value{"name": str("n")}),
		state(ctx, t, map[string]tftypes.Value{"name": str("n")}),
		&fakePrivate{},
		Attributes{{Name: "ghost"}},
	)

	if !diags.HasError() {
		t.Fatal("expected an error for an attribute that is not in the schema")
	}
}

func TestComputeWithoutStatePlansNothing(t *testing.T) {
	ctx := context.Background()

	private := &fakePrivate{data: map[string][]byte{PrivateKey: []byte(`{"v":1,"attrs":["description"]}`)}}

	// Creating: the state is null as a whole.
	pending, diags := Compute(ctx,
		config(ctx, t, map[string]tftypes.Value{"name": str("n")}),
		tfsdk.State{Schema: testSchema(), Raw: tftypes.NewValue(testSchema().Type().TerraformType(ctx), nil)},
		private,
		testAttrs,
	)

	mustNoDiags(t, diags)

	if len(pending) != 0 {
		t.Fatalf("pending = %v, want none", pending.Names())
	}
}

func TestComputePlansRecordedAttributesThatAreNullNowAndSetInState(t *testing.T) {
	ctx := context.Background()

	// Last apply set description, cidr, gateway and vlan_id.
	private := &fakePrivate{data: map[string][]byte{
		PrivateKey: []byte(`{"v":1,"attrs":["cidr","description","gateway","vlan_id"]}`),
	}}

	// Now: description and cidr removed (state still has them) -> pending;
	// gateway removed but state has nothing to clear -> not pending;
	// vlan_id still configured -> not pending;
	// labels null with a server value in state but never recorded -> left
	// to the server.
	pending, diags := Compute(ctx,
		config(ctx, t, map[string]tftypes.Value{"name": str("n"), "vlan_id": num(7)}),
		state(ctx, t, map[string]tftypes.Value{
			"name":        str("n"),
			"description": str("d"),
			"cidr":        str("10.0.0.0/24"),
			"vlan_id":     num(7),
			"labels":      strSet("from-server"),
		}),
		private,
		testAttrs,
	)

	mustNoDiags(t, diags)

	if got := strings.Join(pending.Names(), ","); got != "cidr,description" {
		t.Fatalf("pending = %q, want %q", got, "cidr,description")
	}

	if !pending.Clears("cidr") || pending.Clears("labels") || pending.Clears("gateway") {
		t.Fatalf("Clears() disagrees with Names(): %v", pending)
	}

	if pending["description"].Mode != Clear {
		t.Fatalf("description mode = %v, want Clear", pending["description"].Mode)
	}
}

func TestComputeIgnoresRecordsItCannotRead(t *testing.T) {
	ctx := context.Background()

	for name, raw := range map[string]string{
		"future version": `{"v":2,"attrs":["description"]}`,
		"garbage":        `not json`,
		"empty":          ``,
	} {
		t.Run(name, func(t *testing.T) {
			private := &fakePrivate{data: map[string][]byte{PrivateKey: []byte(raw)}}

			pending, diags := Compute(ctx,
				config(ctx, t, map[string]tftypes.Value{"name": str("n")}),
				state(ctx, t, map[string]tftypes.Value{"name": str("n"), "description": str("d")}),
				private,
				testAttrs,
			)

			mustNoDiags(t, diags)

			if len(pending) != 0 {
				t.Fatalf("pending = %v, want none", pending.Names())
			}
		})
	}
}

func TestRecordThenComputeAcrossApplies(t *testing.T) {
	ctx := context.Background()

	private := &fakePrivate{}

	// Apply 1: the user sets description and gateway.
	first := config(ctx, t, map[string]tftypes.Value{"name": str("n"), "description": str("d"), "gateway": str("g")})
	mustNoDiags(t, Record(ctx, first, private, testAttrs))

	applied := state(ctx, t, map[string]tftypes.Value{"name": str("n"), "description": str("d"), "gateway": str("g")})

	// Plan 2: the user removes description. Only it is pending.
	second := config(ctx, t, map[string]tftypes.Value{"name": str("n"), "gateway": str("g")})

	pending, diags := Compute(ctx, second, applied, private, testAttrs)
	mustNoDiags(t, diags)

	if got := strings.Join(pending.Names(), ","); got != "description" {
		t.Fatalf("plan 2 pending = %q, want %q", got, "description")
	}

	// Apply 2 records the new configuration: description is no longer set.
	mustNoDiags(t, Record(ctx, second, private, testAttrs))

	if got := string(private.data[PrivateKey]); got != `{"v":1,"attrs":["gateway"]}` {
		t.Fatalf("record after apply 2 = %s", got)
	}

	// Plan 3, nothing changed: the server may have put a description back,
	// but the user did not set one, so nothing is pending.
	cleared := state(ctx, t, map[string]tftypes.Value{"name": str("n"), "description": str("server"), "gateway": str("g")})

	pending, diags = Compute(ctx, second, cleared, private, testAttrs)
	mustNoDiags(t, diags)

	if len(pending) != 0 {
		t.Fatalf("plan 3 pending = %v, want none", pending.Names())
	}
}

func TestMarkPlansNullForClearAndUnknownForRelease(t *testing.T) {
	ctx := context.Background()

	p := plan(ctx, t, map[string]tftypes.Value{
		"name":        str("n"),
		"description": str("d"),
		"gateway":     str("g"),
		"cidr":        str("10.0.0.0/24"),
		"vlan_id":     num(7),
		"labels":      strSet("a"),
		"tenant_ids":  numSet(3),
	})

	pending := Pending{
		"description": {Name: "description", Mode: Clear},
		"gateway":     {Name: "gateway", Mode: Release},
		"vlan_id":     {Name: "vlan_id", Mode: Clear},
		"labels":      {Name: "labels", Mode: Clear},
		"tenant_ids":  {Name: "tenant_ids", Mode: Release},
	}

	mustNoDiags(t, Mark(ctx, &p, pending))

	get := func(name string) attr.Value {
		var v attr.Value

		mustNoDiags(t, p.GetAttribute(ctx, path.Root(name), &v))

		return v
	}

	for _, name := range []string{"description", "vlan_id", "labels"} {
		if v := get(name); !v.IsNull() {
			t.Errorf("%s = %v, want null", name, v)
		}
	}

	for _, name := range []string{"gateway", "tenant_ids"} {
		if v := get(name); !v.IsUnknown() {
			t.Errorf("%s = %v, want unknown", name, v)
		}
	}

	// Attributes that are not pending keep their planned value.
	if v, ok := get("cidr").(types.String); !ok || v.ValueString() != "10.0.0.0/24" {
		t.Errorf("cidr = %v, want the planned value", get("cidr"))
	}

	if v, ok := get("name").(types.String); !ok || v.ValueString() != "n" {
		t.Errorf("name = %v, want the planned value", get("name"))
	}
}

func TestMarkRejectsAnAttributeMissingFromTheSchema(t *testing.T) {
	ctx := context.Background()

	p := plan(ctx, t, map[string]tftypes.Value{"name": str("n")})

	diags := Mark(ctx, &p, Pending{"nope": {Name: "nope"}})
	if !diags.HasError() {
		t.Fatal("expected an error for an attribute that is not in the schema")
	}
}

func TestPlanOnDestroyDoesNothing(t *testing.T) {
	ctx := context.Background()

	// A destroy plan is null as a whole while the state is not; every
	// recorded attribute would otherwise look removed.
	nullPlan := tfsdk.Plan{Schema: testSchema(), Raw: tftypes.NewValue(testSchema().Type().TerraformType(ctx), nil)}

	req := resource.ModifyPlanRequest{
		Config: tfsdk.Config{Schema: testSchema(), Raw: nullPlan.Raw},
		State:  state(ctx, t, map[string]tftypes.Value{"name": str("n"), "description": str("d")}),
		Plan:   nullPlan,
	}
	resp := resource.ModifyPlanResponse{Plan: nullPlan}

	// The guard runs before the private state is read: req.Private is nil.
	pending, diags := Plan(ctx, req, &resp, testAttrs)
	mustNoDiags(t, diags)

	if len(pending) != 0 {
		t.Fatalf("pending = %v, want none on destroy", pending.Names())
	}

	if !resp.Plan.Raw.IsNull() {
		t.Fatal("the destroy plan was modified")
	}
}

func TestComputedNamesTheClearableAttributesThatAreComputed(t *testing.T) {
	got := Computed(testSchema().Attributes, testAttrs)

	if want := "description,gateway,labels,tenant_ids,vlan_id"; strings.Join(got, ",") != want {
		t.Fatalf("Computed = %q, want %q", strings.Join(got, ","), want)
	}

	if got := Computed(testSchema().Attributes, Attributes{{Name: "cidr"}, {Name: "ghost"}}); len(got) != 0 {
		t.Fatalf("Optional-only and unknown attributes are not computed, got %v", got)
	}
}

func TestPlanWithoutPrivateStateLeavesThePlanAlone(t *testing.T) {
	ctx := context.Background()

	// The framework hands ModifyPlan a typed nil *ProviderData when a
	// resource has no private state yet (first plan after import or
	// upgrade); Plan must treat that as "no record", not dereference it.
	req := resource.ModifyPlanRequest{
		Config: config(ctx, t, map[string]tftypes.Value{"name": str("n")}),
		State:  state(ctx, t, map[string]tftypes.Value{"name": str("n"), "description": str("d")}),
		Plan:   plan(ctx, t, map[string]tftypes.Value{"name": str("n"), "description": str("d")}),
	}
	resp := resource.ModifyPlanResponse{Plan: req.Plan}

	pending, diags := Plan(ctx, req, &resp, testAttrs)
	mustNoDiags(t, diags)

	if len(pending) != 0 {
		t.Fatalf("pending = %v, want none", pending.Names())
	}

	var v types.String

	mustNoDiags(t, resp.Plan.GetAttribute(ctx, path.Root("description"), &v))

	if v.ValueString() != "d" {
		t.Fatalf("description = %v, want the prior value untouched", v)
	}
}

func TestStringSendsEmptyToClearAndNothingWhenUnset(t *testing.T) {
	pending := Pending{"description": {Name: "description"}}

	if got := pending.String("description", types.StringNull()); got == nil || *got != "" {
		t.Errorf("clearing: got %v, want a pointer to the empty string", got)
	}

	if got := pending.String("gateway", types.StringValue("g")); got == nil || *got != "g" {
		t.Errorf("known: got %v, want a pointer to the value", got)
	}

	if got := pending.String("gateway", types.StringNull()); got != nil {
		t.Errorf("null and not clearing: got %q, want nil", *got)
	}

	if got := pending.String("gateway", types.StringUnknown()); got != nil {
		t.Errorf("unknown: got %q, want nil", *got)
	}
}

// fakeNullable stands in for the SDK's NullableString/NullableInt64: a value
// with a Set method that records what it was given.
type fakeNullable[T any] struct {
	calls []*T
}

func (f *fakeNullable[T]) Set(v *T) {
	f.calls = append(f.calls, v)
}

func TestNullableStringSetsNullToClearAndValueWhenKnown(t *testing.T) {
	pending := Pending{"description": {Name: "description"}}

	var target fakeNullable[string]

	if !NullableString(pending, "description", types.StringNull(), &target) {
		t.Error("clearing: want true")
	}

	if len(target.calls) != 1 || target.calls[0] != nil {
		t.Fatalf("clearing: calls = %v, want one Set(nil)", target.calls)
	}

	target = fakeNullable[string]{}

	if !NullableString(pending, "gateway", types.StringValue("g"), &target) {
		t.Error("known: want true")
	}

	if len(target.calls) != 1 || target.calls[0] == nil || *target.calls[0] != "g" {
		t.Fatalf("known: calls = %v, want one Set(&\"g\")", target.calls)
	}

	target = fakeNullable[string]{}

	if NullableString(pending, "gateway", types.StringNull(), &target) || len(target.calls) != 0 {
		t.Errorf("null and not clearing: want false and no Set, got calls %v", target.calls)
	}

	if NullableString(pending, "gateway", types.StringUnknown(), &target) || len(target.calls) != 0 {
		t.Errorf("unknown: want false and no Set, got calls %v", target.calls)
	}
}

func TestNullableInt64SetsNullToClearAndValueWhenKnown(t *testing.T) {
	pending := Pending{"vlan_id": {Name: "vlan_id"}}

	var target fakeNullable[int64]

	if !NullableInt64(pending, "vlan_id", types.Int64Null(), &target) || len(target.calls) != 1 || target.calls[0] != nil {
		t.Fatalf("clearing: calls = %v, want one Set(nil)", target.calls)
	}

	target = fakeNullable[int64]{}

	if !NullableInt64(pending, "other", types.Int64Value(7), &target) || len(target.calls) != 1 || *target.calls[0] != 7 {
		t.Fatalf("known: calls = %v, want one Set(&7)", target.calls)
	}

	target = fakeNullable[int64]{}

	if NullableInt64(pending, "other", types.Int64Null(), &target) || len(target.calls) != 0 {
		t.Errorf("null and not clearing: want false and no Set, got calls %v", target.calls)
	}
}

func TestStringsSendsAnEmptyListToClear(t *testing.T) {
	ctx := context.Background()

	pending := Pending{"labels": {Name: "labels"}}

	got, send, diags := pending.Strings(ctx, "labels", types.SetNull(types.StringType))
	mustNoDiags(t, diags)

	if !send || got == nil || len(got) != 0 {
		t.Errorf("clearing: got %#v send=%v, want an empty non-nil slice and true", got, send)
	}

	known, diags := types.SetValueFrom(ctx, types.StringType, []string{"a", "b"})
	mustNoDiags(t, diags)

	got, send, diags = pending.Strings(ctx, "other", known)
	mustNoDiags(t, diags)

	if !send || strings.Join(got, ",") != "a,b" {
		t.Errorf("known: got %v send=%v, want [a b] and true", got, send)
	}

	if _, send, _ := pending.Strings(ctx, "other", types.SetNull(types.StringType)); send {
		t.Error("null and not clearing: want false")
	}

	if _, send, _ := pending.Strings(ctx, "other", types.SetUnknown(types.StringType)); send {
		t.Error("unknown: want false")
	}
}

func TestInt64sSendsAnEmptyListToClear(t *testing.T) {
	ctx := context.Background()

	pending := Pending{"tenant_ids": {Name: "tenant_ids"}}

	got, send, diags := pending.Int64s(ctx, "tenant_ids", types.SetNull(types.Int64Type))
	mustNoDiags(t, diags)

	if !send || got == nil || len(got) != 0 {
		t.Errorf("clearing: got %#v send=%v, want an empty non-nil slice and true", got, send)
	}

	known, diags := types.SetValueFrom(ctx, types.Int64Type, []int64{3, 4})
	mustNoDiags(t, diags)

	got, send, diags = pending.Int64s(ctx, "other", known)
	mustNoDiags(t, diags)

	if !send || len(got) != 2 || got[0] != 3 || got[1] != 4 {
		t.Errorf("known: got %v send=%v, want [3 4] and true", got, send)
	}

	if _, send, _ := pending.Int64s(ctx, "other", types.SetNull(types.Int64Type)); send {
		t.Error("null and not clearing: want false")
	}
}

func TestInt64sAcceptsAListAsWell(t *testing.T) {
	ctx := context.Background()

	pending := Pending{"cloud_ids": {Name: "cloud_ids"}}

	got, send, diags := pending.Int64s(ctx, "cloud_ids", types.ListNull(types.Int64Type))
	mustNoDiags(t, diags)

	if !send || got == nil || len(got) != 0 {
		t.Errorf("clearing a list: got %#v send=%v, want an empty non-nil slice and true", got, send)
	}

	known, diags := types.ListValueFrom(ctx, types.Int64Type, []int64{17})
	mustNoDiags(t, diags)

	got, send, diags = pending.Int64s(ctx, "other", known)
	mustNoDiags(t, diags)

	if !send || len(got) != 1 || got[0] != 17 {
		t.Errorf("known list: got %v send=%v, want [17] and true", got, send)
	}
}

func TestCheckAcceptsACompleteClassification(t *testing.T) {
	if problems := Check(testSchema().Attributes, testAttrs, testNot); len(problems) != 0 {
		t.Fatalf("unexpected problems: %v", problems)
	}
}

func TestCheckReportsEveryClassificationMistake(t *testing.T) {
	attrs := testSchema().Attributes

	cases := map[string]struct {
		clearable Attributes
		not       NotClearable
		want      string
	}{
		"unclassified optional attribute": {
			clearable: testAttrs,
			not:       NotClearable{},
			want:      "enabled: optional attribute is neither clearable nor recorded as not clearable",
		},
		"in both lists": {
			clearable: testAttrs,
			not:       NotClearable{"enabled": "boolean", "description": "changed my mind"},
			want:      "description: both clearable and not clearable",
		},
		"declared twice": {
			clearable: append(Attributes{{Name: "description"}}, testAttrs...),
			not:       testNot,
			want:      "description: declared clearable twice",
		},
		"clearable but not in the schema": {
			clearable: append(Attributes{{Name: "ghost"}}, testAttrs...),
			not:       testNot,
			want:      "ghost: declared clearable but not in the schema",
		},
		"clearable write-only": {
			clearable: append(Attributes{{Name: "secret"}}, testAttrs...),
			not:       testNot,
			want:      "secret: write-only attributes are not clearable (they are never in state)",
		},
		"clearable required": {
			clearable: append(Attributes{{Name: "name"}}, testAttrs...),
			not:       testNot,
			want:      "name: a required attribute cannot be removed",
		},
		"clearable computed-only": {
			clearable: append(Attributes{{Name: "id"}}, testAttrs...),
			not:       testNot,
			want:      "id: a computed-only attribute is never configured",
		},
		"not clearable but not in the schema": {
			clearable: testAttrs,
			not:       NotClearable{"enabled": "boolean", "ghost": "gone"},
			want:      "ghost: recorded as not clearable but not in the schema",
		},
		"not clearable without a reason": {
			clearable: testAttrs,
			not:       NotClearable{"enabled": ""},
			want:      "enabled: not clearable without a reason",
		},
		"release on an attribute that is not computed": {
			clearable: append(Attributes{{Name: "cidr", Mode: Release}}, testAttrs[:2]...),
			not:       testNot,
			want:      "cidr: only a computed attribute can be planned unknown (Release)",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			problems := Check(attrs, tc.clearable, tc.not)

			for _, p := range problems {
				if p == tc.want {
					return
				}
			}

			t.Fatalf("problems = %v, want one to be %q", problems, tc.want)
		})
	}
}

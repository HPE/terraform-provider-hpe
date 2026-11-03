// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package wodigest

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var trackerSchema = schema.Schema{Attributes: map[string]schema.Attribute{
	"id":                  schema.Int64Attribute{Computed: true},
	"name":                schema.StringAttribute{Required: true},
	"password_wo":         schema.StringAttribute{Optional: true, WriteOnly: true},
	"password_wo_version": schema.Int64Attribute{Optional: true},
	"key_wo":              schema.StringAttribute{Optional: true, WriteOnly: true},
	"key_wo_version":      schema.Int64Attribute{Optional: true},
	"nested": schema.SingleNestedAttribute{Optional: true, Attributes: map[string]schema.Attribute{
		"secret": schema.StringAttribute{Optional: true, WriteOnly: true},
	}},
	RevisionAttribute: schema.Int64Attribute{Computed: true},
}}

func trackerValue(t *testing.T, attrs map[string]tftypes.Value) tftypes.Value {
	t.Helper()

	typ := trackerSchema.Type().TerraformType(context.Background()).(tftypes.Object)
	all := map[string]tftypes.Value{}

	for name, at := range typ.AttributeTypes {
		if v, ok := attrs[name]; ok {
			all[name] = v
		} else {
			all[name] = tftypes.NewValue(at, nil)
		}
	}

	return tftypes.NewValue(typ, all)
}

func str(s string) tftypes.Value { return tftypes.NewValue(tftypes.String, s) }
func num(n int64) tftypes.Value  { return tftypes.NewValue(tftypes.Number, n) }
func nested(s string) tftypes.Value {
	return tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"secret": tftypes.String}},
		map[string]tftypes.Value{"secret": str(s)})
}

var tracker = Tracker{
	Resource: "thing",
	Attributes: map[string]string{
		"password_wo":   "password_wo_version",
		"key_wo":        "key_wo_version",
		"nested.secret": "",
	},
}

func TestUnitTrackerConfigured(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cfg := tfsdk.Config{Schema: trackerSchema, Raw: trackerValue(t, map[string]tftypes.Value{
		"name": str("x"), "password_wo": str("Secret123!"), "nested": nested("deep"),
		// key_wo unset
	})}

	attrs, diags := tracker.configured(ctx, cfg)
	if diags.HasError() {
		t.Fatal(diags)
	}

	want := Attrs{"password_wo": sha("Secret123!"), "nested.secret": sha("deep")}
	if len(attrs) != len(want) {
		t.Fatalf("attrs = %v, want %v", attrs, want)
	}

	for k, v := range want {
		if attrs[k] != v {
			t.Errorf("%s = %s, want %s", k, attrs[k], v)
		}
	}

	// The digests equal what an importer would put in the import ID.
	if Compute(attrs, nil).All != Compute(Attrs{"password_wo": sha("Secret123!"), "nested.secret": sha("deep")}, nil).All {
		t.Fatal("configured digests must match import-ID digests")
	}
}

func TestUnitTrackerMarkerChanged(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mk := func(pw, key *int64) tfsdk.State {
		attrs := map[string]tftypes.Value{"name": str("x")}
		if pw != nil {
			attrs["password_wo_version"] = num(*pw)
		}

		if key != nil {
			attrs["key_wo_version"] = num(*key)
		}

		return tfsdk.State{Schema: trackerSchema, Raw: trackerValue(t, attrs)}
	}

	one, two := int64(1), int64(2)

	cases := []struct {
		name        string
		plan, state tfsdk.State
		want        bool
	}{
		{"both unset", mk(nil, nil), mk(nil, nil), false},
		{"same", mk(&one, &one), mk(&one, &one), false},
		{"bumped", mk(&two, nil), mk(&one, nil), true},
		{"added", mk(&one, nil), mk(nil, nil), true},
		{"removed", mk(nil, nil), mk(&one, nil), true},
		{"other marker bumped", mk(&one, &two), mk(&one, &one), true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			plan := tfsdk.Plan{Schema: tc.plan.Schema, Raw: tc.plan.Raw}

			got, diags := tracker.markerChanged(ctx, plan, tc.state)
			if diags.HasError() {
				t.Fatal(diags)
			}

			if got != tc.want {
				t.Fatalf("markerChanged = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestUnitTrackerNamesAndKey(t *testing.T) {
	t.Parallel()

	if tracker.Key() != "thing.wo_digests" {
		t.Fatalf("Key = %s", tracker.Key())
	}

	names := tracker.Names()
	if len(names) != 3 {
		t.Fatalf("Names = %v", names)
	}

	// The names are what ParseImportID accepts, nested ones as dotted paths.
	if _, attrs, err := ParseImportID("7#nested.secret="+sha("deep")+"&key_wo="+sha("k"), names); err != nil || len(attrs) != 2 {
		t.Fatalf("ParseImportID with the tracker's names: %v %v", attrs, err)
	}
}

func TestUnitUnknownOf(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	obj, _ := types.ObjectValue(map[string]attr.Type{"a": types.StringType}, map[string]attr.Value{"a": types.StringValue("x")})
	set, _ := types.SetValue(types.Int64Type, []attr.Value{types.Int64Value(1)})
	list, _ := types.ListValue(types.StringType, []attr.Value{types.StringValue("x")})
	m, _ := types.MapValue(types.BoolType, map[string]attr.Value{"k": types.BoolValue(true)})

	// Every kind, primitive and complex, yields the unknown of its own type:
	// a caller that names a collection or object computed attribute in
	// alsoUnknown must get a real unknown, or the apply is inconsistent.
	for name, v := range map[string]attr.Value{
		"string": types.StringValue("x"), "int64": types.Int64Value(1), "bool": types.BoolValue(true),
		"float64": types.Float64Value(1.5), "object": obj, "set": set, "list": list, "map": m,
		"dynamic": types.DynamicValue(types.StringValue("x")),
	} {
		u, err := unknownOf(ctx, v)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}

		if !u.IsUnknown() {
			t.Fatalf("%s: unknownOf returned a known value", name)
		}

		if !u.Type(ctx).Equal(v.Type(ctx)) {
			t.Fatalf("%s: type changed from %s to %s", name, v.Type(ctx), u.Type(ctx))
		}
	}

	if _, err := unknownOf(ctx, nil); err == nil {
		t.Fatal("nil must be an error, not a silent no-op")
	}
}

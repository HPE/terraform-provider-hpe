// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instance

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func tagSet(pairs ...string) types.Set {
	elems := make([]attr.Value, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		elems = append(elems, TagsValue{
			Name:  types.StringValue(pairs[i]),
			Value: types.StringValue(pairs[i+1]),
			state: attr.ValueStateKnown,
		})
	}

	return types.SetValueMust(TagsValue{}.Type(context.Background()), elems)
}

func apiTags(pairs ...string) []apiTag {
	out := make([]apiTag, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		n, v := pairs[i], pairs[i+1]
		out = append(out, apiTag{Name: &n, Value: &v})
	}

	return out
}

// TestUnitManagedTags: the tags written to state are the configured ones as
// Morpheus holds them. Tags Morpheus adds itself are never in state; after a
// write a configured tag Morpheus lacks or changed is kept as configured with
// a warning (the apply is consistent, MORPH-20973); on a refresh the API wins
// so drift shows; omitted stays null and empty stays empty (not null).
func TestUnitManagedTags(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	none := types.SetNull(TagsValue{}.Type(ctx))
	empty := tagSet()
	prod := tagSet("env", "prod")
	prodBackup := tagSet("env", "prod", "backup", "daily")

	cases := []struct {
		name     string
		ref      types.Set
		api      []apiTag
		refresh  bool
		want     types.Set
		warnings int
	}{
		{"omitted stays null whatever Morpheus holds", none, apiTags("env", "prod", "owner", "ops"), false, none, 0},
		{"omitted stays null on refresh", none, apiTags("owner", "ops"), true, none, 0},
		{"empty stays empty after a write", empty, nil, false, empty, 0},
		{"empty stays empty with Morpheus's own tags", empty, apiTags("owner", "ops"), false, empty, 0},
		{"empty stays empty on refresh", empty, apiTags("owner", "ops"), true, empty, 0},
		{"Morpheus holds exactly the configured", prod, apiTags("env", "prod"), false, prod, 0},
		{
			"4 sent, 6 held: the extras are not managed", prodBackup,
			apiTags("env", "prod", "backup", "daily", "owner", "ops", "costcenter", "42"), false, prodBackup, 0,
		},
		{
			"extras are not drift on refresh either", prodBackup,
			apiTags("backup", "daily", "env", "prod", "owner", "ops"), true, prodBackup, 0,
		},
		{"name case changed by Morpheus: configured casing kept", prod, apiTags("Env", "prod"), false, prod, 0},
		{"value case changed by the cloud (vCenter): configured casing kept, no warning", prod, apiTags("env", "Prod"), false, prod, 0},
		{"value case changed by the cloud: not drift on refresh", prod, apiTags("env", "Prod"), true, prod, 0},
		{"configured tag missing after the write: kept, warned", prod, nil, false, prod, 1},
		{"configured value overwritten after the write: kept, warned", prod, apiTags("env", "production"), false, prod, 1},
		{"configured tag reported with a null value after the write: kept, warned", prod, []apiTag{{Name: ptr("env")}}, false, prod, 1},
		{
			"configured tag reported with a null value on refresh: null value stored (drift)", prod,
			[]apiTag{{Name: ptr("env")}},
			true,
			types.SetValueMust(TagsValue{}.Type(ctx), []attr.Value{TagsValue{
				Name: types.StringValue("env"), Value: types.StringNull(), state: attr.ValueStateKnown,
			}}), 0,
		},
		{"configured tag missing on refresh: drift", prodBackup, apiTags("env", "prod"), true, prod, 0},
		{"configured value changed on refresh: drift", prod, apiTags("env", "production"), true, tagSet("env", "production"), 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, diags := managedTags(ctx, 42, tc.ref, tc.api, tc.refresh)
			if diags.HasError() {
				t.Fatalf("errors: %v", diags.Errors())
			}

			if !got.Equal(tc.want) {
				t.Errorf("managedTags() = %s, want %s", got, tc.want)
			}

			if got.IsNull() != tc.want.IsNull() {
				t.Errorf("null-ness: got %v want %v", got.IsNull(), tc.want.IsNull())
			}

			if n := diags.WarningsCount(); n != tc.warnings {
				t.Errorf("warnings = %d, want %d: %v", n, tc.warnings, diags.Warnings())
			}
		})
	}
}

// TestUnitTagChanges: an update sends only the difference between the managed
// tags in the prior state and the plan, so Morpheus's own tags survive.
func TestUnitTagChanges(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	none := types.SetNull(TagsValue{}.Type(ctx))

	names := func(add []string, remove []string) (string, string) {
		return join(add), join(remove)
	}

	cases := []struct {
		name       string
		prior      types.Set
		plan       types.Set
		wantAdd    []string // name=value
		wantRemove []string // name
	}{
		{"no change", tagSet("env", "prod"), tagSet("env", "prod"), nil, nil},
		{"value change", tagSet("env", "prod"), tagSet("env", "staging"), []string{"env=staging"}, nil},
		{"add one, remove one", tagSet("env", "prod", "old", "x"), tagSet("env", "prod", "new", "y"), []string{"new=y"}, []string{"old"}},
		{"first tags on an instance", none, tagSet("env", "prod"), []string{"env=prod"}, nil},
		{
			"tags removed from the configuration: only the managed ones go",
			tagSet("env", "prod", "backup", "daily"), none, nil,
			[]string{"backup", "env"},
		},
		{"emptied: same as removed", tagSet("env", "prod"), tagSet(), nil, []string{"env"}},
		{"case-only rename is not a change", tagSet("Env", "prod"), tagSet("env", "prod"), nil, nil},
		{"case-only value change is sent", tagSet("env", "Prod"), tagSet("env", "prod"), []string{"env=prod"}, nil},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			add, remove := tagChanges(tc.prior, tc.plan)

			var gotAdd, gotRemove []string
			for _, a := range add {
				gotAdd = append(gotAdd, *a.Name+"="+*a.Value)
			}
			for _, r := range remove {
				gotRemove = append(gotRemove, *r.Name)
			}

			wa, wr := names(tc.wantAdd, tc.wantRemove)
			if ga, gr := names(gotAdd, gotRemove); ga != wa || gr != wr {
				t.Errorf("tagChanges() add=%q remove=%q, want add=%q remove=%q", ga, gr, wa, wr)
			}
		})
	}
}

func join(s []string) string {
	out := ""
	for i, v := range s {
		if i > 0 {
			out += ","
		}
		out += v
	}

	return out
}

// tags_all carries everything Morpheus holds, and is an empty set -- never
// null -- when it holds nothing.
func TestUnitTagsAllFromAPI(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	all := tagsAllFromAPI(ctx, apiTags("env", "prod", "owner", "ops"))
	if all.IsNull() || len(all.Elements()) != 2 {
		t.Errorf("tags_all = %s, want two known elements", all)
	}

	if e := tagsAllFromAPI(ctx, nil); e.IsNull() || len(e.Elements()) != 0 {
		t.Errorf("tags_all with none = %s, want an empty known set", e)
	}
}

// A configuration whose tags repeat a name in any case is rejected: the two
// would be one tag to Morpheus.
func TestUnitUniqueTagNamesValidator(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	sch := InstanceResourceSchema(ctx)
	objType := sch.Type().TerraformType(ctx).(tftypes.Object)

	config := func(tags types.Set) tfsdk.Config {
		vals := make(map[string]tftypes.Value, len(objType.AttributeTypes))
		for name, ty := range objType.AttributeTypes {
			vals[name] = tftypes.NewValue(ty, nil)
		}

		raw, err := tags.ToTerraformValue(ctx)
		if err != nil {
			t.Fatal(err)
		}

		vals["tags"] = raw

		return tfsdk.Config{Schema: sch, Raw: tftypes.NewValue(objType, vals)}
	}

	run := func(tags types.Set) diag.Diagnostics {
		resp := &resource.ValidateConfigResponse{}
		uniqueTagNamesValidator{}.ValidateResource(ctx, resource.ValidateConfigRequest{Config: config(tags)}, resp)

		return resp.Diagnostics
	}

	if d := run(tagSet("env", "prod", "owner", "ops")); d.HasError() {
		t.Errorf("distinct names: %v", d)
	}

	if d := run(types.SetNull(TagsValue{}.Type(ctx))); d.HasError() {
		t.Errorf("no tags: %v", d)
	}

	d := run(tagSet("env", "prod", "Env", "staging"))
	if !d.HasError() || !strings.Contains(d.Errors()[0].Detail(), `"env" and "Env"`) {
		t.Errorf("duplicate ignoring case: want an error naming both, got %v", d)
	}
}

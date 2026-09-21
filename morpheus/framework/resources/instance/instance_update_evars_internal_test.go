// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instance

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// ev builds an evarSettings at the appliance defaults (export true, masked
// false) with the given value.
func ev(value string) evarSettings {
	return evarSettings{value: value, export: true, masked: false}
}

func evarsSet(t *testing.T, elems ...EvarsValue) basetypes.SetValue {
	t.Helper()

	ctx := context.Background()
	vals := make([]attr.Value, 0, len(elems))
	for _, e := range elems {
		e.state = attr.ValueStateKnown
		vals = append(vals, e)
	}

	set, diags := types.SetValue(EvarsValue{}.Type(ctx), vals)
	if diags.HasError() {
		t.Fatalf("building evars set: %v", diags)
	}

	return set
}

// entryStrings renders a request list as "name=value[,export=false][,masked=true]"
// so a test can assert the exact payload, including the option flags.
func entryStrings(list []sdk.AddEnvVariablesRequestEnvsInner) []string {
	out := make([]string, 0, len(list))
	for _, e := range list {
		s := e.Name + "=" + e.Value
		if e.Export != nil && !*e.Export {
			s += ",export=false"
		}
		if e.Masked != nil && *e.Masked {
			s += ",masked=true"
		}
		out = append(out, s)
	}

	return out
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

// TestDiffEvars covers the MORPH-612 / MORPH-610 update scenarios as pure map
// diffs -- removing every variable, modifying values while adding a new one,
// the no-change case -- plus the MORPH-17657 option changes: flipping export
// or masked on an otherwise-unchanged variable must upsert it, carrying the
// new flags.
func TestDiffEvars(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name         string
		planned      map[string]evarSettings
		current      map[string]evarSettings
		wantUpsert   []string
		wantRecreate []string
		wantRemove   []string
	}{
		{
			// MORPH-612: the evars block is removed entirely -> every current
			// variable must be deleted, nothing upserted.
			name:       "remove all",
			planned:    map[string]evarSettings{},
			current:    map[string]evarSettings{"ENV": ev("test"), "DEBUG": ev("true")},
			wantUpsert: []string{},
			wantRemove: []string{"DEBUG", "ENV"},
		},
		{
			// MORPH-610: two values change and one variable is added -> all
			// three are upserted (the API replaces by name), nothing removed.
			name:       "modify and add",
			planned:    map[string]evarSettings{"APP_ENV": ev("production"), "LOG_LEVEL": ev("debug"), "API_KEY": ev("test-key-123")},
			current:    map[string]evarSettings{"APP_ENV": ev("development"), "LOG_LEVEL": ev("info")},
			wantUpsert: []string{"API_KEY=test-key-123", "APP_ENV=production", "LOG_LEVEL=debug"},
			wantRemove: []string{},
		},
		{
			name:       "unchanged is a no-op",
			planned:    map[string]evarSettings{"A": ev("1"), "B": ev("2")},
			current:    map[string]evarSettings{"A": ev("1"), "B": ev("2")},
			wantUpsert: []string{},
			wantRemove: []string{},
		},
		{
			name:       "mixed: one changed, one removed, one added",
			planned:    map[string]evarSettings{"KEEP": ev("same"), "CHANGE": ev("new"), "ADD": ev("x")},
			current:    map[string]evarSettings{"KEEP": ev("same"), "CHANGE": ev("old"), "GONE": ev("y")},
			wantUpsert: []string{"ADD=x", "CHANGE=new"},
			wantRemove: []string{"GONE"},
		},
		{
			// MORPH-17657: same value, masked flipped on. The API only updates the
			// value of an existing variable, so an option change is a RECREATE
			// (delete + add), not an upsert.
			name:         "masked toggled on is a recreate",
			planned:      map[string]evarSettings{"SECRET": {value: "s3cr3t", export: true, masked: true}},
			current:      map[string]evarSettings{"SECRET": ev("s3cr3t")},
			wantUpsert:   []string{},
			wantRecreate: []string{"SECRET=s3cr3t,masked=true"},
			wantRemove:   []string{},
		},
		{
			// MORPH-17657: same value, export flipped off -> recreate, for the same reason.
			name:         "export toggled off is a recreate",
			planned:      map[string]evarSettings{"LOCAL": {value: "v", export: false, masked: false}},
			current:      map[string]evarSettings{"LOCAL": ev("v")},
			wantUpsert:   []string{},
			wantRecreate: []string{"LOCAL=v,export=false"},
			wantRemove:   []string{},
		},
		{
			// A default-settings upsert carries explicit export=true/masked=false
			// (rendered as plain name=value here) rather than omitting the flags.
			name:       "first-time set carries the flags",
			planned:    map[string]evarSettings{"A": ev("1")},
			current:    map[string]evarSettings{},
			wantUpsert: []string{"A=1"},
			wantRemove: []string{},
		},
		{
			// A NEW variable with options is a plain add: the add path does honour
			// export/masked, so no delete is needed.
			name:       "new variable with options is an upsert",
			planned:    map[string]evarSettings{"NEW": {value: "n", export: false, masked: true}},
			current:    map[string]evarSettings{},
			wantUpsert: []string{"NEW=n,export=false,masked=true"},
			wantRemove: []string{},
		},
		{
			// Value AND options change together -> still a recreate (options win).
			name:         "value and options changed together is a recreate",
			planned:      map[string]evarSettings{"X": {value: "new", export: true, masked: true}},
			current:      map[string]evarSettings{"X": ev("old")},
			wantUpsert:   []string{},
			wantRecreate: []string{"X=new,masked=true"},
			wantRemove:   []string{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := diffEvars(tc.planned, tc.current)
			if u := entryStrings(got.upsert); !equalStrings(u, tc.wantUpsert) {
				t.Errorf("upsert: got %v, want %v", u, tc.wantUpsert)
			}
			wantRecreate := tc.wantRecreate
			if wantRecreate == nil {
				wantRecreate = []string{}
			}
			if r := entryStrings(got.recreate); !equalStrings(r, wantRecreate) {
				t.Errorf("recreate: got %v, want %v", r, wantRecreate)
			}
			if !equalStrings(got.remove, tc.wantRemove) {
				t.Errorf("remove: got %v, want %v", got.remove, tc.wantRemove)
			}
			wantEmpty := len(tc.wantUpsert) == 0 && len(wantRecreate) == 0 && len(tc.wantRemove) == 0
			if got.empty() != wantEmpty {
				t.Errorf("empty(): got %v, want %v", got.empty(), wantEmpty)
			}
			// Every request must carry both flags explicitly, so the API never
			// falls back to its own defaults for a variable we are setting.
			for _, e := range append(got.upsert, got.recreate...) {
				if e.Export == nil || e.Masked == nil {
					t.Errorf("%s: export/masked must be set explicitly", e.Name)
				}
			}
		})
	}
}

// TestEvarsByName covers the framework-value side: a null or unknown set is
// "no variables" (so a removed block diffs as remove-all); a populated set
// flattens to name -> settings; and null export/masked resolve to the
// appliance defaults rather than to Go zero values.
func TestEvarsByName(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	elemType := EvarsValue{}.Type(ctx)

	t.Run("null set is empty", func(t *testing.T) {
		t.Parallel()
		got, diags := evarsByName(ctx, types.SetNull(elemType))
		if diags.HasError() || len(got) != 0 {
			t.Fatalf("got %v, diags %v", got, diags)
		}
	})

	t.Run("unknown set is empty", func(t *testing.T) {
		t.Parallel()
		got, diags := evarsByName(ctx, types.SetUnknown(elemType))
		if diags.HasError() || len(got) != 0 {
			t.Fatalf("got %v, diags %v", got, diags)
		}
	})

	t.Run("populated set flattens with explicit flags", func(t *testing.T) {
		t.Parallel()
		got, diags := evarsByName(ctx, evarsSet(t,
			EvarsValue{
				Name: types.StringValue("A"), Value: types.StringValue("1"),
				Export: types.BoolValue(false), Masked: types.BoolValue(true),
			},
			EvarsValue{
				Name: types.StringValue("B"), Value: types.StringValue("2"),
				Export: types.BoolValue(true), Masked: types.BoolValue(false),
			},
		))
		if diags.HasError() {
			t.Fatal(diags)
		}
		if got["A"] != (evarSettings{value: "1", export: false, masked: true}) {
			t.Errorf("A: got %+v", got["A"])
		}
		if got["B"] != ev("2") {
			t.Errorf("B: got %+v", got["B"])
		}
	})

	t.Run("null flags resolve to the appliance defaults", func(t *testing.T) {
		t.Parallel()
		got, diags := evarsByName(ctx, evarsSet(t,
			EvarsValue{Name: types.StringValue("A"), Value: types.StringValue("1"), Export: types.BoolNull(), Masked: types.BoolNull()},
		))
		if diags.HasError() {
			t.Fatal(diags)
		}
		if got["A"] != ev("1") {
			t.Errorf("expected export=true, masked=false defaults; got %+v", got["A"])
		}
	})
}

// TestIsAPIUpdateNeededIgnoresEvars pins the division of labour in Update:
// evars are not carried by the instance update request, so an evars-only
// change must NOT trigger that call -- it is applied through the dedicated
// environment-variable operations instead. Adding evars to the gate would
// fire an empty instance update (and its resize/wait machinery) for nothing.
// Ports never reach Update at all: they require replacement.
func TestIsAPIUpdateNeededIgnoresEvars(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	// Zero-valued framework sets (no element type) compare unequal to each other,
	// so every set attribute the gate inspects is given a typed null, as a real
	// plan and state would carry.
	base := InstanceModel{
		Name:               types.StringValue("i"),
		Description:        types.StringValue("d"),
		Tags:               types.SetNull(TagsValue{}.Type(ctx)),
		Volumes:            types.ListNull(VolumesValue{}.Type(ctx)),
		NetworkInterfaces:  types.ListNull(NetworkInterfacesValue{}.Type(ctx)),
		ServicePlanOptions: NewServicePlanOptionsValueNull(),
	}

	// Export/Masked are spelled out as explicit nulls for consistency with the
	// other fixtures (a zero-value BoolValue is already null; this is clarity,
	// not correctness).
	planned := base
	planned.Evars = evarsSet(t, EvarsValue{
		Name: types.StringValue("A"), Value: types.StringValue("1"),
		Export: types.BoolNull(), Masked: types.BoolNull(),
	})
	current := base
	current.Evars = evarsSet(t, EvarsValue{
		Name: types.StringValue("A"), Value: types.StringValue("2"),
		Export: types.BoolNull(), Masked: types.BoolNull(),
	})

	if isAPIUpdateNeeded(planned, current) {
		t.Fatal("an evars-only change must not trigger the instance update call")
	}

	// A change to a field the update request does carry still triggers it.
	planned.Description = types.StringValue("changed")
	if !isAPIUpdateNeeded(planned, current) {
		t.Fatal("a description change must trigger the instance update call")
	}
}

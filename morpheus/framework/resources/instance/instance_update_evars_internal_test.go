// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instance

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func evarsSet(t *testing.T, kv map[string]string) basetypes.SetValue {
	t.Helper()

	ctx := context.Background()
	elemType := EvarsValue{}.Type(ctx)
	elems := make([]attr.Value, 0, len(kv))
	for k, v := range kv {
		elems = append(elems, EvarsValue{
			Name:  types.StringValue(k),
			Value: types.StringValue(v),
			state: attr.ValueStateKnown,
		})
	}

	set, diags := types.SetValue(elemType, elems)
	if diags.HasError() {
		t.Fatalf("building evars set: %v", diags)
	}

	return set
}

func upserts(d evarsDiff) []string {
	out := make([]string, 0, len(d.upsert))
	for _, e := range d.upsert {
		out = append(out, e.Name+"="+e.Value)
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

// TestDiffEvars covers the MORPH-612 / MORPH-610 update scenarios as pure
// map diffs: removing every variable, modifying values while adding a new
// one, and the no-change case.
func TestDiffEvars(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		planned    map[string]string
		current    map[string]string
		wantUpsert []string
		wantRemove []string
	}{
		{
			// MORPH-612: the evars block is removed entirely -> every current
			// variable must be deleted, nothing upserted.
			name:       "remove all",
			planned:    map[string]string{},
			current:    map[string]string{"ENV": "test", "DEBUG": "true"},
			wantUpsert: []string{},
			wantRemove: []string{"DEBUG", "ENV"},
		},
		{
			// MORPH-610: two values change and one variable is added -> all
			// three are upserted (the API replaces by name), nothing removed.
			name:       "modify and add",
			planned:    map[string]string{"APP_ENV": "production", "LOG_LEVEL": "debug", "API_KEY": "test-key-123"},
			current:    map[string]string{"APP_ENV": "development", "LOG_LEVEL": "info"},
			wantUpsert: []string{"API_KEY=test-key-123", "APP_ENV=production", "LOG_LEVEL=debug"},
			wantRemove: []string{},
		},
		{
			name:       "unchanged is a no-op",
			planned:    map[string]string{"A": "1", "B": "2"},
			current:    map[string]string{"A": "1", "B": "2"},
			wantUpsert: []string{},
			wantRemove: []string{},
		},
		{
			name:       "mixed: one changed, one removed, one added",
			planned:    map[string]string{"KEEP": "same", "CHANGE": "new", "ADD": "x"},
			current:    map[string]string{"KEEP": "same", "CHANGE": "old", "GONE": "y"},
			wantUpsert: []string{"ADD=x", "CHANGE=new"},
			wantRemove: []string{"GONE"},
		},
		{
			name:       "first-time set on an instance with none",
			planned:    map[string]string{"A": "1"},
			current:    map[string]string{},
			wantUpsert: []string{"A=1"},
			wantRemove: []string{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := diffEvars(tc.planned, tc.current)
			if u := upserts(got); !equalStrings(u, tc.wantUpsert) {
				t.Errorf("upsert: got %v, want %v", u, tc.wantUpsert)
			}
			if !equalStrings(got.remove, tc.wantRemove) {
				t.Errorf("remove: got %v, want %v", got.remove, tc.wantRemove)
			}
			if got.empty() != (len(tc.wantUpsert) == 0 && len(tc.wantRemove) == 0) {
				t.Errorf("empty(): got %v", got.empty())
			}
		})
	}
}

// TestEvarsByName covers the framework-value side: a null or unknown set is
// "no variables" (so a removed block diffs as remove-all), and a populated set
// flattens to name -> value.
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

	t.Run("populated set flattens", func(t *testing.T) {
		t.Parallel()
		got, diags := evarsByName(ctx, evarsSet(t, map[string]string{"A": "1", "B": "2"}))
		if diags.HasError() {
			t.Fatal(diags)
		}
		if got["A"] != "1" || got["B"] != "2" || len(got) != 2 {
			t.Fatalf("got %v", got)
		}
	})
}

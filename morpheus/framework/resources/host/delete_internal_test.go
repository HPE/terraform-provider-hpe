// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package host

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func knownDeleteOptions(t *testing.T, attrs map[string]attr.Value) DeleteOptionsValue {
	t.Helper()
	full := map[string]attr.Value{
		"force":            types.BoolNull(),
		"preserve_volumes": types.BoolNull(),
		"remove_instances": types.BoolNull(),
		"remove_resources": types.BoolNull(),
	}
	for k, v := range attrs {
		full[k] = v
	}
	v, diags := NewDeleteOptionsValue(
		DeleteOptionsValue{}.AttributeTypes(context.Background()),
		full,
	)
	if diags.HasError() {
		t.Fatalf("building delete_options: %v", diags)
	}

	return v
}

func TestUnitDeleteQueryFlagsUnsetOmitted(t *testing.T) {
	t.Parallel()

	// Null block → no flags.
	if got := deleteQueryFlags(NewDeleteOptionsValueNull()); len(got) != 0 {
		t.Errorf("null delete_options: want no flags, got %v", got)
	}

	// Block with only remove_resources = true → only that flag.
	opts := knownDeleteOptions(t, map[string]attr.Value{
		"remove_resources": types.BoolValue(true),
	})
	got := deleteQueryFlags(opts)
	if len(got) != 1 || got["removeResources"] != "true" {
		t.Errorf("want {removeResources:true}, got %v", got)
	}
}

func TestUnitDeleteQueryFlagsAllSet(t *testing.T) {
	t.Parallel()

	opts := knownDeleteOptions(t, map[string]attr.Value{
		"remove_resources": types.BoolValue(false),
		"preserve_volumes": types.BoolValue(true),
		"force":            types.BoolValue(true),
		"remove_instances": types.BoolValue(false),
	})
	got := deleteQueryFlags(opts)
	want := map[string]string{
		"removeResources": "false",
		"preserveVolumes": "true",
		"force":           "true",
		"removeInstances": "false",
	}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("flag %q = %q, want %q", k, got[k], v)
		}
	}
}

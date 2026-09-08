// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package images

import (
	"context"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// compiled builds compiledFilter values directly, so a test can state the
// filter it means without going through the schema types.
func compiled(t *testing.T, field string, patterns ...string) []compiledFilter {
	t.Helper()

	res := make([]*regexp.Regexp, 0, len(patterns))

	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			t.Fatalf("test pattern %q does not compile: %v", p, err)
		}

		res = append(res, re)
	}

	return []compiledFilter{{field: field, res: res}}
}

// filterSet builds the schema-shaped set of filter blocks that compileFilters
// consumes, so the compilation path is exercised as the framework drives it.
func filterSet(
	t *testing.T,
	ctx context.Context,
	diags *diag.Diagnostics,
	field string,
	values []string,
) types.Set {
	t.Helper()

	vals := make([]attr.Value, 0, len(values))
	for _, v := range values {
		vals = append(vals, types.StringValue(v))
	}

	valueSet, d := types.SetValue(types.StringType, vals)
	diags.Append(d...)

	block, d := NewFilterValue(
		FilterValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"name":   types.StringValue(field),
			"values": valueSet,
		},
	)
	diags.Append(d...)

	set, d := types.SetValue(FilterValue{}.Type(ctx), []attr.Value{block})
	diags.Append(d...)

	return set
}

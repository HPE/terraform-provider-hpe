// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package serviceplan

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// TestUnitCreateOmitsUnknownOptionalScalars: the Optional+Computed scalars
// (max_cores, max_cpu, max_disks, sort_order) are unknown, not null, when the
// configuration leaves them out. The create request must omit them -- sending
// the 0 that ValueInt64Pointer() yields for an unknown made Morpheus store
// maxDisks: 0 and reject every instance on the plan (MORPH-18341).
func TestUnitCreateOmitsUnknownOptionalScalars(t *testing.T) {
	t.Parallel()

	plan := ServicePlanModel{
		Description:      types.StringUnknown(),
		MaxCores:         types.Int64Unknown(),
		MaxCpu:           types.Int64Unknown(),
		MaxDisks:         types.Int64Unknown(),
		SortOrder:        types.Int64Unknown(),
		CoresPerSocket:   types.Int64Unknown(),
		CustomCores:      types.BoolUnknown(),
		CustomCpu:        types.BoolUnknown(),
		CustomMaxMemory:  types.BoolUnknown(),
		CustomMaxStorage: types.BoolUnknown(),
		AddVolumes:       types.BoolUnknown(),
	}

	req := &sdk.AddServicePlansRequestServicePlan{}
	setOptionalScalarsInCreate(&plan, req)

	checks := map[string]bool{
		"description":      req.Description == nil,
		"maxCores":         req.MaxCores == nil,
		"maxCpu":           req.MaxCpu == nil,
		"maxDisks":         req.MaxDisks == nil,
		"sortOrder":        req.SortOrder == nil,
		"coresPerSocket":   req.CoresPerSocket == nil,
		"customCores":      req.CustomCores == nil,
		"customCpu":        req.CustomCpu == nil,
		"customMaxMemory":  req.CustomMaxMemory == nil,
		"customMaxStorage": req.CustomMaxStorage == nil,
		"addVolumes":       req.AddVolumes == nil,
	}
	for field, omitted := range checks {
		if !omitted {
			t.Errorf("%s: unknown planned value was sent in the create request; want omitted", field)
		}
	}
}

// TestUnitCreateOmitsNullOptionalScalars: null (explicitly absent, not
// computed) values are omitted too.
func TestUnitCreateOmitsNullOptionalScalars(t *testing.T) {
	t.Parallel()

	plan := ServicePlanModel{
		Description: types.StringNull(),
		MaxCores:    types.Int64Null(),
		MaxCpu:      types.Int64Null(),
		MaxDisks:    types.Int64Null(),
		SortOrder:   types.Int64Null(),
	}

	req := &sdk.AddServicePlansRequestServicePlan{}
	setOptionalScalarsInCreate(&plan, req)

	if req.Description != nil || req.MaxCores != nil || req.MaxCpu != nil ||
		req.MaxDisks != nil || req.SortOrder != nil {
		t.Errorf("null planned values were sent in the create request; want omitted: %+v", req)
	}
}

// TestUnitCreateSendsKnownOptionalScalars: known values, including an explicit
// 0, are sent as given.
func TestUnitCreateSendsKnownOptionalScalars(t *testing.T) {
	t.Parallel()

	plan := ServicePlanModel{
		Description:      types.StringValue("four cores"),
		MaxCores:         types.Int64Value(4),
		MaxCpu:           types.Int64Value(0),
		MaxDisks:         types.Int64Value(10),
		SortOrder:        types.Int64Value(1000),
		CoresPerSocket:   types.Int64Value(4),
		CustomCores:      types.BoolValue(false),
		CustomCpu:        types.BoolValue(true),
		CustomMaxMemory:  types.BoolValue(true),
		CustomMaxStorage: types.BoolValue(true),
		AddVolumes:       types.BoolValue(true),
	}

	req := &sdk.AddServicePlansRequestServicePlan{}
	setOptionalScalarsInCreate(&plan, req)

	want := map[string]struct{ got, want any }{
		"description":      {deref(req.Description), "four cores"},
		"maxCores":         {deref(req.MaxCores), int64(4)},
		"maxCpu":           {deref(req.MaxCpu), int64(0)},
		"maxDisks":         {deref(req.MaxDisks), int64(10)},
		"sortOrder":        {deref(req.SortOrder), int64(1000)},
		"coresPerSocket":   {deref(req.CoresPerSocket), int64(4)},
		"customCores":      {deref(req.CustomCores), false},
		"customCpu":        {deref(req.CustomCpu), true},
		"customMaxMemory":  {deref(req.CustomMaxMemory), true},
		"customMaxStorage": {deref(req.CustomMaxStorage), true},
		"addVolumes":       {deref(req.AddVolumes), true},
	}
	for field, v := range want {
		if v.got != v.want {
			t.Errorf("%s = %v, want %v", field, v.got, v.want)
		}
	}
}

// deref returns the pointed-to value, or nil for a nil pointer, so the tests
// can compare optional request fields in one place.
func deref[T any](p *T) any {
	if p == nil {
		return nil
	}

	return *p
}

// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package tenancy

import (
	"context"
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	fwschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// fakeChecker is a MasterChecker whose answer (and error) are fixed per test.
// It records whether CallerIsMaster was invoked so tests can assert that the
// guards short-circuit before making the whoami call.
type fakeChecker struct {
	isMaster bool
	err      error
	called   bool
}

func (f *fakeChecker) CallerIsMaster(_ context.Context) (bool, error) {
	f.called = true

	return f.isMaster, f.err
}

// visibilityModel mirrors the single attribute the guard reads.
type visibilityModel struct {
	Visibility types.String `tfsdk:"visibility"`
}

func visibilitySchema() fwschema.Schema {
	return fwschema.Schema{
		Attributes: map[string]fwschema.Attribute{
			"visibility": fwschema.StringAttribute{Optional: true},
		},
	}
}

func makeVisibilityPlan(t *testing.T, v types.String) tfsdk.Plan {
	t.Helper()
	ctx := context.Background()
	p := tfsdk.Plan{Schema: visibilitySchema()}
	if diags := p.Set(ctx, &visibilityModel{Visibility: v}); diags.HasError() {
		t.Fatalf("build plan: %v", diags)
	}

	return p
}

// TestUnitGuardVisibilityPublic_MORPH16419 covers the visibility guard's five
// cases: non-master+public errors, master+public passes, non-master+private
// passes, null/unknown short-circuits without calling whoami, and a checker
// error fails open with a warning.
func TestUnitGuardVisibilityPublic_MORPH16419(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		name       string
		visibility types.String
		checker    fakeChecker
		wantErrors int
		wantWarns  int
		wantCalled bool
	}{
		{
			name:       "non-master public is rejected",
			visibility: types.StringValue("public"),
			checker:    fakeChecker{isMaster: false},
			wantErrors: 1,
			wantCalled: true,
		},
		{
			name:       "master public is allowed",
			visibility: types.StringValue("public"),
			checker:    fakeChecker{isMaster: true},
			wantErrors: 0,
			wantCalled: true,
		},
		{
			name:       "non-master private is allowed",
			visibility: types.StringValue("private"),
			checker:    fakeChecker{isMaster: false},
			wantErrors: 0,
			wantCalled: false,
		},
		{
			name:       "null visibility short-circuits",
			visibility: types.StringNull(),
			checker:    fakeChecker{isMaster: false},
			wantErrors: 0,
			wantCalled: false,
		},
		{
			name:       "unknown visibility short-circuits",
			visibility: types.StringUnknown(),
			checker:    fakeChecker{isMaster: false},
			wantErrors: 0,
			wantCalled: false,
		},
		{
			name:       "checker error fails open with warning",
			visibility: types.StringValue("public"),
			checker:    fakeChecker{err: errors.New("boom")},
			wantErrors: 0,
			wantWarns:  1,
			wantCalled: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			plan := makeVisibilityPlan(t, tc.visibility)
			var diags diag.Diagnostics
			GuardVisibilityPublic(ctx, &tc.checker, plan, &diags)

			if got := diags.ErrorsCount(); got != tc.wantErrors {
				t.Errorf("errors = %d, want %d (%v)", got, tc.wantErrors, diags)
			}
			if got := diags.WarningsCount(); got != tc.wantWarns {
				t.Errorf("warnings = %d, want %d (%v)", got, tc.wantWarns, diags)
			}
			if tc.checker.called != tc.wantCalled {
				t.Errorf("checker called = %v, want %v", tc.checker.called, tc.wantCalled)
			}
		})
	}
}

// TestUnitGuardMasterOnlyAttribute_MORPH16546 covers the master-only attribute
// guard using appliance_name as the attribute under test.
func TestUnitGuardMasterOnlyAttribute_MORPH16546(t *testing.T) {
	ctx := context.Background()
	attr := path.Root("visibility") // reuse the single-attr schema

	cases := []struct {
		name       string
		value      types.String
		checker    fakeChecker
		wantErrors int
		wantWarns  int
		wantCalled bool
	}{
		{
			name:       "non-master set value is rejected",
			value:      types.StringValue("Acme Cloud Platform"),
			checker:    fakeChecker{isMaster: false},
			wantErrors: 1,
			wantCalled: true,
		},
		{
			name:       "master set value is allowed",
			value:      types.StringValue("Acme Cloud Platform"),
			checker:    fakeChecker{isMaster: true},
			wantErrors: 0,
			wantCalled: true,
		},
		{
			name:       "null value short-circuits",
			value:      types.StringNull(),
			checker:    fakeChecker{isMaster: false},
			wantErrors: 0,
			wantCalled: false,
		},
		{
			name:       "unknown value short-circuits",
			value:      types.StringUnknown(),
			checker:    fakeChecker{isMaster: false},
			wantErrors: 0,
			wantCalled: false,
		},
		{
			name:       "checker error fails open with warning",
			value:      types.StringValue("Acme Cloud Platform"),
			checker:    fakeChecker{err: errors.New("boom")},
			wantErrors: 0,
			wantWarns:  1,
			wantCalled: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			plan := makeVisibilityPlan(t, tc.value)
			var diags diag.Diagnostics
			GuardMasterOnlyAttribute(ctx, &tc.checker, plan, attr, &diags)

			if got := diags.ErrorsCount(); got != tc.wantErrors {
				t.Errorf("errors = %d, want %d (%v)", got, tc.wantErrors, diags)
			}
			if got := diags.WarningsCount(); got != tc.wantWarns {
				t.Errorf("warnings = %d, want %d (%v)", got, tc.wantWarns, diags)
			}
			if tc.checker.called != tc.wantCalled {
				t.Errorf("checker called = %v, want %v", tc.checker.called, tc.wantCalled)
			}
		})
	}
}

// TestUnitCheckVisibilityApplied_MORPH16419 covers the apply-time fallback.
func TestUnitCheckVisibilityApplied_MORPH16419(t *testing.T) {
	cases := []struct {
		name    string
		planned types.String
		actual  types.String
		wantErr bool
	}{
		{"public coerced to private errors", types.StringValue("public"), types.StringValue("private"), true},
		{"public stayed public passes", types.StringValue("public"), types.StringValue("public"), false},
		{"private passes", types.StringValue("private"), types.StringValue("private"), false},
		{"null planned passes", types.StringNull(), types.StringValue("private"), false},
		{"unknown planned passes", types.StringUnknown(), types.StringValue("private"), false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d := CheckVisibilityApplied(tc.planned, tc.actual)
			if (d != nil) != tc.wantErr {
				t.Errorf("diagnostic = %v, wantErr %v", d, tc.wantErr)
			}
			if d != nil && d.Severity() != diag.SeverityError {
				t.Errorf("severity = %v, want error", d.Severity())
			}
		})
	}
}

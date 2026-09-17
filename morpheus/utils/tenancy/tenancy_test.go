// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package tenancy

import (
	"context"
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	fwschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// guardClient builds a client whose whoami answer (or failure) is fixed per
// test, backed by a counting test server so tests can assert whether the guard
// consulted whoami at all. Each call builds a fresh client, so the package
// cache cannot leak answers between cases.
func guardClient(t *testing.T, isMaster bool, fail bool) (*sdk.APIClient, *atomic.Int32) {
	t.Helper()

	var hits, status atomic.Int32
	status.Store(http.StatusOK)
	if fail {
		status.Store(http.StatusInternalServerError)
	}

	body := `{}`
	if isMaster {
		body = `{"isMasterAccount": true}`
	}

	server := whoamiServer(t, &status, body, &hits)

	return newTestClient(server.URL), &hits
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

// TestUnitGuardVisibilityPublic_MORPH16419 covers the visibility guard's
// cases: non-master+public errors, master+public passes, non-master+private
// passes, null/unknown short-circuits without calling whoami, and a whoami
// failure fails open with a warning.
func TestUnitGuardVisibilityPublic_MORPH16419(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		name       string
		visibility types.String
		isMaster   bool
		fail       bool
		wantErrors int
		wantWarns  int
		wantCalled bool
	}{
		{
			name:       "non-master public is rejected",
			visibility: types.StringValue("public"),
			isMaster:   false,
			wantErrors: 1,
			wantCalled: true,
		},
		{
			name:       "master public is allowed",
			visibility: types.StringValue("public"),
			isMaster:   true,
			wantErrors: 0,
			wantCalled: true,
		},
		{
			name:       "non-master private is allowed",
			visibility: types.StringValue("private"),
			isMaster:   false,
			wantErrors: 0,
			wantCalled: false,
		},
		{
			name:       "null visibility short-circuits",
			visibility: types.StringNull(),
			isMaster:   false,
			wantErrors: 0,
			wantCalled: false,
		},
		{
			name:       "unknown visibility short-circuits",
			visibility: types.StringUnknown(),
			isMaster:   false,
			wantErrors: 0,
			wantCalled: false,
		},
		{
			name:       "whoami failure fails open with warning",
			visibility: types.StringValue("public"),
			fail:       true,
			wantErrors: 0,
			wantWarns:  1,
			wantCalled: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resetCache(t)

			client, hits := guardClient(t, tc.isMaster, tc.fail)
			plan := makeVisibilityPlan(t, tc.visibility)
			var diags diag.Diagnostics
			GuardVisibilityPublic(ctx, client, plan, &diags)

			if got := diags.ErrorsCount(); got != tc.wantErrors {
				t.Errorf("errors = %d, want %d (%v)", got, tc.wantErrors, diags)
			}
			if got := diags.WarningsCount(); got != tc.wantWarns {
				t.Errorf("warnings = %d, want %d (%v)", got, tc.wantWarns, diags)
			}
			if called := hits.Load() > 0; called != tc.wantCalled {
				t.Errorf("whoami called = %v, want %v", called, tc.wantCalled)
			}
		})
	}
}

// TestUnitGuardVisibilityPublicNilClient verifies a nil client (provider never
// configured, or client construction failed) fails open with a warning rather
// than panicking or blocking the plan.
func TestUnitGuardVisibilityPublicNilClient(t *testing.T) {
	resetCache(t)

	plan := makeVisibilityPlan(t, types.StringValue("public"))
	var diags diag.Diagnostics
	GuardVisibilityPublic(context.Background(), nil, plan, &diags)

	if got := diags.ErrorsCount(); got != 0 {
		t.Errorf("errors = %d, want 0 (%v)", got, diags)
	}
	if got := diags.WarningsCount(); got != 1 {
		t.Errorf("warnings = %d, want 1 (%v)", got, diags)
	}
}

// TestUnitGuardMasterOnlyAttribute_MORPH16546 covers the master-only attribute
// guard using a single-attribute schema.
func TestUnitGuardMasterOnlyAttribute_MORPH16546(t *testing.T) {
	ctx := context.Background()
	attr := path.Root("visibility") // reuse the single-attr schema

	cases := []struct {
		name       string
		value      types.String
		isMaster   bool
		fail       bool
		wantErrors int
		wantWarns  int
		wantCalled bool
	}{
		{
			name:       "non-master set value is rejected",
			value:      types.StringValue("Acme Cloud Platform"),
			isMaster:   false,
			wantErrors: 1,
			wantCalled: true,
		},
		{
			name:       "master set value is allowed",
			value:      types.StringValue("Acme Cloud Platform"),
			isMaster:   true,
			wantErrors: 0,
			wantCalled: true,
		},
		{
			name:       "null value short-circuits",
			value:      types.StringNull(),
			isMaster:   false,
			wantErrors: 0,
			wantCalled: false,
		},
		{
			name:       "unknown value short-circuits",
			value:      types.StringUnknown(),
			isMaster:   false,
			wantErrors: 0,
			wantCalled: false,
		},
		{
			name:       "whoami failure fails open with warning",
			value:      types.StringValue("Acme Cloud Platform"),
			fail:       true,
			wantErrors: 0,
			wantWarns:  1,
			wantCalled: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resetCache(t)

			client, hits := guardClient(t, tc.isMaster, tc.fail)
			plan := makeVisibilityPlan(t, tc.value)
			var diags diag.Diagnostics
			GuardMasterOnlyAttribute(ctx, client, plan, attr, &diags)

			if got := diags.ErrorsCount(); got != tc.wantErrors {
				t.Errorf("errors = %d, want %d (%v)", got, tc.wantErrors, diags)
			}
			if got := diags.WarningsCount(); got != tc.wantWarns {
				t.Errorf("warnings = %d, want %d (%v)", got, tc.wantWarns, diags)
			}
			if called := hits.Load() > 0; called != tc.wantCalled {
				t.Errorf("whoami called = %v, want %v", called, tc.wantCalled)
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

// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package compute

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	morpheus "github.com/HPE/terraform-provider-hpe/internal/sdk/legacy"
)

// TestUnitReconcileTenantIDs_MORPH8207 covers the tenant_ids reconciliation
// that stops the perpetual plan on hpe_morpheus_resource_pool_group. The API
// always reports the owning tenant alongside the shared tenants, so the owner
// is dropped from state unless the user listed it, while real changes made
// outside Terraform (added or removed tenants) still surface as drift.
func TestUnitReconcileTenantIDs_MORPH8207(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		prior []int
		api   []int64
		owner int64
		want  []int
	}{
		{
			name:  "tenant_ids unset: injected owner is dropped",
			prior: nil,
			api:   []int64{1},
			owner: 1,
			want:  []int{},
		},
		{
			name:  "shared tenant without owner in config",
			prior: []int{2},
			api:   []int64{1, 2},
			owner: 1,
			want:  []int{2},
		},
		{
			name:  "owner listed explicitly is kept",
			prior: []int{1, 2},
			api:   []int64{1, 2},
			owner: 1,
			want:  []int{1, 2},
		},
		{
			name:  "owner alone listed explicitly is kept",
			prior: []int{1},
			api:   []int64{1},
			owner: 1,
			want:  []int{1},
		},
		{
			name:  "tenant added outside Terraform surfaces as drift",
			prior: []int{2},
			api:   []int64{1, 2, 3},
			owner: 1,
			want:  []int{2, 3},
		},
		{
			name:  "tenant removed outside Terraform surfaces as drift",
			prior: []int{2, 3},
			api:   []int64{1, 2},
			owner: 1,
			want:  []int{2},
		},
		{
			name:  "import (no prior) yields the shared tenants only",
			prior: nil,
			api:   []int64{1, 5, 6},
			owner: 1,
			want:  []int{5, 6},
		},
		{
			name:  "no tenants reported",
			prior: []int{2},
			api:   nil,
			owner: 1,
			want:  []int{},
		},
		{
			name:  "sub-tenant owner is dropped too",
			prior: nil,
			api:   []int64{42},
			owner: 42,
			want:  []int{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := reconcileTenantIDs(tc.prior, tc.api, tc.owner)
			if got == nil {
				t.Fatal("reconcileTenantIDs returned nil; want a non-nil slice so an empty set is stored explicitly")
			}
			if !slices.Equal(got, tc.want) {
				t.Errorf("reconcileTenantIDs(%v, %v, %d) = %v, want %v", tc.prior, tc.api, tc.owner, got, tc.want)
			}
		})
	}
}

// TestUnitResourcePoolGroupRead_OwnerLookup_MORPH8207 drives Read against a
// fake API with the owner lookup stubbed. When the lookup succeeds the owner
// is reconciled out of tenant_ids (and kept when configured). When it fails,
// Read still completes with every other attribute set, keeps the API's full
// list and raises a warning rather than an error, so refresh, plan and import
// of a resource pool group never depend on whoami being reachable.
//
// Not parallel: it swaps the package-level resourcePoolGroupOwnerID seam.
func TestUnitResourcePoolGroupRead_OwnerLookup_MORPH8207(t *testing.T) {
	const groupJSON = `{"resourcePoolGroup":{"id":7,"name":"g","description":"d",` +
		`"mode":"roundrobin","visibility":"private","pools":[1],` +
		`"tenants":[{"id":1,"name":"owner"},{"id":2,"name":"shared"}],` +
		`"resourcePermission":{"all":true,"sites":[]}}}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/resource-pools/groups/7" {
			http.NotFound(w, r)

			return
		}
		fmt.Fprint(w, groupJSON)
	}))
	defer srv.Close()

	client := morpheus.NewClient(srv.URL, morpheus.SkipLogin()).SetAccessToken("abc", "", 0, "")

	cases := []struct {
		name        string
		lookup      func(*morpheus.Client) (int64, error)
		priorTenant []any
		wantTenants []int
		wantWarning bool
	}{
		{
			name:        "owner reconciled out when not configured",
			lookup:      func(*morpheus.Client) (int64, error) { return 1, nil },
			wantTenants: []int{2},
		},
		{
			name:        "owner kept when configured",
			lookup:      func(*morpheus.Client) (int64, error) { return 1, nil },
			priorTenant: []any{1, 2},
			wantTenants: []int{1, 2},
		},
		{
			name:        "lookup failure keeps the API list and warns",
			lookup:      func(*morpheus.Client) (int64, error) { return 0, errors.New("whoami: 503") },
			wantTenants: []int{1, 2},
			wantWarning: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			orig := resourcePoolGroupOwnerID
			resourcePoolGroupOwnerID = tc.lookup
			t.Cleanup(func() { resourcePoolGroupOwnerID = orig })

			res := ResourceResourcePoolGroup()
			raw := map[string]any{"name": "g", "mode": "roundrobin", "resource_pool_ids": []any{1}}
			if tc.priorTenant != nil {
				raw["tenant_ids"] = tc.priorTenant
			}
			d := schema.TestResourceDataRaw(t, res.Schema, raw)
			d.SetId("7")

			diags := res.ReadContext(context.Background(), d, client)
			if diags.HasError() {
				t.Fatalf("Read returned an error: %v", diags)
			}

			var warned bool
			for _, dg := range diags {
				if dg.Severity == diag.Warning {
					warned = true
				}
			}
			if warned != tc.wantWarning {
				t.Errorf("warning raised = %v, want %v (diags: %v)", warned, tc.wantWarning, diags)
			}

			// Read must have populated the rest of the resource as well, not
			// just tenant_ids.
			if got := d.Get("visibility"); got != "private" {
				t.Errorf("visibility = %v, want private", got)
			}
			if got := d.Get("resource_pool_ids").(*schema.Set).List(); len(got) != 1 || got[0] != 1 {
				t.Errorf("resource_pool_ids = %v, want [1]", got)
			}

			var got []int
			for _, v := range d.Get("tenant_ids").(*schema.Set).List() {
				got = append(got, v.(int))
			}
			slices.Sort(got)
			if !slices.Equal(got, tc.wantTenants) {
				t.Errorf("tenant_ids = %v, want %v", got, tc.wantTenants)
			}
		})
	}
}

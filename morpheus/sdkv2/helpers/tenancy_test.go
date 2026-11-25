// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package helpers_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	sdklegacy "github.com/HPE/terraform-provider-hpe/internal/sdk/legacy"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/client"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"
)

// newTestClient builds a legacy client (token auth, login skipped) pointed at a
// whoami test server. The master-tenant determination is cached process-wide,
// so each case resets it via helpers.ResetTenancyCache.
func newTestClient(t *testing.T, url string) *sdklegacy.Client {
	t.Helper()

	return client.NewLegacyClient(
		context.Background(),
		url,
		"",
		"",
		"",
		"abc",
		sdklegacy.SkipLogin(),
	)
}

// TestUnitVisibilityCustomizeDiff_MORPH16419 verifies the SDKv2 visibility
// CustomizeDiff: a master caller (isMasterAccount present-true) allows public,
// a sub-tenant caller (isMasterAccount absent) rejects it, and a non-public /
// unknown value skips the whoami call entirely.
func TestUnitVisibilityCustomizeDiff_MORPH16419(t *testing.T) {
	cases := []struct {
		name       string
		whoamiBody string
		visibility any
		wantErr    bool
		wantCall   bool
	}{
		{
			name:       "master public allowed",
			whoamiBody: `{"isMasterAccount":true}`,
			visibility: "public",
			wantErr:    false,
			wantCall:   true,
		},
		{
			name:       "sub-tenant public rejected",
			whoamiBody: `{}`,
			visibility: "public",
			wantErr:    true,
			wantCall:   true,
		},
		{
			name:       "private skips whoami",
			whoamiBody: `{}`,
			visibility: "private",
			wantErr:    false,
			wantCall:   false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// The tenancy determination is cached once per process; reset it so
			// each case exercises its own whoami answer.
			helpers.ResetTenancyCache()
			t.Cleanup(helpers.ResetTenancyCache)

			var called bool
			srv := httptest.NewServer(http.HandlerFunc(
				func(w http.ResponseWriter, _ *http.Request) {
					called = true
					fmt.Fprint(w, tc.whoamiBody)
				}))
			defer srv.Close()

			cl := newTestClient(t, srv.URL)

			res := &schema.Resource{
				Schema: map[string]*schema.Schema{
					"visibility": {Type: schema.TypeString, Optional: true},
				},
				CustomizeDiff: helpers.VisibilityCustomizeDiff,
			}

			diff, err := res.Diff(
				context.Background(),
				nil,
				terraform.NewResourceConfigRaw(
					map[string]any{"visibility": tc.visibility},
				),
				cl,
			)
			_ = diff

			if (err != nil) != tc.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tc.wantErr)
			}
			if called != tc.wantCall {
				t.Errorf("whoami called = %v, want %v", called, tc.wantCall)
			}
		})
	}
}

// TestUnitCallerAccountID_MORPH8207 verifies that the caller's tenant ID is
// read from whoami's user.accountId, that it shares the process-wide cache
// with CallerIsMaster (one whoami call serves both), and that a whoami answer
// without an account is reported as an error rather than as tenant 0.
func TestUnitCallerAccountID_MORPH8207(t *testing.T) {
	cases := []struct {
		name       string
		whoamiBody string
		wantID     int64
		wantMaster bool
		wantErr    bool
	}{
		{
			name:       "master caller",
			whoamiBody: `{"user":{"id":7,"accountId":1,"username":"admin"},"isMasterAccount":true}`,
			wantID:     1,
			wantMaster: true,
		},
		{
			name:       "sub-tenant caller",
			whoamiBody: `{"user":{"id":9,"accountId":42,"username":"sub"}}`,
			wantID:     42,
			wantMaster: false,
		},
		{
			name:       "no account reported",
			whoamiBody: `{"user":{"id":9,"username":"sub"}}`,
			wantErr:    true,
		},
		{
			name:       "no user reported",
			whoamiBody: `{"isMasterAccount":true}`,
			wantMaster: true,
			wantErr:    true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			helpers.ResetTenancyCache()
			t.Cleanup(helpers.ResetTenancyCache)

			var calls int
			srv := httptest.NewServer(http.HandlerFunc(
				func(w http.ResponseWriter, _ *http.Request) {
					calls++
					fmt.Fprint(w, tc.whoamiBody)
				}))
			defer srv.Close()

			cl := newTestClient(t, srv.URL)

			id, err := helpers.CallerAccountID(cl)
			if (err != nil) != tc.wantErr {
				t.Fatalf("CallerAccountID err = %v, wantErr %v", err, tc.wantErr)
			}
			if id != tc.wantID {
				t.Errorf("CallerAccountID = %d, want %d", id, tc.wantID)
			}

			// A second lookup, and the master determination, must reuse the
			// cached whoami answer.
			if _, err := helpers.CallerAccountID(cl); (err != nil) != tc.wantErr {
				t.Fatalf("second CallerAccountID err = %v, wantErr %v", err, tc.wantErr)
			}
			isMaster, err := helpers.CallerIsMaster(cl)
			if err != nil {
				t.Fatalf("CallerIsMaster err = %v", err)
			}
			if isMaster != tc.wantMaster {
				t.Errorf("CallerIsMaster = %v, want %v", isMaster, tc.wantMaster)
			}
			if calls != 1 {
				t.Errorf("whoami called %d times, want 1", calls)
			}
		})
	}
}

// TestUnitTenancyCachePerClient_MORPH8207 verifies that the tenancy cache is
// keyed by client: two clients with different credentials resolve to their own
// whoami answers, each answer is fetched exactly once however many lookups are
// made through either entry point, and resetting the cache makes both resolve
// again. A process-wide cache would hand the first client's tenant to the
// second, which in resource_pool_group Read means filtering the wrong owner out
// of tenant_ids.
//
// Not parallel: it shares the package-level cache with the other tenancy tests.
func TestUnitTenancyCachePerClient_MORPH8207(t *testing.T) {
	helpers.ResetTenancyCache()
	t.Cleanup(helpers.ResetTenancyCache)

	type server struct {
		body  string
		calls int
	}
	master := &server{body: `{"user":{"id":7,"accountId":1,"username":"admin"},"isMasterAccount":true}`}
	sub := &server{body: `{"user":{"id":9,"accountId":42,"username":"sub"}}`}

	serve := func(s *server) *httptest.Server {
		return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			s.calls++
			fmt.Fprint(w, s.body)
		}))
	}
	masterSrv := serve(master)
	defer masterSrv.Close()
	subSrv := serve(sub)
	defer subSrv.Close()

	masterClient := newTestClient(t, masterSrv.URL)
	subClient := newTestClient(t, subSrv.URL)

	// Interleave lookups across both clients and both entry points.
	for range 3 {
		id, err := helpers.CallerAccountID(masterClient)
		if err != nil || id != 1 {
			t.Fatalf("master CallerAccountID = %d, %v; want 1, nil", id, err)
		}
		id, err = helpers.CallerAccountID(subClient)
		if err != nil || id != 42 {
			t.Fatalf("sub CallerAccountID = %d, %v; want 42, nil", id, err)
		}
		isMaster, err := helpers.CallerIsMaster(masterClient)
		if err != nil || !isMaster {
			t.Fatalf("master CallerIsMaster = %v, %v; want true, nil", isMaster, err)
		}
		isMaster, err = helpers.CallerIsMaster(subClient)
		if err != nil || isMaster {
			t.Fatalf("sub CallerIsMaster = %v, %v; want false, nil", isMaster, err)
		}
	}

	if master.calls != 1 || sub.calls != 1 {
		t.Errorf("whoami calls after repeated lookups: master=%d sub=%d, want 1 and 1", master.calls, sub.calls)
	}

	// Forgetting the cache makes each client resolve again, still independently.
	helpers.ResetTenancyCache()

	if id, err := helpers.CallerAccountID(subClient); err != nil || id != 42 {
		t.Fatalf("sub CallerAccountID after reset = %d, %v; want 42, nil", id, err)
	}
	if id, err := helpers.CallerAccountID(masterClient); err != nil || id != 1 {
		t.Fatalf("master CallerAccountID after reset = %d, %v; want 1, nil", id, err)
	}
	if master.calls != 2 || sub.calls != 2 {
		t.Errorf("whoami calls after reset: master=%d sub=%d, want 2 and 2", master.calls, sub.calls)
	}
}

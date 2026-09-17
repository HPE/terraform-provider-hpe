// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package tenancy

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// resetCache clears the process-wide tenancy determination, now and at test
// cleanup, so each test exercises its own whoami answer. It routes through the
// exported ResetTenancyCache so there is a single reset implementation.
//
// Tests that use this MUST remain serial: the tenancy cache is a package-level
// global, so calling t.Parallel() in any test here (in caller_test.go or
// tenancy_test.go) would race the reset against another test's CallerIsMaster
// and produce flaky, order-dependent failures.
func resetCache(t *testing.T) {
	t.Helper()

	ResetTenancyCache()
	t.Cleanup(ResetTenancyCache)
}

// newTestClient builds a bare SDK client aimed at a test server. It
// deliberately does not use clientfactory: the tenancy package owns the whole
// tenancy mechanism and must not depend on it.
func newTestClient(url string) *sdk.APIClient {
	cfg := sdk.NewConfiguration()
	cfg.Servers[0].URL = url

	return sdk.NewAPIClient(cfg)
}

// whoamiServer serves /api/whoami with the given body and status, counting
// requests so tests can assert cache behaviour.
func whoamiServer(t *testing.T, status *atomic.Int32, body string, hits *atomic.Int32) *httptest.Server {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/api/whoami" {
				w.WriteHeader(http.StatusNotFound)

				return
			}
			hits.Add(1)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(int(status.Load()))
			_, _ = w.Write([]byte(body))
		}))
	t.Cleanup(server.Close)

	return server
}

// TestUnitCallerIsMasterCachesSuccess verifies a successful determination is
// cached: the second call is answered without a second whoami request.
func TestUnitCallerIsMasterCachesSuccess(t *testing.T) {
	resetCache(t)

	var hits, status atomic.Int32
	status.Store(http.StatusOK)
	server := whoamiServer(t, &status, `{"isMasterAccount": true}`, &hits)

	client := newTestClient(server.URL)
	ctx := context.Background()

	for range 2 {
		isMaster, err := CallerIsMaster(ctx, client)
		if err != nil {
			t.Fatalf("CallerIsMaster: %v", err)
		}
		if !isMaster {
			t.Fatal("CallerIsMaster = false, want true")
		}
	}

	if got := hits.Load(); got != 1 {
		t.Fatalf("whoami requests = %d, want 1 (second call must be a cache hit)", got)
	}
}

// TestUnitCallerIsMasterAbsentFieldMeansNonMaster verifies the whoami
// isMasterAccount omission semantics: the server only serializes the field when
// it is true, so its absence means non-master.
func TestUnitCallerIsMasterAbsentFieldMeansNonMaster(t *testing.T) {
	resetCache(t)

	var hits, status atomic.Int32
	status.Store(http.StatusOK)
	server := whoamiServer(t, &status, `{}`, &hits)

	isMaster, err := CallerIsMaster(context.Background(), newTestClient(server.URL))
	if err != nil {
		t.Fatalf("CallerIsMaster: %v", err)
	}
	if isMaster {
		t.Fatal("CallerIsMaster = true for a response omitting isMasterAccount, want false")
	}
}

// TestUnitCallerIsMasterDoesNotCacheErrors verifies a failed determination is
// retried: a transient whoami failure must not poison later checks.
func TestUnitCallerIsMasterDoesNotCacheErrors(t *testing.T) {
	resetCache(t)

	var hits, status atomic.Int32
	status.Store(http.StatusInternalServerError)
	server := whoamiServer(t, &status, `{"isMasterAccount": true}`, &hits)

	client := newTestClient(server.URL)
	ctx := context.Background()

	if _, err := CallerIsMaster(ctx, client); err == nil {
		t.Fatal("CallerIsMaster: want error on whoami failure")
	}

	status.Store(http.StatusOK)

	isMaster, err := CallerIsMaster(ctx, client)
	if err != nil {
		t.Fatalf("CallerIsMaster after recovery: %v", err)
	}
	if !isMaster {
		t.Fatal("CallerIsMaster = false after recovery, want true")
	}
	if got := hits.Load(); got != 2 {
		t.Fatalf("whoami requests = %d, want 2 (error must not be cached)", got)
	}
}

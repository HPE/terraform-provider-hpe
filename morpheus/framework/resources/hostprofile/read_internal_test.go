// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package hostprofile

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// newTestClient builds a bare SDK client aimed at a test server.
func newTestClient(url string) *sdk.APIClient {
	cfg := sdk.NewConfiguration()
	cfg.Servers[0].URL = url

	return sdk.NewAPIClient(cfg)
}

// hostProfileServer serves GET /api/host-profiles/{id} with the given status
// and body.
func hostProfileServer(t *testing.T, status int, body string) *httptest.Server {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)
			_, _ = w.Write([]byte(body))
		}))
	t.Cleanup(server.Close)

	return server
}

// TestUnitGetHostProfileAsStateNotFound verifies that a 404 (the profile was
// deleted out of band) is reported as a clean "not found" -- found=false with
// no diagnostics -- so the caller removes it from state rather than surfacing a
// refresh error. The generated SDK returns a non-nil error for any status >=
// 300, so this exercises the ordering of the IsNotFound check relative to the
// transport-error check.
func TestUnitGetHostProfileAsStateNotFound(t *testing.T) {
	server := hostProfileServer(t, http.StatusNotFound, `{}`)
	client := newTestClient(server.URL)

	_, found, diags := getHostProfileAsState(
		context.Background(), 123, client, HostProfileModel{},
	)

	if diags.HasError() {
		t.Fatalf("expected no diagnostics on 404, got: %v", diags.Errors())
	}
	if found {
		t.Fatal("expected found=false on 404")
	}
}

// TestUnitGetHostProfileAsStateServerError verifies that a genuine transport /
// server error (500) still surfaces as a diagnostic with found=false, so the
// reordered IsNotFound check does not swallow real failures.
func TestUnitGetHostProfileAsStateServerError(t *testing.T) {
	server := hostProfileServer(t, http.StatusInternalServerError, `{}`)
	client := newTestClient(server.URL)

	_, found, diags := getHostProfileAsState(
		context.Background(), 123, client, HostProfileModel{},
	)

	if !diags.HasError() {
		t.Fatal("expected a diagnostic on 500")
	}
	if found {
		t.Fatal("expected found=false on 500")
	}
}

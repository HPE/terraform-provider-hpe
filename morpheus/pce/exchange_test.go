// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package pce

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	jose "github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"

	"github.com/HPE/terraform-provider-hpe/morpheus/pce/sdk/token/iamversion"
	tokenutil "github.com/HPE/terraform-provider-hpe/morpheus/pce/sdk/token/token-util"
)

const cmpDetailsBody = `{
	"ServiceInstanceID": "18ba6409-ac59-4eac-9414-0147e72d615e",
	"URL": "https://morpheus.example.com",
	"TokenDetails": {
		"access_token": "morpheus-access-token",
		"expires_in": 3600
	}
}`

// testIAMToken returns a signed JWT suitable for use as a pre-generated IAM
// token. Supplying one skips the IAM leg of the exchange, but the token is
// still decoded to check its expiry, so it must be a well formed JWT.
func testIAMToken(t *testing.T) string {
	t.Helper()

	now := time.Now().Unix()

	claims := tokenutil.Token{
		Issuer:   "https://example.invalid/oauth2/default",
		Subject:  "clients/subject",
		Expiry:   now + 3600,
		IssuedAt: now,
		ClientID: "clientID",
		TenantID: "tenantID",
	}

	signer, err := jose.NewSigner(
		jose.SigningKey{Algorithm: jose.HS256, Key: []byte("secret")},
		nil,
	)
	if err != nil {
		t.Fatalf("could not create signer: %v", err)
	}

	token, err := jwt.Signed(signer).Claims(claims).CompactSerialize()
	if err != nil {
		t.Fatalf("could not sign token: %v", err)
	}

	return token
}

// newTestBroker returns a broker stub and a counter of how many exchanges it
// has served.
func newTestBroker(t *testing.T) (*httptest.Server, *atomic.Int64) {
	t.Helper()

	var calls atomic.Int64

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls.Add(1)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(cmpDetailsBody))
	}))

	t.Cleanup(srv.Close)

	return srv, &calls
}

// newRecordingBroker returns a broker stub that records the last request it
// served, so tests can assert how the exchange scoped it.
func newRecordingBroker(t *testing.T) (*httptest.Server, func() (url.Values, http.Header)) {
	t.Helper()

	var (
		mu     sync.Mutex
		query  url.Values
		header http.Header
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		query = r.URL.Query()
		header = r.Header.Clone()
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(cmpDetailsBody))
	}))

	t.Cleanup(srv.Close)

	return srv, func() (url.Values, http.Header) {
		mu.Lock()
		defer mu.Unlock()

		return query, header
	}
}

// Connected PCE scopes the broker request by GLCS space.
func TestTokenExchangeScopesConnectedBySpace(t *testing.T) {
	srv, recorded := newRecordingBroker(t)

	cfg := Config{
		IAMToken:  testIAMToken(t),
		BrokerURL: srv.URL,
		Version:   iamversion.GLCS,
		Location:  "BLR",
		Space:     "scopes-connected",
	}

	if _, _, err := TokenExchange(context.Background(), cfg); err != nil {
		t.Fatalf("TokenExchange() unexpected error: %v", err)
	}

	query, header := recorded()

	if got, want := query.Get("space"), "scopes-connected"; got != want {
		t.Errorf("space query parameter = %q, want %q", got, want)
	}

	if got, want := query.Get("location"), "BLR"; got != want {
		t.Errorf("location query parameter = %q, want %q", got, want)
	}

	// Workspace scoping is GLP only and must not leak into the GLCS exchange.
	if got := query.Get("tenantID"); got != "" {
		t.Errorf("tenantID query parameter = %q, want it to be absent", got)
	}

	if got := header.Get("X-Tenant-ID"); got != "" {
		t.Errorf("X-Tenant-ID header = %q, want it to be absent", got)
	}
}

// Disconnected PCE scopes the broker request by GLP workspace, sending it both
// as a query parameter and as a header. Location is sent for both deployment
// types: the broker rejects a request without one.
func TestTokenExchangeScopesDisconnectedByWorkspace(t *testing.T) {
	srv, recorded := newRecordingBroker(t)

	cfg := Config{
		IAMToken:    testIAMToken(t),
		BrokerURL:   srv.URL,
		Version:     iamversion.GLP,
		Location:    "site-a",
		WorkspaceID: "scopes-disconnected",
	}

	if _, _, err := TokenExchange(context.Background(), cfg); err != nil {
		t.Fatalf("TokenExchange() unexpected error: %v", err)
	}

	query, header := recorded()

	if got, want := query.Get("tenantID"), "scopes-disconnected"; got != want {
		t.Errorf("tenantID query parameter = %q, want %q", got, want)
	}

	if got, want := header.Get("X-Tenant-ID"), "scopes-disconnected"; got != want {
		t.Errorf("X-Tenant-ID header = %q, want %q", got, want)
	}

	// The broker requires a location, so it is always sent.
	if got, want := query.Get("location"), "site-a"; got != want {
		t.Errorf("location query parameter = %q, want %q", got, want)
	}

	// Space scoping is GLCS only and must not leak into the GLP exchange.
	if got := query.Get("space"); got != "" {
		t.Errorf("space query parameter = %q, want it to be absent", got)
	}
}

// Nothing is cached between calls, so a repeated exchange has to reach the
// broker again. This is the counterpart of the removed memoisation: the two
// muxed providers each perform their own exchange.
func TestTokenExchangeIsNotCached(t *testing.T) {
	srv, calls := newTestBroker(t)

	cfg := Config{
		IAMToken:  testIAMToken(t),
		BrokerURL: srv.URL,
		Version:   iamversion.GLCS,
		Location:  "BLR",
		Space:     "not-cached",
	}

	url, token, err := TokenExchange(context.Background(), cfg)
	if err != nil {
		t.Fatalf("first TokenExchange() unexpected error: %v", err)
	}

	if url != "https://morpheus.example.com" || token != "morpheus-access-token" {
		t.Fatalf("unexpected details: url=%q token=%q", url, token)
	}

	if _, _, err := TokenExchange(context.Background(), cfg); err != nil {
		t.Fatalf("second TokenExchange() unexpected error: %v", err)
	}

	if got, want := calls.Load(), int64(2); got != want {
		t.Errorf("broker exchanges = %d, want %d: the exchange must not be cached", got, want)
	}
}

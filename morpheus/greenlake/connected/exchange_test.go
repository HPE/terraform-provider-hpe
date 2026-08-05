// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package connected

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	jose "github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"

	tokenutil "github.com/HPE/terraform-provider-hpe/morpheus/greenlake/sdk/token/token-util"
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

// The framework and SDKv2 providers both exchange from the same configuration,
// so the second call must be served from the cache.
func TestTokenExchangeIsMemoised(t *testing.T) {
	srv, calls := newTestBroker(t)

	// IAMToken short circuits the IAM leg, so only the broker is exercised.
	// The values are unique to this test because the cache is package level.
	cfg := Config{
		IAMToken:  testIAMToken(t),
		BrokerURL: srv.URL,
		Location:  "BLR",
		Space:     "memoised",
	}

	url, token, err := TokenExchange(context.Background(), cfg)
	if err != nil {
		t.Fatalf("first TokenExchange() unexpected error: %v", err)
	}

	if url != "https://morpheus.example.com" || token != "morpheus-access-token" {
		t.Fatalf("unexpected details: url=%q token=%q", url, token)
	}

	cachedURL, cachedToken, err := TokenExchange(context.Background(), cfg)
	if err != nil {
		t.Fatalf("second TokenExchange() unexpected error: %v", err)
	}

	if cachedURL != url || cachedToken != token {
		t.Errorf("cached details differ: url=%q token=%q", cachedURL, cachedToken)
	}

	if got := calls.Load(); got != 1 {
		t.Errorf("broker was called %d times, want 1", got)
	}
}

// Configurations that target different places must not share a cached result.
func TestTokenExchangeIsolatesConfigurations(t *testing.T) {
	srv, calls := newTestBroker(t)

	first := Config{
		IAMToken:  testIAMToken(t),
		BrokerURL: srv.URL,
		Location:  "BLR",
		Space:     "isolated-one",
	}

	second := first
	second.Space = "isolated-two"

	if _, _, err := TokenExchange(context.Background(), first); err != nil {
		t.Fatalf("TokenExchange(first) unexpected error: %v", err)
	}

	if _, _, err := TokenExchange(context.Background(), second); err != nil {
		t.Fatalf("TokenExchange(second) unexpected error: %v", err)
	}

	if got := calls.Load(); got != 2 {
		t.Errorf("broker was called %d times, want 2", got)
	}
}

// A failed exchange must not be cached, so that a later attempt can retry.
func TestTokenExchangeDoesNotCacheFailures(t *testing.T) {
	var calls atomic.Int64

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls.Add(1)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error": "unauthorized"}`))
	}))
	t.Cleanup(srv.Close)

	cfg := Config{
		IAMToken:  testIAMToken(t),
		BrokerURL: srv.URL,
		Location:  "BLR",
		Space:     "failure",
	}

	for i := range 2 {
		if _, _, err := TokenExchange(context.Background(), cfg); err == nil {
			t.Fatalf("attempt %d: expected an error, got nil", i+1)
		}
	}

	if got := calls.Load(); got != 2 {
		t.Errorf("broker was called %d times, want 2", got)
	}
}

// (C) Copyright 2025 Hewlett Packard Enterprise Development LP
package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

// fakeMorpheus is the part of the Morpheus API that authentication touches:
// a password grant at /oauth/token and a bearer-protected API path. Like the
// real appliance (9.0.0 and 9.1.0.2), it refuses a token request that carries
// a bearer token -- any bearer, valid or not: 401 with an empty body and a
// WWW-Authenticate header, which is what made the provider's own token
// refresh impossible -- and it lets a test expire the current token at will.
type fakeMorpheus struct {
	t *testing.T

	mu      sync.Mutex
	issued  int             // tokens issued so far; the current one is tok-<issued>
	expired map[string]bool // bearer values that are no longer accepted

	tokenCalls          atomic.Int32 // password grants received
	tokenCallsWithAuth  atomic.Int32 // ...of which carried an Authorization header
	apiCalls            atomic.Int32 // bearer-protected requests that were accepted
	lastAcceptedAPIBody atomic.Pointer[string]

	rejectAllGrants atomic.Bool // simulate a changed password: every grant 401s
}

func newFakeMorpheus(t *testing.T) (*fakeMorpheus, *httptest.Server) {
	t.Helper()

	f := &fakeMorpheus{t: t, expired: map[string]bool{}}
	srv := httptest.NewServer(http.HandlerFunc(f.handle))
	t.Cleanup(srv.Close)

	return f, srv
}

func (f *fakeMorpheus) currentToken() string {
	f.mu.Lock()
	defer f.mu.Unlock()

	return fmt.Sprintf("tok-%d", f.issued)
}

// expireCurrent makes the token in use stop working, as the real one does
// after its lifetime.
func (f *fakeMorpheus) expireCurrent() {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.expired["Bearer "+fmt.Sprintf("tok-%d", f.issued)] = true
}

func (f *fakeMorpheus) handle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/oauth/token" {
		f.handleToken(w, r)

		return
	}

	auth := r.Header.Get("Authorization")

	f.mu.Lock()
	valid := auth == "Bearer "+fmt.Sprintf("tok-%d", f.issued) && !f.expired[auth]
	f.mu.Unlock()

	if !valid {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = io.WriteString(w, `{"error":"invalid_token","error_description":"Invalid access token"}`)

		return
	}

	body, _ := io.ReadAll(r.Body)
	s := string(body)
	f.lastAcceptedAPIBody.Store(&s)
	f.apiCalls.Add(1)

	w.Header().Set("Content-Type", "application/json")
	_, _ = io.WriteString(w, `{"ok":true}`)
}

func (f *fakeMorpheus) handleToken(w http.ResponseWriter, r *http.Request) {
	f.tokenCalls.Add(1)

	if r.Header.Get("Authorization") != "" {
		// What Morpheus does with a token request that carries a bearer
		// token, valid or not: its authorization-server chain tries to
		// decode the bearer as a JWT and refuses the request before the
		// grant is read -- 401, empty body, the reason in WWW-Authenticate.
		f.tokenCallsWithAuth.Add(1)
		w.Header().Set("WWW-Authenticate", `Bearer error="invalid_token", `+
			`error_description="An error occurred while attempting to decode the Jwt: Malformed token", `+
			`error_uri="https://tools.ietf.org/html/rfc6750#section-3.1"`)
		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	if f.rejectAllGrants.Load() {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = io.WriteString(w, `{"error":"invalid_grant","error_description":"Bad credentials"}`)

		return
	}

	if err := r.ParseForm(); err != nil || r.PostForm.Get("grant_type") != "password" ||
		r.PostForm.Get("username") != "admin" || r.PostForm.Get("password") != "secret" {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = io.WriteString(w, `{"error":"invalid_grant","error_description":"Bad credentials"}`)

		return
	}

	f.mu.Lock()
	f.issued++
	tok := fmt.Sprintf("tok-%d", f.issued)
	f.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"access_token": tok, "token_type": "bearer", "expires_in": 1439, "scope": "write",
	})
}

func newTestClient(t *testing.T, srv *httptest.Server) (*http.Client, *CredsRoundTripper) {
	t.Helper()

	rt := NewCredsRoundTripper(context.Background(), http.DefaultTransport, srv.URL, "admin", "secret")

	creds, ok := rt.(*CredsRoundTripper)
	if !ok {
		t.Fatalf("NewCredsRoundTripper returned %T, want *CredsRoundTripper", rt)
	}

	return &http.Client{Transport: rt}, creds
}

func mustGet(t *testing.T, client *http.Client, url string) *http.Response {
	t.Helper()

	resp, err := client.Get(url)
	if err != nil {
		t.Fatalf("GET %s: %v", url, err)
	}

	t.Cleanup(func() { _ = resp.Body.Close() })

	return resp
}

func TestCredsRoundTripper_FirstRequestLogsIn(t *testing.T) {
	f, srv := newFakeMorpheus(t)
	client, _ := newTestClient(t, srv)

	resp := mustGet(t, client, srv.URL+"/api/whoami")

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}

	if got := f.tokenCalls.Load(); got != 1 {
		t.Fatalf("token requests = %d, want 1", got)
	}

	if got := f.tokenCallsWithAuth.Load(); got != 0 {
		t.Fatalf("token requests carrying an Authorization header = %d, want 0", got)
	}
}

// The bug: after the token expired, the refresh sent the dead bearer token
// with the password grant, Morpheus refused it, and the operation failed with
// "could not authenticate". The refresh must log in without any bearer token
// and the refused request must then succeed with the new one.
func TestCredsRoundTripper_RefreshesExpiredTokenWithoutBearer(t *testing.T) {
	f, srv := newFakeMorpheus(t)
	client, creds := newTestClient(t, srv)

	mustGet(t, client, srv.URL+"/api/whoami")
	f.expireCurrent()

	resp := mustGet(t, client, srv.URL+"/api/processes/75")

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status after token expiry = %d, want 200 (refresh + retry)", resp.StatusCode)
	}

	if got := f.tokenCalls.Load(); got != 2 {
		t.Fatalf("token requests = %d, want 2 (initial login + one refresh)", got)
	}

	if got := f.tokenCallsWithAuth.Load(); got != 0 {
		t.Fatalf("token requests carrying an Authorization header = %d, want 0: the refresh must not send the expired token", got)
	}

	if got, want := creds.cachedAuthHeader(), "Bearer "+f.currentToken(); got != want {
		t.Fatalf("cached header = %q, want %q", got, want)
	}
}

// A PUT or POST refused at the token boundary is retried; the retry must
// carry the body again, not the already-consumed empty one.
func TestCredsRoundTripper_RetryResendsBody(t *testing.T) {
	f, srv := newFakeMorpheus(t)
	client, _ := newTestClient(t, srv)

	mustGet(t, client, srv.URL+"/api/whoami")
	f.expireCurrent()

	const body = `{"virtualSwitch":{"name":"vstest"}}`

	req, err := http.NewRequest(http.MethodPut, srv.URL+"/api/virtual-switches/14/edit", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("PUT: %v", err)
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}

	if got := f.lastAcceptedAPIBody.Load(); got == nil || *got != body {
		t.Fatalf("retried request body = %q, want %q", deref(got), body)
	}
}

// Many requests in flight when the token expires are all refused at once;
// exactly one of them must log in again, and the rest must reuse that token.
func TestCredsRoundTripper_ConcurrentRefreshLogsInOnce(t *testing.T) {
	f, srv := newFakeMorpheus(t)
	client, _ := newTestClient(t, srv)

	mustGet(t, client, srv.URL+"/api/whoami")
	f.expireCurrent()

	const n = 16

	var (
		wg       sync.WaitGroup
		failures atomic.Int32
	)

	for range n {
		wg.Add(1)

		go func() {
			defer wg.Done()

			resp, err := client.Get(srv.URL + "/api/processes/75")
			if err != nil {
				failures.Add(1)

				return
			}

			_ = resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				failures.Add(1)
			}
		}()
	}

	wg.Wait()

	if got := failures.Load(); got != 0 {
		t.Fatalf("%d of %d concurrent requests failed after the token expired", got, n)
	}

	if got := f.tokenCalls.Load(); got != 2 {
		t.Fatalf("token requests = %d, want 2: concurrent 401s must share one refresh", got)
	}

	if got := f.tokenCallsWithAuth.Load(); got != 0 {
		t.Fatalf("token requests carrying an Authorization header = %d, want 0", got)
	}
}

// When logging in again genuinely fails (the password changed), the error
// must say so -- and the dead token must not be kept for the next request to
// retry with.
func TestCredsRoundTripper_FailedRefreshClearsStaleToken(t *testing.T) {
	f, srv := newFakeMorpheus(t)
	client, creds := newTestClient(t, srv)

	mustGet(t, client, srv.URL+"/api/whoami")
	f.expireCurrent()
	f.rejectAllGrants.Store(true)

	_, err := client.Get(srv.URL + "/api/processes/75")
	if err == nil || !strings.Contains(err.Error(), "could not authenticate with the Morpheus API") {
		t.Fatalf("error = %v, want the authentication failure", err)
	}

	if got := creds.cachedAuthHeader(); got != "" {
		t.Fatalf("cached header after a failed refresh = %q, want it cleared", got)
	}

	if got := f.tokenCallsWithAuth.Load(); got != 0 {
		t.Fatalf("token requests carrying an Authorization header = %d, want 0", got)
	}
}

func deref(s *string) string {
	if s == nil {
		return "<nil>"
	}

	return *s
}

// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package containerip

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

func TestUnitReady_EmptyString(t *testing.T) {
	t.Parallel()

	if Ready("") {
		t.Error("expected Ready=false for empty string")
	}
}

func TestUnitReady_ZeroAddress(t *testing.T) {
	t.Parallel()

	if Ready("0.0.0.0") {
		t.Error("expected Ready=false for 0.0.0.0")
	}
}

func TestUnitReady_ValidIPv4(t *testing.T) {
	t.Parallel()

	if !Ready("10.0.1.5") {
		t.Error("expected Ready=true for 10.0.1.5")
	}
}

func TestUnitReady_ValidIPv6(t *testing.T) {
	t.Parallel()

	if !Ready("fe80::1") {
		t.Error("expected Ready=true for fe80::1")
	}
}

func TestUnitReady_WhitespaceTrimmed(t *testing.T) {
	t.Parallel()

	if Ready("  ") {
		t.Error("expected Ready=false for whitespace-only string")
	}
}

func TestUnitReady_WhitespaceAroundValid(t *testing.T) {
	t.Parallel()

	if !Ready("  10.0.1.5  ") {
		t.Error("expected Ready=true for whitespace-padded valid IP")
	}
}

func TestUnitReady_WhitespaceAroundZero(t *testing.T) {
	t.Parallel()

	if Ready("  0.0.0.0  ") {
		t.Error("expected Ready=false for whitespace-padded 0.0.0.0")
	}
}

// newTestClient creates an SDK client pointed at the given test server.
func newTestClient(server *httptest.Server) *sdk.APIClient {
	cfg := sdk.NewConfiguration()
	cfg.Servers = sdk.ServerConfigurations{
		{URL: server.URL},
	}

	return sdk.NewAPIClient(cfg)
}

// --- WaitAny tests ---

func TestUnitWaitAny_Success(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"instance":{"containerDetails":[{"ip":"10.0.0.1"}]}}`)
	}))
	defer ts.Close()

	client := newTestClient(ts)
	warned, err := WaitAny(context.Background(), client, 1, 5*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if warned {
		t.Error("expected warned=false on success")
	}
}

func TestUnitWaitAny_PermanentError(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{"message":"server error"}`)
	}))
	defer ts.Close()

	client := newTestClient(ts)
	warned, err := WaitAny(context.Background(), client, 1, 5*time.Second)
	if err == nil {
		t.Fatal("expected error on permanent API failure, got nil")
	}
	if warned {
		t.Error("expected warned=false on permanent error")
	}
}

func TestUnitWaitAny_Timeout(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Return an instance with no ready IP.
		fmt.Fprint(w, `{"instance":{"containerDetails":[{"ip":"0.0.0.0"}]}}`)
	}))
	defer ts.Close()

	client := newTestClient(ts)
	// Use a very short timeout so the test finishes fast.
	warned, err := WaitAny(context.Background(), client, 1, 1*time.Millisecond)
	if err != nil {
		t.Fatalf("expected nil error on timeout, got: %v", err)
	}
	if !warned {
		t.Error("expected warned=true on timeout")
	}
}

func TestUnitWaitAny_ContextCancelled(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"instance":{"containerDetails":[{"ip":"0.0.0.0"}]}}`)
	}))
	defer ts.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	client := newTestClient(ts)
	warned, err := WaitAny(ctx, client, 1, 30*time.Second)
	if err == nil {
		t.Fatal("expected error on cancelled context, got nil")
	}
	if warned {
		t.Error("expected warned=false on context cancellation")
	}
}

// --- Wait tests ---

func TestUnitWait_Success(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"instance":{"containerDetails":[{"id":42,"ip":"10.0.0.1"}]}}`)
	}))
	defer ts.Close()

	client := newTestClient(ts)
	ip, warned, err := Wait(context.Background(), client, 1, 42, 5*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if warned {
		t.Error("expected warned=false on success")
	}
	if ip != "10.0.0.1" {
		t.Errorf("expected ip=10.0.0.1, got %q", ip)
	}
}

func TestUnitWait_PermanentError(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, `{"message":"forbidden"}`)
	}))
	defer ts.Close()

	client := newTestClient(ts)
	ip, warned, err := Wait(context.Background(), client, 1, 42, 5*time.Second)
	if err == nil {
		t.Fatal("expected error on permanent API failure, got nil")
	}
	if warned {
		t.Error("expected warned=false on permanent error")
	}
	if ip != "" {
		t.Errorf("expected empty ip, got %q", ip)
	}
}

func TestUnitWait_ContainerNotFound(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Container 99 is not in the response.
		fmt.Fprint(w, `{"instance":{"containerDetails":[{"id":1,"ip":"10.0.0.1"}]}}`)
	}))
	defer ts.Close()

	client := newTestClient(ts)
	ip, warned, err := Wait(context.Background(), client, 1, 99, 5*time.Second)
	if err == nil {
		t.Fatal("expected error when container not found, got nil")
	}
	if warned {
		t.Error("expected warned=false on permanent error")
	}
	if ip != "" {
		t.Errorf("expected empty ip, got %q", ip)
	}
}

func TestUnitWait_Timeout(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Container exists but IP not ready.
		fmt.Fprint(w, `{"instance":{"containerDetails":[{"id":42,"ip":"0.0.0.0"}]}}`)
	}))
	defer ts.Close()

	client := newTestClient(ts)
	ip, warned, err := Wait(context.Background(), client, 1, 42, 1*time.Millisecond)
	if err != nil {
		t.Fatalf("expected nil error on timeout, got: %v", err)
	}
	if !warned {
		t.Error("expected warned=true on timeout")
	}
	if ip != "" {
		t.Errorf("expected empty ip on timeout, got %q", ip)
	}
}

func TestUnitWait_ContextCancelled(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"instance":{"containerDetails":[{"id":42,"ip":"0.0.0.0"}]}}`)
	}))
	defer ts.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	client := newTestClient(ts)
	ip, warned, err := Wait(ctx, client, 1, 42, 30*time.Second)
	if err == nil {
		t.Fatal("expected error on cancelled context, got nil")
	}
	if warned {
		t.Error("expected warned=false on context cancellation")
	}
	if ip != "" {
		t.Errorf("expected empty ip, got %q", ip)
	}
}

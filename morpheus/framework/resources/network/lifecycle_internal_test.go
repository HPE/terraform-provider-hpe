// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package network

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// fakeAppliance serves the two calls the lifecycle helpers make — DELETE and
// GET of one network — from scripted answers. Morpheus renders a refused
// delete as HTTP 200 with success:false, which is what the provider must read.
type fakeAppliance struct {
	t        *testing.T
	deletes  atomic.Int32
	gets     atomic.Int32
	onDelete func(n int32) (status int, body string)
	onGet    func(n int32) (status int, body string)
}

func (f *fakeAppliance) client() (*sdk.APIClient, func()) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/api/networks/") {
			f.t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusTeapot)

			return
		}

		var status int

		var body string

		switch r.Method {
		case http.MethodDelete:
			status, body = f.onDelete(f.deletes.Add(1))
		case http.MethodGet:
			status, body = f.onGet(f.gets.Add(1))
		default:
			f.t.Errorf("unexpected method %s", r.Method)
			status = http.StatusMethodNotAllowed
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))

	cfg := sdk.NewConfiguration()
	cfg.Servers[0].URL = srv.URL

	return sdk.NewAPIClient(cfg), srv.Close
}

func networkJSON(status string) string {
	b, _ := json.Marshal(map[string]any{"network": map[string]any{"id": 42, "name": "n", "status": status}})

	return string(b)
}

func TestUnitDeleteRefusedIsReportedNotSwallowed(t *testing.T) { //nolint:paralleltest // adjusts the package cadence
	f := &fakeAppliance{t: t}
	// Every delete: HTTP 200, success:false, no message — exactly what the
	// HVM plugin produces when its host-side delete throws.
	f.onDelete = func(int32) (int, string) { return http.StatusOK, `{"success":false}` }
	f.onGet = func(int32) (int, string) { return http.StatusOK, networkJSON("ok") }

	client, done := f.client()
	defer done()

	t.Cleanup(shortDeleteRetry())

	err := deleteNetworkAndWait(context.Background(), client, 42, time.Minute)
	if err == nil {
		t.Fatal("a refused delete must be an error")
	}

	if !strings.Contains(err.Error(), "Morpheus refused the delete") || !strings.Contains(err.Error(), "no reason given") {
		t.Errorf("error must carry the refusal: %v", err)
	}

	if got := f.deletes.Load(); got != deleteAttempts {
		t.Errorf("a refusal is retried %d times, got %d deletes", deleteAttempts, got)
	}

	if f.gets.Load() != 0 {
		t.Error("nothing to wait for when Morpheus never agreed to delete")
	}
}

func TestUnitDeleteRefusedThenAccepted(t *testing.T) { //nolint:paralleltest // adjusts the package cadence
	f := &fakeAppliance{t: t}
	f.onDelete = func(n int32) (int, string) {
		if n < 3 {
			return http.StatusOK, `{"success":false,"msg":"Failed to delete data network on one or more hosts"}`
		}

		return http.StatusOK, `{"success":true}`
	}
	// Still there on the first read, gone on the second.
	f.onGet = func(n int32) (int, string) {
		if n == 1 {
			return http.StatusOK, networkJSON("ok")
		}

		return http.StatusNotFound, `{"success":false,"msg":"Network not found"}`
	}

	client, done := f.client()
	defer done()

	t.Cleanup(shortDeleteRetry())
	t.Cleanup(shortSettleInterval())

	if err := deleteNetworkAndWait(context.Background(), client, 42, time.Minute); err != nil {
		t.Fatalf("the third delete was accepted and the network went away: %v", err)
	}

	if f.deletes.Load() != 3 || f.gets.Load() != 2 {
		t.Errorf("expected 3 deletes and 2 status reads, got %d and %d", f.deletes.Load(), f.gets.Load())
	}
}

func TestUnitDeleteOfMissingNetworkIsNotAnError(t *testing.T) {
	t.Parallel()

	f := &fakeAppliance{t: t}
	f.onDelete = func(int32) (int, string) { return http.StatusNotFound, `{"success":false,"msg":"Network not found"}` }
	f.onGet = func(int32) (int, string) { return http.StatusNotFound, "" }

	client, done := f.client()
	defer done()

	if err := deleteNetworkAndWait(context.Background(), client, 42, time.Minute); err != nil {
		t.Fatalf("a network that is already gone is deleted: %v", err)
	}
}

func TestUnitRefusalMessage(t *testing.T) {
	t.Parallel()

	cases := []struct {
		body map[string]any
		want string
	}{
		{nil, "no reason given"},
		{map[string]any{"msg": "in use"}, "in use"},
		{map[string]any{"errors": map[string]any{"vlanId": "taken", "name": "dup"}}, "name: dup; vlanId: taken"},
		{map[string]any{"msg": "bad", "errors": map[string]any{"x": 1}}, "bad; x: 1"},
	}

	for _, c := range cases {
		if got := refusalMessage(c.body); got != c.want {
			t.Errorf("refusalMessage(%v) = %q, want %q", c.body, got, c.want)
		}
	}
}

func TestUnitWaitForNetworkSettled(t *testing.T) { //nolint:paralleltest // adjusts the package cadence
	sequence := func(statuses ...string) func(n int32) (int, string) {
		return func(n int32) (int, string) {
			i := int(n) - 1
			if i >= len(statuses) {
				i = len(statuses) - 1
			}

			return http.StatusOK, networkJSON(statuses[i])
		}
	}

	t.Run("leaves provisioning", func(t *testing.T) {
		f := &fakeAppliance{t: t, onGet: sequence("provisioning", "provisioning", "available")}
		client, done := f.client()
		defer done()

		t.Cleanup(shortSettleInterval())

		if diags := waitForNetworkSettled(context.Background(), client, 42); diags.HasError() {
			t.Fatalf("settled network: %v", diags)
		}

		if f.gets.Load() != 3 {
			t.Errorf("expected to stop at the first non-provisioning status, got %d reads", f.gets.Load())
		}
	})

	t.Run("failed is an error", func(t *testing.T) {
		f := &fakeAppliance{t: t, onGet: sequence("provisioning", "failed")}
		client, done := f.client()
		defer done()

		t.Cleanup(shortSettleInterval())

		diags := waitForNetworkSettled(context.Background(), client, 42)
		if !diags.HasError() || !strings.Contains(fmt.Sprint(diags), "provisioning failed") {
			t.Errorf("a failed provisioning must be an error: %v", diags)
		}
	})
}

// The helpers below shorten the lifecycle cadence for the tests that would
// otherwise wait it out; they restore it afterwards.
func shortDeleteRetry() func() {
	old := deleteRetryAfter
	deleteRetryAfter = 5 * time.Millisecond

	return func() { deleteRetryAfter = old }
}

func shortSettleInterval() func() {
	old := settleInterval
	settleInterval = 5 * time.Millisecond

	return func() { settleInterval = old }
}

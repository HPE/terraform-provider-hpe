// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instance

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// fakeAppliance is the smallest Morpheus that instanceDeleter can be driven
// against: one instance, one server, and a script of how each answers.
type fakeAppliance struct {
	mu sync.Mutex
	// instanceStatuses is returned by successive GET /api/instances/{id}; the
	// last entry repeats. "" means 404.
	instanceStatuses []string
	instanceGets     int
	// deleteInstance is the status code and body for DELETE /api/instances/{id}.
	deleteInstanceCode int
	deleteInstanceBody string
	// serverStatuses is returned by successive GET /api/servers/{id}; "" means 404.
	serverStatuses []string
	serverGets     int
	calls          []string
}

func (f *fakeAppliance) handler(t *testing.T) http.Handler {
	t.Helper()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		defer f.mu.Unlock()
		f.calls = append(f.calls, r.Method+" "+r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/instances/1":
			status := f.instanceStatuses[min(f.instanceGets, len(f.instanceStatuses)-1)]
			f.instanceGets++
			if status == "" {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, `{"success":false,"msg":"Instance not found or insufficient access"}`)

				return
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"instance": map[string]any{
				"id": 1, "status": status, "statusMessage": "because",
				"containerDetails": []any{map[string]any{"server": map[string]any{"id": 7}}},
			}})
		case r.Method == http.MethodDelete && r.URL.Path == "/api/instances/1":
			w.WriteHeader(f.deleteInstanceCode)
			fmt.Fprint(w, f.deleteInstanceBody)
		case r.Method == http.MethodGet && r.URL.Path == "/api/servers/7":
			status := f.serverStatuses[min(f.serverGets, len(f.serverStatuses)-1)]
			f.serverGets++
			if status == "" {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, `{"success":false,"msg":"Server not found"}`)

				return
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"server": map[string]any{"id": 7, "status": status}})
		case r.Method == http.MethodPut && r.URL.Path == "/api/servers/7/stop":
			fmt.Fprint(w, `{"success":true}`)
		case r.Method == http.MethodDelete && r.URL.Path == "/api/servers/7":
			// The server path removes the instance and server records.
			f.instanceStatuses = []string{""}
			f.serverStatuses = []string{""}
			fmt.Fprint(w, `{"success":true}`)
		default:
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusTeapot)
		}
	})
}

func (f *fakeAppliance) called(prefix string) int {
	f.mu.Lock()
	defer f.mu.Unlock()
	n := 0
	for _, c := range f.calls {
		if strings.HasPrefix(c, prefix) {
			n++
		}
	}

	return n
}

func newDeleter(t *testing.T, f *fakeAppliance) *instanceDeleter {
	t.Helper()
	srv := httptest.NewServer(f.handler(t))
	t.Cleanup(srv.Close)
	cfg := sdk.NewConfiguration()
	cfg.Servers = sdk.ServerConfigurations{{URL: srv.URL}}

	return &instanceDeleter{client: sdk.NewAPIClient(cfg), id: 1, timeout: 30 * time.Second}
}

// TestUnitDeleteViaInstanceAPI covers the decision the destroy makes between
// the instance path and the server fallback (MORPH-18513).
func TestUnitDeleteViaInstanceAPI(t *testing.T) {
	t.Parallel()

	ok := `{"success":true}`

	tests := []struct {
		name         string
		appliance    *fakeAppliance
		wantRefused  bool
		wantErr      bool
		wantDeletes  int // DELETE /api/instances/1 calls
		wantFallback bool
	}{
		{
			name: "running instance: deleted through the instance API and gone",
			appliance: &fakeAppliance{
				instanceStatuses: []string{"running", "removing", ""}, deleteInstanceCode: 200, deleteInstanceBody: ok,
			},
			wantDeletes: 1,
		},
		{
			name: "busy instance settles, then is deleted",
			appliance: &fakeAppliance{
				instanceStatuses: []string{"resizing", "resizing", "running", "removing", ""}, deleteInstanceCode: 200, deleteInstanceBody: ok,
			},
			wantDeletes: 1,
		},
		{
			name: "Morpheus refuses (400): refused, so the caller falls back",
			appliance: &fakeAppliance{
				instanceStatuses: []string{"running"}, deleteInstanceCode: 400,
				deleteInstanceBody: `{"success":false,"msg":"Unable to remove instance in its current state"}`,
			},
			wantRefused: true, wantDeletes: 1,
		},
		{
			name: "Morpheus declines in a 200 body: refused",
			appliance: &fakeAppliance{
				instanceStatuses: []string{"running"}, deleteInstanceCode: 200,
				deleteInstanceBody: `{"success":false,"msg":"locked"}`,
			},
			wantRefused: true, wantDeletes: 1,
		},
		{
			name: "removal fails after acceptance (status warning): refused",
			appliance: &fakeAppliance{
				instanceStatuses: []string{"running", "removing", "warning"}, deleteInstanceCode: 200, deleteInstanceBody: ok,
			},
			wantRefused: true, wantDeletes: 1,
		},
		{
			name: "already gone: nothing to do",
			appliance: &fakeAppliance{
				instanceStatuses: []string{"running", ""}, deleteInstanceCode: 404,
				deleteInstanceBody: `{"success":false,"msg":"Instance not found or insufficient access"}`,
			},
			wantDeletes: 1,
		},
		{
			name: "server error: hard failure, no fallback",
			appliance: &fakeAppliance{
				instanceStatuses: []string{"running"}, deleteInstanceCode: 500, deleteInstanceBody: `{"success":false}`,
			},
			wantErr: true, wantDeletes: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := newDeleter(t, tt.appliance)
			ctx := context.Background()
			instance, gone, err := d.getInstance(ctx)
			if err != nil || gone {
				t.Fatalf("getInstance: gone=%v err=%v", gone, err)
			}

			err = d.deleteViaInstanceAPI(ctx, instance)
			var refused *instanceDeleteRefused
			gotRefused := err != nil && asRefused(err, &refused)
			switch {
			case tt.wantRefused && !gotRefused:
				t.Errorf("want instanceDeleteRefused, got %v", err)
			case tt.wantErr && (err == nil || gotRefused):
				t.Errorf("want hard error, got %v", err)
			case !tt.wantRefused && !tt.wantErr && err != nil:
				t.Errorf("want success, got %v", err)
			}
			if got := tt.appliance.called("DELETE /api/instances/1"); got != tt.wantDeletes {
				t.Errorf("DELETE /api/instances/1 called %d times, want %d", got, tt.wantDeletes)
			}
			if got := tt.appliance.called("DELETE /api/servers/7"); got != 0 {
				t.Errorf("deleteViaInstanceAPI must never touch the server API; DELETE /api/servers/7 called %d times", got)
			}
		})
	}
}

// TestUnitEnsureServerGone covers phase 2: a server that went with its
// instance, one still deprovisioning, and one left behind.
func TestUnitEnsureServerGone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		serverStatuses []string
		wantForce      bool
	}{
		{"gone with the instance", []string{""}, false},
		{"still deprovisioning, then gone", []string{"deprovisioning", "deprovisioning", ""}, false},
		{"left behind as provisioned: force-deleted", []string{"provisioned"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f := &fakeAppliance{instanceStatuses: []string{""}, serverStatuses: tt.serverStatuses}
			d := newDeleter(t, f)
			if err := d.ensureServerGone(context.Background(), 7); err != nil {
				t.Fatalf("ensureServerGone: %v", err)
			}
			got := f.called("DELETE /api/servers/7") > 0
			if got != tt.wantForce {
				t.Errorf("server force-deleted = %v, want %v (calls: %v)", got, tt.wantForce, f.calls)
			}
		})
	}
}

// asRefused reports whether err is, or wraps, an *instanceDeleteRefused.
func asRefused(err error, target **instanceDeleteRefused) bool {
	return errors.As(err, target)
}

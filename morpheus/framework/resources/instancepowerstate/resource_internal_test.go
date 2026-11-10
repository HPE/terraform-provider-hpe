// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instancepowerstate

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// testInterval is the poll interval the unit tests pass to the polling
// functions. The fake appliance answers instantly, so this is the only wait.
// It is a parameter, not a package-level override: the acceptance test in
// resource_test.go is built into the same test binary and runs the provider
// in-process, and it must keep polling a real appliance at pollInterval.
const testInterval = 20 * time.Millisecond

// fakeAppliance is the smallest Morpheus the power-state resource can be
// driven against: one instance with one server, and a script of how each
// answers. Every request is recorded with its query string so a test can
// assert on exactly what went over the wire.
type fakeAppliance struct {
	mu sync.Mutex
	// instanceStatuses is returned by successive GET /api/instances/1; the
	// last entry repeats.
	instanceStatuses []string
	instanceGets     int
	// servers is the instance's server list; defaults to just server 7.
	// noServers makes the instance report an empty list instead.
	servers   []int64
	noServers bool
	// powerStates is returned by successive GET /api/servers/{id} (any id);
	// the last entry repeats. "" omits powerState from the response.
	powerStates []string
	// serverPowerStates overrides powerStates for one server id.
	serverPowerStates map[int64][]string
	// serverDelay is added to every GET /api/servers/{id} response.
	serverDelay time.Duration
	serverGets  map[string]int // per server path
	calls       []string
}

func (f *fakeAppliance) handler(t *testing.T) http.Handler {
	t.Helper()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if f.serverDelay > 0 && strings.HasPrefix(r.URL.Path, "/api/servers/") {
			time.Sleep(f.serverDelay)
		}

		f.mu.Lock()
		defer f.mu.Unlock()
		f.calls = append(f.calls, r.Method+" "+r.URL.RequestURI())
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/instances/1":
			status := f.instanceStatuses[min(f.instanceGets, len(f.instanceStatuses)-1)]
			f.instanceGets++
			_ = json.NewEncoder(w).Encode(map[string]any{"instance": map[string]any{
				"id": 1, "status": status, "servers": f.serverIDs(),
			}})
		case r.Method == http.MethodPut && strings.HasPrefix(r.URL.Path, "/api/instances/1/"):
			fmt.Fprint(w, `{"success":true,"results":{"1":{"success":true}}}`)
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/api/servers/"):
			server := map[string]any{"id": 7, "status": "provisioned"}
			if f.serverGets == nil {
				f.serverGets = map[string]int{}
			}
			states := f.powerStates
			if id, err := strconv.ParseInt(strings.TrimPrefix(r.URL.Path, "/api/servers/"), 10, 64); err == nil {
				if override, ok := f.serverPowerStates[id]; ok {
					states = override
				}
			}
			if len(states) > 0 {
				ps := states[min(f.serverGets[r.URL.Path], len(states)-1)]
				if ps != "" {
					server["powerState"] = ps
				}
			}
			f.serverGets[r.URL.Path]++
			_ = json.NewEncoder(w).Encode(map[string]any{"server": server})
		default:
			t.Errorf("unexpected request %s %s", r.Method, r.URL.RequestURI())
			w.WriteHeader(http.StatusTeapot)
		}
	})
}

func (f *fakeAppliance) serverIDs() []int64 {
	if f.noServers {
		return []int64{}
	}

	if len(f.servers) > 0 {
		return f.servers
	}

	return []int64{7}
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

func newFakeClient(t *testing.T, f *fakeAppliance) *sdk.APIClient {
	t.Helper()
	srv := httptest.NewServer(f.handler(t))
	t.Cleanup(srv.Close)
	cfg := sdk.NewConfiguration()
	cfg.Servers = sdk.ServerConfigurations{{URL: srv.URL}}

	return sdk.NewAPIClient(cfg)
}

// TestUnitFireActionSendsServerFlag pins the query string of each power
// action. Start and stop must carry server=true: without it Morpheus only
// runs the instance's "Stop Service"/"Start Service" path and, on a layout
// with start/stop-phase scripts, leaves the VM powered on while reporting the
// instance stopped. The generated client never sends the parameter's spec
// default, so the flag has to be explicit. Suspend ignores it server-side.
func TestUnitFireActionSendsServerFlag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desired string
		want    string
	}{
		{desired: "stopped", want: "PUT /api/instances/1/stop?server=true"},
		{desired: "running", want: "PUT /api/instances/1/start?server=true"},
		{desired: "suspended", want: "PUT /api/instances/1/suspend"},
	}

	for _, tt := range tests {
		t.Run(tt.desired, func(t *testing.T) {
			t.Parallel()

			f := &fakeAppliance{}
			client := newFakeClient(t, f)

			if err := fireAction(t.Context(), client, 1, tt.desired); err != nil {
				t.Fatalf("fireAction(%s): %v", tt.desired, err)
			}

			if got := f.called(tt.want); got != 1 {
				t.Errorf("want exactly one %q, got calls %v", tt.want, f.calls)
			}

			// Any other query parameter on the action (e.g. muteMonitoring)
			// would be a change in behaviour worth noticing.
			for _, c := range f.calls {
				if c != tt.want {
					t.Errorf("unexpected request %q (want only %q)", c, tt.want)
				}
			}
		})
	}
}

// TestUnitFireActionRejectsUnknownState guards the validator's contract.
func TestUnitFireActionRejectsUnknownState(t *testing.T) {
	t.Parallel()

	client := newFakeClient(t, &fakeAppliance{})
	if err := fireAction(t.Context(), client, 1, "rebooted"); err == nil {
		t.Fatal("expected an error for an unknown desired_state")
	}
}

// TestUnitPollStatusConfirmsServerPowerState covers the second half of a
// transition: the instance status alone is not trusted, the server's
// powerState has to agree before the resource reports success.
func TestUnitPollStatusConfirmsServerPowerState(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		desired     string
		statuses    []string
		powerStates []string
		noServers   bool
		wantStatus  string
		wantErr     bool
		wantHostGet bool
	}{
		{
			name:        "stopped: waits for the instance, then for powerState off",
			desired:     "stopped",
			statuses:    []string{"stopping", "stopped"},
			powerStates: []string{"on", "off"},
			wantStatus:  "stopped",
			wantHostGet: true,
		},
		{
			name:        "running: powerState on confirms the start",
			desired:     "running",
			statuses:    []string{"running"},
			powerStates: []string{"on"},
			wantStatus:  "running",
			wantHostGet: true,
		},
		{
			name:        "stopped: a server that never powers off is a failure",
			desired:     "stopped",
			statuses:    []string{"stopped"},
			powerStates: []string{"on"},
			wantStatus:  "stopped",
			wantErr:     true,
			wantHostGet: true,
		},
		{
			name:        "stopped: a missing powerState is retried, not trusted",
			desired:     "stopped",
			statuses:    []string{"stopped"},
			powerStates: []string{"", "off"},
			wantStatus:  "stopped",
			wantHostGet: true,
		},
		{
			name:       "stopped: an instance with no servers cannot be confirmed",
			desired:    "stopped",
			statuses:   []string{"stopped"},
			noServers:  true,
			wantStatus: "stopped",
			wantErr:    true,
		},
		{
			name:       "suspended: no servers is fine, the instance status is the signal",
			desired:    "suspended",
			statuses:   []string{"suspended"},
			noServers:  true,
			wantStatus: "suspended",
		},
		{
			name:        "suspended: instance status is the only signal",
			desired:     "suspended",
			statuses:    []string{"suspended"},
			wantStatus:  "suspended",
			wantHostGet: false,
		},
		{
			name:       "an error status is permanent",
			desired:    "running",
			statuses:   []string{"failed"},
			wantStatus: "failed",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f := &fakeAppliance{instanceStatuses: tt.statuses, powerStates: tt.powerStates, noServers: tt.noServers}
			client := newFakeClient(t, f)

			// The fake answers instantly, so the poll interval is the only
			// wait. Each phase needs at most one retry; the generous budget
			// is slack for a loaded machine running the cases in parallel.
			got, err := pollStatus(t.Context(), client, 1, tt.desired, 25*testInterval, testInterval)
			if (err != nil) != tt.wantErr {
				t.Fatalf("pollStatus error = %v, wantErr %v (calls %v)", err, tt.wantErr, f.calls)
			}

			if got != tt.wantStatus {
				t.Errorf("final status = %q, want %q", got, tt.wantStatus)
			}

			if hostGets := f.called("GET /api/servers/7"); (hostGets > 0) != tt.wantHostGet {
				t.Errorf("GET /api/servers/7 called %d times, wantHostGet %v (calls %v)", hostGets, tt.wantHostGet, f.calls)
			}
		})
	}
}

// TestUnitPollStatusHonoursDeadlineAcrossServers pins that the configured
// timeout bounds the whole transition. The first servers power off only on
// their second poll and between them use most of the budget; the many that
// follow are already off, but each costs a round trip and there are far more
// of them than the remaining budget can visit. Once the deadline has passed
// the resource must stop and report the timeout rather than carry on through
// the remaining servers -- each of which would otherwise be granted requests,
// and up to a poll interval, of its own, with the call succeeding well after
// the deadline. The assertions are behavioural (which error, how many servers
// were visited), not wall-clock, so a loaded machine cannot flake them.
func TestUnitPollStatusHonoursDeadlineAcrossServers(t *testing.T) {
	t.Parallel()

	const (
		servers = 200 // far more than the budget left after the slow ones
		slow    = 3   // servers that each cost a poll interval
	)

	f := &fakeAppliance{
		instanceStatuses:  []string{"stopped"},
		powerStates:       []string{"off"},
		serverPowerStates: map[int64][]string{},
		serverDelay:       testInterval / 4, // 200 of these cannot fit in the budget
	}
	for i := range servers {
		id := int64(100 + i)
		f.servers = append(f.servers, id)
		if i < slow {
			f.serverPowerStates[id] = []string{"on", "off"}
		}
	}

	client := newFakeClient(t, f)
	// Room for the slow servers plus slack for a loaded machine; still only
	// a few intervals, which the remaining round trips cannot fit into.
	timeout := (slow + 5) * testInterval

	_, err := pollStatus(t.Context(), client, 1, "stopped", timeout, testInterval)
	if err == nil {
		t.Fatal("expected a timeout error; every server was visited after the deadline")
	}

	if !strings.Contains(err.Error(), "timed out before confirming") {
		t.Errorf("error = %q, want the deadline error", err)
	}

	// The servers past the deadline are never queried.
	if gets := f.called("GET /api/servers/"); gets >= servers {
		t.Errorf("made %d server requests for %d servers; expected the deadline to stop the loop", gets, servers)
	}
}

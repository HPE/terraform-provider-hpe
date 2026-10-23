// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/opsramp/client"
)

func TestUnitClientMonitoringAddons(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name   string
		addons []string
		log    bool
		trace  bool
	}{
		{name: "empty"},
		{name: "ordinary", addons: []string{"Application Management", "Custom Addon"}},
		{name: "log", addons: []string{"Log Management"}, log: true},
		{name: "trace", addons: []string{"Trace Management"}, trace: true},
		{
			name: "mixed", addons: []string{"Adapter Integrations", "Trace Management", "Log Management", "Stream Exports"},
			log: true, trace: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			regular, logEnabled, traceEnabled := clientAddonsForRequest(tc.addons)
			if logEnabled != tc.log || traceEnabled != tc.trace {
				t.Fatalf("monitoring flags = %t, %t; want %t, %t", logEnabled, traceEnabled, tc.log, tc.trace)
			}
			for _, addon := range regular {
				if addon == "Log Management" || addon == "Trace Management" {
					t.Fatalf("special addon leaked into API addOns: %q", addon)
				}
			}

			payload, err := json.Marshal(client.CreateClient{
				Addons: regular, LogMonitoringEnabled: logEnabled, TraceMonitoringEnabled: traceEnabled,
			})
			if err != nil {
				t.Fatal(err)
			}
			var fields map[string]json.RawMessage
			if err := json.Unmarshal(payload, &fields); err != nil {
				t.Fatal(err)
			}
			for key, want := range map[string]bool{"logMonitoringEnabled": tc.log, "traceMonitoringEnabled": tc.trace} {
				wantJSON := "false"
				if want {
					wantJSON = "true"
				}
				if string(fields[key]) != wantJSON {
					t.Fatalf("%s = %s; want explicit %s", key, fields[key], wantJSON)
				}
			}
			if len(regular) == 0 && string(fields["addOns"]) != "[]" {
				t.Fatalf("empty addOns must be sent to clear existing addons: %s", fields["addOns"])
			}

			response := &client.ClientResponse{
				Addons: regular, LogMonitoringEnabled: logEnabled, TraceMonitoringEnabled: traceEnabled,
			}
			gotSet, err := stringSetFromSlice(clientAddonsFromResponse(response))
			if err != nil {
				t.Fatal(err)
			}
			wantSet, err := stringSetFromSlice(tc.addons)
			if err != nil {
				t.Fatal(err)
			}
			if !gotSet.Equal(wantSet) {
				t.Fatalf("round-trip addons = %s; want %s", gotSet, wantSet)
			}
		})
	}
}

func TestUnitClientMonitoringSchema(t *testing.T) {
	t.Parallel()

	var resp resource.SchemaResponse
	(&ClientResource{}).Schema(context.Background(), resource.SchemaRequest{}, &resp)
	for _, name := range []string{"log_monitoring_enabled", "trace_monitoring_enabled"} {
		attribute, ok := resp.Schema.Attributes[name].(schema.BoolAttribute)
		if !ok || !attribute.Computed || attribute.Optional || attribute.Required {
			t.Fatalf("%s must remain a read-only boolean", name)
		}
	}
}

func TestUnitClientMonitoringLifecycle(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	var stored client.ClientResponse
	var requests []client.CreateClient
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch req.Method {
		case http.MethodPost:
			var payload client.CreateClient
			if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)

				return
			}
			requests = append(requests, payload)
			stored = client.ClientResponse{
				Id: 42, UniqueId: "client-test", Name: payload.Name,
				Addons: payload.Addons, Packages: payload.Packages,
				LogMonitoringEnabled: payload.LogMonitoringEnabled, TraceMonitoringEnabled: payload.TraceMonitoringEnabled,
			}
			if err := json.NewEncoder(w).Encode(client.ClientMinimal{Id: 42, UniqueId: stored.UniqueId, Name: stored.Name}); err != nil {
				t.Errorf("encode write response: %v", err)
			}
		case http.MethodGet:
			if err := json.NewEncoder(w).Encode(stored); err != nil {
				t.Errorf("encode read response: %v", err)
			}
		default:
			http.Error(w, "unexpected method", http.StatusMethodNotAllowed)
		}
	}))
	defer server.Close()

	r := &ClientResource{}
	r.apiClient = &client.OpsRampClient{BaseUrl: server.URL, TenantId: "msp-test", Client: server.Client(), Scope: "MSP"}
	var schemaResp resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)
	current := tfsdk.State{Schema: schemaResp.Schema}

	for step, tc := range []struct {
		name    string
		addons  []string
		regular []string
		log     bool
		trace   bool
	}{
		{
			"create both",
			[]string{"Application Management", "Log Management", "Trace Management"},
			[]string{"Application Management"},
			true, true,
		},
		{"disable trace", []string{"Stream Exports", "Log Management"}, []string{"Stream Exports"}, true, false},
		{"enable trace only", []string{"Trace Management"}, []string{}, false, true},
		{"disable all", []string{}, []string{}, false, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			addons, err := stringSetFromSlice(tc.addons)
			if err != nil {
				t.Fatal(err)
			}
			packageNames := []string{"Hybrid Discovery and Monitoring"}
			if step == 0 {
				packageNames = append(packageNames, "Event and Incident Management")
			}
			packages, err := stringSetFromSlice(packageNames)
			if err != nil {
				t.Fatal(err)
			}
			model := ClientModel{Name: types.StringValue("Test Client"), Addons: addons, Packages: packages}
			plan := tfsdk.Plan{Schema: schemaResp.Schema}
			if diags := plan.Set(ctx, &model); diags.HasError() {
				t.Fatal(diags)
			}

			if step == 0 {
				resp := resource.CreateResponse{State: tfsdk.State{Schema: schemaResp.Schema}}
				r.Create(ctx, resource.CreateRequest{Plan: plan}, &resp)
				if resp.Diagnostics.HasError() {
					t.Fatal(resp.Diagnostics)
				}
				current = resp.State
			} else {
				resp := resource.UpdateResponse{State: current}
				r.Update(ctx, resource.UpdateRequest{Plan: plan, State: current}, &resp)
				if resp.Diagnostics.HasError() {
					t.Fatal(resp.Diagnostics)
				}
				current = resp.State
			}

			if len(requests) != step+1 {
				t.Fatalf("write requests = %d; want %d", len(requests), step+1)
			}
			payload := requests[step]
			if !reflect.DeepEqual(payload.Addons, tc.regular) ||
				payload.LogMonitoringEnabled != tc.log || payload.TraceMonitoringEnabled != tc.trace {
				t.Fatalf("unexpected request: %+v", payload)
			}
			assertState := func(state tfsdk.State) {
				t.Helper()
				var got ClientModel
				if diags := state.Get(ctx, &got); diags.HasError() {
					t.Fatal(diags)
				}
				if !got.Addons.Equal(addons) || !got.Packages.Equal(packages) {
					t.Fatalf("state addons/packages = %s / %s; want %s / %s", got.Addons, got.Packages, addons, packages)
				}
				if !got.LogMonitoringEnabled.Equal(types.BoolValue(tc.log)) ||
					!got.TraceMonitoringEnabled.Equal(types.BoolValue(tc.trace)) {
					t.Fatalf("state flags = %s / %s; want %t / %t", got.LogMonitoringEnabled, got.TraceMonitoringEnabled, tc.log, tc.trace)
				}
			}
			assertState(current)

			readResp := resource.ReadResponse{State: current}
			r.Read(ctx, resource.ReadRequest{State: current}, &readResp)
			if readResp.Diagnostics.HasError() {
				t.Fatal(readResp.Diagnostics)
			}
			assertState(readResp.State)
			current = readResp.State

			importResp := resource.ImportStateResponse{State: tfsdk.State{Schema: schemaResp.Schema}}
			r.ImportState(ctx, resource.ImportStateRequest{ID: "client-test"}, &importResp)
			if importResp.Diagnostics.HasError() {
				t.Fatal(importResp.Diagnostics)
			}
			assertState(importResp.State)
		})
	}
}

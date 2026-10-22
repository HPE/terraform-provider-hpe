// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/HPE/terraform-provider-hpe/opsramp/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func servicemapCriteria(t *testing.T, metrics []string, operator string, components types.Set) types.Set {
	t.Helper()
	var criteria []ServicemapMetricCriterionModel
	for _, metric := range metrics {
		criteria = append(criteria, ServicemapMetricCriterionModel{
			MetricName: types.StringValue(metric), Operator: types.StringValue(operator), Components: components,
		})
	}
	value, diags := types.SetValueFrom(context.Background(), servicemapMetricCriterionType, criteria)
	if diags.HasError() {
		t.Fatal(diags)
	}
	return value
}

func assertServicemapMonitorJSON(t *testing.T, payload []byte, alertType string, metrics []string) {
	t.Helper()

	var body struct {
		Monitor map[string]json.RawMessage `json:"serviceAvailabilityMonitor"`
	}
	if err := json.Unmarshal(payload, &body); err != nil {
		t.Fatal(err)
	}
	var gotMetrics []string
	if err := json.Unmarshal(body.Monitor["metrics"], &gotMetrics); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(gotMetrics, metrics) {
		t.Fatalf("metrics = %#v; want %#v", gotMetrics, metrics)
	}
	if string(body.Monitor["alertType"]) != `"`+alertType+`"` || string(body.Monitor["matchType"]) != `"ANY"` {
		t.Fatalf("unexpected monitor: %s", payload)
	}
	criteriaJSON, exists := body.Monitor["metricCriteria"]
	if alertType != "critical-alert" {
		if exists {
			t.Fatalf("metricCriteria must be omitted for %s: %s", alertType, criteriaJSON)
		}

		return
	}
	if !exists || string(criteriaJSON) == "null" {
		t.Fatalf("critical-alert requires metricCriteria: %s", payload)
	}
	var criteria []client.ServicemapMetricCriterion
	if err := json.Unmarshal(criteriaJSON, &criteria); err != nil {
		t.Fatal(err)
	}
	if len(criteria) != len(metrics) {
		t.Fatalf("criteria count = %d; want %d", len(criteria), len(metrics))
	}
	for _, criterion := range criteria {
		found := false
		for _, metric := range metrics {
			found = found || criterion.MetricName == metric
		}
		if !found || criterion.Operator != "-" {
			t.Fatalf("unexpected criterion: %#v", criterion)
		}
	}
}

func TestUnitServicemapMonitorPayload(t *testing.T) {
	t.Parallel()

	for _, alertType := range []string{"critical-alert", "availability", "any-critical-alert"} {
		for _, metrics := range [][]string{{}, {"aggregator_discovery_aggregation_count_total"}, {"metric-a", "metric-b"}} {
			t.Run(fmt.Sprintf("%s/%d-metrics", alertType, len(metrics)), func(t *testing.T) {
				plan := ServicemapModel{AlertType: types.StringValue(alertType)}
				if alertType == "critical-alert" {
					plan.MetricCriteria = servicemapCriteria(t, metrics, "-", types.SetNull(types.StringType))
				}
				for _, metric := range metrics {
					plan.Metrics = append(plan.Metrics, types.StringValue(metric))
				}
				model, diags := translatePlanToModel(context.Background(), plan)
				if diags.HasError() {
					t.Fatal(diags)
				}
				payload, err := json.Marshal(model)
				if err != nil {
					t.Fatal(err)
				}
				wantMetrics := metrics
				if alertType == "any-critical-alert" {
					wantMetrics = []string{}
				}
				if alertType != "critical-alert" || len(metrics) > 0 {
					assertServicemapMonitorJSON(t, payload, alertType, wantMetrics)
				}
			})
		}
	}
	for _, alertType := range []types.String{types.StringNull(), types.StringValue("")} {
		model, diags := translatePlanToModel(context.Background(), ServicemapModel{AlertType: alertType})
		if diags.HasError() {
			t.Fatal(diags)
		}
		if monitor := model.ServiceAvailabilityMonitor; monitor != nil {
			t.Fatalf("unset alert type must not generate a monitor: %#v", monitor)
		}
	}
}

func TestUnitServicemapMonitorLifecycle(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	var requests [][]byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost || req.URL.Path != "/api/v3/tenants/client-test/serviceGroup" {
			http.Error(w, "unexpected request", http.StatusBadRequest)

			return
		}
		var payload json.RawMessage
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}
		requests = append(requests, payload)
		var response client.CreateServicemap
		if err := json.Unmarshal(payload, &response); err != nil {
			t.Errorf("decode service map: %v", err)

			return
		}
		response.Id = "service-test"
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("encode service map: %v", err)
		}
	}))
	defer server.Close()

	r := &ServicemapResource{}
	r.apiClient = &client.OpsRampClient{BaseUrl: server.URL, TenantId: "client-test", Client: server.Client(), Scope: "CLIENT"}
	var schemaResp resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)
	current := tfsdk.State{Schema: schemaResp.Schema}

	for step, tc := range []struct {
		name      string
		alertType string
		metrics   []string
	}{
		{"create critical", "critical-alert", []string{"metric-a"}},
		{"change metrics", "critical-alert", []string{"metric-b", "metric-c"}},
		{"switch to availability", "availability", []string{}},
		{"switch back to critical", "critical-alert", []string{"metric-d"}},
		{"switch to any critical", "any-critical-alert", []string{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			model := ServicemapModel{
				Name:           types.StringValue("test-service"),
				Type:           types.StringValue("Service"),
				AlertType:      types.StringValue(tc.alertType),
				ThresholdType:  types.StringValue("count"),
				ThresholdLimit: types.Int64Value(1),
				MetricCriteria: types.SetNull(servicemapMetricCriterionType),
			}
			if tc.alertType == "critical-alert" {
				model.MetricCriteria = servicemapCriteria(t, tc.metrics, "-", types.SetNull(types.StringType))
			}
			if step > 0 {
				model.Id = types.StringValue("service-test")
			}
			for _, metric := range tc.metrics {
				model.Metrics = append(model.Metrics, types.StringValue(metric))
			}
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
				t.Fatalf("request count = %d; want %d", len(requests), step+1)
			}
			assertServicemapMonitorJSON(t, requests[step], tc.alertType, tc.metrics)
			var state ServicemapModel
			if diags := current.Get(ctx, &state); diags.HasError() {
				t.Fatal(diags)
			}
			if !reflect.DeepEqual(state.Metrics, model.Metrics) || !state.AlertType.Equal(model.AlertType) || !state.MetricCriteria.Equal(model.MetricCriteria) {
				t.Fatalf("monitor state does not match plan: %#v", state)
			}

		})
	}
}

func TestUnitServicemapMonitorValidation(t *testing.T) {
				t.Parallel()
				ctx := context.Background()
				nullComponents := types.SetNull(types.StringType)
				components := types.SetValueMust(types.StringType, []attr.Value{types.StringValue("ads")})
				metrics := types.SetValueMust(types.StringType, []attr.Value{types.StringValue("metric-a")})
				nullCriteria := types.SetNull(servicemapMetricCriterionType)
				for _, tc := range []struct {
					name string
					alertType types.String
					metrics types.Set
					criteria types.Set
					wantError bool
				}{
					{"dash", types.StringValue("critical-alert"), metrics, servicemapCriteria(t, []string{"metric-a"}, "-", nullComponents), false},
					{"equals", types.StringValue("critical-alert"), metrics, servicemapCriteria(t, []string{"metric-a"}, "=", components), false},
					{"missing criteria", types.StringValue("critical-alert"), metrics, nullCriteria, true},
					{"empty criteria", types.StringValue("critical-alert"), metrics, servicemapCriteria(t, nil, "-", nullComponents), true},
					{"missing metrics", types.StringValue("critical-alert"), types.SetNull(types.StringType), servicemapCriteria(t, []string{"metric-a"}, "-", nullComponents), true},
					{"mismatched names", types.StringValue("critical-alert"), metrics, servicemapCriteria(t, []string{"metric-b"}, "-", nullComponents), true},
					{"equals missing components", types.StringValue("critical-alert"), metrics, servicemapCriteria(t, []string{"metric-a"}, "=", nullComponents), true},
					{"equals empty components", types.StringValue("critical-alert"), metrics, servicemapCriteria(t, []string{"metric-a"}, "=", types.SetValueMust(types.StringType, nil)), true},
					{"dash components", types.StringValue("critical-alert"), metrics, servicemapCriteria(t, []string{"metric-a"}, "-", components), true},
					{"bad operator", types.StringValue("critical-alert"), metrics, servicemapCriteria(t, []string{"metric-a"}, ">", nullComponents), true},
					{"any critical nonempty metrics", types.StringValue("any-critical-alert"), metrics, nullCriteria, true},
					{"any critical omitted metrics", types.StringValue("any-critical-alert"), types.SetNull(types.StringType), nullCriteria, false},
					{"any critical empty metrics", types.StringValue("any-critical-alert"), types.SetValueMust(types.StringType, nil), nullCriteria, false},
					{"availability criteria", types.StringValue("availability"), metrics, servicemapCriteria(t, []string{"metric-a"}, "-", nullComponents), true},
					{"unknown alert", types.StringUnknown(), metrics, types.SetUnknown(servicemapMetricCriterionType), false},
					{"unknown criteria", types.StringValue("critical-alert"), metrics, types.SetUnknown(servicemapMetricCriterionType), false},
					{"unknown names", types.StringValue("critical-alert"), types.SetUnknown(types.StringType), servicemapCriteria(t, []string{"metric-a"}, "-", nullComponents), false},
					{"unknown components", types.StringValue("critical-alert"), metrics, servicemapCriteria(t, []string{"metric-a"}, "=", types.SetUnknown(types.StringType)), false},
				} {
					t.Run(tc.name, func(t *testing.T) {
						diags := validateServicemapMonitor(ctx, tc.alertType, tc.metrics, tc.criteria)
						if diags.HasError() != tc.wantError {
							t.Fatalf("validation errors = %v; want error %t", diags, tc.wantError)
						}

								})
				}
				model := ServicemapModel{
					AlertType: types.StringValue("critical-alert"),
					Metrics: []types.String{types.StringValue("metric-a")},
					MetricCriteria: servicemapCriteria(t, []string{"metric-a"}, "=", components),
				}
				request, diags := translatePlanToModel(ctx, model)
				if diags.HasError() {
					t.Fatal(diags)
				}
				payload, err := json.Marshal(request)
				if err != nil {
					t.Fatal(err)
				}
				var decoded client.CreateServicemap
				if err := json.Unmarshal(payload, &decoded); err != nil {
					t.Fatal(err)
				}
				criterion := decoded.ServiceAvailabilityMonitor.MetricCriteria[0]
				if criterion.Operator != "=" || !reflect.DeepEqual(criterion.Components, []string{"ads"}) {
					t.Fatalf("equals components not serialized: %s", payload)
				}
			}

func TestUnitServicemapImportState(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet || req.URL.Path != "/api/v3/tenants/client-test/service-maps/service-test" {
			http.Error(w, "unexpected request", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"id":"service-test","name":"test-service","nodeType":"SERVICE","availabilityThreshold":{"thresholdType":"count","thresholdLimit":1}}`)
	}))
	defer server.Close()

	r := &ServicemapResource{}
	r.apiClient = &client.OpsRampClient{BaseUrl: server.URL, TenantId: "client-test", Client: server.Client(), Scope: "CLIENT"}
	ctx := context.Background()
	var schemaResp resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)
	resp := resource.ImportStateResponse{State: tfsdk.State{Schema: schemaResp.Schema}}
	r.ImportState(ctx, resource.ImportStateRequest{ID: "service-test"}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatal(resp.Diagnostics)
	}
	if resp.State.Raw.IsNull() {
		t.Fatal("import returned no state")
	}
	var state ServicemapModel
	if diags := resp.State.Get(ctx, &state); diags.HasError() {
		t.Fatal(diags)
	}
	if state.Id.ValueString() != "service-test" || !state.MetricCriteria.IsNull() || !state.MetricCriteria.ElementType(ctx).Equal(servicemapMetricCriterionType) {
		t.Fatalf("unexpected imported state: %#v", state)
	}
}

// TestUnitServicemapMapsPopulatedMonitorOnImportAndRead covers importing and refreshing a
// servicemap whose backend service availability monitor already has metric_criteria: both
// should map alert_type/metrics/metric_criteria from the response instead of leaving
// metric_criteria null or stale.
func TestUnitServicemapMapsPopulatedMonitorOnImportAndRead(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	responseBody := `{"id":"service-test","name":"test-service","nodeType":"SERVICE",` +
		`"availabilityThreshold":{"thresholdType":"count","thresholdLimit":1},` +
		`"serviceAvailabilityMonitor":{"alertType":"critical-alert","matchType":"ANY",` +
		`"metrics":["metric-a","metric-b"],` +
		`"metricCriteria":[{"metricName":"metric-a","operator":"=","components":["comp-1"]},` +
		`{"metricName":"metric-b","operator":"-"}]}}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet || req.URL.Path != "/api/v3/tenants/client-test/service-maps/service-test" {
			http.Error(w, "unexpected request", http.StatusBadRequest)

			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, responseBody)
	}))
	defer server.Close()

	r := &ServicemapResource{}
	r.apiClient = &client.OpsRampClient{BaseUrl: server.URL, TenantId: "client-test", Client: server.Client(), Scope: "CLIENT"}
	var schemaResp resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

	assertPopulatedMonitor := func(t *testing.T, state ServicemapModel) {
		t.Helper()

		if state.AlertType.ValueString() != "critical-alert" {
			t.Fatalf("alert_type = %q; want critical-alert", state.AlertType.ValueString())
		}
		if len(state.Metrics) != 2 || state.Metrics[0].ValueString() != "metric-a" || state.Metrics[1].ValueString() != "metric-b" {
			t.Fatalf("unexpected metrics: %#v", state.Metrics)
		}
		if state.MetricCriteria.IsNull() {
			t.Fatal("expected metric_criteria to be populated, got null")
		}

		var criteria []ServicemapMetricCriterionModel
		if diags := state.MetricCriteria.ElementsAs(ctx, &criteria, false); diags.HasError() {
			t.Fatal(diags)
		}
		if len(criteria) != 2 {
			t.Fatalf("expected 2 criteria, got %d", len(criteria))
		}

		byName := make(map[string]ServicemapMetricCriterionModel, len(criteria))
		for _, c := range criteria {
			byName[c.MetricName.ValueString()] = c
		}

		a, ok := byName["metric-a"]
		if !ok || a.Operator.ValueString() != "=" {
			t.Fatalf("unexpected criterion for metric-a: %#v", a)
		}
		var components []types.String
		if diags := a.Components.ElementsAs(ctx, &components, false); diags.HasError() {
			t.Fatal(diags)
		}
		if len(components) != 1 || components[0].ValueString() != "comp-1" {
			t.Fatalf("unexpected components for metric-a: %#v", components)
		}

		b, ok := byName["metric-b"]
		if !ok || b.Operator.ValueString() != "-" {
			t.Fatalf("unexpected criterion for metric-b: %#v", b)
		}
		if !b.Components.IsNull() {
			t.Fatalf("expected metric-b components to be null, got %#v", b.Components)
		}
	}

	t.Run("import", func(t *testing.T) {
		resp := resource.ImportStateResponse{State: tfsdk.State{Schema: schemaResp.Schema}}
		r.ImportState(ctx, resource.ImportStateRequest{ID: "service-test"}, &resp)
		if resp.Diagnostics.HasError() {
			t.Fatal(resp.Diagnostics)
		}

		var state ServicemapModel
		if diags := resp.State.Get(ctx, &state); diags.HasError() {
			t.Fatal(diags)
		}
		assertPopulatedMonitor(t, state)
	})

	t.Run("read", func(t *testing.T) {
		priorState := ServicemapModel{
			Id:             types.StringValue("service-test"),
			Name:           types.StringValue("test-service"),
			Type:           types.StringValue("Service"),
			ThresholdType:  types.StringValue("count"),
			ThresholdLimit: types.Int64Value(1),
			AlertType:      types.StringValue("availability"),
			MetricCriteria: types.SetNull(servicemapMetricCriterionType),
		}
		stateTF := tfsdk.State{Schema: schemaResp.Schema}
		if diags := stateTF.Set(ctx, &priorState); diags.HasError() {
			t.Fatal(diags)
		}

		resp := resource.ReadResponse{State: stateTF}
		r.Read(ctx, resource.ReadRequest{State: stateTF}, &resp)
		if resp.Diagnostics.HasError() {
			t.Fatal(resp.Diagnostics)
		}

		var state ServicemapModel
		if diags := resp.State.Get(ctx, &state); diags.HasError() {
			t.Fatal(diags)
		}
		assertPopulatedMonitor(t, state)
	})
}

// TestUnitServicemapReadPreservesMonitorWhenAPIOmitsIt covers refreshing a servicemap whose
// GET response omits "serviceAvailabilityMonitor" entirely (e.g. some node types don't echo it
// back even though it was accepted on create/update). Read must keep the previously stored
// alert_type/metrics/metric_criteria instead of wiping them out as "removed remotely".
func TestUnitServicemapReadPreservesMonitorWhenAPIOmitsIt(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet || req.URL.Path != "/api/v3/tenants/client-test/service-maps/service-test" {
			http.Error(w, "unexpected request", http.StatusBadRequest)

			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"id":"service-test","name":"test-child","nodeType":"Google",`+
			`"availabilityThreshold":{"thresholdType":"count","thresholdLimit":1}}`)
	}))
	defer server.Close()

	r := &ServicemapResource{}
	r.apiClient = &client.OpsRampClient{BaseUrl: server.URL, TenantId: "client-test", Client: server.Client(), Scope: "CLIENT"}
	var schemaResp resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

	priorState := ServicemapModel{
		Id:             types.StringValue("service-test"),
		Name:           types.StringValue("test-child"),
		Type:           types.StringValue("Resource"),
		ThresholdType:  types.StringValue("count"),
		ThresholdLimit: types.Int64Value(1),
		AlertType:      types.StringValue("critical-alert"),
		MetricCriteria: servicemapCriteria(t, []string{"metric-a"}, "-", types.SetNull(types.StringType)),
	}
	priorState.Metrics = append(priorState.Metrics, types.StringValue("metric-a"))
	stateTF := tfsdk.State{Schema: schemaResp.Schema}
	if diags := stateTF.Set(ctx, &priorState); diags.HasError() {
		t.Fatal(diags)
	}

	resp := resource.ReadResponse{State: stateTF}
	r.Read(ctx, resource.ReadRequest{State: stateTF}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatal(resp.Diagnostics)
	}

	var state ServicemapModel
	if diags := resp.State.Get(ctx, &state); diags.HasError() {
		t.Fatal(diags)
	}
	if !state.AlertType.Equal(priorState.AlertType) || !reflect.DeepEqual(state.Metrics, priorState.Metrics) ||
		!state.MetricCriteria.Equal(priorState.MetricCriteria) {
		t.Fatalf("expected monitor fields to be preserved, got: %#v", state)
	}
}

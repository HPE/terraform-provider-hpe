// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cluster

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/HPE/terraform-provider-hpe/morpheus/model"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clientfactory"
)

// testClusterID is the id the fake appliance hands back from POST /api/clusters.
const testClusterID = 848

// timeoutAttrTypes mirrors timeouts.AttributesAll, which the cluster schema uses.
var timeoutAttrTypes = map[string]attr.Type{
	"create": types.StringType,
	"read":   types.StringType,
	"update": types.StringType,
	"delete": types.StringType,
}

// newTimeouts builds a timeouts.Value with only "create" populated.
func newTimeouts(create string) timeouts.Value {
	return timeouts.Value{
		Object: types.ObjectValueMust(timeoutAttrTypes, map[string]attr.Value{
			"create": types.StringValue(create),
			"read":   types.StringNull(),
			"update": types.StringNull(),
			"delete": types.StringNull(),
		}),
	}
}

// nullOf returns a null attr.Value of the given type. Used to fill in the
// attributes of the nested objects the test does not care about.
func nullOf(ctx context.Context, t *testing.T, typ attr.Type) attr.Value {
	t.Helper()

	v, err := typ.ValueFromTerraform(ctx, tftypes.NewValue(typ.TerraformType(ctx), nil))
	if err != nil {
		t.Fatalf("null value of %s: %v", typ, err)
	}

	return v
}

// attrsWithNulls builds a full attribute map for an object type: every
// attribute null, then overrides applied.
func attrsWithNulls(
	ctx context.Context,
	t *testing.T,
	attrTypes map[string]attr.Type,
	overrides map[string]attr.Value,
) map[string]attr.Value {
	t.Helper()

	attrs := make(map[string]attr.Value, len(attrTypes))
	for name, typ := range attrTypes {
		attrs[name] = nullOf(ctx, t, typ)
	}

	for name, val := range overrides {
		if _, ok := attrTypes[name]; !ok {
			t.Fatalf("override %q is not an attribute of this object", name)
		}

		attrs[name] = val
	}

	return attrs
}

// unknownPaths walks a state value and returns the path of every unknown it
// finds. terraform-plugin-framework rejects a state containing any unknown when
// the response also carries an error diagnostic, so this must come back empty.
func unknownPaths(t *testing.T, v tftypes.Value) []string {
	t.Helper()

	var found []string

	err := tftypes.Walk(v, func(p *tftypes.AttributePath, val tftypes.Value) (bool, error) {
		if !val.IsKnown() {
			found = append(found, p.String())

			// Don't descend into an unknown value.
			return false, nil
		}

		return true, nil
	})
	if err != nil {
		t.Fatalf("walking state: %v", err)
	}

	return found
}

// TestUnitClusterCreateFailureLeavesNoUnknownState drives Create down each of its
// post-POST failure exits and asserts the state handed back alongside the error
// contains no unknown values.
//
// Regression test: Create used to write the whole plan to state as soon as the
// cluster id was known. Every computed-only attribute (uuid, service_url,
// cpu_placement_mode, cluster_type_code) is still unknown at that point, and
// neither error exit cleared them - TaintResourceState only sets "id", and the
// read-back failure path writes no state at all. Those unknowns survived into a
// state returned with an error, which the framework forbids, producing
// "Provider returned invalid result object after apply" and masking the real
// diagnostic.
func TestUnitClusterCreateFailureLeavesNoUnknownState(t *testing.T) {
	ctx := context.Background()
	sch := ClusterResourceSchema(ctx)

	cases := []struct {
		name string
		// status is what GET /api/clusters/{id} keeps reporting.
		status string
		// createTimeout is the HCL create timeout.
		createTimeout string
		// wantDetail is a fragment of the expected error diagnostic.
		wantDetail string
	}{
		{
			// The acceptance failure: the cluster never leaves "provisioning"
			// and the create timeout expires.
			name:          "create timeout while still provisioning",
			status:        clusterStatusProvisioning,
			createTimeout: "2s",
			wantDetail:    "provisioning failed",
		},
		{
			// The cluster reaches a terminal error status, so the poll stops
			// immediately with a permanent error.
			name:          "cluster reaches a terminal error status",
			status:        clusterStatusFailed,
			createTimeout: "5m",
			wantDetail:    "provisioning failed",
		},
		{
			// Provisioning succeeds, but the read-back fails. This exit used to
			// write no state at all, leaving the unknown-laden plan in place.
			name:          "read-back fails after successful provisioning",
			status:        clusterStatusOk,
			createTimeout: "5m",
			wantDetail:    "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// polled distinguishes the polling GETs from the getClusterAsState
			// GET that follows a successful poll, so the last case can fail only
			// the read-back.
			var polled atomic.Bool

			srv := httptest.NewServer(http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")

					if r.Method == http.MethodPost &&
						strings.HasSuffix(r.URL.Path, "/api/clusters") {
						fmt.Fprintf(w, `{"success":true,"cluster":{"id":%d}}`, testClusterID)

						return
					}

					// GET /api/clusters/{id}
					if tc.status == clusterStatusOk && polled.Swap(true) {
						// Second and later GETs are the read-back: fail them.
						w.WriteHeader(http.StatusInternalServerError)
						fmt.Fprint(w, `{"success":false,"msg":"read-back exploded"}`)

						return
					}

					fmt.Fprintf(w,
						`{"cluster":{"id":%d,"status":%q}}`, testClusterID, tc.status)
				}))
			defer srv.Close()

			r := &Resource{}
			cfgResp := &resource.ConfigureResponse{}
			r.Configure(ctx, resource.ConfigureRequest{
				ProviderData: clientfactory.New(model.MorpheusProviderModel{
					URL:         types.StringValue(srv.URL),
					AccessToken: types.StringValue("test-token"),
				}),
			}, cfgResp)
			if cfgResp.Diagnostics.HasError() {
				t.Fatalf("configure: %v", cfgResp.Diagnostics)
			}

			// The plan as the framework would hand it to Create: the practitioner's
			// values are known, every computed-only attribute is unknown.
			//
			// server and config_hvm are populated because the create payload is
			// only marshalable once the server config block has been built.
			serverTypes := NewServerValueNull().AttributeTypes(ctx)
			server := NewServerValueMust(serverTypes, attrsWithNulls(ctx, t, serverTypes,
				map[string]attr.Value{
					"name":                     types.StringValue("unit-test-cluster-worker"),
					"management_net_interface": types.StringValue("eth0"),
					"service_plan_id":          types.Int64Value(1),
					"ssh_username":             types.StringValue("cloud-user"),
					"ssh_port":                 types.Int64Value(22),
				}))

			configHvmTypes := NewConfigHvmValueNull().AttributeTypes(ctx)
			configHvm := NewConfigHvmValueMust(configHvmTypes, attrsWithNulls(ctx, t, configHvmTypes,
				map[string]attr.Value{
					"cpu_arch":     types.StringValue("x86_64"),
					"cpu_model":    types.StringValue("host-model"),
					"power_policy": types.StringValue("balanced"),
				}))

			plan := ClusterModel{
				CloudId:          types.Int64Value(1),
				Description:      types.StringValue("a test cluster"),
				GroupId:          types.Int64Value(265),
				LayoutId:         types.Int64Value(449),
				Name:             types.StringValue("unit-test-cluster"),
				Config:           types.DynamicNull(),
				ConfigHvm:        configHvm,
				Labels:           types.SetNull(types.StringType),
				Server:           server,
				Timeouts:         newTimeouts(tc.createTimeout),
				WorkflowId:       types.Int64Null(),
				Id:               types.Int64Unknown(),
				Uuid:             types.StringUnknown(),
				ServiceUrl:       types.StringUnknown(),
				CpuPlacementMode: types.StringUnknown(),
				ClusterTypeCode:  types.StringUnknown(),
			}

			// The config is the same, minus the unknowns: config never carries
			// values the practitioner did not write.
			cfg := plan
			cfg.Id = types.Int64Null()
			cfg.Uuid = types.StringNull()
			cfg.ServiceUrl = types.StringNull()
			cfg.CpuPlacementMode = types.StringNull()
			cfg.ClusterTypeCode = types.StringNull()

			tfPlan := tfsdk.Plan{Schema: sch}
			if diags := tfPlan.Set(ctx, &plan); diags.HasError() {
				t.Fatalf("build plan: %v", diags)
			}

			// tfsdk.Config has no Set, so round-trip the model through a Plan to
			// get a correctly typed raw value.
			cfgPlan := tfsdk.Plan{Schema: sch}
			if diags := cfgPlan.Set(ctx, &cfg); diags.HasError() {
				t.Fatalf("build config: %v", diags)
			}

			tfConfig := tfsdk.Config{Schema: sch, Raw: cfgPlan.Raw}

			// Sanity check: the plan really does contain the unknowns whose
			// leakage into state this test guards against. Without this the test
			// could pass against a plan that never had any.
			if got := unknownPaths(t, tfPlan.Raw); len(got) == 0 {
				t.Fatal("plan contains no unknowns; the test would prove nothing")
			}

			resp := &resource.CreateResponse{
				State: tfsdk.State{
					Schema: sch,
					Raw:    tftypes.NewValue(sch.Type().TerraformType(ctx), nil),
				},
			}

			r.Create(ctx, resource.CreateRequest{
				Plan:   tfPlan,
				Config: tfConfig,
			}, resp)

			// Every case here is a failure path, so an error diagnostic is
			// expected. If Create ever starts succeeding, the assertions below
			// stop meaning anything.
			if !resp.Diagnostics.HasError() {
				t.Fatalf("expected Create to fail, got diags: %v", resp.Diagnostics)
			}

			if tc.wantDetail != "" {
				var matched bool
				for _, d := range resp.Diagnostics.Errors() {
					if strings.Contains(d.Detail(), tc.wantDetail) {
						matched = true

						break
					}
				}
				if !matched {
					t.Errorf("no error diagnostic containing %q; got %v",
						tc.wantDetail, resp.Diagnostics.Errors())
				}
			}

			// The assertion that matters: state returned alongside an error must
			// contain no unknowns. Nulls are fine; unknowns are not.
			if got := unknownPaths(t, resp.State.Raw); len(got) > 0 {
				t.Errorf("state returned with an error contains %d unknown value(s) at %v; "+
					"terraform-plugin-framework rejects this with "+
					"\"Provider returned invalid result object after apply\"",
					len(got), got)
			}

			// The cluster exists on the appliance, so its id must be recorded -
			// otherwise the fix would trade an invalid state for an orphaned
			// cluster.
			var id types.Int64
			if diags := resp.State.GetAttribute(ctx, path.Root("id"), &id); diags.HasError() {
				t.Fatalf("reading id from state: %v", diags)
			}

			if id.IsNull() || id.ValueInt64() != testClusterID {
				t.Errorf("state id = %v, want %d", id, testClusterID)
			}
		})
	}
}

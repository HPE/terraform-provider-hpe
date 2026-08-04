// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cluster

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
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

// TestUnitClusterCreateFailureLeavesNoUnknownsInState is the regression test for
// the create-error path writing unknown computed values into state.
//
// Create used to persist the whole plan as soon as the cluster ID was known.
// The plan still holds unknown values for the computed-only attributes (uuid,
// service_url, cpu_placement_mode, and cluster_type_code when it is not
// configured), and neither error exit clears them: taintResourceState sets only
// "id", and the getClusterAsState failure path writes no state at all. Terraform
// then rejects the apply with
//
//	Provider returned invalid result object after apply ...
//	.uuid: was null, but now cty.UnknownVal(cty.String)
//
// which additionally masks the real diagnostic explaining why create failed.
//
// This cannot be covered by the acceptance tests: they only reach the taint path
// when provisioning genuinely fails on the appliance, so a healthy appliance
// makes them pass without exercising it at all.
func TestUnitClusterCreateFailureLeavesNoUnknownsInState(t *testing.T) {
	ctx := context.Background()
	sch := ClusterResourceSchema(ctx)

	timeoutsAttrTypes := map[string]attr.Type{
		"create": types.StringType,
		"read":   types.StringType,
		"update": types.StringType,
		"delete": types.StringType,
	}

	// nullOf builds a typed null for any attr.Type, so the nested server and
	// config_hvm objects can be populated without spelling out every one of
	// their attributes.
	nullOf := func(ty attr.Type) attr.Value {
		v, err := ty.ValueFromTerraform(
			ctx, tftypes.NewValue(ty.TerraformType(ctx), nil),
		)
		if err != nil {
			t.Fatalf("null value for %T: %v", ty, err)
		}

		return v
	}

	nullAttrs := func(attrTypes map[string]attr.Type) map[string]attr.Value {
		out := make(map[string]attr.Value, len(attrTypes))
		for name, ty := range attrTypes {
			out[name] = nullOf(ty)
		}

		return out
	}

	// The HVM shape from the failing acceptance test: a config_hvm block, a
	// server block, and no cluster_type_code (it is computed for this path).
	configHvmTypes := ConfigHvmValue{}.AttributeTypes(ctx)
	configHvmAttrs := nullAttrs(configHvmTypes)
	configHvmAttrs["cpu_arch"] = types.StringValue("x86_64")
	configHvmAttrs["cpu_model"] = types.StringValue("host-model")

	serverTypes := ServerValue{}.AttributeTypes(ctx)
	serverAttrs := nullAttrs(serverTypes)
	serverAttrs["service_plan_id"] = types.Int64Value(1)
	serverAttrs["ssh_username"] = types.StringValue("user")

	// The plan as Terraform presents it during create: configured values are
	// known, computed-only attributes are unknown.
	planModel := ClusterModel{
		Name:        types.StringValue("unit-test-cluster"),
		Description: types.StringValue("unit test"),
		CloudId:     types.Int64Value(1),
		GroupId:     types.Int64Value(1),
		LayoutId:    types.Int64Value(2),
		WorkflowId:  types.Int64Null(),
		Config:      types.DynamicNull(),
		ConfigHvm:   NewConfigHvmValueMust(configHvmTypes, configHvmAttrs),
		Server:      NewServerValueMust(serverTypes, serverAttrs),
		Labels:      types.SetNull(types.StringType),
		Timeouts:    timeouts.Value{Object: types.ObjectNull(timeoutsAttrTypes)},

		Id:               types.Int64Unknown(),
		Uuid:             types.StringUnknown(),
		ServiceUrl:       types.StringUnknown(),
		CpuPlacementMode: types.StringUnknown(),
		ClusterTypeCode:  types.StringUnknown(),
	}

	// The configuration never carries unknowns; computed-only attributes are
	// null there.
	configModel := planModel
	configModel.Id = types.Int64Null()
	configModel.Uuid = types.StringNull()
	configModel.ServiceUrl = types.StringNull()
	configModel.CpuPlacementMode = types.StringNull()
	configModel.ClusterTypeCode = types.StringNull()

	plan := tfsdk.Plan{Schema: sch}
	if diags := plan.Set(ctx, &planModel); diags.HasError() {
		t.Fatalf("build plan: %v", diags)
	}

	// tfsdk.Config has no Set, so round-trip the model through a Plan to build
	// the equivalent raw value.
	configRaw := tfsdk.Plan{Schema: sch}
	if diags := configRaw.Set(ctx, &configModel); diags.HasError() {
		t.Fatalf("build config: %v", diags)
	}

	config := tfsdk.Config{Schema: sch, Raw: configRaw.Raw}

	const clusterID = 999

	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			if r.Method == http.MethodPost {
				// Cluster is accepted and gets an ID.
				_, _ = w.Write([]byte(`{"success":true,"cluster":{"id":999}}`))

				return
			}

			// ...and then immediately reports a terminal error status, so the
			// poll gives up at once rather than waiting out the create timeout.
			_, _ = w.Write([]byte(`{"cluster":{"id":999,"status":"failed"}}`))
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

	// The framework hands Create a null state, exactly as fwserver does.
	resp := &resource.CreateResponse{
		State: tfsdk.State{
			Schema: sch,
			Raw:    tftypes.NewValue(sch.Type().TerraformType(ctx), nil),
		},
	}

	r.Create(ctx, resource.CreateRequest{Plan: plan, Config: config}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatalf("expected create to report an error, got: %v", resp.Diagnostics)
	}

	t.Logf("create diagnostics: %v", resp.Diagnostics)

	// The contract Terraform enforces: no unknown values may survive in state
	// once apply returns.
	var unknown []string

	walkErr := tftypes.Walk(resp.State.Raw,
		func(p *tftypes.AttributePath, v tftypes.Value) (bool, error) {
			if !v.IsKnown() {
				unknown = append(unknown, p.String())

				return false, nil
			}

			return true, nil
		})
	if walkErr != nil {
		t.Fatalf("walk state: %v", walkErr)
	}

	if len(unknown) > 0 {
		t.Errorf(
			"state retains %d unknown value(s) after a failed create, which "+
				"Terraform rejects with %q: %s",
			len(unknown),
			"Provider returned invalid result object after apply",
			strings.Join(unknown, ", "),
		)
	}

	// The ID must still be persisted, otherwise the partially created cluster
	// is orphaned instead of being tainted and cleaned up on the next apply.
	var gotID types.Int64
	if diags := resp.State.GetAttribute(
		ctx, path.Root("id"), &gotID,
	); diags.HasError() {
		t.Fatalf("read id from state: %v", diags)
	}

	if gotID.IsNull() || gotID.ValueInt64() != clusterID {
		t.Errorf("state id = %v, want %d (the created cluster must stay tracked)",
			gotID, clusterID)
	}
}

// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkserver

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TestUnitSecretOnlyRotationPlansLastUpdatedUnknown pins the plan of a
// rotation that changes nothing but a write-only value: the framework leaves
// every computed value known (the proposed plan equals the prior state), yet
// the re-send is an API update that advances last_updated, so ModifyPlan must
// mark it unknown beside write_only_revision or the apply is rejected as
// inconsistent. No digests are stored (nil private state), which is the
// "seed" decision: a re-send is planned, as for a changed digest.
func TestUnitSecretOnlyRotationPlansLastUpdatedUnknown(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	r := &networkServerResource{}
	var schemaResp resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)
	sch := schemaResp.Schema

	model := baseNsxPlan()
	model.Id = types.Int64Value(7)
	model.LastUpdated = types.StringValue("2026-11-13T10:00:00Z")
	model.WriteOnlyRevision = types.Int64Value(1700000000)

	state := tfsdk.State{Schema: sch}
	if d := state.Set(ctx, &model); d.HasError() {
		t.Fatal(d)
	}

	// The proposed plan is the prior state: write-only values are never in it.
	plan := tfsdk.Plan{Schema: sch, Raw: state.Raw}

	// The configuration carries the new token.
	cfgModel := model
	cfgModel.ServiceTokenWo = types.StringValue("rotated")
	cfgModel.LastUpdated = types.StringNull()
	cfgModel.WriteOnlyRevision = types.Int64Null()
	cfgModel.Id = types.Int64Null()
	cfgState := tfsdk.State{Schema: sch}
	if d := cfgState.Set(ctx, &cfgModel); d.HasError() {
		t.Fatal(d)
	}
	config := tfsdk.Config{Schema: sch, Raw: cfgState.Raw}

	resp := &resource.ModifyPlanResponse{Plan: plan}
	r.ModifyPlan(ctx, resource.ModifyPlanRequest{Config: config, Plan: plan, State: state}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatal(resp.Diagnostics)
	}

	for _, name := range []string{"write_only_revision", "last_updated"} {
		var v attr.Value
		if d := resp.Plan.GetAttribute(ctx, path.Root(name), &v); d.HasError() {
			t.Fatal(d)
		}

		if !v.IsUnknown() {
			t.Errorf("%s = %v after a secret-only rotation, want unknown", name, v)
		}
	}
}

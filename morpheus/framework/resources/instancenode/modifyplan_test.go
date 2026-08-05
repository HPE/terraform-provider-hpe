// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instancenode

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func makePlan(t *testing.T, m instanceNodeModel) tfsdk.Plan {
	t.Helper()

	ctx := context.Background()
	r := &Resource{}
	var schemaResp resource.SchemaResponse

	r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

	// Build a null timeouts object with the correct attribute types.
	if m.Timeouts.IsNull() || m.Timeouts.IsUnknown() {
		m.Timeouts = timeouts.Value{
			Object: types.ObjectNull(map[string]attr.Type{
				"create": types.StringType,
				"update": types.StringType,
				"delete": types.StringType,
			}),
		}
	}

	p := tfsdk.Plan{Schema: schemaResp.Schema}
	if diags := p.Set(ctx, &m); diags.HasError() {
		t.Fatalf("build plan: %v", diags)
	}

	return p
}

// TestModifyPlan_NoPoolOnNonMetal_NoDiagnostics verifies that ModifyPlan
// produces no diagnostics when resource_pool_id is null (virtual path).
func TestModifyPlan_NoPoolOnNonMetal_NoDiagnostics(t *testing.T) {
	t.Parallel()

	model := instanceNodeModel{
		InstanceID:       types.Int64Value(100),
		ResourcePoolID:   types.Int64Null(),
		PreProvisioned:   types.BoolNull(),
		SelectedServerID: types.Int64Null(),
		SshHost:          types.StringNull(),
		SshUsername:      types.StringNull(),
		SshPassword:      types.StringNull(),
		SshKeyPairID:     types.Int64Null(),
		WaitForIPAddress: types.BoolValue(false),
		ContainerID:      types.Int64Unknown(),
		ServerID:         types.Int64Unknown(),
		IPAddress:        types.StringUnknown(),
	}

	plan := makePlan(t, model)

	r := &Resource{}
	resp := &resource.ModifyPlanResponse{Plan: plan}

	r.ModifyPlan(context.Background(), resource.ModifyPlanRequest{
		Plan: plan,
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Errorf("expected no errors, got: %s", resp.Diagnostics.Errors())
	}
}

// TestModifyPlan_PoolUnknown_NoDiagnostics verifies that ModifyPlan produces
// no diagnostics when resource_pool_id is unknown (e.g. from a data source).
func TestModifyPlan_PoolUnknown_NoDiagnostics(t *testing.T) {
	t.Parallel()

	model := instanceNodeModel{
		InstanceID:       types.Int64Value(100),
		ResourcePoolID:   types.Int64Unknown(),
		PreProvisioned:   types.BoolNull(),
		SelectedServerID: types.Int64Null(),
		SshHost:          types.StringNull(),
		SshUsername:      types.StringNull(),
		SshPassword:      types.StringNull(),
		SshKeyPairID:     types.Int64Null(),
		WaitForIPAddress: types.BoolValue(false),
		ContainerID:      types.Int64Unknown(),
		ServerID:         types.Int64Unknown(),
		IPAddress:        types.StringUnknown(),
	}

	plan := makePlan(t, model)

	r := &Resource{}
	resp := &resource.ModifyPlanResponse{Plan: plan}

	r.ModifyPlan(context.Background(), resource.ModifyPlanRequest{
		Plan: plan,
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Errorf("expected no errors, got: %s", resp.Diagnostics.Errors())
	}
}

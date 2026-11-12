// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package settingbackup

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// planOf builds a tfsdk.Plan/Config carrying the model against the resource
// schema.
func stateValueOf(t *testing.T, m resourceModel) tfsdk.State {
	t.Helper()

	ctx := context.Background()
	s := tfsdk.State{Schema: resourceSchema(ctx)}

	if diags := s.Set(ctx, &m); diags.HasError() {
		t.Fatalf("building state: %v", diags)
	}

	return s
}

// runModifyPlan runs ModifyPlan with the given plan and config and returns the
// resulting plan model.
func runModifyPlan(t *testing.T, plan, config resourceModel) resourceModel {
	t.Helper()

	ctx := context.Background()
	r := &Resource{}

	planState := stateValueOf(t, plan)
	configState := stateValueOf(t, config)

	req := resource.ModifyPlanRequest{
		Plan:   tfsdk.Plan(planState),
		Config: tfsdk.Config(configState),
	}
	resp := resource.ModifyPlanResponse{
		Plan: tfsdk.Plan(planState),
	}

	r.ModifyPlan(ctx, req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("ModifyPlan: %v", resp.Diagnostics)
	}

	var out resourceModel
	if diags := resp.Plan.Get(ctx, &out); diags.HasError() {
		t.Fatalf("reading plan: %v", diags)
	}

	return out
}

// base is a fully-known plan/config with both retention attributes null, the
// starting point each case overrides.
func base() resourceModel {
	return resourceModel{
		SettingBackupModel: SettingBackupModel{
			ScheduledBackups:             types.BoolValue(true),
			CreateBackups:                types.BoolValue(true),
			BackupAppliance:              types.BoolValue(false),
			DefaultBackupStorageBucketId: types.Int64Null(),
			DefaultBackupScheduleId:      types.Int64Null(),
			RetentionCount:               types.Int64Null(),
			RetentionDays:                types.Int64Null(),
		},
		Id: types.StringValue(singletonID),
	}
}

// retention_count set mirrors into retention_days -- over the prior-state
// value UseStateForUnknown has put there, which is the case the mirror exists
// for: without it the plan would carry retention_days = 3 and the read-back
// would make the apply inconsistent.
func TestModifyPlanCountToDays(t *testing.T) {
	t.Parallel()

	config := base()
	config.RetentionCount = types.Int64Value(7)

	plan := config
	plan.RetentionDays = types.Int64Value(3)

	got := runModifyPlan(t, plan, config)

	if got.RetentionDays.ValueInt64() != 7 {
		t.Errorf("retention_days: got %v, want 7", got.RetentionDays)
	}
}

// retention_days set mirrors into retention_count.
func TestModifyPlanDaysToCount(t *testing.T) {
	t.Parallel()

	config := base()
	config.RetentionDays = types.Int64Value(14)

	plan := config
	plan.RetentionCount = types.Int64Value(3) // the prior state's value

	got := runModifyPlan(t, plan, config)

	if got.RetentionCount.ValueInt64() != 14 {
		t.Errorf("retention_count: got %v, want 14", got.RetentionCount)
	}
}

// An unknown retention_count propagates unknown into retention_days: the mirror
// is never more known than its source, so the apply is not inconsistent.
func TestModifyPlanUnknownPropagates(t *testing.T) {
	t.Parallel()

	config := base()
	config.RetentionCount = types.Int64Unknown()

	plan := config
	plan.RetentionDays = types.Int64Value(3) // the prior state's value

	got := runModifyPlan(t, plan, config)

	if !got.RetentionDays.IsUnknown() {
		t.Errorf("retention_days: got %v, want unknown", got.RetentionDays)
	}
}

// Neither configured: the prior-state values are left as they are.
func TestModifyPlanNeitherConfigured(t *testing.T) {
	t.Parallel()

	config := base()

	plan := config
	plan.RetentionCount = types.Int64Value(21)
	plan.RetentionDays = types.Int64Value(21)

	got := runModifyPlan(t, plan, config)

	if got.RetentionCount.ValueInt64() != 21 || got.RetentionDays.ValueInt64() != 21 {
		t.Errorf("retention untouched expected: count=%v days=%v", got.RetentionCount, got.RetentionDays)
	}
}

// Neither configured and the prior state holds only one of the two (a state
// that diverged from what mapSettings writes): the known value is mirrored,
// so the plan never carries one value while settingsBody sends another.
func TestModifyPlanNeitherConfiguredMirrorsTheKnownValue(t *testing.T) {
	t.Parallel()

	config := base()

	plan := config
	plan.RetentionCount = types.Int64Value(21)
	plan.RetentionDays = types.Int64Null()

	got := runModifyPlan(t, plan, config)

	if got.RetentionDays.ValueInt64() != 21 {
		t.Errorf("retention_days: got %v, want 21 (mirrored from retention_count)", got.RetentionDays)
	}

	plan = config
	plan.RetentionCount = types.Int64Null()
	plan.RetentionDays = types.Int64Value(9)

	got = runModifyPlan(t, plan, config)

	if got.RetentionCount.ValueInt64() != 9 {
		t.Errorf("retention_count: got %v, want 9 (mirrored from retention_days)", got.RetentionCount)
	}
}

// With backup_appliance = true planned, the mirror still runs. The guard
// itself is not exercised here: it needs a client, so it runs only once the
// resource is configured, and this harness's Resource is not; the tenancy
// decision is the tenancy package's, tested there.
func TestModifyPlanMirrorsWithBackupApplianceTrueOnUnconfiguredResource(t *testing.T) {
	t.Parallel()

	plan, config := base(), base()
	plan.BackupAppliance = types.BoolValue(true)
	config.BackupAppliance = types.BoolValue(true)
	plan.RetentionCount = types.Int64Value(5)
	config.RetentionCount = types.Int64Value(5)

	out := runModifyPlan(t, plan, config)

	if !out.BackupAppliance.Equal(types.BoolValue(true)) {
		t.Errorf("backup_appliance = %v, want true", out.BackupAppliance)
	}

	if !out.RetentionDays.Equal(types.Int64Value(5)) {
		t.Errorf("retention_days = %v, want 5 (mirrored)", out.RetentionDays)
	}
}

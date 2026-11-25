// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package settingbackup

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// upgradeFixture runs the upgrader on testdata/state_v0.json -- the stored
// instance of the SDKv2 resource -- after edit has changed its attributes. It
// decodes the raw state exactly as the framework does for an upgrade: against
// the version-0 schema, ignoring attributes the schema lacks.
func upgradeFixture(t *testing.T, edit func(attrs map[string]any)) (resourceModel, diag.Diagnostics) {
	t.Helper()

	ctx := context.Background()

	raw, err := os.ReadFile("testdata/state_v0.json")
	if err != nil {
		t.Fatal(err)
	}

	var instance struct {
		Attributes map[string]any `json:"attributes"`
	}

	if err := json.Unmarshal(raw, &instance); err != nil {
		t.Fatal(err)
	}

	if edit != nil {
		edit(instance.Attributes)
	}

	attrs, err := json.Marshal(instance.Attributes)
	if err != nil {
		t.Fatal(err)
	}

	prior := settingBackupSchemaV0()
	priorType := prior.Type().TerraformType(ctx).(tftypes.Object)

	rawState := tfprotov6.RawState{JSON: attrs}

	value, err := rawState.UnmarshalWithOpts(priorType, tfprotov6.UnmarshalOpts{
		ValueFromJSONOpts: tftypes.ValueFromJSONOpts{IgnoreUndefinedAttributes: true},
	})
	if err != nil {
		t.Fatalf("decode of the SDKv2 state against the version-0 schema failed: %v", err)
	}

	req := resource.UpgradeStateRequest{State: &tfsdk.State{Schema: prior, Raw: value}}
	resp := resource.UpgradeStateResponse{State: tfsdk.State{Schema: resourceSchema(ctx)}}

	upgradeStateV0toV1(ctx, req, &resp)

	var upgraded resourceModel
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(resp.State.Get(ctx, &upgraded)...)
	}

	return upgraded, resp.Diagnostics
}

func TestUpgradeStateV0toV1(t *testing.T) {
	t.Parallel()

	got, diags := upgradeFixture(t, nil)
	if diags.HasError() || diags.WarningsCount() > 0 {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	if got.Id.ValueString() != "1" {
		t.Errorf("id: got %v, want \"1\"", got.Id)
	}

	if !got.ScheduledBackups.ValueBool() || !got.CreateBackups.ValueBool() || got.BackupAppliance.ValueBool() {
		t.Errorf("flags: scheduled=%v create=%v appliance=%v",
			got.ScheduledBackups, got.CreateBackups, got.BackupAppliance)
	}

	if got.DefaultBackupStorageBucketId.ValueInt64() != 17 {
		t.Errorf("default_backup_storage_bucket_id: got %v, want 17", got.DefaultBackupStorageBucketId)
	}

	if got.DefaultBackupScheduleId.ValueInt64() != 3 {
		t.Errorf("default_backup_schedule_id: got %v, want 3", got.DefaultBackupScheduleId)
	}

	// The retention is carried into both attributes.
	if got.RetentionCount.ValueInt64() != 21 || got.RetentionDays.ValueInt64() != 21 {
		t.Errorf("retention: count=%v days=%v, want both 21", got.RetentionCount, got.RetentionDays)
	}
}

// An id of 0, the SDKv2 resource's unset value, maps to null.
func TestUpgradeStateV0toV1ZeroIdsToNull(t *testing.T) {
	t.Parallel()

	got, diags := upgradeFixture(t, func(attrs map[string]any) {
		attrs["default_backup_storage_bucket_id"] = 0
		attrs["default_backup_schedule_id"] = 0
	})
	if diags.HasError() {
		t.Fatal(diags)
	}

	if !got.DefaultBackupStorageBucketId.IsNull() {
		t.Errorf("default_backup_storage_bucket_id: got %v, want null", got.DefaultBackupStorageBucketId)
	}

	if !got.DefaultBackupScheduleId.IsNull() {
		t.Errorf("default_backup_schedule_id: got %v, want null", got.DefaultBackupScheduleId)
	}
}

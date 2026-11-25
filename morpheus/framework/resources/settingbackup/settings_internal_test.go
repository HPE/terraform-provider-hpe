// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package settingbackup

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// unconfiguredPlan is a create plan in which nothing is configured: every
// optional+computed attribute is unknown.
func unconfiguredPlan() SettingBackupModel {
	return SettingBackupModel{
		BackupAppliance:              types.BoolUnknown(),
		CreateBackups:                types.BoolUnknown(),
		DefaultBackupScheduleId:      types.Int64Unknown(),
		DefaultBackupStorageBucketId: types.Int64Unknown(),
		RetentionCount:               types.Int64Unknown(),
		RetentionDays:                types.Int64Unknown(),
		ScheduledBackups:             types.BoolUnknown(),
	}
}

// currentState is a resource state: the settings as mapped from a GET.
func currentState() SettingBackupModel {
	return SettingBackupModel{
		ScheduledBackups:             types.BoolValue(true),
		CreateBackups:                types.BoolValue(true),
		BackupAppliance:              types.BoolValue(false),
		DefaultBackupStorageBucketId: types.Int64Value(17),
		DefaultBackupScheduleId:      types.Int64Value(3),
		RetentionCount:               types.Int64Value(21),
		RetentionDays:                types.Int64Value(21),
	}
}

// wire is the backupSettings body as it goes on the wire.
func wire(t *testing.T, body settingsUpdate) map[string]any {
	t.Helper()

	raw, err := json.Marshal(sdk.UpdateBackupSettingsRequest{BackupSettings: &body})
	if err != nil {
		t.Fatal(err)
	}

	var out struct {
		BackupSettings map[string]any `json:"backupSettings"`
	}

	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatal(err)
	}

	if out.BackupSettings == nil {
		return map[string]any{}
	}

	return out.BackupSettings
}

func expectWire(t *testing.T, body settingsUpdate, want map[string]any) {
	t.Helper()

	if got := wire(t, body); !reflect.DeepEqual(got, want) {
		t.Errorf("body on the wire:\n got  %v\n want %v", got, want)
	}
}

// On create only the configured settings are sent; an unconfigured
// optional+computed setting is unknown and left out.
func TestSettingsBodyCreateSendsOnlyTheConfiguredSettings(t *testing.T) {
	t.Parallel()

	plan := unconfiguredPlan()
	plan.ScheduledBackups = types.BoolValue(true)
	plan.RetentionCount = types.Int64Value(21)

	body, send := settingsBody(&plan, nil)
	if !send {
		t.Fatal("a configured setting must be sent")
	}

	expectWire(t, body, map[string]any{
		"backupsEnabled": true,
		"retentionCount": float64(21),
	})
}

// The ids are sent as nested {id} objects when known.
func TestSettingsBodyCreateSendsIds(t *testing.T) {
	t.Parallel()

	plan := unconfiguredPlan()
	plan.DefaultBackupScheduleId = types.Int64Value(3)
	plan.DefaultBackupStorageBucketId = types.Int64Value(17)

	body, send := settingsBody(&plan, nil)
	if !send {
		t.Fatal("configured ids must be sent")
	}

	expectWire(t, body, map[string]any{
		"defaultSchedule":      map[string]any{"id": float64(3)},
		"defaultStorageBucket": map[string]any{"id": float64(17)},
	})
}

// A configuration that sets nothing sends nothing: Morpheus refuses an empty
// update ("No data provided").
func TestSettingsBodyNothingConfigured(t *testing.T) {
	t.Parallel()

	plan := unconfiguredPlan()

	body, send := settingsBody(&plan, nil)
	if send {
		t.Errorf("nothing to send, got %v", wire(t, body))
	}
}

// On an update, a setting is sent only when its planned value differs from the
// state; the ids are omitted when unchanged.
func TestSettingsBodyUpdateSendsOnlyWhatChanged(t *testing.T) {
	t.Parallel()

	state := currentState()

	plan := state
	plan.ScheduledBackups = types.BoolValue(false)
	plan.RetentionCount = types.Int64Value(7)
	plan.RetentionDays = types.Int64Value(7)

	body, send := settingsBody(&plan, &state)
	if !send {
		t.Fatal("changed settings must be sent")
	}

	expectWire(t, body, map[string]any{
		"backupsEnabled": false,
		"retentionCount": float64(7),
	})
}

// An update that changes nothing sends nothing.
func TestSettingsBodyUpdateWithNothingChanged(t *testing.T) {
	t.Parallel()

	state := currentState()
	plan := state

	if body, send := settingsBody(&plan, &state); send {
		t.Errorf("nothing changed, got %v", wire(t, body))
	}
}

// The retention comes from RetentionCount, which ModifyPlan mirrors from
// retention_days: a value set via retention_days reaches the one retentionCount
// setting through RetentionCount.
func TestSettingsBodyRetentionFromRetentionCount(t *testing.T) {
	t.Parallel()

	plan := unconfiguredPlan()
	plan.RetentionCount = types.Int64Value(14)
	plan.RetentionDays = types.Int64Value(14)

	body, send := settingsBody(&plan, nil)
	if !send {
		t.Fatal("a configured retention must be sent")
	}

	expectWire(t, body, map[string]any{"retentionCount": float64(14)})
}

func TestMapSettings(t *testing.T) {
	t.Parallel()

	bucketID, schedID := int64(17), int64(3)
	retention := int64(21)

	s := &backupSettings{
		BackupsEnabled:  boolPtr(true),
		CreateBackups:   boolPtr(true),
		BackupAppliance: boolPtr(false),
	}
	s.DefaultStorageBucket.Set(&sdk.ListBackupSettings200ResponseBackupSettingsDefaultStorageBucket{Id: &bucketID})
	s.DefaultSchedule.Set(&sdk.ListBackupSettings200ResponseBackupSettingsDefaultSchedule{Id: &schedID})
	s.RetentionCount.Set(&retention)

	var got SettingBackupModel
	if diags := mapSettings(&got, s); diags.HasError() {
		t.Fatal(diags)
	}

	want := SettingBackupModel{
		ScheduledBackups:             types.BoolValue(true),
		CreateBackups:                types.BoolValue(true),
		BackupAppliance:              types.BoolValue(false),
		DefaultBackupStorageBucketId: types.Int64Value(17),
		DefaultBackupScheduleId:      types.Int64Value(3),
		RetentionCount:               types.Int64Value(21),
		RetentionDays:                types.Int64Value(21),
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("mapped settings:\n got  %+v\n want %+v", got, want)
	}
}

// A null bucket and schedule map to null ids.
func TestMapSettingsNullObjects(t *testing.T) {
	t.Parallel()

	s := &backupSettings{
		BackupsEnabled: boolPtr(true),
	}

	var got SettingBackupModel
	if diags := mapSettings(&got, s); diags.HasError() {
		t.Fatal(diags)
	}

	if !got.DefaultBackupStorageBucketId.IsNull() {
		t.Errorf("default_backup_storage_bucket_id: got %v, want null", got.DefaultBackupStorageBucketId)
	}

	if !got.DefaultBackupScheduleId.IsNull() {
		t.Errorf("default_backup_schedule_id: got %v, want null", got.DefaultBackupScheduleId)
	}

	// A null retention is mirrored into both attributes.
	if !got.RetentionCount.IsNull() || !got.RetentionDays.IsNull() {
		t.Errorf("retention: count=%v days=%v, want both null", got.RetentionCount, got.RetentionDays)
	}
}

func boolPtr(b bool) *bool { return &b }

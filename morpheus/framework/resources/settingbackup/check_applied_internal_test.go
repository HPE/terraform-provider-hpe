// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package settingbackup

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// applied is a plan whose every setting the appliance reports back as sent.
func applied() SettingBackupModel {
	return SettingBackupModel{
		ScheduledBackups:             types.BoolValue(true),
		CreateBackups:                types.BoolValue(true),
		BackupAppliance:              types.BoolValue(true),
		DefaultBackupStorageBucketId: types.Int64Value(17),
		DefaultBackupScheduleId:      types.Int64Value(3),
		RetentionCount:               types.Int64Value(21),
		RetentionDays:                types.Int64Value(21),
	}
}

// fullBody is an update that sent every setting.
func fullBody() settingsUpdate {
	var body settingsUpdate

	body.BackupsEnabled = boolPtr(true)
	body.CreateBackups = boolPtr(true)
	body.BackupAppliance = boolPtr(true)
	body.RetentionCount.Set(int64Ptr(21))
	body.DefaultSchedule = &sdk.UpdateBackupSettingsRequestBackupSettingsDefaultSchedule{Id: int64Ptr(3)}
	body.DefaultStorageBucket = &sdk.UpdateBackupSettingsRequestBackupSettingsDefaultStorageBucket{Id: int64Ptr(17)}

	return body
}

func int64Ptr(v int64) *int64 { return &v }

// Every setting sent and read back as sent: nothing to report.
func TestCheckAppliedAllApplied(t *testing.T) {
	t.Parallel()

	plan, readBack := applied(), applied()

	if diags := checkApplied(fullBody(), &plan, &readBack, neverCalled(t)); diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
}

// neverCalled is the appliance hint for cases where backup_appliance applied:
// the hint costs a request, so it must not be asked for.
func neverCalled(t *testing.T) func() string {
	t.Helper()

	return func() string {
		t.Error("the backup_appliance hint was evaluated although the value applied")

		return ""
	}
}

// backup_appliance sent as true and read back as false is one error on that
// attribute, carrying the hint for the caller's tenancy -- the tenancy
// explanation for a tenant other than the master, as here.
func TestCheckAppliedBackupApplianceNotApplied(t *testing.T) {
	t.Parallel()

	plan, readBack := applied(), applied()
	readBack.BackupAppliance = types.BoolValue(false)

	diags := checkApplied(fullBody(), &plan, &readBack, func() string { return backupApplianceTenancyDetail })
	if got := diags.ErrorsCount(); got != 1 {
		t.Fatalf("errors = %d, want 1: %v", got, diags)
	}

	d := diags.Errors()[0]
	if !strings.Contains(d.Summary(), "backup_appliance was not applied") {
		t.Errorf("summary = %q", d.Summary())
	}

	if !strings.Contains(d.Detail(), "master tenant") || !strings.Contains(d.Detail(), "reports backup_appliance as false") {
		t.Errorf("detail = %q", d.Detail())
	}
}

// A setting that differs but was not sent is not the update's doing: the
// refresh reconciles it, and checkApplied says nothing.
func TestCheckAppliedIgnoresSettingsNotSent(t *testing.T) {
	t.Parallel()

	plan, readBack := applied(), applied()
	readBack.RetentionCount = types.Int64Value(7)
	readBack.DefaultBackupScheduleId = types.Int64Null()

	var body settingsUpdate
	body.BackupsEnabled = boolPtr(true)

	if diags := checkApplied(body, &plan, &readBack, neverCalled(t)); diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
}

// Each sent setting is checked on its own attribute.
func TestCheckAppliedReportsEachSentSetting(t *testing.T) {
	t.Parallel()

	plan, readBack := applied(), applied()
	readBack.RetentionCount = types.Int64Value(7)
	readBack.DefaultBackupScheduleId = types.Int64Null()
	readBack.DefaultBackupStorageBucketId = types.Int64Value(18)
	readBack.ScheduledBackups = types.BoolValue(false)

	diags := checkApplied(fullBody(), &plan, &readBack, neverCalled(t))
	if got := diags.ErrorsCount(); got != 4 {
		t.Fatalf("errors = %d, want 4: %v", got, diags)
	}

	var summaries []string
	for _, d := range diags.Errors() {
		summaries = append(summaries, d.Summary())
	}

	joined := strings.Join(summaries, "\n")
	for _, name := range []string{
		"scheduled_backups", "retention_count", "default_backup_schedule_id", "default_backup_storage_bucket_id",
	} {
		if !strings.Contains(joined, name+" was not applied") {
			t.Errorf("no error for %s in:\n%s", name, joined)
		}
	}

	if strings.Contains(joined, "backup_appliance") || strings.Contains(joined, "create_backups") {
		t.Errorf("applied settings reported:\n%s", joined)
	}
}

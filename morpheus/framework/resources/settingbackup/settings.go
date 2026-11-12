// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package settingbackup

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/tenancy"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

// backupSettings is the settings object the GET returns; settingsUpdate is the
// one the PUT takes.
type (
	backupSettings = sdk.ListBackupSettings200ResponseBackupSettings
	settingsUpdate = sdk.UpdateBackupSettingsRequestBackupSettings
)

// getSettings reads the backup settings. The error is formatted for a
// diagnostic already; the response is returned for its status code.
func getSettings(ctx context.Context, client *sdk.APIClient) (*backupSettings, *http.Response, error) {
	result, httpResp, err := client.BackupSettingsAPI.ListBackupSettings(ctx).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		return nil, httpResp, err
	}

	if result == nil || result.BackupSettings == nil {
		return nil, httpResp, errors.New("the response carries no backupSettings object")
	}

	return result.BackupSettings, httpResp, nil
}

// updateSettings sends the settings, and the answer must confirm them.
// Morpheus answers an update it refuses -- a schedule or storage bucket that
// does not exist -- with HTTP 200 and {"success": false, "errors": {...}}, so
// the body decides. It renders the success flag on every answer, so an answer
// without it -- no body, or a body that is not Morpheus's -- confirms nothing
// and is an error too. The error is formatted for a diagnostic already.
func updateSettings(ctx context.Context, client *sdk.APIClient, body settingsUpdate) error {
	result, httpResp, err := client.BackupSettingsAPI.UpdateBackupSettings(ctx).
		UpdateBackupSettingsRequest(sdk.UpdateBackupSettingsRequest{BackupSettings: &body}).
		Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		return err
	}

	// The generated client leaves the result nil when the body is empty.
	if result == nil {
		return errors.New("the appliance answered the update with an empty body, " +
			"so the settings cannot be confirmed as applied")
	}

	if result.Success == nil {
		return errors.New("the appliance's answer to the update has no success flag, " +
			"so the settings cannot be confirmed as applied")
	}

	if !*result.Success {
		return fmt.Errorf("the appliance refused the settings: %s",
			errfmt.RefusalMessage(result.AdditionalProperties))
	}

	return nil
}

// known reports whether a value is neither null nor unknown.
func known(v attr.Value) bool {
	return !v.IsNull() && !v.IsUnknown()
}

// settingsBody builds the PUT body from the plan: on create (state nil) every
// setting the plan knows -- the configured ones, since an unconfigured
// optional+computed setting is unknown on create -- and on update the ones
// whose planned value differs from the state. It reports whether it carries
// anything: Morpheus refuses an empty update ("No data provided").
//
// The retention comes from plan.RetentionCount, which ModifyPlan has mirrored
// from retention_days when that was the attribute set: it is the single source
// of truth for the one retentionCount setting. The default schedule and
// storage bucket are sent as nested {id} objects.
func settingsBody(plan, state *SettingBackupModel) (settingsUpdate, bool) {
	var body settingsUpdate

	create := state == nil
	if create {
		state = &SettingBackupModel{}
	}

	send := false

	// pick reports whether a planned value is to be sent: known and, on an
	// update, different from the state.
	pick := func(p, s attr.Value) bool {
		if !known(p) || (!create && p.Equal(s)) {
			return false
		}

		send = true

		return true
	}

	boolean := func(p, s types.Bool) *bool {
		if !pick(p, s) {
			return nil
		}

		return p.ValueBoolPointer()
	}

	body.BackupsEnabled = boolean(plan.ScheduledBackups, state.ScheduledBackups)
	body.CreateBackups = boolean(plan.CreateBackups, state.CreateBackups)
	body.BackupAppliance = boolean(plan.BackupAppliance, state.BackupAppliance)

	if pick(plan.RetentionCount, state.RetentionCount) {
		body.RetentionCount.Set(plan.RetentionCount.ValueInt64Pointer())
	}

	if pick(plan.DefaultBackupScheduleId, state.DefaultBackupScheduleId) {
		body.DefaultSchedule = &sdk.UpdateBackupSettingsRequestBackupSettingsDefaultSchedule{
			Id: plan.DefaultBackupScheduleId.ValueInt64Pointer(),
		}
	}

	if pick(plan.DefaultBackupStorageBucketId, state.DefaultBackupStorageBucketId) {
		body.DefaultStorageBucket = &sdk.UpdateBackupSettingsRequestBackupSettingsDefaultStorageBucket{
			Id: plan.DefaultBackupStorageBucketId.ValueInt64Pointer(),
		}
	}

	return body, send
}

// mapSettings copies the backup settings into the model. The default schedule
// and storage bucket are nested objects: their id when the appliance has one,
// null when it does not. retention_count and retention_days mirror the one
// retentionCount setting, so both carry its value.
func mapSettings(m *SettingBackupModel, s *backupSettings) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ScheduledBackups = convert.BoolToType(s.BackupsEnabled)
	m.CreateBackups = convert.BoolToType(s.CreateBackups)
	m.BackupAppliance = convert.BoolToType(s.BackupAppliance)

	if bucket := s.DefaultStorageBucket.Get(); bucket != nil {
		m.DefaultBackupStorageBucketId = convert.Int64ToType(bucket.Id)
	} else {
		m.DefaultBackupStorageBucketId = types.Int64Null()
	}

	if sched := s.DefaultSchedule.Get(); sched != nil {
		m.DefaultBackupScheduleId = convert.Int64ToType(sched.Id)
	} else {
		m.DefaultBackupScheduleId = types.Int64Null()
	}

	retention := convert.Int64ToType(s.RetentionCount.Get())
	m.RetentionCount = retention
	m.RetentionDays = retention

	return diags
}

// checkApplied compares each setting the update sent with the value the
// appliance reports afterwards, and returns an error diagnostic for every one
// that was not applied. Morpheus answers an update it accepts but does not act
// on with success: true -- backup_appliance from a tenant that is not the
// master tenant is the known case (see ModifyPlan), and the appliance backup
// failing to be created is another -- and without this check Terraform would
// report the mismatch as an inconsistent result with no reason. body is the
// request as sent; plan holds the values asked for and readBack the ones read
// after the update. applianceHint explains an unapplied backup_appliance; it
// is a function because the explanation depends on the caller's tenancy,
// which is only worth asking the appliance about when the value did not take.
func checkApplied(
	body settingsUpdate,
	plan, readBack *SettingBackupModel,
	applianceHint func() string,
) diag.Diagnostics {
	var diags diag.Diagnostics

	check := func(sent bool, name string, want, got attr.Value, hint func() string) {
		if !sent || want.Equal(got) {
			return
		}

		detail := fmt.Sprintf("Morpheus accepted the update but reports %s as %s rather than %s.",
			name, got.String(), want.String())
		if hint != nil {
			detail += " " + hint()
		}

		diags.AddAttributeError(path.Root(name), name+" was not applied", detail)
	}

	check(body.BackupsEnabled != nil, "scheduled_backups", plan.ScheduledBackups, readBack.ScheduledBackups, nil)
	check(body.CreateBackups != nil, "create_backups", plan.CreateBackups, readBack.CreateBackups, nil)
	check(body.BackupAppliance != nil, "backup_appliance", plan.BackupAppliance, readBack.BackupAppliance,
		applianceHint)
	check(body.RetentionCount.IsSet(), "retention_count", plan.RetentionCount, readBack.RetentionCount, nil)
	check(body.DefaultSchedule != nil, "default_backup_schedule_id",
		plan.DefaultBackupScheduleId, readBack.DefaultBackupScheduleId, nil)
	check(body.DefaultStorageBucket != nil, "default_backup_storage_bucket_id",
		plan.DefaultBackupStorageBucketId, readBack.DefaultBackupStorageBucketId, nil)

	return diags
}

// backupApplianceHint explains a backup_appliance the appliance accepted but
// did not apply, by the caller's tenancy: for a tenant other than the master,
// Morpheus applied it to the master tenant's appliance backup (see
// ModifyPlan); for the master tenant, Morpheus failed to create or enable the
// appliance backup, which it only logs. The tenancy is the one ModifyPlan's
// guard cached when it ran, so this rarely costs a request.
func backupApplianceHint(ctx context.Context, client *sdk.APIClient) func() string {
	return func() string {
		isMaster, err := tenancy.CallerIsMaster(ctx, client)

		switch {
		case err != nil:
			return "If the provider is not configured as the master tenant: " + backupApplianceTenancyDetail +
				" Otherwise " + backupApplianceFailedDetail + " (The tenant could not be determined: " +
				err.Error() + ".)"
		case !isMaster:
			return backupApplianceTenancyDetail
		default:
			return backupApplianceFailedDetail
		}
	}
}

// backupApplianceFailedDetail is the hint for the master tenant: Morpheus
// creates or enables the appliance backup inside a try/catch that only logs,
// so the API reports success and the read-back shows it did not happen.
const backupApplianceFailedDetail = "Morpheus did not create or enable the appliance backup (it logs the reason " +
	"and reports success all the same; it may have been unable to read its own database connection details). " +
	"Check the appliance's log and the Backups page."

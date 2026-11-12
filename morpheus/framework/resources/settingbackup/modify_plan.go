// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package settingbackup

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/tenancy"
)

// ModifyPlan does two things.
//
// It mirrors retention_count and retention_days, the two attributes over one
// setting (retentionCount): whichever one the configuration sets, the other is
// planned to the same value, so settingsBody has a single source of truth in
// plan.RetentionCount. The ConfigValidator stops both being set, so at most
// one is configured. The planned value is copied whether known or unknown, so
// a value that is only known after apply leaves the mirror unknown too and the
// apply is never reported inconsistent. When neither is configured both come
// from prior state through UseStateForUnknown (or are unknown on create); if
// only one of them is known, it is mirrored into the other.
//
// It also refuses backup_appliance = true from a tenant that is not the
// master tenant. Morpheus keeps one appliance backup for the whole appliance
// and reports it to each tenant as that tenant's own: another tenant's true
// is accepted (HTTP 200, success true) but applied to the master tenant's
// appliance backup, which the tenant then reads as false, and its false is
// ignored. Sending it would change the master tenant's backup and fail the
// apply on the read-back, so it is refused before anything is sent. The
// guard runs only for a planned true (false is harmless: a tenant that does
// not own the appliance backup reads false and plans nothing) and fails
// open with a warning when the caller's tenancy cannot be determined; the
// read-back check in Create and Update then reports what the appliance did.
func (r *Resource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	// A destroy plan has a null Plan.Raw and nothing to mirror or guard.
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan, config resourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	switch {
	case !config.RetentionCount.IsNull():
		plan.RetentionDays = plan.RetentionCount
	case !config.RetentionDays.IsNull():
		plan.RetentionCount = plan.RetentionDays
	// Neither configured: both come from prior state, which mapSettings and the
	// state upgrade always write equal. Should they ever diverge (a hand-edited
	// state), the known one wins, so settingsBody never sends one value while
	// the state keeps another.
	case known(plan.RetentionCount) && !known(plan.RetentionDays):
		plan.RetentionDays = plan.RetentionCount
	case known(plan.RetentionDays) && !known(plan.RetentionCount):
		plan.RetentionCount = plan.RetentionDays
	}

	resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.guardBackupAppliance(ctx, plan, resp)
}

// guardBackupAppliance refuses a planned backup_appliance = true from a tenant
// that is not the master tenant; see ModifyPlan. The client is built only when
// the guard has something to check, and only when the resource is configured:
// the framework does not configure a resource for every RPC, and the mirror
// above must still run without a client. The check has the contract of the
// tenancy package's guards -- whoami once per provider process, and a warning
// rather than an error when the caller's tenancy cannot be determined -- with
// a message that says what Morpheus does with the value.
func (r *Resource) guardBackupAppliance(
	ctx context.Context,
	plan resourceModel,
	resp *resource.ModifyPlanResponse,
) {
	if !known(plan.BackupAppliance) || !plan.BackupAppliance.ValueBool() || !r.Configured() {
		return
	}

	attr := path.Root("backup_appliance")

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create Morpheus client", err.Error())

		return
	}

	isMaster, err := tenancy.CallerIsMaster(ctx, client)
	if err != nil {
		resp.Diagnostics.AddAttributeWarning(
			attr,
			"Could not validate backup_appliance at plan time",
			"Unable to determine whether the authenticated tenant is the master tenant, so "+
				"backup_appliance = true was not validated. If it is not, Morpheus applies the value to "+
				"the master tenant's appliance backup and the apply fails on the read-back. Error: "+err.Error(),
		)

		return
	}

	if !isMaster {
		resp.Diagnostics.AddAttributeError(
			attr,
			"backup_appliance can only be set to true by the master tenant",
			backupApplianceTenancyDetail,
		)
	}
}

// backupApplianceTenancyDetail explains the refusal of backup_appliance = true
// from a tenant that is not the master tenant, at plan time and when the
// read-back shows the value was not applied.
const backupApplianceTenancyDetail = "Morpheus keeps one appliance backup for the whole appliance and manages it " +
	"for the master tenant only. Another tenant's request to enable it is accepted but applied to the master " +
	"tenant's appliance backup, which that tenant then reads as false, and its request to disable it is ignored. " +
	"Set backup_appliance from the master tenant, or leave it unset."

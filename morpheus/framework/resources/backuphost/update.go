// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package backuphost

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
)

func (r *backupHostResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan BackupHostModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// ssh_password_wo is write-only, so its value is only available from config.
	// Prior state is used to detect whether ssh_password_wo_version changed.
	var config BackupHostModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	var priorState BackupHostModel
	resp.Diagnostics.Append(req.State.Get(ctx, &priorState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.Id.ValueInt64()

	body := sdk.UpdateBackupsRequestBackup{
		Name: plan.Name.ValueStringPointer(),
	}
	if !plan.JobId.IsNull() {
		body.BackupJobId = plan.JobId.ValueInt64Pointer()
	}
	if !plan.Enabled.IsNull() {
		body.Enabled = plan.Enabled.ValueBoolPointer()
	}

	// Attributes removed from the configuration since the last apply, to be
	// cleared: sent when known, null when being cleared, not sent otherwise.
	pending, pdiags := clearable.Compute(ctx, req.Config, req.State, req.Private, clearableAttributes)
	resp.Diagnostics.Append(pdiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	clearable.NullableInt64(pending, "storage_provider_id", plan.StorageProviderId, &body.StorageProviderId)
	clearable.NullableString(pending, "path", plan.Path, &body.TargetPath)
	clearable.NullableString(pending, "ssh_username", plan.SshUsername, &body.TargetUsername)
	// The write-only ssh_password_wo cannot be diffed by Terraform; ModifyPlan
	// marks write_only_revision unknown when its digest changed (or the
	// deprecated marker did), and that is when it is re-sent. See write_only.go.
	resend := wodigest.Resend(ctx, req.Plan)
	if resend && !config.SshPasswordWo.IsNull() && !config.SshPasswordWo.IsUnknown() {
		body.TargetPassword = config.SshPasswordWo.ValueStringPointer()
	}

	_, httpResp, err := client.BackupsAPI.UpdateBackups(ctx, id).UpdateBackupsRequest(sdk.UpdateBackupsRequest{
		Backup: body,
	}).Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(updateOperation, errfmt.ErrMsg(err, httpResp))

		return
	}

	state, diags := getBackupAsState(ctx, id, client, plan)
	state.WriteOnlyRevision = priorState.WriteOnlyRevision
	if resend {
		state.WriteOnlyRevision = wodigest.NewRevision()
	}
	// Always record the digests of the configured values, sent or not.
	resp.Diagnostics.Append(wo.Store(ctx, resp.Private, req.Config)...)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(clearable.Record(ctx, req.Config, resp.Private, clearableAttributes)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

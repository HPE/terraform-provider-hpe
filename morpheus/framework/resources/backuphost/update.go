// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package backuphost

import (
	"context"
	"net/http"

	sdk "github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
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
	if !plan.StorageProviderId.IsUnknown() {
		body.StorageProviderId = *sdk.NewNullableInt64(plan.StorageProviderId.ValueInt64Pointer())
	}
	if !plan.Path.IsNull() {
		body.TargetPath = plan.Path.ValueStringPointer()
	}
	if !plan.SshUsername.IsNull() {
		body.TargetUsername = plan.SshUsername.ValueStringPointer()
	}
	// The write-only ssh_password_wo cannot be diffed by Terraform, so it is only
	// re-sent when ssh_password_wo_version changes.
	if !plan.SshPasswordWoVersion.Equal(priorState.SshPasswordWoVersion) {
		if config.SshPasswordWo.IsNull() {
			resp.Diagnostics.AddError(
				updateOperation,
				"'ssh_password_wo_version' changed, but 'ssh_password_wo' is not set",
			)

			return
		}
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
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

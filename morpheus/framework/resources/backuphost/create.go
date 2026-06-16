// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package backuphost

import (
	"context"
	"net/http"

	sdk "github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
)

func (r *backupHostResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
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
	var config BackupHostModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	host := sdk.BackupServerHost{
		Name:         plan.Name.ValueString(),
		LocationType: "server",
		BackupType:   plan.BackupTypeCode.ValueString(),
		JobAction:    "addTo",
		ServerId:     plan.HostId.ValueInt64(),
		JobId:        plan.JobId.ValueInt64Pointer(),
	}
	if !plan.StorageProviderId.IsNull() && !plan.StorageProviderId.IsUnknown() {
		host.Target = plan.StorageProviderId.ValueInt64Pointer()
	}
	if !plan.Path.IsNull() {
		host.TargetPath = plan.Path.ValueStringPointer()
	}
	if !plan.SshUsername.IsNull() {
		host.SshUsername = plan.SshUsername.ValueStringPointer()
	}
	if !config.SshPasswordWo.IsNull() {
		host.SshPassword = config.SshPasswordWo.ValueStringPointer()
	}

	backup := sdk.AddBackupsRequestBackup{
		BackupServerHost: &host,
	}

	result, httpResp, err := client.BackupsAPI.AddBackups(ctx).AddBackupsRequest(sdk.AddBackupsRequest{
		Backup: backup,
	}).Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(createOperation, errfmt.ErrMsg(err, httpResp))

		return
	}

	if result.Backup == nil || result.Backup.Id == nil {
		resp.Diagnostics.AddError(createOperation, "Backup ID is nil in the create response")

		return
	}
	id := *result.Backup.Id

	state, diags := getBackupAsState(ctx, id, client, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		// The backup was created, but reading it back failed. Taint the resource
		// so the created backup is not leaked from Terraform's perspective.
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "backup_host",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

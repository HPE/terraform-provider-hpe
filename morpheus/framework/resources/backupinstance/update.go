// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package backupinstance

import (
	"context"
	"net/http"

	sdk "github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

func (r *backupInstanceResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan BackupInstanceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
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

	_, httpResp, err := client.BackupsAPI.UpdateBackups(ctx, id).UpdateBackupsRequest(sdk.UpdateBackupsRequest{
		Backup: body,
	}).Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(updateOperation, errfmt.ErrMsg(err, httpResp))

		return
	}

	state, diags := getBackupAsState(ctx, id, client)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

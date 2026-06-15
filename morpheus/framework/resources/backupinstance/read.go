// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package backupinstance

import (
	"context"
	"fmt"
	"net/http"

	sdk "github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

func (r *backupInstanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state BackupInstanceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueInt64()

	state, diags := getBackupAsState(ctx, id, client)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// getBackupAsState performs a read of the backup by ID and returns the
// resulting state. It is used by Create, Read and Update so that state is
// always populated from a fresh read rather than from the create/update
// responses.
func getBackupAsState(
	ctx context.Context,
	id int64,
	client *sdk.APIClient,
) (BackupInstanceModel, diag.Diagnostics) {
	var state BackupInstanceModel
	var diags diag.Diagnostics

	result, hresp, err := client.BackupsAPI.GetBackups(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		diags.AddError(readOperation, errfmt.ErrMsg(err, hresp))

		return state, diags
	}

	b := result.Backup
	if b == nil {
		diags.AddError(
			readOperation,
			fmt.Sprintf("backup %d GET returned no backup", id),
		)

		return state, diags
	}

	state.Id = convert.Int64ToType(b.Id)
	state.Name = convert.StrToType(b.Name)
	state.Enabled = convert.BoolToType(b.Enabled)

	if v := b.ContainerId.Get(); v != nil {
		state.ContainerId = convert.Int64ToType(v)
	}

	if b.BackupType != nil {
		state.BackupTypeCode = convert.StrToType(b.BackupType.Code)
	}

	if b.Instance != nil {
		state.InstanceId = convert.Int64ToType(b.Instance.Id)
	}

	if b.Job != nil {
		state.JobId = convert.Int64ToType(b.Job.Id)
	}

	if b.StorageProvider != nil {
		state.StorageProviderId = convert.Int64ToType(b.StorageProvider.Id)
	}

	return state, diags
}

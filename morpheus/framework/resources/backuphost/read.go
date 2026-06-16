// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package backuphost

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

func (r *backupHostResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var prior BackupHostModel
	resp.Diagnostics.Append(req.State.Get(ctx, &prior)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := prior.Id.ValueInt64()

	state, diags := getBackupAsState(ctx, id, client, prior)
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
	plan BackupHostModel,
) (BackupHostModel, diag.Diagnostics) {
	var state BackupHostModel
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

	if b.BackupType != nil {
		state.BackupTypeCode = convert.StrToType(b.BackupType.Code)
	}

	if b.Job != nil {
		state.JobId = convert.Int64ToType(b.Job.Id)
	}

	if b.StorageProvider != nil {
		state.StorageProviderId = convert.Int64ToType(b.StorageProvider.Id)
	}

	// host_id and path are not present on every backup type, so fall back to the
	// planned value when the API omits them.
	state.HostId = plan.HostId
	if b.Server != nil {
		state.HostId = convert.Int64ToType(b.Server.Id)
	}

	state.Path = plan.Path
	if v := b.TargetPath.Get(); v != nil {
		state.Path = convert.StrToType(v)
	}

	// ssh_username can be read back from the API. ssh_password_wo is write-only
	// (never stored in state) and ssh_password_wo_version is not returned by the
	// API, so preserve it from plan/state.
	state.SshUsername = plan.SshUsername
	if v := b.TargetUsername.Get(); v != nil {
		state.SshUsername = convert.StrToType(v)
	}
	state.SshPasswordWoVersion = plan.SshPasswordWoVersion

	return state, diags
}

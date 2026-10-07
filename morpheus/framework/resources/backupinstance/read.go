// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package backupinstance

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
	"github.com/HPE/terraform-provider-hpe/utils/schemadefaults"
)

func (r *backupInstanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var priorState BackupInstanceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &priorState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := priorState.Id.ValueInt64()

	state, found, diags := getBackupAsState(ctx, id, client, priorState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The backup was deleted out of band (Morpheus returned 404); remove it
	// from state so the next plan recreates it instead of erroring.
	if !found {
		tflog.Info(ctx, fmt.Sprintf("backup instance %d not found, removing from state", id))
		resp.State.RemoveResource(ctx)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	// Fill any schema-declared default the API omitted (null in state) so an
	// imported resource does not plan a change nobody made. MORPH-16192.
	resp.Diagnostics.Append(
		schemadefaults.Apply(ctx, BackupInstanceResourceSchema(ctx), &resp.State)...,
	)
}

// getBackupAsState performs a read of the backup by ID and returns the
// resulting state. It is used by Create, Read and Update so that state is
// always populated from a fresh read rather than from the create/update
// responses.
func getBackupAsState(
	ctx context.Context,
	id int64,
	client *sdk.APIClient,
	plan BackupInstanceModel,
) (BackupInstanceModel, bool, diag.Diagnostics) {
	var state BackupInstanceModel
	var diags diag.Diagnostics

	result, hresp, err := client.BackupsAPI.GetBackups(ctx, id).Execute()
	if errfmt.IsNotFound(hresp) {
		// Signal not-found (without an error) so Read can remove the resource
		// from state. Create/Update treat this as an error via the found flag.
		return state, false, diags
	}
	if err != nil || hresp.StatusCode != http.StatusOK {
		diags.AddError(readOperation, errfmt.ErrMsg(err, hresp))

		return state, false, diags
	}

	b := result.Backup
	if b == nil {
		diags.AddError(
			readOperation,
			fmt.Sprintf("backup %d GET returned no backup", id),
		)

		return state, false, diags
	}

	return mapBackupToState(b, plan), true, diags
}

// mapBackupToState maps an SDK backup struct onto resource state, preserving
// planned/prior values for optional fields the API does not reliably return.
// It is split from getBackupAsState so the null-safety of those optional fields
// can be unit-tested without a live API.
func mapBackupToState(
	b *sdk.GetBackups200ResponseBackup,
	plan BackupInstanceModel,
) BackupInstanceModel {
	var state BackupInstanceModel

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

	// The API may return a null storage provider even when one was requested
	// (it can fall back to the system default), so preserve the configured
	// value rather than nulling it out after apply.
	var storageProviderID *int64
	if b.StorageProvider != nil {
		storageProviderID = b.StorageProvider.Id
	}
	state.StorageProviderId = convert.Int64OrPlan(storageProviderID, plan.StorageProviderId)

	return state
}

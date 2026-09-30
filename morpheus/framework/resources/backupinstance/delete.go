// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package backupinstance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

func (r *backupInstanceResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
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

	// A 404 means the backup is already gone, which is a successful delete;
	// a 200 whose body says success:false is a refusal (MORPH-17987).
	result, httpResp, err := client.BackupsAPI.RemoveBackups(ctx, id).Execute()
	if err := errfmt.CheckDeleted(result, err, httpResp); err != nil {
		resp.Diagnostics.AddError(deleteOperation, err.Error())

		return
	}
}

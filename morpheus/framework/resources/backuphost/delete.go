// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package backuphost

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

func (r *backupHostResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state BackupHostModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueInt64()

	_, httpResp, err := client.BackupsAPI.RemoveBackups(ctx, id).Execute()
	// A 404 means the backup is already gone, so treat it as a successful
	// delete rather than an error.
	if errfmt.IsNotFound(httpResp) {
		return
	}
	if err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(deleteOperation, errfmt.ErrMsg(err, httpResp))

		return
	}
}

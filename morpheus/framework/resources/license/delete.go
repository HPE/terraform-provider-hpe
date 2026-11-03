// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package license

import (
	"context"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/provisionhint"
)

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state LicenseModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	keyID := strings.TrimSpace(state.Id.ValueString())
	if keyID == "" {
		// Without a key id the API deactivates every installed licence; refuse.
		resp.Diagnostics.AddError(deleteOperation,
			"refusing to uninstall without a key id: the API would deactivate every installed licence")

		return
	}

	result, httpResp, err := client.LicenseAPI.UninstallLicense(ctx).KeyId(keyID).Execute()

	// An unknown keyId is reported as HTTP 400 with "No active license found
	// matching the specified keyId": the licence is already gone, which is what
	// a delete wants. CheckDeleted does not cover this 400, so handle it first.
	// Only that message counts: any other 400 ("Uninstall License error: ...")
	// means the licence is still installed and must be reported, not dropped
	// from state (MORPH-17987). The generated client's error text carries only
	// the HTTP status, so the message is looked for in the response body.
	if alreadyUninstalled(err, httpResp) {
		tflog.Info(ctx, "license "+keyID+" already uninstalled")

		return
	}

	if err := errfmt.CheckDeleted(result, err, httpResp); err != nil {
		resp.Diagnostics.AddError(deleteOperation, errfmt.ErrMsg(err, httpResp))

		return
	}
}

// noActiveLicenseMsg is what DELETE /api/license answers, with HTTP 400, when
// no installed licence has the given key id (ApplianceLicenseService.uninstallLicense).
const noActiveLicenseMsg = "No active license found"

// alreadyUninstalled reports whether a failed uninstall failed because the
// licence is not installed (HTTP 400 carrying noActiveLicenseMsg).
func alreadyUninstalled(err error, httpResp *http.Response) bool {
	if httpResp == nil || httpResp.StatusCode != http.StatusBadRequest {
		return false
	}

	return strings.Contains(string(provisionhint.ErrorBody(err, httpResp)), noActiveLicenseMsg)
}

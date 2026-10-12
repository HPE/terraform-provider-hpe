// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package license

import (
	"context"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

// Update makes sure the configured key is installed and records its digest.
// ModifyPlan only plans an update when the key's digest differs from the one
// recorded (or none is recorded, as after a state upgrade or an apply with
// -refresh=false), so this is where a new key is installed.
func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	// key is write-only: read it from config.
	var config LicenseModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	k := strings.TrimSpace(config.Key.ValueString())
	if k == "" {
		resp.Diagnostics.AddError(updateOperation, "key must not be empty")

		return
	}

	configured := keyDigest(k)
	keyID := licenseKeyID(k)

	// Posting a key is a non-stacking install that deactivates every other
	// installed licence, so the configured key is first looked for among the
	// installed licences -- the appliance, not the recorded digest, decides --
	// and is only posted when it is not installed. A key that is already
	// installed (the digest was merely not recorded) just has its digest
	// recorded.
	state, installed, found := getLicenseAsState(ctx, client, configured, keyID, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	if !found {
		_, httpResp, err := client.LicenseAPI.InstallLicense(ctx).
			InstallLicenseRequest(sdk.InstallLicenseRequest{License: k}).
			Execute()
		if err != nil || httpResp.StatusCode != http.StatusOK {
			resp.Diagnostics.AddError(updateOperation, errfmt.ErrMsg(err, httpResp))

			return
		}

		// The new key is installed from here on and the previous one is no
		// longer. Record that before reading the licence back: should the
		// read-back fail, state and the digest describe the installed licence
		// -- not the replaced one, whose absence would otherwise have the next
		// refresh drop the resource and the next apply post the key again --
		// and the following refresh completes the state from the appliance.
		resp.Diagnostics.Append(storeDigest(ctx, resp.Private, configured)...)
		resp.Diagnostics.Append(resp.State.Set(ctx, provisionalState(keyID))...)
		if resp.Diagnostics.HasError() {
			return
		}

		state, installed, found = getLicenseAsState(ctx, client, configured, keyID, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			tflog.Warn(ctx, "license "+keyID+" installed but not read back; state holds the key id only")

			return
		}
		if !found {
			resp.Diagnostics.AddError(updateOperation, "licence "+keyID+
				" was accepted by the appliance but is not among the installed licences")

			return
		}
	}

	// Record the digest from the API's hash (identical to configured when the
	// appliance trimmed the key as this provider does); fall back to the
	// locally computed digest if the API omitted the hash.
	hash := installed.Hash
	if hash == "" {
		hash = configured
	}
	resp.Diagnostics.Append(storeDigest(ctx, resp.Private, hash)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

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

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	// key is write-only: it is null in the plan, so read it from config.
	var config LicenseModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	k := strings.TrimSpace(config.Key.ValueString())
	if k == "" {
		resp.Diagnostics.AddError(createOperation, "key must not be empty")

		return
	}

	// Without InstallAction the POST is non-stacking: every other active licence
	// is deactivated. The SDKv2 resource behaves the same way; keep it.
	_, httpResp, err := client.LicenseAPI.InstallLicense(ctx).
		InstallLicenseRequest(sdk.InstallLicenseRequest{License: k}).
		Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(createOperation, errfmt.ErrMsg(err, httpResp))

		return
	}

	// The licence is installed from here on. Should the read-back below fail,
	// leave the key id and the key's digest behind so Terraform tracks the
	// installed licence as a tainted resource rather than losing it.
	state, installed, found := getLicenseAsState(ctx, client, keyDigest(k), licenseKeyID(k), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		taintInstalled(ctx, k, resp)

		return
	}
	if !found {
		// The appliance accepted the key but does not list it among the
		// installed licences: an appliance inconsistency.
		resp.Diagnostics.AddError(createOperation, "licence "+licenseKeyID(k)+
			" was accepted by the appliance but is not among the installed licences")
		taintInstalled(ctx, k, resp)

		return
	}

	// Record the digest from the API's hash (identical to keyDigest(k) when the
	// appliance trimmed the key as this provider does); fall back to the locally
	// computed digest if the API omitted the hash.
	hash := installed.Hash
	if hash == "" {
		hash = keyDigest(k)
	}
	resp.Diagnostics.Append(storeDigest(ctx, resp.Private, hash)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// taintInstalled records, after a successful install whose read-back failed,
// the key id in state and the key's digest in private state. Returned with
// the error, the partial state makes Terraform keep the installed licence as a
// tainted resource (the same pattern as cleanup.TaintResourceState, which
// cannot be used here because the id is the string key id): the next apply
// replaces it -- an uninstall and a re-install -- unless it is untainted, in
// which case the next refresh completes the state from the appliance.
func taintInstalled(ctx context.Context, k string, resp *resource.CreateResponse) {
	keyID := licenseKeyID(k)

	tflog.Warn(ctx, "license "+keyID+" installed but not read back; setting partial state with the key id only")

	resp.Diagnostics.Append(resp.State.Set(ctx, provisionalState(keyID))...)
	resp.Diagnostics.Append(storeDigest(ctx, resp.Private, keyDigest(k))...)

	resp.Diagnostics.AddWarning(
		"license "+keyID+" installed but not read back",
		"The licence key was installed on the appliance, but its details could not be read back, so the "+
			"resource is marked as tainted. On the next 'terraform apply', Terraform uninstalls and re-installs "+
			"the licence. To keep it installed as it is, run 'terraform untaint <address>' instead; the "+
			"following refresh completes the resource's state from the appliance.",
	)
}

// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package hostprofile

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

// getHostProfileAsState fetches the host profile and maps it onto the
// Terraform model. The returned bool is false when the profile no longer
// exists (HTTP 404), signalling the caller to remove it from state.
func getHostProfileAsState(
	ctx context.Context,
	id int64,
	client *sdk.APIClient,
	plan HostProfileModel,
) (HostProfileModel, bool, diag.Diagnostics) {
	var state HostProfileModel
	var diags diag.Diagnostics

	hpResp, hresp, err := client.HostProfilesAPI.GetHostProfile(ctx, id).Execute()
	// Check for a 404 first. The generated SDK returns a non-nil error for any
	// status >= 300, so a deleted profile arrives here as an error; checking
	// IsNotFound before the transport-error branch lets a genuine 404 return a
	// clean found=false (no diagnostics) so the caller removes it from state
	// rather than surfacing a perpetual refresh error. IsNotFound is nil-safe,
	// so a true transport failure (nil hresp) still falls through to the error
	// branch below.
	if errfmt.IsNotFound(hresp) {
		return state, false, diags
	}
	if err != nil {
		diags.AddError(
			"populate host profile resource",
			fmt.Sprintf("host profile %d GET failed: ", id)+
				errfmt.ErrMsg(err, hresp),
		)

		return state, false, diags
	}
	if hresp == nil || hresp.StatusCode != http.StatusOK {
		diags.AddError(
			"populate host profile resource",
			fmt.Sprintf("host profile %d GET failed: ", id)+
				errfmt.ErrMsg(err, hresp),
		)

		return state, false, diags
	}

	if hpResp == nil {
		diags.AddError(
			"populate host profile resource",
			fmt.Sprintf("host profile %d GET response missing body", id),
		)

		return state, false, diags
	}

	hp := hpResp.HostProfile
	if hp == nil || hp.Id == nil {
		diags.AddError(
			"populate host profile resource",
			fmt.Sprintf("host profile %d GET response missing hostProfile", id),
		)

		return state, false, diags
	}

	state.Id = types.Int64Value(*hp.Id)
	state.Code = convert.StrToType(hp.Code.Get())
	state.Name = convert.StrToType(hp.Name)
	state.Description = convert.StrToType(hp.Description.Get())
	// A description is cleared by storing an empty string (the PUT ignores an
	// explicit null); the API echoes "" back. When the configuration carries no
	// description, treat that echoed "" as unset so the state matches the null
	// configuration and a removed description does not produce a perpetual diff.
	// Only coerce when the plan is itself null, though: an explicitly configured
	// description = "" must round-trip as "" (Create and Update both send it
	// as-is), or the post-apply state would differ from config and raise a
	// "provider produced inconsistent result" error.
	if plan.Description.IsNull() && state.Description.ValueString() == "" {
		state.Description = types.StringNull()
	}
	state.Enabled = convert.BoolToType(hp.Enabled)
	state.ServerCount = convert.Int32ToType(hp.ServerCount)
	state.ComplianceStatus = convert.StrToType(hp.ComplianceStatus)

	// During import the prior state carries only the id (ImportState sets
	// nothing else), so name -- a Required attribute always present in real
	// config and state -- is null. In that case populate config from the API
	// unconditionally; otherwise apply the normal plan-null guard so an
	// unmanaged config stays null and does not produce a perpetual diff.
	importing := !plan.Id.IsNull() && plan.Name.IsNull()

	var cfg ConfigValue
	var cfgDiags diag.Diagnostics
	if importing {
		cfg, cfgDiags = flattenConfigAlways(ctx, hp.Config)
	} else {
		cfg, cfgDiags = flattenConfig(ctx, hp.Config, plan.Config)
	}
	diags.Append(cfgDiags...)
	state.Config = cfg

	versions, verDiags := flattenVersions(ctx, hp.Versions)
	diags.Append(verDiags...)
	state.Versions = versions

	if diags.HasError() {
		return state, false, diags
	}

	return state, true, diags
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var plan HostProfileModel

	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"read host profile resource",
			"failed to create client: "+err.Error(),
		)

		return
	}

	id := plan.Id.ValueInt64()

	state, found, diags := getHostProfileAsState(ctx, id, client, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !found {
		tflog.Info(ctx, fmt.Sprintf(
			"host profile %d not found, removing from state", id,
		))
		resp.State.RemoveResource(ctx)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

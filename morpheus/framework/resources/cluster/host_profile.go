// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cluster

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

// readClusterHostProfile returns the id of the host profile currently attached
// to the cluster, or a null value when none is attached. The read is
// unconditional: it always reflects the profile the appliance actually has
// attached (including one attached out of band), so that drift surfaces as a
// diff rather than being silently masked.
func readClusterHostProfile(
	ctx context.Context,
	clusterID int64,
	client *sdk.APIClient,
) (types.Int64, diag.Diagnostics) {
	var diags diag.Diagnostics

	hpResp, httpResp, err := client.HostProfilesAPI.
		GetClusterHostProfile(ctx, clusterID).Execute()
	if errfmt.IsNotFound(httpResp) {
		return types.Int64Null(), diags
	}
	if err != nil || httpResp == nil || hpResp == nil ||
		httpResp.StatusCode != http.StatusOK {
		diags.AddError(
			"read cluster host profile",
			fmt.Sprintf("cluster %d host-profile GET failed: ", clusterID)+
				errfmt.ErrMsg(err, httpResp),
		)

		return types.Int64Null(), diags
	}

	hp := hpResp.HostProfile.Get()
	if hp == nil || hp.Id == nil {
		return types.Int64Null(), diags
	}

	return convert.Int64ToType(hp.Id), diags
}

// applyClusterHostProfile reconciles the attached host profile with the planned
// host_profile_id. It attaches when a profile id is set and detaches when the
// attribute is cleared. When a *different* profile is already attached it
// detaches that one first: the attach endpoint returns 409 if it is asked to
// replace an existing profile, so switching profiles must go detach-then-attach.
func applyClusterHostProfile(
	ctx context.Context,
	clusterID int64,
	client *sdk.APIClient,
	plan types.Int64,
) diag.Diagnostics {
	var diags diag.Diagnostics

	if plan.IsUnknown() {
		return diags
	}

	if plan.IsNull() {
		return detachClusterHostProfile(ctx, clusterID, client)
	}

	// Determine what is currently attached. If it already matches the plan the
	// attach is a no-op; if a different profile is attached, detach it first to
	// avoid the attach endpoint's 409-on-replace behaviour.
	current, readDiags := readClusterHostProfile(ctx, clusterID, client)
	diags.Append(readDiags...)
	if diags.HasError() {
		return diags
	}

	if current.Equal(plan) {
		return diags
	}

	if !current.IsNull() {
		if diags.Append(
			detachClusterHostProfile(ctx, clusterID, client)...,
		); diags.HasError() {
			return diags
		}
	}

	_, httpResp, err := client.HostProfilesAPI.
		AttachClusterHostProfile(ctx, clusterID).
		AttachClusterHostProfileRequest(sdk.AttachClusterHostProfileRequest{
			ProfileId: plan.ValueInt64(),
		}).
		Execute()
	if err != nil || httpResp == nil || httpResp.StatusCode != http.StatusOK {
		diags.AddError(
			"attach cluster host profile",
			fmt.Sprintf("cluster %d host-profile attach failed: ", clusterID)+
				errfmt.ErrMsg(err, httpResp),
		)
	}

	return diags
}

// detachClusterHostProfile detaches any host profile from the cluster. It is
// idempotent: the server clears the profile and returns success even when none
// is attached, and a missing cluster is treated as already detached.
func detachClusterHostProfile(
	ctx context.Context,
	clusterID int64,
	client *sdk.APIClient,
) diag.Diagnostics {
	var diags diag.Diagnostics

	_, httpResp, err := client.HostProfilesAPI.
		DetachClusterHostProfile(ctx, clusterID).Execute()
	if errfmt.IsNotFound(httpResp) {
		return diags
	}
	if err != nil || httpResp == nil || httpResp.StatusCode != http.StatusOK {
		diags.AddError(
			"detach cluster host profile",
			fmt.Sprintf("cluster %d host-profile detach failed: ", clusterID)+
				errfmt.ErrMsg(err, httpResp),
		)
	}

	return diags
}

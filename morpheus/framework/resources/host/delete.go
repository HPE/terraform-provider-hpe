// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package host

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

// Delete implements resource.Resource.
func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state HostModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteTimeout, diags := state.Timeouts.Delete(ctx, defaultDeleteTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("failed to create API client", err.Error())

		return
	}

	id := state.Id.ValueInt64()

	removeReq := client.HostsAPI.RemoveHost(ctx, sdk.UpdateHostIdParameter{Int64: &id})
	removeReq = applyDeleteOptions(removeReq, state.DeleteOptions)

	result, hresp, err := removeReq.Execute()
	if err := errfmt.CheckDeleted(result, err, hresp); err != nil { // a 404 is success
		resp.Diagnostics.AddError(
			"failed to delete host",
			fmt.Sprintf("host %d: %s", id, err.Error()),
		)

		return
	}

	// Poll until the host is gone (GET returns 404).
	if err := r.pollDeleted(ctx, id, deleteTimeout); err != nil {
		resp.Diagnostics.AddError(
			"host deletion did not complete",
			fmt.Sprintf("host %d: %s", id, err.Error()),
		)

		return
	}

	tflog.Info(ctx, fmt.Sprintf("host %d deleted", id))
}

// deleteQueryFlags returns the DELETE query flags to send for the given
// delete_options. Only set flags appear; an unset flag is omitted so Morpheus
// applies its own type-aware default (removeResources = computeServerType.managed).
func deleteQueryFlags(opts DeleteOptionsValue) map[string]string {
	flags := map[string]string{}

	if opts.IsNull() || opts.IsUnknown() {
		return flags
	}

	if !opts.RemoveResources.IsNull() && !opts.RemoveResources.IsUnknown() {
		flags["removeResources"] = strconv.FormatBool(opts.RemoveResources.ValueBool())
	}

	if !opts.PreserveVolumes.IsNull() && !opts.PreserveVolumes.IsUnknown() {
		flags["preserveVolumes"] = strconv.FormatBool(opts.PreserveVolumes.ValueBool())
	}

	if !opts.Force.IsNull() && !opts.Force.IsUnknown() {
		flags["force"] = strconv.FormatBool(opts.Force.ValueBool())
	}

	if !opts.RemoveInstances.IsNull() && !opts.RemoveInstances.IsUnknown() {
		flags["removeInstances"] = strconv.FormatBool(opts.RemoveInstances.ValueBool())
	}

	return flags
}

// applyDeleteOptions sets the DELETE query flags only when the corresponding
// delete_options field is set.
func applyDeleteOptions(
	req sdk.ApiRemoveHostRequest,
	opts DeleteOptionsValue,
) sdk.ApiRemoveHostRequest {
	flags := deleteQueryFlags(opts)

	if v, ok := flags["removeResources"]; ok {
		req = req.RemoveResources(v)
	}
	if v, ok := flags["preserveVolumes"]; ok {
		req = req.PreserveVolumes(v)
	}
	if v, ok := flags["force"]; ok {
		req = req.Force(v)
	}
	if v, ok := flags["removeInstances"]; ok {
		req = req.RemoveInstances(v)
	}

	return req
}

// pollDeleted polls GetHost until it returns 404.
func (r *Resource) pollDeleted(ctx context.Context, id int64, timeout time.Duration) error {
	client, err := r.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	// The shared rule: 404 ends the wait, a connection failure or a 5xx is
	// retried (capped at errfmt.MaxConsecutiveTransient in a row), anything
	// else is final, and a 200 means the host is still there.
	poller := &errfmt.DeletePoller{What: fmt.Sprintf("host %d", id)}
	check := func() (struct{}, error) {
		_, hresp, err := client.HostsAPI.
			GetHost(ctx, sdk.GetHostIdParameter{Int64: &id}).
			Execute()

		return struct{}{}, poller.Poll(ctx, err, hresp)
	}

	_, err = backoff.Retry(
		ctx,
		check,
		backoff.WithBackOff(backoff.NewConstantBackOff(5*time.Second)),
		backoff.WithMaxElapsedTime(timeout),
	)

	return err
}

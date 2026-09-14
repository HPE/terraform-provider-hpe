// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package tenant

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

const deleteOperation = "delete tenant resource"

const (
	// deleteTimeout bounds how long we wait for the asynchronous tenant
	// deletion job to complete. Tenant delete publishes a background
	// "deleteAccount" job and returns 200 immediately; both name and subdomain
	// are unique-constrained, so we must wait for the tenant to actually
	// disappear before returning, otherwise an immediate recreate would race
	// the still-running deletion. The wait ends early on any error other than
	// the 404 that signals completion.
	deleteTimeout = 5 * time.Minute
	// deletePollInterval is how often we poll GetTenant while waiting for 404.
	deletePollInterval = 5 * time.Second
)

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state TenantModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The master tenant cannot be deleted; the API rejects it with an opaque
	// 400. Surface a clear, named diagnostic instead.
	if state.Master.ValueBool() {
		resp.Diagnostics.AddError(deleteOperation, masterTenantDeleteMsg)

		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	id := state.Id.ValueInt64()
	// remove_resources is null when omitted, which ValueBool() reports as false:
	// remove the tenant from Morpheus but leave its provisioned VMs running.
	removeResources := state.RemoveResources.ValueBool()

	tflog.Debug(ctx, fmt.Sprintf(
		"Deleting tenant %d (remove_resources=%t)", id, removeResources))

	_, hresp, err := client.TenantsAPI.RemoveTenant(ctx, id).
		RemoveResources(removeResources).Execute()
	if errfmt.IsNotFound(hresp) {
		return
	}
	if err := errfmt.CheckResponse(err, hresp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "tenant", "", err, hresp)

		return
	}

	// Delete is asynchronous: poll until the tenant returns 404. A 200 means
	// the background job is still running and is retried until deleteTimeout;
	// any other failure (auth, network, a 500) is permanent, so the real error
	// is reported at once instead of after the full timeout.
	waitForDeleted := func() (struct{}, error) {
		_, ghresp, err := client.TenantsAPI.GetTenant(ctx, id).Execute()
		if errfmt.IsNotFound(ghresp) {
			return struct{}{}, nil
		}
		if err := errfmt.CheckResponse(err, ghresp); err != nil {
			return struct{}{}, backoff.Permanent(
				fmt.Errorf("tenant %d GET failed: %w", id, err))
		}

		return struct{}{}, fmt.Errorf(
			"tenant %d still exists; the background deletion job may still be running on the appliance",
			id)
	}

	if _, err := backoff.Retry(
		ctx,
		waitForDeleted,
		backoff.WithBackOff(backoff.NewConstantBackOff(deletePollInterval)),
		backoff.WithMaxElapsedTime(deleteTimeout),
	); err != nil {
		resp.Diagnostics.AddError(
			deleteOperation,
			fmt.Sprintf("waiting (up to %s) for tenant %d to be deleted: %s",
				deleteTimeout, id, err),
		)

		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Tenant %d deletion confirmed (404)", id))
}

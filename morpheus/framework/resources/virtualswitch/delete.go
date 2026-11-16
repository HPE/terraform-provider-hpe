// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package virtualswitch

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/process"
)

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var data VirtualSwitchModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	id := data.Id.ValueInt64()

	timeout, d := data.Timeouts.Delete(ctx, defaultOpTimeout)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// One budget for the whole delete: the DELETE and its retries, the
	// process wait and the poll until the switch is gone.
	deadline := time.Now().Add(timeout)

	// Morpheus refuses the delete while a sync is in flight (a 400 whose body
	// carries errors.sync, VirtualSwitchService.isSyncInFlight -- transient:
	// retried once the sync ends) or while the switch carries networks (a
	// 400 with errors.networks -- terminal, surfaced at once, MORPH-17987).
	// A refusal rendered as 200 with success:false is caught by inspecting
	// the body (errfmt.CheckDeleted) and is final too.
	var processID int64

	del := func() error {
		result, httpResp, err := client.VirtualSwitchesAPI.
			DeleteVirtualSwitch(ctx, id).Execute()
		if err := syncRefusal(err, httpResp); err != nil {
			return err
		}

		if err := errfmt.CheckDeleted(result, err, httpResp); err != nil {
			return err
		}

		if result != nil && result.ProcessId != nil {
			processID = *result.ProcessId
		}

		return nil
	}

	// Between refused attempts wait for the sync to end; a switch that is
	// gone by then has settled as far as a delete is concerned (the retried
	// DELETE gets its own 404, which is a success).
	settle := func(left time.Duration) error {
		if err := pollUntilNotSyncing(ctx, client, id, left); err != nil && !errors.Is(err, errSwitchGone) {
			return err
		}

		return nil
	}

	if err := retryRefused(ctx, fmt.Sprintf("virtual switch %d", id), deadline,
		refusalAttempts, refusalPause, settle, del); err != nil {
		resp.Diagnostics.AddError(
			"error deleting virtual switch",
			fmt.Sprintf("virtual switch %d DELETE failed: ", id)+err.Error(),
		)

		return
	}

	// A budget the retries used up skips the process wait: the DELETE was
	// accepted, and the confirmation GET below floors its own budget.
	if left := time.Until(deadline); processID != 0 && left > 0 {
		if _, err := process.Wait(ctx, client, processID, left); err != nil {
			resp.Diagnostics.AddError(
				"error deleting virtual switch",
				fmt.Sprintf("virtual switch %d: delete process failed: ", id)+err.Error(),
			)

			return
		}
	}

	// Confirm the switch is gone: the shared rule (errfmt.DeletePoller) ends
	// the wait on a 404, retries a connection failure or a 5xx from a busy
	// appliance (capped in a row), and treats any other answer as final.
	poller := &errfmt.DeletePoller{What: fmt.Sprintf("virtual switch %d", id)}
	waitForDeleted := func() (struct{}, error) {
		_, httpResp, err := client.VirtualSwitchesAPI.GetVirtualSwitch(ctx, id).Execute()

		return struct{}{}, poller.Poll(ctx, err, httpResp)
	}

	left := time.Until(deadline)
	if left <= 0 {
		left = time.Second // the process used the budget; one check of the GET
	}

	if _, err := backoff.Retry(
		ctx,
		waitForDeleted,
		backoff.WithBackOff(backoff.NewConstantBackOff(5*time.Second)),
		backoff.WithMaxElapsedTime(left),
	); err != nil {
		resp.Diagnostics.AddError(
			"error deleting virtual switch",
			fmt.Sprintf("virtual switch %d: DELETE confirmation failed: ", id)+err.Error(),
		)
	}
}

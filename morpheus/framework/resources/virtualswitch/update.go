// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package virtualswitch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/process"
)

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state VirtualSwitchModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	timeout, d := plan.Timeouts.Update(ctx, defaultOpTimeout)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueInt64()

	// Uplink removal is impossible in the API. ModifyPlan refuses it when the
	// uplinks are known at plan time; a set that was only known at apply (a
	// deferred data source) is judged here.
	removed, d := uplinkRemovals(ctx, state, plan)
	if resp.Diagnostics.Append(d...); resp.Diagnostics.HasError() {
		return
	}

	if len(removed) > 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("uplinks"),
			"Cannot remove a host uplink",
			fmt.Sprintf("Morpheus cannot remove a host from a virtual switch (server_id %v). "+
				"Remove the host from the cluster, or recreate the switch, instead of dropping its "+
				"uplink from the configuration.", removed),
		)

		return
	}

	up, d := computeUpdatePlan(ctx, state, plan)
	if resp.Diagnostics.Append(d...); resp.Diagnostics.HasError() {
		return
	}

	if up.empty() {
		// Nothing changed: still refresh computed values.
		newState, diags := readSwitchState(ctx, client, id, plan)
		if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
			return
		}

		resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)

		return
	}

	// Each step is a separate async process; the sync guard forbids overlapping
	// mutations, so each is followed by Wait + poll-until-not-syncing. On any
	// failure, persist what actually happened with a Read, then return the
	// error naming the step.
	// One budget for the whole update: every step's waits share it.
	if err := runUpdateSteps(ctx, client, id, up, time.Now().Add(timeout)); err != nil {
		newState, diags := readSwitchState(ctx, client, id, plan)
		if !diags.HasError() {
			resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
		}

		resp.Diagnostics.AddError("error updating virtual switch", err.Error())

		return
	}

	newState, diags := readSwitchState(ctx, client, id, plan)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func runUpdateSteps(
	ctx context.Context,
	client *sdk.APIClient,
	id int64,
	up updatePlan,
	deadline time.Time,
) error {
	// Morpheus refuses a mutation while a sync is in flight (a 400 whose body
	// carries errors.sync -- VirtualSwitchService.isSyncInFlight); an apply
	// that starts during one would fail on its first step, so wait for the
	// switch to settle first, within the budget. A sync that starts after
	// this wait is handled by retryRefused on each step.
	left := time.Until(deadline)
	if left <= 0 {
		return fmt.Errorf("the update timeout elapsed before the first step")
	}

	if err := pollUntilNotSyncing(ctx, client, id, left); err != nil {
		return fmt.Errorf("before updating: %s", err.Error())
	}

	// Between refused attempts, wait for the sync that refused us to finish.
	settle := func(left time.Duration) error { return pollUntilNotSyncing(ctx, client, id, left) }

	// Step 1: rename.
	if up.renameTo != nil {
		var resp *sdk.UpdateVirtualSwitch200Response

		what := fmt.Sprintf("virtual switch %d rename", id)
		callErr := retryRefused(ctx, what, deadline, refusalAttempts, refusalPause, settle, func() error {
			r, httpResp, err := client.VirtualSwitchesAPI.
				RenameVirtualSwitch(ctx, id).
				RenameVirtualSwitchRequest(buildRename(*up.renameTo)).
				Execute()
			resp = r

			return mutationErr(r, err, httpResp)
		})
		if err := awaitStep(ctx, client, id, "rename", procID(resp), callErr, deadline); err != nil {
			return err
		}
	}

	// Step 2: delete removed segments. The sync guard is a 400 (syncRefusal)
	// and is retried within the budget like the other steps; a refused delete
	// rendered as HTTP 200 with success:false is caught by inspecting the body
	// via errfmt.CheckDeleted (MORPH-17987) and is final.
	for _, segID := range up.removedSegIDs {
		step := fmt.Sprintf("delete segment %d", segID)

		var processID int64

		del := func() error {
			resp, httpResp, err := client.VirtualSwitchesAPI.
				DeleteVirtualSwitchSegment(ctx, id, segID).Execute()
			if err := syncRefusal(err, httpResp); err != nil {
				return err
			}

			if err := errfmt.CheckDeleted(resp, err, httpResp); err != nil {
				return err
			}

			processID = procID(resp)

			return nil
		}

		callErr := retryRefused(ctx, fmt.Sprintf("virtual switch %d segment %d", id, segID),
			deadline, refusalAttempts, refusalPause, settle, del)
		if err := awaitStep(ctx, client, id, step, processID, callErr, deadline); err != nil {
			return err
		}
	}

	// Step 3: composite edit (switchUpdate + newSegments + editedSegments).
	if up.hasCompositeWork() {
		var resp *sdk.UpdateVirtualSwitch200Response

		what := fmt.Sprintf("virtual switch %d edit", id)
		callErr := retryRefused(ctx, what, deadline, refusalAttempts, refusalPause, settle, func() error {
			r, httpResp, err := client.VirtualSwitchesAPI.
				CompositeEditVirtualSwitch(ctx, id).
				CompositeEditVirtualSwitchRequest(buildCompositeEdit(up)).
				Execute()
			resp = r

			return mutationErr(r, err, httpResp)
		})
		if err := awaitStep(ctx, client, id, "edit", procID(resp), callErr, deadline); err != nil {
			return err
		}
	}

	// Step 4: description-only PUT (compositeEdit ignores description).
	if up.descriptionChanged {
		var resp *sdk.UpdateVirtualSwitch200Response

		what := fmt.Sprintf("virtual switch %d description", id)
		callErr := retryRefused(ctx, what, deadline, refusalAttempts, refusalPause, settle, func() error {
			r, httpResp, err := client.VirtualSwitchesAPI.
				UpdateVirtualSwitch(ctx, id).
				UpdateVirtualSwitchRequest(buildDescriptionUpdate(up.description)).
				Execute()
			resp = r

			return mutationErr(r, err, httpResp)
		})
		if err := awaitStep(ctx, client, id, "update description", procID(resp), callErr, deadline); err != nil {
			return err
		}
	}

	return nil
}

// respErr renders a mutation call's HTTP error into a single error, or nil.
func respErr(err error, httpResp *http.Response) error {
	if err != nil || httpResp == nil || httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s", errfmt.ErrMsg(err, httpResp))
	}

	return nil
}

// refusedError is a mutation Morpheus refused: a 400 whose body names the
// sync guard (syncing: a sync is in flight, and the request may succeed once
// it ends -- retryRefused tries again), or a 200 with success:false (final).
type refusedError struct {
	msg     string
	syncing bool
}

func (e refusedError) Error() string { return e.msg }

// syncGuardKey is the per-field error key Morpheus sets when a mutation is
// refused because a sync is in flight (errors.sync =
// gomorpheus.virtualSwitch.error.syncInFlight, VirtualSwitchService).
const syncGuardKey = "sync"

// syncRefusal recognises the sync guard: a 400 whose JSON body has
// errors.sync. It returns a syncing refusedError for one, nil for anything
// else (including other errors, which the caller classifies).
func syncRefusal(err error, httpResp *http.Response) error {
	if err == nil || httpResp == nil || httpResp.StatusCode != http.StatusBadRequest {
		return nil
	}

	// The SDK's GenericOpenAPIError carries the raw body; matched by shape
	// so a test can supply one.
	var apiErr interface{ Body() []byte }
	if !errors.As(err, &apiErr) {
		return nil
	}

	var body map[string]any
	if json.Unmarshal(apiErr.Body(), &body) != nil {
		return nil
	}

	fieldErrs, _ := body["errors"].(map[string]any)
	if _, ok := fieldErrs[syncGuardKey]; !ok {
		return nil
	}

	return refusedError{msg: errfmt.RefusalMessage(body), syncing: true}
}

// refusal returns the error a 200 with success:false carries (Morpheus renders
// some refusals that way, with the reason in msg and per-field errors), or
// nil. Such a refusal is final.
func refusal(success *bool, extra map[string]interface{}) error {
	if success == nil || *success {
		return nil
	}

	return refusedError{msg: errfmt.RefusalMessage(extra)}
}

// refusalAttempts bounds the retries of a mutation refused by the sync guard;
// refusalPause is the least time between attempts (a variable so the
// fake-appliance test can shorten it). Between attempts the switch is polled
// until it leaves syncing, within the operation's deadline, so a long sync is
// waited out rather than failing the step after a fixed number of pauses; a
// settle that returns sooner than the pause is topped up to it.
const refusalAttempts = 5

var refusalPause = 10 * time.Second

// retryRefused runs the mutation and, while the sync guard refuses it, waits
// for the switch to stop syncing (settle; nil falls back to a fixed pause)
// and tries again, up to attempts times and while the deadline allows. Any
// other error -- including a refusal that is not the sync guard -- returns
// at once. Exhaustion returns the last refusal, saying how often it was tried.
func retryRefused(
	ctx context.Context,
	what string,
	deadline time.Time,
	attempts int,
	pause time.Duration,
	settle func(left time.Duration) error,
	call func() error,
) error {
	for attempt := 1; ; attempt++ {
		err := call()

		var refused refusedError
		if err == nil || !errors.As(err, &refused) || !refused.syncing {
			return err
		}

		if attempt >= attempts {
			return fmt.Errorf("%w (still refused after %d attempts)", err, attempt)
		}

		left := time.Until(deadline)
		if left <= pause {
			return fmt.Errorf("%w (the operation's timeout leaves no time to retry)", err)
		}

		tflog.Warn(ctx, what+": Morpheus refused the request because a sync is in flight; retrying",
			map[string]any{"attempt": attempt, "of": attempts, "error": err.Error()})

		// Wait for the sync that refused us to end; when the switch does not
		// report one (the guard and the status field are not known to be in
		// step), the settle returns at once and the pause is slept instead,
		// so the attempts are never spent back to back.
		waited := time.Duration(0)

		if settle != nil {
			started := time.Now()
			if serr := settle(left); serr != nil {
				return fmt.Errorf("%w (while waiting for the sync to finish: %w)", err, serr)
			}

			waited = time.Since(started)
		}

		if waited >= pause {
			continue
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("%w (while waiting to retry: %w)", err, ctx.Err())
		case <-time.After(pause - waited):
		}
	}
}

// mutationErr combines the transport/status error with a success:false body.
func mutationErr(resp *sdk.UpdateVirtualSwitch200Response, err error, httpResp *http.Response) error {
	if e := syncRefusal(err, httpResp); e != nil {
		return e
	}

	if e := respErr(err, httpResp); e != nil {
		return e
	}

	if resp == nil {
		return nil
	}

	return refusal(resp.Success, resp.AdditionalProperties)
}

// awaitStep waits on a mutation's process then polls until the switch is not
// syncing. callErr is the error from the mutation call itself (nil on success).
func awaitStep(
	ctx context.Context,
	client *sdk.APIClient,
	id int64,
	step string,
	processID int64,
	callErr error,
	deadline time.Time,
) error {
	if callErr != nil {
		return fmt.Errorf("%s failed: %s", step, callErr.Error())
	}

	// What is left of the update's budget; a spent one is a timeout, not an
	// unbounded wait (backoff enforces only a positive MaxElapsedTime).
	remaining := func() (time.Duration, error) {
		left := time.Until(deadline)
		if left <= 0 {
			return 0, fmt.Errorf("%s: the update timeout elapsed", step)
		}

		return left, nil
	}

	if processID != 0 {
		left, err := remaining()
		if err != nil {
			return err
		}

		if _, err := process.Wait(ctx, client, processID, left); err != nil {
			return fmt.Errorf("%s process failed: %s", step, err.Error())
		}
	}

	left, err := remaining()
	if err != nil {
		return err
	}

	if err := pollUntilNotSyncing(ctx, client, id, left); err != nil {
		return fmt.Errorf("%s: %s", step, err.Error())
	}

	return nil
}

// procID reads the processId from an async envelope, tolerating a nil response.
func procID(resp *sdk.UpdateVirtualSwitch200Response) int64 {
	if resp == nil || resp.ProcessId == nil {
		return 0
	}

	return *resp.ProcessId
}

// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package virtualswitch

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/HPE/terraform-provider-hpe/utils/cleanup"

	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/constants"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/process"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/versioncheck"
)

const gatedFeature = "HVM virtual switches"

// defaultOpTimeout is used when a timeouts block omits a value; the host
// fan-out can take minutes.
const defaultOpTimeout = 30 * time.Minute

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan VirtualSwitchModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	resp.Diagnostics.Append(versioncheck.Require(
		ctx, client, gatedFeature, constants.VirtualSwitchMinVersion,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, d := plan.Timeouts.Create(ctx, defaultOpTimeout)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	segs, d := segmentInputs(ctx, plan)
	if resp.Diagnostics.Append(d...); resp.Diagnostics.HasError() {
		return
	}

	ups, d := uplinkInputs(ctx, plan)
	if resp.Diagnostics.Append(d...); resp.Diagnostics.HasError() {
		return
	}

	clusterID := plan.ClusterId.ValueInt64()
	name := plan.Name.ValueString()

	body := buildCreateRequest(plan, segs, ups)

	createResp, httpResp, err := client.VirtualSwitchesAPI.
		AddVirtualSwitch(ctx).
		ClusterId(clusterID).
		AddVirtualSwitchRequest(body).
		Execute()
	if err != nil || httpResp == nil || httpResp.StatusCode != http.StatusOK {
		detail := errfmt.ErrMsg(err, httpResp)
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			detail = fmt.Sprintf(
				"cluster %d was not found, or it is not on a virtual-switch-capable "+
					"(2.0+) layout: %s", clusterID, detail)
		}

		resp.Diagnostics.AddError(
			"error creating virtual switch",
			"virtual switch "+name+" POST failed: "+detail,
		)

		return
	}

	// A 200 with success:false is a refusal, with the reason in msg.
	if createResp != nil {
		if err := refusal(createResp.Success, createResp.AdditionalProperties); err != nil {
			resp.Diagnostics.AddError(
				"error creating virtual switch",
				"virtual switch "+name+" was refused: "+err.Error(),
			)

			return
		}
	}

	if createResp == nil || createResp.ProcessId == nil {
		// The create was accepted without a process to wait on; the switch
		// may well exist, so find it before giving up rather than orphan it.
		id, found, lookupDiags := findSwitchIDByNameRetry(ctx, client, clusterID, name, 30*time.Second, true)
		resp.Diagnostics.Append(lookupDiags...)

		if found {
			resp.Diagnostics.AddError(
				"error creating virtual switch",
				"virtual switch "+name+": create response carried no processId; the switch exists and is recorded as tainted",
			)
			cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
				ResourceType: "virtual_switch",
				ResourceID:   id,
				StateWriter:  &resp.State,
				Diagnostics:  &resp.Diagnostics,
			})

			return
		}

		resp.Diagnostics.AddError(
			"error creating virtual switch",
			"virtual switch "+name+": create response carried no processId",
		)

		return
	}

	resp.Diagnostics.Append(process.Warnings("virtual switch", name, createResp.Warnings)...)

	// One budget for everything after the POST: the process wait and the
	// poll until the switch stops syncing share it.
	deadline := time.Now().Add(timeout)

	_, procErr := process.Wait(ctx, client, *createResp.ProcessId, timeout)

	// Whether the process succeeded or failed, the switch may exist: find it by
	// name so that a failed apply still records the id and destroy can clean up.
	// The lookup is retried for a bounded time, since a single failed or
	// not-yet-visible listing would otherwise leave an existing switch
	// untracked and the next apply colliding on the name.
	lookupWait := min(lookupRetry, max(time.Until(deadline), 10*time.Second))
	id, found, findDiags := findSwitchIDByNameRetry(ctx, client, clusterID, name, lookupWait, procErr == nil)
	resp.Diagnostics.Append(findDiags...)

	// taint records the found id with every other attribute null, so a switch
	// that exists on the appliance is never orphaned by a failure after the
	// create: the next apply replaces it instead of colliding on the name.
	taint := func() {
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "virtual_switch",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})
	}

	if procErr != nil {
		if found {
			state, diags := readSwitchState(ctx, client, id, plan)
			resp.Diagnostics.Append(diags...)
			if diags.HasError() {
				taint()
			} else {
				resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			}
		}

		resp.Diagnostics.AddError(
			"error creating virtual switch",
			"virtual switch "+name+": create process failed: "+procErr.Error(),
		)

		return
	}

	if !found {
		if !resp.Diagnostics.HasError() {
			resp.Diagnostics.AddError(
				"error creating virtual switch",
				"virtual switch "+name+": not found by name after create",
			)
		}

		return
	}

	if err := pollUntilNotSyncing(ctx, client, id, time.Until(deadline)); err != nil {
		resp.Diagnostics.AddError(
			"error creating virtual switch",
			"virtual switch "+name+": "+err.Error(),
		)
		taint()

		return
	}

	state, diags := readSwitchState(ctx, client, id, plan)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		taint()

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// findSwitchIDByName lists the switches on a cluster and returns the id of the
// one with the given name (unique per cluster).
func findSwitchIDByName(
	ctx context.Context,
	client *sdk.APIClient,
	clusterID int64,
	name string,
) (int64, bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	listResp, httpResp, err := client.VirtualSwitchesAPI.
		ListVirtualSwitches(ctx).ClusterId(clusterID).Execute()
	if err != nil || httpResp == nil || httpResp.StatusCode != http.StatusOK {
		detail := errfmt.ErrMsg(err, httpResp)
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			detail = fmt.Sprintf(
				"cluster %d was not found, or it is not on a virtual-switch-capable "+
					"(2.0+) layout: %s", clusterID, detail)
		}

		diags.AddError("error listing virtual switches", detail)

		return 0, false, diags
	}

	if listResp == nil {
		return 0, false, diags
	}

	for _, sw := range listResp.VirtualSwitches {
		if sw.Id != nil && sw.Name != nil && *sw.Name == name {
			return *sw.Id, true, diags
		}
	}

	return 0, false, diags
}

// pollUntilNotSyncing polls GET until the switch status is not "syncing".
func pollUntilNotSyncing(
	ctx context.Context,
	client *sdk.APIClient,
	id int64,
	timeout time.Duration,
) error {
	transientRun := 0

	poll := func() (struct{}, error) {
		vsResp, httpResp, err := client.VirtualSwitchesAPI.GetVirtualSwitch(ctx, id).Execute()
		if err != nil || httpResp == nil || httpResp.StatusCode != http.StatusOK {
			if errfmt.IsNotFound(httpResp) {
				// Gone: nothing left to sync. Callers that expect the switch
				// to exist report it; a delete treats it as settled.
				return struct{}{}, backoff.Permanent(fmt.Errorf("virtual switch %d: %w", id, errSwitchGone))
			}

			// A connection failure or a 5xx is transient (errfmt.IsTransient),
			// up to errfmt.MaxConsecutiveTransient in a row as the shared
			// delete poller allows; a 4xx, a refused re-login or a done
			// context is final.
			if errfmt.IsTransient(ctx, err, httpResp) {
				transientRun++
				if transientRun >= errfmt.MaxConsecutiveTransient {
					return struct{}{}, backoff.Permanent(fmt.Errorf(
						"GET during sync poll failed %d times in a row: %s", transientRun, errfmt.ErrMsg(err, httpResp)))
				}

				return struct{}{}, fmt.Errorf("GET during sync poll failed: %s", errfmt.ErrMsg(err, httpResp))
			}

			return struct{}{}, backoff.Permanent(
				fmt.Errorf("GET during sync poll failed: %s", errfmt.ErrMsg(err, httpResp)))
		}

		transientRun = 0

		if vsResp == nil || vsResp.VirtualSwitch == nil || vsResp.VirtualSwitch.Status == nil {
			return struct{}{}, nil
		}

		if *vsResp.VirtualSwitch.Status == statusSyncing {
			return struct{}{}, fmt.Errorf("virtual switch %d: %w", id, errStillSyncing)
		}

		return struct{}{}, nil
	}

	// backoff only enforces a positive MaxElapsedTime; a budget already spent
	// must not become an unbounded poll.
	if timeout <= 0 {
		return fmt.Errorf("no time left in the operation's timeout to wait for virtual switch %d to stop syncing", id)
	}

	_, err := backoff.Retry(
		ctx,
		poll,
		backoff.WithBackOff(backoff.NewConstantBackOff(syncPollInterval)),
		backoff.WithMaxElapsedTime(timeout),
	)
	// backoff gives up with the last poll's error once the budget is spent;
	// a switch that stayed syncing for the whole budget is a timeout.
	if errors.Is(err, errStillSyncing) {
		return fmt.Errorf("timed out after %s waiting for virtual switch %d to stop syncing", timeout.Round(time.Second), id)
	}

	return err
}

// errSwitchGone is the sync poll's final error for a switch that answered 404.
var errSwitchGone = errors.New("not found")

// errStillSyncing is the sync poll's retryable "not yet" error; it is also the
// error backoff returns when the switch was still syncing as the budget ran
// out, which pollUntilNotSyncing reports as a timeout.
var errStillSyncing = errors.New("still syncing")

// syncPollInterval is the pause between GETs while waiting for a switch to
// leave syncing. A variable so the fake-appliance test can shorten it.
var syncPollInterval = 3 * time.Second

const statusSyncing = "syncing"

// lookupRetry bounds the post-create name lookup's retries.
const lookupRetry = 2 * time.Minute

// findSwitchIDByNameRetry is findSwitchIDByName retried, every few seconds for
// up to wait, while the listing fails -- or, when the switch must exist
// (expectExists: the create process succeeded), while it does not show it yet.
// After a failed process a clean listing with no match is final.
func findSwitchIDByNameRetry(
	ctx context.Context,
	client *sdk.APIClient,
	clusterID int64,
	name string,
	wait time.Duration,
	expectExists bool,
) (int64, bool, diag.Diagnostics) {
	deadline := time.Now().Add(wait)

	for {
		id, found, diags := findSwitchIDByName(ctx, client, clusterID, name)
		if found || time.Now().After(deadline) || (!expectExists && !diags.HasError()) {
			return id, found, diags
		}

		select {
		case <-ctx.Done():
			return id, found, diags
		case <-time.After(5 * time.Second):
		}
	}
}

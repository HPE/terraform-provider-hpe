// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package host

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
)

// The ComputeServer statuses Morpheus writes (morpheus-core services), by what
// they mean to a create. An unrecognised status keeps the poller waiting.

// provisioningStatuses mean the host is still being set up.
var provisioningStatuses = []string{"provisioning", "postProvisioning", "pending", "queued", "resizing", ""}

// readyStatuses are terminal: the host exists and is usable. A bare-metal or
// unmanaged node that has been registered but has nothing provisioned on it
// is "available"; a node held back for a later provision is "reserved".
var readyStatuses = []string{"provisioned", "running", "ok", "available", "reserved", "maintenance", "warning"}

// failedStatuses are terminal failures, including a host already on its way out.
var failedStatuses = []string{"failed", "error", "teardown", "removing", "pendingRemoval", "denied", "cancelled"}

// statusClass classifies a host provisioning status.
type statusClass int

const (
	statusProvisioning statusClass = iota
	statusReady
	statusFailed
	statusUnknown
)

// classifyStatus maps a raw host status string to a statusClass used by the
// create poller's decision.
func classifyStatus(status string) statusClass {
	switch {
	case slices.Contains(failedStatuses, status):
		return statusFailed
	case slices.Contains(readyStatuses, status):
		return statusReady
	case slices.Contains(provisioningStatuses, status):
		return statusProvisioning
	default:
		return statusUnknown
	}
}

// Create implements resource.Resource.
func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan HostModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config HostModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createTimeout, diags := plan.Timeouts.Create(ctx, defaultCreateTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("failed to create API client", err.Error())

		return
	}

	resp.Diagnostics.Append(checkCreatable(ctx, client, plan.ComputeServerType.ValueString())...)
	if resp.Diagnostics.HasError() {
		return
	}

	sshPwd, bmaasPwd, omegaPwd := writeOnlyFromConfig(config)

	createReq, bDiags := buildCreateRequest(ctx, plan, sshPwd, bmaasPwd, omegaPwd)
	resp.Diagnostics.Append(bDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createResp, hresp, err := client.HostsAPI.
		AddBaremetalHost(ctx).
		AddBaremetalHostRequest(createReq).
		Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to create host",
			errfmt.ErrMsg(err, hresp),
		)

		return
	}

	if createResp == nil || createResp.Server == nil || createResp.Server.Id == nil {
		resp.Diagnostics.AddError(
			"failed to create host",
			"create response, server, or server ID is nil",
		)

		return
	}

	id := *createResp.Server.Id

	// Record the digests of the configured write-only values now the values
	// have been sent; the appliance-side baseline follows the read-back.
	resp.Diagnostics.Append(wo.Store(ctx, resp.Private, req.Config)...)

	// The host exists from here on; a failure below taints the resource (the
	// taint writes the ID to state) rather than leaving an orphan.
	finalStatus, pollErr := r.pollProvisioned(ctx, id, createTimeout)
	if pollErr != nil {
		resp.Diagnostics.AddError(
			"host provisioning failed",
			fmt.Sprintf("host %d: %s (last status: %s)", id, pollErr.Error(), finalStatus),
		)
		taint(ctx, resp, plan, id)

		return
	}

	// The add-host handler does not bind enabled; a configured value is
	// applied with the update handler, which does.
	if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() {
		enabled := plan.Enabled.ValueBool()
		_, hresp, err := client.HostsAPI.
			UpdateHost(ctx, sdk.UpdateHostIdParameter{Int64: &id}).
			UpdateHostRequest(sdk.UpdateHostRequest{
				Server: &sdk.UpdateHostRequestServer{Enabled: &enabled},
			}).
			Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"failed to set enabled on the new host",
				fmt.Sprintf("host %d: %s", id, errfmt.ErrMsg(err, hresp)),
			)
			taint(ctx, resp, plan, id)

			return
		}
	}

	// Read back the full state.
	getResp, hresp, err := client.HostsAPI.
		GetHost(ctx, sdk.GetHostIdParameter{Int64: &id}).
		Execute()
	if err != nil || getResp == nil || getResp.Server == nil {
		resp.Diagnostics.AddError(
			"failed to read host after create",
			fmt.Sprintf("host %d: %s", id, errfmt.ErrMsg(err, hresp)),
		)
		taint(ctx, resp, plan, id)

		return
	}

	bareMetalCodes := bareMetalCodesOrWarn(ctx, client, &resp.Diagnostics, id)

	getHostAsState(ctx, getResp.Server, bareMetalCodes, &plan)
	plan.WriteOnlyRevision = wodigest.NewRevision()

	// Baseline the appliance-side hash of the stored SSH password, and record
	// which clearable attributes the configuration set, for the next plan.
	resp.Diagnostics.Append(wo.Observe(ctx, resp.Private, getResp.Server)...)
	resp.Diagnostics.Append(clearable.Record(ctx, req.Config, resp.Private, clearableAttributes)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// pollProvisioned polls GetHost until the host leaves provisioning. It returns
// the final status and an error if it reached a failed state or timed out.
func (r *Resource) pollProvisioned(
	ctx context.Context,
	id int64,
	timeout time.Duration,
) (string, error) {
	client, err := r.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create API client: %w", err)
	}

	waitForState := func() (string, error) {
		getResp, hresp, err := client.HostsAPI.
			GetHost(ctx, sdk.GetHostIdParameter{Int64: &id}).
			Execute()
		if err != nil {
			// A connection failure or a 5xx is transient while the appliance is
			// busy provisioning (errfmt.IsTransient); anything else -- 401,
			// 403, 404, a refused re-login -- is final.
			if errfmt.IsTransient(ctx, err, hresp) {
				return "", backoff.RetryAfter(5)
			}

			return "", backoff.Permanent(
				fmt.Errorf("GetHost failed: %s", errfmt.ErrMsg(err, hresp)),
			)
		}

		if getResp == nil || getResp.Server == nil || getResp.Server.Status == nil {
			// Status not yet available; keep waiting.
			return "", backoff.RetryAfter(5)
		}

		status := *getResp.Server.Status

		switch classifyStatus(status) {
		case statusFailed:
			msg := status
			if getResp.Server.StatusMessage.IsSet() && getResp.Server.StatusMessage.Get() != nil {
				msg = *getResp.Server.StatusMessage.Get()
			}

			return status, backoff.Permanent(errors.New("host reached failed status: " + msg))
		case statusReady:
			return status, nil
		case statusProvisioning:
			return status, backoff.RetryAfter(5)
		default:
			// An unrecognised status is neither success nor failure: keep
			// polling, and let the timeout report it as the last status.
			return status, backoff.RetryAfter(5)
		}
	}

	return backoff.Retry(
		ctx,
		waitForState,
		backoff.WithBackOff(backoff.NewConstantBackOff(5*time.Second)),
		backoff.WithMaxElapsedTime(timeout),
	)
}

// taint records the host in state as tainted after a failure past the POST,
// so it is replaced rather than orphaned. Besides the id it keeps the inputs
// Delete reads from state -- delete_options and timeouts -- since the tainted
// host is destroyed from prior state and a configuration change cannot reach
// it.
func taint(ctx context.Context, resp *resource.CreateResponse, plan HostModel, id int64) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("delete_options"), plan.DeleteOptions)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("timeouts"), plan.Timeouts)...)
	cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
		ResourceType: "host",
		ResourceID:   id,
		StateWriter:  &resp.State,
		Diagnostics:  &resp.Diagnostics,
	})
}

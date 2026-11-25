// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instancepowerstate

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

// Interface compliance assertions.
var (
	_ resource.Resource                = &Resource{}
	_ resource.ResourceWithImportState = &Resource{}
	_ resource.ResourceWithModifyPlan  = &Resource{}
)

// No sweeper is needed for this resource because it does not create any
// API objects — it only changes the power state of an existing instance.

// Status constants used by the polling loop.
var (
	targetRunning   = []string{"running"}
	targetStopped   = []string{"stopped"}
	targetSuspended = []string{"suspended"}
	errorStatuses   = []string{"failed", "denied", "cancelled"}
)

// Server power states confirmed after an instance status transition. An
// instance status of "stopped" is set by the Morpheus "Stop Service" path too,
// so the server's own powerState is the only reliable signal that the
// workload has in fact been powered off (or on).
var (
	powerStateOn  = []string{"on"}
	powerStateOff = []string{"off"}
)

const defaultTimeout = 5 * time.Minute

// pollInterval is the pause between status polls against a real appliance.
// The polling functions take the interval as a parameter so that unit tests
// can shorten it for themselves without affecting anything else in the test
// binary -- the acceptance tests run the provider in-process.
const pollInterval = 5 * time.Second

// retryAfter asks backoff to try again after one poll interval.
func retryAfter(interval time.Duration) error {
	return &backoff.RetryAfterError{Duration: interval}
}

// InstancePowerStateModel is the Terraform state model.
type InstancePowerStateModel struct {
	Id           types.Int64    `tfsdk:"id"`
	InstanceId   types.Int64    `tfsdk:"instance_id"`
	DesiredState types.String   `tfsdk:"desired_state"`
	CurrentState types.String   `tfsdk:"current_state"`
	Timeouts     timeouts.Value `tfsdk:"timeouts"`
}

// Resource implements the hpe_instance_power_state resource.
type Resource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &Resource{}
}

// Metadata implements resource.Resource.
func (r *Resource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + "instance_power_state"
}

// Schema implements resource.Resource.
func (r *Resource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages the power state of an existing Morpheus instance. " +
			"This resource does not create or destroy instances — it only transitions " +
			"them between running, stopped, and suspended states.",
		MarkdownDescription: "Manages the power state of an existing Morpheus instance. " +
			"This resource does not create or destroy instances — it only transitions " +
			"them between `running`, `stopped`, and `suspended` states.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				Description:         "Resource identifier (equal to instance_id).",
				MarkdownDescription: "Resource identifier (equal to `instance_id`).",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"instance_id": schema.Int64Attribute{
				Required:            true,
				Description:         "The ID of the instance whose power state is managed.",
				MarkdownDescription: "The ID of the instance whose power state is managed.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"desired_state": schema.StringAttribute{
				Required:            true,
				Description:         "The desired power state: running, stopped, or suspended.",
				MarkdownDescription: "The desired power state: `running`, `stopped`, or `suspended`.",
				Validators: []validator.String{
					stringvalidator.OneOf("running", "stopped", "suspended"),
				},
			},
			"current_state": schema.StringAttribute{
				Computed:            true,
				Description:         "The current power state as reported by the API.",
				MarkdownDescription: "The current power state as reported by the API.",
			},
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Update: true,
			}),
		},
	}
}

// Create implements resource.Resource.
func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan InstancePowerStateModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createTimeout, diags := plan.Timeouts.Create(ctx, defaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to create API client",
			err.Error(),
		)

		return
	}

	instanceID := plan.InstanceId.ValueInt64()
	desiredState := plan.DesiredState.ValueString()

	// Check current state — if already at desired, no-op.
	getResp, hresp, err := client.InstancesAPI.GetInstance(ctx, instanceID).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to read instance",
			fmt.Sprintf("instance %d: %s", instanceID, errfmt.ErrMsg(err, hresp)),
		)

		return
	}

	if getResp.Instance == nil || getResp.Instance.Status == nil {
		resp.Diagnostics.AddError(
			"failed to read instance status",
			fmt.Sprintf("instance %d: response or status is nil", instanceID),
		)

		return
	}

	currentStatus := *getResp.Instance.Status
	if currentStatus == desiredState {
		// Already at desired state — no action needed.
		plan.Id = types.Int64Value(instanceID)
		plan.CurrentState = types.StringValue(currentStatus)
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

		return
	}

	// Fire the appropriate power action.
	if err := fireAction(ctx, client, instanceID, desiredState); err != nil {
		resp.Diagnostics.AddError(
			"failed to execute power action",
			fmt.Sprintf("instance %d → %s: %s", instanceID, desiredState, err.Error()),
		)

		return
	}

	// Poll until target state reached.
	finalStatus, pollErr := pollStatus(ctx, client, instanceID, desiredState, createTimeout, pollInterval)
	if pollErr != nil {
		resp.Diagnostics.AddError(
			"instance power state transition timed out or failed",
			fmt.Sprintf("instance %d → %s: %s (last status: %s)",
				instanceID, desiredState, pollErr.Error(), finalStatus),
		)

		return
	}

	plan.Id = types.Int64Value(instanceID)
	plan.CurrentState = types.StringValue(finalStatus)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read implements resource.Resource.
func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state InstancePowerStateModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to create API client",
			err.Error(),
		)

		return
	}

	instanceID := state.InstanceId.ValueInt64()

	getResp, hresp, err := client.InstancesAPI.GetInstance(ctx, instanceID).Execute()
	if err != nil {
		if hresp != nil && hresp.StatusCode == http.StatusNotFound {
			tflog.Info(ctx, fmt.Sprintf("instance %d not found, removing from state", instanceID))
			resp.State.RemoveResource(ctx)

			return
		}

		resp.Diagnostics.AddError(
			"failed to read instance",
			fmt.Sprintf("instance %d: %s", instanceID, errfmt.ErrMsg(err, hresp)),
		)

		return
	}

	if getResp.Instance == nil {
		resp.Diagnostics.AddError(
			"failed to read instance status",
			fmt.Sprintf("instance %d: response is nil", instanceID),
		)

		return
	}

	state.CurrentState = convert.StrToType(getResp.Instance.Status)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update implements resource.Resource.
func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan InstancePowerStateModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateTimeout, diags := plan.Timeouts.Update(ctx, defaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to create API client",
			err.Error(),
		)

		return
	}

	instanceID := plan.InstanceId.ValueInt64()
	desiredState := plan.DesiredState.ValueString()

	// Confirm the instance is reachable before acting on it.
	_, hresp, getErr := client.InstancesAPI.GetInstance(ctx, instanceID).Execute()
	if getErr != nil {
		resp.Diagnostics.AddError(
			"failed to read instance before update",
			fmt.Sprintf("instance %d: %s", instanceID, errfmt.ErrMsg(getErr, hresp)),
		)

		return
	}

	// Fire the appropriate power action.
	if err := fireAction(ctx, client, instanceID, desiredState); err != nil {
		resp.Diagnostics.AddError(
			"failed to execute power action",
			fmt.Sprintf("instance %d → %s: %s", instanceID, desiredState, err.Error()),
		)

		return
	}

	// Poll until target state reached.
	finalStatus, pollErr := pollStatus(ctx, client, instanceID, desiredState, updateTimeout, pollInterval)
	if pollErr != nil {
		resp.Diagnostics.AddError(
			"instance power state transition timed out or failed",
			fmt.Sprintf("instance %d → %s: %s (last status: %s)",
				instanceID, desiredState, pollErr.Error(), finalStatus),
		)

		return
	}

	plan.Id = types.Int64Value(instanceID)
	plan.CurrentState = types.StringValue(finalStatus)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete implements resource.Resource.
// Removing this resource from state does NOT change the instance power state.
func (r *Resource) Delete(
	ctx context.Context,
	_ resource.DeleteRequest,
	_ *resource.DeleteResponse,
) {
	tflog.Info(ctx, "instance_power_state resource removed from state; "+
		"no power action performed on delete")
}

// ImportState implements resource.ResourceWithImportState.
func (r *Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"invalid import ID",
			fmt.Sprintf("expected numeric instance ID, got %q: %s", req.ID, err.Error()),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("instance_id"), id)...)
}

// ModifyPlan implements resource.ResourceWithModifyPlan.
// If the refreshed current_state differs from desired_state, mark current_state
// as unknown so that the plan shows a diff and triggers an Update.
func (r *Resource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	// Skip on create (no state) or destroy (no plan).
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	var state InstancePowerStateModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan InstancePowerStateModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If current_state (from refreshed state) != desired_state (from plan/config),
	// mark current_state as unknown to force an update.
	if !state.CurrentState.IsNull() && !state.CurrentState.IsUnknown() &&
		!plan.DesiredState.IsNull() && !plan.DesiredState.IsUnknown() &&
		state.CurrentState.ValueString() != plan.DesiredState.ValueString() {
		resp.Diagnostics.Append(
			resp.Plan.SetAttribute(ctx, path.Root("current_state"), types.StringUnknown())...,
		)
	}
}

// fireAction calls the appropriate SDK power action.
//
// Start and stop are sent with server=true so that Morpheus powers the
// underlying server(s) off and on, as its own Stop/Start buttons do. Without
// the flag the API only runs the instance's "Stop Service"/"Start Service"
// path: on a layout with start/stop-phase scripts it executes those scripts,
// marks the instance stopped and leaves the VM running. The generated client
// does not send the parameter's documented default, so it must be explicit.
// Suspend ignores the flag server-side.
func fireAction(
	ctx context.Context,
	client *sdk.APIClient,
	instanceID int64,
	desiredState string,
) error {
	switch desiredState {
	case "running":
		_, hresp, err := client.InstancesAPI.StartInstance(ctx, instanceID).Server(true).Execute()
		if err != nil {
			return fmt.Errorf("StartInstance failed: %s", errfmt.ErrMsg(err, hresp))
		}
	case "stopped":
		_, hresp, err := client.InstancesAPI.StopInstance(ctx, instanceID).Server(true).Execute()
		if err != nil {
			return fmt.Errorf("StopInstance failed: %s", errfmt.ErrMsg(err, hresp))
		}
	case "suspended":
		_, hresp, err := client.InstancesAPI.SuspendInstance(ctx, instanceID).Execute()
		if err != nil {
			return fmt.Errorf("SuspendInstance failed: %s", errfmt.ErrMsg(err, hresp))
		}
	default:
		return fmt.Errorf("unknown desired_state: %s", desiredState)
	}

	return nil
}

// pollStatus polls GetInstance until the target state is reached, then — for
// running and stopped — polls each of the instance's servers until its
// powerState agrees, since the instance status alone does not prove the
// workload was powered on or off. An instance that reports no servers cannot
// be confirmed and is an error. Returns the final instance status.
func pollStatus(
	ctx context.Context,
	client *sdk.APIClient,
	instanceID int64,
	desiredState string,
	timeout time.Duration,
	interval time.Duration,
) (string, error) {
	target := targetForState(desiredState)
	deadline := time.Now().Add(timeout)

	var servers []int64

	waitForState := func() (string, error) {
		getResp, hresp, err := client.InstancesAPI.GetInstance(ctx, instanceID).Execute()
		if err != nil {
			if hresp == nil || hresp.StatusCode != http.StatusOK {
				return "", backoff.Permanent(
					fmt.Errorf("GetInstance failed: %s", errfmt.ErrMsg(err, hresp)),
				)
			}
		}

		if getResp.Instance == nil || getResp.Instance.Status == nil {
			return "", backoff.Permanent(
				fmt.Errorf("instance %d: status is nil", instanceID),
			)
		}

		servers = getResp.Instance.Servers
		status := *getResp.Instance.Status

		return status, checkStatusDone(status, target, errorStatuses, interval)
	}

	status, err := backoff.Retry(
		ctx,
		waitForState,
		backoff.WithBackOff(backoff.NewConstantBackOff(interval)),
		backoff.WithMaxElapsedTime(timeout),
	)
	if err != nil {
		return status, err
	}

	want := powerStateForState(desiredState)
	if want == nil {
		return status, nil
	}

	if len(servers) == 0 {
		return status, fmt.Errorf(
			"instance %d reports status %q but no servers, so the power state cannot be confirmed",
			instanceID, status,
		)
	}

	for _, serverID := range servers {
		remaining := time.Until(deadline)
		if remaining <= 0 {
			return status, fmt.Errorf(
				"instance %d: timed out before confirming the power state of server %d",
				instanceID, serverID,
			)
		}

		if err := pollServerPowerState(ctx, client, serverID, want, remaining, interval); err != nil {
			return status, fmt.Errorf("instance %d: %w", instanceID, err)
		}
	}

	return status, nil
}

// pollServerPowerState polls GetHost until the server reports one of the
// wanted power states or the timeout elapses; a timeout shorter than one
// poll interval allows exactly one attempt.
func pollServerPowerState(
	ctx context.Context,
	client *sdk.APIClient,
	serverID int64,
	want []string,
	timeout time.Duration,
	interval time.Duration,
) error {
	getID := sdk.GetHostIdParameter{Int64: &serverID}

	poll := func() (struct{}, error) {
		ret, hresp, err := client.HostsAPI.GetHost(ctx, getID).Execute()
		if err != nil {
			return struct{}{}, backoff.Permanent(
				fmt.Errorf("GetHost failed: %s", errfmt.ErrMsg(err, hresp)),
			)
		}

		if ret.Server == nil || ret.Server.PowerState == nil {
			return struct{}{}, retryAfter(interval)
		}

		if slices.Contains(want, *ret.Server.PowerState) {
			return struct{}{}, nil
		}

		return struct{}{}, retryAfter(interval)
	}

	_, err := backoff.Retry(
		ctx,
		poll,
		backoff.WithBackOff(backoff.NewConstantBackOff(interval)),
		backoff.WithMaxElapsedTime(timeout),
	)
	if err != nil {
		return fmt.Errorf("waiting for server %d power state %v: %w", serverID, want, err)
	}

	return nil
}

// checkStatusDone returns nil if target is reached, a permanent error if in
// an error state, or a retryable error otherwise.
func checkStatusDone(status string, targetStatuses, errStatuses []string, interval time.Duration) error {
	switch {
	case slices.Contains(errStatuses, status):
		return backoff.Permanent(errors.New("reached error status: " + status))
	case slices.Contains(targetStatuses, status):
		return nil
	default:
		return retryAfter(interval)
	}
}

// targetForState returns the target status slice for polling.
func targetForState(desired string) []string {
	switch desired {
	case "running":
		return targetRunning
	case "stopped":
		return targetStopped
	case "suspended":
		return targetSuspended
	default:
		return nil
	}
}

// powerStateForState returns the server power states that confirm the desired
// instance state, or nil when the instance status is the only signal (suspend
// is reported differently per cloud).
func powerStateForState(desired string) []string {
	switch desired {
	case "running":
		return powerStateOn
	case "stopped":
		return powerStateOff
	default:
		return nil
	}
}

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

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/constants"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
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

const defaultTimeout = 5 * time.Minute

// InstancePowerStateModel is the Terraform state model.
type InstancePowerStateModel struct {
	Id           types.Int64    `tfsdk:"id"`
	InstanceId   types.Int64    `tfsdk:"instance_id"`
	DesiredState types.String   `tfsdk:"desired_state"`
	CurrentState types.String   `tfsdk:"current_state"`
	Timeouts     timeouts.Value `tfsdk:"timeouts"`
}

// Resource implements the hpe_morpheus_instance_power_state resource.
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
	resp.TypeName = req.ProviderTypeName + "_" + constants.SubProviderName + "_instance_power_state"
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
	if err := r.fireAction(ctx, instanceID, desiredState, hresp); err != nil {
		resp.Diagnostics.AddError(
			"failed to execute power action",
			fmt.Sprintf("instance %d → %s: %s", instanceID, desiredState, err.Error()),
		)

		return
	}

	// Poll until target state reached.
	finalStatus, pollErr := r.pollStatus(ctx, instanceID, desiredState, createTimeout)
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

	if getResp.Instance == nil || getResp.Instance.Status == nil {
		resp.Diagnostics.AddError(
			"failed to read instance status",
			fmt.Sprintf("instance %d: response or status is nil", instanceID),
		)

		return
	}

	state.CurrentState = types.StringValue(*getResp.Instance.Status)
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

	// Fire the appropriate power action.
	_, hresp, getErr := client.InstancesAPI.GetInstance(ctx, instanceID).Execute()
	if getErr != nil {
		resp.Diagnostics.AddError(
			"failed to read instance before update",
			fmt.Sprintf("instance %d: %s", instanceID, errfmt.ErrMsg(getErr, hresp)),
		)

		return
	}

	if err := r.fireAction(ctx, instanceID, desiredState, hresp); err != nil {
		resp.Diagnostics.AddError(
			"failed to execute power action",
			fmt.Sprintf("instance %d → %s: %s", instanceID, desiredState, err.Error()),
		)

		return
	}

	// Poll until target state reached.
	finalStatus, pollErr := r.pollStatus(ctx, instanceID, desiredState, updateTimeout)
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
func (r *Resource) fireAction(
	ctx context.Context,
	instanceID int64,
	desiredState string,
	_ *http.Response,
) error {
	client, err := r.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	switch desiredState {
	case "running":
		_, hresp, err := client.InstancesAPI.StartInstance(ctx, instanceID).Execute()
		if err != nil {
			return fmt.Errorf("StartInstance failed: %s", errfmt.ErrMsg(err, hresp))
		}
	case "stopped":
		_, hresp, err := client.InstancesAPI.StopInstance(ctx, instanceID).Execute()
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

// pollStatus polls GetInstance until the target state is reached.
func (r *Resource) pollStatus(
	ctx context.Context,
	instanceID int64,
	desiredState string,
	timeout time.Duration,
) (string, error) {
	client, err := r.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create API client: %w", err)
	}

	target := targetForState(desiredState)

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

		status := *getResp.Instance.Status

		return status, checkStatusDone(status, target, errorStatuses)
	}

	return backoff.Retry(
		ctx,
		waitForState,
		backoff.WithBackOff(backoff.NewConstantBackOff(5*time.Second)),
		backoff.WithMaxElapsedTime(timeout),
	)
}

// checkStatusDone returns nil if target is reached, a permanent error if in
// an error state, or a retryable error otherwise.
func checkStatusDone(status string, targetStatuses, errStatuses []string) error {
	switch {
	case slices.Contains(errStatuses, status):
		return backoff.Permanent(errors.New("reached error status: " + status))
	case slices.Contains(targetStatuses, status):
		return nil
	default:
		return backoff.RetryAfter(5)
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

// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package hostpowerstate

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
	"github.com/hashicorp/terraform-plugin-framework/diag"
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
)

// Interface compliance assertions.
var (
	_ resource.Resource                = &Resource{}
	_ resource.ResourceWithImportState = &Resource{}
)

// No sweeper is needed: this resource creates no API objects, it only changes
// the power state of an existing host.

const defaultTimeout = 10 * time.Minute

var (
	targetOn      = []string{"on", "running"}
	targetOff     = []string{"off", "stopped"}
	errorStatuses = []string{"failed", "error"}
)

// HostPowerStateModel is the Terraform state model.
type HostPowerStateModel struct {
	Id         types.Int64    `tfsdk:"id"`
	HostId     types.Int64    `tfsdk:"host_id"`
	PowerState types.String   `tfsdk:"power_state"`
	Timeouts   timeouts.Value `tfsdk:"timeouts"`
}

// Resource implements the hpe_morpheus_host_power_state resource.
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
	resp.TypeName = req.ProviderTypeName + "_" + "host_power_state"
}

// Schema implements resource.Resource.
func (r *Resource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages the power state of an existing Morpheus host (compute " +
			"server). This resource does not create or destroy hosts — it only " +
			"transitions them between on and off states.",
		MarkdownDescription: "Manages the power state of an existing Morpheus host " +
			"(compute server). This resource does not create or destroy hosts — it " +
			"only transitions them between `on` and `off` states.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				Description:         "Resource identifier (equal to host_id).",
				MarkdownDescription: "Resource identifier (equal to `host_id`).",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"host_id": schema.Int64Attribute{
				Required:            true,
				Description:         "The ID of the host whose power state is managed.",
				MarkdownDescription: "The ID of the host whose power state is managed.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"power_state": schema.StringAttribute{
				Required:            true,
				Description:         "The desired power state: on or off.",
				MarkdownDescription: "The desired power state: `on` or `off`.",
				Validators: []validator.String{
					stringvalidator.OneOf("on", "off"),
				},
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
	var plan HostPowerStateModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Create(ctx, defaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.apply(ctx, &plan, timeout, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Update implements resource.Resource.
func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan HostPowerStateModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Update(ctx, defaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.apply(ctx, &plan, timeout, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// apply fires the power action (if needed) and polls until the host reaches the
// desired power state.
func (r *Resource) apply(
	ctx context.Context,
	plan *HostPowerStateModel,
	timeout time.Duration,
	diags *diag.Diagnostics,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		diags.AddError("failed to create API client", err.Error())

		return
	}

	hostID := plan.HostId.ValueInt64()
	desired := plan.PowerState.ValueString()

	current, hresp, err := getPowerState(ctx, client, hostID)
	if err != nil {
		diags.AddError(
			"failed to read host",
			fmt.Sprintf("host %d: %s", hostID, errfmt.ErrMsg(err, hresp)),
		)

		return
	}

	if !matches(current, desired) {
		if err := fireAction(ctx, client, hostID, desired); err != nil {
			diags.AddError(
				"failed to execute power action",
				fmt.Sprintf("host %d → %s: %s", hostID, desired, err.Error()),
			)

			return
		}

		final, pollErr := pollPower(ctx, client, hostID, desired, timeout)
		if pollErr != nil {
			diags.AddError(
				"host power state transition timed out or failed",
				fmt.Sprintf("host %d → %s: %s (last state: %s)",
					hostID, desired, pollErr.Error(), final),
			)

			return
		}
	}

	plan.Id = types.Int64Value(hostID)
}

// Read implements resource.Resource.
func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state HostPowerStateModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("failed to create API client", err.Error())

		return
	}

	hostID := state.HostId.ValueInt64()

	current, hresp, err := getPowerState(ctx, client, hostID)
	if err != nil {
		if hresp != nil && hresp.StatusCode == http.StatusNotFound {
			tflog.Info(ctx, fmt.Sprintf("host %d not found, removing from state", hostID))
			resp.State.RemoveResource(ctx)

			return
		}

		resp.Diagnostics.AddError(
			"failed to read host",
			fmt.Sprintf("host %d: %s", hostID, errfmt.ErrMsg(err, hresp)),
		)

		return
	}

	// Normalize to on/off for comparison with the configured value. A state
	// this resource does not manage (unknown, suspended, paused...) is not
	// stored: it would plan a change on every run and fire the action each
	// time; the prior value is kept and the practitioner told.
	switch {
	case matches(current, "on"):
		state.PowerState = types.StringValue("on")
	case matches(current, "off"):
		state.PowerState = types.StringValue("off")
	default:
		resp.Diagnostics.AddWarning(
			"Host power state not managed by this resource",
			fmt.Sprintf("host %d reports power state %q, which is neither on nor off; "+
				"power_state is left as %q in state. Resolve the host's state in Morpheus "+
				"if a change is wanted.", hostID, current, state.PowerState.ValueString()),
		)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete is a no-op: removing this resource does not change the host power state.
func (r *Resource) Delete(
	ctx context.Context,
	_ resource.DeleteRequest,
	_ *resource.DeleteResponse,
) {
	tflog.Info(ctx, "host_power_state resource removed from state; "+
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
			fmt.Sprintf("expected numeric host ID, got %q: %s", req.ID, err.Error()),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("host_id"), id)...)
}

// getPowerState reads the host's current power state string.
func getPowerState(
	ctx context.Context,
	client *sdk.APIClient,
	hostID int64,
) (string, *http.Response, error) {
	getResp, hresp, err := client.HostsAPI.
		GetHost(ctx, sdk.GetHostIdParameter{Int64: &hostID}).
		Execute()
	if err != nil {
		return "", hresp, err
	}

	if getResp == nil || getResp.Server == nil {
		return "", hresp, fmt.Errorf("host %d: response or server is nil", hostID)
	}

	// An absent power state is reported as "" and handled by the caller (Read
	// keeps the prior value with a warning; apply treats it as not matching).
	if getResp.Server.PowerState == nil {
		return "", hresp, nil
	}

	return *getResp.Server.PowerState, hresp, nil
}

// fireAction calls StartHost or StopHost.
func fireAction(
	ctx context.Context,
	client *sdk.APIClient,
	hostID int64,
	desired string,
) error {
	switch desired {
	case "on":
		_, hresp, err := client.HostsAPI.
			StartHost(ctx, sdk.UpdateHostIdParameter{Int64: &hostID}).
			Execute()
		if err != nil {
			return fmt.Errorf("StartHost failed: %s", errfmt.ErrMsg(err, hresp))
		}
	case "off":
		_, hresp, err := client.HostsAPI.
			StopHost(ctx, sdk.UpdateHostIdParameter{Int64: &hostID}).
			Execute()
		if err != nil {
			return fmt.Errorf("StopHost failed: %s", errfmt.ErrMsg(err, hresp))
		}
	default:
		return fmt.Errorf("unknown power_state: %s", desired)
	}

	return nil
}

// pollPower polls GetHost until the host reaches the desired power state.
func pollPower(
	ctx context.Context,
	client *sdk.APIClient,
	hostID int64,
	desired string,
	timeout time.Duration,
) (string, error) {
	target := targetOff
	if desired == "on" {
		target = targetOn
	}

	waitForState := func() (string, error) {
		state, hresp, err := getPowerState(ctx, client, hostID)
		if err != nil {
			// A connection failure or a 5xx is transient while the action runs
			// (errfmt.IsTransient); anything else -- 401, 403, 404, a refused
			// re-login -- is final.
			if errfmt.IsTransient(ctx, err, hresp) {
				return "", backoff.RetryAfter(5)
			}

			return "", backoff.Permanent(
				fmt.Errorf("GetHost failed: %s", errfmt.ErrMsg(err, hresp)),
			)
		}

		switch {
		case slices.Contains(errorStatuses, state):
			return state, backoff.Permanent(errors.New("reached error status: " + state))
		case slices.Contains(target, state):
			return state, nil
		default:
			return state, backoff.RetryAfter(5)
		}
	}

	return backoff.Retry(
		ctx,
		waitForState,
		backoff.WithBackOff(backoff.NewConstantBackOff(5*time.Second)),
		backoff.WithMaxElapsedTime(timeout),
	)
}

// matches reports whether the current power state satisfies the desired on/off.
func matches(current, desired string) bool {
	if desired == "on" {
		return slices.Contains(targetOn, current)
	}

	return slices.Contains(targetOff, current)
}

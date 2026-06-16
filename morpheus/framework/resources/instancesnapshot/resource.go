// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instancesnapshot

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &Resource{}
	_ resource.ResourceWithImportState = &Resource{}
)

var (
	snapshotTargetStatuses = []string{"complete"}
	snapshotErrorStatuses  = []string{"failed", "errored"}
)

// InstanceSnapshotModel and InstanceSnapshotResourceSchema are defined in
// the generated schema_gen.go.

func NewResource() resource.Resource {
	return &Resource{}
}

// Resource defines the resource implementation.
type Resource struct {
	configure.ResourceWithMorpheusConfigure
}

func (r *Resource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_instance_snapshot"
}

func (r *Resource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = InstanceSnapshotResourceSchema(ctx)
}

func checkStatusDone(status string, targetStatuses, errorStatuses []string) error {
	for _, s := range errorStatuses {
		if strings.EqualFold(status, s) {
			return backoff.Permanent(errors.New("reached error status: " + status))
		}
	}
	for _, s := range targetStatuses {
		if strings.EqualFold(status, s) {
			return nil
		}
	}

	return fmt.Errorf("snapshot status %q not yet in target set", status)
}

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan InstanceSnapshotModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createTimeout, diags := plan.Timeouts.Create(ctx, 30*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"create instance snapshot",
			"failed to create client: "+err.Error(),
		)

		return
	}

	instanceID := plan.InstanceId.ValueInt64()

	// Step 1: Capture pre-existing snapshot IDs
	existingSnaps, hresp, err := client.InstancesAPI.SnapshotsInstance(ctx, instanceID).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"create instance snapshot",
			fmt.Sprintf("failed to list existing snapshots for instance %d: %s",
				instanceID, errfmt.ErrMsg(err, hresp)),
		)

		return
	}

	preExistingIDs := make(map[int64]bool)
	if existingSnaps != nil {
		for _, snap := range existingSnaps.Snapshots {
			if snap.Id != nil {
				preExistingIDs[*snap.Id] = true
			}
		}
	}

	// Step 2: Fire snapshot creation
	snapshotReq := sdk.SnapshotInstanceRequest{
		Snapshot: &sdk.SnapshotInstanceRequestSnapshot{
			Name:           plan.Name.ValueStringPointer(),
			Description:    plan.Description.ValueStringPointer(),
			MemorySnapshot: plan.MemorySnapshot.ValueBoolPointer(),
			ForExport:      plan.ForExport.ValueBoolPointer(),
		},
	}

	createResp, hresp, err := client.InstancesAPI.SnapshotInstance(ctx, instanceID).
		SnapshotInstanceRequest(snapshotReq).
		Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"create instance snapshot",
			fmt.Sprintf("snapshot request failed for instance %d: %s",
				instanceID, errfmt.ErrMsg(err, hresp)),
		)

		return
	}

	if createResp != nil && createResp.Success != nil && !*createResp.Success {
		resp.Diagnostics.AddError(
			"create instance snapshot",
			fmt.Sprintf("snapshot request returned success=false for instance %d",
				instanceID),
		)

		return
	}

	// Step 3: Poll for the new snapshot to appear and reach "complete"
	snapshotName := plan.Name.ValueString()

	type pollResult struct {
		snap sdk.SnapshotsInstance200ResponseSnapshotsInner
	}

	waitForSnapshot := func() (*pollResult, error) {
		listResp, hresp, err := client.InstancesAPI.SnapshotsInstance(ctx, instanceID).Execute()
		if err != nil {
			return nil, backoff.Permanent(
				fmt.Errorf("failed to list snapshots: %s", errfmt.ErrMsg(err, hresp)),
			)
		}

		if listResp == nil {
			return nil, fmt.Errorf("list snapshots returned nil response")
		}

		for _, snap := range listResp.Snapshots {
			if snap.Id == nil || snap.Name == nil {
				continue
			}
			// Find snapshots matching our name that didn't exist before
			if *snap.Name != snapshotName {
				continue
			}
			if preExistingIDs[*snap.Id] {
				continue
			}

			// Found our new snapshot - check status
			status := ""
			if snap.Status != nil {
				status = *snap.Status
			}

			if err := checkStatusDone(status, snapshotTargetStatuses, snapshotErrorStatuses); err != nil {
				return nil, err
			}

			return &pollResult{snap: snap}, nil
		}

		return nil, fmt.Errorf("snapshot %q not yet found in list", snapshotName)
	}

	result, err := backoff.Retry(
		ctx,
		waitForSnapshot,
		backoff.WithBackOff(backoff.NewConstantBackOff(5*time.Second)),
		backoff.WithMaxElapsedTime(createTimeout),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"create instance snapshot",
			fmt.Sprintf("snapshot %q on instance %d failed or timed out: %v",
				snapshotName, instanceID, errors.Unwrap(err)),
		)

		// Attempt to taint if we can identify the snapshot
		listResp, _, _ := client.InstancesAPI.SnapshotsInstance(ctx, instanceID).Execute()
		if listResp != nil {
			for _, snap := range listResp.Snapshots {
				if snap.Id != nil && snap.Name != nil &&
					*snap.Name == snapshotName && !preExistingIDs[*snap.Id] {
					cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
						ResourceType: "instance_snapshot",
						ResourceID:   *snap.Id,
						StateWriter:  &resp.State,
						Diagnostics:  &resp.Diagnostics,
					})

					return
				}
			}
		}

		return
	}

	// Step 4: Set state from the found snapshot
	snap := result.snap
	snapshotID := *snap.Id

	plan.Id = types.Int64Value(snapshotID)
	plan.Status = types.StringValue("")
	if snap.Status != nil {
		plan.Status = types.StringValue(*snap.Status)
	}
	plan.ExternalId = types.StringNull()
	if snap.ExternalId.IsSet() && snap.ExternalId.Get() != nil {
		plan.ExternalId = types.StringValue(*snap.ExternalId.Get())
	}
	plan.DateCreated = types.StringNull()
	if snap.DateCreated != nil {
		plan.DateCreated = types.StringValue(snap.DateCreated.Format(time.RFC3339))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state InstanceSnapshotModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"read instance snapshot",
			"failed to create client: "+err.Error(),
		)

		return
	}

	snapshotID := state.Id.ValueInt64()

	snapResp, hresp, err := client.InstancesAPI.GetSnapshotInstance(ctx, snapshotID).Execute()
	if err != nil {
		if hresp != nil && hresp.StatusCode == http.StatusNotFound {
			tflog.Warn(ctx, fmt.Sprintf("Snapshot %d not found, removing from state", snapshotID))
			resp.State.RemoveResource(ctx)

			return
		}

		resp.Diagnostics.AddError(
			"read instance snapshot",
			fmt.Sprintf("snapshot %d GET failed: %s", snapshotID, errfmt.ErrMsg(err, hresp)),
		)

		return
	}

	if snapResp == nil || snapResp.Snapshot == nil {
		tflog.Warn(ctx, fmt.Sprintf("Snapshot %d returned nil, removing from state", snapshotID))
		resp.State.RemoveResource(ctx)

		return
	}

	snap := snapResp.Snapshot

	// Update computed fields from API
	if snap.Status != nil {
		state.Status = types.StringValue(*snap.Status)
	}
	if snap.ExternalId.IsSet() && snap.ExternalId.Get() != nil {
		state.ExternalId = types.StringValue(*snap.ExternalId.Get())
	} else {
		state.ExternalId = types.StringNull()
	}
	if snap.DateCreated != nil {
		state.DateCreated = types.StringValue(snap.DateCreated.Format(time.RFC3339))
	}

	// Refresh config-driven fields from the API so out-of-band changes are
	// detected. name, description and memory_snapshot are RequiresReplace, so a
	// difference between the API and the configuration triggers a replacement.
	// instance_id and retain_on_delete are not returned by the API, so they are
	// preserved from prior state.
	state.Name = convert.StrToType(snap.Name)
	state.Description = convert.StrToType(snap.Description.Get())
	state.MemorySnapshot = convert.BoolToType(snap.MemorySnapshot)
	state.ForExport = convert.BoolToType(snap.ForExport)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Only retain_on_delete is updatable (no API call needed).
	// All other attribute changes trigger RequiresReplace.
	var plan, state InstanceSnapshotModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Carry over computed fields from prior state
	plan.Id = state.Id
	plan.Status = state.Status
	plan.ExternalId = state.ExternalId
	plan.DateCreated = state.DateCreated

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state InstanceSnapshotModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If retain_on_delete is true, skip the API delete
	if state.RetainOnDelete.ValueBool() {
		tflog.Info(ctx, fmt.Sprintf(
			"retain_on_delete=true for snapshot %d, removing from state only",
			state.Id.ValueInt64(),
		))

		return
	}

	deleteTimeout, diags := state.Timeouts.Delete(ctx, 10*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"delete instance snapshot",
			"failed to create client: "+err.Error(),
		)

		return
	}

	snapshotID := state.Id.ValueInt64()

	_, hresp, err := client.InstancesAPI.DeleteSnapshotInstance(ctx, snapshotID).Execute()
	if err != nil {
		if hresp != nil && hresp.StatusCode == http.StatusNotFound {
			tflog.Warn(ctx, fmt.Sprintf("Snapshot %d already gone", snapshotID))

			return
		}

		resp.Diagnostics.AddError(
			"delete instance snapshot",
			fmt.Sprintf("snapshot %d DELETE failed: %s", snapshotID, errfmt.ErrMsg(err, hresp)),
		)

		return
	}
}

func (r *Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// Import format: "<instance_id>.<snapshot_id>"
	parts := strings.Split(req.ID, ".")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"import instance snapshot",
			fmt.Sprintf(
				"invalid import ID format %q, expected \"<instance_id>.<snapshot_id>\"",
				req.ID,
			),
		)

		return
	}

	instanceID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"import instance snapshot",
			fmt.Sprintf("invalid instance_id %q: %v", parts[0], err),
		)

		return
	}

	snapshotID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"import instance snapshot",
			fmt.Sprintf("invalid snapshot_id %q: %v", parts[1], err),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("instance_id"), instanceID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), snapshotID)...)
}

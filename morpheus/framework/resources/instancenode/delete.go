// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instancenode

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

const defaultDeleteTimeout = 60 * time.Minute

// Delete implements resource.Resource.
// Removes a node from the instance via the container's remove-node action.
// The server is returned to the pool, not destroyed. Delete completion is
// "the container is gone from the instance", NOT "the server is gone" —
// the plugin's cleanServer returns removeServer:false and the server
// persists as available in its pool.
func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state instanceNodeModel

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

	containerID := state.ContainerID.ValueInt64()
	instanceID := state.InstanceID.ValueInt64()

	// Discover the remove-node action code.
	actionCode, err := discoverRemoveNodeAction(ctx, client, containerID)
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to discover remove-node action",
			fmt.Sprintf("container %d: %s", containerID, err.Error()),
		)

		return
	}

	// Execute the remove-node action. No request body — the controller
	// hard-codes [:].
	_, hresp, err := client.ContainersAPI.
		ExecuteContainerAction(ctx, containerID).
		Code(actionCode).
		Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to remove node",
			fmt.Sprintf("container %d action %q: %s",
				containerID, actionCode, errfmt.ErrMsg(err, hresp)),
		)

		return
	}

	// Poll until the container is gone from the instance.
	if err := waitForContainerRemoval(
		ctx, client, instanceID, containerID, deleteTimeout,
	); err != nil {
		resp.Diagnostics.AddError(
			"node removal timed out",
			fmt.Sprintf("container %d on instance %d: %s",
				containerID, instanceID, err.Error()),
		)

		return
	}

	tflog.Info(ctx, fmt.Sprintf(
		"node removed: container %d gone from instance %d",
		containerID, instanceID))
}

// discoverRemoveNodeAction finds the remove-node action code from the
// container's available actions.
func discoverRemoveNodeAction(
	ctx context.Context,
	client *sdk.APIClient,
	containerID int64,
) (string, error) {
	actResp, hresp, err := client.ContainersAPI.
		GetContainerActions(ctx, containerID).Execute()
	if err != nil {
		return "", fmt.Errorf("GetContainerActions: %s",
			errfmt.ErrMsg(err, hresp))
	}

	if actResp == nil {
		return "", fmt.Errorf("GetContainerActions returned nil")
	}

	for i := range actResp.Actions {
		a := &actResp.Actions[i]
		if a.Name != nil && *a.Name == "Remove Node" {
			if a.Code != nil {
				return *a.Code, nil
			}
		}
	}

	// Fallback: code ending with -remove-node.
	for i := range actResp.Actions {
		a := &actResp.Actions[i]
		if a.Code != nil && strings.HasSuffix(*a.Code, "-remove-node") {
			return *a.Code, nil
		}
	}

	return "", fmt.Errorf("no remove-node action found among %d actions",
		len(actResp.Actions))
}

// waitForContainerRemoval polls the instance until the container is no
// longer present in containerDetails. Do NOT wait for the server to
// disappear — it persists as available in its pool.
func waitForContainerRemoval(
	ctx context.Context,
	client *sdk.APIClient,
	instanceID int64,
	containerID int64,
	timeout time.Duration,
) error {
	poll := func() (struct{}, error) {
		getResp, hresp, err := client.InstancesAPI.
			GetInstance(ctx, instanceID).Execute()
		if err != nil {
			return struct{}{}, backoff.Permanent(
				fmt.Errorf("GetInstance: %s", errfmt.ErrMsg(err, hresp)),
			)
		}

		if getResp == nil || getResp.Instance == nil {
			// Instance gone — container is therefore gone.
			return struct{}{}, nil
		}

		for i := range getResp.Instance.ContainerDetails {
			cd := &getResp.Instance.ContainerDetails[i]
			if cd.Id != nil && *cd.Id == containerID {
				return struct{}{}, fmt.Errorf(
					"container %d still present", containerID)
			}
		}

		// Container is gone — done.
		return struct{}{}, nil
	}

	_, err := backoff.Retry(
		ctx,
		poll,
		backoff.WithBackOff(backoff.NewConstantBackOff(10*time.Second)),
		backoff.WithMaxElapsedTime(timeout),
	)

	return err
}

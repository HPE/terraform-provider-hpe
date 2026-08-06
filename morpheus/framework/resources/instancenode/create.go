// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instancenode

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/containerip"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
)

const defaultCreateTimeout = 90 * time.Minute

// pollForNewContainerTimeout is the maximum time to poll for a new container
// after the add-node action returns. This is intentionally short (60s) because
// addInstanceContainers creates container records synchronously inside the HTTP
// request — only the subsequent provisioning is queued asynchronously. If a new
// container has not appeared within this window, it will never appear; the most
// likely causes are an appliance licence limit, capacity constraint, or
// provisioning policy denial (the endpoint returns HTTP 200 with success:true
// even when the action was silently refused). Do not raise this value.
const pollForNewContainerTimeout = 60 * time.Second

// pollForNewContainerInterval is the retry interval within the bounded poll
// window. Kept short because we are only absorbing read-after-write lag.
const pollForNewContainerInterval = 5 * time.Second

// Create implements resource.Resource.
func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan instanceNodeModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
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

	instanceID := plan.InstanceID.ValueInt64()

	// Step 1: Read instance.
	getResp, hresp, err := client.InstancesAPI.GetInstance(ctx, instanceID).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to read instance",
			fmt.Sprintf("instance %d: %s", instanceID, errfmt.ErrMsg(err, hresp)),
		)

		return
	}

	if getResp == nil || getResp.Instance == nil {
		resp.Diagnostics.AddError(
			"failed to read instance",
			fmt.Sprintf("instance %d: response is nil", instanceID),
		)

		return
	}

	// Metal guard: only when resource_pool_id is set.
	poolSet := !plan.ResourcePoolID.IsNull() && !plan.ResourcePoolID.IsUnknown()
	if poolSet && !IsMetal(getResp.Instance) {
		resp.Diagnostics.AddError(
			notMetalSummary,
			notMetalDetail(instanceID,
				provisionTypeCodeOrUnknown(getResp.Instance)),
		)

		return
	}

	// Step 2: Discover the add-node action code.
	actionCode, err := discoverAddNodeAction(ctx, client, instanceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to discover add-node action",
			fmt.Sprintf("instance %d: %s", instanceID, err.Error()),
		)

		return
	}

	// Snapshot existing container IDs before adding a node. This snapshot
	// is used by both the provisioning and pre-provisioned paths as the
	// baseline for diff-based container identification.
	existingContainerIDs := make(map[int64]bool)
	for i := range getResp.Instance.ContainerDetails {
		cd := &getResp.Instance.ContainerDetails[i]
		if cd.Id != nil {
			existingContainerIDs[*cd.Id] = true
		}
	}

	// Step 3: Build the action envelope and execute.
	envelope := buildAddNodeEnvelope(&plan, actionCode)

	// The instance and container action endpoints discard a top-level JSON
	// request body: InstancesController.action() never merges request.JSON
	// into params, and Grails binds only form-encoded bodies on PUT. The
	// controller does explicitly merge an "action" envelope, so parameters
	// go inside it. Verified against a live appliance on 2026-08-05: an
	// envelope carrying count=2 produced two containers, proving non-code
	// keys are delivered.
	//
	// Do NOT set the typed fields on ExecuteInstanceActionRequest (Code,
	// SelectedResourcePoolId, ...) - they serialise at the top level and
	// would be silently discarded, placing the node in the instance's own
	// pool while returning 200 success. Tracked as MORPH-15280; this stays
	// until the oldest supported appliance carries the fix.
	sdkReq := sdk.ExecuteInstanceActionRequest{
		AdditionalProperties: map[string]any{"action": envelope},
	}

	actionResp, hresp, err := client.InstancesAPI.
		ExecuteInstanceAction(ctx, instanceID).
		Code(actionCode).
		ExecuteInstanceActionRequest(sdkReq).
		Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to add node",
			fmt.Sprintf("instance %d action %q: %s",
				instanceID, actionCode, errfmt.ErrMsg(err, hresp)),
		)

		return
	}

	// Step 4: Extract the new container ID.
	// The response may carry the ID in results[<instanceId>].containers[0].id
	// (fast path), but this is not guaranteed: capacity/policy denials cause
	// scaleInstanceWithProvisioning to return early with no containers, and
	// the pre-provisioned branch (scaleInstanceWithConvertToManaged) never
	// returns container IDs at all. When the fast path fails, we poll the
	// instance for a short bounded window (pollForNewContainerTimeout) to
	// absorb read-after-write lag only — container records are created
	// synchronously, so a longer wait would be pointless.
	containerID, err := extractContainerID(
		actionResp, instanceID, existingContainerIDs,
		ctx, client, pollForNewContainerTimeout,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to identify new container",
			fmt.Sprintf("instance %d: %s", instanceID, err.Error()),
		)

		return
	}

	plan.ContainerID = types.Int64Value(containerID)

	// Step 5: Resolve server_id and ip_address from the instance.
	readErr := refreshNodeState(ctx, client, instanceID, containerID, &plan)
	if readErr != nil {
		// The node exists on the appliance but its state could not be read.
		// Taint rather than drop it: dropping would leak a real node, and
		// leaving it untainted would keep a resource whose state we could
		// not confirm.
		resp.Diagnostics.AddError(
			"node created but state read failed",
			fmt.Sprintf("instance %d container %d was added but the state "+
				"read failed: %s", instanceID, containerID, readErr.Error()),
		)

		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "instance_node",
			ResourceID:   containerID,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	// Step 6: Optionally wait for IP address.
	if plan.WaitForIPAddress.ValueBool() {
		ip, warned, waitErr := containerip.Wait(
			ctx, client, instanceID, containerID, createTimeout,
		)
		if waitErr != nil {
			// The node exists on the appliance, so returning without
			// recording it would leak it.
			resp.Diagnostics.AddError(
				"IP address wait failed",
				fmt.Sprintf("instance %d container %d: %s",
					instanceID, containerID, waitErr.Error()),
			)

			cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
				ResourceType: "instance_node",
				ResourceID:   containerID,
				StateWriter:  &resp.State,
				Diagnostics:  &resp.Diagnostics,
			})

			return
		}

		if !warned && ip != "" {
			plan.IPAddress = types.StringValue(ip)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// discoverAddNodeAction finds the add-node action code from the instance's
// available actions. The code is per-layout: generic-add-node on some,
// ubuntu-add-node on others.
func discoverAddNodeAction(
	ctx context.Context,
	client *sdk.APIClient,
	instanceID int64,
) (string, error) {
	actResp, hresp, err := client.InstancesAPI.
		GetInstanceActions(ctx, instanceID).Execute()
	if err != nil {
		return "", fmt.Errorf("GetInstanceActions: %s", errfmt.ErrMsg(err, hresp))
	}

	if actResp == nil {
		return "", fmt.Errorf("GetInstanceActions returned nil")
	}

	for i := range actResp.Actions {
		a := &actResp.Actions[i]
		if a.Name != nil && *a.Name == "Add Node" {
			if a.Code != nil {
				return *a.Code, nil
			}
		}
	}

	// Fallback: code ending with -add-node.
	for i := range actResp.Actions {
		a := &actResp.Actions[i]
		if a.Code != nil && strings.HasSuffix(*a.Code, "-add-node") {
			return *a.Code, nil
		}
	}

	return "", fmt.Errorf("no add-node action found among %d actions",
		len(actResp.Actions))
}

// buildAddNodeEnvelope constructs the JSON envelope that goes inside
// {"action": ...} in the request body. selectedResourcePoolId is included
// only when resource_pool_id is set in the plan; omitting it entirely is
// the correct path for virtual (non-metal) instances.
func buildAddNodeEnvelope(plan *instanceNodeModel, actionCode string) map[string]any {
	env := map[string]any{
		"code":  actionCode,
		"count": 1,
	}

	// Include selectedResourcePoolId only when resource_pool_id is set.
	if !plan.ResourcePoolID.IsNull() && !plan.ResourcePoolID.IsUnknown() {
		env["selectedResourcePoolId"] = fmt.Sprintf("pool-%d",
			plan.ResourcePoolID.ValueInt64())
	}

	// preProvisioned is evaluated with Groovy truth on the raw value, and
	// the spec's enum is [on, off]. The string "off" is truthy and takes
	// the pre-provisioned branch. When pre_provisioned is false or null,
	// omit the key entirely rather than sending "off".
	if !plan.PreProvisioned.IsNull() && !plan.PreProvisioned.IsUnknown() &&
		plan.PreProvisioned.ValueBool() {
		env["preProvisioned"] = "on"

		if !plan.SelectedServerID.IsNull() && !plan.SelectedServerID.IsUnknown() {
			env["selectedServerId"] = plan.SelectedServerID.ValueInt64()
		}

		if !plan.SshHost.IsNull() && !plan.SshHost.IsUnknown() {
			env["sshHost"] = plan.SshHost.ValueString()
		}

		if !plan.SshUsername.IsNull() && !plan.SshUsername.IsUnknown() {
			env["sshUsername"] = plan.SshUsername.ValueString()
		}

		if !plan.SshPassword.IsNull() && !plan.SshPassword.IsUnknown() {
			env["sshPassword"] = plan.SshPassword.ValueString()
		}

		if !plan.SshKeyPairID.IsNull() && !plan.SshKeyPairID.IsUnknown() {
			env["sshKeyPair"] = map[string]any{
				"id": plan.SshKeyPairID.ValueInt64(),
			}
		}
	}

	return env
}

// extractContainerID retrieves the new container ID from the action response.
// It first tries the fast path: extracting from results[<instanceId>].containers[0].id.
// If that fails for any reason (capacity denial, pre-provisioned path, missing
// data), it falls back to polling the instance until a new container appears
// that was not in the pre-add snapshot.
//
// When multiple new containers appear, the lowest ID is chosen deterministically
// to avoid depending on map or slice iteration order.
func extractContainerID(
	actionResp *sdk.ExecuteInstanceAction200Response,
	instanceID int64,
	existingIDs map[int64]bool,
	ctx context.Context,
	client *sdk.APIClient,
	pollTimeout time.Duration,
) (int64, error) {
	// Fast path: try to extract from the response regardless of preProvisioned.
	if actionResp != nil {
		if results, ok := actionResp.AdditionalProperties["results"]; ok {
			id, err := extractContainerIDFromResults(results, instanceID)
			if err == nil {
				tflog.Debug(ctx, "container ID from response fast path",
					map[string]any{"instance_id": instanceID, "container_id": id})

				return id, nil
			}

			tflog.Debug(ctx, "fast path failed, falling back to poll",
				map[string]any{"instance_id": instanceID, "reason": err.Error()})
		}
	}

	// Poll path: repeatedly read the instance until a new container appears.
	return pollForNewContainer(ctx, client, instanceID, existingIDs, pollTimeout)
}

// pollForNewContainer polls the instance until a container ID appears that is
// not in existingIDs, or until the timeout expires. The timeout should be short
// (pollForNewContainerTimeout) because container records are created
// synchronously — this poll exists only to absorb read-after-write lag.
func pollForNewContainer(
	ctx context.Context,
	client *sdk.APIClient,
	instanceID int64,
	existingIDs map[int64]bool,
	timeout time.Duration,
) (int64, error) {
	tflog.Debug(ctx, "polling for new container via diff",
		map[string]any{"instance_id": instanceID, "timeout": timeout.String()})

	deadline := time.Now().Add(timeout)

	for {
		id, found, err := diffContainers(ctx, client, instanceID, existingIDs)
		if err != nil {
			return 0, err
		}

		if found {
			tflog.Debug(ctx, "new container found via polling",
				map[string]any{"instance_id": instanceID, "container_id": id})

			return id, nil
		}

		if time.Now().Add(pollForNewContainerInterval).After(deadline) {
			return 0, fmt.Errorf(
				"the add-node action returned success but no new container was created on "+
					"instance %d after polling for %s. The Morpheus add-node endpoint returns "+
					"HTTP 200 with success:true even when the action is silently refused. "+
					"The most likely causes are: (1) an appliance licence or capacity limit "+
					"has been reached, (2) a provisioning policy is denying the request, or "+
					"(3) the instance's layout does not support scaling. Check the Morpheus "+
					"appliance activity log for details",
				instanceID, timeout,
			)
		}

		tflog.Debug(ctx, "no new container yet, retrying",
			map[string]any{
				"instance_id": instanceID,
				"interval":    pollForNewContainerInterval.String(),
			})

		select {
		case <-ctx.Done():
			return 0, fmt.Errorf(
				"context cancelled while polling for new container: %w", ctx.Err(),
			)
		case <-time.After(pollForNewContainerInterval):
		}
	}
}

// diffContainers reads the instance and returns the lowest new container ID
// not present in existingIDs. Returns (0, false, nil) if no new container is
// found yet.
func diffContainers(
	ctx context.Context,
	client *sdk.APIClient,
	instanceID int64,
	existingIDs map[int64]bool,
) (int64, bool, error) {
	getResp, hresp, err := client.InstancesAPI.GetInstance(ctx, instanceID).Execute()
	if err != nil {
		return 0, false, fmt.Errorf("failed to re-read instance: %s",
			errfmt.ErrMsg(err, hresp))
	}

	if getResp == nil || getResp.Instance == nil {
		return 0, false, fmt.Errorf(
			"instance %d: response is nil after add-node", instanceID,
		)
	}

	id, found := findNewContainerInDetails(getResp.Instance.ContainerDetails, existingIDs)

	return id, found, nil
}

// findNewContainerInDetails returns the lowest container ID from details that
// is not present in existingIDs. When multiple new containers appear (e.g. a
// race with another caller), the lowest ID is chosen deterministically rather
// than depending on slice iteration order.
func findNewContainerInDetails(
	details []sdk.InstanceContainer2,
	existingIDs map[int64]bool,
) (int64, bool) {
	var lowestNew int64

	found := false

	for i := range details {
		cd := &details[i]
		if cd.Id != nil && !existingIDs[*cd.Id] {
			if !found || *cd.Id < lowestNew {
				lowestNew = *cd.Id
				found = true
			}
		}
	}

	return lowestNew, found
}

// extractContainerIDFromResults parses the results map from the action response.
// Expected shape: {"<instanceId>": {"containers": [{"id": <int>}]}}.
func extractContainerIDFromResults(results any, instanceID int64) (int64, error) {
	resultsMap, ok := results.(map[string]any)
	if !ok {
		return 0, fmt.Errorf("results is not a map: %T", results)
	}

	instanceKey := fmt.Sprintf("%d", instanceID)

	instanceResult, ok := resultsMap[instanceKey]
	if !ok {
		return 0, fmt.Errorf("no results for instance %d", instanceID)
	}

	instanceMap, ok := instanceResult.(map[string]any)
	if !ok {
		return 0, fmt.Errorf("instance result is not a map: %T", instanceResult)
	}

	containers, ok := instanceMap["containers"]
	if !ok {
		return 0, fmt.Errorf("no containers in instance result")
	}

	containerList, ok := containers.([]any)
	if !ok {
		return 0, fmt.Errorf("containers is not a list: %T", containers)
	}

	if len(containerList) == 0 {
		return 0, fmt.Errorf("containers list is empty")
	}

	containerMap, ok := containerList[0].(map[string]any)
	if !ok {
		return 0, fmt.Errorf("container is not a map: %T", containerList[0])
	}

	idVal, ok := containerMap["id"]
	if !ok {
		return 0, fmt.Errorf("container has no id field")
	}

	switch v := idVal.(type) {
	case float64:
		return int64(v), nil
	case json.Number:
		return v.Int64()
	case int64:
		return v, nil
	default:
		return 0, fmt.Errorf("container id has unexpected type: %T", idVal)
	}
}

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

	// Snapshot existing container IDs for the pre-provisioned diff path.
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
	containerID, err := extractContainerID(
		actionResp, instanceID, existingContainerIDs,
		plan.PreProvisioned.ValueBool(), ctx, client,
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
// For the provisioning branch, it comes from results[<instanceId>].containers[0].id.
// For the pre-provisioned branch, it diffs containerDetails before/after.
func extractContainerID(
	actionResp *sdk.ExecuteInstanceAction200Response,
	instanceID int64,
	existingIDs map[int64]bool,
	preProvisioned bool,
	ctx context.Context,
	client *sdk.APIClient,
) (int64, error) {
	if !preProvisioned && actionResp != nil {
		// Try to extract from the response's AdditionalProperties.
		if results, ok := actionResp.AdditionalProperties["results"]; ok {
			return extractContainerIDFromResults(results, instanceID)
		}
	}

	// Pre-provisioned path or response didn't carry container IDs:
	// diff containerDetails.
	tflog.Debug(ctx, "extracting container ID via diff",
		map[string]any{"instance_id": instanceID})

	getResp, hresp, err := client.InstancesAPI.GetInstance(ctx, instanceID).Execute()
	if err != nil {
		return 0, fmt.Errorf("failed to re-read instance: %s",
			errfmt.ErrMsg(err, hresp))
	}

	if getResp == nil || getResp.Instance == nil {
		return 0, fmt.Errorf("instance %d: response is nil after add-node", instanceID)
	}

	for i := range getResp.Instance.ContainerDetails {
		cd := &getResp.Instance.ContainerDetails[i]
		if cd.Id != nil && !existingIDs[*cd.Id] {
			return *cd.Id, nil
		}
	}

	return 0, fmt.Errorf("no new container found after add-node")
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

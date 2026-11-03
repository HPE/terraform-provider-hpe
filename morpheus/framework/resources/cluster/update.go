// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cluster

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	updateOperation = "update cluster resource"
)

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state ClusterModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get timeout from HCL if set, the default is 45 minutes
	updateTimeout, diags := plan.Timeouts.Update(ctx, 45*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, updateTimeout)
	defer cancel()

	id := state.Id.ValueInt64()

	updateClusterReq := &sdk.UpdateClusterRequest{}
	cluster := &sdk.UpdateClusterRequestCluster{}

	// name
	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		cluster.Name = plan.Name.ValueStringPointer()
	}

	// Attributes removed from the configuration, to be cleared: sent when
	// known, sent empty when being cleared, not sent otherwise.
	pending, pdiags := clearable.Compute(ctx, req.Config, req.State, req.Private, clearableAttributes)
	resp.Diagnostics.Append(pdiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cluster.Description = pending.String("description", plan.Description)

	labels, sendLabels, ldiags := pending.Strings(ctx, "labels", plan.Labels)
	resp.Diagnostics.Append(ldiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	switch {
	case sendLabels && len(labels) == 0:
		// The API rejects an empty labels list as an unrecognised field (its
		// handler only renames a non-empty one to labelArray, the key it
		// binds), so an empty list is sent under that key directly.
		if cluster.AdditionalProperties == nil {
			cluster.AdditionalProperties = map[string]interface{}{}
		}

		cluster.AdditionalProperties["labelArray"] = []string{}
	case sendLabels:
		cluster.Labels = labels
	}

	switch {
	case !plan.Config.IsNull() && !plan.Config.IsUnknown():
		configValue := plan.Config.UnderlyingValue()
		configAny, err := convert.ValueToAny(ctx, configValue)
		if err != nil {
			resp.Diagnostics.AddError(
				updateOperation,
				"cluster: failed to convert config: "+
					err.Error(),
			)

			return
		}

		configDataMap, ok := configAny.(map[string]any)
		if !ok {
			resp.Diagnostics.AddError(
				"error updating cluster",
				"could not parse config value",
			)

			return
		}

		cluster.Config = &sdk.UpdateClusterRequestClusterConfig{
			AdditionalProperties: configDataMap,
		}
	// HVM config updatable fields
	case !plan.ConfigHvm.IsNull() && !plan.ConfigHvm.IsUnknown():
		config := &sdk.UpdateClusterRequestClusterConfig{
			AdditionalProperties: make(map[string]any),
		}

		if !plan.ConfigHvm.CpuArch.IsNull() && !plan.ConfigHvm.CpuArch.IsUnknown() {
			config.AdditionalProperties["cpuArch"] = plan.ConfigHvm.CpuArch.ValueString()
		}

		if !plan.ConfigHvm.CpuModel.IsNull() && !plan.ConfigHvm.CpuModel.IsUnknown() {
			config.AdditionalProperties["cpuModel"] = plan.ConfigHvm.CpuModel.ValueString()
		}

		if !plan.ConfigHvm.DynamicPlacement.IsNull() && !plan.ConfigHvm.DynamicPlacement.IsUnknown() {
			config.DynamicPlacementMode = convert.BoolTypeToStringPointerOnOff(plan.ConfigHvm.DynamicPlacement)
		}

		if !plan.ConfigHvm.VcpuPlacementMode.IsNull() && !plan.ConfigHvm.VcpuPlacementMode.IsUnknown() {
			config.AdditionalProperties["vcpuPlacementMode"] = plan.ConfigHvm.VcpuPlacementMode.ValueString()
		}

		if !plan.ConfigHvm.PowerPolicy.IsNull() && !plan.ConfigHvm.PowerPolicy.IsUnknown() {
			config.AdditionalProperties["powerPolicy"] = plan.ConfigHvm.PowerPolicy.ValueString()
		}

		// HVM 2.0 network configuration. virt_switch_name is requires_replace
		// (virtual switches are provisioned per host at create/scale-out time),
		// so an in-place update never changes it; it is re-sent here only so a
		// config rebuild does not drop it. The remaining network fields are
		// reconciled onto the hosts by the daily cluster drift-correction sync,
		// so a change is an in-place update that converges on the next sync.
		if !plan.ConfigHvm.VirtSwitchName.IsNull() && !plan.ConfigHvm.VirtSwitchName.IsUnknown() {
			config.AdditionalProperties["virtSwitchName"] = plan.ConfigHvm.VirtSwitchName.ValueString()
		}

		if !plan.ConfigHvm.Dns1.IsNull() && !plan.ConfigHvm.Dns1.IsUnknown() {
			config.AdditionalProperties["dns1"] = plan.ConfigHvm.Dns1.ValueString()
		}

		if !plan.ConfigHvm.Dns2.IsNull() && !plan.ConfigHvm.Dns2.IsUnknown() {
			config.AdditionalProperties["dns2"] = plan.ConfigHvm.Dns2.ValueString()
		}

		if !plan.ConfigHvm.Dns3.IsNull() && !plan.ConfigHvm.Dns3.IsUnknown() {
			config.AdditionalProperties["dns3"] = plan.ConfigHvm.Dns3.ValueString()
		}

		if !plan.ConfigHvm.SearchDomain.IsNull() && !plan.ConfigHvm.SearchDomain.IsUnknown() {
			config.AdditionalProperties["searchDomain"] = plan.ConfigHvm.SearchDomain.ValueString()
		}

		if !plan.ConfigHvm.Ntp1.IsNull() && !plan.ConfigHvm.Ntp1.IsUnknown() {
			config.AdditionalProperties["ntp1"] = plan.ConfigHvm.Ntp1.ValueString()
		}

		if !plan.ConfigHvm.Ntp2.IsNull() && !plan.ConfigHvm.Ntp2.IsUnknown() {
			config.AdditionalProperties["ntp2"] = plan.ConfigHvm.Ntp2.ValueString()
		}

		if !plan.ConfigHvm.Proxy.IsNull() && !plan.ConfigHvm.Proxy.IsUnknown() {
			config.AdditionalProperties["proxy"] = plan.ConfigHvm.Proxy.ValueString()
		}

		if !plan.ConfigHvm.NoProxy.IsNull() && !plan.ConfigHvm.NoProxy.IsUnknown() {
			config.AdditionalProperties["noProxy"] = plan.ConfigHvm.NoProxy.ValueString()
		}

		cluster.Config = config

	}

	updateClusterReq.Cluster = cluster

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	// Gate the HVM 2.0 network fields on the selected layout's clusterVersion.
	// layout_id is requires_replace, so the layout cannot change in place --
	// but the network fields (other than virt_switch_name) are now in-place
	// updatable, so a user could add them to a cluster whose fixed layout
	// predates clusterVersion 2.0, where Morpheus would silently ignore them.
	// Refuse that here, mirroring the create-path gate (see create.go).
	if !plan.ConfigHvm.IsNull() && !plan.ConfigHvm.IsUnknown() &&
		!plan.LayoutId.IsNull() && !plan.LayoutId.IsUnknown() {
		resp.Diagnostics.Append(requireHvm2Layout(
			ctx, client, plan.LayoutId.ValueInt64(), plan.ConfigHvm,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	_, httpResp, err := client.ClustersAPI.UpdateCluster(ctx, id).
		UpdateClusterRequest(*updateClusterReq).Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			updateOperation,
			fmt.Sprintf("cluster %d PUT failed: ", id)+errfmt.ErrMsg(err, httpResp),
		)

		return
	}

	// The SSH password cannot be updated (a change replaces the cluster); an
	// update planned through write_only_revision only records the digests of
	// the configured values -- the one-time upgrade from a provider without
	// them, or a deprecated marker being removed.
	plan.WriteOnlyRevision = state.WriteOnlyRevision
	if wodigest.Resend(ctx, req.Plan) {
		plan.WriteOnlyRevision = wodigest.NewRevision()
	}

	updatedState, found, diag := getClusterAsState(ctx, id, client, plan)
	if resp.Diagnostics.Append(diag...); resp.Diagnostics.HasError() {
		return
	}
	if !found {
		resp.Diagnostics.AddError(readOperation, fmt.Sprintf("cluster %d not found after update", id))

		return
	}

	// Record the digests once the read-back succeeded, so private state and
	// state move together, and re-baseline the appliance-hash drift detection
	// off a fresh read of the cluster's hosts (see write_only.go).
	resp.Diagnostics.Append(wo.StoreObserved(ctx, resp.Private, req.Config,
		clusterHostHashes(ctx, client, id))...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

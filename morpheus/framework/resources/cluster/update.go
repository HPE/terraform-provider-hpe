// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cluster

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
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

	updateClusterReq := sdk.NewUpdateClusterRequest()
	cluster := sdk.NewUpdateClusterRequestCluster()

	// name
	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		cluster.SetName(plan.Name.ValueString())
	}

	// description
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		cluster.SetDescription(plan.Description.ValueString())
	}

	// labels
	if !plan.Labels.IsNull() && !plan.Labels.IsUnknown() {
		labels, err := convert.SetToStrSlice(plan.Labels)
		if err != nil {
			resp.Diagnostics.AddError(
				updateOperation,
				"cluster "+plan.Name.ValueString()+": failed to parse label: "+err.Error(),
			)

			return
		}

		cluster.SetLabels(labels)
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

		cluster.SetConfig(sdk.UpdateClusterRequestClusterConfig{
			AdditionalProperties: configDataMap,
		})
	// HVM config updatable fields
	case !plan.ConfigHvm.IsNull() && !plan.ConfigHvm.IsUnknown():
		config := sdk.NewUpdateClusterRequestClusterConfig()

		config.AdditionalProperties = make(map[string]any)

		if !plan.ConfigHvm.CpuArch.IsNull() && !plan.ConfigHvm.CpuArch.IsUnknown() {
			config.AdditionalProperties["cpuArch"] = plan.ConfigHvm.CpuArch.ValueString()
		}

		if !plan.ConfigHvm.CpuModel.IsNull() && !plan.ConfigHvm.CpuModel.IsUnknown() {
			config.AdditionalProperties["cpuModel"] = plan.ConfigHvm.CpuModel.ValueString()
		}

		if !plan.ConfigHvm.DynamicPlacement.IsNull() && !plan.ConfigHvm.DynamicPlacement.IsUnknown() {
			config.SetDynamicPlacementMode(
				*convert.BoolTypeToStringPointerOnOff(plan.ConfigHvm.DynamicPlacement),
			)
		}

		if !plan.ConfigHvm.VcpuPlacementMode.IsNull() && !plan.ConfigHvm.VcpuPlacementMode.IsUnknown() {
			config.AdditionalProperties["vcpuPlacementMode"] = plan.ConfigHvm.VcpuPlacementMode.ValueString()
		}

		if !plan.ConfigHvm.PowerPolicy.IsNull() && !plan.ConfigHvm.PowerPolicy.IsUnknown() {
			config.AdditionalProperties["powerPolicy"] = plan.ConfigHvm.PowerPolicy.ValueString()
		}

		cluster.SetConfig(*config)

	}

	updateClusterReq.SetCluster(*cluster)

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
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

	updatedState, diag := getClusterAsState(ctx, id, client, plan)
	if resp.Diagnostics.Append(diag...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

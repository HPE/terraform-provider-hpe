// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cluster

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/sdkfuncs"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	createOperation = "create cluster resource"
)

var (
	CreateTargetStatuses = []string{clusterStatusOk}

	CreateErrorStatuses = []string{
		clusterStatusCancelled,
		clusterStatusDenied,
		clusterStatusFailed,
		clusterStatusDeprovisioned,
		clusterStatusDeprovisioning,
		clusterStatusPendingRemoval,
		clusterStatusStopping,
		clusterStatusRemoved,
		clusterStatusRemoving,
		clusterStatusWarning,
	}
)

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, config ClusterModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get timeout from HCL if set, the default is 45 minutes
	createTimeout, diags := plan.Timeouts.Create(ctx, 45*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	addClusterReq := &sdk.AddClusterRequest{}
	addClusterReq.Cluster = &sdk.AddClusterRequestCluster{}

	if !plan.CloudId.IsNull() && !plan.CloudId.IsUnknown() {
		addClusterReq.Cluster.Cloud.Id = plan.CloudId.ValueInt64Pointer()
	}

	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		addClusterReq.Cluster.Name = plan.Name.ValueString()
	}

	// Cluster type code
	// Read from TF config if using generic config,
	// otherwise set it based off the static config used.
	switch {
	case !plan.Config.IsNull():
		if !plan.ClusterTypeCode.IsNull() && !plan.ClusterTypeCode.IsUnknown() {
			v := plan.ClusterTypeCode.ValueString()
			addClusterReq.Cluster.Type = sdk.AddClusterRequestClusterType{String: &v}
		}
	case !plan.ConfigHvm.IsNull() && !plan.ConfigHvm.IsUnknown():
		v := clusterTypeCodeMVM
		addClusterReq.Cluster.Type = sdk.AddClusterRequestClusterType{String: &v}
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		addClusterReq.Cluster.Description = plan.Description.ValueStringPointer()
	}

	if !plan.GroupId.IsNull() && !plan.GroupId.IsUnknown() {
		addClusterReq.Cluster.Group = &sdk.AddClusterRequestClusterGroup{
			Id: plan.GroupId.ValueInt64Pointer(),
		}
	}

	if !plan.LayoutId.IsNull() && !plan.LayoutId.IsUnknown() {
		addClusterReq.Cluster.Layout.Id = plan.LayoutId.ValueInt64Pointer()
	}

	if !plan.WorkflowId.IsNull() && !plan.WorkflowId.IsUnknown() {
		addClusterReq.Cluster.TaskSetId = plan.WorkflowId.ValueInt64Pointer()
	}

	// server
	if !plan.Server.IsNull() && !plan.Server.IsUnknown() {
		diags := buildCreateClusterServer(ctx, plan, config, addClusterReq)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)

			return
		}
	}

	// labels
	if !plan.Labels.IsNull() && !plan.Labels.IsUnknown() {
		labels, err := convert.SetToStrSlice(plan.Labels)
		if err != nil {
			resp.Diagnostics.AddError(
				createOperation,
				"cluster "+plan.Name.ValueString()+": failed to parse label: "+err.Error(),
			)

			return
		}

		addClusterReq.Cluster.Labels = labels
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	clusterResp, httpResp, err := client.ClustersAPI.AddCluster(ctx).
		AddClusterRequest(*addClusterReq).Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(createOperation, errfmt.ErrMsg(err, httpResp))

		return
	}

	clusterId := *clusterResp.Cluster.Id
	plan.Id = convert.Int64ToType(&clusterId)

	// write the ID now
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	taintResourceState := func(id int64) {
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "cluster",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})
	}

	// Polling
	waitForReady := func() (string, error) {
		resp, hresp, err := client.ClustersAPI.GetCluster(ctx, clusterId).Execute()
		if err != nil {
			if hresp == nil || hresp.StatusCode != http.StatusOK {
				return "", backoff.Permanent(err)
			}
		}

		// Get cluster
		cluster := resp.Cluster
		if cluster == nil {
			return "", backoff.Permanent(fmt.Errorf("cluster %d: GET returned empty cluster", clusterId))
		}

		// Get status
		status := cluster.Status
		if status == nil {
			return "", backoff.Permanent(fmt.Errorf("cluster %d: GET returned empty status", clusterId))
		}

		return *status, checkStatusDone(
			*status,
			CreateTargetStatuses,
			CreateErrorStatuses,
		)
	}

	if status, err := backoff.Retry(
		ctx,
		waitForReady,
		backoff.WithBackOff(backoff.NewConstantBackOff(5*time.Second)),
		backoff.WithMaxElapsedTime(createTimeout),
	); err != nil {
		if status == "" {
			errUnwrapped := errors.Unwrap(err)
			if errUnwrapped != nil {
				resp.Diagnostics.AddError(
					"create cluster resource",
					fmt.Sprintf(
						"cluster %d: provisioning failed: %v",
						clusterId,
						errUnwrapped,
					),
				)
			} else {
				resp.Diagnostics.AddError(
					"create cluster resource",
					fmt.Sprintf(
						"cluster %d: provisioning failed - unknown error",
						clusterId,
					),
				)
			}
		} else {
			resp.Diagnostics.AddError(
				"create cluster resource",
				fmt.Sprintf(
					"cluster %d: provisioning failed current status is: %s",
					clusterId,
					status,
				),
			)
		}
		taintResourceState(clusterId)

		return
	}

	state, diag := getClusterAsState(ctx, clusterId, client, plan)
	if resp.Diagnostics.Append(diag...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func buildCreateClusterServer(
	ctx context.Context,
	plan ClusterModel,
	config ClusterModel,
	req *sdk.AddClusterRequest,
) diag.Diagnostics {
	diags := diag.Diagnostics{}
	name := plan.Name.ValueString()
	server := &sdk.AddClusterRequestClusterServer{}

	// the basic types
	if !plan.Server.Hostname.IsNull() && !plan.Server.Hostname.IsUnknown() {
		server.Hostname = *sdk.NewNullableString(plan.Server.Hostname.ValueStringPointer())
	}

	if !plan.Server.Name.IsNull() && !plan.Server.Name.IsUnknown() {
		server.Name = plan.Server.Name.ValueString()
	}

	if !plan.Server.NetworkDomain.IsNull() && !plan.Server.NetworkDomain.IsUnknown() {
		server.NetworkDomain = *sdk.NewNullableString(plan.Server.NetworkDomain.ValueStringPointer())
	}

	if !plan.Server.ServicePlanId.IsNull() && !plan.Server.ServicePlanId.IsUnknown() {
		server.Plan = sdk.AddClusterRequestClusterServerPlan{
			Id: plan.Server.ServicePlanId.ValueInt64Pointer(),
		}
	}

	if !plan.Server.SshPort.IsNull() && !plan.Server.SshPort.IsUnknown() {
		server.SshPort = plan.Server.SshPort.ValueInt64Pointer()
	}

	// need config for getting the value of ssh password - it's write-only so isn't in plan
	if !config.Server.SshPasswordWo.IsNull() && !config.Server.SshPasswordWo.IsUnknown() {
		server.SshPassword = *sdk.NewNullableString(config.Server.SshPasswordWo.ValueStringPointer())
	}

	if !plan.Server.SshUsername.IsNull() && !plan.Server.SshUsername.IsUnknown() {
		server.SshUsername = plan.Server.SshUsername.ValueStringPointer()
	}

	if !plan.Server.SshKeyPairId.IsNull() && !plan.Server.SshKeyPairId.IsUnknown() {
		server.SshKeyPair = &sdk.AddClusterRequestClusterServerSshKeyPair{
			Id: plan.Server.SshKeyPairId.ValueInt64Pointer(),
		}
	}

	if !plan.Server.UserGroupId.IsNull() && !plan.Server.UserGroupId.IsUnknown() {
		server.UserGroup = &sdk.AddClusterRequestClusterServerUserGroup{
			Id: plan.Server.UserGroupId.ValueInt64Pointer(),
		}
	}

	if !plan.Server.Visibility.IsNull() && !plan.Server.Visibility.IsUnknown() {
		server.Visibility = plan.Server.Visibility.ValueStringPointer()
	}

	if !plan.Server.DataDevice.IsNull() && !plan.Server.DataDevice.IsUnknown() {
		server.DataDevice = plan.Server.DataDevice.ValueStringPointer()
	}

	if !plan.Server.LvmEnabled.IsNull() && !plan.Server.LvmEnabled.IsUnknown() {
		server.LvmEnabled = plan.Server.LvmEnabled.ValueBoolPointer()
	}

	if !plan.Server.ManagementNetInterface.IsNull() && !plan.Server.ManagementNetInterface.IsUnknown() {
		server.Network = &sdk.AddClusterRequestClusterServerNetwork{
			Name: plan.Server.ManagementNetInterface.ValueString(),
		}
	}

	// the sets

	// network interfaces
	if !plan.Server.NetworkInterfaces.IsNull() && !plan.Server.NetworkInterfaces.IsUnknown() {
		nis, diags := convert.FromSetType(ctx, plan.Server.NetworkInterfaces, createNetworkInterfacesMapper)
		if diags.HasError() {
			return diags
		}

		server.NetworkInterfaces = nis

	}

	// security groups
	if !plan.Server.SecurityGroups.IsNull() && !plan.Server.SecurityGroups.IsUnknown() {
		securityGroups, err := convert.SetToStrSlice(plan.Server.SecurityGroups)
		if err != nil {
			diags.AddError("failed to convert security_groups to slice of strings", err.Error())

			return diags
		}

		server.SecurityGroups = securityGroups
	}

	// ssh hosts
	if !plan.Server.SshHosts.IsNull() && !plan.Server.SshHosts.IsUnknown() {
		hosts, diags := convert.FromSetType(ctx, plan.Server.SshHosts, createSSHHostsMapper)
		if diags.HasError() {
			return diags
		}

		server.SshHosts = hosts
	}

	// tags
	if !plan.Server.Tags.IsNull() && !plan.Server.Tags.IsUnknown() {
		var tags []TagsValue
		diags := plan.Server.Tags.ElementsAs(ctx, &tags, false)
		if diags.HasError() {
			return diags
		}

		var addTags []sdk.AddClusterRequestClusterServerTagsInner

		for _, v := range tags {
			addTags = append(addTags, sdk.AddClusterRequestClusterServerTagsInner{
				Name:  v.Name.ValueStringPointer(),
				Value: v.Value.ValueStringPointer(),
			})
		}

		server.Tags = addTags

	}

	// volumes
	if !plan.Server.Volumes.IsNull() && !plan.Server.Volumes.IsUnknown() {
		volumes, diags := convert.FromSetType(ctx, plan.Server.Volumes, createVolumeMapper)
		if diags.HasError() {
			return diags
		}

		server.Volumes = volumes
	}

	// set the server.config based off which config block is used...
	switch {
	// Generic config
	case !plan.Config.IsNull() && !plan.Config.IsUnknown():
		configValue := plan.Config.UnderlyingValue()
		configValueAny, err := convert.ValueToAny(ctx, configValue)
		if err != nil {
			diags.AddError(
				createOperation,
				"cluster "+name+": failed to convert config: "+err.Error(),
			)

			return diags
		}

		configDataMap, ok := configValueAny.(map[string]any)
		if !ok {
			diags.AddError(
				createOperation,
				"cluster "+name+": could not parse config value",
			)

			return diags
		}

		server.Config.MapmapOfStringAny = &configDataMap
	// HVM
	case !plan.ConfigHvm.IsNull() && !plan.ConfigHvm.IsUnknown():
		cfg := sdkfuncs.NewHvmClusterServerConfig()

		if !plan.ConfigHvm.CpuArch.IsNull() && !plan.ConfigHvm.CpuArch.IsUnknown() {
			cfg.CpuArch = plan.ConfigHvm.CpuArch.ValueStringPointer()
		}

		if !plan.ConfigHvm.CpuModel.IsNull() && !plan.ConfigHvm.CpuModel.IsUnknown() {
			cfg.CpuModel = plan.ConfigHvm.CpuModel.ValueStringPointer()
		}

		if !plan.ConfigHvm.DynamicPlacement.IsNull() && !plan.ConfigHvm.DynamicPlacement.IsUnknown() {
			cfg.DynamicPlacementMode = convert.BoolTypeToStringPointerOnOff(plan.ConfigHvm.DynamicPlacement)
		}

		if !plan.ConfigHvm.PowerPolicy.IsNull() && !plan.ConfigHvm.PowerPolicy.IsUnknown() {
			cfg.PowerPolicy = plan.ConfigHvm.PowerPolicy.ValueStringPointer()
		}

		if !plan.ConfigHvm.StorageInterfaceName.IsNull() && !plan.ConfigHvm.StorageInterfaceName.IsUnknown() {
			cfg.StorageInterfaceName = plan.ConfigHvm.StorageInterfaceName.ValueStringPointer()
		}

		if !plan.ConfigHvm.ComputeInterfaceName.IsNull() && !plan.ConfigHvm.ComputeInterfaceName.IsUnknown() {
			cfg.ComputeInterfaceName = plan.ConfigHvm.ComputeInterfaceName.ValueStringPointer()
		}

		if !plan.ConfigHvm.ComputeVlans.IsNull() && !plan.ConfigHvm.ComputeVlans.IsUnknown() {
			cfg.ComputeVlans = plan.ConfigHvm.ComputeVlans.ValueStringPointer()
		}

		if !plan.ConfigHvm.OverlayInterfaceName.IsNull() && !plan.ConfigHvm.OverlayInterfaceName.IsUnknown() {
			cfg.OverlayInterfaceName = plan.ConfigHvm.OverlayInterfaceName.ValueStringPointer()
		}

		if !plan.ConfigHvm.CreateUser.IsNull() && !plan.ConfigHvm.CreateUser.IsUnknown() {
			cfg.CreateUser = plan.ConfigHvm.CreateUser.ValueBoolPointer()
		}

		cfgAnyOf := sdkfuncs.NewHvmClusterServerConfigAsAnyOf(cfg)
		server.Config.AddClusterRequestClusterServerConfigAnyOf = &cfgAnyOf
	}

	if req.Cluster == nil {
		diags.AddError("build cluster server payload", "cluster is nil")

		return diags
	}

	req.Cluster.Server = *server

	return nil
}

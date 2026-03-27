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

	addClusterReq := sdk.NewAddClusterRequestWithDefaults()
	addClusterReq.Cluster = sdk.NewAddClusterRequestClusterWithDefaults()

	if !plan.CloudId.IsNull() && !plan.CloudId.IsUnknown() {
		addClusterReq.Cluster.Cloud.SetId(plan.CloudId.ValueInt64())
	}

	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		addClusterReq.Cluster.SetName(plan.Name.ValueString())
	}

	// Cluster type code
	// Read from TF config if using generic config,
	// otherwise set it based off the static config used.
	switch {
	case !plan.Config.IsNull():
		if !plan.ClusterTypeCode.IsNull() && !plan.ClusterTypeCode.IsUnknown() {
			v := sdk.StringAsAddClusterRequestClusterType(plan.ClusterTypeCode.ValueStringPointer())
			addClusterReq.Cluster.SetType(v)
		}
	case !plan.ConfigHvm.IsNull() && !plan.ConfigHvm.IsUnknown():
		v := clusterTypeCodeMVM
		addClusterReq.Cluster.SetType(sdk.StringAsAddClusterRequestClusterType(&v))
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		addClusterReq.Cluster.SetDescription(plan.Description.ValueString())
	}

	if !plan.GroupId.IsNull() && !plan.GroupId.IsUnknown() {
		addClusterReq.Cluster.Group = sdk.NewAddClusterRequestClusterGroupWithDefaults()
		addClusterReq.Cluster.Group.SetId(plan.GroupId.ValueInt64())
	}

	if !plan.LayoutId.IsNull() && !plan.LayoutId.IsUnknown() {
		addClusterReq.Cluster.Layout.SetId(plan.LayoutId.ValueInt64())
	}

	if !plan.WorkflowId.IsNull() && !plan.WorkflowId.IsUnknown() {
		addClusterReq.Cluster.SetTaskSetId(plan.WorkflowId.ValueInt64())
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

		addClusterReq.Cluster.SetLabels(labels)
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

	clusterId := clusterResp.Cluster.GetId()
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
		cluster, ok := resp.GetClusterOk()
		if !ok || cluster == nil {
			return "", backoff.Permanent(fmt.Errorf("cluster %d: GET returned empty cluster", clusterId))
		}

		// Get status
		status, ok := cluster.GetStatusOk()
		if !ok || status == nil {
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
	server := sdk.NewAddClusterRequestClusterServerWithDefaults()

	// the basic types
	if !plan.Server.Hostname.IsNull() && !plan.Server.Hostname.IsUnknown() {
		server.SetHostname(plan.Server.Hostname.ValueString())
	}

	if !plan.Server.Name.IsNull() && !plan.Server.Name.IsUnknown() {
		server.SetName(plan.Server.Name.ValueString())
	}

	if !plan.Server.NetworkDomain.IsNull() && !plan.Server.NetworkDomain.IsUnknown() {
		server.SetNetworkDomain(plan.Server.NetworkDomain.ValueString())
	}

	if !plan.Server.ServicePlanId.IsNull() && !plan.Server.ServicePlanId.IsUnknown() {
		server.Plan = *sdk.NewAddClusterRequestClusterServerPlan()
		server.Plan.SetId(plan.Server.ServicePlanId.ValueInt64())
	}

	if !plan.Server.SshPort.IsNull() && !plan.Server.SshPort.IsUnknown() {
		server.SetSshPort(plan.Server.SshPort.ValueInt64())
	}

	// need config for getting the value of ssh password - it's write-only so isn't in plan
	if !config.Server.SshPasswordWo.IsNull() && !config.Server.SshPasswordWo.IsUnknown() {
		server.SetSshPassword(config.Server.SshPasswordWo.ValueString())
	}

	if !plan.Server.SshUsername.IsNull() && !plan.Server.SshUsername.IsUnknown() {
		server.SetSshUsername(plan.Server.SshUsername.ValueString())
	}

	if !plan.Server.SshKeyPairId.IsNull() && !plan.Server.SshKeyPairId.IsUnknown() {
		server.SetSshKeyPair(sdk.AddClusterRequestClusterServerSshKeyPair{
			Id: plan.Server.SshKeyPairId.ValueInt64Pointer(),
		})
	}

	if !plan.Server.UserGroupId.IsNull() && !plan.Server.UserGroupId.IsUnknown() {
		server.UserGroup = sdk.NewAddClusterRequestClusterServerUserGroup()
		server.UserGroup.SetId(plan.Server.UserGroupId.ValueInt64())
	}

	if !plan.Server.Visibility.IsNull() && !plan.Server.Visibility.IsUnknown() {
		server.SetVisibility(plan.Server.Visibility.ValueString())
	}

	if !plan.Server.DataDevice.IsNull() && !plan.Server.DataDevice.IsUnknown() {
		server.SetDataDevice(plan.Server.DataDevice.ValueString())
	}

	if !plan.Server.LvmEnabled.IsNull() && !plan.Server.LvmEnabled.IsUnknown() {
		server.SetLvmEnabled(plan.Server.LvmEnabled.ValueBool())
	}

	if !plan.Server.ManagementNetInterface.IsNull() && !plan.Server.ManagementNetInterface.IsUnknown() {
		server.Network = sdk.NewAddClusterRequestClusterServerNetwork(plan.Server.ManagementNetInterface.ValueString())
	}

	// the sets

	// network interfaces
	if !plan.Server.NetworkInterfaces.IsNull() && !plan.Server.NetworkInterfaces.IsUnknown() {
		nis, diags := convert.FromSetType(ctx, plan.Server.NetworkInterfaces, createNetworkInterfacesMapper)
		if diags.HasError() {
			return diags
		}

		server.SetNetworkInterfaces(nis)

	}

	// security groups
	if !plan.Server.SecurityGroups.IsNull() && !plan.Server.SecurityGroups.IsUnknown() {
		securityGroups, err := convert.SetToStrSlice(plan.Server.SecurityGroups)
		if err != nil {
			diags.AddError("failed to convert security_groups to slice of strings", err.Error())

			return diags
		}

		server.SetSecurityGroups(securityGroups)
	}

	// ssh hosts
	if !plan.Server.SshHosts.IsNull() && !plan.Server.SshHosts.IsUnknown() {
		hosts, diags := convert.FromSetType(ctx, plan.Server.SshHosts, createSSHHostsMapper)
		if diags.HasError() {
			return diags
		}

		server.SetSshHosts(hosts)
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

		server.SetTags(addTags)

	}

	// volumes
	if !plan.Server.Volumes.IsNull() && !plan.Server.Volumes.IsUnknown() {
		volumes, diags := convert.FromSetType(ctx, plan.Server.Volumes, createVolumeMapper)
		if diags.HasError() {
			return diags
		}

		server.SetVolumes(volumes)
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
			cfg.SetCpuArch(plan.ConfigHvm.CpuArch.ValueString())
		}

		if !plan.ConfigHvm.CpuModel.IsNull() && !plan.ConfigHvm.CpuModel.IsUnknown() {
			cfg.SetCpuModel(plan.ConfigHvm.CpuModel.ValueString())
		}

		if !plan.ConfigHvm.DynamicPlacement.IsNull() && !plan.ConfigHvm.DynamicPlacement.IsUnknown() {
			cfg.SetDynamicPlacementMode(*convert.BoolTypeToStringPointerOnOff(plan.ConfigHvm.DynamicPlacement))
		}

		if !plan.ConfigHvm.PowerPolicy.IsNull() && !plan.ConfigHvm.PowerPolicy.IsUnknown() {
			cfg.SetPowerPolicy(plan.ConfigHvm.PowerPolicy.ValueString())
		}

		if !plan.ConfigHvm.StorageInterfaceName.IsNull() && !plan.ConfigHvm.StorageInterfaceName.IsUnknown() {
			cfg.SetStorageInterfaceName(plan.ConfigHvm.StorageInterfaceName.ValueString())
		}

		if !plan.ConfigHvm.ComputeInterfaceName.IsNull() && !plan.ConfigHvm.ComputeInterfaceName.IsUnknown() {
			cfg.SetComputeInterfaceName(plan.ConfigHvm.ComputeInterfaceName.ValueString())
		}

		if !plan.ConfigHvm.ComputeVlans.IsNull() && !plan.ConfigHvm.ComputeVlans.IsUnknown() {
			cfg.SetComputeVlans(plan.ConfigHvm.ComputeVlans.ValueString())
		}

		if !plan.ConfigHvm.OverlayInterfaceName.IsNull() && !plan.ConfigHvm.OverlayInterfaceName.IsUnknown() {
			cfg.SetOverlayInterfaceName(plan.ConfigHvm.OverlayInterfaceName.ValueString())
		}

		if !plan.ConfigHvm.CreateUser.IsNull() && !plan.ConfigHvm.CreateUser.IsUnknown() {
			cfg.SetCreateUser(plan.ConfigHvm.CreateUser.ValueBool())
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

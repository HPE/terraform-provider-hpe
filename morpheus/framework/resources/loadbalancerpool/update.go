// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancerpool

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, currentState LoadBalancerPoolModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &currentState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	pool := sdk.NewUpdateLoadBalancerPoolRequestLoadBalancerPoolWithDefaults()

	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		pool.SetName(plan.Name.ValueString())
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		pool.SetDescription(plan.Description.ValueString())
	}

	if !plan.VipBalance.IsNull() && !plan.VipBalance.IsUnknown() {
		pool.SetVipBalance(plan.VipBalance.ValueString())
	}

	if !plan.MinActive.IsNull() && !plan.MinActive.IsUnknown() {
		pool.SetMinActive(plan.MinActive.ValueInt64())
	}

	if !plan.Port.IsNull() && !plan.Port.IsUnknown() {
		portVal, parseErr := strconv.ParseInt(plan.Port.ValueString(), 10, 64)
		if parseErr != nil {
			resp.Diagnostics.AddError(
				"error updating load balancer pool",
				fmt.Sprintf("invalid port value %q: %s", plan.Port.ValueString(), parseErr),
			)

			return
		}

		pool.SetPort(portVal)
	}

	if !plan.VipSticky.IsNull() && !plan.VipSticky.IsUnknown() {
		pool.SetVipSticky(plan.VipSticky.ValueString())
	}

	if !plan.VipClientIpMode.IsNull() && !plan.VipClientIpMode.IsUnknown() {
		pool.SetVipClientIpMode(plan.VipClientIpMode.ValueString())
	}

	if !plan.Partition.IsNull() && !plan.Partition.IsUnknown() {
		pool.SetPartition(plan.Partition.ValueString())
	}

	if err := setUpdateConfig(ctx, pool, plan); err != nil {
		resp.Diagnostics.AddError("error updating load balancer pool", err.Error())

		return
	}

	loadBalancerID := currentState.LoadBalancerId.ValueInt64()
	id := currentState.Id.ValueInt64()

	updateReq := sdk.NewUpdateLoadBalancerPoolRequestWithDefaults()
	updateReq.SetLoadBalancerPool(*pool)

	_, httpResp, err := client.LoadBalancersAPI.
		UpdateLoadBalancerPool(ctx, loadBalancerID, id).
		UpdateLoadBalancerPoolRequest(*updateReq).
		Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"error updating load balancer pool",
			"load balancer pool "+plan.Name.ValueString()+" PUT failed: "+
				errfmt.ErrMsg(err, httpResp),
		)

		return
	}

	state, diags := getLoadBalancerPoolAsState(
		ctx, loadBalancerID, id, client,
	)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// Preserve config from plan since the API may not return it in the same shape.
	state.Config = plan.Config
	state.ConfigNsxt = plan.ConfigNsxt

	// Partition is write-only; preserve from plan.
	state.Partition = plan.Partition

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func setUpdateConfig(
	ctx context.Context,
	pool *sdk.UpdateLoadBalancerPoolRequestLoadBalancerPool,
	plan LoadBalancerPoolModel,
) error {
	if !plan.ConfigNsxt.IsNull() && !plan.ConfigNsxt.IsUnknown() {
		nsxConfig := sdk.NewNSXTLoadBalancerPoolConfigObject1WithDefaults()

		if !plan.ConfigNsxt.ActiveMonitorPaths.IsNull() && !plan.ConfigNsxt.ActiveMonitorPaths.IsUnknown() {
			nsxConfig.SetActiveMonitorPaths(plan.ConfigNsxt.ActiveMonitorPaths.ValueInt64())
		}

		if !plan.ConfigNsxt.PassiveMonitorPath.IsNull() && !plan.ConfigNsxt.PassiveMonitorPath.IsUnknown() {
			nsxConfig.SetPassiveMonitorPath(plan.ConfigNsxt.PassiveMonitorPath.ValueInt64())
		}

		if !plan.ConfigNsxt.SnatTranslationType.IsNull() && !plan.ConfigNsxt.SnatTranslationType.IsUnknown() {
			nsxConfig.SetSnatTranslationType(plan.ConfigNsxt.SnatTranslationType.ValueString())
		}

		if !plan.ConfigNsxt.SnatIpAddresses.IsNull() && !plan.ConfigNsxt.SnatIpAddresses.IsUnknown() {
			var addrs []string
			for _, elem := range plan.ConfigNsxt.SnatIpAddresses.Elements() {
				if sv, ok := elem.(types.String); ok && !sv.IsNull() {
					addrs = append(addrs, sv.ValueString())
				}
			}

			nsxConfig.SetSnatIpAddresses(addrs)
		}

		if !plan.ConfigNsxt.TcpMultiplexing.IsNull() && !plan.ConfigNsxt.TcpMultiplexing.IsUnknown() {
			nsxConfig.SetTcpMultiplexing(plan.ConfigNsxt.TcpMultiplexing.ValueBool())
		}

		if !plan.ConfigNsxt.TcpMultiplexingNumber.IsNull() && !plan.ConfigNsxt.TcpMultiplexingNumber.IsUnknown() {
			nsxConfig.SetTcpMultiplexingNumber(plan.ConfigNsxt.TcpMultiplexingNumber.ValueInt64())
		}

		// Build nested memberGroup if any member_group_* fields are set.
		hasMemberGroup := (!plan.ConfigNsxt.MemberGroupPath.IsNull() && !plan.ConfigNsxt.MemberGroupPath.IsUnknown()) ||
			(!plan.ConfigNsxt.MemberGroupIpRevisionFilter.IsNull() &&
				!plan.ConfigNsxt.MemberGroupIpRevisionFilter.IsUnknown()) ||
			(!plan.ConfigNsxt.MemberGroupMaxIpListSize.IsNull() &&
				!plan.ConfigNsxt.MemberGroupMaxIpListSize.IsUnknown()) ||
			(!plan.ConfigNsxt.MemberGroupPort.IsNull() && !plan.ConfigNsxt.MemberGroupPort.IsUnknown())

		if hasMemberGroup {
			mg := sdk.NewNSXTLoadBalancerPoolConfigObject1MemberGroupWithDefaults()

			if !plan.ConfigNsxt.MemberGroupPath.IsNull() && !plan.ConfigNsxt.MemberGroupPath.IsUnknown() {
				mg.SetPath(plan.ConfigNsxt.MemberGroupPath.ValueString())
			}

			if !plan.ConfigNsxt.MemberGroupIpRevisionFilter.IsNull() &&
				!plan.ConfigNsxt.MemberGroupIpRevisionFilter.IsUnknown() {
				mg.SetIpRevisionFilter(plan.ConfigNsxt.MemberGroupIpRevisionFilter.ValueString())
			}

			if !plan.ConfigNsxt.MemberGroupMaxIpListSize.IsNull() &&
				!plan.ConfigNsxt.MemberGroupMaxIpListSize.IsUnknown() {
				mg.SetMaxIpListSize(plan.ConfigNsxt.MemberGroupMaxIpListSize.ValueInt64())
			}

			if !plan.ConfigNsxt.MemberGroupPort.IsNull() && !plan.ConfigNsxt.MemberGroupPort.IsUnknown() {
				mg.SetPort(plan.ConfigNsxt.MemberGroupPort.ValueInt64())
			}

			nsxConfig.SetMemberGroup(*mg)
		}

		cfg := sdk.UpdateLoadBalancerPoolRequestLoadBalancerPoolConfig{
			NSXTLoadBalancerPoolConfigObject1: nsxConfig,
		}
		pool.SetConfig(cfg)

		return nil
	}

	if plan.Config.IsNull() || plan.Config.IsUnknown() {
		return nil
	}

	configValue := plan.Config.UnderlyingValue()
	configMap, err := convert.ValueToAny(ctx, configValue)
	if err != nil {
		return fmt.Errorf("failed to convert config: %w", err)
	}

	configDataMap, ok := configMap.(map[string]any)
	if !ok {
		return fmt.Errorf("config must be a valid object/map")
	}

	cfg := sdk.UpdateLoadBalancerPoolRequestLoadBalancerPoolConfig{
		MapmapOfStringAny: &configDataMap,
	}
	pool.SetConfig(cfg)

	return nil
}

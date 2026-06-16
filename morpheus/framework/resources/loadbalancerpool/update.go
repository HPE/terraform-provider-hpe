// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancerpool

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

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

	pool := &sdk.UpdateLoadBalancerPoolRequestLoadBalancerPool{}

	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		pool.Name = sdk.PtrString(plan.Name.ValueString())
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		pool.Description = sdk.PtrString(plan.Description.ValueString())
	}

	if !plan.VipBalance.IsNull() && !plan.VipBalance.IsUnknown() {
		pool.VipBalance = sdk.PtrString(plan.VipBalance.ValueString())
	}

	if !plan.MinActive.IsNull() && !plan.MinActive.IsUnknown() {
		pool.MinActive = sdk.PtrInt64(plan.MinActive.ValueInt64())
	}

	if !plan.Port.IsNull() && !plan.Port.IsUnknown() {
		pool.Port = sdk.PtrInt64(plan.Port.ValueInt64())
	}

	if !plan.VipSticky.IsNull() && !plan.VipSticky.IsUnknown() {
		pool.VipSticky = sdk.PtrString(plan.VipSticky.ValueString())
	}

	if !plan.VipClientIpMode.IsNull() && !plan.VipClientIpMode.IsUnknown() {
		pool.VipClientIpMode = sdk.PtrString(plan.VipClientIpMode.ValueString())
	}

	if !plan.Partition.IsNull() && !plan.Partition.IsUnknown() {
		pool.Partition = sdk.PtrString(plan.Partition.ValueString())
	}

	if err := setUpdateConfig(ctx, pool, plan); err != nil {
		resp.Diagnostics.AddError("error updating load balancer pool", err.Error())

		return
	}

	loadBalancerID := currentState.LoadBalancerId.ValueInt64()
	id := currentState.Id.ValueInt64()

	updateReq := &sdk.UpdateLoadBalancerPoolRequest{}
	updateReq.LoadBalancerPool = pool

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
		nsxConfig := &sdk.NSXTLoadBalancerPoolConfigObject1{}

		if !plan.ConfigNsxt.ActiveMonitorPaths.IsNull() && !plan.ConfigNsxt.ActiveMonitorPaths.IsUnknown() {
			nsxConfig.ActiveMonitorPaths.Set(sdk.PtrInt64(plan.ConfigNsxt.ActiveMonitorPaths.ValueInt64()))
		}

		if !plan.ConfigNsxt.PassiveMonitorPath.IsNull() && !plan.ConfigNsxt.PassiveMonitorPath.IsUnknown() {
			nsxConfig.PassiveMonitorPath.Set(sdk.PtrInt64(plan.ConfigNsxt.PassiveMonitorPath.ValueInt64()))
		}

		if !plan.ConfigNsxt.SnatTranslationType.IsNull() && !plan.ConfigNsxt.SnatTranslationType.IsUnknown() {
			nsxConfig.SnatTranslationType = sdk.PtrString(plan.ConfigNsxt.SnatTranslationType.ValueString())
		}

		if !plan.ConfigNsxt.SnatIpAddresses.IsNull() && !plan.ConfigNsxt.SnatIpAddresses.IsUnknown() {
			var addrs []string
			for _, elem := range plan.ConfigNsxt.SnatIpAddresses.Elements() {
				if sv, ok := elem.(types.String); ok && !sv.IsNull() {
					addrs = append(addrs, sv.ValueString())
				}
			}

			nsxConfig.SnatIpAddresses = addrs
		}

		if !plan.ConfigNsxt.TcpMultiplexing.IsNull() && !plan.ConfigNsxt.TcpMultiplexing.IsUnknown() {
			nsxConfig.TcpMultiplexing = sdk.PtrBool(plan.ConfigNsxt.TcpMultiplexing.ValueBool())
		}

		if !plan.ConfigNsxt.TcpMultiplexingNumber.IsNull() && !plan.ConfigNsxt.TcpMultiplexingNumber.IsUnknown() {
			nsxConfig.TcpMultiplexingNumber.Set(sdk.PtrInt64(plan.ConfigNsxt.TcpMultiplexingNumber.ValueInt64()))
		}

		// Build nested memberGroup if member_group is set.
		if !plan.ConfigNsxt.MemberGroup.IsNull() && !plan.ConfigNsxt.MemberGroup.IsUnknown() {
			var memberGroup MemberGroupValue

			diags := plan.ConfigNsxt.MemberGroup.As(ctx, &memberGroup, basetypes.ObjectAsOptions{})
			if diags.HasError() {
				return fmt.Errorf("failed to extract member_group: %s", diags.Errors()[0].Detail())
			}

			mg := &sdk.NSXTLoadBalancerPoolConfigObject1MemberGroup{}

			if !memberGroup.Path.IsNull() && !memberGroup.Path.IsUnknown() {
				mg.Path = sdk.PtrString(memberGroup.Path.ValueString())
			}

			if !memberGroup.IpRevisionFilter.IsNull() && !memberGroup.IpRevisionFilter.IsUnknown() {
				mg.IpRevisionFilter = sdk.PtrString(memberGroup.IpRevisionFilter.ValueString())
			}

			if !memberGroup.MaxIpListSize.IsNull() && !memberGroup.MaxIpListSize.IsUnknown() {
				mg.MaxIpListSize.Set(sdk.PtrInt64(memberGroup.MaxIpListSize.ValueInt64()))
			}

			if !memberGroup.Port.IsNull() && !memberGroup.Port.IsUnknown() {
				mg.Port.Set(sdk.PtrInt64(memberGroup.Port.ValueInt64()))
			}

			nsxConfig.MemberGroup = mg
		}

		cfg := sdk.UpdateLoadBalancerPoolRequestLoadBalancerPoolConfig{
			NSXTLoadBalancerPoolConfigObject1: nsxConfig,
		}
		pool.Config = &cfg

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
	pool.Config = &cfg

	return nil
}

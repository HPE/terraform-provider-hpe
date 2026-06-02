// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancer

import (
	"context"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan LoadBalancerModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"create load balancer resource",
			"failed to create client: "+err.Error(),
		)

		return
	}

	name := plan.Name.ValueString()

	createLB := &sdk.CreateLoadBalancerRequestLoadBalancer{
		Name: &name,
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		desc := plan.Description.ValueString()
		createLB.Description = &desc
	}

	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		vis := plan.Visibility.ValueString()
		createLB.Visibility = &vis
	}

	if !plan.GroupId.IsNull() && !plan.GroupId.IsUnknown() {
		groupID := plan.GroupId.ValueInt64()
		createLB.Site = &sdk.CreateLoadBalancerRequestLoadBalancerSite{
			Id: &groupID,
		}
	}

	if !plan.CloudId.IsNull() && !plan.CloudId.IsUnknown() {
		cloudID := plan.CloudId.ValueInt64()
		createLB.Zone = &sdk.CreateLoadBalancerRequestLoadBalancerZone{
			Id: &cloudID,
		}
	}

	if !plan.NetworkServerId.IsNull() && !plan.NetworkServerId.IsUnknown() {
		nsID := plan.NetworkServerId.ValueInt64()
		createLB.NetworkServerId = &nsID
	}

	if err := setCreateConfig(ctx, createLB, plan); err != nil {
		resp.Diagnostics.AddError("create load balancer resource", err.Error())

		return
	}

	if err := setCreateTenants(ctx, createLB, plan); err != nil {
		resp.Diagnostics.AddError("create load balancer resource", err.Error())

		return
	}

	if err := setCreatePermissions(ctx, createLB, plan); err != nil {
		resp.Diagnostics.AddError("create load balancer resource", err.Error())

		return
	}

	createReq := &sdk.CreateLoadBalancerRequest{
		LoadBalancer: createLB,
	}

	lb, hresp, err := client.LoadBalancersAPI.CreateLoadBalancer(ctx).
		CreateLoadBalancerRequest(*createReq).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"create load balancer resource",
			fmt.Sprintf("load balancer %s POST failed: %s",
				name, errfmt.ErrMsg(err, hresp)),
		)

		return
	}

	if lb.LoadBalancer.Id == nil {
		resp.Diagnostics.AddError(
			"create load balancer resource",
			"load balancer "+name+": id is nil",
		)

		return
	}

	id := *lb.LoadBalancer.Id
	plan.Id = types.Int64Value(id)

	taintResourceState := func(id int64) {
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "load_balancer",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})
	}

	state, err := getLoadBalancerAsState(ctx, id, client, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to read load balancer state",
			fmt.Sprintf("Load balancer %d was created but could not be read: %s", id, err),
		)
		taintResourceState(id)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError(
			"failed to set load balancer state",
			fmt.Sprintf("Load balancer %d was created but state could not be saved", id),
		)
		taintResourceState(id)

		return
	}
}

func setCreateConfig(
	ctx context.Context,
	createLB *sdk.CreateLoadBalancerRequestLoadBalancer,
	plan LoadBalancerModel,
) error {
	switch {
	case !plan.ConfigHaproxy.IsNull() && !plan.ConfigHaproxy.IsUnknown():
		typeCode := typeCodeHAProxy
		createLB.Type = &typeCode

		planID := plan.ConfigHaproxy.PlanId.ValueInt64()
		planObj := &sdk.HAProxyLoadBalancerConfigObjectPlan{
			Id: &planID,
		}

		poolID := plan.ConfigHaproxy.Pool.ValueString()
		poolObj := &sdk.HAProxyLoadBalancerConfigObjectPool{
			Id: &poolID,
		}

		haproxyConfig := &sdk.HAProxyLoadBalancerConfigObject{
			Plan: planObj,
			Pool: poolObj,
		}

		cfg := sdk.CreateLoadBalancerRequestLoadBalancerConfig{}
		cfg.HAProxyLoadBalancerConfigObject = haproxyConfig
		createLB.Config = &cfg

	case !plan.ConfigNsxt.IsNull() && !plan.ConfigNsxt.IsUnknown():
		if !plan.TypeCode.IsNull() && !plan.TypeCode.IsUnknown() {
			typeCode := plan.TypeCode.ValueString()
			createLB.Type = &typeCode
		} else {
			typeCode := typeCodeNSXT
			createLB.Type = &typeCode
		}

		nsxtConfig := &sdk.NSXTLoadBalancerConfigObject{}

		if !plan.ConfigNsxt.AdminState.IsNull() && !plan.ConfigNsxt.AdminState.IsUnknown() {
			adminState := plan.ConfigNsxt.AdminState.ValueBool()
			nsxtConfig.AdminState = &adminState
		}

		if !plan.ConfigNsxt.LogLevel.IsNull() && !plan.ConfigNsxt.LogLevel.IsUnknown() {
			logLevel := plan.ConfigNsxt.LogLevel.ValueString()
			nsxtConfig.Loglevel = &logLevel
		}

		if !plan.ConfigNsxt.Size.IsNull() && !plan.ConfigNsxt.Size.IsUnknown() {
			size := plan.ConfigNsxt.Size.ValueString()
			nsxtConfig.Size = &size
		}

		if !plan.ConfigNsxt.Tier1Gateway.IsNull() && !plan.ConfigNsxt.Tier1Gateway.IsUnknown() {
			tier1 := plan.ConfigNsxt.Tier1Gateway.ValueString()
			nsxtConfig.Tier1 = &tier1
		}

		cfg := sdk.CreateLoadBalancerRequestLoadBalancerConfig{}
		cfg.NSXTLoadBalancerConfigObject = nsxtConfig
		createLB.Config = &cfg

	case !plan.Config.IsNull() && !plan.Config.IsUnknown():
		typeCode := plan.TypeCode.ValueString()
		createLB.Type = &typeCode

		configValue := plan.Config.UnderlyingValue()
		configMap, err := convert.ValueToAny(ctx, configValue)
		if err != nil {
			return fmt.Errorf("failed to convert config: %w", err)
		}

		configDataMap, ok := configMap.(map[string]any)
		if !ok {
			return fmt.Errorf("config must be a valid object/map")
		}

		cfg := sdk.CreateLoadBalancerRequestLoadBalancerConfig{}
		cfg.MapmapOfStringAny = &configDataMap
		createLB.Config = &cfg
	}

	return nil
}

func setCreateTenants(
	ctx context.Context,
	createLB *sdk.CreateLoadBalancerRequestLoadBalancer,
	plan LoadBalancerModel,
) error {
	if plan.Tenants.IsNull() || plan.Tenants.IsUnknown() {
		return nil
	}

	var tenantObjs []TenantsValue
	if diags := plan.Tenants.ElementsAs(ctx, &tenantObjs, false); diags.HasError() {
		return fmt.Errorf("failed to parse tenants: %s", diags.Errors())
	}

	var tenants []sdk.CreateLoadBalancerRequestLoadBalancerTenantsInner
	for _, t := range tenantObjs {
		if !t.Id.IsNull() && !t.Id.IsUnknown() {
			tenantID := t.Id.ValueInt64()
			tenant := sdk.CreateLoadBalancerRequestLoadBalancerTenantsInner{
				Id: &tenantID,
			}
			tenants = append(tenants, tenant)
		}
	}

	if len(tenants) > 0 {
		createLB.Tenants = tenants
	}

	return nil
}

func setCreatePermissions(
	ctx context.Context,
	createLB *sdk.CreateLoadBalancerRequestLoadBalancer,
	plan LoadBalancerModel,
) error {
	if plan.Permissions.IsNull() || plan.Permissions.IsUnknown() {
		return nil
	}

	perms := &sdk.CreateLoadBalancerRequestLoadBalancerResourcePermissions{}

	if !plan.Permissions.All.IsNull() && !plan.Permissions.All.IsUnknown() {
		allVal := plan.Permissions.All.ValueBool()
		perms.All = &allVal
	}

	if !plan.Permissions.Groups.IsNull() && !plan.Permissions.Groups.IsUnknown() {
		falseVal := false
		perms.All = &falseVal
		var groupIDs []int64
		if diags := plan.Permissions.Groups.ElementsAs(ctx, &groupIDs, false); diags.HasError() {
			return fmt.Errorf("failed to parse permission groups: %s", diags.Errors())
		}

		perms.Sites = groupIDs
	}

	createLB.ResourcePermissions = perms

	return nil
}

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

	createLB := sdk.NewCreateLoadBalancerRequestLoadBalancerWithDefaults()
	createLB.SetName(name)

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		createLB.SetDescription(plan.Description.ValueString())
	}

	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		createLB.SetVisibility(plan.Visibility.ValueString())
	}

	if !plan.GroupId.IsNull() && !plan.GroupId.IsUnknown() {
		site := sdk.NewCreateLoadBalancerRequestLoadBalancerSite()
		site.SetId(plan.GroupId.ValueInt64())
		createLB.SetSite(*site)
	}

	if !plan.CloudId.IsNull() && !plan.CloudId.IsUnknown() {
		zone := sdk.NewCreateLoadBalancerRequestLoadBalancerZone()
		zone.SetId(plan.CloudId.ValueInt64())
		createLB.SetZone(*zone)
	}

	if !plan.NetworkServerId.IsNull() && !plan.NetworkServerId.IsUnknown() {
		createLB.SetNetworkServerId(plan.NetworkServerId.ValueInt64())
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

	createReq := sdk.NewCreateLoadBalancerRequest()
	createReq.SetLoadBalancer(*createLB)

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

	if lb.GetLoadBalancer().Id == nil {
		resp.Diagnostics.AddError(
			"create load balancer resource",
			"load balancer "+name+": id is nil",
		)

		return
	}

	id := *lb.GetLoadBalancer().Id
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

	switch {
	case !plan.ConfigHaproxy.IsNull() && !plan.ConfigHaproxy.IsUnknown():
		state.ConfigHaproxy = plan.ConfigHaproxy
	case !plan.ConfigNsxt.IsNull() && !plan.ConfigNsxt.IsUnknown():
		state.ConfigNsxt = plan.ConfigNsxt
	case !plan.Config.IsNull() && !plan.Config.IsUnknown():
		state.Config = plan.Config
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
		createLB.SetType(typeCodeHAProxy)

		haproxyConfig := sdk.NewHAProxyLoadBalancerConfigObject()

		planObj := sdk.NewHAProxyLoadBalancerConfigObjectPlan()
		planObj.SetId(plan.ConfigHaproxy.PlanId.ValueInt64())
		haproxyConfig.SetPlan(*planObj)

		poolObj := sdk.NewHAProxyLoadBalancerConfigObjectPool()
		poolObj.SetId(plan.ConfigHaproxy.Pool.ValueString())
		haproxyConfig.SetPool(*poolObj)

		cfg := sdk.CreateLoadBalancerRequestLoadBalancerConfig{}
		cfg.HAProxyLoadBalancerConfigObject = haproxyConfig
		createLB.SetConfig(cfg)

	case !plan.ConfigNsxt.IsNull() && !plan.ConfigNsxt.IsUnknown():
		if !plan.TypeCode.IsNull() && !plan.TypeCode.IsUnknown() {
			createLB.SetType(plan.TypeCode.ValueString())
		} else {
			createLB.SetType(typeCodeNSXT)
		}

		configDataMap := configNsxtToMap(plan.ConfigNsxt)

		cfg := sdk.CreateLoadBalancerRequestLoadBalancerConfig{}
		cfg.MapmapOfStringAny = &configDataMap
		createLB.SetConfig(cfg)

	case !plan.Config.IsNull() && !plan.Config.IsUnknown():
		createLB.SetType(plan.TypeCode.ValueString())

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
		createLB.SetConfig(cfg)
	}

	return nil
}

func configNsxtToMap(config ConfigNsxtValue) map[string]any {
	configMap := map[string]any{}

	if !config.AdminState.IsNull() && !config.AdminState.IsUnknown() {
		configMap["adminState"] = config.AdminState.ValueBool()
	}

	if !config.LogLevel.IsNull() && !config.LogLevel.IsUnknown() {
		configMap["loglevel"] = config.LogLevel.ValueString()
	}

	if !config.Size.IsNull() && !config.Size.IsUnknown() {
		configMap["size"] = config.Size.ValueString()
	}

	if !config.Tier1Gateway.IsNull() && !config.Tier1Gateway.IsUnknown() {
		configMap["tier1"] = config.Tier1Gateway.ValueString()
	}

	return configMap
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
			tenant := sdk.CreateLoadBalancerRequestLoadBalancerTenantsInner{}
			tenant.SetId(t.Id.ValueInt64())
			tenants = append(tenants, tenant)
		}
	}

	if len(tenants) > 0 {
		createLB.SetTenants(tenants)
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

	perms := sdk.NewCreateLoadBalancerRequestLoadBalancerResourcePermissions()

	if !plan.Permissions.All.IsNull() && !plan.Permissions.All.IsUnknown() {
		perms.SetAll(plan.Permissions.All.ValueBool())
	}

	if !plan.Permissions.Groups.IsNull() && !plan.Permissions.Groups.IsUnknown() {
		perms.SetAll(false)
		var groupIDs []int64
		if diags := plan.Permissions.Groups.ElementsAs(ctx, &groupIDs, false); diags.HasError() {
			return fmt.Errorf("failed to parse permission groups: %s", diags.Errors())
		}

		perms.SetSites(groupIDs)
	}

	createLB.SetResourcePermissions(*perms)

	return nil
}

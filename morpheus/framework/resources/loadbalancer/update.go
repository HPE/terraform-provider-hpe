// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancer

import (
	"context"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state LoadBalancerModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueInt64()

	updateLB := sdk.NewUpdateLoadBalancerRequestLoadBalancerWithDefaults()

	updateLB.SetName(plan.Name.ValueString())

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		updateLB.SetDescription(plan.Description.ValueString())
	}

	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		updateLB.SetVisibility(plan.Visibility.ValueString())
	}

	if err := setUpdateConfig(ctx, updateLB, plan); err != nil {
		resp.Diagnostics.AddError("update load balancer resource", err.Error())

		return
	}

	if err := setUpdateTenants(ctx, updateLB, plan); err != nil {
		resp.Diagnostics.AddError("update load balancer resource", err.Error())

		return
	}

	if err := setUpdatePermissions(ctx, updateLB, plan); err != nil {
		resp.Diagnostics.AddError("update load balancer resource", err.Error())

		return
	}

	updateReq := sdk.NewUpdateLoadBalancerRequest()
	updateReq.SetLoadBalancer(*updateLB)

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"update load balancer resource",
			"failed to create client: "+err.Error(),
		)

		return
	}

	_, hresp, err := client.LoadBalancersAPI.UpdateLoadBalancer(ctx, id).
		UpdateLoadBalancerRequest(*updateReq).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"update load balancer resource",
			fmt.Sprintf("load balancer %d UPDATE failed: %s",
				id, errfmt.ErrMsg(err, hresp)),
		)

		return
	}

	newState, err := getLoadBalancerAsState(ctx, id, client, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"update load balancer resource",
			fmt.Sprintf("load balancer %d: failed to read from api: %s", id, err),
		)

		return
	}

	switch {
	case !plan.ConfigHaproxy.IsNull() && !plan.ConfigHaproxy.IsUnknown():
		newState.ConfigHaproxy = plan.ConfigHaproxy
	case !plan.Config.IsNull() && !plan.Config.IsUnknown():
		newState.Config = plan.Config
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func setUpdateConfig(
	ctx context.Context,
	updateLB *sdk.UpdateLoadBalancerRequestLoadBalancer,
	plan LoadBalancerModel,
) error {
	// The update SDK's SetConfig accepts map[string]interface{} (not a typed union
	// like the create SDK), so we build the map directly for HAProxy config.
	switch {
	case !plan.ConfigHaproxy.IsNull() && !plan.ConfigHaproxy.IsUnknown():
		configMap := map[string]interface{}{
			"plan": map[string]interface{}{
				"id": plan.ConfigHaproxy.PlanId.ValueInt64(),
			},
			"pool": map[string]interface{}{
				"id": plan.ConfigHaproxy.Pool.ValueString(),
			},
		}
		updateLB.SetConfig(configMap)

	case !plan.Config.IsNull() && !plan.Config.IsUnknown():
		configValue := plan.Config.UnderlyingValue()
		configMap, err := convert.ValueToAny(ctx, configValue)
		if err != nil {
			return fmt.Errorf("failed to convert config: %w", err)
		}

		configDataMap, ok := configMap.(map[string]any)
		if !ok {
			return fmt.Errorf("config must be a valid object/map")
		}

		updateLB.SetConfig(configDataMap)
	}

	return nil
}

func setUpdateTenants(
	ctx context.Context,
	updateLB *sdk.UpdateLoadBalancerRequestLoadBalancer,
	plan LoadBalancerModel,
) error {
	if plan.Tenants.IsNull() || plan.Tenants.IsUnknown() {
		return nil
	}

	var tenantObjs []TenantsValue
	if diags := plan.Tenants.ElementsAs(ctx, &tenantObjs, false); diags.HasError() {
		return fmt.Errorf("failed to parse tenants: %s", diags.Errors())
	}

	var tenants []sdk.UpdateLoadBalancerRequestLoadBalancerTenantsInner
	for _, t := range tenantObjs {
		if !t.Id.IsNull() && !t.Id.IsUnknown() {
			tenant := sdk.UpdateLoadBalancerRequestLoadBalancerTenantsInner{}
			tenant.SetId(t.Id.ValueInt64())
			tenants = append(tenants, tenant)
		}
	}

	if len(tenants) > 0 {
		updateLB.SetTenants(tenants)
	}

	return nil
}

func setUpdatePermissions(
	ctx context.Context,
	updateLB *sdk.UpdateLoadBalancerRequestLoadBalancer,
	plan LoadBalancerModel,
) error {
	if plan.Permissions.IsNull() || plan.Permissions.IsUnknown() {
		return nil
	}

	perms := sdk.NewUpdateLoadBalancerRequestLoadBalancerResourcePermission()

	if !plan.Permissions.All.IsNull() && !plan.Permissions.All.IsUnknown() {
		perms.SetAll(plan.Permissions.All.ValueBool())
	} else if !plan.Permissions.Groups.IsNull() && !plan.Permissions.Groups.IsUnknown() {
		perms.SetAll(false)
	}

	if !plan.Permissions.Groups.IsNull() && !plan.Permissions.Groups.IsUnknown() {
		var groupIDs []int64
		if diags := plan.Permissions.Groups.ElementsAs(ctx, &groupIDs, false); diags.HasError() {
			return fmt.Errorf("failed to parse permission groups: %s", diags.Errors())
		}

		perms.SetSites(groupIDs)
	}

	updateLB.SetResourcePermission(*perms)

	return nil
}

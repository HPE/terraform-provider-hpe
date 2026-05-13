// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancervirtualserver

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
	var plan, currentState LoadBalancerVirtualServerModel

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
		resp.Diagnostics.AddError("update load balancer virtual server", "failed to create client: "+err.Error())

		return
	}

	lbID, err := loadBalancerIDFromInt64(plan.LoadBalancerId)
	if err != nil {
		resp.Diagnostics.AddError("update load balancer virtual server", err.Error())

		return
	}

	currentLBID, err := loadBalancerIDFromInt64(currentState.LoadBalancerId)
	if err != nil {
		resp.Diagnostics.AddError("update load balancer virtual server", err.Error())

		return
	}

	if lbID != currentLBID {
		resp.Diagnostics.AddError(
			"update load balancer virtual server",
			"changing load_balancer_id is not supported; destroy and recreate the resource instead",
		)

		return
	}

	id := currentState.Id.ValueInt64()

	instance := sdk.NewUpdateLoadBalancerVirtualServerRequestLoadBalancerInstanceWithDefaults()

	if !plan.VipName.IsNull() && !plan.VipName.IsUnknown() {
		instance.SetVipName(plan.VipName.ValueString())
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		instance.SetDescription(plan.Description.ValueString())
	}

	if !plan.VipAddress.IsNull() && !plan.VipAddress.IsUnknown() {
		instance.SetVipAddress(plan.VipAddress.ValueString())
	}

	if !plan.VipPort.IsNull() && !plan.VipPort.IsUnknown() {
		instance.SetVipPort(plan.VipPort.ValueInt64())
	}

	if !plan.VipProtocol.IsNull() && !plan.VipProtocol.IsUnknown() {
		instance.SetVipProtocol(plan.VipProtocol.ValueString())
	}

	if !plan.VipHostname.IsNull() && !plan.VipHostname.IsUnknown() {
		instance.SetVipHostname(plan.VipHostname.ValueString())
	}

	if !plan.SslCert.IsNull() && !plan.SslCert.IsUnknown() {
		instance.SetSslCert(plan.SslCert.ValueInt64())
	}

	if !plan.SslServerCert.IsNull() && !plan.SslServerCert.IsUnknown() {
		instance.SetSslServerCert(plan.SslServerCert.ValueInt64())
	}

	if err := setUpdateConfig(ctx, instance, plan); err != nil {
		resp.Diagnostics.AddError("update load balancer virtual server", err.Error())

		return
	}

	updateReq := sdk.NewUpdateLoadBalancerVirtualServerRequestWithDefaults()
	updateReq.SetLoadBalancerInstance(*instance)

	_, hresp, err := client.LoadBalancersAPI.
		UpdateLoadBalancerVirtualServer(ctx, lbID, id).
		UpdateLoadBalancerVirtualServerRequest(*updateReq).
		Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"error updating load balancer virtual server",
			fmt.Sprintf("load balancer %d virtual server %d PUT failed: %s",
				lbID, id, errfmt.ErrMsg(err, hresp)),
		)

		return
	}

	state, _, _, diags := getVirtualServerAsState(ctx, lbID, id, client)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// Preserve fields the API does not return on GET.
	state.SslServerCert = plan.SslServerCert
	state.Config = plan.Config
	state.ConfigNsxt = plan.ConfigNsxt

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func setUpdateConfig(
	ctx context.Context,
	instance *sdk.UpdateLoadBalancerVirtualServerRequestLoadBalancerInstance,
	plan LoadBalancerVirtualServerModel,
) error {
	if !plan.ConfigNsxt.IsNull() && !plan.ConfigNsxt.IsUnknown() {
		nsxConfig := sdk.NewNSXVirtualServerConfigObject1()
		if !plan.ConfigNsxt.ApplicationProfile.IsNull() && !plan.ConfigNsxt.ApplicationProfile.IsUnknown() {
			nsxConfig.SetApplicationProfile(plan.ConfigNsxt.ApplicationProfile.ValueString())
		}

		cfg := sdk.NSXVirtualServerConfigObject1AsUpdateLoadBalancerVirtualServerRequestLoadBalancerInstanceConfig(nsxConfig)
		instance.SetConfig(cfg)

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

	cfg := sdk.MapmapOfStringAnyAsUpdateLoadBalancerVirtualServerRequestLoadBalancerInstanceConfig(&configDataMap)
	instance.SetConfig(cfg)

	return nil
}

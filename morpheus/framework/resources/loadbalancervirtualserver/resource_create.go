// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancervirtualserver

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan LoadBalancerVirtualServerModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("create load balancer virtual server", "failed to create client: "+err.Error())

		return
	}

	lbID, err := loadBalancerIDFromInt64(plan.LoadBalancerId)
	if err != nil {
		resp.Diagnostics.AddError("create load balancer virtual server", err.Error())

		return
	}

	instance := sdk.NewCreateLoadBalancerVirtualServerRequestLoadBalancerInstanceWithDefaults()

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

	if !plan.VipType.IsNull() && !plan.VipType.IsUnknown() {
		instance.SetVipType(plan.VipType.ValueString())
	}

	if !plan.VipPool.IsNull() && !plan.VipPool.IsUnknown() {
		instance.SetVipPool(plan.VipPool.ValueInt64())
	}

	if !plan.SslCert.IsNull() && !plan.SslCert.IsUnknown() {
		instance.SetSslCert(plan.SslCert.ValueInt64())
	}

	if !plan.SslServerCert.IsNull() && !plan.SslServerCert.IsUnknown() {
		instance.SetSslServerCert(plan.SslServerCert.ValueInt64())
	}

	if err := setCreateConfig(ctx, instance, plan); err != nil {
		resp.Diagnostics.AddError("create load balancer virtual server", err.Error())

		return
	}

	createReq := sdk.NewCreateLoadBalancerVirtualServerRequestWithDefaults()
	createReq.SetLoadBalancerInstance(*instance)

	createResp, hresp, err := client.LoadBalancersAPI.
		CreateLoadBalancerVirtualServer(ctx, lbID).
		CreateLoadBalancerVirtualServerRequest(*createReq).
		Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"error creating load balancer virtual server",
			fmt.Sprintf("load balancer %d virtual server POST failed: %s",
				lbID, errfmt.ErrMsg(err, hresp)),
		)

		return
	}

	createdInstance := createResp.GetLoadBalancerInstance()
	if createdInstance.Id == nil {
		resp.Diagnostics.AddError(
			"error creating load balancer virtual server",
			"response did not contain an id",
		)

		return
	}

	state, _, _, diags := getVirtualServerAsState(ctx, lbID, *createdInstance.Id, client)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// Preserve fields the API does not return in the expected schema type.
	state.VipPool = plan.VipPool
	state.Config = plan.Config
	state.ConfigNsxt = plan.ConfigNsxt

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func setCreateConfig(
	ctx context.Context,
	instance *sdk.CreateLoadBalancerVirtualServerRequestLoadBalancerInstance,
	plan LoadBalancerVirtualServerModel,
) error {
	if !plan.ConfigNsxt.IsNull() && !plan.ConfigNsxt.IsUnknown() {
		nsxConfig := sdk.NewNSXVirtualServerConfigObject()

		if !plan.PoolId.IsNull() && !plan.PoolId.IsUnknown() {
			nsxConfig.SetPool(strconv.FormatInt(plan.PoolId.ValueInt64(), 10))
		}

		if !plan.ConfigNsxt.ApplicationProfile.IsNull() && !plan.ConfigNsxt.ApplicationProfile.IsUnknown() {
			nsxConfig.SetApplicationProfile(plan.ConfigNsxt.ApplicationProfile.ValueInt64())
		}

		if !plan.ConfigNsxt.Persistence.IsNull() && !plan.ConfigNsxt.Persistence.IsUnknown() {
			nsxConfig.SetPersistence(plan.ConfigNsxt.Persistence.ValueString())
		}

		if !plan.ConfigNsxt.PersistenceProfile.IsNull() && !plan.ConfigNsxt.PersistenceProfile.IsUnknown() {
			nsxConfig.SetPersistenceProfile(plan.ConfigNsxt.PersistenceProfile.ValueInt64())
		}

		if !plan.ConfigNsxt.SslClientProfile.IsNull() && !plan.ConfigNsxt.SslClientProfile.IsUnknown() {
			nsxConfig.SetSslClientProfile(plan.ConfigNsxt.SslClientProfile.ValueInt64())
		}

		if !plan.ConfigNsxt.SslServerProfile.IsNull() && !plan.ConfigNsxt.SslServerProfile.IsUnknown() {
			nsxConfig.SetSslServerProfile(plan.ConfigNsxt.SslServerProfile.ValueInt64())
		}

		cfg := sdk.CreateLoadBalancerVirtualServerRequestLoadBalancerInstanceConfig{
			NSXVirtualServerConfigObject: nsxConfig,
		}
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

	cfg := sdk.CreateLoadBalancerVirtualServerRequestLoadBalancerInstanceConfig{
		MapmapOfStringAny: &configDataMap,
	}
	instance.SetConfig(cfg)

	return nil
}

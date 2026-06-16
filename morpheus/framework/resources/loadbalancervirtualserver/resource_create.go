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
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
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

	instance := &sdk.CreateLoadBalancerVirtualServerRequestLoadBalancerInstance{}

	if !plan.VipName.IsNull() && !plan.VipName.IsUnknown() {
		instance.VipName = plan.VipName.ValueStringPointer()
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		instance.Description = plan.Description.ValueStringPointer()
	}

	if !plan.VipAddress.IsNull() && !plan.VipAddress.IsUnknown() {
		instance.VipAddress = plan.VipAddress.ValueStringPointer()
	}

	if !plan.VipPort.IsNull() && !plan.VipPort.IsUnknown() {
		instance.VipPort = plan.VipPort.ValueInt64Pointer()
	}

	if !plan.VipProtocol.IsNull() && !plan.VipProtocol.IsUnknown() {
		instance.VipProtocol = plan.VipProtocol.ValueStringPointer()
	}

	if !plan.VipHostname.IsNull() && !plan.VipHostname.IsUnknown() {
		instance.VipHostname = plan.VipHostname.ValueStringPointer()
	}

	if !plan.VipType.IsNull() && !plan.VipType.IsUnknown() {
		instance.VipType = plan.VipType.ValueStringPointer()
	}

	if !plan.VipPool.IsNull() && !plan.VipPool.IsUnknown() {
		poolVal := plan.VipPool.ValueInt64()
		instance.VipPool.Set(&poolVal)
	}

	if !plan.SslCert.IsNull() && !plan.SslCert.IsUnknown() {
		instance.SslCert = plan.SslCert.ValueInt64Pointer()
	}

	if !plan.SslServerCert.IsNull() && !plan.SslServerCert.IsUnknown() {
		instance.SslServerCert = plan.SslServerCert.ValueInt64Pointer()
	}

	if err := setCreateConfig(ctx, instance, plan); err != nil {
		resp.Diagnostics.AddError("create load balancer virtual server", err.Error())

		return
	}

	createReq := &sdk.CreateLoadBalancerVirtualServerRequest{}
	createReq.LoadBalancerInstance = instance

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

	createdInstance := createResp.LoadBalancerInstance
	if createdInstance == nil || createdInstance.Id == nil {
		resp.Diagnostics.AddError(
			"error creating load balancer virtual server",
			"response did not contain an id",
		)

		return
	}

	state, _, _, diags := getVirtualServerAsState(ctx, lbID, *createdInstance.Id, client)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "loadbalancervirtualserver",
			ResourceID:   *createdInstance.Id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	// Preserve fields the API does not return in the expected schema type.
	state.VipPool = plan.VipPool
	state.Config = plan.Config
	state.ConfigNsxt = plan.ConfigNsxt

	// Pool ID: the API sends pool inside config for NSX-T but may not return
	// a top-level pool object in the GET response. Preserve from plan if needed.
	if state.PoolId.IsNull() && !plan.PoolId.IsNull() {
		state.PoolId = plan.PoolId
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func setCreateConfig(
	ctx context.Context,
	instance *sdk.CreateLoadBalancerVirtualServerRequestLoadBalancerInstance,
	plan LoadBalancerVirtualServerModel,
) error {
	if !plan.ConfigNsxt.IsNull() && !plan.ConfigNsxt.IsUnknown() {
		nsxConfig := &sdk.NSXVirtualServerConfigObject{}

		if !plan.PoolId.IsNull() && !plan.PoolId.IsUnknown() {
			poolStr := strconv.FormatInt(plan.PoolId.ValueInt64(), 10)
			nsxConfig.Pool.Set(&poolStr)
		}

		if !plan.ConfigNsxt.ApplicationProfile.IsNull() && !plan.ConfigNsxt.ApplicationProfile.IsUnknown() {
			val := plan.ConfigNsxt.ApplicationProfile.ValueInt64()
			nsxConfig.ApplicationProfile.Set(&val)
		}

		if !plan.ConfigNsxt.Persistence.IsNull() && !plan.ConfigNsxt.Persistence.IsUnknown() {
			val := plan.ConfigNsxt.Persistence.ValueString()
			nsxConfig.Persistence.Set(&val)
		}

		if !plan.ConfigNsxt.PersistenceProfile.IsNull() && !plan.ConfigNsxt.PersistenceProfile.IsUnknown() {
			val := plan.ConfigNsxt.PersistenceProfile.ValueInt64()
			nsxConfig.PersistenceProfile.Set(&val)
		}

		if !plan.ConfigNsxt.SslClientProfile.IsNull() && !plan.ConfigNsxt.SslClientProfile.IsUnknown() {
			val := plan.ConfigNsxt.SslClientProfile.ValueInt64()
			nsxConfig.SslClientProfile.Set(&val)
		}

		if !plan.ConfigNsxt.SslServerProfile.IsNull() && !plan.ConfigNsxt.SslServerProfile.IsUnknown() {
			val := plan.ConfigNsxt.SslServerProfile.ValueInt64()
			nsxConfig.SslServerProfile.Set(&val)
		}

		cfg := sdk.CreateLoadBalancerVirtualServerRequestLoadBalancerInstanceConfig{
			NSXVirtualServerConfigObject: nsxConfig,
		}
		instance.Config = &cfg

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
	instance.Config = &cfg

	return nil
}

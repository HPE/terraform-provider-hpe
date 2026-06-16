// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package network

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
	var plan NetworkModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"create network resource",
			"failed to create client: "+err.Error(),
		)

		return
	}

	name := plan.Name.ValueString()

	groupId := plan.GroupId.ValueInt64()
	site := sdk.CreateNetworksRequestNetworkSite{
		Id: &groupId,
	}

	createNetwork := &sdk.CreateNetworksRequestNetwork{
		Name: name,
		Site: site,
		Zone: sdk.CreateNetworksRequestNetworkZone{
			Id: plan.CloudId.ValueInt64(),
		},
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		description := plan.Description.ValueString()
		createNetwork.Description = *sdk.NewNullableString(&description)
	}

	if !plan.DisplayName.IsNull() && !plan.DisplayName.IsUnknown() {
		createNetwork.DisplayName = plan.DisplayName.ValueStringPointer()
	}

	if !plan.Active.IsNull() && !plan.Active.IsUnknown() {
		createNetwork.Active = plan.Active.ValueBoolPointer()
	}

	if !plan.Cidr.IsNull() && !plan.Cidr.IsUnknown() {
		createNetwork.Cidr = plan.Cidr.ValueStringPointer()
	}

	if !plan.CidrIpv6.IsNull() && !plan.CidrIpv6.IsUnknown() {
		createNetwork.CidrIPv6 = plan.CidrIpv6.ValueStringPointer()
	}

	if !plan.Gateway.IsNull() && !plan.Gateway.IsUnknown() {
		createNetwork.Gateway = plan.Gateway.ValueStringPointer()
	}

	if !plan.GatewayIpv6.IsNull() && !plan.GatewayIpv6.IsUnknown() {
		gatewayIpv6 := plan.GatewayIpv6.ValueString()
		createNetwork.GatewayIPv6 = *sdk.NewNullableString(&gatewayIpv6)
	}

	if !plan.DnsPrimary.IsNull() && !plan.DnsPrimary.IsUnknown() {
		createNetwork.DnsPrimary = plan.DnsPrimary.ValueStringPointer()
	}

	if !plan.DnsSecondary.IsNull() && !plan.DnsSecondary.IsUnknown() {
		createNetwork.DnsSecondary = plan.DnsSecondary.ValueStringPointer()
	}

	if !plan.DnsPrimaryIpv6.IsNull() && !plan.DnsPrimaryIpv6.IsUnknown() {
		dnsPrimaryIpv6 := plan.DnsPrimaryIpv6.ValueString()
		createNetwork.DnsPrimaryIPv6 = *sdk.NewNullableString(&dnsPrimaryIpv6)
	}

	if !plan.DnsSecondaryIpv6.IsNull() &&
		!plan.DnsSecondaryIpv6.IsUnknown() {
		dnsSecondaryIpv6 := plan.DnsSecondaryIpv6.ValueString()
		createNetwork.DnsSecondaryIPv6 = *sdk.NewNullableString(&dnsSecondaryIpv6)
	}

	if !plan.DhcpServer.IsNull() && !plan.DhcpServer.IsUnknown() {
		createNetwork.DhcpServer = plan.DhcpServer.ValueBoolPointer()
	}

	if !plan.DhcpServerIpv6.IsNull() &&
		!plan.DhcpServerIpv6.IsUnknown() {
		createNetwork.DhcpServerIPv6 = plan.DhcpServerIpv6.ValueBoolPointer()
	}

	if !plan.AllowStaticOverride.IsNull() &&
		!plan.AllowStaticOverride.IsUnknown() {
		createNetwork.AllowStaticOverride = plan.AllowStaticOverride.ValueBoolPointer()
	}

	if !plan.AssignPublicIp.IsNull() &&
		!plan.AssignPublicIp.IsUnknown() {
		createNetwork.AssignPublicIp = plan.AssignPublicIp.ValueBoolPointer()
	}

	if !plan.ApplianceUrlProxyBypass.IsNull() &&
		!plan.ApplianceUrlProxyBypass.IsUnknown() {
		createNetwork.ApplianceUrlProxyBypass = plan.ApplianceUrlProxyBypass.ValueBoolPointer()
	}

	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		createNetwork.Visibility = plan.Visibility.ValueStringPointer()
	}

	if !plan.VlanId.IsNull() && !plan.VlanId.IsUnknown() {
		createNetwork.VlanId = plan.VlanId.ValueInt64Pointer()
	}

	if !plan.PoolId.IsNull() && !plan.PoolId.IsUnknown() {
		poolId := plan.PoolId.ValueInt64()
		createNetwork.Pool = *sdk.NewNullableInt64(&poolId)
	}

	if !plan.PoolIpv6Id.IsNull() && !plan.PoolIpv6Id.IsUnknown() {
		poolIpv6Id := plan.PoolIpv6Id.ValueInt64()
		createNetwork.PoolIPv6 = *sdk.NewNullableInt64(&poolIpv6Id)
	}

	if !plan.ZonePoolId.IsNull() && !plan.ZonePoolId.IsUnknown() {
		zonePoolId := plan.ZonePoolId.ValueInt64()
		createNetwork.ZonePool = &sdk.CreateNetworksRequestNetworkZonePool{
			Id: &zonePoolId,
		}
	}

	if !plan.Ipv4enabled.IsNull() && !plan.Ipv4enabled.IsUnknown() {
		createNetwork.Ipv4Enabled = plan.Ipv4enabled.ValueBoolPointer()
	}

	if !plan.Ipv6enabled.IsNull() && !plan.Ipv6enabled.IsUnknown() {
		createNetwork.Ipv6Enabled = plan.Ipv6enabled.ValueBoolPointer()
	}

	if !plan.NetmaskIpv6.IsNull() && !plan.NetmaskIpv6.IsUnknown() {
		netmaskIpv6 := plan.NetmaskIpv6.ValueString()
		createNetwork.NetmaskIPv6 = *sdk.NewNullableString(&netmaskIpv6)
	}

	if !plan.NoProxy.IsNull() && !plan.NoProxy.IsUnknown() {
		noProxy := plan.NoProxy.ValueString()
		createNetwork.NoProxy = *sdk.NewNullableString(&noProxy)
	}

	if !plan.SearchDomains.IsNull() && !plan.SearchDomains.IsUnknown() {
		createNetwork.SearchDomains = plan.SearchDomains.ValueStringPointer()
	}

	if !plan.SwitchId.IsNull() && !plan.SwitchId.IsUnknown() {
		createNetwork.SwitchId = plan.SwitchId.ValueStringPointer()
	}

	if !plan.TypeId.IsNull() && !plan.TypeId.IsUnknown() {
		createNetwork.Type = &sdk.CreateNetworksRequestNetworkType{
			Id: plan.TypeId.ValueInt64(),
		}
	}

	if !plan.NetworkDomainId.IsNull() &&
		!plan.NetworkDomainId.IsUnknown() {
		networkDomainId := plan.NetworkDomainId.ValueInt64()
		createNetwork.NetworkDomain = &sdk.CreateNetworksRequestNetworkNetworkDomain{
			Id: &networkDomainId,
		}
	}

	if !plan.NetworkProxyId.IsNull() &&
		!plan.NetworkProxyId.IsUnknown() {
		networkProxyId := plan.NetworkProxyId.ValueInt64()
		createNetwork.NetworkProxy = &sdk.CreateNetworksRequestNetworkNetworkProxy{
			Id: &networkProxyId,
		}
	}

	if !plan.Labels.IsNull() && !plan.Labels.IsUnknown() {
		labels, err := convert.SetToStrSlice(plan.Labels)
		if err != nil {
			resp.Diagnostics.AddError(
				"create network resource",
				"network "+name+": failed to parse labels: "+
					err.Error(),
			)

			return
		}
		createNetwork.Labels = labels
	}

	if !plan.Config.IsNull() && !plan.Config.IsUnknown() {
		configValue := plan.Config.UnderlyingValue()
		configMap, err := convert.ValueToAny(ctx, configValue)
		if err != nil {
			resp.Diagnostics.AddError(
				"create network resource",
				"network "+name+": failed to convert config: "+
					err.Error(),
			)

			return
		}

		configDataMap, ok := configMap.(map[string]any)
		if ok {
			networkConfig := sdk.CreateNetworksRequestNetworkConfig{}
			networkConfig.MapmapOfStringAny = &configDataMap
			createNetwork.Config = &networkConfig
		} else {
			resp.Diagnostics.AddError(
				"create network resource",
				"network "+name+": config must be a valid object/map",
			)

			return
		}
	}

	var tenants []sdk.CreateNetworksRequestNetworkTenantsInner
	if !plan.TenantIds.IsNull() && !plan.TenantIds.IsUnknown() {
		var tenantIDs []types.Int64
		diags := plan.TenantIds.ElementsAs(ctx, &tenantIDs, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		for _, idVal := range tenantIDs {
			if !idVal.IsNull() {
				tenantId := idVal.ValueInt64()
				tenant := sdk.
					CreateNetworksRequestNetworkTenantsInner{
					Id: &tenantId,
				}
				tenants = append(tenants, tenant)
			}
		}
		if len(tenants) > 0 {
			createNetwork.Tenants = tenants
		}
	}

	createNetworkReq := &sdk.CreateNetworksRequest{
		Network: createNetwork,
	}

	network, hresp, err := client.NetworksAPI.CreateNetworks(ctx).
		CreateNetworksRequest(*createNetworkReq).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"create network resource",
			fmt.Sprintf("network %s POST failed: %s",
				name, errfmt.ErrMsg(err, hresp)),
		)

		return
	}

	if network.Network == nil {
		resp.Diagnostics.AddError("API returned nil", "Network is nil in the response")

		return
	}

	if network.Network.Id == nil {
		resp.Diagnostics.AddError(
			"create network resource",
			"network "+name+": id is nil",
		)

		return
	}

	id := *network.Network.Id
	plan.Id = types.Int64Value(id)

	// Helper to taint the resource state on an error after the POST request
	taintResourceState := func(id int64) {
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "network",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})
	}

	state, pdiags := getNetworkAsState(ctx, id, client, plan)
	if pdiags.HasError() {
		resp.Diagnostics.Append(pdiags...)
		resp.Diagnostics.AddError(
			"failed to read network state",
			fmt.Sprintf("Network %d was created but could not be read", id),
		)
		taintResourceState(id)

		return
	}

	if !plan.Config.IsNull() && !plan.Config.IsUnknown() {
		state.Config = plan.Config
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError(
			"failed to set network state",
			fmt.Sprintf("Network %d was created but state could not be saved", id),
		)
		taintResourceState(id)

		return
	}
}
